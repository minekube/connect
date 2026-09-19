package ws

import (
	"context"
	"sync"
	"testing"

	"go.minekube.com/connect"
	"go.minekube.com/connect/internal/testutil"
)

// orderTestAddr is bound by the forced-order suite below. It is a different
// port than the one ws/suite_test.go uses so both can be part of the same run.
const orderTestAddr = ":8081"

// TestWatchRejectServerFirstObservation pins the interleaving that reds the
// build of an unrelated change: the server records its observation of the
// rejection first, the client records its own only afterwards
// (minekube/connect#170, run 35408785186, TestSuite/TestWatchReject).
//
// The suite observes one exchange from both sides: the server records its
// observation as soon as it read the rejection frame, while the client records
// "rejection sent" only after its write returned. Under runner load the server
// side wins that race. This test forces the interleaving deterministically
// instead of hoping for a busy scheduler: the client's Reject does not return
// before the server side recorded its observation.
func TestWatchRejectServerFirstObservation(t *testing.T) {
	serverObserved := make(chan struct{})
	s := &testutil.Suite{
		Endpoint: holdRejectEndpoint{
			Endpoint: ClientOptions{URL: "ws://" + orderTestAddr},
			observed: serverObserved,
		},
		StartWatchServer: func(ctx context.Context, ln connect.EndpointListener) error {
			return startServer(ctx, orderTestAddr, ServerOptions{}.EndpointHandler(
				&observeFirstRejectionListener{EndpointListener: ln, observed: serverObserved}))
		},
	}
	// Install the *testing.T context and drive the single method under test
	// instead of suite.Run: the suite's other method (TestTunnel) needs no
	// forced order and is covered by the transport's own suite test.
	s.SetT(t)
	s.TestWatchReject()
}

// holdRejectEndpoint wraps a connect.Endpoint and holds back the return of
// SessionProposal.Reject until the server has observed the rejection.
type holdRejectEndpoint struct {
	connect.Endpoint
	observed <-chan struct{}
}

func (e holdRejectEndpoint) Watch(ctx context.Context, propose connect.ReceiveProposal) error {
	return e.Endpoint.Watch(ctx, func(proposal connect.SessionProposal) error {
		return propose(holdRejectProposal{SessionProposal: proposal, observed: e.observed})
	})
}

type holdRejectProposal struct {
	connect.SessionProposal
	observed <-chan struct{}
}

func (p holdRejectProposal) Reject(ctx context.Context, reason *connect.RejectionReason) error {
	if err := p.SessionProposal.Reject(ctx, reason); err != nil {
		return err
	}
	// Wait until the server side recorded the rejection, so this side records
	// its own observation strictly after the server's.
	select {
	case <-p.observed:
	case <-ctx.Done():
	}
	return nil
}

// observeFirstRejectionListener forwards an accepted watch to the wrapped
// listener and closes observed once that listener returned, which is after the
// server recorded its observation of the rejection.
type observeFirstRejectionListener struct {
	connect.EndpointListener
	observed chan struct{}
	once     sync.Once
}

func (l *observeFirstRejectionListener) AcceptEndpoint(ctx context.Context, watch connect.EndpointWatch) error {
	rejections := make(chan *connect.SessionRejection)
	go func() {
		defer close(rejections)
		for rejection := range watch.Rejections() {
			select {
			case rejections <- rejection:
			case <-ctx.Done():
				return
			}
		}
	}()

	err := l.EndpointListener.AcceptEndpoint(ctx, forwardRejectionsWatch{EndpointWatch: watch, rejections: rejections})
	l.once.Do(func() { close(l.observed) })
	return err
}

// forwardRejectionsWatch serves the rejections of a wrapped watch through a
// channel of its own so the wrapper can observe them too.
type forwardRejectionsWatch struct {
	connect.EndpointWatch
	rejections <-chan *connect.SessionRejection
}

func (w forwardRejectionsWatch) Rejections() <-chan *connect.SessionRejection { return w.rejections }

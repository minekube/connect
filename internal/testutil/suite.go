package testutil

import (
	"context"
	"fmt"
	"io"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.minekube.com/connect"
)

type Suite struct {
	suite.Suite

	Endpoint          connect.Endpoint
	StartWatchServer  func(ctx context.Context, ln connect.EndpointListener) error
	StartTunnelServer func(ctx context.Context, ln connect.TunnelListener) error
}

func (suite *Suite) SetupTest() {
	suite.NotNil(suite.Endpoint)
	suite.NotNil(suite.StartWatchServer)
}

func (suite *Suite) TestWatchReject() {
	proposal := &connect.Session{Id: "abc"}
	rejection := status.New(codes.Aborted, "don't want this").Proto()

	var (
		clientProposal  = "c: got proposal " + proposal.GetId()
		clientSent      = "c: rejection sent"
		serverRejection = "s: got rejection " + rejection.String()
	)

	// The exchange is observed from both sides and only the wire synchronizes
	// them: the server records its observation of the rejection as soon as it
	// read the rejecting frame, while the client records "c: rejection sent"
	// only after its write returned. Both interleavings are therefore valid
	// outcomes of a correct exchange - depending on how the scheduler runs the
	// two sides - so those two observations are compared as an unordered pair
	// instead of pinning one of them (a loaded runner records the server
	// observation first, see minekube/connect#170).
	//
	// Every other observation of the exchange is causally ordered:
	//
	//	"s: got watcher" -> "c: got proposal" -> "s: got rejection"
	//
	// and the client records its two observations in between them.
	var (
		allEvents    = []string{"s: got watcher", clientProposal, clientSent, serverRejection}
		causalEvents = []string{"s: got watcher", clientProposal, serverRejection}
		clientEvents = []string{clientProposal, clientSent}
	)
	var seq sequence

	ctx, stop := context.WithTimeout(context.TODO(), time.Second*3)

	ln := acceptEndpoint(func(ctx context.Context, watch connect.EndpointWatch) error {
		seq.Add("s: got watcher")
		suite.NoError(watch.Propose(ctx, proposal))
		for rej := range watch.Rejections() {
			seq.Add("s: got rejection " + rej.GetReason().String())
			break
		}
		time.Sleep(time.Millisecond * 100) // let "c: rejection sent"

		return nil
	})

	go func() { suite.NoError(suite.StartWatchServer(ctx, ln)) }()

	time.Sleep(time.Millisecond * 100) // Wait for server
	go func() {
		err := suite.Endpoint.Watch(ctx, func(proposal connect.SessionProposal) error {
			seq.Add("c: got proposal " + proposal.Session().GetId())
			suite.NoError(proposal.Reject(ctx, rejection))
			seq.Add("c: rejection sent")
			return nil
		})
		suite.NotNil(err)
		suite.Contains(err.Error(), "closed serverside")
		stop()
	}()

	<-ctx.Done()
	suite.Assert().ErrorIs(ctx.Err(), context.Canceled)

	// All observations are recorded, in the order they were made.
	suite.Assert().ElementsMatch(allEvents, seq.Get(), "all observations")
	// The causally ordered observations are exactly in this order, no matter
	// which side recorded its rejection observation first.
	suite.Assert().Equal(causalEvents, seq.Excluding(clientSent), "causally ordered observations")
	// The client records its own observations in the order it made them.
	suite.Assert().Equal(clientEvents, seq.WithPrefix("c: "), "client observations")
}

func (suite *Suite) TestTunnel() {
	ctx, stop := context.WithTimeout(context.TODO(), time.Second*3)

	toClientMsg := []byte("hello client")
	toServerMsg := []byte("hello server")

	var expSeq = []string{
		"c: tunnel opened Tunnel(remote=unknown via connect, local=unknown via connect)",
		"s: got tunnel Tunnel(remote=unknown via connect, local=unknown via connect)",
		"s: read",
		"c: read",
		"c: no read",
	}
	var seq sequence

	ln := acceptTunnel(func(ctx context.Context, tunnel connect.Tunnel) error {
		time.Sleep(time.Millisecond * 100) // let "c: tunnel opened"
		seq.Add("s: got tunnel " + fmt.Sprint(tunnel))

		// client -> server
		b := make([]byte, 100)
		n, err := tunnel.Read(b)
		suite.NoError(err)
		suite.Equal(len(toServerMsg), n)
		suite.Equal(toServerMsg, b[:n])
		seq.Add("s: read")

		// client <- server
		n, err = tunnel.Write(toClientMsg)
		suite.NoError(err)
		suite.Equal(len(toClientMsg), n)

		// Close server side
		suite.NoError(tunnel.Close())
		return nil
	})

	go func() { suite.NoError(suite.StartTunnelServer(ctx, ln)) }()

	time.Sleep(time.Millisecond * 100) // Wait for server
	go func() {
		tunnel, err := suite.Endpoint.Tunnel(ctx)
		suite.NoError(err)
		seq.Add("c: tunnel opened " + fmt.Sprint(tunnel))

		// client -> server
		n, err := tunnel.Write(toServerMsg)
		suite.NoError(err)
		suite.Equal(len(toServerMsg), n)

		// client <- server
		b := make([]byte, 100)
		n, err = tunnel.Read(b)
		suite.NoError(err)
		suite.Equal(len(toClientMsg), n)
		suite.Equal(toClientMsg, b[:n])
		seq.Add("c: read")

		// Should be closed server side by now
		b = make([]byte, 100)
		for i := 0; i < 5; i++ {
			n, err = tunnel.Read(b)
			suite.Empty(n)
			if err == nil {
				continue // retry
			}
		}
		suite.ErrorIs(err, io.EOF)
		seq.Add("c: no read")

		_ = tunnel.Close()
		stop()
	}()

	<-ctx.Done()
	suite.Assert().ErrorIs(ctx.Err(), context.Canceled)
	suite.Assert().Equal(expSeq, seq.Get())
}

// sequence records the observations of a suite in the order they were made.
// Its observations are appended by concurrently running sides, so access is
// serialized.
type sequence struct {
	mu sync.Mutex
	v  []string
}

// Add records an observation.
func (s *sequence) Add(str string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.v = append(s.v, str)
}

// Get returns a snapshot of all recorded observations in record order.
func (s *sequence) Get() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return slices.Clone(s.v)
}

// Excluding returns the recorded observations without those equal to any of
// the given events, preserving record order. It drops observations whose order
// relative to the other side's observations depends on scheduling.
func (s *sequence) Excluding(events ...string) []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	kept := make([]string, 0, len(s.v))
	for _, observation := range s.v {
		if !slices.Contains(events, observation) {
			kept = append(kept, observation)
		}
	}
	return kept
}

// WithPrefix returns the recorded observations of the side with the given
// prefix, preserving record order.
func (s *sequence) WithPrefix(prefix string) []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	side := make([]string, 0, len(s.v))
	for _, observation := range s.v {
		if strings.HasPrefix(observation, prefix) {
			side = append(side, observation)
		}
	}
	return side
}

type acceptEndpoint func(ctx context.Context, watch connect.EndpointWatch) error

func (fn acceptEndpoint) AcceptEndpoint(ctx context.Context, watch connect.EndpointWatch) error {
	return fn(ctx, watch)
}

type acceptTunnel func(ctx context.Context, tunnel connect.Tunnel) error

func (fn acceptTunnel) AcceptTunnel(ctx context.Context, tunnel connect.Tunnel) error {
	return fn(ctx, tunnel)
}

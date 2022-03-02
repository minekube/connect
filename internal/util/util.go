package util

import (
	"context"
	"errors"

	"go.minekube.com/connect"
)

type SessionProposal struct {
	Proposal *connect.Session
	RejectFn func(ctx context.Context, r *connect.RejectionReason) error
}

func (p *SessionProposal) Session() *connect.Session {
	return p.Proposal
}

func (p *SessionProposal) Reject(ctx context.Context, r *connect.RejectionReason) error {
	return p.RejectFn(ctx, r)
}

type EndpointWatch struct {
	ProposeFn      func(session *connect.Session) error
	RejectionsChan chan *connect.SessionRejection

	Receive func() (*connect.WatchRequest, error)
}

func (w *EndpointWatch) Propose(session *connect.Session) error {
	if session == nil {
		return errors.New("session must not be nil")
	}
	if session.GetId() == "" {
		return errors.New("missing session id")
	}
	return w.ProposeFn(session)
}
func (w *EndpointWatch) Rejections() <-chan *connect.SessionRejection {
	return w.RejectionsChan
}

func (w *EndpointWatch) StartReceiveRejections(ctx context.Context) {
	defer close(w.RejectionsChan)
	for {
		req, err := w.Receive()
		if err != nil {
			return // drop error
		}
		if req.GetSessionRejection() == nil {
			continue // Unexpected
		}
		select {
		case w.RejectionsChan <- req.GetSessionRejection():
		case <-ctx.Done():
			return // stream closed, done
		}
	}
}

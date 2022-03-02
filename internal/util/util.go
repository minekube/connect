package util

import (
	"context"

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

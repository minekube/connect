package connect

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"io"
)

// SessionProposal specifies an incoming session proposal.
// Use the Session to create the connection tunnel or reject the session with an optional reason.
type SessionProposal interface {
	Session() *Session                    // The session that wants to connect.
	Reject(reason *RejectionReason) error // Rejects the session proposal with an optional reason.
}

// WatchOptions are the options for a call to Watch.
type WatchOptions struct {
	Cli      WatchServiceClient                   // The WatchServiceClient to use for watching.
	Callback func(proposal SessionProposal) error // The callback that is called when receiving session proposals.
}

// Watch registers the endpoint and watches for sessions being proposed
// by the WatchService and runs the WatchOptions.Callback with them.
// To stop Watch and return early cancel the passed context.
// If the callback function returns a non-nil non-EOF error Watch returns it.
func Watch(
	ctx context.Context,
	watchOpts WatchOptions,
	opts ...grpc.CallOption,
) error {
	// Validate options
	if watchOpts.Cli == nil {
		return errors.New("missing Cli in WatchOptions")
	}
	if watchOpts.Callback == nil {
		return errors.New("missing Callback in WatchOptions")
	}
	// Start watch
	watchCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	stream, err := watchOpts.Cli.Watch(watchCtx, opts...)
	if err != nil {
		return err
	}
	// Watch for session proposals
	var res *WatchResponse
	for {
		res, err = stream.Recv()
		if err != nil {
			break
		}
		if res.GetSession() == nil {
			continue
		}
		proposal := &sessionProposal{
			s:      res.GetSession(),
			stream: stream,
		}
		if err = watchOpts.Callback(proposal); err != nil {
			break
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return err
}

type sessionProposal struct {
	s      *Session
	stream WatchService_WatchClient
}

func (s *sessionProposal) Session() *Session { return s.s }
func (s *sessionProposal) Reject(reason *RejectionReason) error {
	return s.stream.Send(&WatchRequest{SessionRejection: &SessionRejection{
		Id:     s.s.GetId(),
		Reason: reason,
	}})
}

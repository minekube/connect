package connect

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"sync"
)

// WatchOptions are the options for a call to Watch.
type WatchOptions struct {
	Cli      WatchServiceClient                   // The WatchServiceClient to use for watching.
	Endpoint string                               // The endpoint name to register the watcher as.
	Callback func(proposal SessionProposal) error // The callback that is called when receiving session proposals.
}

// Watch registers the endpoint and watches for sessions being proposed
// by the WatchService and runs the WatchOptions.Callback with them.
// To stop Watch and return early cancel the passed context.
// If the callback function returns a non-nil error Watch returns it.
func Watch(
	ctx context.Context,
	watchOpts WatchOptions,
	opts ...grpc.CallOption,
) error {
	// Validate options
	if watchOpts.Endpoint == "" {
		return errors.New("missing Endpoint in WatchOptions")
	}
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
	// Register watcher as an endpoint.
	e := &Endpoint{Name: watchOpts.Endpoint}
	err = stream.Send(&WatchRequest{Endpoint: e})
	if err != nil {
		return err
	}
	var res *WatchResponse
	for {
		res, err = stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		if res.GetSession() == nil {
			continue
		}
		proposal := &sessionProposal{
			e:      e,
			s:      res.GetSession(),
			stream: stream,
		}
		if err = watchOpts.Callback(proposal); err != nil {
			return err
		}
	}
}

// SessionProposal specifies an incoming session proposal.
// Either prepare the session connection or reject it.
type SessionProposal interface {
	Session() *Session              // The session that wants to connect.
	Reject(reason *RejectionReason) // Rejects the session proposal with an optional reason.
	Endpoint() *Endpoint            // The endpoint details this session proposes to connect to.
}

type sessionProposal struct {
	once   sync.Once
	s      *Session
	e      *Endpoint
	stream WatchService_WatchClient
}

func (s *sessionProposal) Endpoint() *Endpoint { return s.e }
func (s *sessionProposal) Session() *Session   { return s.s }
func (s *sessionProposal) Reject(reason *RejectionReason) {
	s.once.Do(func() {
		if reason == nil {
			reason = status.New(codes.PermissionDenied, "").Proto()
		}
		_ = s.stream.Send(&WatchRequest{SessionRejection: &SessionRejection{
			Id:     s.s.GetId(),
			Reason: reason,
		}})
	})
}

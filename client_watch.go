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

// Watch registers the endpoint and watches for new sessions received from the WatchService and callbacks fn with it.
// The context must be canceled to stop watching and make Watch to return.
func Watch(
	ctx context.Context,
	cli WatchServiceClient,
	endpoint string,
	fn func(req SessionRequest) error,
	opts ...grpc.CallOption,
) error {
	stream, err := cli.Watch(ctx, opts...)
	if err != nil {
		return err
	}
	e := &Endpoint{Name: endpoint}
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
		req := &sessionRequest{
			e:      e,
			s:      res.GetSession(),
			stream: stream,
		}
		if err = fn(req); err != nil {
			return err
		}
	}
}

// SessionRequest specifies an incoming session intent.
// Either prepare for the session connection or deny it.
type SessionRequest interface {
	Session() *Session       // The session that wants to connect.
	Deny(reason *DenyReason) // Denies the session request with an optional reason.
	Endpoint() *Endpoint     // The endpoint details this session intents to connect to.
}

type sessionRequest struct {
	once   sync.Once
	s      *Session
	e      *Endpoint
	stream WatchService_WatchClient
}

func (s *sessionRequest) Endpoint() *Endpoint { return s.e }
func (s *sessionRequest) Session() *Session   { return s.s }
func (s *sessionRequest) Deny(reason *DenyReason) {
	s.once.Do(func() {
		if reason == nil {
			reason = status.New(codes.FailedPrecondition, "").Proto()
		}
		_ = s.stream.Send(&WatchRequest{
			DenySession: &DenySession{
				Id:     s.s.GetId(),
				Reason: reason,
			},
		})
	})
}

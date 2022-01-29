package connect

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"net"
	"time"
)

// Watcher represents a watching endpoint.
type Watcher interface {
	Context() context.Context             // The stream context.
	Endpoint() *Endpoint                  // The endpoint registered by the Watcher itself.
	Propose(*Session) error               // Propose proposes a session to take over by the watcher's endpoint.
	Rejections() <-chan *SessionRejection // Rejections returns a channel for listening to rejected session proposals.
}

// WatchService serves as a simple-to-use reference implementation for the WatchServiceServer.
type WatchService struct {
	// StartWatch is called when a new endpoint starts watching for sessions.
	// Watcher is the client that called the WatchService.
	// The watch lasts as long as the function is blocking and the watch is closed serverside on return.
	// Watchers context signals when the watcher disconnects.
	// The returned error is sent to the Watcher.
	StartWatch func(Watcher) error

	// ReceiveEndpointTimeout is the timeout to wait for the first stream
	// message that must be an Endpoint defining the watcher endpoint.
	//
	// Defaults to DefaultReceiveEndpointTimeout if <= 0.
	ReceiveEndpointTimeout time.Duration

	UnimplementedWatchServiceServer
}

// DefaultReceiveEndpointTimeout is the default timeout to wait for the first
// stream message that must be an Endpoint defining the watcher endpoint.
const DefaultReceiveEndpointTimeout = time.Second * 30

// Serve creates a gRPC server with the specified options and serves on the given listener.
// Signal the stop channel to stop the server and return.
// Remember that ctx.Done() can be passed as the stop argument.
func (s *WatchService) Serve(stop <-chan struct{}, ln net.Listener, opts ...grpc.ServerOption) error {
	if err := s.Valid(); err != nil {
		return err
	}
	svr := grpc.NewServer(opts...)
	s.Register(svr)
	go func() { <-stop; svr.Stop() }()
	return svr.Serve(ln)
}

// Register registers the WatchService with a grpc.ServiceRegistrar such as a grpc.Server.
func (s *WatchService) Register(r grpc.ServiceRegistrar) { RegisterWatchServiceServer(r, s) }

// Valid validates the WatchService struct fields.
// This is already called internally by Serve.
func (s *WatchService) Valid() error {
	if s == nil {
		return errors.New("nil pointer dereference")
	}
	if s.StartWatch == nil {
		return errors.New("missing StartWatch callback")
	}
	if s.ReceiveEndpointTimeout <= 0 {
		s.ReceiveEndpointTimeout = DefaultReceiveEndpointTimeout
	}
	return nil
}

// Watch implements WatchServiceServer.
// See the proto definition for more documentation.
func (s *WatchService) Watch(stream WatchService_WatchServer) error {
	done := make(chan error, 1)
	timeoutRecv := time.AfterFunc(s.ReceiveEndpointTimeout, func() {
		// Did not receive endpoint in time.
		// Returning the rpc cancels the fn below.
		done <- nil
	})
	fn := func() error {
		// Receive endpoint
		req, err := stream.Recv()
		timeoutRecv.Stop()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		// Validate request
		if req.GetEndpoint() == nil {
			return status.Error(codes.InvalidArgument, "first stream message must set the endpoint")
		}
		endpoint := req.GetEndpoint().GetName()
		if endpoint == "" {
			return status.Error(codes.InvalidArgument, "missing endpoint name")
		}
		// Prepare watcher
		w := &watcher{
			e:          req.GetEndpoint(),
			stream:     stream,
			rejections: make(chan *SessionRejection),
		}
		go w.receiveRejections()
		// Start blocking watch callback
		return s.StartWatch(w)
	}
	go func() { done <- fn() }()
	return <-done
}

var _ WatchServiceServer = (*WatchService)(nil)

type watcher struct {
	e          *Endpoint
	stream     WatchService_WatchServer
	rejections chan *SessionRejection
}

var _ Watcher = (*watcher)(nil)

func (w *watcher) Context() context.Context { return w.stream.Context() }
func (w *watcher) Endpoint() *Endpoint      { return w.e }
func (w *watcher) Propose(session *Session) error {
	if session == nil {
		return errors.New("session must not be nil")
	}
	if session.GetId() == "" {
		return errors.New("missing session id")
	}
	return w.stream.Send(&WatchResponse{Session: session})
}
func (w *watcher) Rejections() <-chan *SessionRejection { return w.rejections }
func (w *watcher) receiveRejections() {
	defer close(w.rejections)
	for {
		req, err := w.stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return // stream closed, done
			}
			return // drop error
		}
		if req.GetSessionRejection() == nil {
			// Unexpected
			continue
		}
		select {
		case w.rejections <- req.GetSessionRejection():
		case <-w.Context().Done():
			return // stream closed, done
		}
	}
}

package connect

import (
	"context"
	"errors"
	"net"

	"google.golang.org/grpc"
)

// Watcher represents a watching endpoint.
type Watcher interface {
	Context() context.Context             // The stream context that also carries metadata.
	Propose(*Session) error               // Propose proposes a session to the Watcher.
	Rejections() <-chan *SessionRejection // Rejections returns a channel for listening to rejected session proposals.
}

// WatchService serves as a reference implementation for the WatchServiceServer.
type WatchService struct {
	// StartWatch is called when a new endpoint starts watching for sessions.
	// Watcher is the client that called the WatchService.
	// The watch lasts as long as the function is blocking and the watch is closed serverside on return.
	// Watchers context signals when the watcher disconnects.
	// The returned error is sent to the client that called Watch.
	StartWatch func(Watcher) error

	UnimplementedWatchServiceServer
}

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
	return nil
}

// Watch implements WatchServiceServer.
// See the proto definition for more documentation.
func (s *WatchService) Watch(stream WatchService_WatchServer) error {
	// Prepare watcher
	w := &watcher{
		stream:     stream,
		rejections: make(chan *SessionRejection),
	}
	go w.startRecvRejections()
	// Start blocking watch callback
	return s.StartWatch(w)
}

var _ WatchServiceServer = (*WatchService)(nil)

type watcher struct {
	stream     WatchService_WatchServer
	rejections chan *SessionRejection
}

var _ Watcher = (*watcher)(nil)

func (w *watcher) Propose(session *Session) error {
	if session == nil {
		return errors.New("session must not be nil")
	}
	if session.GetId() == "" {
		return errors.New("missing session id")
	}
	return w.stream.Send(&WatchResponse{Session: session})
}
func (w *watcher) Context() context.Context             { return w.stream.Context() }
func (w *watcher) Rejections() <-chan *SessionRejection { return w.rejections }
func (w *watcher) startRecvRejections() {
	defer close(w.rejections)
	for {
		req, err := w.stream.Recv()
		if err != nil {
			return // drop error
		}
		if req.GetSessionRejection() == nil {
			continue // Unexpected
		}
		select {
		case w.rejections <- req.GetSessionRejection():
		case <-w.stream.Context().Done():
			return // stream closed, done
		}
	}
}

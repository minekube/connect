package connect

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"io"
	"net"
	"sync"
	"time"
)

// InboundTunnel represents an inbound tunnel.
type InboundTunnel interface {
	Context() context.Context // The stream context.
	SessionID() string        // The session id of the tunnel.
	Conn() TunnelConn         // The tunnel connection.
	io.Closer                 // Closes the tunnel.
}

type TunnelService struct {
	// AcceptTunnel is called when a new tunnel is inbound.
	// Unlike WatchService this function can return immediately
	// without the tunnel being closed. To close the tunnel
	// use InboundTunnel.Close() or InboundTunnel.Conn().Close().
	// The tunnel is also closed when the context is canceled.
	// If AcceptTunnel returns an error the tunnel is closed.
	AcceptTunnel func(InboundTunnel) error

	// ReceiveSessionTimeout is the timeout to wait for the
	// first stream message that must specify the session id.
	//
	// Defaults to DefaultReceiveSessionTimeout if <= 0.
	ReceiveSessionTimeout time.Duration

	// LocalAddr is the LocalAddr for InboundTunnel.
	// If not set becomes the net.Listener's addr passed to Serve.
	LocalAddr net.Addr

	UnimplementedTunnelServiceServer
}

// DefaultReceiveSessionTimeout is the default timeout to wait for the
// first stream message that must specify the session id.
const DefaultReceiveSessionTimeout = time.Second * 5

// Serve creates a gRPC server with the specified options and serves on the given listener.
// Signal the stop channel to stop the server and return.
// Remember that ctx.Done() can be passed as the stop argument.
func (s *TunnelService) Serve(stop <-chan struct{}, ln net.Listener, opts ...grpc.ServerOption) error {
	if s.LocalAddr == nil {
		s.LocalAddr = ln.Addr()
	}
	if err := s.Valid(); err != nil {
		return err
	}
	svr := grpc.NewServer(opts...)
	s.Register(svr)
	go func() { <-stop; svr.Stop() }()
	return svr.Serve(ln)
}

// Register registers the WatchService with a grpc.ServiceRegistrar such as a grpc.Server.
func (s *TunnelService) Register(r grpc.ServiceRegistrar) { RegisterTunnelServiceServer(r, s) }

// Valid validates the WatchService struct fields.
// This is already called internally by Serve.
func (s *TunnelService) Valid() error {
	if s == nil {
		return errors.New("nil pointer dereference")
	}
	if s.AcceptTunnel == nil {
		return errors.New("missing AcceptTunnel callback")
	}
	if s.LocalAddr == nil {
		return errors.New("missing LocalAddr")
	}
	if s.ReceiveSessionTimeout <= 0 {
		s.ReceiveSessionTimeout = DefaultReceiveSessionTimeout
	}
	return nil
}

// Tunnel implements TunnelServiceServer.
// See the proto definition for more documentation.
func (s *TunnelService) Tunnel(stream TunnelService_TunnelServer) error {
	done := make(chan error, 1)
	timeoutRecv := time.AfterFunc(s.ReceiveSessionTimeout, func() {
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
		if req.GetSessionId() == "" {
			return status.Error(codes.InvalidArgument, "first stream message must specify the session id")
		}

		// Prepare inbound tunnel
		closeTunnel := make(chan struct{})
		var closeOnce, connOnce sync.Once
		var conn TunnelConn
		closeFn := func() error { closeOnce.Do(func() { close(closeTunnel) }); return nil }
		tunnel := &inboundTunnel{
			sessionID: func() string { return req.GetSessionId() },
			ctx:       func() context.Context { return stream.Context() },
			close:     closeFn,
			conn: func() TunnelConn {
				// initialize tunnel read-/writer once we need it
				connOnce.Do(func() {
					p, _ := peer.FromContext(stream.Context())
					r, w := tunnelServerRW(stream.Context(), stream)
					conn = newTunnelConn(s.LocalAddr, p.Addr, r, w, closeFn)
				})
				return conn
			},
		}

		// Call inbound tunnel callback
		if err = s.AcceptTunnel(tunnel); err != nil {
			return err
		}

		// Block until tunnel closing
		select {
		case <-stream.Context().Done():
			return stream.Context().Err()
		case <-closeTunnel:
			return nil
		}
	}
	go func() { done <- fn() }()
	return <-done
}

var _ TunnelServiceServer = (*TunnelService)(nil)

type inboundTunnel struct {
	sessionID func() string
	close     func() error
	conn      func() TunnelConn
	ctx       func() context.Context
}

var _ InboundTunnel = (*inboundTunnel)(nil)

func (i *inboundTunnel) Close() error             { return i.close() }
func (i *inboundTunnel) Context() context.Context { return i.ctx() }
func (i *inboundTunnel) SessionID() string        { return i.sessionID() }
func (i *inboundTunnel) Conn() TunnelConn         { return i.conn() }

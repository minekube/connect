package v0

import (
	"context"
	"errors"
	"io"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// InboundTunnel represents an inbound tunnel.
type InboundTunnel interface {
	Context() context.Context // The stream context.
	Conn() TunnelConn         // The tunnel connection.
	io.Closer                 // Closes the tunnel. Same as Conn().Close()
}

// TunnelService serves as a simple-to-use reference implementation for the TunnelServiceServer.
type TunnelService struct {
	// AcceptTunnel is called when a new tunnel is inbound.
	// Unlike WatchService this function can return immediately
	// without the tunnel being closed. To close the tunnel
	// use InboundTunnel.Close() or InboundTunnel.Conn().Close().
	// The tunnel is also closed when the context is canceled.
	// If AcceptTunnel returns an error the tunnel is closed.
	AcceptTunnel func(InboundTunnel) error

	// LocalAddr is the LocalAddr for InboundTunnel.
	// If not set becomes the net.Listener's addr passed to Serve.
	LocalAddr net.Addr

	UnimplementedTunnelServiceServer
}

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
	return nil
}

// Tunnel implements TunnelServiceServer.
// See the proto definition for more documentation.
func (s *TunnelService) Tunnel(stream TunnelService_TunnelServer) error {
	// Get address of client
	var remoteAddr net.Addr
	if p, ok := peer.FromContext(stream.Context()); !ok {
		return status.Error(codes.Internal, "could not resolve address of client")
	} else {
		remoteAddr = p.Addr
	}
	// Create inbound tunnel from stream
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()
	var initConnOnce sync.Once
	var conn TunnelConn
	closeFn := func() error { cancel(); return nil }
	tunnel := &inboundTunnel{
		ctx:   func() context.Context { return ctx },
		close: closeFn,
		conn: func() TunnelConn {
			// initialize tunnel conn once we need it
			initConnOnce.Do(func() {
				r, w := tunnelServerRW(ctx, stream)
				conn = newTunnelConn(s.LocalAddr, remoteAddr, r, w, closeFn)
			})
			return conn
		},
	}
	// Call inbound tunnel callback
	if err := s.AcceptTunnel(tunnel); err != nil {
		return err
	}
	// Block until tunnel closed
	<-ctx.Done()
	return nil
}

var _ TunnelServiceServer = (*TunnelService)(nil)

type inboundTunnel struct {
	close func() error
	conn  func() TunnelConn
	ctx   func() context.Context
}

var _ InboundTunnel = (*inboundTunnel)(nil)

func (i *inboundTunnel) Close() error             { return i.close() }
func (i *inboundTunnel) Context() context.Context { return i.ctx() }
func (i *inboundTunnel) Conn() TunnelConn         { return i.conn() }

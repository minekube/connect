package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TunnelOptions are the options for a call to Tunnel.
type TunnelOptions struct {
	TunnelClient TunnelServiceClient // The TunnelServiceClient to use for tunneling.
	// The local address as reflected by TunnelConn.LocalAddr().
	LocalAddr string // It is recommended to use the endpoint name.
	// The remote address as reflected by TunnelConn.RemoteAddr().
	RemoteAddr string // It is recommended to use the player address.
}

// Tunnel establishes a tunnel connection for a new session by dialing the
// TunnelService and returning the TunnelConn implementing net.Conn.
// The tunnel connection can be closed either by cancelling the passed
// context or by calling TunnelConn.Close().
func Tunnel(
	ctx context.Context,
	tunnelOpts TunnelOptions,
	callOpts ...grpc.CallOption,
) (TunnelConn, error) {
	// Validate options
	if tunnelOpts.TunnelClient == nil {
		return nil, errors.New("missing TunnelClient in TunnelOptions")
	}
	if tunnelOpts.LocalAddr == "" {
		return nil, errors.New("missing LocalAddr in TunnelOptions")
	}
	if tunnelOpts.RemoteAddr == "" {
		return nil, errors.New("missing RemoteAddr in TunnelOptions")
	}
	// Make tunnel connection
	tunnelCtx, tunnelCancel := context.WithCancel(ctx)
	tunnelStream, err := tunnelOpts.TunnelClient.Tunnel(tunnelCtx, callOpts...)
	if err != nil {
		tunnelCancel()
		return nil, status.Errorf(statusErr(err).Code(), "%s: %v", "could not create tunnel", err)
	}
	// Return tunnel connection ready to use
	r, w := tunnelClientRW(tunnelCtx, tunnelStream)
	return newTunnelConn(
		connectAddr(tunnelOpts.LocalAddr), connectAddr(tunnelOpts.RemoteAddr),
		r, w, func() error { tunnelCancel(); return nil },
	), nil
}

func statusErr(err error) *status.Status {
	s := status.FromContextError(err)
	if s.Code() == codes.Unknown {
		s = status.Convert(err)
	}
	return s
}

type connectAddr string

func (a connectAddr) String() string  { return string(a) }
func (a connectAddr) Network() string { return "connect" }

package connect

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TunnelOptions are the options for a call to Tunnel.
type TunnelOptions struct {
	TunnelCli TunnelServiceClient // The TunnelServiceClient to use for tunneling.
	SessionID string              // The session id to tunnel the connection of.
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
	if tunnelOpts.SessionID == "" {
		return nil, errors.New("missing SessionID in TunnelOptions")
	}
	if tunnelOpts.TunnelCli == nil {
		return nil, errors.New("missing TunnelCli in TunnelOptions")
	}
	if tunnelOpts.LocalAddr == "" {
		return nil, errors.New("missing LocalAddr in TunnelOptions")
	}
	if tunnelOpts.RemoteAddr == "" {
		return nil, errors.New("missing RemoteAddr in TunnelOptions")
	}

	wrap := func(err error, msg string) error {
		return status.Errorf(statusErr(err).Code(), "%s: %v", msg, err)
	}
	// Make tunnel connection
	tunnelCtx, tunnelCancel := context.WithCancel(ctx)
	tunnelStream, err := tunnelOpts.TunnelCli.Tunnel(tunnelCtx, callOpts...)
	if err != nil {
		tunnelCancel()
		return nil, wrap(err, "could not create tunnel")
	}
	// Send session id before sending and receiving data
	err = tunnelStream.Send(&TunnelRequest{SessionId: tunnelOpts.SessionID})
	if err != nil {
		tunnelCancel()
		return nil, wrap(err, "could not send session id for new tunnel")
	}
	// Return tunnel connection ready to use
	r, w := tunnelClientRW(tunnelCtx, tunnelStream)
	return newTunnelConn(
		tcpAddr(tunnelOpts.LocalAddr), tcpAddr(tunnelOpts.RemoteAddr),
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

type tcpAddr string

func (a tcpAddr) String() string  { return string(a) }
func (a tcpAddr) Network() string { return "tcp" }

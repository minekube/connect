package connect

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

// DefaultTunnelDialTimeout is the default dial timeout for connecting to a TunnelService in Tunnel.
var DefaultTunnelDialTimeout = time.Second * 10

// Tunnel establishes a tunnel connection for a new session by dialing the
// TunnelService and returning the TunnelConn implementing net.Conn.
// The context must be canceled to close the tunnel connection gracefully.
func Tunnel(
	ctx context.Context,
	req SessionRequest,
	dialOpts []grpc.DialOption,
	callOpts []grpc.CallOption,
) (TunnelConn, error) {
	if req.Session().GetTunnelServiceAddr() == "" {
		return nil, errors.New("missing tunnel service address in session request")
	}

	wrap := func(err error, msg string) error {
		stErr := statusErr(err)
		st, _ := status.New(stErr.Code(), msg).
			WithDetails(&errdetails.ErrorInfo{
				Reason: stErr.String(),
				Domain: req.Session().GetTunnelServiceAddr(),
				Metadata: map[string]string{
					"sessionId":  req.Session().GetId(),
					"playerId":   req.Session().GetPlayer().GetProfile().GetId(),
					"playerName": req.Session().GetPlayer().GetProfile().GetName(),
					"playerAddr": req.Session().GetPlayer().GetAddr(),
				},
			})
		req.Deny(st.Proto())
		return fmt.Errorf("%s: %w", msg, err)
	}

	// Connect to tunnel service
	dialCtx, dialCancel := context.WithTimeout(ctx, DefaultTunnelDialTimeout)
	svcConn, err := grpc.DialContext(dialCtx, req.Session().GetTunnelServiceAddr(), dialOpts...)
	dialCancel()
	if err != nil {
		return nil, wrap(err, "could not dial tunnel service")
	}

	// Make tunnel connection
	tunnelCtx, tunnelCancel := context.WithCancel(ctx)
	tunnelClose := func() error { tunnelCancel(); return svcConn.Close() }
	tunnelStream, err := NewTunnelServiceClient(svcConn).Tunnel(tunnelCtx, callOpts...)
	if err != nil {
		_ = tunnelClose()
		return nil, wrap(err, "could not create tunnel")
	}

	// Send session id before before sending and receiving data
	err = tunnelStream.Send(&TunnelRequest{SessionId: req.Session().GetId()})
	if err != nil {
		_ = tunnelClose()
		return nil, wrap(err, "could not send session id for new tunnel")
	}

	// Return connection interface
	r, w := tunnelClientRW(tunnelCtx, tunnelStream)
	return newTunnelConn(
		req.Session(),
		tcpAddr(req.Endpoint().GetName()),
		tcpAddr(req.Session().GetPlayer().GetAddr()),
		r, w,
		tunnelClose,
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

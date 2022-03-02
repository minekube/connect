package v0

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"nhooyr.io/websocket"
)

type TunnelWebsocketOptions struct {
	// The local address as reflected by TunnelConn.LocalAddr().
	LocalAddr string // It is recommended to use the endpoint name.
	// The remote address as reflected by TunnelConn.RemoteAddr().
	RemoteAddr string // It is recommended to use the player address.

	HTTPHeader http.Header // The HTTP headers included in the handshake request.
	// HTTPClient is used for the connection.
	// Its Transport must return writable bodies for WebSocket handshakes.
	HTTPClient      *http.Client
	HandshakeResult func(res *http.Response) error
}

func TunnelWebsocket(ctx context.Context, url string, tunnelOpts TunnelWebsocketOptions) (TunnelConn, error) {
	// Validate options
	if tunnelOpts.LocalAddr == "" {
		return nil, errors.New("missing LocalAddr in TunnelOptions")
	}
	if tunnelOpts.RemoteAddr == "" {
		return nil, errors.New("missing RemoteAddr in TunnelOptions")
	}
	// Dial service // TODO move outside for specifying dial ctx timeout
	conn, resp, err := websocket.Dial(ctx, url, &websocket.DialOptions{
		HTTPClient:   tunnelOpts.HTTPClient,
		HTTPHeader:   tunnelOpts.HTTPHeader,
		Subprotocols: wsTunnelSvcOpts.accept.Subprotocols,
	})
	if err != nil {
		return nil, fmt.Errorf("error dialing %q: %w", url, err)
	}
	if tunnelOpts.HandshakeResult != nil {
		if err = tunnelOpts.HandshakeResult(resp); err != nil {
			return nil, err
		}
	}
	// Make tunnel connection
	ctx, cancel := context.WithCancel(ctx)
	nc := websocket.NetConn(ctx, conn, websocket.MessageBinary)
	return &connWithAddr{
		TunnelConn: &connWithClose{
			TunnelConn: nc,
			close:      func() error { cancel(); return nc.Close() },
		},
		remoteAddr: connectAddr(tunnelOpts.RemoteAddr),
		localAddr:  connectAddr(tunnelOpts.LocalAddr),
	}, nil
}

type connWithClose struct {
	TunnelConn
	close func() error
}

func (c *connWithClose) Close() error { return c.close() }

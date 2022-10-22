package ws

import (
	"context"
	"errors"
	"io"
	"net/http"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wspb"

	"go.minekube.com/connect"
	"go.minekube.com/connect/internal/ctxutil"
	"go.minekube.com/connect/internal/netutil"
	"go.minekube.com/connect/internal/util"
)

// ClientOptions for Watch and Tunnel, implementing connect.Tunneler.
type ClientOptions struct {
	URL         string                // The URL of the WebSocket server
	DialContext context.Context       // Optional WebSocket dial context
	DialOptions websocket.DialOptions // Optional WebSocket dial options
	Handshake   HandshakeResponse     // Optionally run after successful WebSocket handshake
}

// HandshakeResponse is called after receiving the
// WebSocket handshake response from the server.
type HandshakeResponse func(ctx context.Context, res *http.Response) (context.Context, error)

// Tunnel implements connect.Tunneler and creates a connection over a WebSocket.
// On error a http.Response may be provided by DialErrorResponse.
func (o ClientOptions) Tunnel(ctx context.Context) (connect.Tunnel, error) {
	ctx, ws, err := o.dial(ctx)
	if err != nil {
		return nil, err
	}

	// Extract additional options
	opts := ctxutil.TunnelOptionsOrDefault(ctx)

	// Return connection
	ctx, cancel := context.WithCancel(ctx)
	wsConn := websocket.NetConn(ctx, ws, websocket.MessageBinary)
	return &netutil.Conn{
		Conn:        wsConn,
		CloseFn:     func() error { cancel(); return wsConn.Close() },
		VLocalAddr:  opts.LocalAddr,
		VRemoteAddr: opts.RemoteAddr,
	}, nil
}

// Watch implements connect.Watcher and watches for session proposals.
// On error a http.Response may be provided by DialErrorResponse.
func (o ClientOptions) Watch(ctx context.Context, propose connect.ReceiveProposal) error {
	ctx, ws, err := o.dial(ctx)
	if err != nil {
		return err
	}
	// Watch for session proposals
	for {
		res := new(connect.WatchResponse)
		err = wspb.Read(ctx, ws, res)
		if err != nil {
			break
		}
		if res.GetSession() == nil {
			continue
		}
		id := res.GetSession().GetId()
		rejectFn := func(ctx context.Context, r *connect.RejectionReason) error {
			return wspb.Write(ctx, ws, &connect.WatchRequest{
				SessionRejection: &connect.SessionRejection{
					Id:     id,
					Reason: r,
				},
			})
		}
		proposal := &util.SessionProposal{
			Proposal: res.GetSession(),
			RejectFn: rejectFn,
		}
		if err = propose(proposal); err != nil {
			break
		}
	}
	_ = ws.Close(websocket.StatusNormalClosure, err.Error())
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return err
}

var _ connect.Tunneler = (*ClientOptions)(nil)
var _ connect.Watcher = (*ClientOptions)(nil)

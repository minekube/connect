package ws

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
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
	Handshake   HandshakeResponse     // Optional callback
}

// HandshakeResponse is called after receiving the
// WebSocket handshake response from the server.
type HandshakeResponse func(ctx context.Context, res *http.Response) (context.Context, error)

// Tunnel implements connect.Tunneler and creates a connection over a WebSocket.
func (o ClientOptions) Tunnel(ctx context.Context) (connect.Tunnel, error) {
	ctx, ws, err := o.dial(ctx)
	if err != nil {
		return nil, err
	}

	// Extract additional options
	opts := ctxutil.TunnelOptions(ctx)
	if opts.LocalAddr == nil {
		opts.LocalAddr = connect.Addr("unknown")
	}
	if opts.RemoteAddr == nil {
		if p, ok := peer.FromContext(ctx); ok {
			opts.RemoteAddr = p.Addr
		} else {
			opts.RemoteAddr = connect.Addr("unknown")
		}
	}

	// Return connection
	ctx, cancel := context.WithCancel(ctx)
	conn := websocket.NetConn(ctx, ws, websocket.MessageBinary)
	return &netutil.Conn{
		Conn:        conn,
		CloseFn:     func() error { cancel(); return conn.Close() },
		VLocalAddr:  opts.LocalAddr,
		VRemoteAddr: opts.RemoteAddr,
	}, nil
}

// Watch implements connect.Watcher and watches for session proposals.
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

func (o *ClientOptions) dial(ctx context.Context) (context.Context, *websocket.Conn, error) {
	if o.URL == "" {
		return nil, nil, errors.New("missing websocket url")
	}

	// Add metadata to websocket handshake request header
	md, _ := metadata.FromOutgoingContext(ctx)
	if o.DialOptions.HTTPHeader == nil {
		o.DialOptions.HTTPHeader = http.Header(md)
	} else {
		header := metadata.Join(metadata.MD(o.DialOptions.HTTPHeader), md)
		o.DialOptions.HTTPHeader = http.Header(header)
	}

	// Dial service
	if o.DialContext == nil {
		o.DialContext = ctx
	}
	ws, res, err := websocket.Dial(o.DialContext, o.URL, &o.DialOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("error handshaking with websocket server: %w", err)
	}

	// Callback for handshake response
	if o.Handshake != nil {
		ctx, err = o.Handshake(ctx, res)
		if err != nil {
			_ = ws.Close(websocket.StatusNormalClosure, fmt.Sprintf(
				"handshake response rejected: %v", err))
			return nil, nil, err
		}
	}

	return ctx, ws, nil
}

var _ connect.Tunneler = (*ClientOptions)(nil)
var _ connect.Watcher = (*ClientOptions)(nil)

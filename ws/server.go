package ws

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"google.golang.org/grpc/metadata"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wspb"

	"go.minekube.com/connect"
	"go.minekube.com/connect/internal/ctxutil"
	"go.minekube.com/connect/internal/netutil"
	"go.minekube.com/connect/internal/util"
)

// ServerOptions for TunnelHandler and EndpointHandler.
type ServerOptions struct {
	AcceptOptions websocket.AcceptOptions // Optional WebSocket accept options
}

// RequestFromContext returns the accepted WebSocket request from the context.
func RequestFromContext(ctx context.Context) *http.Request {
	r, _ := ctx.Value(httpRequestContextKey{}).(*http.Request)
	return r
}

// TunnelHandler returns a new http.Handler for accepting WebSocket requests for tunneling.
func (o ServerOptions) TunnelHandler(ln connect.TunnelListener) (http.Handler, error) {
	fn := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		// Accept WebSocket
		ws, err := websocket.Accept(w, r, &o.AcceptOptions)
		if err != nil {
			return err
		}

		// Extract additional options
		opts := ctxutil.TunnelOptionsOrDefault(ctx)

		// Create tunnel from WebSocket
		ctx, cancel := context.WithCancel(ctx)
		wsConn := websocket.NetConn(ctx, ws, websocket.MessageBinary)
		conn := &netutil.Conn{
			Conn:        wsConn,
			CloseFn:     func() error { cancel(); return wsConn.Close() },
			VLocalAddr:  opts.LocalAddr,
			VRemoteAddr: opts.RemoteAddr,
		}
		defer conn.Close()

		// Add WebSocket handshake request header to ctx metadata
		md, _ := metadata.FromIncomingContext(ctx)
		md = metadata.Join(md, metadata.MD(r.Header))
		ctx = metadata.NewIncomingContext(ctx, md)

		// Add http request to ctx
		ctx = withRequest(ctx, r)

		// Accept tunnel
		if err = ln.AcceptTunnel(ctx, conn); err != nil {
			// Specify meaningful close error
			_ = ws.Close(websocket.StatusProtocolError, fmt.Sprintf(
				"did not accept tunnel: %v", err))
			return err
		}

		// Block handler until tunnel closure
		<-ctx.Done()
		return nil
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Dropping this error as http.Error(...) would be already called
		// at this point by our WebSocket library.
		_ = fn(r.Context(), w, r)
	}), nil
}

// EndpointHandler returns a new http.Handler for accepting WebSocket requests for watching endpoints.
func (o ServerOptions) EndpointHandler(ln connect.EndpointListener) (http.Handler, error) {
	fn := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		// Accept WebSocket
		ws, err := websocket.Accept(w, r, &o.AcceptOptions)
		if err != nil {
			return err
		}

		// Prepare endpoint watch
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		ew := &util.EndpointWatch{
			ProposeFn: func(session *connect.Session) error {
				return wspb.Write(ctx, ws, &connect.WatchResponse{Session: session})
			},
			Receive: func() (*connect.WatchRequest, error) {
				req := new(connect.WatchRequest)
				return req, wspb.Read(ctx, ws, req)
			},
			RejectionsChan: make(chan *connect.SessionRejection),
		}

		// Receive session rejections from endpoint
		go func() { ew.StartReceiveRejections(ctx); cancel() }()
		go pingLoop(ctx, pingInterval, ws)

		// Add http request to ctx
		ctx = withRequest(ctx, r)

		// Start blocking watch callback
		if err = ln.AcceptEndpoint(ctx, ew); err != nil {
			// Specify meaningful close error
			_ = ws.Close(websocket.StatusProtocolError, fmt.Sprintf(
				"did not accept endpoint: %v", err))
			return err
		}

		return nil
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Dropping this error as http.Error(...) would be already called
		// at this point by our WebSocket library.
		_ = fn(r.Context(), w, r)
	}), nil
}

const pingInterval = time.Second * 50

// Send periodic pings to keep the WebSocket active since some web proxies
// timeout the connection after 60-100 seconds.
// https://community.cloudflare.com/t/cloudflare-websocket-timeout/5865/2
func pingLoop(ctx context.Context, d time.Duration, ws *websocket.Conn) {
	t := time.NewTicker(d)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			_ = ws.Ping(ctx)
		case <-ctx.Done():
			return
		}
	}
}

type httpRequestContextKey struct{}

func withRequest(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, httpRequestContextKey{}, r)
}

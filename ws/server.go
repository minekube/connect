package ws

import (
	"context"
	"net/http"

	"nhooyr.io/websocket"

	"go.minekube.com/connect"
)

type ServerOptions struct {
	AcceptOptions websocket.AcceptOptions
	Handshake     HandshakeRequest // Optional callback
}

// HandshakeRequest is called before accepting a
// WebSocket handshake request from the client.
type HandshakeRequest func(ctx context.Context, res *http.Response) (context.Context, error)

func (o ServerOptions) EndpointHandler(ln connect.EndpointListener) (http.Handler, error) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	}), nil
}

func (o ServerOptions) TunnelHandler(ln connect.TunnelListener) (http.Handler, error) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ws, err := websocket.Accept(w, r, &o.AcceptOptions)
		if err != nil {
			return
		}

	}), nil
}

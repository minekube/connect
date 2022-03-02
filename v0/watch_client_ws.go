package v0

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wspb"
)

type WatchWebsocketOptions struct {
	Callback func(proposal SessionProposal) error // The callback that is called when receiving session proposals.

	HTTPHeader http.Header // The HTTP headers included in the handshake request.
	// HTTPClient is used for the connection.
	// Its Transport must return writable bodies for WebSocket handshakes.
	HTTPClient      *http.Client
	HandshakeResult func(res *http.Response) error
}

func WatchWebsocket(ctx context.Context, url string, watchOpts WatchWebsocketOptions) error {
	// Dial service
	conn, resp, err := websocket.Dial(ctx, url, &websocket.DialOptions{
		HTTPClient:   watchOpts.HTTPClient,
		HTTPHeader:   watchOpts.HTTPHeader,
		Subprotocols: wsWatchSvcOpts.accept.Subprotocols,
	})
	if err != nil {
		return fmt.Errorf("error dialing %q: %w", url, err)
	}
	if watchOpts.HandshakeResult != nil {
		if err = watchOpts.HandshakeResult(resp); err != nil {
			_ = conn.Close(websocket.StatusNormalClosure, err.Error())
			return err
		}
	}
	// Watch for session proposals
	for {
		res := new(WatchResponse)
		err = wspb.Read(ctx, conn, res)
		if err != nil {
			break
		}
		if res.GetSession() == nil {
			continue
		}
		proposal := &websocketSessionProposal{
			s:    res.GetSession(),
			conn: conn,
			ctx:  ctx,
		}
		if err = watchOpts.Callback(proposal); err != nil {
			break
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	if err != nil {
		_ = conn.Close(websocket.StatusNormalClosure, err.Error())
	} else {
		_ = conn.Close(websocket.StatusNormalClosure, "EOF")
	}
	return err
}

type websocketSessionProposal struct {
	s    *Session
	ctx  context.Context
	conn *websocket.Conn
}

func (s *websocketSessionProposal) Session() *Session { return s.s }
func (s *websocketSessionProposal) Reject(reason *RejectionReason) error {
	return wspb.Write(s.ctx, s.conn, &WatchRequest{SessionRejection: &SessionRejection{
		Id:     s.s.GetId(),
		Reason: reason,
	}})
}

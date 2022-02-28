package connect

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wspb"
)

type websocketOptions struct {
	requiredSubProtocol string
	accept              *websocket.AcceptOptions
}

var wsWatchSvcOpts = &websocketOptions{
	accept: &websocket.AcceptOptions{
		Subprotocols:       []string{"connect"},
		InsecureSkipVerify: true,
	},
	requiredSubProtocol: "connect",
}

func acceptWebsocket(w http.ResponseWriter, r *http.Request, opts *websocketOptions) (*websocket.Conn, bool) {
	conn, err := websocket.Accept(w, r, wsWatchSvcOpts.accept)
	if err != nil {
		return nil, false
	}
	if conn.Subprotocol() != wsWatchSvcOpts.requiredSubProtocol {
		_ = conn.Close(websocket.StatusProtocolError,
			fmt.Sprintf("only supporting protocol: %s", wsWatchSvcOpts.requiredSubProtocol))
		return nil, false
	}
	return conn, true
}

// ServeHTTP accepts a websocket for the bidirectional streaming.
func (s *WatchService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, ok := acceptWebsocket(w, r, wsWatchSvcOpts)
	if !ok {
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	ww := &websocketWatcher{
		ctx:        r.Context(),
		conn:       conn,
		rejections: make(chan *SessionRejection),
	}
	go ww.startRecvRejections()
	// Start blocking watch callback
	if err := s.StartWatch(ww); err != nil {
		// todo map grpc status codes to websocket codes?
		//  https://www.iana.org/assignments/websocket/websocket.xhtml#close-code-number
		_ = conn.Close(websocket.StatusProtocolError, err.Error())
	}
}

type websocketWatcher struct {
	ctx        context.Context
	conn       *websocket.Conn
	rejections chan *SessionRejection
}

var _ Watcher = (*websocketWatcher)(nil)

func (w *websocketWatcher) Propose(session *Session) error {
	if session == nil {
		return errors.New("session must not be nil")
	}
	if session.GetId() == "" {
		return errors.New("missing session id")
	}
	return wspb.Write(w.ctx, w.conn, &WatchResponse{Session: session})
}
func (w *websocketWatcher) Context() context.Context             { return w.ctx }
func (w *websocketWatcher) Rejections() <-chan *SessionRejection { return w.rejections }
func (w *websocketWatcher) startRecvRejections() {
	defer close(w.rejections)
	for {
		req := new(WatchRequest)
		err := wspb.Read(w.ctx, w.conn, req)
		if err != nil {
			return // drop error
		}
		if req.GetSessionRejection() == nil {
			continue // Unexpected
		}
		select {
		case w.rejections <- req.GetSessionRejection():
		case <-w.ctx.Done():
			return // stream closed, done
		}
	}
}

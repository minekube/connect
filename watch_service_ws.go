package connect

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

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

func acceptWebsocket(w http.ResponseWriter, r *http.Request, opts *websocketOptions) (*websocket.Conn, error) {
	conn, err := websocket.Accept(w, r, opts.accept)
	if err != nil {
		return nil, err
	}
	if conn.Subprotocol() != opts.requiredSubProtocol {
		err = fmt.Errorf("only supporting protocol: %s", opts.requiredSubProtocol)
		_ = conn.Close(websocket.StatusProtocolError, err.Error())
		return nil, err
	}
	return conn, nil
}

// ServeHTTP accepts a websocket for the bidirectional streaming.
func (s *WatchService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := acceptWebsocket(w, r, wsWatchSvcOpts)
	if err != nil {
		fmt.Println("WatchService acceptWebsocket error", err)
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "closed serverside by watch service")

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	ww := &websocketWatcher{
		ctx:        ctx, // TODO don't use request context (never done)
		conn:       conn,
		rejections: make(chan *SessionRejection),
	}
	go ww.startRecvRejections(cancel)

	// Send periodic pings to keep the websocket active since some web proxies
	// timeout the connection after 60-100 seconds.
	// https://community.cloudflare.com/t/cloudflare-websocket-timeout/5865/2
	go func() {
		t := time.NewTicker(time.Minute)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				_ = conn.Ping(ctx)
			case <-ctx.Done():
				return
			}
		}
	}()

	// Start blocking watch callback
	if err = s.StartWatch(ww); err != nil {
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
func (w *websocketWatcher) startRecvRejections(connClosed context.CancelFunc) {
	defer connClosed()
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

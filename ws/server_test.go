package ws

import (
	"context"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/coder/websocket"

	"go.minekube.com/connect"
	"go.minekube.com/connect/internal/wspb"
)

type endpointListenerFunc func(context.Context, connect.EndpointWatch) error

func (f endpointListenerFunc) AcceptEndpoint(ctx context.Context, watch connect.EndpointWatch) error {
	return f(ctx, watch)
}

// TestEndpointHandlerSeparatesConnectionAndRequestContexts exercises the
// concurrent rejection reader while the handler attaches the HTTP request to
// the listener context. Run it with -race: the reader must never observe a
// reassigned context interface.
func TestEndpointHandlerSeparatesConnectionAndRequestContexts(t *testing.T) {
	var receivedRequest atomic.Bool
	rejections := make(chan *connect.SessionRejection, 1)
	listener := endpointListenerFunc(func(ctx context.Context, watch connect.EndpointWatch) error {
		receivedRequest.Store(RequestFromContext(ctx) != nil)
		select {
		case rejection, ok := <-watch.Rejections():
			if ok {
				rejections <- rejection
			}
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})

	server := httptest.NewServer(ServerOptions{}.EndpointHandler(listener))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ws, _, err := websocket.Dial(ctx, "ws"+strings.TrimPrefix(server.URL, "http"), nil)
	if err != nil {
		t.Fatalf("dial endpoint watch: %v", err)
	}
	defer ws.Close(websocket.StatusNormalClosure, "test complete")

	want := &connect.SessionRejection{Id: "race-regression"}
	if err := wspb.Write(ctx, ws, &connect.WatchRequest{SessionRejection: want}); err != nil {
		t.Fatalf("write rejection: %v", err)
	}

	select {
	case got := <-rejections:
		if got.GetId() != want.GetId() {
			t.Fatalf("received rejection ID = %q, want %q", got.GetId(), want.GetId())
		}
	case <-ctx.Done():
		t.Fatalf("did not receive rejection: %v", ctx.Err())
	}
	if !receivedRequest.Load() {
		t.Fatal("endpoint listener did not receive the HTTP request context")
	}
}

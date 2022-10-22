package ws

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"google.golang.org/grpc/metadata"
	"nhooyr.io/websocket"
)

// DialErrorResponse returns the HTTP response from the WebSocket handshake error, if any.
func DialErrorResponse(err error) (*http.Response, bool) {
	var e *dialErr
	if errors.As(err, &e) {
		return e.res, true
	}
	return nil, false
}

func (o *ClientOptions) dial(ctx context.Context) (context.Context, *websocket.Conn, error) {
	if o.URL == "" {
		return nil, nil, errors.New("missing websocket url")
	}

	header := metadata.Join(
		metadata.MD(o.DialOptions.HTTPHeader),
		mdFromContext(o.DialContext),
		mdFromContext(ctx),
	)
	if o.DialContext != nil {
		ctx = o.DialContext
	}

	// Dial service
	ws, res, err := websocket.Dial(ctx, o.URL, &websocket.DialOptions{
		HTTPClient:           o.DialOptions.HTTPClient,
		HTTPHeader:           http.Header(header),
		Subprotocols:         o.DialOptions.Subprotocols,
		CompressionMode:      o.DialOptions.CompressionMode,
		CompressionThreshold: o.DialOptions.CompressionThreshold,
	})
	if err != nil {
		if res != nil {
			err = &dialErr{error: err, res: res}
		}
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

type dialErr struct {
	error
	res *http.Response
}

func (e *dialErr) Error() string {
	return fmt.Sprintf("%s (%d): %v", e.res.Status, e.res.StatusCode, e.error)
}
func (e *dialErr) Unwrap() error { return e.error }

func mdFromContext(ctx context.Context) metadata.MD {
	if ctx == nil {
		return nil
	}
	md, _ := metadata.FromOutgoingContext(ctx)
	return md
}

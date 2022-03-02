package ws

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"

	"go.minekube.com/connect"
	"go.minekube.com/connect/internal/testutil"
)

func TestSuite(t *testing.T) {
	suite.Run(t, &testutil.Suite{
		Endpoint: ClientOptions{
			URL: "ws://:8080",
		},
		StartWatchServer: func(ctx context.Context, ln connect.EndpointListener) error {
			return startServer(ctx, ":8080", ServerOptions{}.EndpointHandler(ln))
		},
		StartTunnelServer: func(ctx context.Context, ln connect.TunnelListener) error {
			return startServer(ctx, ":8080", ServerOptions{}.TunnelHandler(ln))
		},
	})
}

func startServer(ctx context.Context, addr string, handler http.Handler) error {
	svr := &http.Server{Addr: addr, Handler: handler}
	go func() { _ = svr.ListenAndServe() }()
	<-ctx.Done()
	return svr.Close()
}

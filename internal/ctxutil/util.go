package ctxutil

import (
	"context"

	"google.golang.org/grpc/peer"

	"go.minekube.com/connect"
	"go.minekube.com/connect/internal/ctxkey"
)

func tunnelOptions(ctx context.Context) connect.TunnelOptions {
	opts, _ := ctx.Value(ctxkey.TunnelOptions{}).(connect.TunnelOptions)
	return opts
}

func TunnelOptionsOrDefault(ctx context.Context) connect.TunnelOptions {
	opts := tunnelOptions(ctx)
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
	return opts
}

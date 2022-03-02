package ctxutil

import (
	"context"

	"go.minekube.com/connect"
	"go.minekube.com/connect/internal/ctxkey"
)

func TunnelOptions(ctx context.Context) connect.TunnelOptions {
	opts, _ := ctx.Value(ctxkey.TunnelOptsCtxKey{}).(connect.TunnelOptions)
	return opts
}

package connect

import (
	"context"
	"net"

	"go.minekube.com/connect/internal/ctxkey"
)

// Addr is an address in the "connect" network.
type Addr string

func (a Addr) String() string  { return string(a) }
func (a Addr) Network() string { return "connect" }

// TunnelOptions are common options for a call to Tunneler.Tunnel.
type TunnelOptions struct {
	// LocalAddr fakes the local address of the connection when specified.
	LocalAddr net.Addr // It is recommended to use the endpoint name.
	// The remote address as reflected by Tunnel.RemoteAddr().
	RemoteAddr net.Addr // It is recommended to use the player address.
}

func WithTunnelOptions(ctx context.Context, opts TunnelOptions) context.Context {
	return context.WithValue(ctx, ctxkey.TunnelOptsCtxKey{}, opts)
}

package connect

import (
	"context"
	"net"

	"go.minekube.com/connect/internal/ctxkey"
)

// Tunnel represents an outbound only tunnel initiated by
// an Endpoint for a specific SessionProposal.
type Tunnel net.Conn

// Tunneler creates a Tunnel.
type Tunneler interface {
	Tunnel(context.Context) (Tunnel, error)
}

// Watcher registers the calling endpoint and watches for sessions
// proposed by the WatchService. To stop watching cancel the context.
// If ReceiveProposal returns a non-nil non-EOF error Watch returns it.
type Watcher interface {
	Watch(context.Context, ReceiveProposal) error
}

// ReceiveProposal is called when Watcher receives a SessionProposal.
type ReceiveProposal func(proposal SessionProposal) error

// SessionProposal specifies an incoming session proposal.
// Use the Session to create the connection tunnel or reject the session with an optional reason.
type SessionProposal interface {
	Session() *Session                              // The session proposed to connect to the Endpoint.
	Reject(context.Context, *RejectionReason) error // Rejects the session proposal with an optional reason.
}

// Endpoint is an endpoint that listens for
// sessions to either reject them or establish
// a tunnel for receiving the connection.
type Endpoint interface {
	Watcher
	Tunneler
}

// TunnelListener is a network listener for tunnel connections.
type TunnelListener interface {
	AcceptTunnel(context.Context, Tunnel) error
}

// EndpointListener is a network listener for endpoint watches.
type EndpointListener interface {
	AcceptEndpoint(context.Context, EndpointWatch) error
}

// EndpointWatch is a watching Endpoint.
type EndpointWatch interface {
	// Propose proposes a session to the Endpoint.
	// The Endpoint either rejects the proposal or initiates
	// a Tunnel to receive the session connection.
	Propose(*Session) error
	Rejections() <-chan *SessionRejection // Rejections receives rejected session proposals from the Endpoint.
}

// Listener listens for watching endpoints and tunnel connections from endpoints.
type Listener interface {
	EndpointListener
	TunnelListener
}

// TunnelOptions are options for Tunneler and TunnelListener.
// Use WithTunnelOptions to propagate TunnelOptions in context.
type TunnelOptions struct {
	// LocalAddr fakes the local address of the Tunnel when specified.
	//
	// If this TunnelOptions destine to Tunneler
	// it is recommended to use the Endpoint address/name.
	//
	// If this TunnelOptions destine to TunnelListener
	// it is recommended to use the underlying network listener address.
	LocalAddr net.Addr
	// RemoteAddr fakes the remote address of the Tunnel when specified.
	//
	// If this TunnelOptions destine to Tunneler
	// it is recommended to use the underlying connection remote address.
	//
	// If this TunnelOptions destine to TunnelListener
	// it is recommended to use the Endpoint address/name.
	RemoteAddr net.Addr // It is recommended to use the player address.
}

// WithTunnelOptions stores TunnelOptions in a context.
func WithTunnelOptions(ctx context.Context, opts TunnelOptions) context.Context {
	return context.WithValue(ctx, ctxkey.TunnelOptions{}, opts)
}

// Addr is an address in the "connect" network.
type Addr string

func (a Addr) String() string  { return string(a) }
func (a Addr) Network() string { return "connect" }

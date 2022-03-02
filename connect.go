package connect

import (
	"context"
	"net"
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

// Listener listens for watching endpoints and tunnel sessions from endpoints
type Listener interface {
	EndpointListener
	TunnelListener
}

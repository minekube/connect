package connect

import (
	"net"
	"net/http"

	"google.golang.org/grpc/peer"
	"nhooyr.io/websocket"
)

var wsTunnelSvcOpts = &websocketOptions{
	accept: &websocket.AcceptOptions{
		Subprotocols:       []string{"minecraft-java"},
		InsecureSkipVerify: true,
	},
	requiredSubProtocol: "minecraft-java",
}

// ServeHTTP accepts a websocket for the bidirectional streaming.
func (s *TunnelService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get address of client
	var remoteAddr net.Addr
	if p, ok := peer.FromContext(r.Context()); !ok { // TODO use a remoteAddr fn instead
		http.Error(w, "could not resolve address of client", http.StatusBadRequest)
		return
	} else {
		remoteAddr = p.Addr
	}

	conn, ok := acceptWebsocket(w, r, wsTunnelSvcOpts)
	if !ok {
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	nc := &connWithAddr{
		TunnelConn: websocket.NetConn(r.Context(), conn, websocket.MessageBinary),
		remoteAddr: remoteAddr,
		localAddr:  s.LocalAddr,
	}
	tunnel := &inboundTunnel{
		close: nc.Close,
		conn:  func() TunnelConn { return nc },
		ctx:   r.Context,
	}

	// Call inbound tunnel callback
	if err := s.AcceptTunnel(tunnel); err != nil {
		_ = conn.Close(websocket.StatusProtocolError, err.Error())
		return
	}
	// Block until tunnel closed
	<-r.Context().Done()
}

// overload address methods
type connWithAddr struct {
	TunnelConn
	remoteAddr net.Addr
	localAddr  net.Addr
}

func (c *connWithAddr) RemoteAddr() net.Addr { return c.remoteAddr }
func (c *connWithAddr) LocalAddr() net.Addr  { return c.localAddr }

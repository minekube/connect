package connect

import (
	"context"
	"fmt"
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

	conn, err := acceptWebsocket(w, r, wsTunnelSvcOpts)
	if err != nil {
		fmt.Println(2, err)
		return
	}

	// Create inbound tunnel from websocket
	ctx, cancel := context.WithCancel(r.Context())
	nc := &connWithAddr{
		TunnelConn: websocket.NetConn(ctx, conn, websocket.MessageBinary),
		remoteAddr: remoteAddr,
		localAddr:  s.LocalAddr,
	}
	closeFn := func() error {
		fmt.Println("tunnel conn closed")
		cancel()
		return conn.Close(websocket.StatusNormalClosure, "closed serverside by tunnel service")
	}
	tunnel := &inboundTunnel{
		close: closeFn,
		conn: func() TunnelConn {
			return &connWithClose{
				TunnelConn: nc,
				close:      closeFn,
			}
		},
		ctx: func() context.Context { return ctx },
	}
	defer tunnel.close()

	// Call inbound tunnel callback
	if err = s.AcceptTunnel(tunnel); err != nil {
		_ = conn.Close(websocket.StatusProtocolError, err.Error())
		return
	}
	// Block until tunnel closed
	fmt.Println("tunnel waiting ctx done")
	<-ctx.Done()
	fmt.Println("tunnel req ctx done")
}

// overload address methods
type connWithAddr struct {
	TunnelConn
	remoteAddr net.Addr
	localAddr  net.Addr
}

func (c *connWithAddr) RemoteAddr() net.Addr { return c.remoteAddr }
func (c *connWithAddr) LocalAddr() net.Addr  { return c.localAddr }

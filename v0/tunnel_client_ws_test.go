package v0

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/peer"
)

func TestTunnelWebsocket(t *testing.T) {
	const addr = ":8080"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	go func() {
		conn, err := TunnelWebsocket(ctx, "ws://"+addr, TunnelWebsocketOptions{
			LocalAddr:  "endpointName",
			RemoteAddr: "playerIP",
		})
		require.NoError(t, err)
		n, err := conn.Write([]byte("hello"))
		require.NoError(t, err)
		require.Equal(t, 5, n)
		require.NoError(t, conn.Close())
		fmt.Println("closed")
	}()

	ts := &TunnelService{
		AcceptTunnel: func(tunnel InboundTunnel) error {
			fmt.Println("accepted tunnel", tunnel.Conn().RemoteAddr())
			b := make([]byte, 10)
			fmt.Println(tunnel.Conn().RemoteAddr())
			fmt.Println(tunnel.Conn().LocalAddr())
			n, err := tunnel.Conn().Read(b)
			require.NoError(t, err)
			require.Equal(t, 5, n)
			require.Equal(t, "hello", string(b[:5]))

			// for reading client close frames
			for {
				_, err = tunnel.Conn().Read(b)
				if err != nil {
					break
				}
			}
			return nil
		},
		LocalAddr: connectAddr("localhost"),
	}
	err := http.ListenAndServe(addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(peer.NewContext(r.Context(), &peer.Peer{Addr: connectAddr("remoteAddr")}))
		ts.ServeHTTP(w, r)
	}))
	require.NoError(t, err)
}

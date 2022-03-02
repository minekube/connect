package v0

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWatchWS(t *testing.T) {
	const addr = ":8080"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	go func() {
		var proposals int
		err := WatchWebsocket(ctx, "ws://"+addr, WatchWebsocketOptions{
			Callback: func(proposal SessionProposal) error {
				fmt.Println("got proposal", proposal.Session().GetId())
				proposals++
				if proposals == 2 {
					return io.EOF
				}
				return nil
			},
			HandshakeResult: func(res *http.Response) error {
				fmt.Println("handshaked")
				return nil
			},
		})
		require.NoError(t, err)
	}()

	ws := &WatchService{
		StartWatch: func(w Watcher) error {
			fmt.Println("new watcher")
			err := w.Propose(&Session{Id: "hello"})
			require.NoError(t, err)
			err = w.Propose(&Session{Id: "hello2"})
			require.NoError(t, err)
			<-w.Context().Done()
			return nil
		},
	}
	err := http.ListenAndServe(addr, ws)
	require.NoError(t, err)
}

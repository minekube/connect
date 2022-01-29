package connect

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"net"
	"testing"
	"time"
)

func TestWatchService_Serve(t *testing.T) {
	ln, err := net.Listen("tcp", ":8443")
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)

	stopped := make(chan struct{})
	go func() {
		defer close(stopped)
		defer ln.Close()
		ws := &WatchService{
			StartWatch: func(w Watcher) error {
				require.Equal(t, "foo", w.Endpoint().GetName())
				// just propose some sessions and then close
				err = w.Propose(&Session{Id: "abc"})
				require.NoError(t, err)
				err = w.Propose(&Session{Id: "abc"})
				require.NoError(t, err)
				return w.Propose(&Session{Id: "abc"})
			},
		}
		err = ws.Serve(ctx.Done(), ln)
		require.NoError(t, err)
	}()

	cliConn, err := grpc.DialContext(ctx, ":8443", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	watchCli := NewWatchServiceClient(cliConn)
	var proposals int
	err = Watch(ctx, WatchOptions{
		Cli:      watchCli,
		Endpoint: "foo",
		Callback: func(proposal SessionProposal) error {
			proposals++
			require.Equal(t, "abc", proposal.Session().GetId())
			require.Equal(t, "foo", proposal.Endpoint().GetName())
			if proposals == 3 {
				cancel()
			}
			return nil
		},
	})
	if !errors.Is(err, context.Canceled) && status.Code(err) != codes.Canceled {
		require.NoError(t, err)
	}
	require.Equal(t, 3, proposals)

	select {
	case <-stopped:
	case <-time.After(time.Second):
		require.Fail(t, "WatchService not shutdown gracefully")
	}
}

func TestSessionProposal_Reject(t *testing.T) {
	ln, err := net.Listen("tcp", ":8444")
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)

	stopped := make(chan struct{})
	go func() {
		defer close(stopped)
		defer ln.Close()
		ws := &WatchService{
			StartWatch: func(w Watcher) error {
				go func() {
					err = w.Propose(&Session{Id: "abc"})
					require.NoError(t, err)
					err = w.Propose(&Session{Id: "abc"})
					require.NoError(t, err)
				}()

				var rejections int
				for rejection := range w.Rejections() {
					rejections++
					require.Equal(t, "abc", rejection.GetId())
					if rejections == 2 {
						break
					}
				}
				return nil // stops watcher
			},
		}
		err = ws.Serve(ctx.Done(), ln)
		require.NoError(t, err)
	}()

	cliConn, err := grpc.DialContext(ctx, ":8444", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	watchCli := NewWatchServiceClient(cliConn)
	var proposals int
	err = Watch(ctx, WatchOptions{
		Cli:      watchCli,
		Endpoint: "foo",
		Callback: func(proposal SessionProposal) error {
			proposals++
			proposal.Reject(nil) // reject all
			return nil
		},
	})
	require.NoError(t, err)
	require.Equal(t, 2, proposals)

	cancel()

	select {
	case <-stopped:
	case <-time.After(time.Second):
		require.Fail(t, "WatchService not shutdown gracefully")
	}
}

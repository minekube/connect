package testutil

import (
	"context"
	"fmt"
	"io"
	"sync/atomic"
	"time"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.minekube.com/connect"
)

type Suite struct {
	suite.Suite

	Endpoint          connect.Endpoint
	StartWatchServer  func(ctx context.Context, ln connect.EndpointListener) error
	StartTunnelServer func(ctx context.Context, ln connect.TunnelListener) error
}

func (suite *Suite) SetupTest() {
	suite.NotNil(suite.Endpoint)
	suite.NotNil(suite.StartWatchServer)
}

func (suite *Suite) TestWatchReject() {
	proposal := &connect.Session{Id: "abc"}
	rejection := status.New(codes.Aborted, "don't want this").Proto()

	var expSeq = []string{
		"s: got watcher",
		"c: got proposal " + proposal.GetId(),
		"c: rejection sent",
		"s: got rejection " + rejection.String(),
	}
	var seq sequence

	ctx, stop := context.WithTimeout(context.TODO(), time.Second*3)

	ln := acceptEndpoint(func(ctx context.Context, watch connect.EndpointWatch) error {
		seq.Add("s: got watcher")
		suite.NoError(watch.Propose(ctx, proposal))
		for rej := range watch.Rejections() {
			seq.Add("s: got rejection " + rej.GetReason().String())
			break
		}
		time.Sleep(time.Millisecond * 100) // let "c: rejection sent"

		return nil
	})

	go func() { suite.NoError(suite.StartWatchServer(ctx, ln)) }()

	time.Sleep(time.Millisecond * 100) // Wait for server
	go func() {
		err := suite.Endpoint.Watch(ctx, func(proposal connect.SessionProposal) error {
			seq.Add("c: got proposal " + proposal.Session().GetId())
			suite.NoError(proposal.Reject(ctx, rejection))
			seq.Add("c: rejection sent")
			return nil
		})
		suite.NotNil(err)
		suite.Contains(err.Error(), "closed serverside")
		stop()
	}()

	<-ctx.Done()
	suite.Assert().ErrorIs(ctx.Err(), context.Canceled)
	suite.Assert().Equal(expSeq, seq.Get())
}

func (suite *Suite) TestTunnel() {
	ctx, stop := context.WithTimeout(context.TODO(), time.Second*3)

	toClientMsg := []byte("hello client")
	toServerMsg := []byte("hello server")

	var expSeq = []string{
		"c: tunnel opened Tunnel(remote=unknown via connect, local=unknown via connect)",
		"s: got tunnel Tunnel(remote=unknown via connect, local=unknown via connect)",
		"s: read",
		"c: read",
		"c: no read",
	}
	var seq sequence

	ln := acceptTunnel(func(ctx context.Context, tunnel connect.Tunnel) error {
		time.Sleep(time.Millisecond * 100) // let "c: tunnel opened"
		seq.Add("s: got tunnel " + fmt.Sprint(tunnel))

		// client -> server
		b := make([]byte, 100)
		n, err := tunnel.Read(b)
		suite.NoError(err)
		suite.Equal(len(toServerMsg), n)
		suite.Equal(toServerMsg, b[:n])
		seq.Add("s: read")

		// client <- server
		n, err = tunnel.Write(toClientMsg)
		suite.NoError(err)
		suite.Equal(len(toClientMsg), n)

		// Close server side
		suite.NoError(tunnel.Close())
		return nil
	})

	go func() { suite.NoError(suite.StartTunnelServer(ctx, ln)) }()

	time.Sleep(time.Millisecond * 100) // Wait for server
	go func() {
		tunnel, err := suite.Endpoint.Tunnel(ctx)
		suite.NoError(err)
		seq.Add("c: tunnel opened " + fmt.Sprint(tunnel))

		// client -> server
		n, err := tunnel.Write(toServerMsg)
		suite.NoError(err)
		suite.Equal(len(toServerMsg), n)

		// client <- server
		b := make([]byte, 100)
		n, err = tunnel.Read(b)
		suite.NoError(err)
		suite.Equal(len(toClientMsg), n)
		suite.Equal(toClientMsg, b[:n])
		seq.Add("c: read")

		// Should be closed server side by now
		b = make([]byte, 100)
		for i := 0; i < 5; i++ {
			n, err = tunnel.Read(b)
			suite.Empty(n)
			if err == nil {
				continue // retry
			}
		}
		suite.ErrorIs(err, io.EOF)
		seq.Add("c: no read")

		_ = tunnel.Close()
		stop()
	}()

	<-ctx.Done()
	suite.Assert().ErrorIs(ctx.Err(), context.Canceled)
	suite.Assert().Equal(expSeq, seq.Get())
}

type sequence struct {
	v atomic.Value
}

func (s *sequence) Add(str string) {
	s.v.Store(append(s.Get(), str))
}

func (s *sequence) Get() []string {
	str, _ := s.v.Load().([]string)
	return str
}

type acceptEndpoint func(ctx context.Context, watch connect.EndpointWatch) error

func (fn acceptEndpoint) AcceptEndpoint(ctx context.Context, watch connect.EndpointWatch) error {
	return fn(ctx, watch)
}

type acceptTunnel func(ctx context.Context, tunnel connect.Tunnel) error

func (fn acceptTunnel) AcceptTunnel(ctx context.Context, tunnel connect.Tunnel) error {
	return fn(ctx, tunnel)
}

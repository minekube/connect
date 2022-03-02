package v0

import (
	"context"
	"net"
	"time"
)

// TunnelConn is a tunnel connection.
type TunnelConn net.Conn

func newTunnelConn(
	localAddr, remoteAddr net.Addr,
	r TunnelReader, w TunnelWriter,
	closeFn func() error,
) TunnelConn {
	return &tunnelConn{
		closeFn:    closeFn,
		localAddr:  localAddr,
		remoteAddr: remoteAddr,
		r:          r,
		w:          w,
	}
}

type tunnelConn struct {
	closeFn               func() error
	localAddr, remoteAddr net.Addr
	r                     TunnelReader
	w                     TunnelWriter
}

func (c *tunnelConn) Read(b []byte) (n int, err error)  { return c.r.Read(b) }
func (c *tunnelConn) Write(b []byte) (n int, err error) { return c.w.Write(b) }
func (c *tunnelConn) Close() error                      { return c.closeFn() }
func (c *tunnelConn) LocalAddr() net.Addr               { return c.localAddr }
func (c *tunnelConn) RemoteAddr() net.Addr              { return c.remoteAddr }
func (c *tunnelConn) SetDeadline(t time.Time) error {
	_ = c.SetWriteDeadline(t)
	return c.SetReadDeadline(t)
}
func (c *tunnelConn) SetReadDeadline(t time.Time) error  { return c.r.SetDeadline(t) }
func (c *tunnelConn) SetWriteDeadline(t time.Time) error { return c.w.SetDeadline(t) }

func tunnelServerRW(ctx context.Context, ss TunnelService_TunnelServer) (r TunnelReader, w TunnelWriter) {
	return newTunnelReader(ctx, func() ([]byte, error) { msg, err := ss.Recv(); return msg.GetData(), err }),
		newTunnelWriter(ctx, func(b []byte) (err error) { return ss.Send(&TunnelResponse{Data: b}) })
}

func tunnelClientRW(ctx context.Context, cs TunnelService_TunnelClient) (r TunnelReader, w TunnelWriter) {
	return newTunnelReader(ctx, func() ([]byte, error) { msg, err := cs.Recv(); return msg.GetData(), err }),
		newTunnelWriter(ctx, func(b []byte) (err error) { return cs.Send(&TunnelRequest{Data: b}) })
}

package netutil

import "net"

type Conn struct {
	net.Conn
	CloseFn                 func() error
	VLocalAddr, VRemoteAddr net.Addr
}

func (c *Conn) Close() error         { return c.CloseFn() }
func (c *Conn) LocalAddr() net.Addr  { return c.VLocalAddr }
func (c *Conn) RemoteAddr() net.Addr { return c.VRemoteAddr }

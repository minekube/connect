package netutil

import (
	"fmt"
	"net"
)

type Conn struct {
	net.Conn
	CloseFn                 func() error
	VLocalAddr, VRemoteAddr net.Addr
}

func (c *Conn) Close() error         { return c.CloseFn() }
func (c *Conn) LocalAddr() net.Addr  { return c.VLocalAddr }
func (c *Conn) RemoteAddr() net.Addr { return c.VRemoteAddr }
func (c *Conn) String() string {
	if s, ok := c.Conn.(fmt.Stringer); ok {
		return s.String()
	}
	return fmt.Sprintf("Tunnel(remote=%s via %s, local=%s via %s)",
		c.RemoteAddr().String(), c.RemoteAddr().Network(),
		c.LocalAddr().String(), c.LocalAddr().Network())
}

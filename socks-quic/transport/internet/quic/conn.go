package quic

import (
	"net"

	quic "github.com/lucas-clemente/quic-go"
)

type sysConn struct {
	net.PacketConn
}

// TODO implement ReadFrom, WriteTo with header stuff.
//func (c *sysConn) ReadFrom(p []byte) (int, net.Addr, error) {
//	return c.conn.ReadFrom(p)
//}
//
//func (c *sysConn) WriteTo(p []byte, addr net.Addr) (int, error) {
//	return c.conn.WriteTo(p, addr)
//}

type interConn struct {
	quic.Stream
	local  net.Addr
	remote net.Addr
}

// TODO why multiBuffer? do we really need this?
//func (c *interConn) WriteMultiBuffer(mb buf.MultiBuffer) error {
//	mb = buf.Compact(mb)
//	mb, err := buf.WriteMultiBuffer(c, mb)
//	buf.ReleaseMulti(mb)
//	return err
//}

func (c *interConn) LocalAddr() net.Addr {
	return c.local
}

func (c *interConn) RemoteAddr() net.Addr {
	return c.remote
}

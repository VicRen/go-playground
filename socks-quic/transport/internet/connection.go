package internet

import "net"

type Connection interface {
	net.Conn
}

type ConnHandler func(Connection)

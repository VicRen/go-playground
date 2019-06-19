package main

import (
	"context"
	"fmt"
	"net"
	"net/url"

	"golang.org/x/net/proxy"
)

type dialContext interface {
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}

type socksConn interface {
	BoundAddr() net.Addr
}

func main() {
	u, err := url.Parse("socks5://127.0.0.1:10080")
	if err != nil {
		panic(err)
	}
	dial, err := proxy.FromURL(u, proxy.Direct)
	if err != nil {
		panic(err)
	}
	if d, ok := dial.(dialContext); ok {
		c, err := d.DialContext(context.Background(), "tcp", "39.96.21.158:5000")
		if err != nil {
			panic(err)
		}
		if d, ok := c.(socksConn); ok {
			fmt.Println("yes!", d.BoundAddr())
		}
		c.Write([]byte{0x00})
		buf := make([]byte, 1024)
		n, err := c.Read(buf)
		if err != nil {
			panic(err)
		}
		fmt.Println("data read:", string(buf[:n]))
	}
}

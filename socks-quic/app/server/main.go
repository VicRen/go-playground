package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/VicRen/go-play-ground/socks-quic/socks"
	"github.com/VicRen/go-play-ground/socks-quic/transport/internet"
	"github.com/VicRen/go-play-ground/socks-quic/transport/internet/quic"
)

var host = flag.String("host", "0.0.0.0", "host")
var port = flag.String("port", "8221", "port")

func main() {
	addr := net.JoinHostPort(*host, *port)
	_, err := quic.Listen(context.Background(), addr, func(connection internet.Connection) {
		go func() {
			fmt.Println("accept conn from:", connection.RemoteAddr())
			defer connection.Close()
			dstAddr, err := readDstAddr(connection)
			if err != nil {
				fmt.Println("error handling request:", err)
				return
			}

			_, err = connection.Write([]byte{0x00})
			if err != nil {
				fmt.Println("error responding request:", err)
				return
			}
			fmt.Println("try proxying to:", dstAddr)
			conn, err := socks.DefaultDialer(context.Background(), "udp", dstAddr)
			if err != nil {
				fmt.Println("error handling request:", err)
				return
			}
			conn.SetDeadline(time.Now().Add(10 * time.Second))
			defer conn.Close()
			fmt.Println("proxying to:", conn.RemoteAddr())
			finishCh := make(chan struct{})
			go func() {
				n, _ := io.Copy(connection, conn)
				fmt.Println("remote->server copy finished:", n)
				finishCh <- struct{}{}
			}()
			n, _ := io.Copy(conn, connection)
			fmt.Println("server->remote copy finished:", n)
			<-finishCh
		}()
	})
	if err != nil {
		return
	}
	fmt.Println("server listening on:", addr)
	select {}
}

func readDstAddr(reader io.Reader) (string, error) {
	buf := make([]byte, 1)
	_, err := io.ReadFull(reader, buf)
	if err != nil {
		return "", err
	}
	if buf[0] != byte(0x73) {
		return "", errors.New("invalid request")
	}
	_, dstAddr, err := socks.ReadAddrSpec(reader)
	if err != nil {
		return "", err
	}
	return dstAddr.Address(), nil
}

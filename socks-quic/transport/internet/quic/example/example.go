package main

import (
	"context"
	"fmt"
	"io"
	"net"

	"github.com/VicRen/go-play-ground/socks-quic/transport/internet"
	"github.com/VicRen/go-play-ground/socks-quic/transport/internet/quic"
)

func main() {
	var l *quic.Listener
	defer func() {
		if l != nil {
			fmt.Println("closing server")
			l.Close()
		}
	}()
	go func() {
		l, _ = echoServer()
	}()

	clientMain()
}

func clientMain() error {
	conn, err := quic.Dial(context.Background(), &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8080})
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Write([]byte{0x00})
	if err != nil {
		return err
	}
	b := make([]byte, 10)
	n, err := conn.Read(b)
	if err != nil {
		return nil
	}
	fmt.Println("recv from server:", b[:n])
	return nil
}

func echoServer() (*quic.Listener, error) {
	l, err := quic.Listen(context.Background(), "127.0.0.1:8080", func(connection internet.Connection) {
		n, err := io.Copy(connection, connection)
		if err != nil {
			fmt.Println("server conn finished:", connection.RemoteAddr(), connection.LocalAddr(), n)
		}
		connection.Close()
	})
	if err != nil {
		return nil, err
	}
	fmt.Println("server listening")
	return l, nil
}

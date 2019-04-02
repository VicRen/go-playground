package main

import (
	"fmt"
	"io"
	"net"
)

func main() {
	go func() {
		Serve()
	}()
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		println(err)
	}

	conn.Write([]byte("hello"))

	if tcpConn, ok := conn.(*net.TCPConn); ok {
		err := tcpConn.CloseWrite()
		if err != nil {
			fmt.Println(err)
		}
	}

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("read data from server:", string(buf[:n]))
	}
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		err := tcpConn.CloseRead()
		if err != nil {
			fmt.Println(err)
		}
	}
	err = conn.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func Serve() {
	l, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		fmt.Println(err)
	}
	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
	}
	n, err := io.Copy(conn, conn)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("data received:", n)

}

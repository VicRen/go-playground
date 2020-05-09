package main

import (
	"fmt"
	"net"
)

func main() {
	ua, err := net.ResolveUDPAddr("udp", "127.0.0.1:20000")
	if err != nil {
		panic(err)
	}

	uc, err := net.ListenUDP("udp", ua)
	if err != nil {
		panic(err)
	}

	go func() {
		conn, err := net.Dial("udp", ua.String())
		if err != nil {
			panic(err)
		}
		b := make([]byte, 9*1024)
		n, err := conn.Write(b)
		if err != nil {
			panic(err)
		}
		fmt.Println("------->data written:", n)
		n, err = conn.Write(b)
		if err != nil {
			panic(err)
		}
		fmt.Println("------->data written2:", n)
	}()

	buf := make([]byte, 10000)
	for {
		n, err := uc.Read(buf)
		if err != nil {
			panic(err)
		}
		fmt.Println("------>data read", n)
	}
}

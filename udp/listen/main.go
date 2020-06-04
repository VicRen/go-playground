package main

import (
	"flag"
	"fmt"
	"net"
)

var (
	l = flag.String("l", "", "addr listening on")
)

func main() {
	flag.Parse()
	if *l == "" {
		panic("parse listen addr")
	}

	udpAddr, err := net.ResolveUDPAddr("udp", *l)
	if err != nil {
		panic(err)
	}
	fmt.Println("listening on", udpAddr)

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			panic(err)
		}
		fmt.Printf("reading data from [%v]: %v\n", addr, buf[:n])
	}
}

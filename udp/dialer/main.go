package main

import (
	"flag"
	"fmt"
	"net"
)

var (
	target = flag.String("t", "", "Target to dial.")
)

func main() {
	flag.Parse()

	if *target == "" {
		panic("target is empty")
	}

	addr, err := net.ResolveUDPAddr("udp", *target)
	if err != nil {
		panic(fmt.Sprintf("failed to resolve addr: %s", addr))
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		panic(fmt.Sprintf("failed to dial addr: %s", addr))
	}

	_, _ = conn.Write([]byte("hello world"))

	_ = conn.Close()
}

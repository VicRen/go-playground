package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
)

var (
	host  = flag.String("host", "47.102.45.160", "host to connect")
	port  = flag.Int("port", 10999, "port to connect")
	port2 = flag.Int("port2", 10998, "port to connect")
)

func main() {
	flag.Parse()
	addr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(*host, strconv.Itoa(*port)))
	if err != nil {
		panic(err)
	}
	addr2, err := net.ResolveUDPAddr("udp", net.JoinHostPort(*host, strconv.Itoa(*port2)))
	if err != nil {
		panic(err)
	}
	dial(addr, addr2)
}

func dial(addr, addr2 *net.UDPAddr) {
	c, err := net.ListenUDP("udp", nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("listening on:%v\n", c.LocalAddr())

	_, err = c.WriteTo([]byte("hello world"), addr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("client send: %v->%v\n", c.LocalAddr(), addr)

	_, err = c.WriteTo([]byte("hello world 2"), addr2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("client send: %v->%v\n", c.LocalAddr(), addr2)

	a, _ := net.ResolveUDPAddr("udp", "34.87.83.113:31310")
	_, err = c.WriteTo([]byte("hello world 3 to SG"), a)
	if err != nil {
		panic(err)
	}
	fmt.Printf("client send: %v->%v\n", c.LocalAddr(), addr2)

	for {
		var b [65535]byte
		n, a, err := c.ReadFrom(b[:])
		if err != nil {
			fmt.Println("client read:", err)
			panic(err)
		}
		d := b[:n]
		fmt.Printf("client recv: %v->%v: %v\n", a, c.LocalAddr(), string(d))
	}
}

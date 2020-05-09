package main

import (
	"flag"
	"fmt"
	"net"
)

var (
	port  = flag.Int("port", 10999, "port to listen")
	port2 = flag.Int("port2", 10998, "port to listen")
)

func main() {
	flag.Parse()
	a := &net.UDPAddr{Port: *port}
	b := &net.UDPAddr{Port: *port2}
	go serve(b)
	serve(a)
}

func serve(a *net.UDPAddr) {
	c, e := net.ListenUDP("udp", a)
	if e != nil {
		panic(e)
	}
	defer c.Close()
	fmt.Printf("listing on %v\n", c.LocalAddr())
	for {
		var b [65535]byte
		n, addr, err := c.ReadFrom(b[:])
		if err != nil {
			fmt.Println("server read:", err)
			panic(err)
		}
		d := b[:n]
		fmt.Printf("server recv: %v->%v: %v\n", addr, c.LocalAddr(), string(d))

		_, err = c.WriteTo(d, addr)
		if err != nil {
			panic(err)
		}
		fmt.Printf("1server send: %v->%v: %v\n", c.LocalAddr(), addr, string(d))

		pc, err := net.ListenUDP("udp", nil)
		if err != nil {
			panic(err)
		}
		fmt.Printf("server listening on: %v\n", pc.LocalAddr())
		_, err = pc.WriteTo([]byte("hello again world"), addr)
		if err != nil {
			panic(err)
		}
		fmt.Printf("2server send: %v->%v: %v\n", pc.LocalAddr(), addr, string(d))
	}
}

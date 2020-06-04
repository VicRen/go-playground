package main

import (
	"flag"
	"fmt"
	"net"
)

var (
	s    = flag.Bool("s", false, "is server")
	port = flag.String("p", "8081", "port")
	host = flag.String("h", "127.0.0.1", "host")
)

func main() {
	flag.Parse()

	addr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(*host, *port))
	if err != nil {
		panic(err)
	}
	if *s {
		runServer(addr)
	} else {
		runClient(addr)
	}
}

func runServer(addr *net.UDPAddr) {
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	fmt.Println("listening on", addr)
	buf := make([]byte, 1024)
	for {
		n, remote, err := conn.ReadFrom(buf)
		if err != nil {
			panic(err)
		}
		fmt.Printf("reading data from [%v]: %v\n", remote, buf[:n])
		_, err = conn.WriteTo(buf[:n], remote)
		if err != nil {
			panic(err)
		}
		fmt.Println("write data back")
	}
}

func runClient(addr *net.UDPAddr) {
	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	go func() {
		_, err = conn.WriteTo([]byte("Hello world"), addr)
		if err != nil {
			panic(err)
		}
	}()
	buf := make([]byte, 1024)
	for {
		n, remote, err := conn.ReadFrom(buf)
		if err != nil {
			panic(err)
		}
		fmt.Printf("reading data from [%v]: %v\n", remote, buf[:n])
	}
}

package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	//go func() {
	//	log.Fatal(echoServer())
	//}()

	log.Fatal(client())
}

func client() error {
	c, e := net.Dial("udp", "39.96.21.158:6060")
	if e != nil {
		return e
	}

	_, err := c.Write([]byte("testing"))
	if err != nil {
		return err
	}
	fmt.Printf("client send: %v->%v\n", c.LocalAddr(), c.RemoteAddr())

	for {
		var b [65535]byte
		n, err := c.Read(b[:])
		if err != nil {
			fmt.Println("client read:", err)
			return err
		}
		d := b[:n]
		fmt.Printf("client recv: %v->%v: %v\n", c.RemoteAddr(), c.LocalAddr(), d)
	}
}

func echoServer() error {
	c, e := net.ListenUDP("udp", &net.UDPAddr{Port: 6060})
	if e != nil {
		return e
	}
	for {
		var b [65535]byte
		n, addr, err := c.ReadFrom(b[:])
		if err != nil {
			fmt.Println("server read:", err)
			return err
		}
		d := b[:n]
		fmt.Printf("server recv: %v->%v: %v\n", addr, c.LocalAddr(), d)

		//conn, err := net.Dial("udp", addr.String())
		//if err != nil {
		//	return err
		//}
		_, err = c.WriteTo(d, addr)
		if err != nil {
			return err
		}
		fmt.Printf("server send: %v->%v: %v\n", c.LocalAddr(), addr, d)
	}
}

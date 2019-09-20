package main

import (
	"fmt"
	"net"
)

func main() {
	c, err := net.Dial("tcp", "39.96.21.158:8221")
	if err != nil {
		panic(err)
	}
	var message = []byte("Hello galaxy")
	_, err = c.Write(message)
	if err != nil {
		panic(fmt.Sprintf("failed to write data: %s", err))
	}
	fmt.Println("Client Sent:", message)

	var b [2048]byte
	n, err := c.Read(b[:])
	if err != nil {
		panic(err)
	}
	fmt.Println("Client Got:", b[:n])
}

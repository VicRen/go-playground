package main

import (
	"fmt"
	"net"
)

func main() {
	l, err := net.Listen("tcp", ":8221")
	if err != nil {
		panic(err)
	}

	for {
		c, err := l.Accept()
		if err != nil {
			panic(err)
		}

		var b [2048]byte
		n, err := c.Read(b[:])
		if err != nil {
			fmt.Println("failed to read data:", err)
			continue
		}
		fmt.Println("Server Got:", b[:n])

		_, err = c.Write(b[:n])
		if err != nil {
			fmt.Println("failed to write data:", err)
		}
		fmt.Println("Server Sent:", b[:n])
	}
}

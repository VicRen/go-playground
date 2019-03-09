package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	server, client := net.Pipe()

	go func() {
		for {
			fmt.Println("Client Send: Hello Galaxy")
			_, err := client.Write([]byte("Hello Galaxy"))
			if err != nil {
				fmt.Println(err)
				break
			}
			time.Sleep(1 * time.Second)
		}
	}()

	for {
		buf := make([]byte, 1024)
		n, err := server.Read(buf)
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println("Server Got: ", string(buf[:n]))
	}
}

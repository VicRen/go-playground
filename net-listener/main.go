package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	l, err := net.Listen("tcp", ":8081")
	if err != nil {
		panic(err)
	}
	go func() {
		time.Sleep(time.Second)
		_ = l.Close()
	}()
	for {
		_, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
	}

}

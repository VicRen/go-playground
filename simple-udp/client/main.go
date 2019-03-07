package main

import (
	"flag"
	"fmt"
	"net"
	"time"
)

var host = flag.String("host", "", "host")
var port = flag.String("port", "8081", "port")

func main() {
	flag.Parse()

	addr := net.JoinHostPort(*host, *port)
	for {
		time.Sleep(time.Second)
		sendData(addr)
	}
}

func sendData(addr string) {
	conn, err := net.Dial("udp", addr)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	_, err = conn.Write([]byte("Hello Galaxy"))
	if err != nil {
		fmt.Println("Client failed to write UDP msg: ", err.Error())
		return
	}
	data := make([]byte, 1024)
	n, err := conn.Read(data)
	if err != nil {
		fmt.Println("Client failed to read UDP msg: ", err.Error())
		return
	}

	fmt.Printf("Client Get \"%s\"\n", string(data[:n]))
}

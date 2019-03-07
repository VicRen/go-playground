package main

import (
	"flag"
	"fmt"
	"net"
)

var host = flag.String("host", "", "host")
var port = flag.String("port", "8081", "port")

func main() {
	flag.Parse()

	addr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(*host, *port))
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	for {
		handleClient(conn)
	}
}

func handleClient(conn *net.UDPConn) {
	data := make([]byte, 1024)
	n, remoteAddr, err := conn.ReadFromUDP(data)
	if err != nil {
		fmt.Println("Server failed to read UDP msg: ", err.Error())
		return
	}

	fmt.Printf("Server Got \"%s\" from %s\n", string(data[:n]), remoteAddr)
	_, err = conn.WriteToUDP(data[:n], remoteAddr)
	if err != nil {
		fmt.Println("Server failed to write UDP msg:", err.Error())
	}
}

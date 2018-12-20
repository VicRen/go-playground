package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	addr1 := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9981}
	addr2 := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9982}
	addr3 := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9983}

	go func() {
		listener1, err := net.ListenUDP("udp", addr1)
		if err != nil {
			fmt.Println(err)
			return
		}

		go read(listener1)

		for {
			time.Sleep(5 * time.Second)

			_, err = listener1.WriteToUDP([]byte("ping to #2: "+addr2.String()), addr2)
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	go func() {
		listener1, err := net.ListenUDP("udp", addr2)
		if err != nil {
			fmt.Println(err)
			return
		}
		go read(listener1)
		time.Sleep(5 * time.Second)

		_, err = listener1.WriteToUDP([]byte("ping to #1: " + addr1.String()), addr2)
		if err != nil {
			fmt.Println(err)
		}
	}()
	listener2 , err := net.ListenUDP("udp", addr3)
	if err != nil {
		fmt.Println(err)
		return
	}
	read(listener2)
}

func read(conn *net.UDPConn) {
	for {
		data := make([]byte, 1024)
		n, removeAddr, err := conn.ReadFromUDP(data)
		if err != nil {
			fmt.Printf("error during read: %s", err)
		}
		fmt.Printf("receive %s from <%s>\n", data[:n], removeAddr)
	}
}

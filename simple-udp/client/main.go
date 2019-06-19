package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"time"
)

var host = flag.String("host", "47.101.180.117", "host")
var port = flag.String("port", "6062", "port")

func main() {
	flag.Parse()

	addr := net.JoinHostPort(*host, *port)
	toAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		panic(err)
	}
	conn := newUDPConn()
	count := 0
	for {
		if count > 5 {
			conn.Close()
			conn = newUDPConn()
		}
		time.Sleep(time.Second)
		sendData(conn, toAddr, count)
		count++
	}
}

func sendData(conn *net.UDPConn, addr net.Addr, count int) {
	_, err := conn.WriteTo([]byte("Hello Galaxy:"+strconv.Itoa(count)), addr)
	fmt.Printf("Client send msg to %v from %v\n", addr, conn.LocalAddr())
	if err != nil {
		fmt.Println("Client failed to write UDP msg: ", err.Error())
		return
	}
	ch := make(chan []byte)
	go recvData(conn, ch)

	select {
	case <-time.After(1 * time.Second):
		fmt.Println("read timeout")
	case d := <-ch:
		fmt.Printf("Client Get \"%s\"\n", string(d))
	}
}

func recvData(conn net.Conn, dataCh chan []byte) {
	defer close(dataCh)
	data := make([]byte, 1024)
	n, err := conn.Read(data)
	if err != nil {
		fmt.Println("Client failed to read UDP msg: ", err.Error())
		return
	}
	dataCh <- data[:n]
}

func newUDPConn() *net.UDPConn {
	udpAddr, err := net.ResolveUDPAddr("udp", "")
	if err != nil {
		panic(err)
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		panic(err)
	}
	fmt.Println("listening on addr:", conn.LocalAddr())
	return conn
}

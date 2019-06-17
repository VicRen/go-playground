package main

import (
	"log"
	"net"
	"time"
)

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", ":8080")
	if err != nil {
		panic(err)
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		panic(err)
	}
	defer func() {
		conn.Close()
		log.Println("closing listen")
	}()

	data := byte(0x00)
	go func() {
		for {
			time.Sleep(1 * time.Second)
			conn, err := net.Dial("udp", ":8080")
			if err != nil {
				panic(err)
			}
			log.Println("writing data")
			conn.Write([]byte{data})
			conn.Close()
			data++
		}
	}()

	count := 1
	for {
		buf := make([]byte, 1024)
		n, src, err := conn.ReadFrom(buf)
		if err != nil {
			panic(err)
		}
		log.Println("data read:", n, src, buf[:n], conn.RemoteAddr())
		if count == 3 {
			break
		}
		count++
	}
	time.Sleep(3 * time.Second)
	count = 1
	for {
		buf := make([]byte, 1024)
		n, src, err := conn.ReadFrom(buf)
		if err != nil {
			panic(err)
		}
		log.Println("data read2:", n, src, buf[:n], conn.RemoteAddr())
		if count == 3 {
			break
		}
		count++
	}
	time.Sleep(3 * time.Second)
}

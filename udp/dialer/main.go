package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"strings"
)

var (
	target = flag.String("t", "", "Target to dial.")
)

func main() {
	flag.Parse()

	if *target == "" {
		panic("target is empty")
	}

	addr, err := net.ResolveUDPAddr("udp", *target)
	if err != nil {
		panic(fmt.Sprintf("failed to resolve addr: %s", addr))
	}

	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		panic(fmt.Sprintf("failed to dial addr: %s", addr))
	}

	go func() {
		var str string

		for i := 0; i < 10000; i++ {
			str = str + string(rand.Intn(120))
		}

		var r io.Reader = strings.NewReader(str)

		wt, n, err := copyString(addr, conn, r)
		if err != nil {
			log.Fatal("error sending data: %s", err)
		}

		log.Printf("sent %d bytes, times %d", n, wt)
	}()

	buf := make([]byte, 2048)
	count := 0
	for {
		n, err := conn.Read(buf)
		count += n
		if err != nil {
			fmt.Println("------>write error:", err)
			break
		}
	}
	log.Println("------>read count:", count)

	_ = conn.Close()
}

func copyString(addr net.Addr, dst net.PacketConn, src io.Reader) (int, int, error) {
	buf := make([]byte, 1024)
	n := 0
	wt := 0
	for {
		nr, err := src.Read(buf)
		if err != nil && err != io.EOF {
			return wt, n, err
		}
		nw, errw := dst.WriteTo(buf[:nr], addr)
		if errw != nil {
			log.Println("write err:", err)
			return wt, n, errw
		}
		n += nw
		wt++
		if err == io.EOF {
			break
		}
	}
	return wt, n, nil
}

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/VicRen/go-play-ground/socks-quic/socks"
	"github.com/VicRen/go-play-ground/socks-quic/transport/internet/quic"
)

var quicDialer = func(ctx context.Context, network, addr string) (net.Conn, error) {
	ip, ports, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	portI, err := strconv.Atoi(ports)
	if err != nil {
		return nil, err
	}
	dstAddr := &socks.AddrSpec{IP: net.ParseIP(ip), Port: portI}
	dstAddrBuf, err := dstAddr.Bytes()
	if err != nil {
		return nil, err
	}
	header := append([]byte{0x73}, dstAddrBuf...)

	proxyPortI, err := strconv.Atoi(*port)
	conn, err := quic.Dial(context.Background(), &net.UDPAddr{IP: net.ParseIP(*host), Port: proxyPortI})
	if err != nil {
		return nil, err
	}

	fmt.Println("proxying to", addr)
	_, err = conn.Write(header)

	var b [1]byte
	// shake timeout
	err = conn.SetDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		return nil, errors.New("failed to set handshake timeout")
	}
	_, err = conn.Read(b[:])
	if err != nil {
		return nil, err
	}
	if b[0] != byte(0x00) {
		return nil, errors.New("failed to handshake")
	}
	// cancel handle shake timeout
	err = conn.SetDeadline(time.Time{})
	if err != nil {
		return nil, errors.New("failed to set handshake timeout")
	}
	fmt.Println("handshake ok, start proxying")
	return conn, nil
}

var host = flag.String("host", "192.168.50.7", "host")
var port = flag.String("port", "8221", "port")

func main() {
	// Create a SOCKS5 server
	conf := &socks.Config{
		Dial: quicDialer,
	}
	server, err := socks.NewSOCKSServer(conf)
	if err != nil {
		panic(err)
	}

	// Create SOCKS5 proxy on localhost port 10080
	if err := server.Serve("127.0.0.1:10080"); err != nil {
		panic(err)
	}
}

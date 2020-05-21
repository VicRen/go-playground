package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"

	addr "github.com/VicRen/go-playground/addr-util"
)

func main() {
	udpCh := make(chan UDPItem, 50000)
	addr := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 6060}
	c, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	go echoSever()
	go UDPConnDaemon(udpCh)
	go func() {
		for {
			buf := make([]byte, 2048)
			n, srcAddr, err := c.ReadFromUDP(buf)
			if err != nil {
				fmt.Println("error reading UDP:", err)
				return
			}
			packet := buf[:n]
			fmt.Println("data read:", packet)
			udpItem := UDPItem{
				packet:    &packet,
				localAddr: addr,
				srcAddr:   srcAddr,
			}
			udpCh <- udpItem
		}
	}()
	time.Sleep(1 * time.Second)
	UDPSender()
	time.Sleep(1 * time.Second)
}

type UDPItem struct {
	packet    *[]byte
	localAddr *net.UDPAddr
	srcAddr   *net.UDPAddr
}

func echoSever() {
	udpAddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:6000")
	c, e := net.ListenUDP("udp", udpAddr)
	if e != nil {
		panic(e)
	}
	n, err := io.Copy(c, c)
	if err != nil && err != io.EOF {
		panic(err)
	}
	fmt.Println("echoServer finished:", n)
}

func UDPConnDaemon(udpCh chan UDPItem) {
	for {
		item := <-udpCh
		fmt.Printf("udp item recv: %v -> %v: %v\n", item.srcAddr, item.localAddr, *item.packet)
		bufReader := bytes.NewReader(*item.packet)
		var addrLength uint16
		var bodyLength uint16
		err := binary.Read(bufReader, binary.BigEndian, &addrLength)
		if err != nil {
			fmt.Println("error reading udp item:", err)
			break
		}
		_srcAddr := make([]byte, addrLength)
		n, err := bufReader.Read(_srcAddr)
		if err != nil {
			fmt.Println("error reading udp item:", err)
			break
		}
		if n != int(addrLength) {
			break
		}

		_, i, p, err := addr.UnmarshalAddr(_srcAddr)
		if err != nil {
			fmt.Println("error unmarshal addr:", err)
			break
		}

		err = binary.Read(bufReader, binary.BigEndian, &bodyLength)
		if err != nil {
			fmt.Println("error reading udp item:", err)
			break
		}
		packet := make([]byte, bodyLength)
		n, err = bufReader.Read(packet)
		if err != nil {
			break
		}
		if n != int(bodyLength) {
			break
		}
		fmt.Println("udp packet body:", packet, "from:", net.JoinHostPort(i.String(), strconv.Itoa(int(p))))

		outConn, err := GetOutConn()
		if err != nil {
			panic(err)
		}
		writer := bufio.NewWriter(outConn)
		dd, e := PackUDP(item.srcAddr.String(), *item.packet)
		if e != nil {
			break
		}
		writer.Write(dd)
		writer.Flush()
	}
}

func UDPSender() {
	laddr := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 6061}
	addr := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 6060}
	c, err := net.DialUDP("udp", laddr, addr)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			buf := make([]byte, 1024)
			n, err := c.Read(buf)
			if err != nil {
				continue
			}
			fmt.Println("\n\nsender read:", buf[:n])
		}

	}()

	count := 0
	for {
		d, e := PackUDP(laddr.String(), []byte{uint8(count)})
		if e != nil {
			continue
		}
		fmt.Println(count, "written")
		w := bufio.NewWriter(c)
		w.Write(d)
		err := w.Flush()
		if err != nil {
			break
		}
		if count > 4 {
			break
		}
		count++
	}

}

func PackUDP(srcAddr string, packet []byte) ([]byte, error) {
	addrBuf, err := addr.MarshalAddr(srcAddr)
	if err != nil {
		return nil, err
	}
	addrLen := uint16(len(addrBuf))
	bodyLen := uint16(len(packet))
	pkg := new(bytes.Buffer)
	binary.Write(pkg, binary.BigEndian, addrLen)
	binary.Write(pkg, binary.BigEndian, addrBuf)
	binary.Write(pkg, binary.BigEndian, bodyLen)
	binary.Write(pkg, binary.BigEndian, packet)
	return pkg.Bytes(), nil
}

var conn net.Conn

func GetOutConn() (outConn net.Conn, err error) {
	if conn == nil {
		conn, err = net.Dial("udp", "127.0.0.1:6000")
		if err == nil {
			outConn = conn
			return
		}
	}
	outConn = conn
	return
}

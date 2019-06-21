package addr

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strconv"
)

const (
	IPv4Address = uint8(1)
	fqdnAddress = uint8(3)
	IPv6Address = uint8(4)
)

// MarshalAddr addr convert to byte
func MarshalAddr(addr string) ([]byte, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	ip := net.ParseIP(host)
	portI, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}
	var addrType uint8
	var addrBody []byte
	var addrPort uint16
	switch {
	case ip.To4() != nil:
		addrType = IPv4Address
		addrBody = []byte(ip.To4())
		addrPort = uint16(portI)
	case ip.To16() != nil:
		addrType = IPv6Address
		addrBody = []byte(ip.To16())
		addrPort = uint16(portI)
	default:
		return nil, fmt.Errorf("failed to formate address: %v", addr)
	}
	ret := make([]byte, len(addrBody)+3)
	ret[0] = addrType
	copy(ret[1:], addrBody)
	ret[1+len(addrBody)] = byte(addrPort >> 8)
	ret[1+len(addrBody)+1] = byte(addrPort)
	return ret, nil
}

func UnmarshalAddr(addrBuf []byte) (addrType uint8, ip net.IP, port uint16, err error) {
	if uint8(len(addrBuf)) < 7 {
		err = errors.New("some error")
	}
	addrType = addrBuf[0]
	ipLen := uint8(len(addrBuf)) - 2 - 1
	ipBuf := addrBuf[1 : ipLen+1]
	portBuf := addrBuf[ipLen+1:]
	ip = net.IP(ipBuf)
	port = binary.BigEndian.Uint16(portBuf)
	return
}

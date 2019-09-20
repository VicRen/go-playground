package socks

import (
	"fmt"
	"io"
	"net"
	"strconv"
)

// AddrSpec is used to return the target AddrSpec
// which may be specified as IPv4, IPv6, or a FQDN
type AddrSpec struct {
	FQDN string
	IP   net.IP
	Port int
}

func (a *AddrSpec) String() string {
	if a.FQDN != "" {
		return fmt.Sprintf("%s (%s):%d", a.FQDN, a.IP, a.Port)
	}
	return fmt.Sprintf("%s:%d", a.IP, a.Port)
}

// Address returns a string suitable to dial; prefer returning IP-based
// address, fallback to FQDN
func (a AddrSpec) Address() string {
	if 0 != len(a.IP) {
		return net.JoinHostPort(a.IP.String(), strconv.Itoa(a.Port))
	}
	return net.JoinHostPort(a.FQDN, strconv.Itoa(a.Port))
}

// ReadAddrSpec is used to read AddrSpec.
// Expects an address type byte, followed by the address and port
func ReadAddrSpec(r io.Reader) (int, *AddrSpec, error) {
	addrLen := 0
	d := &AddrSpec{}

	// Get the address type
	addrType := make([]byte, lenAddrType)
	if _, err := io.ReadFull(r, addrType); err != nil {
		return 0, nil, err
	}
	// Handle on a per type basis
	switch addrType[0] {
	case ipv4Address:
		addrLen = int(lenIPv4)
		addr := make([]byte, lenIPv4)
		if _, err := io.ReadFull(r, addr); err != nil {
			return 0, nil, err
		}
		d.IP = net.IP(addr)
	case ipv6Address:
		addrLen = int(lenIPv6)
		addr := make([]byte, lenIPv6)
		if _, err := io.ReadFull(r, addr); err != nil {
			return 0, nil, err
		}
		d.IP = net.IP(addr)
	case fqdnAddress:
		if _, err := io.ReadFull(r, addrType); err != nil {
			return 0, nil, err
		}
		addrLen = int(addrType[0])
		domainName := make([]byte, addrLen)
		if _, err := io.ReadFull(r, domainName); err != nil {
			return 0, nil, err
		}
		d.FQDN = string(domainName)
		// when addr is domain name, one byte is the length of addr.
		addrLen += int(lenAddrType)
	default:
		return 0, nil, errUnrecognizedAddrType
	}
	port := make([]byte, lenPort)
	if _, err := io.ReadFull(r, port); err != nil {
		return 0, nil, err
	}
	d.Port = (int(port[0]) << 8) | int(port[1])
	return addrLen + 2, d, nil
}

func (a *AddrSpec) Bytes() ([]byte, error) {
	var addrType uint8
	var addrBody []byte
	var addrPort uint16
	switch {
	case a.FQDN != "":
		addrType = fqdnAddress
		addrBody = append([]byte{byte(len(a.FQDN))}, a.FQDN...)
		addrPort = uint16(a.Port)
	case a.IP.To4() != nil:
		addrType = ipv4Address
		addrBody = []byte(a.IP.To4())
		addrPort = uint16(a.Port)
	case a.IP.To16() != nil:
		addrType = ipv6Address
		addrBody = []byte(a.IP.To16())
		addrPort = uint16(a.Port)
	default:
		return nil, fmt.Errorf(errFormatAddrFmt, a)
	}
	ret := make([]byte, len(addrBody)+3)
	ret[0] = addrType
	copy(ret[1:], addrBody)
	ret[1+len(addrBody)] = byte(addrPort >> 8)
	ret[1+len(addrBody)+1] = byte(addrPort & 0xff)
	return ret, nil
}

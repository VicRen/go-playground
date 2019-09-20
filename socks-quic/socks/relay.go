package socks

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	cache "github.com/patrickmn/go-cache"
)

type Relay struct {
	flag byte
	// UDPAddr of the the network that sent the request
	RemoteAddr *net.UDPAddr
	// AddrSpec of the desired destination
	DstAddr *AddrSpec
	// Raw data of this relay
	Packet []byte
}

const (
	errBadSrc        = "bad source"
	errBadRequestFmt = "bad request: %v"
	errIgnoreFlagFmt = "ignore flag: %d"
)

var errTooShort = errors.New("too short")

func NewRelay(srcAddr *net.UDPAddr, data []byte) (*Relay, error) {
	if srcAddr == nil {
		return nil, fmt.Errorf(errBadRequestFmt, errBadSrc)
	}
	addr, d, err := newDatagramFromBytes(data)
	if err != nil {
		return nil, fmt.Errorf(errBadRequestFmt, err)
	}

	return &Relay{
		flag:       d.frag,
		RemoteAddr: srcAddr,
		DstAddr:    addr,
		Packet:     d.data,
	}, nil
}

func (s *Server) handleRelay(relay *Relay) error {
	if relay.flag != 0x00 {
		return fmt.Errorf(errIgnoreFlagFmt, relay.flag)
	}
	// TODO debug logger
	send := func(ue *udpExchange, data []byte) error {
		_, err := ue.RemoteConn.Write(data)
		if err != nil {
			return err
		}
		return nil
	}
	var ue *udpExchange
	iue, ok := s.udpExchangeCache.Get(relay.RemoteAddr.String())
	if ok {
		ue = iue.(*udpExchange)
		return send(ue, relay.Packet)
	}

	c, err := s.config.Dial(context.Background(), "udp", relay.DstAddr.Address())
	if err != nil {
		s.deleteAssociate(relay.DstAddr.String())
		return err
	}
	// A UDP association terminates when the TCP connection that the UDP
	// ASSOCIATE request arrived on terminates.
	ue = &udpExchange{
		ClientAddr: relay.RemoteAddr,
		DstAddr:    relay.DstAddr,
		RemoteConn: c,
	}
	if err := send(ue, relay.Packet); err != nil {
		s.deleteAssociate(relay.DstAddr.String())
		_ = ue.RemoteConn.Close()
		return err
	}
	s.udpExchangeCache.Set(ue.ClientAddr.String(), ue, cache.DefaultExpiration)
	go func(ue *udpExchange) {
		defer func() {
			s.deleteAssociate(ue.DstAddr.String())
			s.udpExchangeCache.Delete(ue.ClientAddr.String())
			_ = ue.RemoteConn.Close()
		}()
		b := make([]byte, 65535)
		for {
			if s.config.UDPDeadline != 0 {
				if err := ue.RemoteConn.SetDeadline(time.Now().Add(time.Duration(s.config.UDPDeadline) * time.Second)); err != nil {
					log.Println(err)
					break
				}
			}
			n, err := ue.RemoteConn.Read(b)
			if n > 0 {
				packet := b[:n]
				d, err := newDatagram(ue.DstAddr, packet)
				if err != nil {
					log.Println("udp exchange: create datagram", err)
					break
				}
				if _, err := s.udpConn.WriteToUDP(d.Bytes(), ue.ClientAddr); err != nil {
					log.Println("udp exchange: write datagram to client", err)
					break
				}
			}
			if err != nil && err != io.EOF {
				if e, ok := err.(timeout); ok && e.Timeout() {
					// timeout error will not be logger
					break
				}
				log.Println("udp exchange: read datagram from remote", err)
				break
			}
		}
	}(ue)
	return nil
}

// TODO move to common
type timeout interface {
	Timeout() bool
}

func (s *Server) deleteAssociate(addr string) {
	v, ok := s.associateCache.Get(addr)
	if ok {
		ch := v.(chan struct{})
		ch <- struct{}{}
		s.associateCache.Delete(addr)
	}
}

type datagram struct {
	rsv  []byte
	frag byte
	atyp byte
	addr []byte
	port []byte
	data []byte
}

func newDatagramFromBytes(b []byte) (*AddrSpec, *datagram, error) {
	n := len(b)
	if n < minDatagramLen {
		return nil, nil, errTooShort
	}
	bufReader := bytes.NewReader(b)

	rsvBuf := make([]byte, lenAssociateRSV)
	if _, err := io.ReadFull(bufReader, rsvBuf); err != nil {
		return nil, nil, errTooShort
	}

	flag, err := bufReader.ReadByte()
	if err != nil {
		return nil, nil, err
	}

	addrLen, dstAddr, err := ReadAddrSpec(bufReader)
	if err != nil {
		return nil, nil, err
	}

	addrBuf, err := dstAddr.Bytes()
	if err != nil {
		return nil, nil, err
	}

	return dstAddr, &datagram{
		rsv:  rsvBuf,
		frag: flag,
		atyp: addrBuf[0],
		addr: addrBuf[lenAddrType : len(addrBuf)-int(lenPort)],
		port: addrBuf[lenPort:],
		data: b[int(lenAssociateRSV+lenAssociateFRAG+lenAddrType)+addrLen:],
	}, nil
}

func newDatagram(addr *AddrSpec, data []byte) (*datagram, error) {
	addrBuf, err := addr.Bytes()
	if err != nil {
		return nil, err
	}
	datagram := &datagram{
		rsv:  []byte{0x00, 0x00},
		frag: 0x00,
		atyp: addrBuf[0],
		addr: addrBuf[lenAddrType : len(addrBuf)-int(lenPort)],
		port: addrBuf[len(addrBuf)-int(lenPort):],
		data: data,
	}
	return datagram, nil
}

// Bytes return []byte
func (d *datagram) Bytes() []byte {
	b := make([]byte, 0)
	b = append(b, d.rsv...)
	b = append(b, d.frag)
	b = append(b, d.atyp)
	b = append(b, d.addr...)
	b = append(b, d.port...)
	b = append(b, d.data...)
	return b
}

type udpExchange struct {
	ClientAddr *net.UDPAddr
	DstAddr    *AddrSpec
	RemoteConn net.Conn
}

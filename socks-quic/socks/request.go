package socks

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	cache "github.com/patrickmn/go-cache"
)

const errFormatAddrFmt = "failed to format address: %v"

var (
	errUnrecognizedAddrType = errors.New("unrecognized address type")
)

// A Request represents request received by a server
type Request struct {
	// Protocol version
	Version uint8
	// Requested command
	Command uint8
	// AuthContext provided during negotiation
	AuthContext *AuthContext
	// AddrSpec of the the network that sent the request
	RemoteAddr *AddrSpec
	// AddrSpec of the desired destination
	DstAddr *AddrSpec
	bufConn io.Reader
}

type conn interface {
	Write([]byte) (int, error)
	RemoteAddr() net.Addr
}

// NewRequest creates a new Request from the tcp connection
func NewRequest(bufConn io.Reader) (*Request, error) {
	// Read the version byte
	header := []byte{0, 0, 0}
	if _, err := io.ReadAtLeast(bufConn, header, 3); err != nil {
		return nil, fmt.Errorf("failed to get command version: %v", err)
	}

	// Ensure we are compatible
	if header[0] != socks5Version {
		return nil, fmt.Errorf("unsupported command version: %v", header[0])
	}

	// Read in the destination address
	_, dest, err := ReadAddrSpec(bufConn)
	if err != nil {
		return nil, err
	}

	request := &Request{
		Version: socks5Version,
		Command: header[1],
		DstAddr: dest,
		bufConn: bufConn,
	}

	return request, nil
}

// handleRequest is used for request processing after authentication
func (s *Server) handleRequest(req *Request, conn conn) error {
	ctx := context.Background()

	// Resolve the address if we have a FQDN
	dest := req.DstAddr
	if dest.FQDN != "" {
		ctx_, addr, err := s.config.Resolver.Resolve(ctx, dest.FQDN)
		if err != nil {
			if err := sendReply(conn, hostUnreachable, nil); err != nil {
				return fmt.Errorf("failed to send reply: %v", err)
			}
			return fmt.Errorf("failed to resolve destination '%v': %v", dest.FQDN, err)
		}
		ctx = ctx_
		dest.IP = addr
	}

	// Switch on the command
	switch req.Command {
	case ConnectCommand:
		return s.handleConnect(ctx, conn, req)
	case BindCommand:
		return s.handleBind(ctx, conn, req)
	case AssociateCommand:
		return s.handleAssociate(ctx, conn, req)
	default:
		if err := sendReply(conn, commandNotSupported, nil); err != nil {
			return fmt.Errorf("failed to send reply: %v", err)
		}
		return fmt.Errorf("unsupported command: %v", req.Command)
	}
}

// handleConnect is used to handle a connect command
func (s *Server) handleConnect(ctx context.Context, conn conn, req *Request) error {
	// Check if this is allowed
	if ctx_, ok := s.config.Rules.Allow(ctx, req); !ok {
		if err := sendReply(conn, ruleFailure, nil); err != nil {
			return fmt.Errorf("failed to send reply: %v", err)
		}
		return fmt.Errorf("connect to %v blocked by rules", req.DstAddr)
	} else {
		ctx = ctx_
	}
	// Attempt to connect
	target, err := s.config.Dial(ctx, "tcp", req.DstAddr.Address())
	if err != nil {
		msg := err.Error()
		resp := hostUnreachable
		if strings.Contains(msg, "refused") {
			resp = connectionRefused
		} else if strings.Contains(msg, "network is unreachable") {
			resp = networkUnreachable
		}
		if err := sendReply(conn, resp, nil); err != nil {
			return fmt.Errorf("failed to send reply: %v", err)
		}
		return fmt.Errorf("connect to %v failed: %v", req.DstAddr, err)
	}
	defer target.Close()
	if s.config.TCPTimeout != 0 {
		if err := target.(*net.TCPConn).SetKeepAlivePeriod(time.Duration(s.config.TCPTimeout) * time.Second); err != nil {
			return fmt.Errorf("failed to set tcp timeout: %v", err)
		}
	}
	if s.config.TCPDeadline != 0 {
		if err := target.SetDeadline(time.Now().Add(time.Duration(s.config.TCPDeadline) * time.Second)); err != nil {
			return fmt.Errorf("failed to set tcp deadline: %v", err)
		}
	}

	// Send success
	local := target.LocalAddr().(*net.TCPAddr)
	bind := AddrSpec{IP: local.IP, Port: local.Port}
	if err := sendReply(conn, successReply, &bind); err != nil {
		return fmt.Errorf("failed to send reply: %v", err)
	}

	// Start proxying
	errCh := make(chan error, 2)
	go proxy(target, req.bufConn, errCh)
	go proxy(conn, target, errCh)

	// Wait
	for i := 0; i < 2; i++ {
		e := <-errCh
		if e != nil {
			// return from this function closes target (and conn).
			return e
		}
	}
	return nil
}

// handleBind is used to handle a connect command
func (s *Server) handleBind(ctx context.Context, conn conn, req *Request) error {
	// Check if this is allowed
	if ctx_, ok := s.config.Rules.Allow(ctx, req); !ok {
		if err := sendReply(conn, ruleFailure, nil); err != nil {
			return fmt.Errorf("failed to send reply: %v", err)
		}
		return fmt.Errorf("bind to %v blocked by rules", req.DstAddr)
	} else {
		ctx = ctx_
	}

	// TODO: Support bind
	if err := sendReply(conn, commandNotSupported, nil); err != nil {
		return fmt.Errorf("failed to send reply: %v", err)
	}
	return nil
}

// handleAssociate is used to handle a connect command
func (s *Server) handleAssociate(ctx context.Context, conn conn, req *Request) error {
	// Check if this is allowed
	if ctx_, ok := s.config.Rules.Allow(ctx, req); !ok {
		if err := sendReply(conn, ruleFailure, nil); err != nil {
			return fmt.Errorf("failed to send reply: %v", err)
		}
		return fmt.Errorf("associate to %v blocked by rules", req.DstAddr)
	} else {
		ctx = ctx_
	}

	if err := sendReply(conn, successReply, s.udpBindAddr); err != nil {
		return fmt.Errorf("failed to send reply: %v", err)
	}

	_, p, err := net.SplitHostPort(req.DstAddr.Address())
	if err != nil {
		return err
	}
	if p == "0" {
		time.Sleep(time.Duration(s.config.UDPSessionTime) * time.Second)
		return nil
	}

	ch := make(chan struct{})
	s.associateCache.Set(req.DstAddr.Address(), ch, cache.DefaultExpiration)
	<-ch
	return nil
}

// sendReply is used to send a reply message
func sendReply(w io.Writer, resp uint8, addr *AddrSpec) error {
	// Format the address
	var addrType uint8
	var addrBody []byte
	var addrPort uint16
	switch {
	case addr == nil:
		addrType = ipv4Address
		addrBody = []byte{0, 0, 0, 0}
		addrPort = 0

	case addr.FQDN != "":
		addrType = fqdnAddress
		addrBody = append([]byte{byte(len(addr.FQDN))}, addr.FQDN...)
		addrPort = uint16(addr.Port)

	case addr.IP.To4() != nil:
		addrType = ipv4Address
		addrBody = []byte(addr.IP.To4())
		addrPort = uint16(addr.Port)

	case addr.IP.To16() != nil:
		addrType = ipv6Address
		addrBody = []byte(addr.IP.To16())
		addrPort = uint16(addr.Port)

	default:
		return fmt.Errorf("failed to format address: %v", addr)
	}

	// Format the message
	msg := make([]byte, 6+len(addrBody))
	msg[0] = socks5Version
	msg[1] = resp
	msg[2] = 0 // Reserved
	msg[3] = addrType
	copy(msg[4:], addrBody)
	msg[4+len(addrBody)] = byte(addrPort >> 8)
	msg[4+len(addrBody)+1] = byte(addrPort & 0xff)

	// Send the message
	_, err := w.Write(msg)
	return err
}

type closeWriter interface {
	CloseWrite() error
}

// proxy is used to shuffle data from src to destination, and sends errors
// down a dedicated channel
func proxy(dst io.Writer, src io.Reader, errCh chan error) {
	_, err := io.Copy(dst, src)
	if tcpConn, ok := dst.(closeWriter); ok {
		_ = tcpConn.CloseWrite()
	}
	errCh <- err
}

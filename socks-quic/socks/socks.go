package socks

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"

	cache "github.com/patrickmn/go-cache"
)

// Config is used to setup and configure a SOCKS server.
type Config struct {
	// AuthMethods can be provided to implement custom authentication
	// By default, "auth-less" mode is enabled.
	// For password-based auth use UserPassAuthenticator.
	AuthMethods []Authenticator

	// If provided, username/password authentication is enabled,
	// by appending a UserPassAuthenticator to AuthMethods. If not provided,
	// and AUthMethods is nil, then "auth-less" mode is enabled.
	Credentials CredentialStore

	// Resolver can be provided to do custom name resolution.
	// Defaults to DNSResolver if not provided.
	Resolver NameResolver

	// Rules is provided to enable custom logic around permitting
	// various commands. If not provided, PermitAll is used.
	Rules RuleSet

	// BindIP is used for bind or udp associate
	BindIP net.IP

	// TCPTimeout is the keepalive period of a request.
	TCPTimeout int

	// TCPDeadline is the deadline of a request.
	TCPDeadline int

	// UDPDeadline is the deadline of a udp associate.
	UDPDeadline int

	// UDPSessionTime is the fixed time for udp associate which client send all zero address.
	UDPSessionTime int

	// Optional function for dialing out
	Dial func(ctx context.Context, network, addr string) (net.Conn, error)
}

// Server is responsible for accepting connections and handling
// the details of the SOCKS5 protocol
type Server struct {
	config      *Config
	authMethods map[uint8]Authenticator

	udpBindAddr *AddrSpec
	udpConn     *net.UDPConn

	associateCache   *cache.Cache
	udpExchangeCache *cache.Cache
}

// DefaultDialer is the default Dial func for socks-quic.
var DefaultDialer = func(ctx context.Context, network, addr string) (net.Conn, error) {
	if network == "tcp" {
		return net.Dial("tcp", addr)
	} else if network == "udp" {
		udpAddr, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			return nil, err
		}
		c, err := net.ListenUDP("udp", nil)
		if err != nil {
			return nil, err
		}
		return &packetConn{c, udpAddr}, nil
	}
	return nil, errors.New("unsupported network")
}

type packetConn struct {
	*net.UDPConn
	udpAddr *net.UDPAddr
}

func (c *packetConn) Read(p []byte) (n int, err error) {
	n, _, err = c.ReadFrom(p)
	return
}

func (c *packetConn) Write(p []byte) (n int, err error) {
	return c.WriteTo(p, c.udpAddr)
}

func (c *packetConn) RemoteAddr() net.Addr {
	return c.udpAddr
}

func NewSOCKSServer(config *Config) (*Server, error) {
	// Ensure we have at least one authentication method enabled
	if len(config.AuthMethods) == 0 {
		if config.Credentials != nil {
			config.AuthMethods = []Authenticator{&UserPassAuthenticator{config.Credentials}}
		} else {
			config.AuthMethods = []Authenticator{&NoAuthAuthenticator{}}
		}
	}

	// Ensure we have a DNS resolver
	if config.Resolver == nil {
		config.Resolver = DNSResolver{}
	}

	// Ensure we have a rule set
	if config.Rules == nil {
		config.Rules = PermitAll()
	}

	if config.BindIP == nil {
		config.BindIP = net.ParseIP("127.0.0.1")
	}

	// set a default udp deadline, in case of UDP associate not release
	if config.UDPDeadline == 0 {
		config.UDPDeadline = 30
	}

	if config.Dial == nil {
		config.Dial = DefaultDialer
	}

	server := &Server{
		config: config,
	}
	server.authMethods = make(map[uint8]Authenticator)
	server.associateCache = cache.New(cache.NoExpiration, cache.NoExpiration)
	server.udpExchangeCache = cache.New(cache.NoExpiration, cache.NoExpiration)

	for _, a := range config.AuthMethods {
		server.authMethods[a.GetCode()] = a
	}
	return server, nil
}

// Serve is used to create a listener and serve on it
func (s *Server) Serve(addr string) error {
	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}
	portI, err := strconv.Atoi(port)
	s.udpBindAddr = &AddrSpec{IP: s.config.BindIP, Port: portI}
	a := s.udpBindAddr.Address()
	errCh := make(chan error)
	go func() {
		errCh <- s.ListenAndServeTCP(a)
	}()
	go func() {
		errCh <- s.ListenAndServeUDP(a)
	}()
	return <-errCh
}

// Serve is used to serve connections from a listener
func (s *Server) ListenAndServeTCP(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go s.ServeTCPConn(conn)
	}
}

func (s *Server) ListenAndServeUDP(addr string) error {
	// TODO add udp udpBindAddr to Server
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}
	s.udpConn, err = net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}
	defer s.udpConn.Close()
	for {
		// TODO use leaky buffer
		b := make([]byte, 65536)
		n, src, err := s.udpConn.ReadFromUDP(b)
		if err != nil {
			return err
		}
		s.ServeUDPDatagram(src, b[:n])
	}
}

// ServeConn is used to serve a single connection.
func (s *Server) ServeTCPConn(conn net.Conn) {
	defer conn.Close()
	bufConn := bufio.NewReader(conn)

	// Read the version byte
	version := []byte{0}
	if _, err := bufConn.Read(version); err != nil {
		log.Println("socks-quic error", fmt.Errorf("failed to get version bytes: %v", err))
		return
	}

	// Ensure we are compatible
	if version[0] != socks5Version {
		log.Println("socks-quic error", fmt.Errorf("unsupported SOCKS version: %v", version))
		return
	}

	// Authenticate the connection
	authContext, err := s.authenticate(conn, bufConn)
	if err != nil {
		log.Println("socks-quic failed to authenticate", err)
		return
	}

	request, err := NewRequest(bufConn)
	if err != nil {
		if err == errUnrecognizedAddrType {
			if err := sendReply(conn, addrTypeNotSupported, nil); err != nil {
				log.Println("socks-quic failed to send reply", err)
				return
			}
		}
		log.Println("socks-quic failed to read destination address", err)
		return
	}
	request.AuthContext = authContext
	if client, ok := conn.RemoteAddr().(*net.TCPAddr); ok {
		request.RemoteAddr = &AddrSpec{IP: client.IP, Port: client.Port}
	}

	// Process the client request
	if err := s.handleRequest(request, conn); err != nil {
		// error from handling request will not be log into file.
		//log.Println("socks-quic failed to handle request", err)
		return
	}
}

func (s *Server) ServeUDPDatagram(srcAddr *net.UDPAddr, b []byte) {
	r, err := NewRelay(srcAddr, b)
	if err != nil {
		log.Println("socks-quic failed to read udp datagram", err)
		return
	}
	if err := s.handleRelay(r); err != nil {
		// error from handling relay will not be log into file.
		//log.Println("socks-quic failed to handle relay", err)
		return
	}
}

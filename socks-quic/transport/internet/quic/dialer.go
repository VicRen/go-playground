package quic

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"sync"

	"github.com/VicRen/go-play-ground/socks-quic/transport/internet"
	quic "github.com/lucas-clemente/quic-go"
)

var errSessionClosed = errors.New("session closed")

type sessionContext struct {
	session quic.Session
}

func (s *sessionContext) openStreamSync(dstAddr net.Addr) (*interConn, error) {
	if !isActive(s.session) {
		return nil, errSessionClosed
	}
	stream, err := s.session.OpenStreamSync()
	if err != nil {
		return nil, err
	}
	conn := &interConn{
		Stream: stream,
		local:  s.session.LocalAddr(),
		remote: dstAddr,
	}
	return conn, nil
}

// TODO move to common
type destination struct {
	addr string
}

type clientSessions struct {
	access   sync.Mutex
	sessions map[destination][]*sessionContext
}

func openStream(sessions []*sessionContext, dstAddr net.Addr) *interConn {
	for _, s := range sessions {
		if !isActive(s.session) {
			continue
		}

		conn, err := s.openStreamSync(dstAddr)
		if err != nil {
			continue
		}

		return conn
	}

	return nil
}

func isActive(s quic.Session) bool {
	select {
	case <-s.Context().Done():
		return false
	default:
		return true
	}
}

func removeInactiveSession(sessions []*sessionContext) []*sessionContext {
	activeSessions := make([]*sessionContext, 0, len(sessions))
	for _, s := range sessions {
		if isActive(s.session) {
			activeSessions = append(activeSessions, s)
			continue
		}
		if err := s.session.Close(); err != nil {
			// TODO logger error failed to close session
		}
	}
	if len(activeSessions) < len(sessions) {
		return activeSessions
	}
	return sessions
}

var client clientSessions

func (c *clientSessions) cleanSessions() error {
	c.access.Lock()
	defer c.access.Unlock()

	if len(c.sessions) == 0 {
		return nil
	}

	newSessionMap := make(map[destination][]*sessionContext)
	for dst, sessions := range c.sessions {
		sessions = removeInactiveSession(sessions)
		if len(sessions) > 0 {
			newSessionMap[dst] = sessions
		}
	}

	c.sessions = newSessionMap
	return nil
}

func (c *clientSessions) openConnection(dstAddr net.Addr) (internet.Connection, error) {
	c.access.Lock()
	defer c.access.Unlock()

	if c.sessions == nil {
		c.sessions = make(map[destination][]*sessionContext)
	}

	dst := destination{dstAddr.String()}
	var sessions []*sessionContext
	if s, found := c.sessions[dst]; found {
		sessions = s
	}

	conn := openStream(sessions, dstAddr)
	if conn != nil {
		return conn, nil
	}

	sessions = removeInactiveSession(sessions)

	session, err := quic.DialAddr(dstAddr.String(), &tls.Config{InsecureSkipVerify: true},
		&quic.Config{KeepAlive: true, MaxIncomingStreams: MaxIncomingStreams})
	if err != nil {
		return nil, err
	}
	sc := &sessionContext{
		session: session,
	}
	c.sessions[dst] = append(sessions, sc)
	return sc.openStreamSync(dstAddr)
}

func init() {
	client.sessions = make(map[destination][]*sessionContext)
}

func Dial(ctx context.Context, addr net.Addr) (net.Conn, error) {
	return client.openConnection(addr)
}

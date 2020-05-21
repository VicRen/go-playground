package quic

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"net"

	"time"

	"github.com/VicRen/go-playground/socks-quic/common/signal/done"
	"github.com/VicRen/go-playground/socks-quic/transport/internet"
	quic "github.com/lucas-clemente/quic-go"
)

// MaxIncomingStreams ...
const MaxIncomingStreams = 1<<16 - 1 //65535

func Listen(ctx context.Context, addr string, handler internet.ConnHandler) (*Listener, error) {
	l, err := quic.ListenAddr(addr, generateTLSConfig(), &quic.Config{MaxIncomingStreams: MaxIncomingStreams})
	if err != nil {
		return nil, err
	}

	listener := &Listener{
		done:     done.New(),
		listener: l,
		addConn:  handler,
	}
	go listener.keepAccepting()
	return listener, nil
}

type Listener struct {
	listener quic.Listener
	done     *done.Instance
	addConn  internet.ConnHandler
}

func (l *Listener) keepAccepting() {
	for {
		conn, err := l.listener.Accept()
		if err != nil {
			// TODO failed to accept QUIC sessions
			//newError("failed to accept QUIC sessions").Base(err).WriteToLog()
			if l.done.Done() {
				break
			}
			// TODO why should I sleep here
			time.Sleep(time.Second)
			continue
		}
		go l.acceptStreams(conn)
	}
}

func (l *Listener) acceptStreams(session quic.Session) {
	for {
		stream, err := session.AcceptStream()
		if err != nil {
			// TODO logger failed to accept stream
			//newError("failed to accept stream").Base(err).WriteToLog()
			select {
			case <-session.Context().Done():
				return
			case <-l.done.Wait():
				_ = session.Close()
				return
			default:
				time.Sleep(time.Second)
				continue
			}
		}

		conn := &interConn{
			Stream: stream,
			local:  session.LocalAddr(),
			remote: session.RemoteAddr(),
		}

		l.addConn(conn)
	}

}

// Addr implements internet.Listener.Addr.
func (l *Listener) Addr() net.Addr {
	return l.listener.Addr()
}

// Close implements internet.Listener.Close.
func (l *Listener) Close() error {
	l.done.Close()
	l.listener.Close()
	return nil
}

// Setup a bare-bones TLS config for the rudp
func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{Certificates: []tls.Certificate{tlsCert}}
}

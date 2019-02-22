package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/big"

	quic_chan "github.com/VicRen/go-play-ground/quic-chan"

	quic "github.com/lucas-clemente/quic-go"
)

const addr = "localhost:4242"

func main() {
	err := echoServer()
	if err != nil {
		println(err)
	}
}

// Start a server that echos all data on the first stream opened by the client
func echoServer() error {
	listener, err := quic.ListenAddr(addr, generateTLSConfig(), nil)
	if err != nil {
		return err
	}
	for {
		sess, err := listener.Accept()
		if err != nil {
			break
		}
		go func() {
			for {
				stream, err := sess.AcceptStream()
				if err != nil {
					panic(err)
				}
				ret := quic_chan.GenerateStreamHandler(stream)
				go receiveData(ret.ReadChan, ret.WriteChan, ret.CloseChan)
				go printError(ret.StreamErrChan)
			}
		}()
	}
	return err
}

func receiveData(readChan <-chan []byte, writeChan chan<- []byte, closeChan chan struct{}) {
	for buf := range readChan {
		fmt.Printf("Server: Got '%s'\n", buf)
		writeChan <- buf
	}
	closeChan <- struct{}{}
	fmt.Println("printData finished")
}

func printError(errChan <-chan error) {
	for err := range errChan {
		fmt.Println(err)
	}
	fmt.Println("printError finished")
}

type listenerHandler struct {
	sessChan chan quic.Session
	errChan  chan error
}

func manageListener(listener quic.Listener) *listenerHandler {
	ret := &listenerHandler{
		make(chan quic.Session),
		make(chan error),
	}
	go func() {
		for {
			sess, err := listener.Accept()
			if err != nil {
				ret.errChan <- err
				return
			}
			ret.sessChan <- sess
		}
	}()
	return ret
}

type sessionHandler struct {
	streamChan chan quic.Stream
	errChan    chan error
}

func accecptSession(session quic.Session) *sessionHandler {
	ret := &sessionHandler{
		make(chan quic.Stream),
		make(chan error),
	}
	go func() {
		for {
			stream, err := session.AcceptStream()
			if err != nil {
				ret.errChan <- err
				return
			}
			ret.streamChan <- stream
		}
	}()
	return ret
}

type streamHandler struct {
	writeChan chan []byte
	readChan  chan []byte
	errChan   chan error
	closeChan chan struct{}
}

// Setup a bare-bones TLS config for the server
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

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
	dataChan := make(chan *quic_chan.RecvData)
	if err := echoServer(dataChan); err != nil {
		panic(err)
	}
	for d := range dataChan {
		fmt.Printf("Server: Got '%s', SessionID=%d, StreamID=%d\n", d.Data, d.SessionID, d.StreamID)
	}
}

// Start a server that echos all data on the first stream opened by the client
func echoServer(dataChan chan *quic_chan.RecvData) error {
	listener, err := quic.ListenAddr(addr, generateTLSConfig(), nil)
	if err != nil {
		return err
	}

	handleSession := func(handler *quic_chan.SessionHandler) {
		// TODO session ID valid?
		sessID := handler.SessionID
		for stream := range handler.StreamChan {
			streamHandler := quic_chan.GenerateStreamHandler(stream)
			err := quic_chan.PutStreamHandler(sessID, streamHandler)
			if err != nil {
				streamHandler.CloseChan <- struct{}{}
			} else {
				go echoData(streamHandler.ReadChan, streamHandler.WriteChan, streamHandler.CloseChan, dataChan, streamHandler.StreamID, sessID)
				go printError(streamHandler.StreamErrChan)
			}
		}
	}
	handleListener := func(handler *quic_chan.ListenerHandler) {
		for sess := range handler.SessChan {
			sessHandler := quic_chan.GenerateSessionHandler(sess)
			quic_chan.PutSessionHandler(sessHandler)
			go handleSession(sessHandler)
			go func() {
				for err := range sessHandler.SessErrChan {
					fmt.Println("SessionErr:", err)
				}
				fmt.Println("PrintSessionErr finished")
			}()
		}
	}

	listenerHandler := quic_chan.GenerateListenerHandler(listener)
	go handleListener(listenerHandler)
	return nil
}

func echoData(readChan <-chan []byte, writeChan chan<- []byte,
	closeChan chan struct{}, dataChan chan *quic_chan.RecvData, streamID quic.StreamID, sessionID uint64) {
	for buf := range readChan {
		writeChan <- buf
		dataChan <- &quic_chan.RecvData{
			SessionID: sessionID,
			StreamID:  uint64(streamID),
			Data:      buf,
		}
	}
	closeChan <- struct{}{}
}

func printError(errChan <-chan error) {
	for err := range errChan {
		fmt.Println(err)
	}
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

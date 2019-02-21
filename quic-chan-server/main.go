package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"sync"

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
				ret := &streamHandler{
					make(chan []byte),
					make(chan []byte),
					make(chan error),
					make(chan error),
					make(chan struct{}),
				}
				go writeStream(stream, ret.writeChan, ret.writeErrChan, ret.closeChan)
				go readStream(stream, ret.readChan, ret.readErrChan)
				go receiveData(ret.readChan, ret.writeChan)
				go printError(merge(ret.writeErrChan, ret.readErrChan))
			}
		}()
	}
	return err
}

func acceptSession() {

}

type streamHandler struct {
	writeChan    chan []byte
	readChan     chan []byte
	writeErrChan chan error
	readErrChan  chan error
	closeChan    chan struct{}
}

func receiveData(readChan <-chan []byte, writeChan chan<- []byte) {
	for buf := range readChan {
		fmt.Printf("Client: Got '%s'\n", buf)
		writeChan <- buf
	}
	fmt.Println("printData finished")
}

func printError(errChan <-chan error) {
	for err := range errChan {
		fmt.Println(err)
	}
	fmt.Println("printError finished")
}

func merge(errChan ...<-chan error) <-chan error {
	out := make(chan error)
	var wg sync.WaitGroup
	collect := func(in <-chan error) {
		defer wg.Done()
		for n := range in {
			out <- n
		}
	}
	wg.Add(len(errChan))
	for _, c := range errChan {
		go collect(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func writeStream(stream quic.Stream, writeChan <-chan []byte, errChan chan<- error, quitChan <-chan struct{}) {
	for {
		select {
		case data := <-writeChan:
			_, err := stream.Write(data)
			if err != nil {
				errChan <- err
				close(errChan)
			}
		case <-quitChan:
			err := stream.Close()
			if err != nil {
				// stream close error
				fmt.Println(err)
			}
		}
	}
}

func readStream(stream quic.Stream, readChan chan<- []byte, errChan chan<- error) {
	for {
		buf := make([]byte, 1024)
		n, err := stream.Read(buf)
		if n > 0 {
			readChan <- buf[:n]
		}
		if err == io.EOF {
			errChan <- err
			close(errChan)
			close(readChan)
			break
		} else if err != nil {
			errChan <- err
			close(errChan)
			close(readChan)
			break
		}
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

package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"math/big"
	"strconv"
	"time"

	quic "github.com/lucas-clemente/quic-go"
)

const addr = "localhost:4242"

const message = "foobar"

func main() {
	go func() { log.Fatal(echoServer(addr)) }()

	err := clientMain()
	if err != nil {
		fmt.Println(err)
	}
	select {}
}

func clientMain() error {
	dc := make(chan string)
	session, err := quic.DialAddr(addr, &tls.Config{InsecureSkipVerify: true}, nil)
	if err != nil {
		return err
	}

	stream, err := session.OpenStreamSync()
	if err != nil {
		return err
	}

	stream2, err := session.OpenStreamSync()
	if err != nil {
		return err
	}

	go func() {
		for i := 10; i > 0; i-- {
			m := message + strconv.Itoa(i)
			fmt.Printf("Client1: Sending '%s'\n", m)
			_, err = stream.Write([]byte(m))
			if err != nil {
				fmt.Printf("Client writing data to stream1: %v", err)
			}
			time.Sleep(100 * time.Millisecond)
		}
		time.Sleep(1000 * time.Millisecond)
		session.Close()
	}()

	go func() {
		for i := 10; i > 0; i-- {
			m := message + strconv.Itoa(i)
			fmt.Printf("Client2: Sending '%s'\n", m)
			_, err = stream2.Write([]byte(m))
			if err != nil {
				fmt.Printf("Client writing data to stream1: %v", err)
			}
			time.Sleep(100 * time.Millisecond)
		}
		time.Sleep(1000 * time.Millisecond)
		session.Close()
	}()

	go func() {
		for {
			buf := make([]byte, len(message)+2)
			n, err := stream.Read(buf)
			if err != nil {
				fmt.Printf("Client1 reading error: %v\n", err)
				break
			}
			dc <- "client1:" + string(buf[:n])
		}
		//close(dc)
	}()

	go func() {
		for {
			buf := make([]byte, len(message)+2)
			n, err := stream2.Read(buf)
			if err != nil {
				fmt.Printf("Client2 reading error: %v\n", err)
				break
			}
			dc <- "client2:" + string(buf[:n])
		}
		//close(dc)
	}()

	for d := range dc {
		fmt.Printf("Client: Got '%s'\n", d)
	}
	return nil
}

// Start a server that echos all data on the first stream opened by the client
func echoServer(addr string) error {
	listener, err := quic.ListenAddr(addr, generateTLSConfig(), nil)
	if err != nil {
		return err
	}
	sessChan := make(chan quic.Session)
	go func() {
		for {
			sess, err := listener.Accept()
			if err != nil {
				fmt.Printf("Server reading error: %v\n", err)
				break
			}
			sessChan <- sess
		}
		close(sessChan)
	}()
	if err != nil {
		return err
	}
	streamChan := make(chan quic.Stream)
	go func() {
		for sess := range sessChan {
			go func() {
				for {
					stream, err := sess.AcceptStream()
					if err != nil {
						fmt.Printf("Server reading stream error: %v\n", err)
						break
					}
					fmt.Printf("Server Got one stream\n\n")
					streamChan <- stream
				}
			}()
		}
		close(streamChan)
	}()
	dataChan := make(chan string)
	go func() {
		count1 := 0
		for stream := range streamChan {
			count1++
			go readData(stream, count1, dataChan)
		}
		close(dataChan)
	}()
	count := 0
	for b := range dataChan {
		fmt.Printf("------->Server: Number %d Got '%s'\n", count, string(b))
		count++
	}
	// Echo through the loggingWriter
	//_, err = io.Copy(loggingWriter{stream}, stream)
	return fmt.Errorf("Server end reading error: %v\n", err)
}

func readData(stream quic.Stream, count1 int, dc chan string) {
	fmt.Printf("Server reading one stream\n\n")
	count := 0
	for {
		fmt.Printf("\n\nServer%s reading one stream loop\n", strconv.Itoa(count1))
		buf := make([]byte, len(message)+2)
		n, err := stream.Read(buf)
		go writeData(stream, []byte(strconv.Itoa(count)))
		if err != nil {
			fmt.Printf("Server reading buffer error: %v\n", err)
			break
		}
		d := "server" + strconv.Itoa(count1) + ":" + string(buf[:n])
		fmt.Printf("Server receving: %s\n", d)
		dc <- d
		count++
	}
	//Echo through the loggingWriter
	//_, err = io.Copy(loggingWriter{stream}, stream)
}

func writeData(stream quic.Stream, data []byte) {
	stream.Write(data)
}

// A wrapper for io.Writer that also logs the message.
type loggingWriter struct{ io.Writer }

func (w loggingWriter) Write(b []byte) (int, error) {
	fmt.Printf("Server: Got '%s'\n", string(b))
	return w.Writer.Write(b)
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

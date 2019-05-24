package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"

	"github.com/VicRen/go-play-ground/route"

	quic "github.com/lucas-clemente/quic-go"
)

func main() {
	log.Fatal(serverMain())
}

func serverMain() error {
	sessOut, err := quic.DialAddr(route.PeerB, &tls.Config{InsecureSkipVerify: true}, nil)
	if err != nil {
		return err
	}
	fmt.Println("connected to", route.PeerB)
	streamOut, err := sessOut.OpenStreamSync()
	if err != nil {
		return err
	}
	listener, err := quic.ListenAddr(route.Server1, generateTLSConfig(), nil)
	fmt.Println("Listening on", route.Server1)
	sess, err := listener.Accept()
	if err != nil {
		return err
	}
	for {
		stream, err := sess.AcceptStream()
		if err != nil {
			return err
		}
		go func() {
			buf := make([]byte, 100)
			count := 0
			for {
				n, err := stream.Read(buf)
				if err != nil {
					break
				}
				streamOut.Write(buf[:n])
				fmt.Printf("Server: Got '%s' from Stream '%d'\n", string(buf[:n]), stream.StreamID())
				if count == 5 {
					fmt.Println("Server end receive")
					sess.Close()
					count = 0
					return
				}
				count++
			}
		}()
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

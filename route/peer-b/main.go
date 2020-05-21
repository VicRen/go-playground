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

	"github.com/VicRen/go-playground/route"

	quic "github.com/lucas-clemente/quic-go"
)

func main() {
	log.Fatal(serverMain())
}

func serverMain() error {
	listener, err := quic.ListenAddr(route.PeerB, generateTLSConfig(), nil)
	if err != nil {
		return err
	}
	fmt.Println("Listening on", route.PeerB)
	for {
		sess, err := listener.Accept()
		if err != nil {
			return err
		}
		go listenSession(sess)
	}
}

func listenSession(session quic.Session) {
	for {
		stream, err := session.AcceptStream()
		if err != nil {
			panic(err)
		}
		go func() {
			buf := make([]byte, 100)
			for {
				n, err := stream.Read(buf)
				if err != nil {
					break
				}
				fmt.Printf("PeerB: Got '%s' from Stream '%d'\n", string(buf[:n]), stream.StreamID())
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

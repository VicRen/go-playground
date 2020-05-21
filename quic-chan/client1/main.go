package main

import (
	"crypto/tls"
	"fmt"
	"time"

	quic_chan "github.com/VicRen/go-playground/quic-chan"

	quic "github.com/lucas-clemente/quic-go"
)

const addr = "localhost:4242"

const message = "foobar"

func main() {
	err := clientMain()
	if err != nil {
		panic(err)
	}
	err = clientMain()
	if err != nil {
		panic(err)
	}
}

func clientMain() error {
	session, err := quic.DialAddr(addr, &tls.Config{InsecureSkipVerify: true}, nil)
	if err != nil {
		return err
	}

	stream, err := session.OpenStreamSync()
	if err != nil {
		return err
	}
	ret := quic_chan.GenerateStreamHandler(stream)
	go writeData(ret.WriteChan, ret.CloseChan, []byte(message), 3)
	go printError(ret.StreamErrChan, ret.CloseChan)
	printData(ret.ReadChan)
	return nil
}

func writeData(writeChan chan<- []byte, quitChan chan<- struct{}, data []byte, n int) {
	for i := n; i > 0; i-- {
		time.Sleep(100 * time.Millisecond)
		writeChan <- data
	}
	quitChan <- struct{}{}
}

func printData(readChan <-chan []byte) {
	for buf := range readChan {
		fmt.Printf("Client: Got '%s'\n", buf)
	}
	fmt.Println("printData finished")
}

func printError(errChan <-chan error, quitChan chan struct{}) {
	for err := range errChan {
		fmt.Println(err)
		// tell write goroutine to close stream
		quitChan <- struct{}{}
	}
	fmt.Println("printError finished")
}

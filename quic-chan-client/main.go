package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"sync"
	"time"

	quic "github.com/lucas-clemente/quic-go"
)

const addr = "localhost:4242"

const message = "foobar"

func main() {
	go clientMain()
	err := clientMain()
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
	ret := manageStream(stream)
	go writeData(ret.writeChan, ret.closeChan, []byte(message), 3)
	go printError(ret.errChan)
	printData(ret.readChan)
	return nil
}

func writeData(writeChan chan<- []byte, closeChan chan<- struct{}, data []byte, n int) {
	for i := n; i > 0; i-- {
		time.Sleep(100 * time.Millisecond)
		writeChan <- data
	}
	closeChan <- struct{}{}
}

type streamHandler struct {
	writeChan chan []byte
	readChan  chan []byte
	errChan   chan error
	closeChan chan struct{}
}

func printData(readChan <-chan []byte) {
	for buf := range readChan {
		fmt.Printf("Client: Got '%s'\n", buf)
	}
	fmt.Println("printData finished")
}

func printError(errChan <-chan error) {
	for err := range errChan {
		fmt.Println(err)
	}
	fmt.Println("printError finished")
}

func merge(errChan ...chan error) chan error {
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

func manageStream(stream quic.Stream) *streamHandler {
	writeErrChan := make(chan error)
	readErrChan := make(chan error)
	ret := &streamHandler{
		make(chan []byte),
		make(chan []byte),
		merge(writeErrChan, readErrChan),
		make(chan struct{}),
	}
	go writeStream(stream, ret.writeChan, writeErrChan, ret.closeChan)
	go readStream(stream, ret.readChan, readErrChan)
	return ret
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
		if err != nil {
			close(readChan)
			if err == io.EOF {
				errChan <- err
			}
			close(errChan)
		}
		if err == io.EOF {
			close(readChan)
			close(errChan)
			break
		} else if err != nil {
			close(readChan)
			errChan <- err
			close(errChan)
			break
		}
	}
}

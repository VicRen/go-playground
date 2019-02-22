package main

import (
	"crypto/tls"
	"fmt"
	"time"

	quic_chan "github.com/VicRen/go-play-ground/quic-chan"
	quic "github.com/lucas-clemente/quic-go"
)

const addr = "localhost:4242"

const message = "foobar"

func main() {
	go runInstance()
	runInstance()
}

func runInstance() {
	dataChan := make(chan *quic_chan.RecvData)
	err := clientMain(dataChan)
	if err != nil {
		panic(err)
	}
	for d := range dataChan {
		fmt.Printf("Client: Got '%s', SessionID=%d, StreamID=%d\n", d.Data, d.SessionID, d.StreamID)
	}
}

func clientMain(dataChan chan *quic_chan.RecvData) error {
	listenSession := func(sessHandler *quic_chan.SessionHandler) {
		for stream := range sessHandler.StreamChan {
			handler := quic_chan.GenerateStreamHandler(stream)
			err := quic_chan.PutStreamHandler(sessHandler.SessionID, handler)
			if err != nil {
				handler.CloseChan <- struct{}{}
			} else {
				go printData(handler.ReadChan, dataChan, sessHandler.SessionID, handler.StreamID)
				go printError(handler.StreamErrChan)
			}
		}
	}
	session, err := quic.DialAddr(addr, &tls.Config{InsecureSkipVerify: true}, nil)
	if err != nil {
		return err
	}
	sessHandler := quic_chan.GenerateSessionHandler(session)
	quic_chan.PutSessionHandler(sessHandler)

	go func() {
		stream, err := session.OpenStreamSync()
		if err != nil {
			// TODO error
			return
		}
		handler := quic_chan.GenerateStreamHandler(stream)
		err = quic_chan.PutStreamHandler(sessHandler.SessionID, handler)
		if err != nil {
			handler.CloseChan <- struct{}{}
		} else {
			go writeData(handler.WriteChan, handler.CloseChan, []byte(message), 3)
			go printData(handler.ReadChan, dataChan, sessHandler.SessionID, handler.StreamID)
			go printError(handler.StreamErrChan)
		}
	}()
	go listenSession(sessHandler)
	return nil
}

func writeData(writeChan chan<- []byte, quitChan chan<- struct{}, data []byte, n int) {
	for i := n; i > 0; i-- {
		time.Sleep(100 * time.Millisecond)
		writeChan <- data
	}
	quitChan <- struct{}{}
}

func printData(readChan <-chan []byte, dataChan chan<- *quic_chan.RecvData, sessionID uint64, streamID quic.StreamID) {
	for buf := range readChan {
		dataChan <- &quic_chan.RecvData{
			SessionID: sessionID,
			StreamID:  uint64(streamID),
			Data:      buf,
		}
	}
}

func printError(errChan <-chan error) {
	for err := range errChan {
		fmt.Println(err)
	}
}

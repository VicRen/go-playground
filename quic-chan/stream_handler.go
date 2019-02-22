package quic_chan

import (
	"fmt"
	"io"
	"sync"

	quic "github.com/lucas-clemente/quic-go"
)

// StreamHandler provides ways to control Stream in a much easy way.
type StreamHandler struct {
	// WriteChan provides a way the write to the Stream.
	// It will be closed.
	WriteChan chan []byte
	// ReadChan returns data reading from the Stream.
	// When the read side of the Stream is closed or error occur during reading, it will be closed.
	ReadChan chan []byte
	// StreamErrChan returns the error of both side of the Stream.
	// When the Stream is closed from both side, it will be closed.
	StreamErrChan chan error
	// CloseChan closes the write side of a Stream.
	// When Stream is a SendStream, StreamHandler will be destroyed.
	// When Stream is a ReceiveStream, StreamHandler will wait the read side of the Stream to be cancelled.
	// Write to WriteChan will not be allowed after CloseChan sent.
	CloseChan chan struct{}
}

// GenerateStreamHandler returns a StreamHandler that provides channels to
// write or read from a Stream. Also there's a channel to close Stream.
func GenerateStreamHandler(stream quic.Stream) *StreamHandler {
	writeErrChan := make(chan error)
	readErrChan := make(chan error)
	ret := &StreamHandler{
		make(chan []byte),
		make(chan []byte),
		merge(writeErrChan, readErrChan),
		make(chan struct{}),
	}
	go writeStream(stream, ret.WriteChan, writeErrChan, ret.CloseChan)
	go readStream(stream, ret.ReadChan, readErrChan)
	return ret
}

// GenerateSendStreamHandler returns a StreamHandler that provides channel to
// write to a Stream. Also there's a channel to close Stream.
func GenerateSendStreamHandler(stream quic.SendStream) *StreamHandler {
	ret := &StreamHandler{
		make(chan []byte),
		nil,
		make(chan error),
		make(chan struct{}),
	}
	go writeStream(stream, ret.WriteChan, ret.StreamErrChan, ret.CloseChan)
	return ret
}

// GenerateReceiveStreamHandler returns a StreamHandler that provides channels to
// read from a Steam. Also there's a channel to close Stream.
func GenerateReceiveStreamHandler(stream quic.ReceiveStream) *StreamHandler {
	ret := &StreamHandler{
		nil,
		make(chan []byte),
		make(chan error),
		make(chan struct{}),
	}
	go readStream(stream, ret.ReadChan, ret.StreamErrChan)
	return ret
}

func writeStream(stream io.WriteCloser, writeChan chan []byte, errChan chan<- error, quitChan <-chan struct{}) {
	for {
		select {
		case data := <-writeChan:
			_, err := stream.Write(data)
			if err != nil {
				errChan <- err
				close(errChan)
			}
		case <-quitChan:
			close(errChan)
			err := stream.Close()
			if err != nil {
				// stream close error
				fmt.Println(err)
			}
		}
	}
}

func readStream(stream io.Reader, readChan chan<- []byte, errChan chan<- error) {
	for {
		buf := make([]byte, 1024)
		n, err := stream.Read(buf)
		if n > 0 {
			readChan <- buf[:n]
		}
		// err occur
		if err != nil {
			// err is not EOF, notify
			if err != io.EOF {
				errChan <- err
			}
			// close errChan and readChan.
			close(errChan)
			close(readChan)
			return
		}
	}
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

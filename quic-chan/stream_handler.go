package quic_chan

import (
	"io"

	quic "github.com/lucas-clemente/quic-go"
)

// StreamHandler provides ways to control Stream in a much easy way.
type StreamHandler struct {
	// StreamID is the unique ID for Stream.
	StreamID quic.StreamID
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
	// When Stream is a ReceiveStreamChan, StreamHandler will wait the read side of the Stream to be cancelled.
	// Write to WriteChan will not be allowed after CloseChan sent.
	CloseChan chan struct{}
}

// GenerateStreamHandler returns a StreamHandler that provides channels to
// write or read from a Stream. Also there's a channel to close Stream.
func GenerateStreamHandler(stream quic.Stream) *StreamHandler {
	writeErrChan := make(chan error)
	readErrChan := make(chan error)
	ret := &StreamHandler{
		stream.StreamID(),
		make(chan []byte),
		make(chan []byte),
		MergeErrChan(writeErrChan, readErrChan),
		make(chan struct{}),
	}
	go writeToStream(stream, ret.WriteChan, writeErrChan, ret.CloseChan)
	go readFromStream(stream, ret.ReadChan, readErrChan)
	return ret
}

func writeToStream(stream io.WriteCloser, writeChan chan []byte, errChan chan<- error, quitChan <-chan struct{}) {
	for {
		select {
		case data := <-writeChan:
			_, err := stream.Write(data)
			if err != nil {
				errChan <- err
				close(errChan)
				return
			}
		case <-quitChan:
			err := stream.Close()
			if err != nil {
				// stream close error
				errChan <- err
			}
			close(errChan)
			return
		}
	}
}

func readFromStream(stream io.Reader, readChan chan<- []byte, errChan chan<- error) {
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

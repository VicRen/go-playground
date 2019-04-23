package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/VicRen/go-play-ground/route"

	quic "github.com/lucas-clemente/quic-go"
)

const message = "hello galaxy"

var nodes = []string{route.Server2, route.Server1, route.PeerB}

func main() {
	err := clientMain()
	if err != nil {
		panic(err)
	}
}

type Foo struct {
	index int
	quic.Session
	quic.Stream
}

func (f *Foo) init() {
	f.nextSession()
}

func (f *Foo) Write(p []byte) (n int, err error) {
	if f.Session == nil {
		err = f.nextSession()
		if err != nil {
			return
		}
	}
	if f.Stream == nil {
		f.Stream, err = f.Session.OpenStreamSync()
		if err != nil {
			err = f.nextSession()
			return f.Write(p)
		}
	}
	n, err = f.Stream.Write(p)
	if err != nil {
		fmt.Println("failed write", string(p), " on stream", f.StreamID())
		f.Stream = nil
		return f.Write(p)
	}
	fmt.Println("write on stream", f.StreamID())
	return
}

func (f *Foo) nextSession() error {
	if f.Session != nil {
		f.Session.Close()
	}
	if f.index > len(nodes) {
		panic("client terminated")
	}
	session, err := quic.DialAddr(nodes[f.index], &tls.Config{InsecureSkipVerify: true}, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}
	f.Session = session
	f.index++
	return nil
}

func clientMain() error {
	r, w := io.Pipe()
	go func() {
		n := 0
		for {
			time.Sleep(1000 * time.Millisecond)
			s := message + " " + strconv.Itoa(n)
			_, err := w.Write([]byte(s))
			if err != nil {
				fmt.Println("error sending message: ", n, err)
				return
			}
			n++
		}
	}()
	f := &Foo{}
	f.init()
	_, err := copyBuffer(loggingWriter{f}, r, nil)
	fmt.Println("end of copy")
	if err != nil {
		return err
	}
	return nil
}

// A wrapper for io.Writer that also logs the message.
type loggingWriter struct{ io.Writer }

func (w loggingWriter) Write(b []byte) (int, error) {
	fmt.Printf("Client: Sending '%s' \n", string(b))
	return w.Writer.Write(b)
}

// copyBuffer is the actual implementation of Copy and CopyBuffer.
// if buf is nil, one is allocated.
func copyBuffer(dst io.Writer, src io.Reader, buf []byte) (written int64, err error) {
	// If the reader has a WriteTo method, use it to do the copy.
	// Avoids an allocation and a copy.
	if wt, ok := src.(io.WriterTo); ok {
		return wt.WriteTo(dst)
	}
	// Similarly, if the writer has a ReadFrom method, use it to do the copy.
	if rt, ok := dst.(io.ReaderFrom); ok {
		return rt.ReadFrom(src)
	}
	if buf == nil {
		size := 32 * 1024
		if l, ok := src.(*io.LimitedReader); ok && int64(size) > l.N {
			if l.N < 1 {
				size = 1
			} else {
				size = int(l.N)
			}
		}
		buf = make([]byte, size)
	}
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				fmt.Println("short write nr", nr, "nw", nw)
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written, err
}

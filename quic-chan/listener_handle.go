package quic_chan

import (
	quic "github.com/lucas-clemente/quic-go"
)

type ListenerHandler struct {
	SessChan  chan quic.Session
	ErrChan   chan error
	CloseChan chan struct{}
}

func GenerateListenerHandler(listener quic.Listener) *ListenerHandler {
	acceptErrChan := make(chan error)
	closeErrChan := make(chan error)
	ret := &ListenerHandler{
		make(chan quic.Session),
		MergeErrChan(acceptErrChan, closeErrChan),
		make(chan struct{}),
	}
	go acceptFromListener(listener, ret.SessChan, acceptErrChan)
	go closeListener(listener, closeErrChan, ret.CloseChan)
	return ret
}

func acceptFromListener(listener quic.Listener, sessChan chan quic.Session, errChan chan error) {
	for {
		sess, err := listener.Accept()
		if err != nil {
			errChan <- err
			close(errChan)
			return
		}
		sessChan <- sess
	}
}

func closeListener(listener quic.Listener, errChan chan error, closeChan chan struct{}) {
	select {
	case <-closeChan:
		err := listener.Close()
		if err != nil {
			errChan <- err
		}
		close(errChan)
	}
}

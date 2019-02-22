package quic_chan

import (
	"fmt"

	quic "github.com/lucas-clemente/quic-go"
)

type SessionHandler struct {
	SessionID   uint64
	StreamChan  chan quic.Stream
	SessErrChan chan error
	CloseChan   chan struct{}
}

func GenerateSessionHandler(session quic.Session) *SessionHandler {
	acceptErrChan := make(chan error)
	closeErrChan := make(chan error)
	openStreamErrChan := make(chan error)
	ret := &SessionHandler{
		GenerateSessionID(),
		make(chan quic.Stream),
		MergeErrChan(acceptErrChan, closeErrChan, openStreamErrChan),
		make(chan struct{}),
	}
	go closeSession(session, closeErrChan, ret.CloseChan)
	go acceptFromSession(session, ret.StreamChan, acceptErrChan, ret.CloseChan)
	return ret
}

func acceptFromSession(session quic.Session, acceptChan chan<- quic.Stream, errChan chan<- error, closeChan chan<- struct{}) {
	for {
		stream, err := session.AcceptStream()
		if err != nil {
			errChan <- err
			close(errChan)
			closeChan <- struct{}{}
			return
		}
		acceptChan <- stream
	}
}

func closeSession(session quic.Session, errChan chan<- error, closeChan <-chan struct{}) {
	select {
	case <-closeChan:
		err := session.Close()
		if err != nil {
			fmt.Println(err)
			errChan <- err
		}
		close(errChan)
	}
}

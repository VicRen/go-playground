package quic_chan

import (
	"errors"
	"sync"
)

type Sess struct {
	sessHandler      *SessionHandler
	streamHandlerMap map[uint64]*StreamHandler
}

var sessHandlerMap = map[uint64]*Sess{}
var sessHandlerMapMutex sync.Mutex

func PutSessionHandler(handler *SessionHandler) {
	sessHandlerMapMutex.Lock()
	defer sessHandlerMapMutex.Unlock()
	old, ok := sessHandlerMap[handler.SessionID]
	oldStreamMap := map[uint64]*StreamHandler{}
	if ok {
		oldStreamMap = old.streamHandlerMap
	}
	sessHandlerMap[handler.SessionID] = &Sess{
		handler,
		oldStreamMap,
	}
}

func PutStreamHandler(sessID uint64, stream *StreamHandler) error {
	sessHandlerMapMutex.Lock()
	defer sessHandlerMapMutex.Unlock()
	sess, ok := sessHandlerMap[sessID]
	if !ok {
		return errors.New("cannot find sessionID")
	}
	sess.streamHandlerMap[uint64(stream.StreamID)] = stream
	return nil
}

package quic_chan

import (
	"sync"
)

var (
	sessIDLock sync.Mutex
	sessIDBase uint64
)

func GenerateSessionID() uint64 {
	sessIDLock.Lock()
	defer sessIDLock.Unlock()
	sessIDBase++
	return sessIDBase
}

func IsSessionIDValid(sessID uint64) bool {
	return sessID > 0
}

func MergeErrChan(errChan ...chan error) chan error {
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

func FanOutSignalChan(signalChan chan struct{}, n int) []chan struct{} {
	ret := make([]chan struct{}, 0)
	for n > 0 {
		s := make(chan struct{})
		ret = append(ret, s)
		n--
	}
	go func() {
		for signal := range signalChan {
			for _, v := range ret {
				v <- signal
			}
		}
	}()
	return ret
}

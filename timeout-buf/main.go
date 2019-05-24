package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"
)

func main() {
	exit := false
	tb := NewTimeoutBuffer([]byte("testing"), func() {
		fmt.Println("calling after func---------")
		exit = true
	})
	go func() {
		for {
			buf := make([]byte, 2)
			n, err := tb.Read(buf)
			if err != nil {
				fmt.Println("error reading from timeout buffer:", err)
				return
			}
			fmt.Println("data read:", string(buf[:n]))
		}
	}()
	count := 0
	data := []byte(strconv.Itoa(count))
	for {
		if count > 4 {
			time.Sleep(4 * time.Second)
		}
		if count > 3 {
			time.Sleep(3 * time.Second)
		}
		fmt.Println("---------loop--------", exit)
		if exit {
			fmt.Println("timeout, exiting with data:", data)
			return
		}
		err := tb.Append(data)
		if err != nil {
			fmt.Println("error writing data:", err)
			continue
		}
		data = append(data, []byte("2")...)
		count++
	}
}

type TimeoutBuffer struct {
	reader    *bytes.Buffer
	dataCh    chan []byte
	closeCh   chan struct{}
	timer     *time.Timer
	dataMutex sync.RWMutex
	closed    bool
}

func NewTimeoutBuffer(data []byte, afterFunc func()) *TimeoutBuffer {
	ret := &TimeoutBuffer{}
	ret.reader = bytes.NewBuffer(data)
	ret.dataCh = make(chan []byte)
	ret.closeCh = make(chan struct{})
	ret.timer = time.AfterFunc(4*time.Second, func() {
		afterFunc()
		ret.closeCh <- struct{}{}
	})
	return ret
}

func (t *TimeoutBuffer) Read(p []byte) (n int, err error) {
	n, err = t.reader.Read(p)
	if err == io.EOF && !t.isClosed() {
		fmt.Println("-----Read EOF-----")
		select {
		case data := <-t.dataCh:
			t.resetTimer()
			t.reader.Write(data)
		case <-t.closeCh:
			t.dataMutex.Lock()
			t.closed = true
			t.dataMutex.Unlock()
			close(t.dataCh)
		}
		return t.Read(p)
	}
	return
}

func (t *TimeoutBuffer) Append(data []byte) error {
	fmt.Println("-----> append data:", data)
	if t.isClosed() {
		return errors.New("timeout")
	}
	t.dataCh <- data
	return nil
}

func (t *TimeoutBuffer) resetTimer() bool {
	ret := t.timer.Stop()
	t.timer.Reset(4 * time.Second)
	return ret
}

func (t *TimeoutBuffer) isClosed() bool {
	t.dataMutex.RLock()
	defer t.dataMutex.RUnlock()
	return t.closed
}

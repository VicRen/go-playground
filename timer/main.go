package main

import (
	"fmt"
	"time"
)

func main() {
	//timer()
	afterTimer()
}

func timer() {
	tn := time.Now().UnixNano()
	stopChannel := make(chan struct{})
	t := time.NewTimer(3 * time.Second)
	fmt.Println("timer started")
	afterFunc := func() {
		select {
		case <-t.C:
			fmt.Println("timeout:", time.Now().UnixNano()-tn)
		case <-stopChannel:
			fmt.Println("recv stop chan")
		}
		fmt.Println("afterFunc called")
	}
	go afterFunc()
	stopFunc := func() {
		if t.Stop() {
			fmt.Println("sending stop chan")
			stopChannel <- struct{}{}
		}
		time.Sleep(1 * time.Second)
		t.Reset(1 * time.Second)
		go afterFunc()
	}
	time.Sleep(1 * time.Second)
	stopFunc()
	fmt.Println("stop func called")
	time.Sleep(10 * time.Second)
}

func afterTimer() {
	stopChannel := make(chan struct{})
	t := time.AfterFunc(3*time.Second, func() {
		fmt.Println("afterFunc called")
	})
	stopFunc := func() {
		if t.Stop() {
			fmt.Println("sending stop chan")
			stopChannel <- struct{}{}
		}
		time.Sleep(1 * time.Second)
		t.Reset(1 * time.Second)
	}
	go func() {
		select {
		case <-t.C:
			fmt.Println("t.C read")
		case <-stopChannel:
			fmt.Println("recv stop chan")
		}
		fmt.Println("go func called")
	}()
	time.Sleep(1 * time.Second)
	go stopFunc()
	time.Sleep(5 * time.Second)
}

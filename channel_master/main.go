package main

import (
	"fmt"
	"time"
)

func main() {
	practiceSync()
}

// 只向 chan 中发送了三个 struct，而有四个 receiver，因此最后一个 receiver 的 ok 为 false.
// 由于 receiver 的接收顺序的不确定性，输出值并不确定，但至少有个一 hello channel x false.
func practiceOne() {
	c := make(chan struct{}, 2)
	go func() {
		fmt.Println("hello channel")
		c <- struct{}{}
		c <- struct{}{}
		c <- struct{}{}
		close(c)
	}()
	go func() {
		_, ok := <-c
		fmt.Printf("hello channel 5 %v\n", ok)
	}()
	go func() {
		_, ok := <-c
		fmt.Printf("hello channel 4 %v\n", ok)
	}()
	go func() {
		_, ok := <-c
		fmt.Printf("hello channel 3 %v\n", ok)
	}()
	<-c
	fmt.Println("hello channel again")
	time.Sleep(1 * time.Second)
}

func practiceBlock() {
	s := []int{1, 2, 3, 5, 8, 13, 21, 24}
	c := make(chan int)
	go sum(s[:len(s)/2], c)
	go sum(s[len(s)/2:], c)
	x, y := <-c, <-c
	fmt.Println(x, y, x+y)
}

func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum
}

func practiceRange() {
	c := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			c <- i
		}
		close(c)
	}()

	for i := range c {
		fmt.Println(i)
	}
	fmt.Println("Finished")
}

func practiceSelect() {
	c := make(chan int)
	quit := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(<-c)
		}
		quit <- 0
	}()
	fibonacci(c, quit)
}

func fibonacci(c, quit chan int) {
	x, y := 0, 1
	for {
		select {
		case c <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}

func practiceTimeout() {
	c := make(chan string)
	go func() {
		time.Sleep(2 * time.Second)
		c <- "result1"
	}()
	select {
	case res := <-c:
		fmt.Println(res)
	case <-time.After(1 * time.Second):
		fmt.Println("time out")
	}
}

func practiceTimer() {
	timer := time.NewTimer(1 * time.Second)
	go func() {
		<-timer.C
		fmt.Println("timer expired")
	}()
	time.Sleep(2 * time.Second)
	ok := timer.Stop()
	if ok {
		fmt.Println("timer stopped")
	} else {
		fmt.Println("timer stop failed")
	}
}

func practiceTicker() {
	quit := make(chan struct{})
	ticker := time.NewTicker(500 * time.Millisecond)
	go func() {
		time.Sleep(2 * time.Second)
		quit <- struct{}{}
		ticker.Stop()
	}()
	//for t := range ticker.C {
	//	fmt.Println("Tick at ", t)
	//}

	for {
		select {
		case t, ok := <-ticker.C:
			fmt.Println("Tick at", t, ok)
		case <-quit:
			fmt.Println("ticker stopped")
			return
		}
	}
}

func practiceSendToClosedChan() {
	c := make(chan int, 10)
	c <- 1
	c <- 2
	close(c)
	c <- 3
}

func practiceReadFromClosedChan() {
	c := make(chan int, 10)
	c <- 1
	c <- 2
	close(c)
	fmt.Println(<-c) // 1
	fmt.Println(<-c) // 2
	fmt.Println(<-c) // 0
	fmt.Println(<-c) // 0
	fmt.Println(<-c) // 0
}

func practiceReadFromClosedChanWithRange() {
	c := make(chan int, 10)
	c <- 1
	c <- 2
	close(c)
	for i := range c {
		fmt.Println(i)
	}
}

func practiceSync() {
	c := make(chan struct{})
	go func() {
		time.Sleep(2 * time.Second)
		c <- struct{}{}
	}()

	<-c
}

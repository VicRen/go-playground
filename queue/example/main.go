package main

import (
	"fmt"
	"time"

	"github.com/VicRen/go-play-ground/queue"
)

func main() {
	q := queue.New(10)

	go func() {
		for {
			time.Sleep(1 * time.Second)

			q.Dispose()
			fmt.Println("putting")
			err := q.Put(20)
			if err != nil {
				fmt.Println("exiting putting")
			}
		}
	}()

	i, err := q.Get(2)
	if err != nil {
		fmt.Println("err Get", err)
		return
	}
	fmt.Println("Get from queue", i)

	err = q.Put(20)
	if err != nil {
		fmt.Println("err Put", err)
		return
	}

	q.Dispose()
	err = q.Put(20)
	if err != nil {
		fmt.Println("err Put", err)
		return
	}
}

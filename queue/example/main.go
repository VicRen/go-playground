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
			fmt.Println("putting")
			err := q.Put(20)
			err = q.Put(20)
			if err != nil {
				fmt.Println("exiting putting")
			}
		}
	}()

	i, err := q.Get(1)
	if err != nil {
		fmt.Println("err Get", err)
		return
	}
	fmt.Println("Get from queue", i)
	q.Dispose()

	i, err = q.Get(1)
	if err != nil {
		fmt.Println("err Get", err, i)
		return
	}
	fmt.Println("Get from queue", i)
}

package main

import (
	"log"
	"time"
)

func main() {
	ch := make(chan struct{})
	go func() {
		close(ch)
	}()
	d, ok := <-ch
	log.Println("chan:", d, ok)
	time.Sleep(1 * time.Second)
}

package main

import (
	"log"
	"time"
)

func main() {
	ch := make(chan []byte, 10)
	go func() {
		ch <- []byte{0x00}
		ch <- []byte{0x01}
		ch <- []byte{0x02}
		close(ch)
	}()
	d, ok := <-ch
	log.Println("chan:", d, ok)
	d, ok = <-ch
	log.Println("chan:", d, ok)
	d, ok = <-ch
	log.Println("chan:", d, ok)
	d, ok = <-ch
	log.Println("chan:", d, ok)
	time.Sleep(1 * time.Second)
}

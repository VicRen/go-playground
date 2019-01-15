package main

import (
	"fmt"
)

func main() {
	done := make(chan struct{}, 2)
	go func() {
		fmt.Println("Series 1")
		done <- struct{}{}
		fmt.Println("Will done go routine")
		done <- struct{}{}
		fmt.Println("Done go routine")
	}()
	fmt.Println("Series 2")
	<- done
	fmt.Println("will Done the whole application")
	<- done
	fmt.Println("Done the whole application")
}



package main

import "fmt"

func main() {
	squares := square(naturals())
	for x := range squares {
		fmt.Println(x)
	}
}

func naturals() chan int {
	naturals := make(chan int)

	go func() {
		for i := 0; i < 10; i++ {
			naturals <- i
		}
		close(naturals)
	}()
	return naturals
}

func square(naturals chan int) chan int {
	squares := make(chan int)
	go func() {
		for {
			x, ok := <-naturals
			if !ok {
				break
			}
			squares <- x * x
		}
		close(squares)
	}()
	return squares
}

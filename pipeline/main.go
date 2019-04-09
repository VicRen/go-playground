package main

import (
	"fmt"
	"math"
	"sync"
	"time"
)

func main() {
	n1()
	n2()
}

func n1() {
	t := time.Now()
	n := naturals(10)
	f := intToFloat64(square(n))
	for x := range squareRoot(f) {
		fmt.Println(x)
	}
	fmt.Println(time.Since(t))
}

func n2() {
	t := time.Now()
	n := naturals(10)
	f := intToFloat64(square(n))
	s1 := squareRoot(f)
	s2 := squareRoot(f)
	s3 := squareRoot(f)
	out := merge(s1, s2, s3)
	for x := range out {
		fmt.Println(x)
	}
	fmt.Println(time.Since(t))
}

func naturals(count int) <-chan int {
	naturals := make(chan int)

	go func() {
		for i := 0; i < count; i++ {
			naturals <- i
		}
		close(naturals)
	}()
	return naturals
}

func square(naturals <-chan int) <-chan int {
	squares := make(chan int)
	go func() {
		defer close(squares)
		for x := range naturals {
			squares <- x * x
		}
	}()
	return squares
}

func intToFloat64(intChan <-chan int) <-chan float64 {
	c := make(chan float64)
	go func() {
		defer close(c)
		for x := range intChan {
			c <- float64(x)
		}
	}()
	return c
}

func squareRoot(number <-chan float64) <-chan int {
	squareRoots := make(chan int)
	go func() {
		defer close(squareRoots)
		for x := range number {
			squareRoots <- int(math.Sqrt(x))
		}
	}()
	return squareRoots
}

func merge(cs ...<-chan int) <-chan int {
	out := make(chan int)

	var wg sync.WaitGroup
	collect := func(in <-chan int) {
		defer wg.Done()
		for n := range in {
			out <- n
		}
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go collect(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

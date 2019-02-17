package main

import "fmt"

func main() {
	defer fmt.Println("Third")
	fmt.Println("First")
	fmt.Println("Second")
}

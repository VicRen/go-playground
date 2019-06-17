package main

import (
	"bytes"
	"fmt"
)

func main() {
	buffer := bytes.NewBuffer([]byte{0x00})
	count := 0
	for {
		buf := make([]byte, 10)
		n, err := buffer.Read(buf)
		if count < 10 {
			buffer.Write(buf[:n])
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("data read:", buf[:n], n, count)
		count++
	}
}

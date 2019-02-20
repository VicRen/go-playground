package main

import (
	md52 "crypto/md5"
	"encoding/hex"
	"fmt"
)

func main() {
	md5 := &MD{"testing"}
	fmt.Printf("md5 is %v, addr is %v\n", md5, &md5)
	m := md52.New()
	b := m.Sum([]byte("hello world"))
	fmt.Printf("md5 is %v\nhash is %v\n", m, hex.EncodeToString(b[:]))
}

type MD struct {
	Text string
}

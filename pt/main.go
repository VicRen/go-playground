package main

import "fmt"

type pt struct {
}

func (p *pt) XXX() {
	fmt.Println(&p)
}

func main() {
	p := &pt{}
	p2 := p
	fmt.Println("------>p", &p, &p2)
	p.XXX()
	p2.XXX()
}

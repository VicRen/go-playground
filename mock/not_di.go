package mock

import (
	"fmt"
	"strconv"
)

func main() {
	a := NewA()
	s := NewSomeStruct(*a)
	fmt.Println(s.SomeMethod())
}

type B struct {
	a A
}

func (s B) SomeMethod() string {
	return strconv.Itoa(s.a.SomeMethod())
}

func NewSomeStruct(a A) *B {
	return &B{a: a}
}

type A struct {
}

func NewA() *A {
	return &A{}
}

func (a A) SomeMethod() int {
	return a.add(3, 4)
}

func (A) add(a, b int) int {
	return a + b
}

package mock

import (
	"fmt"
	"strconv"
)

func DoSomething() {
	dia := NewDiA()
	dib := NewSomeStructDiB(dia)
	fmt.Println(dib.SomeMethod())
}

type DiB struct {
	a DiA
}

func (s DiB) SomeMethod() string {
	ret, err := s.a.SomeMethodDiA()
	if err != nil {
		return ""
	}
	return strconv.Itoa(ret)
}

func NewSomeStructDiB(a DiA) *DiB {
	return &DiB{a: a}
}

type DiA interface {
	SomeMethodDiA() (int, error)
}

type dia struct {

}

func NewDiA() *dia {
	return &dia{}
}

func (a dia) SomeMethodDiA() (int, error) {
	return a.addDiA(3, 4), nil
}

func (dia) addDiA(a, b int) int {
	return a + b
}

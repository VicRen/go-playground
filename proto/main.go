package main

import (
	"fmt"
	"reflect"

	"github.com/VicRen/go-playground/proto/person"
	"github.com/VicRen/go-playground/proto/phone"
	"github.com/VicRen/go-playground/proto/serial"
	"github.com/golang/protobuf/proto"
)

func main() {
	p := &person.Person{
		Name:  "vic",
		Id:    1,
		Email: "vic.ren@vic.ren",
		Phones: []*phone.PhoneNumber{
			{Number: "123", Type: phone.PhoneType_HOME},
		},
	}

	fmt.Println("person:\n", p)

	t, err := proto.Marshal(p)
	if err != nil {
		panic(err)
	}

	fmt.Println("person marshal:", t)

	a, err := serial.MarshalAny(p)
	if err != nil {
		panic(err)
	}

	fmt.Println("person Any:", a, a.Value, a.TypeUrl)

	pback, err := serial.UnmarshalAny(a)
	if err != nil {
		panic(err)
	}

	fmt.Println("person back string:\n", pback)

	if reflect.DeepEqual(pback, p) {
		fmt.Println("p and pback are the same")
	}

	p2 := &person.Person{}
	if err := proto.Unmarshal(t, p2); err != nil {
		panic(err)
	}

	fmt.Println("person2:", p2)
}

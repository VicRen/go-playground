package main

import "fmt"

var vb *Bird

func main() {
	parrot := &Bird{1, "Blue"}
	vb = &Bird{2, "Red"}
	fmt.Printf("原始的Bird:\t\t %+v, \t\t内存地址:%p\n", parrot, &parrot)
	passV(*parrot)
	fmt.Printf("调用后原始的Bird:\t\t %+v, \t\t内存地址:%p\n", parrot, &parrot)


	fmt.Printf("原始的vb:\t\t %+v, \t\t内存地址:%p\n", vb, &vb)
	vb2 := getBird()
	fmt.Printf("原始的vb2:\t\t %+v, \t\t内存地址:%p\n", vb2, &vb2)
	vb2.Name = "Great" + vb2.Name
	fmt.Printf("调用后原始的vb:\t\t %+v, \t\t内存地址:%p\n", vb, &vb)
	fmt.Printf("调用后原始的vb2:\t\t %+v, \t\t内存地址:%p\n", vb2, &vb2)
	vb3 := getBird()
	fmt.Printf("原始的vb3:\t\t %+v, \t\t内存地址:%p\n", vb3, &vb3)
}

type Bird struct {
	Age  int
	Name string
}

func passV(b Bird) {
	b.Age ++
	b.Name = "Great" + b.Name
	fmt.Printf("传入修改后的Bird:\t %+v, \t内存地址:%p\n", b, &b)
}

func getBird() *Bird {
	return vb
}

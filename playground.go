package main

import "fmt"

type Books struct {
	title  string
	author string
	price  int
}

var (
	testing2 = 2
	myBook   = Books{title: "Computer", price: 10}
)

func main() {
	fmt.Println("Hello World")
	number := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	// number := make([]int, 2)
	PrintSlice(number)
	if number == nil {
		fmt.Println("Slice is empty")
	}

	number = append(number, 0)
	PrintSlice(number)

	number = append(number, 1, 2, 3)
	PrintSlice(number)

	for _, num := range number {
		fmt.Println("Slice number:", num)
	}

	kvs := map[string]string{"a": "apple", "b": "banana"}

	for k, v := range kvs {
		fmt.Println("Map value is ", k, v)
	}

	i := 15
	fmt.Printf("%d 的阶乘是 %d\n", i, Factorial(uint64(i)))

	var j uint64
	for j = 1; j < 10; j++ {
		fmt.Printf("%d\t", fibonacci(j))
	}

	var phone Phone
	phone = new(Nokia)
	phone.call()

	phone = new(IPhone)
	phone.call()

	if result, errorMessage := Divide(100, 10); errorMessage == "" {
		fmt.Println("100/10 = ", result)
	}

	if _, errorMessage := Divide(100, 0); errorMessage != "" {
		fmt.Println("errorMessage is " + errorMessage)
	}
}

func PrintSlice(x []int) {
	fmt.Printf("Slice len=%d, cap=%d, slice=%v\n", len(x), cap(x), x)
}

func SayMyName() {
	testing2++
	fmt.Println("Heisenberg")
}

func SayMyNameAgain(myName string) (string, bool) {
	fmt.Println("Your name is " + myName)
	return myName, true
}

func Factorial(number uint64) (result uint64) {
	if number > 0 {
		result = number * Factorial(number-1)
		return result
	}

	return 1
}

func fibonacci(number uint64) (result uint64) {
	if number < 2 {
		return number
	}
	return fibonacci(number-2) + fibonacci(number-1)
}

type Phone interface {
	call()
}

type Nokia struct {
}

func (nokiaPhone Nokia) call() {
	fmt.Println("I'm calling with Nokia Phone")
}

type IPhone struct {
}

func (iPhone IPhone) call() {
	fmt.Println("I'm calling with iPhone")
}

type DivideError struct {
	dividee int
	divider int
}

func (de *DivideError) Error() string {
	strFormat := `
    Cannot proceed, the divider is zero.
    dividee: %d
    divider: 0
`
	return fmt.Sprintf(strFormat, de.dividee)
}

func Divide(varDividee int, varDivider int) (result int, errorMessage string) {
	if (varDivider == 0) {
		dData := DivideError{
			dividee: varDividee,
			divider: varDivider,
		}
		errorMessage = dData.Error()
	} else {
		return varDividee / varDivider, ""
	}
	return
}

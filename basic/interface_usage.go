package main

import "fmt"

type Int int

func (thisVal Int) Less(b Int) bool {
	return thisVal < b
}

func (thisRef *Int) Add(b Int) {
	*thisRef += b
}

type LessAdder interface {
	Less(b Int) bool
	Add(b Int)
}

func main() {
	var a Int = 1
	var ib LessAdder = &a
	ib.Add(2)
	fmt.Println(ib.Less(2))
	fmt.Println(a)

	// 接口查询
	if ia, ok := interface{}(&a).(LessAdder); ok {
		fmt.Println("a is LessAdder")
		ia.Add(4)
	}

	var v Int = 1
	var v1 interface{} = v
	var v2 interface{} = &v1
	fmt.Println(v1)
	fmt.Println(v2)
}

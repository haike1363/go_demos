package main

import (
	"fmt"
)

type st struct {
}

func testLimiter() {

}

func main() {
	var m map[string]interface{}
	if m == nil {
		fmt.Println("m is nil")
	}
	var l []string
	if l == nil {
		fmt.Println(" ls is nil")
	}
	mm := map[string]interface{}{}
	fmt.Println(mm)

	ll := []string{}
	fmt.Println(ll)

}

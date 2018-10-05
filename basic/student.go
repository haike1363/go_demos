package main

import (
	"fmt"
)

type People interface {
	Say(a int)
}

type Student struct {
	Name string
	Year int
}

func (self *Student) Nice() {
	fmt.Println("nice ", self.Name)
}

func (self *Student) Say(a int) {
}

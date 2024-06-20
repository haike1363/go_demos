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

func (thisRef *Student) Nice() {
    fmt.Println("nice ", thisRef.Name)
}

func (thisRef *Student) Say(a int) {
}

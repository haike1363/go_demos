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
    fmt.Println("nice ", ref.Name)
}

func (self *Student) Say(a int) {
}

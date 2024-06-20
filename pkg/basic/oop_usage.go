package main

import (
	"fmt"
)

type Integer int

func (thisRef Integer) Less(b Integer) bool {
    return thisRef < b
}

func (thisRef *Integer) Add(b Integer) {
    *thisRef += b
}

type Rect struct {
    x, y int
    w, h int
}

type Base struct {
    firstName string
    lastName  string
}

func (thisRef *Base) FirstName() string {
    return thisRef.firstName
}
func (thisRef *Base) LastName() string {
    return thisRef.lastName
}

type Drive struct {
    age int
    Base
}

type Drive2 struct {
    *Base
    isMan bool
}

func (thisRef *Drive) FirstName() string {
    return thisRef.Base.FirstName() + " " + thisRef.LastName()
}

func testInherit() {
    base1 := &Base{firstName: "baseFirst", lastName: "baseLast"}
    drive1 := &Drive{Base: Base{"dFirst", "dFirst"}, age: 1}
    drive1.firstName = "driveBaseFirst"
    drive1.lastName = "driveBaseLast"
    fmt.Println(base1.FirstName())
    fmt.Println(drive1.FirstName())
    if true {
        fmt.Println(drive1.LastName())
    }
    drive2 := &Drive2{&Base{"first2", "last2"}, true}

    fmt.Println(drive2)
}

func testStruct() {
    var val Integer = 3
    fmt.Println(val.Less(4))
    val.Add(5)
    fmt.Println(val)

    var a = [3]int{1, 2, 3}
    var b = a // 值传递
    b[0] = 1999
    fmt.Println(a)
    fmt.Println(b)

    var refb = &a
    refb[0] = 1998
    fmt.Println(a)

    // slice, map,  channel, interface 为引用语义
    // 其他的都是值传递语义

    // 返回引用
    rect1 := new(Rect)
    rect2 := rect1
    rect2.h = 1998
    fmt.Println(rect1)

    // 返回值
    rect3 := Rect{x: 0, y: 1, w: 2}
    fmt.Println(rect3)
}

func main() {
    testStruct()
    testInherit()
}

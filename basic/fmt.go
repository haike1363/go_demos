package main

import (
    "fmt"
)

func GetName() (firstName string, lastName string, age int) {
    return "firstName", "", 0
}

func TestString() {
    str := "Hello,世界"
    n := len(str)
    fmt.Println("len: ", n)
    for i := 0; i < n; i++ {
        ch := str[i]
        fmt.Println(i, ch)
    }

    for i, ch := range str {
        fmt.Println(i, ch)
    }
}

func main() {

    TestString()

    fmt.Println(GetName())
    /*
    var bl bool
    var by byte

    var ii int
    var ui uint
    var up uintptr

    var i8 int8
    var i16 int16
    var i32 int32
    var i64 int64
    var f32 float32
    var f64 float64
    var com64 complex64
    var comp128 complex128
    var ss string
    var cc rune
    var err error
    var ptr types.Pointer
    var arr types.Array
    var slice types.Slice
    var mp map[string]int
    var tmp types.Map
    var ch chan
    var tch types.Chan
    var struct1 struct {
        age int
        name string
    }
    var interf interface {

    }


    var scores [10]int
    var part_scores []int
    var student struct {
        age  int
        name string
    }
    _, student.name, student.age = GetName()
    var money *int
    var yearMap map[string]int
    var hoho func(a int) int
        var (
        y1 int
        n1 string
        y2 int
        n2 string
    )

    y1, n1 = year, name

    */

    // 定义变量
    var year int
    var name string
    // 赋值变量
    year = 24
    name = "lilei"
    fmt.Println(year, name)

    // 定义赋值变量
    var hmmLastName string = "meimei"
    // 定义赋值推导变量
    var hmmEge = 23
    // 定义赋值推导变量
    hmmFirstName := "han"
    fmt.Println(hmmEge, hmmFirstName, hmmLastName)
}

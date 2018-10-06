package main

import (
    "fmt"
    "reflect"
)

func myAdd(lastName string, adders ...int) (int, error) {
    sum := 0
    for _, arg := range adders {
        sum += arg
    }
    return sum, nil
}

func myPrintf(args ...interface{}) {

    for _, arg := range args {
        switch arg.(type) {

        case int:
            fmt.Println("int ", arg)
        case string:
            fmt.Println("string ", arg)

        default:
            fmt.Println("type ", reflect.TypeOf(arg))
        }

    }
}

func main() {
    sum, _ := myAdd("lastName", 1, 2, 3)
    fmt.Println(sum)

    myPrintf("abc", 1, 2.3+3i)

    myFun := func(x, y int) int {
        return x*y + 1
    }
    fmt.Println(myFun(2, 3))
    added := func(x, y int) int {
        return x + y
    }(1, 2)
    fmt.Println(added)

    // j的修改作用于闭包调用
    var j int = 5
    funa := func() {
        var i int = 10
        fmt.Println("i, j: ", i, j)
    }
    funa()
    j = 10
    funa()
}

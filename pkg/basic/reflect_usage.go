package main

import (
	"fmt"
    "reflect"
)

func main() {
    var x float64 = 3.14
    var t = reflect.TypeOf(x)
    var v = reflect.ValueOf(x)
    fmt.Println("type:", t)
    fmt.Println("value: ", v)

    p := reflect.ValueOf(&x)
    fmt.Println("type of p: ", p.Type())
    fmt.Println("can set of p: ", p.CanSet())
    v = p.Elem()
    fmt.Println("can set of v: ", v.CanSet())
    v.SetFloat(1.2)
    fmt.Println(x)
}

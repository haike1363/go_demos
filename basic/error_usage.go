package main

import "fmt"

type PathError struct {
    Op string
}

func (self *PathError) Error() string {
    return self.Op
}

func process() {
    defer func() {
        fmt.Println("defer called 1")
    }()
    defer func() {
        fmt.Println("defer called 2")
    }()
    defer func() {
        // 在defer中处理panic
        if r := recover(); r != nil {
            fmt.Println("revocer ", r)
        }
    }()
    panic("panic throw")
}

func main() {
    pathError := &PathError{"path err"}
    fmt.Println(pathError.Error())
    process()
}

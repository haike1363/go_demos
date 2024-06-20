package main

import "fmt"

const Pi float64 = 3.1415926

const (
    size int64 = 1024
    eof        = -1
)

const a, b, c = "a", 1, "c"

const (
    Sunday = iota
    Monday
    Tuesday
    Wednesday
    Thursday
    Friday
    Saturday
    numberOfDays  // 小写变量外部包不可见
)

func main() {

    fmt.Println(Pi, size, eof)
    fmt.Println(a, b, c)
    fmt.Println(Monday, Sunday, Tuesday, Wednesday, Thursday, Friday, Saturday)
}

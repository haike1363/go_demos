package main

import "fmt"

func main() {
    var values [2][3]int

    fmt.Println(values)
    fmt.Println("array len:", len(values))
    fmt.Println("array len:", len(values[0]))

}

package main

import "fmt"

func main() {
    Num := 5
    switch {
    case Num == 0:
        fmt.Println(0)
    case 1 <= Num && Num <= 4:
        fmt.Println(Num)
    case Num == 5:
        fmt.Println(5)
        fallthrough
    case Num == 6:
        fmt.Println("6")
    default:
        fmt.Println("default ", Num)
    }

    sum := 0
    for {
        sum++
        if sum > 10 {
            break
        }
    }

    a := []int{1, 2, 3, 4, 5}
    for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
        a[i], a[j] = a[j], a[i]
    }
    fmt.Println(a)

JLoop:
    for j := 0; j < 5; j++ {
    ILoop:
        for i := 0; i < 10; i++ {
            for k := 0; k < 3; k++ {
                if i == 2 {
                    fmt.Println("break ILoop")
                    break ILoop
                }
                if j == 2 {
                    fmt.Println("break JLoop")
                    break JLoop
                }
                if k == 2 {
                    break
                }
                fmt.Println(j, i, k)
            }
        }
    }
}

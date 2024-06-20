package main

import (
	"fmt"
    "sync"
    "time"
)

func Input(input_chan chan<- int) {
    input_chan <- 1999
}

func Output(output_chan <-chan int) {
    fmt.Println("out put ", <-output_chan)
}

var once sync.Once

func main() {
    once.Do(func() {
        fmt.Println("once exec")
    })
    once.Do(func() {
        fmt.Println("once exec")
    })

    var int_sigal chan int = make(chan int, 1)
    go Input(int_sigal)
    Output(int_sigal)

    var float_signal chan float32 = make(chan float32, 1)

    select {
    case val := <-int_sigal:
        fmt.Println("int signal ", val)
    case val := <-float_signal:
        fmt.Println("float32 signal ", val)
    default:
        fmt.Println("no signals")
    }
    defer func() {
        close(float_signal)
        close(int_sigal)
        if x, ok := <-int_sigal; !ok {
            fmt.Println("not close int_signal ", x)
        }
        fmt.Println("close chans")
    }()


    timeout := make(chan bool, 1)
    go func() {
        time.Sleep(1e9)
        timeout <- true
    }()
    ch := make(chan int, 0)

    select {
    case ch <- 0:
        fmt.Println("send 0 ok")
    case <-timeout:
        fmt.Println("timeout")
    }
}

package main

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"time"
)

func main() {
	pool, _ := ants.NewPool(3)
	for i := 0; i < 4; i++ {
		pool.Submit(func() {
			fmt.Println("sleep")
			time.Sleep(10000 * time.Second)
			fmt.Println("sleep ok")
		})
	}
	fmt.Println("ok")
}

package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	fmt.Println("start")
	runtime.GOMAXPROCS(1)
	mutex := sync.Mutex{}
	// mapVar := map[string]int{}

	for i := 0; i < 10; i++ {
		if i < 9 {
			go func() {
				mutex.Lock()
				fmt.Println("lock it")
				time.Sleep(time.Duration(1) * time.Hour)
			}()
		}else {
			go func() {
				for {
					fmt.Println("work it")
					time.Sleep(time.Duration(3) * time.Second)
					// key := fmt.Sprint(rand.Int() % 100)
					// mapVar[key] = 1
				}
			}()
		}
	}
	time.Sleep(time.Duration(1) * time.Hour)
}

package main

import (
	"fmt"
	"sync"
)

func main() {
	errList := make(chan error)

	var wg sync.WaitGroup
	for i := 0; i< 3;i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errList <- fmt.Errorf("err ")
		}()
	}
	wg.Wait()
	fmt.Println(len(errList))
}

package main

import (
	"fmt"
	"sync"
)

func main() {

	result := make([]int, 10)
	var i = 0
	var wg sync.WaitGroup
	var lock sync.Mutex
	for k := 0; k < len(result); k++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			fmt.Println(index)
			lock.Lock()
			defer lock.Unlock()
			result[index] = index
		}(i)
		i++
	}
	wg.Wait()

	fmt.Println(result)
}

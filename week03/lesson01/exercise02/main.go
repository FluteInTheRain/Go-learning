package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	for i := 1; i < 10; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			fmt.Println(i, i*i)
		}(i)
	}
	wg.Wait()
	fmt.Println("Xong")
}

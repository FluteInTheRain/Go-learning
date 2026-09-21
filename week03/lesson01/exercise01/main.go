package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i <= 10; i++ {
			fmt.Println(i)
		}
	}()
	wg.Wait()
	fmt.Println("Xong")
}

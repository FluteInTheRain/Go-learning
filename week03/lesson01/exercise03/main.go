package main

import (
	"fmt"
	"runtime"
	"sync"
)

func printRange(start int, end int, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := start; i <= end; i++ {
			fmt.Println(i)
		}
	}()
}

func main() {
	var wg sync.WaitGroup
	printRange(1, 3, &wg)
	printRange(4, 6, &wg)
	printRange(7, 10, &wg)
	fmt.Println("Goroutines:", runtime.NumGoroutine())
	wg.Wait()
	fmt.Println("Xong")
}

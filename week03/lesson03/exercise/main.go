package main

import (
	"fmt"
	"sync"
)

func sumEven(n int, out chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	partial := 0
	for i := 2; i <= n; i += 2 {
		partial += i
	}
	out <- partial
}

func sumOdd(n int, out chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	partial := 0
	for i := 1; i <= n; i += 2 {
		partial += i
	}
	out <- partial
}

func main() {
	n := 100
	const workers = 2
	var wg sync.WaitGroup
	results := make(chan int, workers)
	wg.Add(2)
	go sumEven(n, results, &wg)
	go sumOdd(n, results, &wg)
	go func() {
		wg.Wait()
		close(results)
	}()
	total := 0
	for partial := range results {
		total += partial
	}
	fmt.Printf("Tong cac so tu 1 den %d = %d\n", n, total)
	fmt.Printf("Tong cong thuc = %d", n*(n+1)/2)
}

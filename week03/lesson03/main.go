package main

import (
	"fmt"
	"sync"
)

func sumChunk(nums []int, out chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	partial := 0
	for _, n := range nums {
		partial += n
	}
	out <- partial
}

func main() {
	const size = 1_000_000
	nums := make([]int, size)
	for i := 0; i < size; i++ {
		nums[i] = i + 1
	}
	const workers = 4
	chunkSize := (size + workers - 1) / workers
	results := make(chan int, workers)
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		start := w * chunkSize
		end := start + chunkSize
		if end > size {
			end = size
		}
		if start >= end {
			continue
		}
		wg.Add(1)
		go sumChunk(nums[start:end], results, &wg)
	}
	go func() {
		wg.Wait()
		close(results)
	}()
	total := 0
	for partial := range results {
		total += partial
	}
	fmt.Println("Tong song song = ", total)
	fmt.Println("Tong cong thuc = ", size*(size+1)/2)
}

package main

import "fmt"

func countEven(numbers []int) int {
	count := 0
	for _, number := range numbers {
		if number%2 == 0 {
			count++
		}
	}
	return count
}

func main() {
	numbers := []int{3, 8, 12, 7, 4, 9, 10}
	count_even := countEven(numbers)
	fmt.Printf("Array: %v\n", numbers)
	fmt.Print("So luong so chan: ", count_even)
}

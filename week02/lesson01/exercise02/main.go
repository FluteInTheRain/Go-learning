package main

import (
	"fmt"
)

func removeAt(s []int, i int) []int {
	return append(s[:i], s[i+1:]...)
}

func main() {
	numbers := []int{10, 20, 30, 40, 50}
	fmt.Println(numbers)
	numbers = removeAt(numbers, 2)
	fmt.Printf("After remove: %v\n", numbers)
}

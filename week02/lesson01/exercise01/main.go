package main

import (
	"fmt"
	"slices"
)

func main() {
	numbers := []int{3, 1, 4, 1, 5, 9, 2, 6}
	slices.Sort(numbers)
	fmt.Println(numbers)
	fmt.Printf("Len: %d\n", len(numbers))
	fmt.Printf("Cap: %d\n", cap(numbers))
}

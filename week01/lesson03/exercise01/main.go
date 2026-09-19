package main

import "fmt"

func isOdd(n int) bool {
	return n%2 != 0
}

func main() {
	var number int
	fmt.Scan(&number)
	fmt.Printf("La so le: %t", isOdd(number))
}

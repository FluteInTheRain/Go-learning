package main

import "fmt"

func main() {
	var a int
	var b int
	var op string
	fmt.Scanln(&a, &b)
	fmt.Scanln(&op)
	switch op {
	case "+":
		fmt.Printf("%d + %d = %d\n", a, b, a+b)
	case "-":
		fmt.Printf("%d - %d = %d\n", a, b, a-b)
	case "*":
		fmt.Printf("%d * %d = %d\n", a, b, a*b)
	case "/":
		if b == 0 {
			fmt.Println("Khong the chia cho 0")
		} else {
			fmt.Printf("%d / %d = %d\n", a, b, a/b)
		}
	default:
		fmt.Println("Toan tu khong hop le!")
	}
}

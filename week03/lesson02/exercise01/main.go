package main

import (
	"demo/geometry"
	"fmt"
)

func main() {
	radius := 5.0
	width, height := 4.0, 6.0

	// Use exported functions from the geometry package
	fmt.Printf("Diện tích hình tròn (r=%.1f): %.2f\n", radius, geometry.CircleArea(radius))
	fmt.Printf("Chu vi hình tròn (r=%.1f): %.2f\n", radius, geometry.CirclePerimeter(radius))
	fmt.Printf("Diện tích hình chữ nhật: %.2f\n", geometry.RectangleArea(width, height))
	fmt.Printf("Hằng số Pi trong package: %.5f\n", geometry.Pi)

	// Handle the error-returning function properly
	area, err := geometry.SafeCircleArea(-1)
	if err != nil {
		fmt.Println("Lỗi:", err)
	} else {
		fmt.Printf("Diện tích an toàn: %.2f\n", area)
	}
}

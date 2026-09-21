package geometry

import (
	"errors"
	"math"
)

const Pi = 3.14159265358979

var ErrNegativeInput = errors.New("geometry: dimension must be positive")

func CircleArea(radius float64) float64 {
	return math.Pi * radius * radius
}

func CirclePerimeter(radius float64) float64 {
	return 2 * math.Pi * radius
}

func RectangleArea(width, height float64) float64 {
	return width * height
}

func RectanglePerimeter(width, height float64) float64 {
	return 2 * (width + height)
}

func validatePositive(value float64) bool {
	return value > 0
}

func TriangleArea(base, height float64) float64 {
	if !validatePositive(base) || !validatePositive(height) {
		return 0 // invalid input returns 0
	}
	return (base * height) / 2
}

func SafeCircleArea(radius float64) (float64, error) {
	if radius <= 0 {
		return 0, ErrNegativeInput
	}
	return math.Pi * radius * radius, nil
}

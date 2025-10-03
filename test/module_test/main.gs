package main

import "math"

func main() {
	// Test math module functions
	result1 := math.add(3, 4)
	result2 := math.multiply(result1, 6)
	result3 := math.square(result2)
	result4 := math.fibonacci(3)
	
	return result4 + result3
}
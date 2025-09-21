package main

func main() {
	// Test range statement with slice
	nums := []int{1, 2, 3, 4, 5}
	
	// Simple range - just iterate over elements
	for range nums {
		// This is a simplified range statement
		// In a full implementation, we would capture the value
	}
	
	// Another range test
	items := []int{10, 20, 30}
	sum := 0
	for range items {
		sum += 1  // Count items
	}
	
	return sum
}
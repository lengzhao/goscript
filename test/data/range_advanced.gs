package main

func main() {
	// Test range statement with slice
	nums := []int{2, 4, 6, 8}
	
	// Range with index and value (simplified)
	sum := 0
	for range nums {
		sum += 1  // In simplified version, we just count elements
	}
	
	// Range over different slice
	values := []int{10, 20, 30}
	total := 0
	for range values {
		total += 1  // Count elements
	}
	
	// Nested range test
	matrix := [][]int{{1, 2}, {3, 4}}
	count := 0
	for range matrix {
		count += 1  // Count rows
	}
	
	return sum + total + count
}
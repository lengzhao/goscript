package main

func main() {
	// Test various composite literals
	items := []int{1, 2, 3, 4, 5}
	
	// Test range with different types
	strings := []string{"hello", "world"}
	count := 0
	for _, s := range strings {
		count += len(s)
	}
	
	// Test nested operations
	sum := 0
	for _, it := range items {
		sum += it
	}
	
	return sum + count
}
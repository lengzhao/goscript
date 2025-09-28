package main

func main() {
	// Test slice literal
	items := []int{1, 2, 3}
	
	// Test range with slice literal
	sum := 0
	for _, it := range items {
		sum += it
	}
	
	return sum
}
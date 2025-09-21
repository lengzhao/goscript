package main

func main() {
	// Simple range test
	items := []int{1, 2, 3}
	count := 0
	
	// Simple range statement
	for range items {
		count++
	}
	
	return count
}
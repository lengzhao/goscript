package main

func main() {
	// Simple range test
	items := []int{1, 2, 3}
	count := 0
	
	// Simple range statement
	for it := range items {
		count+=it
	}
	
	return count
}
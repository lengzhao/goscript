package test

import (
	"testing"

	goscript "github.com/lengzhao/goscript"
)

func TestRangeSlice(t *testing.T) {
	scriptSource := `package main

func main() {
	// Create a slice
	slice := []int{1, 2, 3, 4, 5}
	
	// Sum all elements using range
	sum := 0
	for _, value := range slice {
		sum += value
	}
	
	return sum  // Should return 15
}`

	script := goscript.NewScript([]byte(scriptSource))
	result, err := script.Run()
	if err != nil {
		t.Fatalf("Failed to run script: %v", err)
	}

	if result != 15 {
		t.Errorf("Expected 15, got %v", result)
	}
}

func TestRangeSliceWithIndex(t *testing.T) {
	scriptSource := `package main

func main() {
	// Create a slice
	slice := []int{10, 20, 30}
	
	// Sum indices and values using range
	sum := 0
	for index, value := range slice {
		sum += index + value
	}
	
	return sum  // Should return (0+10) + (1+20) + (2+30) = 10 + 21 + 32 = 63
}`

	script := goscript.NewScript([]byte(scriptSource))
	result, err := script.Run()
	if err != nil {
		t.Fatalf("Failed to run script: %v", err)
	}

	if result != 63 {
		t.Errorf("Expected 63, got %v", result)
	}
}

func TestRangeString(t *testing.T) {
	scriptSource := `package main

func main() {
	// Create a string
	str := "abc"
	
	// Count characters using range
	count := 0
	for range str {
		count++
	}
	
	return count  // Should return 3
}`

	script := goscript.NewScript([]byte(scriptSource))
	result, err := script.Run()
	if err != nil {
		t.Fatalf("Failed to run script: %v", err)
	}

	if result != 3 {
		t.Errorf("Expected 3, got %v", result)
	}
}

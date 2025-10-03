package test

import (
	"fmt"
	"testing"

	goscript "github.com/lengzhao/goscript"
)

func TestFor(t *testing.T) {
	// Test nested loops
	script := goscript.NewScript([]byte(`
	package main

	func multiplyTables() int {
		result := 0
		for i := 1; i <= 2; i++ {
			for j := 1; j <= 2; j++ {
				result += i * j
			}
		}
		return result
	}

	func main() {
		result := multiplyTables()  // (1*1 + 1*2) + (2*1 + 2*2) = 1 + 2 + 2 + 4 = 9
		return result
	}
	`))

	// Enable debug mode to see instructions
	script.SetDebug(true)

	result, err := script.Run()
	if err != nil {
		fmt.Printf("Failed to run script: %v\n", err)
		return
	}

	fmt.Printf("Result: %v\n", result)
	if result != 9 {
		t.Errorf("Expected result 9, got %v", result)
	}
}

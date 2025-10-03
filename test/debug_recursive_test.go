package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/lengzhao/goscript"
)

func TestRecursiveDebug(t *testing.T) {
	// Create a simple recursive function to debug
	source := `package main

func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n-1)
}

func main() {
    result := factorial(3)  // Should be 6
    return result
}`

	script := goscript.NewScript([]byte(source))
	script.SetDebug(true) // Enable debug mode

	// Get the VM to inspect instructions
	vm := script.GetVM()

	// Execute the script
	result, err := script.RunContext(context.Background())
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	// Print all instructions for debugging
	fmt.Println("Instructions:")
	instructions := vm.GetInstructions()
	for i, instr := range instructions {
		fmt.Printf("%d: %s\n", i, instr.String())
	}

	// Check the result
	expected := 6 // factorial(3) = 3 * 2 * 1 = 6
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

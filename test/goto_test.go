package test

import (
	"testing"

	"github.com/lengzhao/goscript"
)

func TestGotoStatement(t *testing.T) {
	// Test script with goto statement
	script := `
package main

func main() {
	i := 0
	goto label1
	
	// This code should be skipped
	i = 100
	
label1:
	// This code should be executed
	i = 42
	
	return i
}
`

	// Create a new GoScript VM
	s := goscript.NewScript([]byte(script))

	// Build the script
	err := s.Build()
	if err != nil {
		t.Fatalf("Failed to build script: %v", err)
	}

	// Execute the script
	result, err := s.Run()
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	// Check the result
	if result != 42 {
		t.Errorf("Expected result to be 42, got %v", result)
	}
}

func TestGotoLoop(t *testing.T) {
	// Test script with goto used for looping
	script := `
package main

func main() {
	sum := 0
	i := 0
	
loop:
	if i >= 5 {
		goto end
	}
	sum = sum + i
	i = i + 1
	goto loop
	
end:
	return sum
}
`

	// Create a new GoScript VM
	s := goscript.NewScript([]byte(script))

	// Build the script
	err := s.Build()
	if err != nil {
		t.Fatalf("Failed to build script: %v", err)
	}

	// Execute the script
	result, err := s.Run()
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	// Check the result (0+1+2+3+4 = 10)
	if result != 10 {
		t.Errorf("Expected result to be 10, got %v", result)
	}
}

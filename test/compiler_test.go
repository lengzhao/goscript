package test

import (
	"testing"

	"github.com/lengzhao/goscript"
)

func TestCompilerBasic(t *testing.T) {
	// Test basic arithmetic
	script := goscript.NewScript([]byte(`
	package main

	func main() {
		return 1 + 2
	}
	`))

	result, err := script.Run()
	if err != nil {
		t.Fatalf("Failed to run script: %v", err)
	}

	if result != 3 {
		t.Errorf("Expected 3, got %v", result)
	}
}

func TestCompilerVariables(t *testing.T) {
	// Test variable assignment
	script := goscript.NewScript([]byte(`
	package main

	func main() {
		x := 10
		y := 20
		return x + y
	}
	`))

	result, err := script.Run()
	if err != nil {
		t.Fatalf("Failed to run script: %v", err)
	}

	if result != 30 {
		t.Errorf("Expected 30, got %v", result)
	}
}

func TestCompilerFunctionCall(t *testing.T) {
	// Test function call
	script := goscript.NewScript([]byte(`
	package main

	func main() {
		return add(5, 3)
	}
	`))

	// Add the add function
	script.AddFunction("add", goscript.NewSimpleFunction("add", func(args ...interface{}) (interface{}, error) {
		a := args[0].(int)
		b := args[1].(int)
		return a + b, nil
	}))

	result, err := script.Run()
	if err != nil {
		t.Fatalf("Failed to run script: %v", err)
	}

	if result != 8 {
		t.Errorf("Expected 8, got %v", result)
	}
}

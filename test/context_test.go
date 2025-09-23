package test

import (
	"context"
	"testing"

	"github.com/lengzhao/goscript"
)

func TestContextManagement(t *testing.T) {
	// Create a new script with a simple function that demonstrates context management
	source := `
package main

var globalVar = "global"

func main() {
	x := 10
	y := 20
	
	// Nested block to test scope management
	{
		z := 30
		x = x + z  // x should be accessible from parent scope
	}
	
	// Function call to test function scope
	result := add(x, y)
	return result
}

func add(a, b int) int {
	return a + b
}
`

	script := goscript.NewScript([]byte(source))

	// Execute the script
	result, err := script.RunContext(context.Background())
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	// Check the result
	expected := 60 // (10 + 30) + 20
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}

	// Check that context management worked correctly
	// This would involve inspecting the VM's context map and stack
	// For now, we're just verifying that the script executes without errors
}

func TestContextKeyGeneration(t *testing.T) {
	// Test that the compiler generates correct scope keys
	t.Skip("Skipping")
	source := `
package main

func main() {
	x := 10
	
	func() {
		y := 20
		x = x + y
	}()
	
	return x
}

type Calculator struct{}

func (c *Calculator) Add(a, b int) int {
	return a + b
}

func (c Calculator) Multiply(a, b int) int {
	return a * b
}
`

	script := goscript.NewScript([]byte(source))

	// TODO: Add instruction inspection when we have access to compiled instructions

	// Execute the script to verify it works
	result, err := script.RunContext(context.Background())
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	expected := 30 // 10 + 20
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

func TestVariableIsolation(t *testing.T) {
	// t.Skip("Skipping test until we have access to the execution context")
	// Test that variables in different scopes are properly isolated
	source := `
package main

var globalVar = "global"

func main() {
	x := "main"
	
	{
		x := "block"  // This should shadow the main scope variable
		globalVar = x // This should modify the global variable
	}
	
	return x  // Should return "main", not "block"
}
`

	script := goscript.NewScript([]byte(source))
	script.SetDebug(true) // Enable debug mode

	result, err := script.RunContext(context.Background())
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	// Check that the main scope variable was not affected by the block scope
	if result != "main" {
		t.Errorf("Expected 'main', got '%v'", result)
	}

	// TODO: Add global variable inspection when we have access to the execution context
}

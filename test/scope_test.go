package test

import (
	"testing"

	"github.com/lengzhao/goscript/compiler"
	"github.com/lengzhao/goscript/context"
	"github.com/lengzhao/goscript/parser"
	"github.com/lengzhao/goscript/vm"
)

func TestScopeManagement(t *testing.T) {
	// Create a simple script with nested scopes
	script := `
package main

func main() {
	x := 10
	{
		x := 20
		{
			x := 30
			return x
		}
		// This should not be reached
		return x
	}
	// This should not be reached
	return x
}
`

	// Parse the script
	parser := parser.New()
	astFile, err := parser.Parse("test.go", []byte(script), 0)
	if err != nil {
		t.Fatalf("Failed to parse script: %v", err)
	}

	// Create VM and context
	vmInstance := vm.NewVM()
	context := context.NewExecutionContext()

	// Create compiler and compile
	compiler := compiler.NewCompiler(vmInstance, context)
	err = compiler.Compile(astFile)
	if err != nil {
		t.Fatalf("Failed to compile script: %v", err)
	}

	// Execute the VM
	result, err := vmInstance.Execute(nil)
	if err != nil {
		t.Fatalf("Failed to execute VM: %v", err)
	}

	// Check result: should be 30 (innermost scope)
	if result != 30 {
		t.Errorf("Expected result 30, got %v", result)
	}
}

func TestFunctionScope(t *testing.T) {
	// Create a script with function scopes
	script := `
package main

func add(a, b int) int {
	// Local variables in function scope
	x := a
	y := b
	return x + y
}

func main() {
	result := add(5, 3)
	return result
}
`

	// Parse the script
	parser := parser.New()
	astFile, err := parser.Parse("test.go", []byte(script), 0)
	if err != nil {
		t.Fatalf("Failed to parse script: %v", err)
	}

	// Create VM and context
	vmInstance := vm.NewVM()
	context := context.NewExecutionContext()

	// Create compiler and compile
	compiler := compiler.NewCompiler(vmInstance, context)
	err = compiler.Compile(astFile)
	if err != nil {
		t.Fatalf("Failed to compile script: %v", err)
	}

	// Execute the VM
	result, err := vmInstance.Execute(nil)
	if err != nil {
		t.Fatalf("Failed to execute VM: %v", err)
	}

	// Check result: should be 8 (5 + 3)
	if result != 8 {
		t.Errorf("Expected result 8, got %v", result)
	}
}

func TestTypeScope(t *testing.T) {
	// Create a script with type scopes
	script := `
package main

type Person struct {
	name string
	age  int
}

func main() {
	// Create a new person using struct literal
	person := Person{name: "Alice", age: 30}
	
	// Access fields
	return person.age
}
`

	// Parse the script
	parser := parser.New()
	astFile, err := parser.Parse("test.go", []byte(script), 0)
	if err != nil {
		t.Fatalf("Failed to parse script: %v", err)
	}

	// Create VM and context
	vmInstance := vm.NewVM()
	context := context.NewExecutionContext()

	// Create compiler and compile
	compiler := compiler.NewCompiler(vmInstance, context)
	err = compiler.Compile(astFile)
	if err != nil {
		t.Fatalf("Failed to compile script: %v", err)
	}

	// Execute the VM
	result, err := vmInstance.Execute(nil)
	if err != nil {
		t.Fatalf("Failed to execute VM: %v", err)
	}

	// Check result: person.age should be 30
	if result != 30 {
		t.Errorf("Expected result 30, got %v", result)
	}
}

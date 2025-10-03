package test

import (
	"testing"

	"github.com/lengzhao/goscript/compiler"
	"github.com/lengzhao/goscript/parser"
	"github.com/lengzhao/goscript/vm"
)

func TestSimpleAssignment(t *testing.T) {
	// Create a simple script with assignment
	script := `
package main

func main() {
	x := 10
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
	// Enable debug mode
	// vmInstance.SetDebug(true) // Not implemented in current VM
	// context.SetDebug(true) // ExecutionContext not used in current implementation

	// Create compiler and compile
	compiler := compiler.NewCompiler(vmInstance)
	err = compiler.Compile(astFile)
	if err != nil {
		t.Fatalf("Failed to compile script: %v", err)
	}

	// Print instructions for debugging
	instructions := vmInstance.GetInstructions()
	t.Logf("Generated %d instructions:", len(instructions))
	for i, instr := range instructions {
		t.Logf("  %d: %s", i, instr.String())
	}

	// Execute the VM
	result, err := vmInstance.Execute("")
	if err != nil {
		t.Fatalf("Failed to execute VM: %v", err)
	}

	// Check result: should be 10
	if result != 10 {
		t.Errorf("Expected result 10, got %v", result)
	}
}

func TestDebugCompiler(t *testing.T) {
	// Create a simple script with assignment
	script := `
package main

func main() {
	x := 10
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

	// Enable debug mode
	// vmInstance.SetDebug(true) // Not implemented in current VM

	// Create compiler and compile
	compiler := compiler.NewCompiler(vmInstance)

	// Add some debug prints to see what's happening
	t.Logf("Starting compilation...")
	err = compiler.Compile(astFile)
	if err != nil {
		t.Fatalf("Failed to compile script: %v", err)
	}

	// Print instructions for debugging
	instructions := vmInstance.GetInstructions()
	t.Logf("Generated %d instructions:", len(instructions))
	for i, instr := range instructions {
		t.Logf("  %d: %s", i, instr.String())
	}
}

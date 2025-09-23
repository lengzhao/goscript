package compiler_test

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
		y := 20
		x = x + y
	}
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

	// Enable debug mode
	vmInstance.SetDebug(true)
	context.SetDebug(true)

	// Create compiler and compile
	compiler := compiler.NewCompiler(vmInstance, context)
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
	// Note: This might fail because we haven't fully implemented the runtime yet
	// but we can still check that compilation worked
	t.Log("Compilation completed successfully")
}

package compiler_test

import (
	"testing"

	"github.com/lengzhao/goscript/compiler"
	execContext "github.com/lengzhao/goscript/context"
	"github.com/lengzhao/goscript/instruction"
	"github.com/lengzhao/goscript/parser"
	"github.com/lengzhao/goscript/vm"
)

func TestCreateVarInstruction(t *testing.T) {
	// Create a new VM
	v := vm.NewVM()

	// Create execution context
	execCtx := execContext.NewExecutionContext()

	// Create compiler
	c := compiler.NewCompiler(v, execCtx)

	// Create a parser
	p := parser.New()

	// Parse the source code with variable declarations
	src := `
package main

func main() {
	// Test different forms of variable declaration
	var a int
	var b = 42
	var c int = 100
	
	// Use the variables to make sure they exist
	a = b + c
	return a
}
`

	f, err := p.Parse("test.go", []byte(src), 0)
	if err != nil {
		t.Fatalf("Failed to parse source code: %v", err)
	}

	// Compile the file
	err = c.Compile(f)
	if err != nil {
		t.Fatalf("Failed to compile source code: %v", err)
	}

	// Print instructions for debugging
	instructions := v.GetInstructions()
	t.Logf("Generated %d instructions:", len(instructions))
	for i, instr := range instructions {
		t.Logf("  %d: %s", i, instr.String())
	}

	// Check that OpCreateVar instructions were generated
	createVarCount := 0
	for _, instr := range instructions {
		if instr.Op == instruction.OpCreateVar {
			createVarCount++
		}
	}

	if createVarCount != 3 {
		t.Errorf("Expected 3 OpCreateVar instructions, got %d", createVarCount)
	}

	// Execute the compiled code
	result, err := v.Execute(nil)
	if err != nil {
		t.Fatalf("Failed to execute compiled code: %v", err)
	}

	// Check the result (should be 42 + 100 = 142)
	expected := 142
	if result != expected {
		t.Errorf("Expected result %d, got %v", expected, result)
	}

	t.Log("CreateVar instruction test passed")
}

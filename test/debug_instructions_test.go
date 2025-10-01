package test

import (
	"testing"

	"github.com/lengzhao/goscript/compiler"
	"github.com/lengzhao/goscript/parser"
	"github.com/lengzhao/goscript/vm"
)

func TestDebugInstructions(t *testing.T) {
	// Simple function to test
	code := `package test

func add(a, b) {
	return a + b
}`

	// Create a VM
	testVM := vm.NewVM()

	// Parse and compile the code
	parserInstance := parser.New()
	astFile, err := parserInstance.Parse("test.gs", []byte(code), 0)
	if err != nil {
		t.Fatalf("Failed to parse code: %v", err)
	}

	// Compile the code
	compilerInstance := compiler.NewCompiler(testVM)
	err = compilerInstance.Compile(astFile)
	if err != nil {
		t.Fatalf("Failed to compile code: %v", err)
	}

	// Print all instruction sets
	instructionSets := testVM.GetInstructionSets()
	for key, instructions := range instructionSets {
		t.Logf("Instructions for key %s:", key)
		for i, instr := range instructions {
			t.Logf("  %d: %s", i, instr.String())
		}
	}
}

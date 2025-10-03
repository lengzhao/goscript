package test

import (
	"fmt"
	"testing"

	"github.com/lengzhao/goscript"
	"github.com/lengzhao/goscript/compiler"
	"github.com/lengzhao/goscript/parser"
	"github.com/lengzhao/goscript/vm"
)

func TestSwitchWithGoto(t *testing.T) {
	// Test script with switch statement
	script := `package main

func main() {
	x := 2
	result := 0
	
	switch x {
	case 1:
		result = 10
	case 2:
		result = 20
	case 3:
		result = 30
	default:
		result = 0
	}
	
	return result
}`

	// Create a parser
	p := parser.New()

	// Parse the source code into an AST
	astFile, err := p.Parse("test.go", []byte(script), 0)
	if err != nil {
		t.Fatalf("Failed to parse source code: %v", err)
	}

	// Create a VM and compiler
	vm := vm.NewVM()
	c := compiler.NewCompiler(vm)

	// Compile the AST to bytecode
	err = c.Compile(astFile)
	if err != nil {
		t.Fatalf("Failed to compile AST: %v", err)
	}

	// Print instructions for debugging
	instructions := vm.GetAllInstructionSets()
	for key, instrs := range instructions {
		fmt.Printf("Instructions for key: %s\n", key)
		for i, instr := range instrs {
			fmt.Printf("  %d: %s\n", i, instr.String())
		}
	}

	// Execute the script
	result, err := vm.Execute("")
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	// Check the result
	if result != 20 {
		t.Errorf("Expected result to be 20, got %v", result)
	}
}

func TestSwitchWithDefault(t *testing.T) {
	// Test script with switch statement that hits default case
	script := `package main

func main() {
	x := 5
	result := 0
	
	switch x {
	case 1:
		result = 10
	case 2:
		result = 20
	case 3:
		result = 30
	default:
		result = 100
	}
	
	return result
}`

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
	if result != 100 {
		t.Errorf("Expected result to be 100, got %v", result)
	}
}

func TestSwitchWithoutDefault(t *testing.T) {
	// Test script with switch statement that has no default case
	script := `package main

func main() {
	x := 5
	result := 0
	
	switch x {
	case 1:
		result = 10
	case 2:
		result = 20
	case 3:
		result = 30
	}
	
	return result
}`

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

	// Check the result (should be 0 as no case matches and no assignment happens)
	if result != 0 {
		t.Errorf("Expected result to be 0, got %v", result)
	}
}

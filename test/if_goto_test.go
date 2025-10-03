package test

import (
	"fmt"
	"testing"

	"github.com/lengzhao/goscript/compiler"
	"github.com/lengzhao/goscript/parser"
	"github.com/lengzhao/goscript/vm"
)

func TestIfWithGoto(t *testing.T) {
	// Test script with if statement
	script := `package main

func main() {
	x := 5
	result := 0
	
	if x > 3 {
		result = 10
	} else {
		result = 20
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
	vmInstance := vm.NewVM()
	c := compiler.NewCompiler(vmInstance)

	// Compile the AST to bytecode
	err = c.Compile(astFile)
	if err != nil {
		t.Fatalf("Failed to compile AST: %v", err)
	}

	// Print instructions for debugging
	instructions := vmInstance.GetAllInstructionSets()
	for key, instrs := range instructions {
		fmt.Printf("Instructions for key: %s\n", key)
		for i, instr := range instrs {
			fmt.Printf("  %d: %s\n", i, instr.String())
		}
	}

	// Execute the script
	result, err := vmInstance.Execute("")
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	// Check the result
	if result != 10 {
		t.Errorf("Expected result to be 10, got %v", result)
	}
}

func TestIfWithoutElse(t *testing.T) {
	// Test script with if statement but no else
	script := `package main

func main() {
	x := 5
	result := 0
	
	if x > 3 {
		result = 10
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
	vmInstance := vm.NewVM()
	c := compiler.NewCompiler(vmInstance)

	// Compile the AST to bytecode
	err = c.Compile(astFile)
	if err != nil {
		t.Fatalf("Failed to compile AST: %v", err)
	}

	// Execute the script
	result, err := vmInstance.Execute("")
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	// Check the result
	if result != 10 {
		t.Errorf("Expected result to be 10, got %v", result)
	}
}

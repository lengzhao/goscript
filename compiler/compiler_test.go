package compiler

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/lengzhao/goscript/vm"
)

func TestCompilerWithKeyBasedInstructions(t *testing.T) {
	// Create a new VM
	vm := vm.NewVM()

	// Create a new compiler
	compiler := NewCompiler(vm)

	// Simple test code with a function call
	code := `
package main

func add(a, b int) int {
	return a + b
}

func main() {
	result := add(1, 2)
}
`

	// Parse the code
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, "", code, 0)
	if err != nil {
		t.Fatalf("Failed to parse code: %v", err)
	}

	// Compile the code
	err = compiler.Compile(astFile)
	if err != nil {
		t.Fatalf("Failed to compile code: %v", err)
	}

	// Check that instructions were generated
	instructions := vm.GetInstructions()
	if len(instructions) == 0 {
		t.Error("No instructions were generated")
	}

	t.Logf("Generated %d instructions", len(instructions))
}

func TestCompilerWithGlobalVariables(t *testing.T) {
	// Create a new VM
	vm := vm.NewVM()

	// Create a new compiler
	compiler := NewCompiler(vm)

	// Test code with global variables
	code := `
package testpkg

var globalVar = 42

func GetGlobal() int {
	return globalVar
}

func main() {
	localVar := globalVar + 1
}
`

	// Parse the code
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, "", code, 0)
	if err != nil {
		t.Fatalf("Failed to parse code: %v", err)
	}

	// Compile the code
	err = compiler.Compile(astFile)
	if err != nil {
		t.Fatalf("Failed to compile code: %v", err)
	}

	// Check that instructions were generated
	instructions := vm.GetInstructions()
	if len(instructions) == 0 {
		t.Error("No instructions were generated")
	}

	t.Logf("Generated %d instructions for global variables test", len(instructions))
}

func TestCompilerWithCustomPackageName(t *testing.T) {
	// Create a new VM
	vm := vm.NewVM()

	// Create a new compiler
	compiler := NewCompiler(vm)

	// Test code with custom package name
	code := `
package mypackage

func Calculate(x, y int) int {
	return x + y
}

func main() {
	result := Calculate(10, 20)
}
`

	// Parse the code
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, "", code, 0)
	if err != nil {
		t.Fatalf("Failed to parse code: %v", err)
	}

	// Compile the code
	err = compiler.Compile(astFile)
	if err != nil {
		t.Fatalf("Failed to compile code: %v", err)
	}

	// Check that instructions were generated
	instructions := vm.GetInstructions()
	if len(instructions) == 0 {
		t.Error("No instructions were generated")
	}

	t.Logf("Generated %d instructions for custom package name test", len(instructions))
}

package test

import (
	"testing"

	"github.com/lengzhao/goscript/compiler"
	"github.com/lengzhao/goscript/parser"
	"github.com/lengzhao/goscript/vm"
)

func TestSimpleGSModule(t *testing.T) {
	// Create a simple module
	moduleCode := `package math

func add(a, b) {
	return a + b
}`

	// Create a VM for the module
	moduleVM := vm.NewVM()

	// Parse and compile the module
	parserInstance := parser.New()
	astFile, err := parserInstance.Parse("math.gs", []byte(moduleCode), 0)
	if err != nil {
		t.Fatalf("Failed to parse module: %v", err)
	}

	// Compile the module
	compilerInstance := compiler.NewCompiler(moduleVM)
	err = compilerInstance.Compile(astFile)
	if err != nil {
		t.Fatalf("Failed to compile module: %v", err)
	}

	// Test calling the function directly through the module VM
	result, err := moduleVM.Execute("math.func.add", 3, 4)
	if err != nil {
		t.Fatalf("Failed to execute function: %v", err)
	}

	if result != 7 {
		t.Errorf("Expected 7, got %v", result)
	}
}

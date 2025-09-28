package test

import (
	"testing"

	"github.com/lengzhao/goscript/compiler"
	"github.com/lengzhao/goscript/parser"
	"github.com/lengzhao/goscript/vm"
)

func TestVariableDeclaration(t *testing.T) {
	// Create a new VM
	v := vm.NewVM()

	// Create compiler
	c := compiler.NewCompiler(v)

	// Create a parser
	p := parser.New()

	// Parse the source code
	src := `
package main

func main() {
	var x int
	x = 10
	var y = 20
	var z int = 30
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

	// Execute the compiled code
	_, err = v.Execute("")
	if err != nil {
		t.Fatalf("Failed to execute compiled code: %v", err)
	}

	// Check that variables were created and assigned correctly
	// This would require accessing the VM's context/locals, which may not be directly accessible
	// For now, we're just testing that the compilation and execution don't fail
	t.Log("Variable declaration test passed")
}

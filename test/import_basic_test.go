package test

import (
	"testing"

	"github.com/lengzhao/goscript/compiler"
	"github.com/lengzhao/goscript/parser"
	"github.com/lengzhao/goscript/vm"
)

// TestBasicImport tests basic import functionality
func TestBasicImport(t *testing.T) {
	// Create a new VM
	vmInstance := vm.NewVM()
	vmInstance.SetDebug(false) // Disable debug output for cleaner test logs

	// Create a new compiler
	compiler := compiler.NewCompiler(vmInstance)

	// Parse the test script
	parser := parser.New()
	script := `
package test

import "strings"

func main() {
    s := "Hello, World!"
    result := strings.Contains(s, "World")
    return result
}
`

	ast, err := parser.Parse("test.gs", []byte(script), 0)
	if err != nil {
		t.Fatalf("Failed to parse script: %v", err)
	}

	// Compile the AST
	err = compiler.Compile(ast)
	if err != nil {
		t.Fatalf("Failed to compile script: %v", err)
	}

	// Check that the imported module is tracked
	// This is an internal check - we can't directly access the compiler's importedModules
	// but we can verify that the compilation succeeded

	// Execute the compiled code
	result, err := vmInstance.Execute("")
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	// Expected result: true (strings.Contains("Hello, World!", "World") should return true)
	expected := true
	if result != expected {
		t.Errorf("Expected result to be %v, got %v", expected, result)
	}
}

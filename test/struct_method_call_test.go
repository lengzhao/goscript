package test

import (
	"testing"

	"github.com/lengzhao/goscript/compiler"
	execContext "github.com/lengzhao/goscript/context"
	"github.com/lengzhao/goscript/parser"
	"github.com/lengzhao/goscript/vm"
)

// TestStructMethodCall tests struct method calls
func TestStructMethodCall(t *testing.T) {
	// Create a new VM
	vmInstance := vm.NewVM()
	vmInstance.SetDebug(true) // Enable debug output

	// Create a new execution context
	execCtx := execContext.NewExecutionContext()

	// Create a new compiler
	compiler := compiler.NewCompiler(vmInstance, execCtx)

	// Parse the test script
	parser := parser.New()
	script := `
package test

type Rectangle struct {
	width  int
	height int
}

func (r Rectangle) Area() int {
	return r.width * r.height
}

// Value receiver - modifications should be ineffective
func (r Rectangle) SetWidth(width int) {
	r.width = width
}

// Pointer receiver - modifications should be effective
func (r *Rectangle) SetHeight(height int) {
	r.height = height
}

func main() {
	// Create a rectangle
	rect := Rectangle{width: 10, height: 5}
	
	// Test value receiver method (should not modify the original)
	rect.SetWidth(20)
	
	// Test pointer receiver method (should modify the original)
	rect.SetHeight(20)
	
	// Calculate area
	area := rect.Area()
	
	return area
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

	// Print instructions for debugging
	instructions := vmInstance.GetInstructions()
	t.Logf("Instructions:")
	for i, instr := range instructions {
		t.Logf("%d: %s", i, instr.String())
	}

	// Execute the compiled code
	result, err := vmInstance.Execute(nil)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	// With pointer receiver modifying height to 20 and width unchanged at 10,
	// area should be 10 * 20 = 200
	expected := 200
	if result != expected {
		t.Errorf("Expected result to be %d, got %d", expected, result)

		// Print stack for debugging
		t.Logf("Stack size: %d", vmInstance.StackSize())
	}
}

package test

import (
	"testing"

	"github.com/lengzhao/goscript/builtin"
	"github.com/lengzhao/goscript/compiler"
	"github.com/lengzhao/goscript/parser"
	"github.com/lengzhao/goscript/vm"
)

// TestModuleVsMethodCall tests the distinction between module calls and method calls
func TestModuleVsMethodCall(t *testing.T) {
	// Create a new VM
	vmInstance := vm.NewVM()
	vmInstance.SetDebug(false) // Disable debug output for cleaner test logs

	// Register builtin functions
	for name, fn := range builtin.BuiltInFunctions {
		// Convert builtin.Function to vm.ScriptFunction
		vmInstance.RegisterFunction(name, func(f builtin.Function) vm.ScriptFunction {
			return func(args ...interface{}) (interface{}, error) {
				return f(args...)
			}
		}(fn))
	}

	// Register builtin modules
	modules := []string{"math"}
	for _, moduleName := range modules {
		moduleExecutor, exists := builtin.GetModuleExecutor(moduleName)
		if exists {
			vmInstance.RegisterModule(moduleName, moduleExecutor)
		}
	}

	// Create a new compiler
	compiler := compiler.NewCompiler(vmInstance)

	// Parse the test script
	parser := parser.New()
	script := `
package test

import "math"

type Rectangle struct {
	width  int
	height int
}

func (r Rectangle) Area() int {
	return r.width * r.height
}

func (r Rectangle) MaxDimension() int {
	// This method has the same name as a math module function to test distinction
	if r.width > r.height {
		return r.width
	}
	return r.height
}

func main() {
	// Create a rectangle
	rect := Rectangle{width: 10, height: 5}
	
	// Test method call
	area := rect.Area()
	
	// Test method call that has the same name as a module function
	maxDim := rect.MaxDimension()
	
	// Test module function call
	maxValue := math.Max(10.0, 20.0)  // Use float64 for math.Max
	
	// Return sum without type conversion
	return int(area + maxDim) + int(maxValue)
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

	// Execute the compiled code
	result, err := vmInstance.Execute("")
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	// Expected result:
	// area = 10 * 5 = 50
	// maxDim = max(10, 5) = 10
	// maxValue = math.Max(10.0, 20.0) = 20.0
	// total = 50 + 10 + 20 = 80
	expected := 80
	if result != expected {
		t.Errorf("Expected result to be %d, got %d", expected, result)
	}
}

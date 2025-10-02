package test

import (
	"testing"

	"github.com/lengzhao/goscript/builtin"
	"github.com/lengzhao/goscript/compiler"
	"github.com/lengzhao/goscript/parser"
	"github.com/lengzhao/goscript/vm"
)

// TestUnifiedCallHandling tests the unified handling of module calls and method calls
func TestUnifiedCallHandling(t *testing.T) {
	// Create a new VM
	vmInstance := vm.NewVM()
	vmInstance.SetDebug(false) // Disable debug output for cleaner test logs

	// Register builtin modules
	modules := []string{"math"}
	for _, moduleName := range modules {
		moduleExecutor, exists := builtin.GetModuleExecutor(moduleName)
		if exists {
			vmInstance.RegisterModule(moduleName, moduleExecutor)
		}
	}

	// Register builtin functions
	for name, fn := range builtin.BuiltInFunctions {
		// Convert builtin.Function to vm.ScriptFunction
		vmInstance.RegisterFunction(name, func(f builtin.Function) vm.ScriptFunction {
			return func(args ...interface{}) (interface{}, error) {
				return f(args...)
			}
		}(fn))
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

func (r Rectangle) Scale(factor int) Rectangle {
	return Rectangle{width: r.width * factor, height: r.height * factor}
}

func main() {
	// Create a rectangle
	rect := Rectangle{width: 10, height: 5}
	
	// Test method call
	area := rect.Area()
	
	// Test method call with return value
	scaledRect := rect.Scale(2)
	scaledArea := scaledRect.Area()
	
	// Test module function call
	maxValue := math.Max(10.0, 20.0)  // Use float64 for math.Max
	
	// Return sum of all values
	return area + scaledArea + int(maxValue)
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
	t.Logf("Generated %d instructions:", len(instructions))
	for i, instr := range instructions {
		t.Logf("  %d: %s", i, instr.String())
	}

	// Execute the compiled code
	result, err := vmInstance.Execute("")
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	// Expected result:
	// area = 10 * 5 = 50
	// scaledArea = (10*2) * (5*2) = 20 * 10 = 200
	// maxValue = math.Max(10.0, 20.0) = 20.0
	// total = 50 + 200 + 20 = 270
	expected := 270
	if result != expected {
		t.Errorf("Expected result to be %d, got %d", expected, result)
	}
}

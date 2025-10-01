package vm

import (
	"testing"

	"github.com/lengzhao/goscript/builtin"
)

func TestModuleIntegration(t *testing.T) {
	vm := NewVM()

	// Register builtin modules with simplified interface
	modules := []string{"strings", "fmt", "math", "json"}
	for _, moduleName := range modules {
		if moduleFuncs, exists := builtin.GetModuleFunctions(moduleName); exists {
			// Create a module executor that delegates to the builtin functions
			moduleExecutor := func(entrypoint string, args ...interface{}) (interface{}, error) {
				if fn, ok := moduleFuncs[entrypoint]; ok {
					return fn(args...)
				}
				return nil, nil
			}
			vm.RegisterModule(moduleName, moduleExecutor)
		}
	}

	// Test strings module function
	toUpperFn, exists := vm.GetFunction("strings.ToUpper")
	if !exists {
		t.Errorf("Function 'strings.ToUpper' should exist")
	}

	result, err := toUpperFn("hello")
	if err != nil {
		t.Errorf("Function call should not return error: %v", err)
	}

	if result != "HELLO" {
		t.Errorf("Function should return 'HELLO', got %v", result)
	}

	// Test math module function
	absFn, exists := vm.GetFunction("math.Abs")
	if !exists {
		t.Errorf("Function 'math.Abs' should exist")
	}

	result, err = absFn(-5)
	if err != nil {
		t.Errorf("Function call should not return error: %v", err)
	}

	if result != 5 {
		t.Errorf("Function should return 5, got %v", result)
	}

	// Test that non-existent module function doesn't exist
	_, exists = vm.GetFunction("nonexistent.NonExistent")
	if exists {
		t.Errorf("Function 'nonexistent.NonExistent' should not exist")
	}
}

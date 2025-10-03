package vm

import (
	"testing"
)

func TestModuleRegistration(t *testing.T) {
	vm := NewVM()

	// Create a test module executor
	testModule := func(entrypoint string, args ...interface{}) (interface{}, error) {
		switch entrypoint {
		case "testFunc":
			return "test result", nil
		case "add":
			if len(args) != 2 {
				return nil, nil
			}
			a, ok1 := args[0].(int)
			b, ok2 := args[1].(int)
			if !ok1 || !ok2 {
				return nil, nil
			}
			return a + b, nil
		default:
			return nil, nil
		}
	}

	// Register the module
	vm.RegisterModule("test", testModule)

	// Check if module is registered
	module, exists := vm.GetModule("test")
	if !exists {
		t.Errorf("Module 'test' should exist after registration")
	}

	if module == nil {
		t.Errorf("Module should not be nil")
	}

	// Check if we can get a function from the module
	fn, exists := vm.GetFunction("test.testFunc")
	if !exists {
		t.Errorf("Function 'test.testFunc' should exist")
	}

	// Test the function
	result, err := fn()
	if err != nil {
		t.Errorf("Function call should not return error: %v", err)
	}

	if result != "test result" {
		t.Errorf("Function should return 'test result', got %v", result)
	}

	// Check if we can get another function from the module
	addFn, exists := vm.GetFunction("test.add")
	if !exists {
		t.Errorf("Function 'test.add' should exist")
	}

	// Test the add function
	result, err = addFn(2, 3)
	if err != nil {
		t.Errorf("Function call should not return error: %v", err)
	}

	if result != 5 {
		t.Errorf("Function should return 5, got %v", result)
	}

	// Check that non-existent module doesn't exist
	_, exists = vm.GetModule("nonexistent")
	if exists {
		t.Errorf("Module 'nonexistent' should not exist")
	}

	// Note: Even non-existent entrypoints in existing modules will return a wrapper function
	// because we can't know at registration time what entrypoints are valid
	_, exists = vm.GetFunction("test.nonexistent")
	if !exists {
		t.Errorf("Function 'test.nonexistent' should exist as a wrapper")
	}
}

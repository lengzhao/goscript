package builtin

import (
	"testing"
)

func TestBuiltinModuleExecutor(t *testing.T) {
	// Test strings module
	moduleExecutor, exists := GetModuleExecutor("strings")
	if !exists {
		t.Errorf("strings module should exist")
	}

	// Test ToUpper function
	result, err := moduleExecutor("ToUpper", "hello")
	if err != nil {
		t.Errorf("Failed to call ToUpper: %v", err)
	}

	if result != "HELLO" {
		t.Errorf("Expected 'HELLO', got '%v'", result)
	}

	// Test Contains function
	result, err = moduleExecutor("Contains", "hello world", "world")
	if err != nil {
		t.Errorf("Failed to call Contains: %v", err)
	}

	if result != true {
		t.Errorf("Expected true, got '%v'", result)
	}

	// Test math module
	moduleExecutor, exists = GetModuleExecutor("math")
	if !exists {
		t.Errorf("math module should exist")
	}

	// Test Abs function
	result, err = moduleExecutor("Abs", -5)
	if err != nil {
		t.Errorf("Failed to call Abs: %v", err)
	}

	if result != 5 {
		t.Errorf("Expected 5, got '%v'", result)
	}

	// Test non-existent module
	_, exists = GetModuleExecutor("nonexistent")
	if exists {
		t.Errorf("nonexistent module should not exist")
	}

	// Test non-existent function
	_, err = moduleExecutor("NonExistent", "test")
	if err == nil {
		t.Errorf("NonExistent function should return an error")
	}
}
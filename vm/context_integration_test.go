package vm

import (
	"testing"
)

func TestContextIntegration(t *testing.T) {
	// Create a new VM with context
	vm := NewVM()

	// Test that global context is properly initialized
	globalCtx := vm.GetGlobalContext()
	if globalCtx == nil {
		t.Fatal("Expected non-nil global context")
	}

	if globalCtx.GetPathKey() != "global" {
		t.Errorf("Expected global context path key 'global', got '%s'", globalCtx.GetPathKey())
	}

	// Test entering a new scope
	newCtx := vm.EnterScope("test.function")
	if newCtx == nil {
		t.Fatal("Expected non-nil new context")
	}

	if newCtx.GetPathKey() != "test.function" {
		t.Errorf("Expected new context path key 'test.function', got '%s'", newCtx.GetPathKey())
	}

	// Test that current context is the new context
	currentCtx := vm.GetCurrentContext()
	if currentCtx != newCtx {
		t.Error("Expected current context to be the new context")
	}

	// Test setting a variable in the new context
	err := vm.SetVariable("testVar", 42)
	if err != nil {
		t.Fatalf("Unexpected error setting variable: %v", err)
	}

	// Test getting the variable from the same context
	value, exists := vm.GetVariable("testVar")
	if !exists {
		t.Fatal("Expected variable to exist")
	}

	if value != 42 {
		t.Errorf("Expected value 42, got %v", value)
	}

	// Test exiting the scope
	parentCtx := vm.ExitScope()
	if parentCtx != globalCtx {
		t.Error("Expected parent context to be the global context")
	}

	// Test that current context is back to global context
	currentCtx = vm.GetCurrentContext()
	if currentCtx != globalCtx {
		t.Error("Expected current context to be the global context")
	}
}

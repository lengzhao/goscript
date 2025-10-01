package vm

import (
	"testing"

	"github.com/lengzhao/goscript/builtin"
)

func TestVMWithModuleExecutor(t *testing.T) {
	// Create a new VM
	vm := NewVM()

	// Test registering a module using ModuleExecutor
	moduleExecutor, exists := builtin.GetModuleExecutor("strings")
	if !exists {
		t.Errorf("strings module should exist")
	}

	// Register the module
	vm.RegisterModule("strings", moduleExecutor)

	// Verify the module was registered
	registeredModule, exists := vm.GetModule("strings")
	if !exists {
		t.Errorf("strings module should be registered")
	}

	if registeredModule == nil {
		t.Errorf("registered module should not be nil")
	}

	// Test calling a function through GetFunction
	toUpperFn, exists := vm.GetFunction("strings.ToUpper")
	if !exists {
		t.Errorf("strings.ToUpper function should exist")
	}

	result, err := toUpperFn("hello")
	if err != nil {
		t.Errorf("Failed to call strings.ToUpper: %v", err)
	}

	if result != "HELLO" {
		t.Errorf("Expected 'HELLO', got '%v'", result)
	}

	// Test another function
	containsFn, exists := vm.GetFunction("strings.Contains")
	if !exists {
		t.Errorf("strings.Contains function should exist")
	}

	result, err = containsFn("hello world", "world")
	if err != nil {
		t.Errorf("Failed to call strings.Contains: %v", err)
	}

	if result != true {
		t.Errorf("Expected true, got '%v'", result)
	}
}
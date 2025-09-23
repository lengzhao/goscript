package vm

import (
	"testing"
)

func TestScopeManagement(t *testing.T) {
	vm := NewVM()

	// Set a global variable
	err := vm.SetVariable("globalVar", "globalValue")
	if err != nil {
		t.Fatalf("Unexpected error setting global variable: %v", err)
	}

	// Enter a function scope
	functionCtx := vm.EnterScope("main.myFunction")
	if functionCtx == nil {
		t.Fatal("Expected non-nil function context")
	}

	// Set a local variable in the function scope
	err = vm.SetVariable("localVar", "localValue")
	if err != nil {
		t.Fatalf("Unexpected error setting local variable: %v", err)
	}

	// Set another variable with the same name as the global one (should shadow it)
	err = vm.SetVariable("globalVar", "localGlobalValue")
	if err != nil {
		t.Fatalf("Unexpected error setting local global variable: %v", err)
	}

	// Test variable lookup in the function scope
	// Should find the local variable
	localVar, exists := vm.GetVariable("localVar")
	if !exists {
		t.Fatal("Expected localVar to exist")
	}
	if localVar != "localValue" {
		t.Errorf("Expected localVar to be 'localValue', got '%v'", localVar)
	}

	// Should find the local global variable (shadowing the global one)
	globalVar, exists := vm.GetVariable("globalVar")
	if !exists {
		t.Fatal("Expected globalVar to exist")
	}
	if globalVar != "localGlobalValue" {
		t.Errorf("Expected globalVar to be 'localGlobalValue', got '%v'", globalVar)
	}

	// Exit the function scope
	vm.ExitScope()

	// Now we should be back in the global scope
	// Should find the original global variable
	globalVar, exists = vm.GetVariable("globalVar")
	if !exists {
		t.Fatal("Expected globalVar to exist")
	}
	if globalVar != "globalValue" {
		t.Errorf("Expected globalVar to be 'globalValue', got '%v'", globalVar)
	}

	// Should not find the local variable anymore
	_, exists = vm.GetVariable("localVar")
	if exists {
		t.Error("Expected localVar to not exist in global scope")
	}
}

func TestNestedScopeManagement(t *testing.T) {
	vm := NewVM()

	// Set a global variable
	err := vm.SetVariable("x", 10)
	if err != nil {
		t.Fatalf("Unexpected error setting global variable: %v", err)
	}

	// Enter function scope
	vm.EnterScope("main.calculate")

	// Set a local variable
	err = vm.SetVariable("x", 20)
	if err != nil {
		t.Fatalf("Unexpected error setting local variable: %v", err)
	}

	// Enter nested block scope
	vm.EnterScope("main.calculate.block")

	// Set another local variable
	err = vm.SetVariable("x", 30)
	if err != nil {
		t.Fatalf("Unexpected error setting nested variable: %v", err)
	}

	// Test variable lookup - should find the most local one
	x, exists := vm.GetVariable("x")
	if !exists {
		t.Fatal("Expected x to exist")
	}
	if x != 30 {
		t.Errorf("Expected x to be 30, got %v", x)
	}

	// Exit nested block scope
	vm.ExitScope()

	// Now we should find the function scope variable
	x, exists = vm.GetVariable("x")
	if !exists {
		t.Fatal("Expected x to exist")
	}
	if x != 20 {
		t.Errorf("Expected x to be 20, got %v", x)
	}

	// Exit function scope
	vm.ExitScope()

	// Now we should find the global variable
	x, exists = vm.GetVariable("x")
	if !exists {
		t.Fatal("Expected x to exist")
	}
	if x != 10 {
		t.Errorf("Expected x to be 10, got %v", x)
	}
}
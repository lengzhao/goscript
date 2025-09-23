package vm

import (
	"testing"
)

func TestNewVMWithContex(t *testing.T) {
	vm := NewVMWithContex()

	if vm == nil {
		t.Fatal("Expected non-nil VM")
	}

	if vm.globalCtx == nil {
		t.Error("Expected non-nil global context")
	}

	if vm.currentCtx == nil {
		t.Error("Expected non-nil current context")
	}

	if vm.currentCtx != vm.globalCtx {
		t.Error("Expected current context to be global context initially")
	}

	if vm.stack == nil {
		t.Error("Expected non-nil stack")
	}

	if vm.instructions == nil {
		t.Error("Expected non-nil instructions")
	}

	if vm.functionRegistry == nil {
		t.Error("Expected non-nil function registry")
	}

	if vm.scriptFunctions == nil {
		t.Error("Expected non-nil script functions")
	}

	if vm.typeSystem == nil {
		t.Error("Expected non-nil type system")
	}

	if vm.handlers == nil {
		t.Error("Expected non-nil handlers")
	}
}

func TestVMContextOperations(t *testing.T) {
	vm := NewVMWithContex()

	// Test EnterScope and ExitScope
	initialCtx := vm.GetCurrentContext()
	if initialCtx.GetPathKey() != "global" {
		t.Errorf("Expected initial context path key 'global', got '%s'", initialCtx.GetPathKey())
	}

	// Enter a new scope
	newCtx := vm.EnterScope("test.scope")
	if newCtx == nil {
		t.Fatal("Expected non-nil new context")
	}

	if newCtx.GetPathKey() != "test.scope" {
		t.Errorf("Expected new context path key 'test.scope', got '%s'", newCtx.GetPathKey())
	}

	if vm.GetCurrentContext() != newCtx {
		t.Error("Expected current context to be the new context")
	}

	// Exit the scope
	parentCtx := vm.ExitScope()
	if parentCtx != initialCtx {
		t.Error("Expected parent context to be the initial context")
	}

	if vm.GetCurrentContext() != initialCtx {
		t.Error("Expected current context to be back to initial context")
	}
}

func TestVMVariableOperations(t *testing.T) {
	vm := NewVMWithContex()

	// Test SetVariable and GetVariable
	err := vm.SetVariable("testVar", 42)
	if err != nil {
		t.Fatalf("Unexpected error setting variable: %v", err)
	}

	value, exists := vm.GetVariable("testVar")
	if !exists {
		t.Fatal("Expected variable to exist")
	}

	if value != 42 {
		t.Errorf("Expected value 42, got %v", value)
	}

	// Test HasVariable
	if !vm.HasVariable("testVar") {
		t.Error("Expected HasVariable to return true")
	}

	// Test DeleteVariable
	vm.DeleteVariable("testVar")
	_, exists = vm.GetVariable("testVar")
	if exists {
		t.Error("Expected variable to be deleted")
	}
}

func TestVMContextHierarchy(t *testing.T) {
	vm := NewVMWithContex()

	// Set a variable in global context
	vm.SetVariable("globalVar", "globalValue")

	// Enter a new scope
	vm.EnterScope("function.main")

	// Set a variable in the function context
	vm.SetVariable("localVar", "localValue")

	// Test variable lookup in hierarchy
	// Local variable should be found in current context
	value, exists := vm.GetVariable("localVar")
	if !exists {
		t.Fatal("Expected local variable to exist")
	}
	if value != "localValue" {
		t.Errorf("Expected 'localValue', got '%v'", value)
	}

	// Global variable should be found in parent context
	value, exists = vm.GetVariable("globalVar")
	if !exists {
		t.Fatal("Expected global variable to be found in hierarchy")
	}
	if value != "globalValue" {
		t.Errorf("Expected 'globalValue', got '%v'", value)
	}
}

func TestVMRun(t *testing.T) {
	vm := NewVMWithContex()

	// Add some simple instructions
	vm.AddInstruction(NewInstruction(OpLoadConst, 42, nil))
	vm.AddInstruction(NewInstruction(OpLoadConst, 8, nil))
	vm.AddInstruction(NewInstruction(OpBinaryOp, OpAdd, nil))
	vm.AddInstruction(NewInstruction(OpReturn, nil, nil))

	// Test Run in global context
	result, err := vm.Run(vm.globalCtx, 0, len(vm.instructions))
	if err != nil {
		t.Fatalf("Unexpected error running VM: %v", err)
	}

	// The result should be 50 (42 + 8)
	if result != 50 {
		t.Errorf("Expected result 50, got %v", result)
	}
}

func TestVMExecute(t *testing.T) {
	vm := NewVMWithContex()

	// Add some simple instructions
	vm.AddInstruction(NewInstruction(OpLoadConst, 10, nil))
	vm.AddInstruction(NewInstruction(OpLoadConst, 5, nil))
	vm.AddInstruction(NewInstruction(OpBinaryOp, OpMul, nil))
	vm.AddInstruction(NewInstruction(OpReturn, nil, nil))

	// Test Execute
	result, err := vm.Execute(nil)
	if err != nil {
		t.Fatalf("Unexpected error executing VM: %v", err)
	}

	// The result should be 50 (10 * 5)
	if result != 50 {
		t.Errorf("Expected result 50, got %v", result)
	}
}

func TestVMGetGlobalContext(t *testing.T) {
	vm := NewVMWithContex()

	globalCtx := vm.GetGlobalContext()
	if globalCtx == nil {
		t.Fatal("Expected non-nil global context")
	}

	if globalCtx.GetPathKey() != "global" {
		t.Errorf("Expected global context path key 'global', got '%s'", globalCtx.GetPathKey())
	}
}

func TestVMGetCurrentContext(t *testing.T) {
	vm := NewVMWithContex()

	currentCtx := vm.GetCurrentContext()
	if currentCtx == nil {
		t.Fatal("Expected non-nil current context")
	}

	if currentCtx.GetPathKey() != "global" {
		t.Errorf("Expected current context path key 'global', got '%s'", currentCtx.GetPathKey())
	}

	// Enter a new scope
	vm.EnterScope("test.scope")
	currentCtx = vm.GetCurrentContext()
	if currentCtx.GetPathKey() != "test.scope" {
		t.Errorf("Expected current context path key 'test.scope', got '%s'", currentCtx.GetPathKey())
	}
}

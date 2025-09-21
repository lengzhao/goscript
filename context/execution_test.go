package context

import (
	"testing"
	"time"
)

func TestNewExecutionContext(t *testing.T) {
	ec := NewExecutionContext()

	if ec == nil {
		t.Fatal("Expected non-nil ExecutionContext")
	}

	if ec.Context == nil {
		t.Error("Expected non-nil context")
	}

	if ec.Cancel == nil {
		t.Error("Expected non-nil cancel function")
	}

	if ec.ScopeManager == nil {
		t.Error("Expected non-nil scope manager")
	}

	if ec.ModuleName != "main" {
		t.Errorf("Expected module name 'main', got '%s'", ec.ModuleName)
	}

	if ec.Parent != nil {
		t.Error("Expected nil parent")
	}

	if ec.Security == nil {
		t.Error("Expected non-nil security context")
	}
}

func TestNewExecutionContextWithTimeout(t *testing.T) {
	timeout := 100 * time.Millisecond
	ec := NewExecutionContextWithTimeout(timeout)

	if ec == nil {
		t.Fatal("Expected non-nil ExecutionContext")
	}

	if ec.Security.MaxExecutionTime != timeout {
		t.Errorf("Expected MaxExecutionTime %v, got %v", timeout, ec.Security.MaxExecutionTime)
	}
}

func TestNewChildContext(t *testing.T) {
	parent := NewExecutionContext()
	child := parent.NewChildContext("child")

	if child == nil {
		t.Fatal("Expected non-nil child context")
	}

	if child.Parent != parent {
		t.Error("Expected child's parent to be the parent context")
	}

	if child.ModuleName != "child" {
		t.Errorf("Expected module name 'child', got '%s'", child.ModuleName)
	}

	if child.Security != parent.Security {
		t.Error("Expected child and parent to share the same security context")
	}

	if child.ScopeManager != parent.ScopeManager {
		t.Error("Expected child and parent to share the same scope manager")
	}
}

func TestExitContext(t *testing.T) {
	parent := NewExecutionContext()
	child := parent.NewChildContext("child")

	// Exit child context
	exited := child.ExitContext()

	if exited != parent {
		t.Error("Expected to exit to parent context")
	}
}

func TestSetValueGetValue(t *testing.T) {
	ec := NewExecutionContext()

	// Set a value
	err := ec.SetValue("testVar", 42)
	if err != nil {
		t.Fatalf("Unexpected error setting value: %v", err)
	}

	// Get the value
	value, exists := ec.GetValue("testVar")
	if !exists {
		t.Fatal("Expected variable to exist")
	}

	if value != 42 {
		t.Errorf("Expected value 42, got %v", value)
	}

	// Try to get a non-existent value
	_, exists = ec.GetValue("nonExistent")
	if exists {
		t.Error("Expected non-existent variable to not exist")
	}
}

func TestUpdateValue(t *testing.T) {
	ec := NewExecutionContext()

	// Set initial value
	ec.SetValue("testVar", 42)

	// Update the value
	err := ec.UpdateValue("testVar", 100)
	if err != nil {
		t.Fatalf("Unexpected error updating value: %v", err)
	}

	// Get the updated value
	value, _ := ec.GetValue("testVar")
	if value != 100 {
		t.Errorf("Expected updated value 100, got %v", value)
	}

	// Try to update a non-existent value
	err = ec.UpdateValue("nonExistent", 50)
	if err == nil {
		t.Error("Expected error when updating non-existent variable")
	}
}

func TestDeleteValue(t *testing.T) {
	ec := NewExecutionContext()

	// Set a value
	ec.SetValue("testVar", 42)

	// Delete the value
	ec.DeleteValue("testVar")

	// Try to get the deleted value
	_, exists := ec.GetValue("testVar")
	if exists {
		t.Error("Expected deleted variable to not exist")
	}
}

func TestSetGlobalValue(t *testing.T) {
	parent := NewExecutionContext()
	child := parent.NewChildContext("child")

	// Set a global value from child context
	err := child.SetGlobalValue("globalVar", "globalValue")
	if err != nil {
		t.Fatalf("Unexpected error setting global value: %v", err)
	}

	// Get the value from parent context
	value, exists := parent.GetValue("globalVar")
	if !exists {
		t.Fatal("Expected global variable to exist in parent context")
	}

	if value != "globalValue" {
		t.Errorf("Expected global value 'globalValue', got %v", value)
	}
}

func TestGetRootContext(t *testing.T) {
	parent := NewExecutionContext()
	child := parent.NewChildContext("child")
	grandchild := child.NewChildContext("grandchild")

	root := grandchild.GetRootContext()
	if root != parent {
		t.Error("Expected root context to be the top-level context")
	}
}

func TestCheckModuleAccess(t *testing.T) {
	ec := NewExecutionContext()

	// Test allowed module
	err := ec.CheckModuleAccess("fmt")
	if err != nil {
		t.Errorf("Expected fmt module to be allowed, got error: %v", err)
	}

	// Test forbidden keyword
	// Note: This test may need to be adjusted based on actual implementation
	// as the current implementation doesn't seem to check forbidden keywords
	// in the CheckModuleAccess method
}

func TestIsCancelled(t *testing.T) {
	ec := NewExecutionContext()

	// Initially should not be cancelled
	if ec.IsCancelled() {
		t.Error("Expected context to not be cancelled initially")
	}

	// Cancel the context
	ec.Cancel()

	// Should now be cancelled
	if !ec.IsCancelled() {
		t.Error("Expected context to be cancelled after calling Cancel()")
	}
}

func TestDeadline(t *testing.T) {
	// Test with regular context
	ec := NewExecutionContext()
	_, hasDeadline := ec.Deadline()
	if hasDeadline {
		t.Error("Expected no deadline for regular context")
	}

	// Test with timeout context
	timeout := 100 * time.Millisecond
	ecWithTimeout := NewExecutionContextWithTimeout(timeout)
	deadline, hasDeadline := ecWithTimeout.Deadline()
	if !hasDeadline {
		t.Error("Expected deadline for context with timeout")
	}

	if deadline.IsZero() {
		t.Error("Expected non-zero deadline")
	}
}

func TestRegisterAndGetFunction(t *testing.T) {
	ec := NewExecutionContext()

	// Create a test function
	testFunc := func(args ...interface{}) (interface{}, error) {
		return "test", nil
	}

	// Register the function
	err := ec.RegisterFunction("testFunc", testFunc)
	if err != nil {
		t.Fatalf("Unexpected error registering function: %v", err)
	}

	// Get the function
	fn, exists := ec.GetFunction("testFunc")
	if !exists {
		t.Fatal("Expected function to exist")
	}

	if fn == nil {
		t.Error("Expected non-nil function")
	}

	// Try to get a non-existent function
	_, exists = ec.GetFunction("nonExistent")
	if exists {
		t.Error("Expected non-existent function to not exist")
	}
}

func TestGetAllFunctions(t *testing.T) {
	ec := NewExecutionContext()

	// Register a function
	testFunc := func(args ...interface{}) (interface{}, error) {
		return "test", nil
	}
	ec.RegisterFunction("testFunc", testFunc)

	// Get all functions
	functions := ec.GetAllFunctions()
	if len(functions) != 1 {
		t.Errorf("Expected 1 function, got %d", len(functions))
	}

	if functions["testFunc"] == nil {
		t.Error("Expected testFunc to be in the functions map")
	}
}

func TestWithContextTimeout(t *testing.T) {
	timeout := 100 * time.Millisecond
	ec := NewExecutionContext()
	ecWithTimeout := ec.WithTimeout(timeout)

	if ecWithTimeout == nil {
		t.Fatal("Expected non-nil ExecutionContext with timeout")
	}

	if ecWithTimeout.Security.MaxExecutionTime != timeout {
		t.Errorf("Expected MaxExecutionTime %v, got %v", timeout, ecWithTimeout.Security.MaxExecutionTime)
	}
}

func TestString(t *testing.T) {
	ec := NewExecutionContext()

	// Test string representation without deadline
	str := ec.String()
	if str == "" {
		t.Error("Expected non-empty string representation")
	}

	// Test string representation with deadline
	timeout := 100 * time.Millisecond
	ecWithTimeout := NewExecutionContextWithTimeout(timeout)
	strWithTimeout := ecWithTimeout.String()
	if strWithTimeout == "" {
		t.Error("Expected non-empty string representation with timeout")
	}
}

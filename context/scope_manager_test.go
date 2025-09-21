package context

import (
	"testing"
	"time"
)

func TestNewScopeManager(t *testing.T) {
	sm := NewScopeManager()

	if sm == nil {
		t.Fatal("Expected non-nil ScopeManager")
	}

	if sm.root == nil {
		t.Error("Expected non-nil root scope")
	}

	if sm.current != sm.root {
		t.Error("Expected current scope to be root scope")
	}

	if sm.functions == nil {
		t.Error("Expected non-nil functions map")
	}
}

func TestNewScopeManagerWithTimeout(t *testing.T) {
	timeout := 100 * time.Millisecond
	sm := NewScopeManagerWithTimeout(timeout)

	if sm == nil {
		t.Fatal("Expected non-nil ScopeManager")
	}

	if sm.deadline.IsZero() {
		t.Error("Expected non-zero deadline")
	}
}

func TestNewScopeAndExitScope(t *testing.T) {
	sm := NewScopeManager()

	// Initially should be at root level
	if sm.GetScopeLevel() != 1 {
		t.Errorf("Expected initial scope level 1, got %d", sm.GetScopeLevel())
	}

	// Create a new scope
	sm.NewScope()

	if sm.GetScopeLevel() != 2 {
		t.Errorf("Expected scope level 2 after NewScope, got %d", sm.GetScopeLevel())
	}

	// Exit the scope
	sm.ExitScope()

	if sm.GetScopeLevel() != 1 {
		t.Errorf("Expected scope level 1 after ExitScope, got %d", sm.GetScopeLevel())
	}
}

func TestSetAndGetVariable(t *testing.T) {
	sm := NewScopeManager()

	// Set a variable
	sm.SetVariable("testVar", 42, "int")

	// Get the variable
	variable, exists := sm.GetVariable("testVar")
	if !exists {
		t.Fatal("Expected variable to exist")
	}

	if variable == nil {
		t.Fatal("Expected non-nil variable")
	}

	if variable.Name != "testVar" {
		t.Errorf("Expected variable name 'testVar', got '%s'", variable.Name)
	}

	if variable.Value != 42 {
		t.Errorf("Expected variable value 42, got %v", variable.Value)
	}

	if variable.Type != "int" {
		t.Errorf("Expected variable type 'int', got '%s'", variable.Type)
	}

	// Try to get a non-existent variable
	_, exists = sm.GetVariable("nonExistent")
	if exists {
		t.Error("Expected non-existent variable to not exist")
	}
}

func TestUpdateVariable(t *testing.T) {
	sm := NewScopeManager()

	// Set initial value
	sm.SetVariable("testVar", 42, "int")

	// Update the value
	err := sm.UpdateVariable("testVar", 100)
	if err != nil {
		t.Fatalf("Unexpected error updating variable: %v", err)
	}

	// Get the updated value
	variable, _ := sm.GetVariable("testVar")
	if variable.Value != 100 {
		t.Errorf("Expected updated value 100, got %v", variable.Value)
	}

	// Try to update a non-existent variable
	err = sm.UpdateVariable("nonExistent", 50)
	if err == nil {
		t.Error("Expected error when updating non-existent variable")
	}
}

func TestDeleteVariable(t *testing.T) {
	sm := NewScopeManager()

	// Set a variable
	sm.SetVariable("testVar", 42, "int")

	// Delete the variable
	sm.DeleteVariable("testVar")

	// Try to get the deleted variable
	_, exists := sm.GetVariable("testVar")
	if exists {
		t.Error("Expected deleted variable to not exist")
	}
}

func TestGetCurrentScopeVariables(t *testing.T) {
	sm := NewScopeManager()

	// Set variables
	sm.SetVariable("var1", 1, "int")
	sm.SetVariable("var2", "test", "string")

	// Get current scope variables
	variables := sm.GetCurrentScopeVariables()

	if len(variables) != 2 {
		t.Errorf("Expected 2 variables, got %d", len(variables))
	}

	if variables["var1"] == nil {
		t.Error("Expected var1 to be in variables map")
	}

	if variables["var2"] == nil {
		t.Error("Expected var2 to be in variables map")
	}
}

func TestGetAllVariables(t *testing.T) {
	sm := NewScopeManager()

	// Set variables in root scope
	sm.SetVariable("globalVar", "global", "string")

	// Create new scope and set variables
	sm.NewScope()
	sm.SetVariable("localVar", 42, "int")

	// Get all variables (should include both global and local)
	allVariables := sm.GetAllVariables()

	if len(allVariables) != 2 {
		t.Errorf("Expected 2 variables, got %d", len(allVariables))
	}

	if allVariables["globalVar"] == nil {
		t.Error("Expected globalVar to be in all variables map")
	}

	if allVariables["localVar"] == nil {
		t.Error("Expected localVar to be in all variables map")
	}
}

func TestScopeManagerRegisterAndGetFunction(t *testing.T) {
	sm := NewScopeManager()

	// Create a test function
	testFunc := func(args ...interface{}) (interface{}, error) {
		return "test", nil
	}

	// Register the function
	sm.RegisterFunction("testFunc", testFunc)

	// Get the function
	fn, exists := sm.GetFunction("testFunc")
	if !exists {
		t.Fatal("Expected function to exist")
	}

	if fn == nil {
		t.Error("Expected non-nil function")
	}

	// Try to get a non-existent function
	_, exists = sm.GetFunction("nonExistent")
	if exists {
		t.Error("Expected non-existent function to not exist")
	}
}

func TestScopeManagerGetAllFunctions(t *testing.T) {
	sm := NewScopeManager()

	// Register a function
	testFunc := func(args ...interface{}) (interface{}, error) {
		return "test", nil
	}
	sm.RegisterFunction("testFunc", testFunc)

	// Get all functions
	functions := sm.GetAllFunctions()
	if len(functions) != 1 {
		t.Errorf("Expected 1 function, got %d", len(functions))
	}

	if functions["testFunc"] == nil {
		t.Error("Expected testFunc to be in the functions map")
	}
}

func TestScopeManagerIsCancelled(t *testing.T) {
	// Test with regular scope manager
	sm := NewScopeManager()
	if sm.IsCancelled() {
		t.Error("Expected scope manager to not be cancelled initially")
	}

	// Test with timeout scope manager
	timeout := 100 * time.Millisecond
	smWithTimeout := NewScopeManagerWithTimeout(timeout)
	// Initially should not be cancelled
	if smWithTimeout.IsCancelled() {
		t.Error("Expected scope manager with timeout to not be cancelled initially")
	}
}

func TestSetDeadlineAndGetDeadline(t *testing.T) {
	sm := NewScopeManager()

	// Initially should have no deadline
	_, hasDeadline := sm.GetDeadline()
	if hasDeadline {
		t.Error("Expected no deadline initially")
	}

	// Set a deadline
	deadline := time.Now().Add(100 * time.Millisecond)
	sm.SetDeadline(deadline)

	// Get the deadline
	gotDeadline, hasDeadline := sm.GetDeadline()
	if !hasDeadline {
		t.Error("Expected to have deadline after setting it")
	}

	if !gotDeadline.Equal(deadline) {
		t.Errorf("Expected deadline %v, got %v", deadline, gotDeadline)
	}
}

func TestScopeLevel(t *testing.T) {
	sm := NewScopeManager()

	// Initially should be at level 1 (root)
	if sm.GetScopeLevel() != 1 {
		t.Errorf("Expected initial scope level 1, got %d", sm.GetScopeLevel())
	}

	// Create nested scopes
	sm.NewScope()
	if sm.GetScopeLevel() != 2 {
		t.Errorf("Expected scope level 2, got %d", sm.GetScopeLevel())
	}

	sm.NewScope()
	if sm.GetScopeLevel() != 3 {
		t.Errorf("Expected scope level 3, got %d", sm.GetScopeLevel())
	}

	// Exit scopes
	sm.ExitScope()
	if sm.GetScopeLevel() != 2 {
		t.Errorf("Expected scope level 2 after exit, got %d", sm.GetScopeLevel())
	}

	sm.ExitScope()
	if sm.GetScopeLevel() != 1 {
		t.Errorf("Expected scope level 1 after exit, got %d", sm.GetScopeLevel())
	}
}

func TestStringRepresentation(t *testing.T) {
	sm := NewScopeManager()

	// Test string representation
	str := sm.String()
	if str == "" {
		t.Error("Expected non-empty string representation")
	}

	// Test debug string representation
	debugStr := sm.DebugString()
	if debugStr == "" {
		t.Error("Expected non-empty debug string representation")
	}
}

func TestVariableScopeIsolation(t *testing.T) {
	sm := NewScopeManager()

	// Set a variable in the root scope
	sm.SetVariable("globalVar", "globalValue", "string")

	// Create a new scope
	sm.NewScope()

	// Set a variable with the same name in the new scope
	sm.SetVariable("globalVar", "localValue", "string")

	// Get the variable - should return the local value
	variable, _ := sm.GetVariable("globalVar")
	if variable.Value != "localValue" {
		t.Errorf("Expected local value 'localValue', got %v", variable.Value)
	}

	// Exit the scope
	sm.ExitScope()

	// Get the variable - should now return the global value
	variable, _ = sm.GetVariable("globalVar")
	if variable.Value != "globalValue" {
		t.Errorf("Expected global value 'globalValue', got %v", variable.Value)
	}
}

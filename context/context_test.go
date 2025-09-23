package context

import (
	"testing"
)

func TestNewContext(t *testing.T) {
	parent := NewContext("parent", nil)
	ctx := NewContext("child", parent)

	if ctx == nil {
		t.Fatal("Expected non-nil context")
	}

	if ctx.GetPathKey() != "child" {
		t.Errorf("Expected path key 'child', got '%s'", ctx.GetPathKey())
	}

	if ctx.GetParent() != parent {
		t.Error("Expected parent to be set correctly")
	}

	if len(ctx.GetAllVariables()) != 0 {
		t.Error("Expected empty variables map")
	}
}

func TestContextHierarchy(t *testing.T) {
	// Create parent context
	parent := NewContext("parent", nil)
	parent.CreateVariableWithType("parentVar", "parentValue", "string")

	// Create child context
	child := NewContext("child", parent)
	child.CreateVariableWithType("childVar", "childValue", "string")

	// Test variable lookup in hierarchy
	// Child variable should be found in child context
	value, exists := child.GetVariable("childVar")
	if !exists {
		t.Fatal("Expected child variable to exist")
	}
	if value != "childValue" {
		t.Errorf("Expected 'childValue', got '%v'", value)
	}

	// Parent variable should be found in parent context
	value, exists = child.GetVariable("parentVar")
	if !exists {
		t.Fatal("Expected parent variable to be found in hierarchy")
	}
	if value != "parentValue" {
		t.Errorf("Expected 'parentValue', got '%v'", value)
	}

	// Non-existent variable should not be found
	_, exists = child.GetVariable("nonExistent")
	if exists {
		t.Error("Expected non-existent variable to not be found")
	}
}

func TestContextChildren(t *testing.T) {
	parent := NewContext("parent", nil)
	child1 := NewContext("child1", parent)
	child2 := NewContext("child2", parent)

	// Add children to parent
	parent.AddChild(child1)
	parent.AddChild(child2)

	// Test GetChild
	retrievedChild1, exists := parent.GetChild("child1")
	if !exists {
		t.Fatal("Expected child1 to exist")
	}
	if retrievedChild1 != child1 {
		t.Error("Expected retrieved child to be the same as original")
	}

	// Test GetChildren
	children := parent.GetChildren()
	if len(children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(children))
	}

	// Test RemoveChild
	parent.RemoveChild("child1")
	_, exists = parent.GetChild("child1")
	if exists {
		t.Error("Expected child1 to be removed")
	}
}

func TestContextVariablesWithTypes(t *testing.T) {
	ctx := NewContext("test", nil)

	// Test CreateVariableWithType
	ctx.CreateVariableWithType("typedVar", 42, "int")

	// Test GetVariableType
	varType, exists := ctx.GetVariableType("typedVar")
	if !exists {
		t.Fatal("Expected variable type to exist")
	}
	if varType != "int" {
		t.Errorf("Expected type 'int', got '%s'", varType)
	}

	// Test GetAllVariablesWithTypes
	vars, types := ctx.GetAllVariablesWithTypes()
	if len(vars) != 1 {
		t.Errorf("Expected 1 variable, got %d", len(vars))
	}
	if len(types) != 1 {
		t.Errorf("Expected 1 type, got %d", len(types))
	}

	if vars["typedVar"] != 42 {
		t.Errorf("Expected value 42, got %v", vars["typedVar"])
	}

	if types["typedVar"] != "int" {
		t.Errorf("Expected type 'int', got '%s'", types["typedVar"])
	}
}

func TestMustGetVariable(t *testing.T) {
	ctx := NewContext("test", nil)
	ctx.CreateVariableWithType("existingVar", "value", "string")

	// Test that existing variable can be retrieved
	value := ctx.MustGetVariable("existingVar")
	if value != "value" {
		t.Errorf("Expected 'value', got '%v'", value)
	}

	// Test that non-existent variable panics
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected MustGetVariable to panic for non-existent variable")
		}
	}()
	ctx.MustGetVariable("nonExistent")
}

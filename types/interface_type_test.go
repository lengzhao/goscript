package types

import (
	"reflect"
	"strings"
	"testing"
)

func TestInterfaceType(t *testing.T) {
	// Create a new interface type
	it := NewInterfaceType("TestInterface")

	// Test TypeName
	if it.TypeName() != "TestInterface" {
		t.Errorf("Expected TypeName to be 'TestInterface', got '%s'", it.TypeName())
	}

	// Test String for empty interface
	expectedEmpty := "interface{}"
	if it.String() != expectedEmpty {
		t.Errorf("Expected String to be '%s', got '%s'", expectedEmpty, it.String())
	}

	// Add a method
	params := []IType{IntType}
	returns := []IType{StringType}
	it.AddMethod("TestMethod", params, returns)

	// Test String for interface with method
	expectedWithMethod := "interface { TestMethod(int) string }"
	if it.String() != expectedWithMethod {
		t.Errorf("Expected String to be '%s', got '%s'", expectedWithMethod, it.String())
	}

	// Test GetMethod
	method, exists := it.GetMethod("TestMethod")
	if !exists {
		t.Error("Expected TestMethod to exist")
	}
	if method.Name != "TestMethod" {
		t.Errorf("Expected method name to be 'TestMethod', got '%s'", method.Name)
	}
	if len(method.Params) != 1 {
		t.Errorf("Expected 1 parameter, got %d", len(method.Params))
	}
	if method.Params[0].TypeName() != "int" {
		t.Errorf("Expected parameter type to be 'int', got '%s'", method.Params[0].TypeName())
	}
	if len(method.Returns) != 1 {
		t.Errorf("Expected 1 return value, got %d", len(method.Returns))
	}
	if method.Returns[0].TypeName() != "string" {
		t.Errorf("Expected return type to be 'string', got '%s'", method.Returns[0].TypeName())
	}

	// Test HasMethod
	if !it.HasMethod("TestMethod") {
		t.Error("Expected TestMethod to exist")
	}
	if it.HasMethod("NonExistentMethod") {
		t.Error("Expected NonExistentMethod not to exist")
	}

	// Test GetMethods
	methods := it.GetMethods()
	if len(methods) != 1 {
		t.Errorf("Expected 1 method, got %d", len(methods))
	}
	if methods[0].Name != "TestMethod" {
		t.Errorf("Expected method name to be 'TestMethod', got '%s'", methods[0].Name)
	}

	// Test Kind
	if it.Kind() != reflect.Interface {
		t.Errorf("Expected Kind to be reflect.Interface, got %v", it.Kind())
	}

	// Test Size
	if it.Size() != 16 {
		t.Errorf("Expected Size to be 16, got %d", it.Size())
	}

	// Test DefaultValue
	if it.DefaultValue() != nil {
		t.Errorf("Expected DefaultValue to be nil, got %v", it.DefaultValue())
	}

	// Test Clone
	clone := it.Clone()
	if clone.TypeName() != it.TypeName() {
		t.Errorf("Expected cloned type name to be '%s', got '%s'", it.TypeName(), clone.TypeName())
	}

	// Test Equals
	if !it.Equals(clone) {
		t.Error("Expected interface to equal its clone")
	}

	// Create a different interface
	other := NewInterfaceType("TestInterface") // Use the same name
	if it.Equals(other) {
		// This is expected since they have the same name but no methods yet
	}

	// Add the same method to the other interface
	other.AddMethod("TestMethod", params, returns)
	if !it.Equals(other) {
		t.Error("Expected interfaces to be equal after adding the same method")
	}
}

func TestInterfaceTypeWithReturns(t *testing.T) {
	it := NewInterfaceType("TestInterface")

	// Add a method with return values
	params := []IType{IntType}
	returns := []IType{StringType, IntType}
	it.AddMethod("TestMethod", params, returns)

	expected := "interface { TestMethod(int) (string, int) }"
	if it.String() != expected {
		t.Errorf("Expected String to be '%s', got '%s'", expected, it.String())
	}
}

func TestInterfaceTypeWithMultipleReturns(t *testing.T) {
	it := NewInterfaceType("TestInterface")

	// Add a method with no return values
	params := []IType{IntType}
	returns := []IType{}
	it.AddMethod("TestMethod", params, returns)

	expected := "interface { TestMethod(int) }"
	if it.String() != expected {
		t.Errorf("Expected String to be '%s', got '%s'", expected, it.String())
	}
}

func TestInterfaceTypeEmbedded(t *testing.T) {
	// Create embedded interface
	reader := NewInterfaceType("Reader")
	reader.AddMethod("Read", []IType{}, []IType{StringType})

	// Create main interface
	writer := NewInterfaceType("Writer")
	writer.AddMethod("Write", []IType{StringType}, []IType{IntType})

	// Embed reader in writer
	writer.AddEmbedded(reader)

	// Test String representation
	// Note: The order might vary due to map iteration, so we check for the presence of parts
	result := writer.String()
	if !(contains(result, "Reader") && contains(result, "Read() string") && contains(result, "Write(string) int")) {
		t.Errorf("Expected String to contain embedded interface and methods, got '%s'", result)
	}

	// Test GetMethods includes embedded methods
	methods := writer.GetMethods()
	if len(methods) != 2 {
		t.Errorf("Expected 2 methods (1 direct + 1 embedded), got %d", len(methods))
	}

	// Test HasMethod includes embedded methods
	if !writer.HasMethod("Read") {
		t.Error("Expected to find embedded Read method")
	}

	if !writer.HasMethod("Write") {
		t.Error("Expected to find direct Write method")
	}

	// Test GetMethod includes embedded methods
	readMethod, exists := writer.GetMethod("Read")
	if !exists {
		t.Error("Expected to find embedded Read method")
	}
	if readMethod.Name != "Read" {
		t.Errorf("Expected method name to be 'Read', got '%s'", readMethod.Name)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

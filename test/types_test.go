package test

import (
	"testing"

	"github.com/lengzhao/goscript/types"
)

// TestStructType tests struct type functionality
func TestStructType(t *testing.T) {
	// Create a struct type
	personStruct := types.NewStructType("Person")
	personStruct.AddField("name", types.StringType, "")
	personStruct.AddField("age", types.IntType, "")

	// Test field access
	if personStruct.TypeName() != "Person" {
		t.Errorf("Expected type name to be 'Person', got '%s'", personStruct.TypeName())
	}

	nameType, exists := personStruct.GetField("name")
	if !exists {
		t.Error("Expected to find 'name' field")
	}

	if nameType.TypeName() != "string" {
		t.Errorf("Expected 'name' field to be string, got '%s'", nameType.TypeName())
	}

	ageType, exists := personStruct.GetField("age")
	if !exists {
		t.Error("Expected to find 'age' field")
	}

	if ageType.TypeName() != "int" {
		t.Errorf("Expected 'age' field to be int, got '%s'", ageType.TypeName())
	}

	// Test method access
	personStruct.AddMethod("GetName", []types.IType{}, []types.IType{types.StringType})
	personStruct.AddMethod("SetAge", []types.IType{types.IntType}, []types.IType{})
	personStruct.AddMethod("GetAge", []types.IType{}, []types.IType{types.IntType})

	if !personStruct.HasMethod("GetName") {
		t.Error("Expected struct to have GetName method")
	}

	if !personStruct.HasMethod("SetAge") {
		t.Error("Expected struct to have SetAge method")
	}

	if !personStruct.HasMethod("GetAge") {
		t.Error("Expected struct to have GetAge method")
	}

	getNameMethod, exists := personStruct.GetMethod("GetName")
	if !exists {
		t.Error("Expected to find GetName method")
	}

	if len(getNameMethod.Returns) != 1 || getNameMethod.Returns[0].TypeName() != "string" {
		t.Error("Expected GetName method to return string")
	}

	// Test GetMethods
	methods := personStruct.GetMethods()
	if len(methods) != 3 {
		t.Errorf("Expected 3 methods, got %d", len(methods))
	}

	// Test string representation
	str := personStruct.String()
	if str == "" {
		t.Error("Expected non-empty string representation")
	}

	// Test default value
	defaultValue := personStruct.DefaultValue()
	if defaultValue == nil {
		t.Error("Expected non-nil default value")
	}

	// Test equality
	clone := personStruct.Clone()
	if !personStruct.Equals(clone) {
		t.Error("Expected struct to equal its clone")
	}
}

// TestInterfaceType tests interface type functionality
func TestInterfaceType(t *testing.T) {
	// Create an interface type
	shaperInterface := types.NewInterfaceType("Shaper")
	shaperInterface.AddMethod("Area", []types.IType{}, []types.IType{types.Float64Type})
	shaperInterface.AddMethod("Perimeter", []types.IType{}, []types.IType{types.Float64Type})

	// Test basic properties
	if shaperInterface.TypeName() != "Shaper" {
		t.Errorf("Expected type name to be 'Shaper', got '%s'", shaperInterface.TypeName())
	}

	// Test string representation
	str := shaperInterface.String()
	if str == "" {
		t.Error("Expected non-empty string representation")
	}

	// Test default value
	defaultValue := shaperInterface.DefaultValue()
	// Interface types should have nil as default value
	if defaultValue != nil {
		t.Error("Expected nil default value for interface type")
	}

	// Test equality
	clone := shaperInterface.Clone()
	if !shaperInterface.Equals(clone) {
		t.Error("Expected interface to equal its clone")
	}
}

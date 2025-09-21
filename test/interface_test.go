package test

import (
	"testing"

	"github.com/lengzhao/goscript/types"
)

// TestInterfaceBasic tests basic interface functionality
func TestInterfaceBasic(t *testing.T) {
	// Create an interface type
	shaperInterface := types.NewInterfaceType("Shaper")
	shaperInterface.AddMethod("Area", []types.IType{}, []types.IType{types.Float64Type})
	shaperInterface.AddMethod("Perimeter", []types.IType{}, []types.IType{types.Float64Type})

	// Test method access
	if !shaperInterface.HasMethod("Area") {
		t.Error("Expected interface to have Area method")
	}

	if !shaperInterface.HasMethod("Perimeter") {
		t.Error("Expected interface to have Perimeter method")
	}

	if shaperInterface.HasMethod("NonExistent") {
		t.Error("Expected interface to not have NonExistent method")
	}

	// Test GetMethod
	areaMethod, exists := shaperInterface.GetMethod("Area")
	if !exists {
		t.Error("Expected to find Area method")
	}

	if len(areaMethod.Returns) == 0 || areaMethod.Returns[0].TypeName() != "float64" {
		t.Error("Expected Area method to return float64")
	}

	// Test GetMethods
	methods := shaperInterface.GetMethods()
	if len(methods) != 2 {
		t.Errorf("Expected 2 methods, got %d", len(methods))
	}
}

// TestInterfaceEmbedded tests embedded interface functionality
func TestInterfaceEmbedded(t *testing.T) {
	// Create interfaces
	readerInterface := types.NewInterfaceType("Reader")
	readerInterface.AddMethod("Read", []types.IType{}, []types.IType{types.StringType})

	writerInterface := types.NewInterfaceType("Writer")
	writerInterface.AddMethod("Write", []types.IType{types.StringType}, []types.IType{types.IntType})

	// Create an interface that embeds other interfaces
	readWriterInterface := types.NewInterfaceType("ReadWriter")
	readWriterInterface.AddEmbedded(readerInterface)
	readWriterInterface.AddEmbedded(writerInterface)
	readWriterInterface.AddMethod("Close", []types.IType{}, []types.IType{types.IntType})

	// Test embedded methods
	if !readWriterInterface.HasMethod("Read") {
		t.Error("Expected ReadWriter to have Read method from embedded interface")
	}

	if !readWriterInterface.HasMethod("Write") {
		t.Error("Expected ReadWriter to have Write method from embedded interface")
	}

	if !readWriterInterface.HasMethod("Close") {
		t.Error("Expected ReadWriter to have Close method")
	}

	// Test GetMethods with embedded interfaces
	allMethods := readWriterInterface.GetMethods()
	// Should have Read, Write, and Close methods
	if len(allMethods) != 3 {
		t.Errorf("Expected 3 methods (including embedded), got %d", len(allMethods))
	}
}

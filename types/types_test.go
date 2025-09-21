package types

import (
	"testing"
)

func TestTypeExtensions(t *testing.T) {
	// Create a struct type
	personStruct := NewStructType("Person")
	personStruct.AddField("name", StringType, "")
	personStruct.AddField("age", IntType, "")

	// Create an interface type
	shaperInterface := NewInterfaceType("Shaper")
	shaperInterface.AddMethod("Area", []IType{}, []IType{Float64Type})

	// Test IsStructType
	if !IsStructType(personStruct) {
		t.Error("Expected personStruct to be recognized as a struct type")
	}

	if IsStructType(IntType) {
		t.Error("Expected IntType not to be recognized as a struct type")
	}

	// Test IsInterfaceType
	if !IsInterfaceType(shaperInterface) {
		t.Error("Expected shaperInterface to be recognized as an interface type")
	}

	if IsInterfaceType(IntType) {
		t.Error("Expected IntType not to be recognized as an interface type")
	}

	// Test AsStructType
	structType, ok := AsStructType(personStruct)
	if !ok {
		t.Error("Expected to be able to convert personStruct to StructType")
	}
	if structType.TypeName() != "Person" {
		t.Errorf("Expected struct type name to be 'Person', got '%s'", structType.TypeName())
	}

	_, ok = AsStructType(IntType)
	if ok {
		t.Error("Expected not to be able to convert IntType to StructType")
	}

	// Test AsInterfaceType
	interfaceType, ok := AsInterfaceType(shaperInterface)
	if !ok {
		t.Error("Expected to be able to convert shaperInterface to InterfaceType")
	}
	if interfaceType.TypeName() != "Shaper" {
		t.Errorf("Expected interface type name to be 'Shaper', got '%s'", interfaceType.TypeName())
	}

	_, ok = AsInterfaceType(IntType)
	if ok {
		t.Error("Expected not to be able to convert IntType to InterfaceType")
	}

	// Test struct field access
	nameType, exists := structType.GetField("name")
	if !exists {
		t.Error("Expected 'name' field to exist in Person struct")
	}
	if nameType.TypeName() != "string" {
		t.Errorf("Expected 'name' field type to be 'string', got '%s'", nameType.TypeName())
	}

	ageType, exists := structType.GetField("age")
	if !exists {
		t.Error("Expected 'age' field to exist in Person struct")
	}
	if ageType.TypeName() != "int" {
		t.Errorf("Expected 'age' field type to be 'int', got '%s'", ageType.TypeName())
	}

	// Test interface method access
	areaMethod, exists := interfaceType.GetMethod("Area")
	if !exists {
		t.Error("Expected 'Area' method to exist in Shaper interface")
	}
	if areaMethod.Name != "Area" {
		t.Errorf("Expected method name to be 'Area', got '%s'", areaMethod.Name)
	}
	if len(areaMethod.Returns) != 1 {
		t.Errorf("Expected 'Area' method to have 1 return value, got %d", len(areaMethod.Returns))
	}
	if areaMethod.Returns[0].TypeName() != "float64" {
		t.Errorf("Expected 'Area' method return type to be 'float64', got '%s'", areaMethod.Returns[0].TypeName())
	}
}

func TestMethodAccess(t *testing.T) {
	// Create a struct type with methods
	employeeStruct := NewStructType("Employee")
	employeeStruct.AddField("name", StringType, "")
	employeeStruct.AddField("salary", Float64Type, "")

	// Add methods to the struct
	employeeStruct.AddMethod("GetName", []IType{}, []IType{StringType})
	employeeStruct.AddMethod("SetSalary", []IType{Float64Type}, []IType{})
	employeeStruct.AddMethod("GetSalary", []IType{}, []IType{Float64Type})

	// Test HasMethod
	if !employeeStruct.HasMethod("GetName") {
		t.Error("Expected Employee struct to have GetName method")
	}

	if employeeStruct.HasMethod("NonExistentMethod") {
		t.Error("Expected Employee struct not to have NonExistentMethod")
	}

	// Test GetMethod
	getNameMethod, exists := employeeStruct.GetMethod("GetName")
	if !exists {
		t.Error("Expected to find GetName method")
	}
	if getNameMethod.Name != "GetName" {
		t.Errorf("Expected method name to be 'GetName', got '%s'", getNameMethod.Name)
	}
	if len(getNameMethod.Params) != 0 {
		t.Errorf("Expected GetName method to have 0 parameters, got %d", len(getNameMethod.Params))
	}
	if len(getNameMethod.Returns) != 1 {
		t.Errorf("Expected GetName method to have 1 return value, got %d", len(getNameMethod.Returns))
	}
	if getNameMethod.Returns[0].TypeName() != "string" {
		t.Errorf("Expected GetName method return type to be 'string', got '%s'", getNameMethod.Returns[0].TypeName())
	}

	// Test GetMethods
	methods := employeeStruct.GetMethods()
	if len(methods) != 3 {
		t.Errorf("Expected Employee struct to have 3 methods, got %d", len(methods))
	}

	// Create an interface with embedded interface
	readerInterface := NewInterfaceType("Reader")
	readerInterface.AddMethod("Read", []IType{}, []IType{StringType})

	writerInterface := NewInterfaceType("Writer")
	writerInterface.AddMethod("Write", []IType{StringType}, []IType{IntType})

	// Create an interface that embeds other interfaces
	readWriterInterface := NewInterfaceType("ReadWriter")
	readWriterInterface.AddEmbedded(readerInterface)
	readWriterInterface.AddEmbedded(writerInterface)
	readWriterInterface.AddMethod("Close", []IType{}, []IType{IntType})

	// Test HasMethod with embedded interfaces
	if !readWriterInterface.HasMethod("Read") {
		t.Error("Expected ReadWriter interface to have Read method from embedded Reader")
	}

	if !readWriterInterface.HasMethod("Write") {
		t.Error("Expected ReadWriter interface to have Write method from embedded Writer")
	}

	if !readWriterInterface.HasMethod("Close") {
		t.Error("Expected ReadWriter interface to have Close method")
	}

	// Test GetMethod with embedded interfaces
	readMethod, exists := readWriterInterface.GetMethod("Read")
	if !exists {
		t.Error("Expected to find Read method from embedded interface")
	}
	if readMethod.Name != "Read" {
		t.Errorf("Expected method name to be 'Read', got '%s'", readMethod.Name)
	}

	// Test GetMethods with embedded interfaces
	allMethods := readWriterInterface.GetMethods()
	if len(allMethods) != 3 {
		t.Errorf("Expected ReadWriter interface to have 3 methods (including embedded), got %d", len(allMethods))
	}
}

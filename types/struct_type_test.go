package types

import (
	"reflect"
	"strings"
	"testing"
)

func TestStructType(t *testing.T) {
	// Create a new struct type
	person := NewStructType("Person")
	person.AddField("name", StringType.Clone(), "")
	person.AddField("age", IntType.Clone(), "")

	// Test TypeName
	if person.TypeName() != "Person" {
		t.Errorf("Expected TypeName to be 'Person', got '%s'", person.TypeName())
	}

	// Test String - check that it contains the expected fields (order may vary)
	actualStr := person.String()
	if !strings.Contains(actualStr, "name string") || !strings.Contains(actualStr, "age int") {
		t.Errorf("Expected String to contain 'name string' and 'age int', got '%s'", actualStr)
	}

	// Test Size
	if person.Size() != 24 {
		t.Errorf("Expected Size to be 24, got %d", person.Size())
	}

	// Test Kind
	if person.Kind() != reflect.Struct {
		t.Errorf("Expected Kind to be reflect.Struct, got %v", person.Kind())
	}

	// Test DefaultValue
	defaultValue := person.DefaultValue()
	if defaultValue == nil {
		t.Error("Expected DefaultValue to not be nil")
	}

	defaultMap, ok := defaultValue.(map[string]interface{})
	if !ok {
		t.Error("Expected DefaultValue to be a map")
	}

	if len(defaultMap) != 2 {
		t.Errorf("Expected DefaultValue to have 2 fields, got %d", len(defaultMap))
	}

	if defaultMap["name"] != "" {
		t.Error("Expected default name to be empty string")
	}

	if defaultMap["age"] != 0 {
		t.Error("Expected default age to be 0")
	}

	// Test Clone
	clone := person.Clone()
	if !person.Equals(clone) {
		t.Error("Expected clone to be equal to original")
	}

	// Test GetField
	nameType, exists := person.GetField("name")
	if !exists {
		t.Error("Expected 'name' field to exist")
	}
	if nameType.TypeName() != "string" {
		t.Errorf("Expected 'name' field type to be 'string', got '%s'", nameType.TypeName())
	}

	// Test HasField
	if !person.HasField("name") {
		t.Error("Expected HasField('name') to return true")
	}
	if person.HasField("nonexistent") {
		t.Error("Expected HasField('nonexistent') to return false")
	}

	// Test GetFieldNames
	fieldNames := person.GetFieldNames()
	if len(fieldNames) != 2 {
		t.Errorf("Expected 2 field names, got %d", len(fieldNames))
	}
}

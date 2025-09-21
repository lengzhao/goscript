package test

import (
	"testing"

	goscript "github.com/lengzhao/goscript"
	"github.com/lengzhao/goscript/builtin"
	execContext "github.com/lengzhao/goscript/context"
)

func TestStringsModule(t *testing.T) {
	// Test strings module functions from builtin package
	stringsFuncs, exists := builtin.GetModuleFunctions("strings")
	if !exists {
		t.Fatal("strings module should exist")
	}

	// Test Contains function
	containsFunc, exists := stringsFuncs["Contains"]
	if !exists {
		t.Fatal("contains function should exist in strings module")
	}

	result1, err := containsFunc("hello world", "world")
	if err != nil {
		t.Fatalf("Failed to call contains function: %v", err)
	}
	if result1 != true {
		t.Errorf("Expected contains('hello world', 'world') to be true, got %v", result1)
	}

	result2, err := containsFunc("hello world", "test")
	if err != nil {
		t.Fatalf("Failed to call contains function: %v", err)
	}
	if result2 != false {
		t.Errorf("Expected contains('hello world', 'test') to be false, got %v", result2)
	}

	// Test ToLower function
	toLowerFunc, exists := stringsFuncs["ToLower"]
	if !exists {
		t.Fatal("toLower function should exist in strings module")
	}

	result3, err := toLowerFunc("HELLO WORLD")
	if err != nil {
		t.Fatalf("Failed to call toLower function: %v", err)
	}
	if result3 != "hello world" {
		t.Errorf("Expected toLower('HELLO WORLD') to be 'hello world', got %v", result3)
	}

	// Test Split function
	splitFunc, exists := stringsFuncs["Split"]
	if !exists {
		t.Fatal("split function should exist in strings module")
	}

	result4, err := splitFunc("a,b,c", ",")
	if err != nil {
		t.Fatalf("Failed to call split function: %v", err)
	}
	if slice, ok := result4.([]interface{}); ok {
		if len(slice) != 3 || slice[0] != "a" || slice[1] != "b" || slice[2] != "c" {
			t.Errorf("Expected split('a,b,c', ',') to be ['a', 'b', 'c'], got %v", slice)
		}
	} else {
		t.Errorf("Expected split to return []interface{}, got %T", result4)
	}

	t.Log("Strings module functions test passed")
}

func TestJsonModule(t *testing.T) {
	// Test json module functions from builtin package
	jsonFuncs, exists := builtin.GetModuleFunctions("json")
	if !exists {
		t.Fatal("json module should exist")
	}

	// Test Marshal function
	marshalFunc, exists := jsonFuncs["Marshal"]
	if !exists {
		t.Fatal("marshal function should exist in json module")
	}

	// Test with a map
	testMap := map[string]interface{}{
		"name": "John",
		"age":  30,
	}
	result1, err := marshalFunc(testMap)
	if err != nil {
		t.Fatalf("Failed to call marshal function: %v", err)
	}
	if result1 != `{"age":30,"name":"John"}` && result1 != `{"name":"John","age":30}` {
		t.Errorf("Expected marshal to return JSON string, got %v", result1)
	}

	// Test Unmarshal function
	unmarshalFunc, exists := jsonFuncs["Unmarshal"]
	if !exists {
		t.Fatal("unmarshal function should exist in json module")
	}

	jsonStr := `{"name":"John","age":30}`
	result2, err := unmarshalFunc(jsonStr)
	if err != nil {
		t.Fatalf("Failed to call unmarshal function: %v", err)
	}

	// Check if result is a map with the expected values
	if resultMap, ok := result2.(map[string]interface{}); ok {
		if resultMap["name"] != "John" || resultMap["age"] != float64(30) {
			t.Errorf("Expected unmarshal to return map with correct values, got %v", resultMap)
		}
	} else {
		t.Errorf("Expected unmarshal to return map[string]interface{}, got %T", result2)
	}

	t.Log("JSON module functions test passed")
}

func TestModuleFunctionalityInScript(t *testing.T) {
	// Test using module functions through the script API
	script := goscript.NewScript([]byte(""))

	// Set security context to allow strings and json modules
	securityCtx := &execContext.SecurityContext{
		AllowedModules: []string{"strings", "json"},
	}
	script.SetSecurityContext(securityCtx)

	// Import modules explicitly
	err := script.ImportModule("strings", "json")
	if err != nil {
		t.Fatalf("Failed to import modules: %v", err)
	}

	// Test calling strings module functions through the script
	result1, err := script.CallFunction("strings.ToLower", "HELLO WORLD")
	if err != nil {
		t.Fatalf("Failed to call strings.ToLower function: %v", err)
	}
	if result1 != "hello world" {
		t.Errorf("Expected strings.ToLower('HELLO WORLD') to be 'hello world', got %v", result1)
	}

	result2, err := script.CallFunction("strings.Contains", "hello world", "world")
	if err != nil {
		t.Fatalf("Failed to call strings.Contains function: %v", err)
	}
	if result2 != true {
		t.Errorf("Expected strings.Contains('hello world', 'world') to be true, got %v", result2)
	}

	// Test calling json module functions through the script
	testMap := map[string]interface{}{
		"name": "John",
		"age":  30,
	}
	result3, err := script.CallFunction("json.Marshal", testMap)
	if err != nil {
		t.Fatalf("Failed to call json.Marshal function: %v", err)
	}
	if result3 != `{"age":30,"name":"John"}` && result3 != `{"name":"John","age":30}` {
		t.Errorf("Expected json.Marshal to return JSON string, got %v", result3)
	}

	jsonStr := `{"name":"John","age":30}`
	result4, err := script.CallFunction("json.Unmarshal", jsonStr)
	if err != nil {
		t.Fatalf("Failed to call json.Unmarshal function: %v", err)
	}

	// Check if result is a map with the expected values
	if resultMap, ok := result4.(map[string]interface{}); ok {
		if resultMap["name"] != "John" || resultMap["age"] != float64(30) {
			t.Errorf("Expected json.Unmarshal to return map with correct values, got %v", resultMap)
		}
	} else {
		t.Errorf("Expected json.Unmarshal to return map[string]interface{}, got %T", result4)
	}

	t.Log("Module functionality in script test passed")
}

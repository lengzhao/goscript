package builtin

import (
	"testing"
)

func TestLen(t *testing.T) {
	// Test string length
	result, err := Len("hello")
	if err != nil {
		t.Errorf("Len failed for string: %v", err)
	}
	if result != 5 {
		t.Errorf("Expected length 5 for 'hello', got %v", result)
	}

	// Test slice length
	slice := []interface{}{1, 2, 3}
	result, err = Len(slice)
	if err != nil {
		t.Errorf("Len failed for slice: %v", err)
	}
	if result != 3 {
		t.Errorf("Expected length 3 for slice, got %v", result)
	}

	// Test map length
	m := map[string]interface{}{"a": 1, "b": 2}
	result, err = Len(m)
	if err != nil {
		t.Errorf("Len failed for map: %v", err)
	}
	if result != 2 {
		t.Errorf("Expected length 2 for map, got %v", result)
	}

	// Test with wrong number of arguments
	_, err = Len("hello", "world")
	if err == nil {
		t.Error("Expected error for wrong number of arguments")
	}
}

func TestMake(t *testing.T) {
	// Test make slice
	result, err := Make("slice", 3)
	if err != nil {
		t.Errorf("Make failed: %v", err)
	}
	if slice, ok := result.([]interface{}); !ok {
		t.Errorf("Expected slice, got %T", result)
	} else if len(slice) != 3 {
		t.Errorf("Expected slice of length 3, got %d", len(slice))
	}

	// Test make slice with capacity
	result, err = Make("slice", 2, 5)
	if err != nil {
		t.Errorf("Make failed with capacity: %v", err)
	}
	if slice, ok := result.([]interface{}); !ok {
		t.Errorf("Expected slice, got %T", result)
	} else if len(slice) != 5 {
		t.Errorf("Expected slice of length 5, got %d", len(slice))
	}

	// Test with no arguments
	_, err = Make()
	if err == nil {
		t.Error("Expected error for no arguments")
	}
}

func TestCopy(t *testing.T) {
	// Test copy
	src := []interface{}{1, 2, 3, 4, 5}
	dst := make([]interface{}, 3)
	result, err := Copy(dst, src)
	if err != nil {
		t.Errorf("Copy failed: %v", err)
	}
	if result != 3 {
		t.Errorf("Expected to copy 3 elements, got %v", result)
	}
	if dst[0] != 1 || dst[1] != 2 || dst[2] != 3 {
		t.Errorf("Copy did not work correctly: %v", dst)
	}

	// Test with wrong number of arguments
	_, err = Copy(dst)
	if err == nil {
		t.Error("Expected error for wrong number of arguments")
	}

	// Test with non-slice arguments
	_, err = Copy("not a slice", src)
	if err == nil {
		t.Error("Expected error for non-slice arguments")
	}
}

func TestPrint(t *testing.T) {
	// Test print (this will output to stdout)
	_, err := Print("hello", "world", 123)
	if err != nil {
		t.Errorf("Print failed: %v", err)
	}
}

func TestInt(t *testing.T) {
	// Test int conversion from int
	result, err := Int(42)
	if err != nil {
		t.Errorf("Int failed for int: %v", err)
	}
	if result != 42 {
		t.Errorf("Expected 42 for int 42, got %v", result)
	}

	// Test int conversion from float64
	result, err = Int(3.14)
	if err != nil {
		t.Errorf("Int failed for float64: %v", err)
	}
	if result != 3 {
		t.Errorf("Expected 3 for float64 3.14, got %v", result)
	}

	// Test int conversion from string
	result, err = Int("123")
	if err != nil {
		t.Errorf("Int failed for string: %v", err)
	}
	// For now, we return 0 for string conversion
	if result != 0 {
		t.Errorf("Expected 0 for string '123', got %v", result)
	}

	// Test with wrong number of arguments
	_, err = Int(1, 2)
	if err == nil {
		t.Error("Expected error for wrong number of arguments")
	}

	// Test with unsupported type
	_, err = Int([]int{1, 2, 3})
	if err == nil {
		t.Error("Expected error for unsupported type")
	}
}

func TestStringsModule(t *testing.T) {
	// Test contains function
	containsFunc, exists := StringsModule["Contains"]
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

	// Test toLower function
	toLowerFunc, exists := StringsModule["ToLower"]
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

	// Test split function
	splitFunc, exists := StringsModule["Split"]
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
}

func TestJSONModule(t *testing.T) {
	// Test marshal function
	marshalFunc, exists := JSONModule["Marshal"]
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

	// Test unmarshal function
	unmarshalFunc, exists := JSONModule["Unmarshal"]
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
}

func TestGetModuleFunctions(t *testing.T) {
	// Test getting strings module functions
	stringsFuncs, exists := GetModuleFunctions("strings")
	if !exists {
		t.Error("strings module should exist")
	}
	if len(stringsFuncs) == 0 {
		t.Error("strings module should have functions")
	}

	// Test getting json module functions
	jsonFuncs, exists := GetModuleFunctions("json")
	if !exists {
		t.Error("json module should exist")
	}
	if len(jsonFuncs) == 0 {
		t.Error("json module should have functions")
	}

	// Test getting non-existent module
	_, exists = GetModuleFunctions("nonexistent")
	if exists {
		t.Error("nonexistent module should not exist")
	}
}

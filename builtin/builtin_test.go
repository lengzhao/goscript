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

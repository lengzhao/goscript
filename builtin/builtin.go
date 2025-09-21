// Package builtin provides built-in functions for GoScript
package builtin

import (
	"fmt"
	"reflect"

	"github.com/lengzhao/goscript/types"
)

// Function represents a built-in function
// Note: This must match the Function type in runtime/runtime.go
type Function = types.Function

// BuiltInFunctions holds all built-in functions
var BuiltInFunctions = map[string]Function{
	"len":   Len,
	"make":  Make,
	"copy":  Copy,
	"print": Print,
	"int":   Int,
}

// Len returns the length of a string, array, slice, or map
func Len(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("len expects 1 argument, got %d", len(args))
	}

	switch v := args[0].(type) {
	case string:
		return len(v), nil
	case []interface{}:
		return len(v), nil
	case map[string]interface{}:
		return len(v), nil
	default:
		// Use reflection for other types
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Slice, reflect.Array, reflect.Map, reflect.String:
			return rv.Len(), nil
		default:
			return nil, fmt.Errorf("len: unsupported type %T", v)
		}
	}
}

// Make creates a slice, map, or channel
func Make(args ...interface{}) (interface{}, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("make expects at least 1 argument, got %d", len(args))
	}

	// For now, we only support slice creation
	// In a full implementation, we would need type information
	switch args[0].(type) {
	case string:
		// Create a slice
		if len(args) >= 2 {
			if size, ok := args[1].(int); ok {
				slice := make([]interface{}, size)
				// If capacity is provided, use it
				if len(args) >= 3 {
					if cap, ok := args[2].(int); ok && cap > size {
						// Extend the slice to capacity
						for i := size; i < cap; i++ {
							slice = append(slice, nil)
						}
					}
				}
				return slice, nil
			}
		}
		return make([]interface{}, 0), nil
	default:
		return nil, fmt.Errorf("make: unsupported type %T", args[0])
	}
}

// Copy copies elements from a source slice to a destination slice
func Copy(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("copy expects 2 arguments, got %d", len(args))
	}

	dst, ok1 := args[0].([]interface{})
	src, ok2 := args[1].([]interface{})
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("copy: both arguments must be slices")
	}

	// Copy elements
	count := 0
	for i := 0; i < len(dst) && i < len(src); i++ {
		dst[i] = src[i]
		count++
	}

	return count, nil
}

// Print prints the arguments to stdout
func Print(args ...interface{}) (interface{}, error) {
	for i, arg := range args {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(arg)
	}
	fmt.Println()
	return nil, nil
}

// Int converts a value to an integer
func Int(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("int expects 1 argument, got %d", len(args))
	}

	switch v := args[0].(type) {
	case int:
		return v, nil
	case float64:
		return int(v), nil
	case string:
		// In a full implementation, we would parse the string
		// For now, we'll just return 0
		return 0, nil
	default:
		return 0, fmt.Errorf("int: unsupported type %T", v)
	}
}

package test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	goscript "github.com/lengzhao/goscript"
)

// TestScriptsInDataFolder tests all .gs scripts in the test/data folder
func TestScriptsInDataFolder(t *testing.T) {
	// Get all .gs files in test/data folder
	dataDir := "./data"
	files, err := os.ReadDir(dataDir)
	if err != nil {
		t.Fatalf("Failed to read data directory: %v", err)
	}

	runCase := "complex.gs"
	// Test each .gs file
	for _, file := range files {
		if runCase != "" && runCase != file.Name() {
			continue
		}
		if filepath.Ext(file.Name()) == ".gs" {
			t.Run(file.Name(), func(t *testing.T) {
				testScriptFile(t, filepath.Join(dataDir, file.Name()))
			})
		}
	}
}

// testScriptFile tests a single script file
func testScriptFile(t *testing.T, filePath string) {
	// Read the script file
	source, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read script file %s: %v", filePath, err)
	}

	// Create a new script
	script := goscript.NewScript(source)
	script.SetDebug(true)

	// // Register required functions based on script name
	// switch filepath.Base(filePath) {
	// case "add.gs":
	// 	// Register add function
	// 	script.AddFunction("add", goscript.NewSimpleFunction("add", func(args ...interface{}) (interface{}, error) {
	// 		a := args[0].(int)
	// 		b := args[1].(int)
	// 		return a + b, nil
	// 	}))
	// case "sub.gs":
	// 	// Register sub function
	// 	script.AddFunction("sub", goscript.NewSimpleFunction("sub", func(args ...interface{}) (interface{}, error) {
	// 		a := args[0].(int)
	// 		b := args[1].(int)
	// 		return a - b, nil
	// 	}))
	// case "complex.gs":
	// 	// Register factorial function
	// 	script.AddFunction("factorial", goscript.NewSimpleFunction("factorial", func(args ...interface{}) (interface{}, error) {
	// 		n := args[0].(int)
	// 		if n <= 1 {
	// 			return 1, nil
	// 		}
	// 		result := 1
	// 		for i := 2; i <= n; i++ {
	// 			result *= i
	// 		}
	// 		return result, nil
	// 	}))

	// 	// Register fibonacci function
	// 	script.AddFunction("fibonacci", goscript.NewSimpleFunction("fibonacci", func(args ...interface{}) (interface{}, error) {
	// 		n := args[0].(int)
	// 		if n <= 1 {
	// 			return n, nil
	// 		}
	// 		a, b := 0, 1
	// 		for i := 2; i <= n; i++ {
	// 			a, b = b, a+b
	// 		}
	// 		return b, nil
	// 	}))
	// case "conditional.gs":
	// 	// Register max function
	// 	script.AddFunction("max", goscript.NewSimpleFunction("max", func(args ...interface{}) (interface{}, error) {
	// 		a := args[0].(int)
	// 		b := args[1].(int)
	// 		if a > b {
	// 			return a, nil
	// 		}
	// 		return b, nil
	// 	}))
	// case "function_call.gs":
	// 	// Register calculate function
	// 	script.AddFunction("calculate", goscript.NewSimpleFunction("calculate", func(args ...interface{}) (interface{}, error) {
	// 		return 15, nil // 10 + 5
	// 	}))
	// case "loop.gs":
	// 	// Register sum function
	// 	script.AddFunction("sum", goscript.NewSimpleFunction("sum", func(args ...interface{}) (interface{}, error) {
	// 		n := args[0].(int)
	// 		total := 0
	// 		for i := 1; i <= n; i++ {
	// 			total += i
	// 		}
	// 		return total, nil
	// 	}))
	// }

	// Run the script
	result, err := script.Run()
	if err != nil {
		t.Fatalf("Failed to run script %s: %v", filePath, err)
	}

	// Log the result
	t.Logf("Script %s executed successfully, result: %v", filePath, result)

	// Perform basic validation based on script name
	validateScriptResult(t, filePath, result)
}

// validateScriptResult performs basic validation on script results
func validateScriptResult(t *testing.T, filePath string, result interface{}) {
	switch filepath.Base(filePath) {
	case "hello.gs":
		if result != "Hello, World!" {
			t.Errorf("Expected 'Hello, World!', got %v", result)
		}
	case "add.gs":
		if result != 8 {
			t.Errorf("Expected 8, got %v", result)
		}
	case "loop.gs":
		if result != 55 {
			t.Errorf("Expected 55, got %v", result)
		}
	case "conditional.gs":
		if result != 20 {
			t.Errorf("Expected 20, got %v", result)
		}
	case "function_call.gs":
		if result != 15 {
			t.Errorf("Expected 15, got %v", result)
		}
	case "complex.gs":
		// factorial(5) = 120, fibonacci(10) = 55, result = 120 + 55 = 175
		if result != 175 {
			t.Errorf("Expected 175, got %v", result)
		}
	default:
		// For other scripts, just ensure they executed without error
		t.Logf("Script %s executed successfully with result: %v", filePath, result)
	}
}

// Example of testing with custom functions
func TestScriptWithCustomFunction(t *testing.T) {
	// Simple script that uses a custom function
	source := []byte(`
package main

func main() {
    result := power(2, 3)
    return result
}
`)

	// Create script
	script := goscript.NewScript(source)

	// Define a custom power function
	powerFunc := goscript.NewSimpleFunction("power", func(args ...interface{}) (interface{}, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("power function requires 2 arguments")
		}
		base, ok1 := args[0].(int)
		exp, ok2 := args[1].(int)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("power function requires integer arguments")
		}
		result := 1
		for i := 0; i < exp; i++ {
			result *= base
		}
		return result, nil
	})

	// Register the custom function
	err := script.AddFunction("power", powerFunc)
	if err != nil {
		t.Fatalf("Failed to register custom function: %v", err)
	}

	// Run the script
	result, err := script.Run()
	if err != nil {
		t.Fatalf("Failed to run script with custom function: %v", err)
	}

	// Validate result (2^3 = 8)
	if result != 8 {
		t.Errorf("Expected 8, got %v", result)
	}
}

// Example of testing with variables
func TestScriptWithVariables(t *testing.T) {
	// Simple script that uses variables
	source := []byte(`
package main

func main() {
    x := 10
    y := 20
    return x + y
}
`)

	// Create script
	script := goscript.NewScript(source)

	// Run the script
	result, err := script.Run()
	if err != nil {
		t.Fatalf("Failed to run script with variables: %v", err)
	}

	// Validate result (10 + 20 = 30)
	if result != 30 {
		t.Errorf("Expected 30, got %v", result)
	}
}

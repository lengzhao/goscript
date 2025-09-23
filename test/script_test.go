package test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	goscript "github.com/lengzhao/goscript"
	execContext "github.com/lengzhao/goscript/context"
)

// TestScriptsInDataFolder tests all .gs scripts in the test/data folder
func TestScriptsInDataFolder(t *testing.T) {
	// Get all .gs files in test/data folder
	dataDir := "./data"
	files, err := os.ReadDir(dataDir)
	if err != nil {
		t.Fatalf("Failed to read data directory: %v", err)
	}

	runCase := "" // Run all test cases
	// runCase := "range_simple.gs" // Uncomment to run a specific test case

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
	script.SetDebug(true) // Enable debug mode

	// 设置较大的指令数限制
	securityCtx := &execContext.SecurityContext{
		MaxExecutionTime:  5 * time.Second,
		MaxMemoryUsage:    10 * 1024 * 1024, // 10MB
		AllowedModules:    []string{"fmt", "math", "strings", "json"},
		ForbiddenKeywords: []string{"unsafe"},
		AllowCrossModule:  true,
		MaxInstructions:   10000, // 设置较小的指令数限制用于简单脚本
	}
	script.SetSecurityContext(securityCtx)
	// script.ImportModule(builtin.ListAllModules()...)

	ctx := context.Background()
	// ctx1, cancel := context.WithTimeout(ctx, 2*time.Second)
	// defer cancel()
	// Run the script
	result, err := script.RunContext(ctx)
	if err != nil {
		t.Fatalf("Failed to run script %s: %v", filePath, err)
	}

	// Log the result
	t.Logf("Script %s executed successfully, result: %v", filePath, result)

	// Perform basic validation based on script name
	validateScriptResult(t, filePath, result)
}

// validateScriptResult performs validation on script results using expected results from JSON
func validateScriptResult(t *testing.T, filePath string, result interface{}) {
	data, err := os.ReadFile("./data/result.json")
	if err != nil {
		t.Fatalf("Failed to load expected results: %v", err)
	}
	expectedResults := make(map[string]interface{})
	err = json.Unmarshal(data, &expectedResults)
	if err != nil {
		t.Fatalf("Failed to unmarshal expected results: %v", err)
	}
	// Get the base file name
	baseName := filepath.Base(filePath)

	expected, exists := expectedResults[baseName]
	if !exists {
		t.Errorf("No expected result found for %s", baseName)
		return
	}
	v1 := fmt.Sprint(result)
	v2 := fmt.Sprint(expected)
	if v1 != v2 {
		t.Errorf("Expected %v, got %v", v2, v1)
	}
}

// Example of testing with custom functions
func TestScriptWithCustomFunction(t *testing.T) {
	// Simple script that uses a custom function
	source := []byte(`
package test

func main() {
    result := power(2, 3)
    return result
}
`)

	// Create script
	script := goscript.NewScript(source)

	// Register the custom function
	err := script.AddFunction("power", func(args ...interface{}) (interface{}, error) {
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
package test

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

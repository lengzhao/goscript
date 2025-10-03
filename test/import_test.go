package test

import (
	"testing"

	goscript "github.com/lengzhao/goscript"
	"github.com/lengzhao/goscript/builtin"
)

func TestImportInstruction(t *testing.T) {
	// Simple script that uses import
	source := []byte(`
package main

import "strings"

func main() {
    s := "Hello, World!"
    if strings.HasPrefix(s, "Hello") {
        return "YES"
    } else {
        return "NO"
    }
}
`)

	// Create script
	script := goscript.NewScript(source)
	script.SetDebug(true) // Disable debug mode for cleaner output

	// Run the script
	result, err := script.Run()
	if err != nil {
		t.Fatalf("Failed to run script with import: %v", err)
	}

	// Validate result
	if result != "YES" {
		t.Errorf("Expected YES, got %v", result)
	}
}

func TestImportMultipleModules(t *testing.T) {
	// Script that imports multiple modules
	source := []byte(`
package main

import "strings"
import "fmt"

func main() {
    s := "hello world"
    upper := strings.ToUpper(s)
    msg := fmt.Printf("Uppercase: %s", upper)
    return msg
}
`)

	// Create script
	script := goscript.NewScript(source)
	script.SetDebug(false) // Disable debug mode for cleaner output

	// Register builtin modules
	modules := []string{"strings", "fmt"}
	for _, moduleName := range modules {
		moduleExecutor, exists := builtin.GetModuleExecutor(moduleName)
		if exists {
			script.RegisterModule(moduleName, moduleExecutor)
		}
	}

	// Run the script
	result, err := script.Run()
	if err != nil {
		t.Fatalf("Failed to run script with multiple imports: %v", err)
	}

	// Validate result
	expected := "Uppercase: HELLO WORLD"
	if result != expected {
		t.Errorf("Expected %s, got %v", expected, result)
	}
}

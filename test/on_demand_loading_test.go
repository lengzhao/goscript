package test

import (
	"testing"

	goscript "github.com/lengzhao/goscript"
)

// TestOnDemandModuleLoading tests that modules are loaded on-demand
func TestOnDemandModuleLoading(t *testing.T) {
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
	script.SetDebug(false) // Disable debug mode for cleaner output

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

// TestMultipleModuleFunctions tests that only used functions are loaded
func TestMultipleModuleFunctions(t *testing.T) {
	// Script that imports multiple modules but only uses some functions
	source := []byte(`
package main

import "strings"
import "fmt"

func main() {
    s := "hello world"
    upper := strings.ToUpper(s)
    // Note: we don't use fmt.Printf, so it shouldn't be loaded
    return upper
}
`)

	// Create script
	script := goscript.NewScript(source)
	script.SetDebug(false) // Disable debug mode for cleaner output

	// Run the script
	result, err := script.Run()
	if err != nil {
		t.Fatalf("Failed to run script with multiple imports: %v", err)
	}

	// Validate result
	expected := "HELLO WORLD"
	if result != expected {
		t.Errorf("Expected %s, got %v", expected, result)
	}
}

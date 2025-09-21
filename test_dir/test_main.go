package main

import (
	"fmt"
	"os"

	"github.com/lengzhao/goscript"
)

func main() {
	// Test simple script
	fmt.Println("=== Testing simple script ===")
	testScript("../test_simple.gs")
	
	// Test complex script
	fmt.Println("\n=== Testing complex script ===")
	testScript("../test_complex.gs")
}

func testScript(scriptPath string) {
	// Read the test script
	scriptContent, err := os.ReadFile(scriptPath)
	if err != nil {
		fmt.Printf("Failed to read script: %v\n", err)
		return
	}

	// Create a new script
	script := goscript.NewScript(scriptContent)
	script.SetDebug(true) // Set to true for detailed output

	// Run the script
	result, err := script.Run()
	if err != nil {
		fmt.Printf("Failed to run script: %v\n", err)
		return
	}

	fmt.Printf("Script result: %v\n", result)
}
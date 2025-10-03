// Package main provides an advanced example of using GoScript
package main

import (
	"fmt"
	"os"

	"github.com/lengzhao/goscript"
)

func main() {
	// Read the script file
	scriptFile := "script.gs"
	source, err := os.ReadFile(scriptFile)
	if err != nil {
		fmt.Printf("Failed to read script file: %v\n", err)
		return
	}

	// Create a script
	script := goscript.NewScript(source)
	script.SetDebug(true) // Enable debug mode

	// Execute the script
	result, err := script.Run()
	if err != nil {
		fmt.Printf("Execution error: %v\n", err)
		return
	}

	fmt.Printf("Script result: %v\n", result)
}

package main

import (
	"fmt"
	"os"

	goscript "github.com/lengzhao/goscript"
)

func main() {
	// Read script file from command line argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <script_file>")
		return
	}

	scriptFile := os.Args[1]
	source, err := os.ReadFile(scriptFile)
	if err != nil {
		fmt.Printf("Failed to read script file: %v\n", err)
		return
	}

	// Create a new script
	script := goscript.NewScript(source)
	script.SetDebug(true) // Enable debug mode

	// Run the script
	result, err := script.Run()
	if err != nil {
		fmt.Printf("Failed to run script: %v\n", err)
		return
	}

	// Print the result
	fmt.Printf("Script executed successfully, result: %v\n", result)
}

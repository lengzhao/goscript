package main

import (
	"fmt"
	"os"
	"time"

	goscript "github.com/lengzhao/goscript"
	execContext "github.com/lengzhao/goscript/context"
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

	// Set security context with higher instruction limit
	securityCtx := &execContext.SecurityContext{
		MaxExecutionTime:  5 * time.Second,
		MaxMemoryUsage:    10 * 1024 * 1024, // 10MB
		AllowedModules:    []string{"fmt", "math"},
		ForbiddenKeywords: []string{"unsafe"},
		AllowCrossModule:  true,
		MaxInstructions:   100000, // Increased instruction limit
	}
	script.SetSecurityContext(securityCtx)

	// Run the script
	result, err := script.Run()
	if err != nil {
		fmt.Printf("Failed to run script: %v\n", err)
		return
	}

	// Print the result
	fmt.Printf("Script executed successfully, result: %v\n", result)
}

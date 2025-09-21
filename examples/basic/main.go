// Package main provides a basic example of using GoScript
package main

import (
	"fmt"

	"github.com/lengzhao/goscript"
)

func main() {
	// Create a script
	source := `
package main

func main() {
	x := 10
	y := 20
	return x + y
}
`

	script := goscript.NewScript([]byte(source))

	// // Set security context
	// securityCtx := &goscript.SecurityContext{
	// 	MaxExecutionTime: 5 * time.Second,
	// 	MaxMemoryUsage:   10 * 1024 * 1024, // 10MB
	// 	AllowedModules:   []string{"fmt"},
	// }

	// script.SetSecurityContext(securityCtx)

	// Add a custom variable
	script.AddVariable("customVar", 42)

	// Execute the script
	result, err := script.Run()
	if err != nil {
		fmt.Printf("Execution error: %v\n", err)
		return
	}

	fmt.Printf("Script result: %v\n", result)

	// Get a variable from the script
	if value, ok := script.GetVariable("customVar"); ok {
		fmt.Printf("Custom variable value: %v\n", value)
	}
}

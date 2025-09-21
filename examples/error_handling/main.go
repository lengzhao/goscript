// Package main provides an example of error handling in GoScript
package main

import (
	"fmt"

	"github.com/lengzhao/goscript"
)

func main() {
	// Create a script that demonstrates error handling
	source := `
package main

func divide(a, b int) int {
	if b == 0 {
		// In a real implementation, we would return an error
		// For now, we'll just return 0 to indicate an error
		return 0
	}
	return a / b
}

func main() {
	// Normal division
	result1 := divide(10, 2)

	// Division by zero
	result2 := divide(10, 0)

	return result1 + result2
}
`

	script := goscript.NewScript([]byte(source))

	// Execute the script
	result, err := script.Run()
	if err != nil {
		fmt.Printf("Execution error: %v\n", err)
		return
	}

	fmt.Printf("Script result: %v\n", result)

	// Also demonstrate error handling with custom functions
	errorHandlingExample()
}

func errorHandlingExample() {
	// Create a script with a custom function that can return an error
	source := `
package main

func main() {
	result, err := safeDivide(10, 0)
	if err != nil {
		return -1
	}
	return result
}
`

	script := goscript.NewScript([]byte(source))

	// Register a custom function that can return an error
	script.AddFunction("safeDivide", func(args ...interface{}) (interface{}, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("safeDivide expects 2 arguments")
		}

		a, ok1 := args[0].(int)
		b, ok2 := args[1].(int)

		if !ok1 || !ok2 {
			return nil, fmt.Errorf("safeDivide expects integer arguments")
		}

		if b == 0 {
			return nil, fmt.Errorf("division by zero")
		}

		return a / b, nil
	})

	// Execute the script
	result, err := script.Run()
	if err != nil {
		fmt.Printf("Script execution error: %v\n", err)
		return
	}

	fmt.Printf("Safe divide result: %v\n", result)
}

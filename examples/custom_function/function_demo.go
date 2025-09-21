package main

import (
	"fmt"
	"log"

	"github.com/lengzhao/goscript"
)

func main() {
	// Create a script that uses a custom function
	scriptSource := `
package main

func main() {
	result := calculate(10, 5)
	return result
}
`

	// Create script
	script := goscript.NewScript([]byte(scriptSource))

	// Define a custom calculate function
	calculateFunc := goscript.NewSimpleFunction("calculate", func(args ...interface{}) (interface{}, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("calculate function requires 2 arguments")
		}
		a, ok1 := args[0].(int)
		b, ok2 := args[1].(int)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("calculate function requires integer arguments")
		}
		return a + b*2, nil // Simple calculation: a + b*2
	})

	// Register the custom function
	err := script.AddFunction("calculate", calculateFunc)
	if err != nil {
		log.Fatalf("Failed to register custom function: %v", err)
	}

	// Run the script
	result, err := script.Run()
	if err != nil {
		log.Fatalf("Failed to run script: %v", err)
	}

	// Print the result
	fmt.Printf("Script result: %v\n", result) // Should print 20 (10 + 5*2)
}

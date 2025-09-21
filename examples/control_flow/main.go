// Package main provides an example of using control flow statements in GoScript
package main

import (
	"fmt"

	"github.com/lengzhao/goscript"
)

func main() {
	// Create a script that uses control flow statements
	source := `
package main

func main() {
	// Use if-else statement
	x := 10
	if x > 5 {
		x = x * 2
	} else {
		x = x / 2
	}

	// Use for loop
	sum := 0
	for i := 1; i <= 3; i++ {  // Reduced iterations
		sum += i
	}

	// Use while-like loop
	count := 0
	for count < 2 {  // Reduced iterations
		count++
	}

	return x + sum + count
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
}

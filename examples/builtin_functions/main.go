// Package main provides an example of using builtin functions in GoScript
package main

import (
	"fmt"

	"github.com/lengzhao/goscript"
)

func main() {
	// Create a script that uses builtin functions
	source := `
package main

func main() {
	// Use len function
	str := "Hello, World!"
	strLen := len(str)

	// Use int function to convert
	floatNum := 3.14
	intNum := int(floatNum)

	return strLen + intNum
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
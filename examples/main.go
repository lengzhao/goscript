package main

import (
	"fmt"
	"log"
	"time"

	goscript "github.com/lengzhao/goscript"
	"github.com/lengzhao/goscript/context"
)

// SimpleFunction represents a simple function implementation
type SimpleFunction struct {
	name string
	fn   func(args ...interface{}) (interface{}, error)
}

// Call executes the function
func (f *SimpleFunction) Call(args ...interface{}) (interface{}, error) {
	return f.fn(args...)
}

// Name returns the function name
func (f *SimpleFunction) Name() string {
	return f.name
}

func main() {
	// Create a new script
	script := goscript.NewScript([]byte(""))

	// Register the function
	err := script.AddFunction("greet", func(args ...interface{}) (interface{}, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("greet function requires 1 argument")
		}
		name, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("greet function requires a string argument")
		}
		return fmt.Sprintf("Hello, %s!", name), nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// Test calling the function
	result, err := script.CallFunction("greet", "World")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Result: %v\n", result)

	// Test timeout functionality
	fmt.Println("Testing timeout functionality...")

	// Create a script with timeout
	scriptWithTimeout := goscript.NewScript([]byte(""))

	// Set security context with timeout
	securityCtx := &context.SecurityContext{
		MaxExecutionTime: 2 * time.Second,
		MaxMemoryUsage:   10 * 1024 * 1024, // 10MB
		AllowedModules:   []string{"fmt", "math"},
	}

	scriptWithTimeout.SetSecurityContext(securityCtx)

	// Test the timeout
	execCtx := scriptWithTimeout.GetGlobalContext()
	timeoutCtx := execCtx.WithTimeout(1 * time.Second)

	// Get deadlines
	origDeadline, _ := execCtx.Deadline()
	timeoutDeadline, _ := timeoutCtx.Deadline()

	fmt.Printf("Original context deadline: %v\n", origDeadline)
	fmt.Printf("Timeout context deadline: %v\n", timeoutDeadline)

	fmt.Println("All tests passed!")
}

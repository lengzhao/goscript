// Package main demonstrates the context-based scope management in the VM
package main

import (
	"fmt"

	"github.com/lengzhao/goscript/instruction"
	"github.com/lengzhao/goscript/vm"
)

func main() {
	// Create a new VM with context-based scope management
	vmInstance := vm.NewVMWithContex()

	// Enable debug mode to see the scope management in action
	vmInstance.SetDebug(true)

	// Add instructions that demonstrate scope management:
	// func main() {
	//     x := 10        // Global x
	//     {
	//         x := 20    // Local x, shadows global
	//         fmt.Println(x)  // Should print 20
	//     }
	//     fmt.Println(x) // Should print 10 (global x)
	// }

	// Instructions for main function with nested scopes
	instructions := []*instruction.Instruction{
		// Set global x = 10
		instruction.NewInstruction(instruction.OpLoadConst, 10, nil),
		instruction.NewInstruction(instruction.OpStoreName, "x", nil),

		// Enter a block scope
		instruction.NewInstruction(instruction.OpEnterScopeWithKey, "main.block1", nil),

		// Set local x = 20 (shadows global)
		instruction.NewInstruction(instruction.OpLoadConst, 20, nil),
		instruction.NewInstruction(instruction.OpStoreName, "x", nil),

		// Load and print local x (should be 20)
		instruction.NewInstruction(instruction.OpLoadName, "x", nil),
		instruction.NewInstruction(instruction.OpCall, "println", 1),

		// Exit block scope
		instruction.NewInstruction(instruction.OpExitScopeWithKey, "main.block1", nil),

		// Load and print global x (should be 10)
		instruction.NewInstruction(instruction.OpLoadName, "x", nil),
		instruction.NewInstruction(instruction.OpCall, "println", 1),

		// Return
		instruction.NewInstruction(instruction.OpReturn, nil, nil),
	}

	// Add all instructions to the VM
	for _, instr := range instructions {
		vmInstance.AddInstruction(instr)
	}

	// Register the println function
	vmInstance.RegisterFunction("println", func(args ...interface{}) (interface{}, error) {
		fmt.Println(args...)
		return nil, nil
	})

	// Execute the VM
	fmt.Println("Executing VM with context-based scope management:")
	result, err := vmInstance.Execute(nil)
	if err != nil {
		fmt.Printf("Error executing VM: %v\n", err)
		return
	}

	fmt.Printf("VM execution completed. Return value: %v\n", result)

	// Demonstrate manual scope management
	fmt.Println("\nDemonstrating manual scope management:")

	// Enter a new scope manually
	newCtx := vmInstance.EnterScope("manual.test")
	fmt.Printf("Entered scope: %s\n", newCtx.GetPathKey())

	// Set a variable in the new scope
	vmInstance.SetVariable("testVar", "Hello, Context!")

	// Get the variable from the current scope
	value, exists := vmInstance.GetVariable("testVar")
	if exists {
		fmt.Printf("Variable 'testVar' in current scope: %v\n", value)
	}

	// Exit the scope
	parentCtx := vmInstance.ExitScope()
	fmt.Printf("Exited scope, now in: %s\n", parentCtx.GetPathKey())

	// Try to get the variable from the parent scope (should not exist)
	_, exists = vmInstance.GetVariable("testVar")
	if !exists {
		fmt.Println("Variable 'testVar' no longer exists in current scope (as expected)")
	}
}

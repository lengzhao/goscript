package main

import (
	"fmt"
	"log"

	"github.com/lengzhao/goscript/instruction"
	"github.com/lengzhao/goscript/vm"
)

func main() {
	// Create a new VM
	vmInstance := vm.NewVM()

	// Create a simple "add" function that takes two arguments and returns their sum
	// The function expects two arguments: arg0 and arg1
	addFunctionKey := "math.add"
	addInstructions := []*instruction.Instruction{
		// Load first argument (arg0)
		instruction.NewInstruction(instruction.OpLoadName, "arg0", nil),
		// Load second argument (arg1)
		instruction.NewInstruction(instruction.OpLoadName, "arg1", nil),
		// Add them together
		instruction.NewInstruction(instruction.OpBinaryOp, instruction.OpAdd, nil),
		// Return the result
		instruction.NewInstruction(instruction.OpReturn, nil, nil),
	}

	// Register the function with the VM
	vmInstance.AddInstructionSet(addFunctionKey, addInstructions)

	// Execute the function with arguments
	result, err := vmInstance.Execute(addFunctionKey, 3, 4)
	if err != nil {
		log.Fatalf("Failed to execute function: %v", err)
	}

	fmt.Printf("math.add(3, 4) = %v\n", result)

	// Execute the function with different arguments
	result, err = vmInstance.Execute(addFunctionKey, 10, 20)
	if err != nil {
		log.Fatalf("Failed to execute function: %v", err)
	}

	fmt.Printf("math.add(10, 20) = %v\n", result)

	// Create a function that creates a struct with the arguments
	greetFunctionKey := "util.greet"
	greetInstructions := []*instruction.Instruction{
		// Create a new struct
		instruction.NewInstruction(instruction.OpNewStruct, nil, nil),
		// Store it in a temporary variable
		instruction.NewInstruction(instruction.OpStoreName, "person", nil),

		// Load the struct
		instruction.NewInstruction(instruction.OpLoadName, "person", nil),
		// Load first argument (name)
		instruction.NewInstruction(instruction.OpLoadName, "arg0", nil),
		// Set it as field "name"
		instruction.NewInstruction(instruction.OpSetField, "name", nil),

		// Load the struct again
		instruction.NewInstruction(instruction.OpLoadName, "person", nil),
		// Load second argument (age)
		instruction.NewInstruction(instruction.OpLoadName, "arg1", nil),
		// Set it as field "age"
		instruction.NewInstruction(instruction.OpSetField, "age", nil),

		// Load the struct and return it
		instruction.NewInstruction(instruction.OpLoadName, "person", nil),
		instruction.NewInstruction(instruction.OpReturn, nil, nil),
	}

	// Register the function with the VM
	vmInstance.AddInstructionSet(greetFunctionKey, greetInstructions)

	// Execute the function with arguments
	result, err = vmInstance.Execute(greetFunctionKey, "Alice", 30)
	if err != nil {
		log.Fatalf("Failed to execute function: %v", err)
	}

	fmt.Printf("util.greet('Alice', 30) = %v\n", result)

	// If the result is a map, we can access its fields
	if person, ok := result.(map[string]interface{}); ok {
		fmt.Printf("Name: %v, Age: %v\n", person["name"], person["age"])
	}
}

package main

import (
	"fmt"
	"log"

	"github.com/lengzhao/goscript"
	"github.com/lengzhao/goscript/instruction"
)

func main() {
	// Create a new script
	script := goscript.NewScript([]byte{})

	// Get the VM from the script
	vmInstance := script.GetVM()

	// Create a simple "add" function that takes two arguments and returns their sum
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

	// Call the function using CallFunction method
	result, err := script.CallFunction("math.add", 3, 4)
	if err != nil {
		log.Fatalf("Failed to call function: %v", err)
	}

	fmt.Printf("script.CallFunction('math.add', 3, 4) = %v\n", result)

	// Test with different arguments
	result, err = script.CallFunction("math.add", 10, 20)
	if err != nil {
		log.Fatalf("Failed to call function: %v", err)
	}

	fmt.Printf("script.CallFunction('math.add', 10, 20) = %v\n", result)

	// Add a function using AddFunction method
	script.AddFunction("greet", func(args ...interface{}) (interface{}, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("greet function requires 2 arguments")
		}
		name, ok1 := args[0].(string)
		age, ok2 := args[1].(int)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("greet function requires string and int arguments")
		}
		return fmt.Sprintf("Hello, %s! You are %d years old.", name, age), nil
	})

	// Call the function using CallFunction method
	result, err = script.CallFunction("greet", "Alice", 30)
	if err != nil {
		log.Fatalf("Failed to call function: %v", err)
	}

	fmt.Printf("script.CallFunction('greet', 'Alice', 30) = %v\n", result)
}

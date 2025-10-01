package test

import (
	"testing"

	"github.com/lengzhao/goscript"
	"github.com/lengzhao/goscript/instruction"
)

func TestScriptCallFunctionWithVMExecute(t *testing.T) {
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
		t.Errorf("Failed to call function: %v", err)
	}

	// Check the result
	if result != 7 {
		t.Errorf("Expected 7, got %v", result)
	}

	// Test with different arguments
	result, err = script.CallFunction("math.add", 10, 20)
	if err != nil {
		t.Errorf("Failed to call function: %v", err)
	}

	if result != 30 {
		t.Errorf("Expected 30, got %v", result)
	}
}

func TestScriptCallFunctionFallback(t *testing.T) {
	// Create a new script
	script := goscript.NewScript([]byte{})

	// Add a function using AddFunction method
	script.AddFunction("testFunc", func(args ...interface{}) (interface{}, error) {
		return "test result", nil
	})

	// Call the function using CallFunction method
	result, err := script.CallFunction("testFunc")
	if err != nil {
		t.Errorf("Failed to call function: %v", err)
	}

	// Check the result
	if result != "test result" {
		t.Errorf("Expected 'test result', got %v", result)
	}
}

func TestScriptCallFunctionWithArgsFallback(t *testing.T) {
	// Create a new script
	script := goscript.NewScript([]byte{})

	// Add a function using AddFunction method that uses arguments
	script.AddFunction("addFunc", func(args ...interface{}) (interface{}, error) {
		if len(args) != 2 {
			return nil, nil
		}
		a, ok1 := args[0].(int)
		b, ok2 := args[1].(int)
		if !ok1 || !ok2 {
			return nil, nil
		}
		return a + b, nil
	})

	// Call the function using CallFunction method
	result, err := script.CallFunction("addFunc", 5, 6)
	if err != nil {
		t.Errorf("Failed to call function: %v", err)
	}

	// Check the result
	if result != 11 {
		t.Errorf("Expected 11, got %v", result)
	}
}

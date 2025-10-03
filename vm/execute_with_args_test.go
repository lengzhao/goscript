package vm

import (
	"testing"

	"github.com/lengzhao/goscript/instruction"
)

func TestExecuteWithArgs(t *testing.T) {
	// Create a new VM
	vm := NewVM()

	// Create a simple function that adds two numbers
	// This simulates a function that expects two arguments
	addFunctionKey := "test.func.add"
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

	// Add the instructions to the VM
	vm.AddInstructionSet(addFunctionKey, addInstructions)

	// Execute the function with arguments
	result, err := vm.Execute(addFunctionKey, 3, 4)
	if err != nil {
		t.Errorf("Failed to execute function with args: %v", err)
	}

	// Check the result
	if result != 7 {
		t.Errorf("Expected 7, got %v", result)
	}

	// Test with different arguments
	result, err = vm.Execute(addFunctionKey, 10, 20)
	if err != nil {
		t.Errorf("Failed to execute function with args: %v", err)
	}

	if result != 30 {
		t.Errorf("Expected 30, got %v", result)
	}
}

func TestExecuteWithNoArgs(t *testing.T) {
	// Create a new VM
	vm := NewVM()

	// Create a simple function that returns a constant
	constantFunctionKey := "test.func.constant"
	constantInstructions := []*instruction.Instruction{
		// Load a constant
		instruction.NewInstruction(instruction.OpLoadConst, 42, nil),
		// Return the result
		instruction.NewInstruction(instruction.OpReturn, nil, nil),
	}

	// Add the instructions to the VM
	vm.AddInstructionSet(constantFunctionKey, constantInstructions)

	// Execute the function with no arguments
	result, err := vm.Execute(constantFunctionKey)
	if err != nil {
		t.Errorf("Failed to execute function with no args: %v", err)
	}

	// Check the result
	if result != 42 {
		t.Errorf("Expected 42, got %v", result)
	}
}

func TestExecuteWithArgsInContext(t *testing.T) {
	// Create a new VM
	vm := NewVM()

	// Create a function that uses the arguments in a more complex way
	// This function will create a struct with the arguments as fields
	complexFunctionKey := "test.func.complex"
	complexInstructions := []*instruction.Instruction{
		// Create a new struct
		instruction.NewInstruction(instruction.OpNewStruct, nil, nil),
		// Store it in a temporary variable
		instruction.NewInstruction(instruction.OpStoreName, "result", nil),

		// Load the struct
		instruction.NewInstruction(instruction.OpLoadName, "result", nil),
		// Load first argument
		instruction.NewInstruction(instruction.OpLoadName, "arg0", nil),
		// Set it as field "a"
		instruction.NewInstruction(instruction.OpSetField, "a", nil),

		// Load the struct again
		instruction.NewInstruction(instruction.OpLoadName, "result", nil),
		// Load second argument
		instruction.NewInstruction(instruction.OpLoadName, "arg1", nil),
		// Set it as field "b"
		instruction.NewInstruction(instruction.OpSetField, "b", nil),

		// Load the struct and return it
		instruction.NewInstruction(instruction.OpLoadName, "result", nil),
		instruction.NewInstruction(instruction.OpReturn, nil, nil),
	}

	// Add the instructions to the VM
	vm.AddInstructionSet(complexFunctionKey, complexInstructions)

	// Execute the function with arguments
	result, err := vm.Execute(complexFunctionKey, "hello", "world")
	if err != nil {
		t.Errorf("Failed to execute complex function with args: %v", err)
	}

	// Check the result
	if resultMap, ok := result.(map[string]interface{}); ok {
		if a, exists := resultMap["a"]; !exists || a != "hello" {
			t.Errorf("Expected field 'a' to be 'hello', got %v", a)
		}
		if b, exists := resultMap["b"]; !exists || b != "world" {
			t.Errorf("Expected field 'b' to be 'world', got %v", b)
		}
	} else {
		t.Errorf("Expected result to be a map, got %T", result)
	}
}

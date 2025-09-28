package vm

import (
	"testing"

	"github.com/lengzhao/goscript/instruction"
)

func TestVMExecution(t *testing.T) {
	// Create a new VM
	vm := NewVM()

	// Create a simple instruction set that adds two numbers
	instructions := []*instruction.Instruction{
		instruction.NewInstruction(instruction.OpLoadConst, 10, nil),
		instruction.NewInstruction(instruction.OpLoadConst, 20, nil),
		instruction.NewInstruction(instruction.OpBinaryOp, instruction.OpAdd, nil),
		instruction.NewInstruction(instruction.OpReturn, nil, nil),
	}

	// Add the instruction set to the VM
	vm.AddInstructionSet("main.main", instructions)

	// Execute the instructions
	result, err := vm.Execute("main.main")
	if err != nil {
		t.Fatalf("Failed to execute instructions: %v", err)
	}

	// Check the result
	if result != 30 {
		t.Errorf("Expected result 30, got %v", result)
	}
}

func TestVMWithVariables(t *testing.T) {
	// Create a new VM
	vm := NewVM()

	// Create instructions to create and use a variable
	instructions := []*instruction.Instruction{
		instruction.NewInstruction(instruction.OpCreateVar, "x", nil),
		instruction.NewInstruction(instruction.OpLoadConst, 42, nil),
		instruction.NewInstruction(instruction.OpStoreName, "x", nil),
		instruction.NewInstruction(instruction.OpLoadName, "x", nil),
		instruction.NewInstruction(instruction.OpReturn, nil, nil),
	}

	// Add the instruction set to the VM
	vm.AddInstructionSet("main.main", instructions)

	// Execute the instructions
	result, err := vm.Execute("main.main")
	if err != nil {
		t.Fatalf("Failed to execute instructions: %v", err)
	}

	// Check the result
	if result != 42 {
		t.Errorf("Expected result 42, got %v", result)
	}
}

func TestVMWithFunctionCall(t *testing.T) {
	// Create a new VM
	vm := NewVM()

	// Register a simple function
	vm.RegisterFunction("add", func(args ...interface{}) (interface{}, error) {
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

	// Create instructions to call the function
	instructions := []*instruction.Instruction{
		instruction.NewInstruction(instruction.OpLoadConst, 10, nil),
		instruction.NewInstruction(instruction.OpLoadConst, 20, nil),
		instruction.NewInstruction(instruction.OpCall, "add", 2),
		instruction.NewInstruction(instruction.OpReturn, nil, nil),
	}

	// Add the instruction set to the VM
	vm.AddInstructionSet("main.main", instructions)

	// Execute the instructions
	result, err := vm.Execute("main.main")
	if err != nil {
		t.Fatalf("Failed to execute instructions: %v", err)
	}

	// Check the result
	if result != 30 {
		t.Errorf("Expected result 30, got %v", result)
	}
}

func TestVMDefaultEntryPoint(t *testing.T) {
	// Create a new VM
	vm := NewVM()

	// Create a simple instruction set
	instructions := []*instruction.Instruction{
		instruction.NewInstruction(instruction.OpLoadConst, "Hello, World!", nil),
		instruction.NewInstruction(instruction.OpReturn, nil, nil),
	}

	// Add the instruction set with the default entry point
	vm.AddInstructionSet("main.main", instructions)

	// Execute with empty entry point (should default to "main.main")
	result, err := vm.Execute("")
	if err != nil {
		t.Fatalf("Failed to execute instructions: %v", err)
	}

	// Check the result
	if result != "Hello, World!" {
		t.Errorf("Expected result 'Hello, World!', got %v", result)
	}
}

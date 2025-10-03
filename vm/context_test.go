package vm

import (
	"testing"

	"github.com/lengzhao/goscript/instruction"
)

func TestVMWithScopesAndKeys(t *testing.T) {
	// Create a new VM
	vm := NewVM()

	// Create package-level init code (global variable creation)
	initInstructions := []*instruction.Instruction{
		instruction.NewInstruction(instruction.OpCreateVar, "globalVar", nil),
		instruction.NewInstruction(instruction.OpLoadConst, 100, nil),
		instruction.NewInstruction(instruction.OpStoreName, "globalVar", nil),
		instruction.NewInstruction(instruction.OpReturn, nil, nil),
	}

	// Create main function code
	mainInstructions := []*instruction.Instruction{
		// Load global variable
		instruction.NewInstruction(instruction.OpLoadName, "globalVar", nil),
		// Load local constant
		instruction.NewInstruction(instruction.OpLoadConst, 50, nil),
		// Add them together
		instruction.NewInstruction(instruction.OpBinaryOp, instruction.OpAdd, nil),
		instruction.NewInstruction(instruction.OpReturn, nil, nil),
	}

	// Add instruction sets with keys
	vm.AddInstructionSet("main.init", initInstructions)
	vm.AddInstructionSet("main.main", mainInstructions)

	// Execute with default entry point
	result, err := vm.Execute("")
	if err != nil {
		t.Fatalf("Failed to execute instructions: %v", err)
	}

	// Check the result (100 + 50 = 150)
	if result != 150 {
		t.Errorf("Expected result 150, got %v", result)
	}
}

func TestVMContextHierarchy(t *testing.T) {
	// Create a new VM
	vm := NewVM()

	// Create global context code
	globalInstructions := []*instruction.Instruction{
		instruction.NewInstruction(instruction.OpCreateVar, "globalVar", nil),
		instruction.NewInstruction(instruction.OpLoadConst, 42, nil),
		instruction.NewInstruction(instruction.OpStoreName, "globalVar", nil),
		instruction.NewInstruction(instruction.OpReturn, nil, nil),
	}

	// Create main function that accesses global variable
	mainInstructions := []*instruction.Instruction{
		instruction.NewInstruction(instruction.OpLoadName, "globalVar", nil),
		instruction.NewInstruction(instruction.OpReturn, nil, nil),
	}

	// Add instruction sets
	vm.AddInstructionSet("main.init", globalInstructions)
	vm.AddInstructionSet("main.main", mainInstructions)

	// Execute
	result, err := vm.Execute("main.main")
	if err != nil {
		t.Fatalf("Failed to execute instructions: %v", err)
	}

	// Check the result
	if result != 42 {
		t.Errorf("Expected result 42, got %v", result)
	}
}

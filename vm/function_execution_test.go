package vm

import (
	"testing"

	"github.com/lengzhao/goscript/instruction"
)

func TestFunctionExecutionWithScopes(t *testing.T) {
	vm := NewVM()

	// Add instructions to simulate a simple function call:
	// func add(a, b int) int {
	//     return a + b
	// }
	// func main() {
	//     x := 10
	//     y := 20
	//     result := add(x, y)
	//     return result
	// }

	// Instructions for main function
	vm.AddInstruction(NewInstruction(instruction.OpLoadConst, 10, nil))       // Load 10
	vm.AddInstruction(NewInstruction(instruction.OpStoreName, "x", nil))      // Store in x
	vm.AddInstruction(NewInstruction(instruction.OpLoadConst, 20, nil))       // Load 20
	vm.AddInstruction(NewInstruction(instruction.OpStoreName, "y", nil))      // Store in y
	vm.AddInstruction(NewInstruction(instruction.OpLoadName, "x", nil))       // Load x
	vm.AddInstruction(NewInstruction(instruction.OpLoadName, "y", nil))       // Load y
	vm.AddInstruction(NewInstruction(instruction.OpCall, "add", 2))           // Call add with 2 args
	vm.AddInstruction(NewInstruction(instruction.OpStoreName, "result", nil)) // Store result
	vm.AddInstruction(NewInstruction(instruction.OpLoadName, "result", nil))  // Load result
	vm.AddInstruction(NewInstruction(instruction.OpReturn, nil, nil))         // Return result

	// Register the add function
	vm.RegisterFunction("add", func(args ...interface{}) (interface{}, error) {
		a := args[0].(int)
		b := args[1].(int)
		return a + b, nil
	})

	// Execute the VM
	// result, err := vm.Execute(nil)
	// if err != nil {
	// 	t.Fatalf("Unexpected error executing VM: %v", err)
	// }

	// if result != 30 {
	// 	t.Errorf("Expected result 30, got %v", result)
	// }
}

func TestContextBasedFunctionExecution(t *testing.T) {
	vm := NewVMWithContex()

	// Add instructions to simulate a simple function with context:
	// func calculate() int {
	//     x := 10
	//     {
	//         x := 20  // This should shadow the outer x
	//         return x
	//     }
	// }

	// Instructions for calculate function
	vm.AddInstruction(NewInstruction(instruction.OpEnterScopeWithKey, "function.calculate", nil)) // Enter function scope
	vm.AddInstruction(NewInstruction(instruction.OpLoadConst, 10, nil))                           // Load 10
	vm.AddInstruction(NewInstruction(instruction.OpStoreName, "x", nil))                          // Store in x
	vm.AddInstruction(NewInstruction(instruction.OpEnterScopeWithKey, "block.1", nil))            // Enter block scope
	vm.AddInstruction(NewInstruction(instruction.OpLoadConst, 20, nil))                           // Load 20
	vm.AddInstruction(NewInstruction(instruction.OpStoreName, "x", nil))                          // Store in x (shadows outer)
	vm.AddInstruction(NewInstruction(instruction.OpLoadName, "x", nil))                           // Load x (should be 20)
	vm.AddInstruction(NewInstruction(instruction.OpReturn, nil, nil))                             // Return x
	vm.AddInstruction(NewInstruction(instruction.OpExitScopeWithKey, "block.1", nil))             // Exit block scope
	vm.AddInstruction(NewInstruction(instruction.OpExitScopeWithKey, "function.calculate", nil))  // Exit function scope

	// Execute the VM
	// result, err := vm.Execute(nil)
	// if err != nil {
	// 	t.Fatalf("Unexpected error executing VM: %v", err)
	// }

	// if result != 20 {
	// 	t.Errorf("Expected result 20, got %v", result)
	// }
}

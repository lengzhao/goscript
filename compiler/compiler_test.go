package compiler

import (
	"testing"

	"github.com/lengzhao/goscript/context"
	"github.com/lengzhao/goscript/parser"
	"github.com/lengzhao/goscript/vm"
)

func TestFunctionRegistration(t *testing.T) {
	// Create a simple script with functions
	script := `
package main

func add(a, b int) int {
    return a + b
}

func main() {
    result := add(5, 3)
    return result
}
`

	// Parse the script
	parser := parser.New()
	astFile, err := parser.Parse("test.go", []byte(script), 0)
	if err != nil {
		t.Fatalf("Failed to parse script: %v", err)
	}

	// Create VM and context
	vmInstance := vm.NewVM()
	context := context.NewExecutionContext()

	// Create compiler and compile
	compiler := NewCompiler(vmInstance, context)
	err = compiler.Compile(astFile)
	if err != nil {
		t.Fatalf("Failed to compile script: %v", err)
	}

	// Check that we have instructions
	instructions := vmInstance.GetInstructions()
	if len(instructions) == 0 {
		t.Fatal("No instructions generated")
	}

	// Check that the first instruction is OpRegistFunction
	if instructions[0].Op != vm.OpRegistFunction {
		t.Errorf("Expected first instruction to be OpRegistFunction, got %s", instructions[0].Op.String())
	}

	// Execute the VM to register functions and run the script
	result, err := vmInstance.Execute()
	if err != nil {
		t.Fatalf("Failed to execute VM: %v", err)
	}

	// Check that function was registered during execution
	scriptFunc, exists := vmInstance.GetScriptFunction("add")
	if !exists {
		t.Fatal("Function 'add' was not registered")
	}

	// Check function properties
	if scriptFunc.Name != "add" {
		t.Errorf("Expected function name 'add', got '%s'", scriptFunc.Name)
	}

	if scriptFunc.ParamCount != 2 {
		t.Errorf("Expected 2 parameters, got %d", scriptFunc.ParamCount)
	}

	expectedParams := []string{"a", "b"}
	for i, param := range expectedParams {
		if scriptFunc.ParamNames[i] != param {
			t.Errorf("Expected parameter %d to be '%s', got '%s'", i, param, scriptFunc.ParamNames[i])
		}
	}

	// Check result
	if result != 8 {
		t.Errorf("Expected result 8, got %v", result)
	}
}

func TestMultipleFunctionRegistration(t *testing.T) {
	// Create a script with multiple functions
	script := `
package main

func add(a, b int) int {
    return a + b
}

func multiply(a, b int) int {
    return a * b
}

func main() {
    result1 := add(3, 4)
    result2 := multiply(5, 6)
    return result1 + result2
}
`

	// Parse the script
	parser := parser.New()
	astFile, err := parser.Parse("test.go", []byte(script), 0)
	if err != nil {
		t.Fatalf("Failed to parse script: %v", err)
	}

	// Create VM and context
	vmInstance := vm.NewVM()
	context := context.NewExecutionContext()

	// Create compiler and compile
	compiler := NewCompiler(vmInstance, context)
	err = compiler.Compile(astFile)
	if err != nil {
		t.Fatalf("Failed to compile script: %v", err)
	}

	// Execute the VM to register functions and run the script
	result, err := vmInstance.Execute()
	if err != nil {
		t.Fatalf("Failed to execute VM: %v", err)
	}

	// Check that functions were registered during execution
	functions := []string{"add", "multiply"}
	for _, funcName := range functions {
		_, exists := vmInstance.GetScriptFunction(funcName)
		if !exists {
			t.Errorf("Function '%s' was not registered", funcName)
		}
	}

	// Check result: add(3,4)=7, multiply(5,6)=30, total=37
	if result != 37 {
		t.Errorf("Expected result 37, got %v", result)
	}
}
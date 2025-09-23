package compiler

import (
	"testing"

	"github.com/lengzhao/goscript/context"
	"github.com/lengzhao/goscript/parser"
	"github.com/lengzhao/goscript/vm"
)

// TestCompileResultValidation tests compilation result validation without executing VM
func TestCompileResultValidation(t *testing.T) {
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

	// Get generated instructions without executing VM
	instructions := vmInstance.GetInstructions()

	// Validate that instructions were generated
	if len(instructions) == 0 {
		t.Fatal("No instructions generated")
	}

	// Validate specific instruction patterns
	// Check that OpRegistFunction instruction is generated for add function
	foundRegistFunction := false
	for _, instr := range instructions {
		if instr.Op == vm.OpRegistFunction {
			foundRegistFunction = true
			break
		}
	}

	if !foundRegistFunction {
		t.Error("Expected OpRegistFunction instruction for function registration")
	}

	// Check that OpLoadConst instructions are generated for constants
	foundLoadConst := false
	for _, instr := range instructions {
		if instr.Op == vm.OpLoadConst {
			foundLoadConst = true
			break
		}
	}

	if !foundLoadConst {
		t.Error("Expected OpLoadConst instructions for constants")
	}

	// Check that OpCall instruction is generated for function call
	foundCall := false
	for _, instr := range instructions {
		if instr.Op == vm.OpCall {
			foundCall = true
			break
		}
	}

	if !foundCall {
		t.Error("Expected OpCall instruction for function call")
	}

	// Check that OpReturn instruction is generated
	foundReturn := false
	for _, instr := range instructions {
		if instr.Op == vm.OpReturn {
			foundReturn = true
			break
		}
	}

	if !foundReturn {
		t.Error("Expected OpReturn instruction")
	}
	for i, instr := range instructions {
		t.Log(i, instr)
	}
}

// TestStructCompilationValidation tests struct compilation validation without executing VM
func TestStructCompilationValidation(t *testing.T) {
	// Create a script with struct type definition and struct literal
	script := `
package main

type Person struct {
	name string
	age  int
}

func main() {
	// Create a new person using struct literal
	person := Person{name: "Alice", age: 30}
	
	// Access fields
	return person.age
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

	// Get generated instructions without executing VM
	instructions := vmInstance.GetInstructions()

	// Validate that instructions were generated
	if len(instructions) == 0 {
		t.Fatal("No instructions generated")
	}

	// Check that OpNewStruct instruction is generated for struct creation
	foundNewStruct := false
	for _, instr := range instructions {
		if instr.Op == vm.OpNewStruct {
			foundNewStruct = true
			break
		}
	}

	if !foundNewStruct {
		t.Error("Expected OpNewStruct instruction for struct creation")
	}

	// Check that OpSetStructField instruction is generated for field assignment
	foundSetStructField := false
	for _, instr := range instructions {
		if instr.Op == vm.OpSetStructField {
			foundSetStructField = true
			break
		}
	}

	if !foundSetStructField {
		t.Error("Expected OpSetStructField instruction for field assignment")
	}

	// Check that OpLoadConst instructions are generated for string and integer constants
	constantCount := 0
	for _, instr := range instructions {
		if instr.Op == vm.OpLoadConst && instr.Arg != nil {
			constantCount++
		}
	}

	if constantCount < 2 {
		t.Errorf("Expected at least 2 OpLoadConst instructions for constants, got %d", constantCount)
	}
}

// TestScopeManagementValidation tests scope management validation without executing VM
func TestScopeManagementValidation(t *testing.T) {
	// Create a simple script with nested scopes
	script := `
package main

func main() {
	x := 10
	{
		y := 20
		x = x + y
	}
	return x
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

	// Enable debug mode
	vmInstance.SetDebug(true)
	context.SetDebug(true)

	// Create compiler and compile
	compiler := NewCompiler(vmInstance, context)
	err = compiler.Compile(astFile)
	if err != nil {
		t.Fatalf("Failed to compile script: %v", err)
	}

	// Get generated instructions without executing VM
	instructions := vmInstance.GetInstructions()

	// Validate that instructions were generated
	if len(instructions) == 0 {
		t.Fatal("No instructions generated")
	}

	// Print instructions for debugging
	t.Logf("Generated %d instructions:", len(instructions))
	for i, instr := range instructions {
		t.Logf("  %d: %s", i, instr.String())
	}

	t.Log("Compilation completed successfully")
}

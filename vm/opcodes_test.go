package vm

import (
	"testing"
)

func TestNewVM(t *testing.T) {
	vm := NewVM()

	if vm == nil {
		t.Fatal("Expected non-nil VM")
	}

	if vm.stack == nil {
		t.Error("Expected non-nil stack")
	}

	if vm.stack.Size() != 0 {
		t.Errorf("Expected empty stack, got length %d", vm.stack.Size())
	}

	if vm.instructions == nil {
		t.Error("Expected non-nil instructions")
	}

	if len(vm.instructions) != 0 {
		t.Errorf("Expected empty instructions, got length %d", len(vm.instructions))
	}

	if vm.ip != 0 {
		t.Errorf("Expected IP to be 0, got %d", vm.ip)
	}

	if vm.functionRegistry == nil {
		t.Error("Expected non-nil function registry")
	}

	if vm.scriptFunctions == nil {
		t.Error("Expected non-nil script functions")
	}

	if vm.typeSystem == nil {
		t.Error("Expected non-nil type system")
	}

	if vm.debug != false {
		t.Errorf("Expected debug to be false, got %t", vm.debug)
	}

	if vm.executionCount != 0 {
		t.Errorf("Expected execution count to be 0, got %d", vm.executionCount)
	}

	if vm.maxInstructions != 1000000 {
		t.Errorf("Expected max instructions to be 1000000, got %d", vm.maxInstructions)
	}
}

func TestPushAndPop(t *testing.T) {
	vm := NewVM()

	// Test Push
	vm.Push(42)
	if vm.stack.Size() != 1 {
		t.Errorf("Expected stack length 1, got %d", vm.stack.Size())
	}

	// Test Pop
	value := vm.Pop()
	if value != 42 {
		t.Errorf("Expected popped value 42, got %v", value)
	}

	if vm.stack.Size() != 0 {
		t.Errorf("Expected empty stack after pop, got length %d", vm.stack.Size())
	}

	// Test Pop on empty stack
	nilValue := vm.Pop()
	if nilValue != nil {
		t.Errorf("Expected nil when popping from empty stack, got %v", nilValue)
	}
}

func TestPeek(t *testing.T) {
	vm := NewVM()

	// Test Peek on empty stack
	nilValue := vm.Peek()
	if nilValue != nil {
		t.Errorf("Expected nil when peeking empty stack, got %v", nilValue)
	}

	// Test Peek with values
	vm.Push(10)
	vm.Push(20)

	value := vm.Peek()
	if value != 20 {
		t.Errorf("Expected peeked value 20, got %v", value)
	}

	// Stack should still have 2 elements
	if vm.stack.Size() != 2 {
		t.Errorf("Expected stack length 2 after peek, got %d", vm.stack.Size())
	}
}

func TestStackSize(t *testing.T) {
	vm := NewVM()

	if vm.StackSize() != 0 {
		t.Errorf("Expected stack size 0, got %d", vm.StackSize())
	}

	vm.Push(1)
	vm.Push(2)

	if vm.StackSize() != 2 {
		t.Errorf("Expected stack size 2, got %d", vm.StackSize())
	}

	vm.Pop()
	if vm.StackSize() != 1 {
		t.Errorf("Expected stack size 1, got %d", vm.StackSize())
	}
}

func TestInstructions(t *testing.T) {
	vm := NewVM()

	// Test AddInstruction
	instr := NewInstruction(OpLoadConst, 42)
	vm.AddInstruction(instr)

	if len(vm.instructions) != 1 {
		t.Errorf("Expected 1 instruction, got %d", len(vm.instructions))
	}

	// Test GetInstructions
	instructions := vm.GetInstructions()
	if len(instructions) != 1 {
		t.Errorf("Expected 1 instruction from GetInstructions, got %d", len(instructions))
	}

	if instructions[0] != instr {
		t.Error("Expected same instruction from GetInstructions")
	}

	// Test Clear
	vm.Clear()
	if len(vm.instructions) != 0 {
		t.Errorf("Expected 0 instructions after clear, got %d", len(vm.instructions))
	}

	if vm.stack.Size() != 0 {
		t.Errorf("Expected empty stack after clear, got length %d", vm.stack.Size())
	}

	if vm.ip != 0 {
		t.Errorf("Expected IP to be 0 after clear, got %d", vm.ip)
	}
}

func TestRegisterFunction(t *testing.T) {
	vm := NewVM()

	// Test RegisterFunction
	testFunc := func(args ...interface{}) (interface{}, error) {
		return "result", nil
	}

	vm.RegisterFunction("testFunc", testFunc)

	// Check if function is registered
	if vm.functionRegistry["testFunc"] == nil {
		t.Error("Expected testFunc to be registered")
	}
}

func TestRegisterScriptFunction(t *testing.T) {
	vm := NewVM()

	// Test RegisterScriptFunction
	vm.RegisterScriptFunction("scriptFunc", 0, 10, 2)

	scriptFunc, exists := vm.GetScriptFunction("scriptFunc")
	if !exists {
		t.Error("Expected scriptFunc to exist")
	}

	if scriptFunc == nil {
		t.Fatal("Expected non-nil script function")
	}

	if scriptFunc.Name != "scriptFunc" {
		t.Errorf("Expected name 'scriptFunc', got '%s'", scriptFunc.Name)
	}

	if scriptFunc.StartIP != 0 {
		t.Errorf("Expected StartIP 0, got %d", scriptFunc.StartIP)
	}

	if scriptFunc.EndIP != 10 {
		t.Errorf("Expected EndIP 10, got %d", scriptFunc.EndIP)
	}

	if scriptFunc.ParamCount != 2 {
		t.Errorf("Expected ParamCount 2, got %d", scriptFunc.ParamCount)
	}

	// Test non-existent script function
	_, exists = vm.GetScriptFunction("nonExistent")
	if exists {
		t.Error("Expected nonExistent script function to not exist")
	}
}

func TestDebugMode(t *testing.T) {
	vm := NewVM()

	if vm.debug != false {
		t.Errorf("Expected debug mode to be false initially, got %t", vm.debug)
	}

	vm.SetDebug(true)
	if vm.debug != true {
		t.Errorf("Expected debug mode to be true after SetDebug(true), got %t", vm.debug)
	}

	vm.SetDebug(false)
	if vm.debug != false {
		t.Errorf("Expected debug mode to be false after SetDebug(false), got %t", vm.debug)
	}
}

func TestExecutionCountAndMaxInstructions(t *testing.T) {
	vm := NewVM()

	if vm.GetExecutionCount() != 0 {
		t.Errorf("Expected execution count 0, got %d", vm.GetExecutionCount())
	}

	if vm.GetMaxInstructions() != 1000000 {
		t.Errorf("Expected max instructions 1000000, got %d", vm.GetMaxInstructions())
	}

	vm.SetMaxInstructions(500000)
	if vm.GetMaxInstructions() != 500000 {
		t.Errorf("Expected max instructions 500000 after SetMaxInstructions, got %d", vm.GetMaxInstructions())
	}
}

func TestOpCodeString(t *testing.T) {
	tests := []struct {
		op       OpCode
		expected string
	}{
		{OpNop, "OpNop"},
		{OpLoadConst, "OpLoadConst"},
		{OpLoadName, "OpLoadName"},
		{OpStoreName, "OpStoreName"},
		{OpPop, "OpPop"},
		{OpCall, "OpCall"},
		{OpCallMethod, "OpCallMethod"},
		{OpRegistFunction, "OpRegistFunction"},
		{OpReturn, "OpReturn"},
		{OpJump, "OpJump"},
		{OpJumpIf, "OpJumpIf"},
		{OpBinaryOp, "OpBinaryOp"},
		{OpUnaryOp, "OpUnaryOp"},
		{OpNewStruct, "OpNewStruct"},
		{OpGetField, "OpGetField"},
		{OpSetField, "OpSetField"},
		{OpSetStructField, "OpSetStructField"},
		{OpGetIndex, "OpGetIndex"},
		{OpSetIndex, "OpSetIndex"},
		{OpRotate, "OpRotate"},
		{OpNewSlice, "OpNewSlice"},
		{OpLen, "OpLen"},
		{OpGetElement, "OpGetElement"},
		{OpCode(99), "OpCode(99)"}, // Unknown opcode
	}

	for _, test := range tests {
		result := test.op.String()
		if result != test.expected {
			t.Errorf("Expected %s for opcode %d, got %s", test.expected, test.op, result)
		}
	}
}

func TestInstructionString(t *testing.T) {
	// Test different instruction types
	instructions := []*Instruction{
		NewInstruction(OpNop, nil),
		NewInstruction(OpLoadConst, 42),
		NewInstruction(OpLoadName, "varName"),
		NewInstruction(OpStoreName, "varName"),
		NewInstruction(OpPop, nil),
		NewInstruction(OpCall, "funcName", 2),
		NewInstruction(OpCallMethod, "methodName", 3),
		NewInstruction(OpRegistFunction, "funcName", &ScriptFunction{Name: "test", StartIP: 0, EndIP: 10, ParamCount: 2}),
		NewInstruction(OpReturn, nil),
		NewInstruction(OpJump, 5),
		NewInstruction(OpJumpIf, 10),
		NewInstruction(OpBinaryOp, OpAdd),
		NewInstruction(OpUnaryOp, OpNeg),
		NewInstruction(OpNewStruct, "structType"),
		NewInstruction(OpGetField, "fieldName"),
		NewInstruction(OpSetField, "fieldName"),
		NewInstruction(OpSetStructField, "fieldName"),
		NewInstruction(OpGetIndex, 2),
		NewInstruction(OpSetIndex, 3),
		NewInstruction(OpRotate, 1),
		NewInstruction(OpNewSlice, 5),
		NewInstruction(OpLen, nil),
		NewInstruction(OpGetElement, 1),
		NewInstruction(OpCode(99), "unknown"), // Unknown opcode
	}

	// Just verify that none of them panic and all return non-empty strings
	for i, instr := range instructions {
		result := instr.String()
		if result == "" {
			t.Errorf("Instruction %d returned empty string", i)
		}
	}
}

func TestTypeSystem(t *testing.T) {
	vm := NewVM()

	// Test RegisterType and GetType
	if vm.typeSystem["int"] == nil {
		t.Error("Expected int type to be pre-registered")
	}

	if vm.typeSystem["float64"] == nil {
		t.Error("Expected float64 type to be pre-registered")
	}

	if vm.typeSystem["string"] == nil {
		t.Error("Expected string type to be pre-registered")
	}

	if vm.typeSystem["bool"] == nil {
		t.Error("Expected bool type to be pre-registered")
	}

	// Test GetFunctionRegistry and SetFunctionRegistry
	registry := vm.GetFunctionRegistry()
	if registry == nil {
		t.Error("Expected non-nil function registry from GetFunctionRegistry")
	}

	newRegistry := make(map[string]func(args ...interface{}) (interface{}, error))
	vm.SetFunctionRegistry(newRegistry)

	updatedRegistry := vm.GetFunctionRegistry()
	if len(updatedRegistry) != 0 {
		t.Errorf("Expected empty registry after SetFunctionRegistry, got %d functions", len(updatedRegistry))
	}
}

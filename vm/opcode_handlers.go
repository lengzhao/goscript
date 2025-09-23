// Package vm provides optimized opcode handlers using dispatch table
package vm

import (
	"fmt"
	"strings"

	"github.com/lengzhao/goscript/instruction"
	"github.com/lengzhao/goscript/types"
)

// OpHandler represents a function that handles a specific opcode
type OpHandler func(v *VM, instr *Instruction) error

// VmError represents a VM execution error with context
type VmError struct {
	Type    ErrorType
	Message string
	IP      int
	OpCode  OpCode
}

type ErrorType int

const (
	ErrorTypeStackUnderflow ErrorType = iota
	ErrorTypeTypeMismatch
	ErrorTypeUndefinedVariable
	ErrorTypeDivisionByZero
	ErrorTypeInvalidOperation
	ErrorTypeInvalidJumpTarget
	ErrorTypeUndefinedFunction
)

func (e *VmError) Error() string {
	return fmt.Sprintf("%v at IP %d (opcode %v): %s", e.Type, e.IP, e.OpCode, e.Message)
}

// registerHandlers registers all opcode handlers
func (vm *VM) registerHandlers() {
	vm.handlers = make(map[OpCode]OpHandler)

	// Basic operations
	vm.handlers[OpNop] = vm.handleNop
	vm.handlers[OpLoadConst] = vm.handleLoadConst
	vm.handlers[OpLoadName] = vm.handleLoadName
	vm.handlers[OpStoreName] = vm.handleStoreName
	vm.handlers[OpPop] = vm.handlePop
	vm.handlers[OpSwap] = vm.handleSwap
	vm.handlers[instruction.OpCreateVar] = vm.handleCreateVar

	// Control flow
	vm.handlers[OpCall] = vm.handleCall
	vm.handlers[OpCallMethod] = vm.handleCallMethod
	vm.handlers[OpRegistFunction] = vm.handleRegistFunction
	vm.handlers[OpReturn] = vm.handleReturn
	vm.handlers[OpJump] = vm.handleJump
	vm.handlers[OpJumpIf] = vm.handleJumpIf

	// Arithmetic operations
	vm.handlers[OpBinaryOp] = vm.handleBinaryOp
	vm.handlers[OpUnaryOp] = vm.handleUnaryOp

	// Data structures
	vm.handlers[OpNewStruct] = vm.handleNewStruct
	vm.handlers[OpNewSlice] = vm.handleNewSlice
	vm.handlers[OpGetField] = vm.handleGetField
	vm.handlers[OpSetField] = vm.handleSetField
	vm.handlers[OpSetStructField] = vm.handleSetStructField
	vm.handlers[OpGetIndex] = vm.handleGetIndex
	vm.handlers[OpSetIndex] = vm.handleSetIndex
	vm.handlers[OpRotate] = vm.handleRotate
	vm.handlers[OpLen] = vm.handleLen
	vm.handlers[OpGetElement] = vm.handleGetElement
	vm.handlers[OpImport] = vm.handleImport

	// Context-based scope management (new)
	vm.handlers[instruction.OpEnterScopeWithKey] = vm.handleEnterScopeWithKey
	vm.handlers[instruction.OpExitScopeWithKey] = vm.handleExitScopeWithKey

	// Loop control
	vm.handlers[instruction.OpBreak] = vm.handleBreak
}

// Basic operation handlers

func (vm *VM) handleNop(v *VM, instr *Instruction) error {
	// Do nothing
	return nil
}

func (vm *VM) handleLoadConst(v *VM, instr *Instruction) error {
	v.Push(instr.Arg)
	return nil
}

// handleLoadName handles loading a variable by name from the context hierarchy
func (vm *VM) handleLoadName(v *VM, instr *Instruction) error {
	name, ok := instr.Arg.(string)
	if !ok {
		return &VmError{
			Type:    ErrorTypeTypeMismatch,
			Message: "OpLoadName requires string argument",
			IP:      v.ip,
			OpCode:  OpLoadName,
		}
	}

	// First try to get from context-based approach (new)
	if v.currentCtx != nil {
		if value, exists := v.currentCtx.GetVariable(name); exists {
			v.Push(value)
			return nil
		}
	}

	return &VmError{
		Type:    ErrorTypeUndefinedVariable,
		Message: fmt.Sprintf("undefined variable: %s", name),
		IP:      v.ip,
		OpCode:  OpLoadName,
	}
}

// handleStoreName handles storing a value to a variable by name in the current context
func (vm *VM) handleStoreName(v *VM, instr *Instruction) error {
	name, ok := instr.Arg.(string)
	if !ok {
		return &VmError{
			Type:    ErrorTypeTypeMismatch,
			Message: "OpStoreName requires string argument",
			IP:      v.ip,
			OpCode:  OpStoreName,
		}
	}

	if v.stack.IsEmpty() {
		return &VmError{
			Type:    ErrorTypeStackUnderflow,
			Message: "expected 1 value for STORE_NAME",
			IP:      v.ip,
			OpCode:  OpStoreName,
		}
	}

	value := v.Pop()

	// Use context-based approach (new) if available
	if v.currentCtx != nil {
		v.currentCtx.SetVariable(name, value)
		return nil // Return nil on success
	}

	// If no context is available, return an error
	return &VmError{
		Type:    ErrorTypeUndefinedVariable,
		Message: "no context available",
		IP:      v.ip,
		OpCode:  OpStoreName,
	}
}

func (vm *VM) handlePop(v *VM, instr *Instruction) error {
	if v.stack.IsEmpty() {
		return &VmError{
			Type:    ErrorTypeStackUnderflow,
			Message: "expected 1 value for POP",
			IP:      v.ip,
			OpCode:  OpPop,
		}
	}
	v.Pop()
	return nil
}

// handleSwap swaps the top two elements of the stack
func (vm *VM) handleSwap(v *VM, instr *Instruction) error {
	if v.stack.Size() < 2 {
		return &VmError{
			Type:    ErrorTypeStackUnderflow,
			Message: "expected 2 values for SWAP",
			IP:      v.ip,
			OpCode:  OpSwap,
		}
	}

	return v.stack.Swap()
}

// handleCreateVar handles creating a new variable in the current context
func (vm *VM) handleCreateVar(v *VM, instr *Instruction) error {
	name, ok := instr.Arg.(string)
	if !ok {
		return &VmError{
			Type:    ErrorTypeTypeMismatch,
			Message: "OpCreateVar requires string variable name",
			IP:      v.ip,
			OpCode:  instruction.OpCreateVar,
		}
	}

	// Create the variable in the current context with a nil value
	if v.currentCtx != nil {
		v.currentCtx.CreateVariableWithType(name, nil, "unknown")
	} else {
		// If no context is available, return an error
		return &VmError{
			Type:    ErrorTypeUndefinedVariable,
			Message: "no context available to create variable",
			IP:      v.ip,
			OpCode:  instruction.OpCreateVar,
		}
	}

	return nil
}

// Control flow handlers

// handleCall handles function calls
func (vm *VM) handleCall(v *VM, instr *Instruction) error {
	fnName, ok := instr.Arg.(string)
	if !ok {
		return &VmError{
			Type:    ErrorTypeTypeMismatch,
			Message: "OpCall requires string function name",
			IP:      v.ip,
			OpCode:  OpCall,
		}
	}

	argCount, ok := instr.Arg2.(int)
	if !ok {
		return &VmError{
			Type:    ErrorTypeTypeMismatch,
			Message: "OpCall requires int argument count",
			IP:      v.ip,
			OpCode:  OpCall,
		}
	}

	// Check if we have enough arguments on the stack
	if v.stack.Size() < argCount {
		return &VmError{
			Type:    ErrorTypeStackUnderflow,
			Message: fmt.Sprintf("expected %d arguments, got %d", argCount, v.stack.Size()),
			IP:      v.ip,
			OpCode:  OpCall,
		}
	}

	// Pop arguments from stack (in reverse order to maintain correct parameter order)
	args := make([]interface{}, argCount)
	for i := argCount - 1; i >= 0; i-- {
		args[i] = v.Pop()
	}

	// Handle receiver parameter for script functions
	if scriptFunc, exists := v.scriptFunctions[fnName]; exists {
		if scriptFunc.ReceiverType == "value" && len(args) > 0 {
			if objMap, ok := args[0].(map[string]interface{}); ok {
				args[0] = deepCopyMap(objMap)
			}
		}
	}

	// Execute function from registry
	result, err := v.executeFunction(fnName, args...)
	if err != nil {
		return &VmError{
			Type:    ErrorTypeUndefinedFunction,
			Message: err.Error(),
			IP:      v.ip,
			OpCode:  OpCall,
		}
	}

	v.Push(result)
	return nil
}

func (vm *VM) handleCallMethod(v *VM, instr *Instruction) error {
	fnName, ok := instr.Arg.(string)
	if !ok {
		return &VmError{
			Type:    ErrorTypeTypeMismatch,
			Message: "OpCallMethod requires string function name",
			IP:      v.ip,
			OpCode:  OpCallMethod,
		}
	}

	argCount, ok := instr.Arg2.(int)
	if !ok {
		return &VmError{
			Type:    ErrorTypeTypeMismatch,
			Message: "OpCallMethod requires int argument count",
			IP:      v.ip,
			OpCode:  OpCallMethod,
		}
	}

	// Check if we have enough arguments on the stack
	if v.stack.Size() < argCount {
		return &VmError{
			Type:    ErrorTypeStackUnderflow,
			Message: fmt.Sprintf("expected %d arguments, got %d", argCount, v.stack.Size()),
			IP:      v.ip,
			OpCode:  OpCallMethod,
		}
	}

	// Pop arguments from stack (in reverse order to maintain correct parameter order)
	args := make([]interface{}, argCount)
	for i := argCount - 1; i >= 0; i-- {
		args[i] = v.Pop()
	}

	// Store the original receiver for pointer receiver methods
	var originalReceiver interface{}
	if scriptFunc, exists := v.scriptFunctions[fnName]; exists && scriptFunc.ReceiverType == "pointer" && len(args) > 0 {
		originalReceiver = args[0]
	}

	// Handle receiver parameter for script functions
	if scriptFunc, exists := v.scriptFunctions[fnName]; exists {
		if scriptFunc.ReceiverType == "value" && len(args) > 0 {
			if objMap, ok := args[0].(map[string]interface{}); ok {
				args[0] = deepCopyMap(objMap)
			}
		}
	}

	// Execute function from registry
	result, err := v.executeFunction(fnName, args...)
	if err != nil {
		return &VmError{
			Type:    ErrorTypeUndefinedFunction,
			Message: err.Error(),
			IP:      v.ip,
			OpCode:  OpCallMethod,
		}
	}

	// For method calls, handle the return value based on receiver type
	if scriptFunc, exists := v.scriptFunctions[fnName]; exists {
		if scriptFunc.ReceiverType == "pointer" {
			// For pointer receiver methods, we need to update the original receiver
			// with any modifications made during the method execution
			if originalReceiver != nil && result != nil {
				// If the method returned a modified receiver, use that
				// Otherwise, the original receiver might have been modified in-place
				if modifiedReceiver, ok := result.(map[string]interface{}); ok {
					v.Push(modifiedReceiver)
				} else {
					v.Push(originalReceiver)
				}
			} else if originalReceiver != nil {
				// Push the original receiver which may have been modified in-place
				v.Push(originalReceiver)
			} else {
				// Push the result if no original receiver
				v.Push(result)
			}
		} else {
			// For value receiver methods, push the result
			v.Push(result)
		}
	} else {
		// For regular functions, push the result
		v.Push(result)
	}

	return nil
}

func (vm *VM) handleRegistFunction(v *VM, instr *Instruction) error {
	fnName, ok := instr.Arg.(string)
	if !ok {
		return &VmError{
			Type:    ErrorTypeTypeMismatch,
			Message: "OpRegistFunction requires string function name",
			IP:      v.ip,
			OpCode:  OpRegistFunction,
		}
	}

	funcInfo, ok := instr.Arg2.(*ScriptFunction)
	if !ok {
		return &VmError{
			Type:    ErrorTypeTypeMismatch,
			Message: "OpRegistFunction requires *ScriptFunction argument",
			IP:      v.ip,
			OpCode:  OpRegistFunction,
		}
	}

	v.scriptFunctions[fnName] = funcInfo
	if v.debug {
		fmt.Printf("Registered function: %s with %d params\n", fnName, funcInfo.ParamCount)
	}
	return nil
}

func (vm *VM) handleReturn(v *VM, instr *Instruction) error {
	if v.stack.IsEmpty() {
		return &VmError{
			Type:    ErrorTypeStackUnderflow,
			Message: "expected 1 value for RETURN",
			IP:      v.ip,
			OpCode:  OpReturn,
		}
	}
	v.retval = v.Pop()
	if v.debug {
		fmt.Printf("Return value: %v\n", v.retval)
	}
	// Set IP to end of instructions to exit the execution loop
	v.ip = len(v.instructions)
	return nil
}

func (vm *VM) handleJump(v *VM, instr *Instruction) error {
	target, ok := instr.Arg.(int)
	if !ok {
		return &VmError{
			Type:    ErrorTypeTypeMismatch,
			Message: "OpJump requires int target",
			IP:      v.ip,
			OpCode:  OpJump,
		}
	}

	if target < 0 || target >= len(v.instructions) {
		return &VmError{
			Type:    ErrorTypeInvalidJumpTarget,
			Message: fmt.Sprintf("invalid jump target: %d", target),
			IP:      v.ip,
			OpCode:  OpJump,
		}
	}

	v.ip = target - 1 // -1 because we increment at the end of the loop
	return nil
}

func (vm *VM) handleJumpIf(v *VM, instr *Instruction) error {
	if v.stack.IsEmpty() {
		return &VmError{
			Type:    ErrorTypeStackUnderflow,
			Message: "expected 1 value for JUMP_IF",
			IP:      v.ip,
			OpCode:  OpJumpIf,
		}
	}

	condition := v.Pop()
	target, ok := instr.Arg.(int)
	if !ok {
		return &VmError{
			Type:    ErrorTypeTypeMismatch,
			Message: "OpJumpIf requires int target",
			IP:      v.ip,
			OpCode:  OpJumpIf,
		}
	}

	// Jump if condition is FALSE (negate the condition)
	if !isTruthy(condition) {
		if target < 0 || target >= len(v.instructions) {
			return &VmError{
				Type:    ErrorTypeInvalidJumpTarget,
				Message: fmt.Sprintf("invalid jump target: %d", target),
				IP:      v.ip,
				OpCode:  OpJumpIf,
			}
		}
		v.ip = target - 1 // -1 because we increment at the end of the loop
	}

	return nil
}

// Arithmetic operation handlers

func (vm *VM) handleBinaryOp(v *VM, instr *Instruction) error {
	if v.stack.Size() < 2 {
		return &VmError{
			Type:    ErrorTypeStackUnderflow,
			Message: "expected 2 values for BINARY_OP",
			IP:      v.ip,
			OpCode:  OpBinaryOp,
		}
	}

	right := v.Pop()
	left := v.Pop()

	op, ok := instr.Arg.(BinaryOp)
	if !ok {
		return &VmError{
			Type:    ErrorTypeTypeMismatch,
			Message: "OpBinaryOp requires BinaryOp argument",
			IP:      v.ip,
			OpCode:  OpBinaryOp,
		}
	}

	result, err := v.executeBinaryOp(op, left, right)
	if err != nil {
		return &VmError{
			Type:    ErrorTypeInvalidOperation,
			Message: err.Error(),
			IP:      v.ip,
			OpCode:  OpBinaryOp,
		}
	}

	v.Push(result)
	return nil
}

func (vm *VM) handleUnaryOp(v *VM, instr *Instruction) error {
	if v.stack.Size() < 1 {
		return &VmError{
			Type:    ErrorTypeStackUnderflow,
			Message: "expected 1 value for UNARY_OP",
			IP:      v.ip,
			OpCode:  OpUnaryOp,
		}
	}

	value := v.Pop()

	op, ok := instr.Arg.(UnaryOp)
	if !ok {
		return &VmError{
			Type:    ErrorTypeTypeMismatch,
			Message: "OpUnaryOp requires UnaryOp argument",
			IP:      v.ip,
			OpCode:  OpUnaryOp,
		}
	}

	result, err := v.executeUnaryOp(op, value)
	if err != nil {
		return &VmError{
			Type:    ErrorTypeInvalidOperation,
			Message: err.Error(),
			IP:      v.ip,
			OpCode:  OpUnaryOp,
		}
	}

	v.Push(result)
	return nil
}

// Data structure handlers

func (vm *VM) handleNewStruct(v *VM, instr *Instruction) error {
	var structInstance map[string]interface{}

	// If a type name is provided, try to create a struct with default values
	if typeName, ok := instr.Arg.(string); ok && typeName != "" {
		if structType, exists := v.typeSystem[typeName]; exists {
			// Create a struct with default values
			if structTypeDef, ok := structType.(*types.StructType); ok {
				structInstance = make(map[string]interface{})
				// Initialize with default values
				for fieldName, fieldType := range structTypeDef.GetFields() {
					structInstance[fieldName] = fieldType.DefaultValue()
				}
			} else {
				// Fallback to empty map
				structInstance = make(map[string]interface{})
			}
		} else {
			// Type not found, create empty map
			structInstance = make(map[string]interface{})
		}
	} else {
		// No type specified, create empty map
		structInstance = make(map[string]interface{})
	}

	v.Push(structInstance)
	return nil
}

func (vm *VM) handleNewSlice(v *VM, instr *Instruction) error {
	elementCount, ok := instr.Arg.(int)
	if !ok {
		elementCount = 0
	}

	// Check if we have enough elements on the stack
	if v.stack.Size() < elementCount {
		return &VmError{
			Type:    ErrorTypeStackUnderflow,
			Message: fmt.Sprintf("expected %d elements, got %d", elementCount, v.stack.Size()),
			IP:      v.ip,
			OpCode:  OpNewSlice,
		}
	}

	// Create a slice with the elements
	slice := make([]interface{}, elementCount)
	// Pop elements from stack in reverse order to maintain correct order
	for i := elementCount - 1; i >= 0; i-- {
		slice[i] = v.Pop()
	}

	v.Push(slice)
	return nil
}

func (vm *VM) handleGetField(v *VM, instr *Instruction) error {
	if v.stack.Size() < 2 {
		return &VmError{
			Type:    ErrorTypeStackUnderflow,
			Message: "expected 2 values for GET_FIELD",
			IP:      v.ip,
			OpCode:  OpGetField,
		}
	}

	fieldName := v.Pop()
	obj := v.Pop()

	// Convert fieldName to string if it's not already
	var fieldNameStr string
	if s, ok := fieldName.(string); ok {
		fieldNameStr = s
	} else {
		fieldNameStr = fmt.Sprintf("%v", fieldName)
	}

	// For now, we assume obj is a map
	if objMap, ok := obj.(map[string]interface{}); ok {
		// Try to find the field, including in embedded structs
		field := getFieldFromStruct(objMap, fieldNameStr)
		v.Push(field)
	} else {
		// If obj is not a map, return nil
		v.Push(nil)
	}

	return nil
}

func (vm *VM) handleSetField(v *VM, instr *Instruction) error {
	if v.stack.Size() < 3 {
		return &VmError{
			Type:    ErrorTypeStackUnderflow,
			Message: "expected 3 values for SET_FIELD",
			IP:      v.ip,
			OpCode:  OpSetField,
		}
	}

	// Pop in reverse order to get the correct values
	// Stack order should be: [object, fieldName, value]
	value := v.Pop()
	fieldName := v.Pop()
	obj := v.Pop()

	// Convert fieldName to string if it's not already
	var fieldNameStr string
	if s, ok := fieldName.(string); ok {
		fieldNameStr = s
	} else {
		fieldNameStr = fmt.Sprintf("%v", fieldName)
	}

	// For now, we assume obj is a map
	if objMap, ok := obj.(map[string]interface{}); ok {
		if v.debug {
			fmt.Printf("Setting field '%s' to '%v' in object: %v\n", fieldNameStr, value, objMap)
		}
		objMap[fieldNameStr] = value
		if v.debug {
			fmt.Printf("Object after setting field: %v\n", objMap)
		}
		// Push the modified object back onto the stack
		v.Push(objMap)
	} else {
		if v.debug {
			fmt.Printf("Object is not a map, type: %T, value: %v, pushing back unchanged\n", obj, obj)
		}
		// If obj is not a map, push it back unchanged
		v.Push(obj)
	}

	return nil
}

func (vm *VM) handleSetStructField(v *VM, instr *Instruction) error {
	if v.stack.Size() < 3 {
		return &VmError{
			Type:    ErrorTypeStackUnderflow,
			Message: "expected 3 values for SET_STRUCT_FIELD",
			IP:      v.ip,
			OpCode:  OpSetStructField,
		}
	}

	// Pop in reverse order to get the correct values
	// Stack order should be: [object, fieldName, value]
	value := v.Pop()
	fieldName := v.Pop()
	obj := v.Pop()

	// Convert fieldName to string if it's not already
	var fieldNameStr string
	if s, ok := fieldName.(string); ok {
		fieldNameStr = s
	} else {
		fieldNameStr = fmt.Sprintf("%v", fieldName)
	}

	// For now, we assume obj is a map
	if objMap, ok := obj.(map[string]interface{}); ok {
		if v.debug {
			fmt.Printf("Setting struct field '%s' to '%v' in object: %v\n", fieldNameStr, value, objMap)
		}
		objMap[fieldNameStr] = value
		if v.debug {
			fmt.Printf("Object after setting struct field: %v\n", objMap)
		}
		// Push the modified object back onto the stack
		v.Push(objMap)
	} else {
		if v.debug {
			fmt.Printf("Object is not a map, type: %T, value: %v, pushing back unchanged\n", obj, obj)
		}
		// If obj is not a map, push it back unchanged
		v.Push(obj)
	}

	return nil
}

func (vm *VM) handleGetIndex(v *VM, instr *Instruction) error {
	if v.stack.Size() < 2 {
		return &VmError{
			Type:    ErrorTypeStackUnderflow,
			Message: "expected 2 values for GET_INDEX",
			IP:      v.ip,
			OpCode:  OpGetIndex,
		}
	}

	index := v.Pop()
	array := v.Pop()

	// Convert index to int if it's not already
	var indexInt int
	if i, ok := index.(int); ok {
		indexInt = i
	} else {
		return &VmError{
			Type:    ErrorTypeTypeMismatch,
			Message: fmt.Sprintf("index must be an integer, got %T", index),
			IP:      v.ip,
			OpCode:  OpGetIndex,
		}
	}

	// For now, we assume array is a slice of interface{}
	if arraySlice, ok := array.([]interface{}); ok {
		if indexInt >= 0 && indexInt < len(arraySlice) {
			v.Push(arraySlice[indexInt])
		} else {
			v.Push(nil)
		}
	} else {
		// If array is not a slice, return nil
		v.Push(nil)
	}

	return nil
}

func (vm *VM) handleSetIndex(v *VM, instr *Instruction) error {
	if v.stack.Size() < 3 {
		return &VmError{
			Type:    ErrorTypeStackUnderflow,
			Message: "expected 3 values for SET_INDEX",
			IP:      v.ip,
			OpCode:  OpSetIndex,
		}
	}

	value := v.Pop()
	index := v.Pop()
	array := v.Pop()

	// Convert index to int if it's not already
	var indexInt int
	if i, ok := index.(int); ok {
		indexInt = i
	} else {
		return &VmError{
			Type:    ErrorTypeTypeMismatch,
			Message: fmt.Sprintf("index must be an integer, got %T", index),
			IP:      v.ip,
			OpCode:  OpSetIndex,
		}
	}

	// For now, we assume array is a slice of interface{}
	if arraySlice, ok := array.([]interface{}); ok {
		if indexInt >= 0 && indexInt < len(arraySlice) {
			arraySlice[indexInt] = value
			// Push the modified array back onto the stack
			v.Push(arraySlice)
		} else {
			// Index out of bounds, push the array back unchanged
			v.Push(arraySlice)
		}
	} else {
		// If array is not a slice, push it back unchanged
		v.Push(array)
	}

	return nil
}

func (vm *VM) handleRotate(v *VM, instr *Instruction) error {
	if v.stack.Size() < 3 {
		return &VmError{
			Type:    ErrorTypeStackUnderflow,
			Message: "expected 3 values for ROTATE",
			IP:      v.ip,
			OpCode:  OpRotate,
		}
	}

	// Rotate the top three elements: [a, b, c] -> [b, c, a]
	return v.stack.Rotate(3)
}

func (vm *VM) handleLen(v *VM, instr *Instruction) error {
	if v.stack.Size() < 1 {
		return &VmError{
			Type:    ErrorTypeStackUnderflow,
			Message: "expected 1 value for LEN",
			IP:      v.ip,
			OpCode:  OpLen,
		}
	}

	value := v.Pop()

	// Get the length of the value if it's a slice or array
	switch val := value.(type) {
	case []interface{}:
		v.Push(len(val))
	default:
		// For other types, return 0
		v.Push(0)
	}

	return nil
}

func (vm *VM) handleGetElement(v *VM, instr *Instruction) error {
	if v.stack.Size() < 2 {
		return &VmError{
			Type:    ErrorTypeStackUnderflow,
			Message: "expected 2 values for GET_ELEMENT",
			IP:      v.ip,
			OpCode:  OpGetElement,
		}
	}

	index := v.Pop()
	array := v.Pop()

	// Convert index to int if it's not already
	var indexInt int
	if i, ok := index.(int); ok {
		indexInt = i
	} else {
		return &VmError{
			Type:    ErrorTypeTypeMismatch,
			Message: fmt.Sprintf("index must be an integer, got %T", index),
			IP:      v.ip,
			OpCode:  OpGetElement,
		}
	}

	// Get the element at the index if array is a slice
	if arraySlice, ok := array.([]interface{}); ok {
		if indexInt >= 0 && indexInt < len(arraySlice) {
			v.Push(arraySlice[indexInt])
		} else {
			v.Push(nil)
		}
	} else {
		v.Push(nil)
	}

	return nil
}

func (vm *VM) handleImport(v *VM, instr *Instruction) error {
	// Handle import instruction
	// Pop the module name from the stack
	if v.stack.Size() < 1 {
		return &VmError{
			Type:    ErrorTypeStackUnderflow,
			Message: "expected 1 value for IMPORT",
			IP:      v.ip,
			OpCode:  OpImport,
		}
	}

	modulePath := v.Pop()

	// Convert module path to string
	var modulePathStr string
	if s, ok := modulePath.(string); ok {
		modulePathStr = s
	} else {
		modulePathStr = fmt.Sprintf("%v", modulePath)
	}

	// Extract module name from path (last part)
	parts := strings.Split(modulePathStr, "/")
	moduleName := parts[len(parts)-1]

	// Store the module name as a global variable so it can be referenced
	if v.currentCtx != nil {
		// Create the module variable in the global context
		v.currentCtx.CreateVariableWithType(moduleName, moduleName, "module")
	}

	// Debug output
	if v.debug {
		fmt.Printf("IMPORT: Imported module %s\n", moduleName)
	}

	return nil
}

// handleEnterScopeWithKey handles entering a scope with a specific key
func (vm *VM) handleEnterScopeWithKey(v *VM, instr *Instruction) error {
	pathKey, ok := instr.Arg.(string)
	if !ok {
		return &VmError{
			Type:    ErrorTypeTypeMismatch,
			Message: "OpEnterScopeWithKey requires string path key",
			IP:      v.ip,
			OpCode:  instruction.OpEnterScopeWithKey,
		}
	}

	// Enter the new scope
	if v.currentCtx != nil {
		v.EnterScope(pathKey)
	} else {
		// Fallback behavior if context not initialized
		if v.debug {
			fmt.Printf("Entering scope: %s (context not initialized)\n", pathKey)
		}
	}

	return nil
}

// handleExitScopeWithKey handles exiting a scope with a specific key
func (vm *VM) handleExitScopeWithKey(v *VM, instr *Instruction) error {
	pathKey, ok := instr.Arg.(string)
	if !ok {
		return &VmError{
			Type:    ErrorTypeTypeMismatch,
			Message: "OpExitScopeWithKey requires string path key",
			IP:      v.ip,
			OpCode:  instruction.OpExitScopeWithKey,
		}
	}

	// Exit the current scope
	if v.currentCtx != nil {
		currentPath := v.currentCtx.GetPathKey()
		if currentPath != pathKey {
			if v.debug {
				fmt.Printf("Warning: Exiting scope %s but current scope is %s\n", pathKey, currentPath)
			}
		}
		v.ExitScope()
	} else {
		// Fallback behavior if context not initialized
		if v.debug {
			fmt.Printf("Exiting scope: %s (context not initialized)\n", pathKey)
		}
	}

	return nil
}

// handleBreak handles break statement
func (vm *VM) handleBreak(v *VM, instr *Instruction) error {
	// For break statement, we need to jump to the end of the current loop
	// Looking at the instruction sequence, we can see that:
	// - The loop body ends with a JUMP instruction that goes back to the start
	// - We need to skip that JUMP and go to the next instruction after it

	// Find the next JUMP instruction after the current IP and jump to the instruction after that
	for i := v.ip + 1; i < len(v.instructions); i++ {
		if v.instructions[i].Op == OpJump {
			// Jump to the instruction after the JUMP
			v.ip = i
			return nil
		}
	}

	// If no JUMP found, jump to end of instructions
	v.ip = len(v.instructions) - 1
	return nil
}

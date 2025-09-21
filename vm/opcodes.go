// Package vm provides the virtual machine implementation
package vm

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/lengzhao/goscript/instruction"
	"github.com/lengzhao/goscript/types"
)

// OpCode represents an operation code for the virtual machine
type OpCode = instruction.OpCode

const (
	// No operation
	OpNop = instruction.OpNop

	// Load a constant onto the stack
	OpLoadConst = instruction.OpLoadConst

	// Load a variable by name
	OpLoadName = instruction.OpLoadName

	// Store a value to a variable by name
	OpStoreName = instruction.OpStoreName

	// Pop a value from the stack (discard it)
	OpPop = instruction.OpPop

	// Call a function
	OpCall = instruction.OpCall

	// Call a struct method
	OpCallMethod = instruction.OpCallMethod

	// Register a script-defined function
	OpRegistFunction = instruction.OpRegistFunction

	// Return from function
	OpReturn = instruction.OpReturn

	// Unconditional jump
	OpJump = instruction.OpJump

	// Conditional jump
	OpJumpIf = instruction.OpJumpIf

	// Binary operation (add, sub, mul, div, etc.)
	OpBinaryOp = instruction.OpBinaryOp

	// Unary operation (neg, not, etc.)
	OpUnaryOp = instruction.OpUnaryOp

	// Create a new struct instance
	OpNewStruct = instruction.OpNewStruct

	// Access a field of a struct
	OpGetField = instruction.OpGetField

	// Set a field of a struct
	OpSetField = instruction.OpSetField

	// Set a field of a struct with explicit stack order
	OpSetStructField = instruction.OpSetStructField

	// Access an element of an array/slice by index
	OpGetIndex = instruction.OpGetIndex

	// Set an element of an array/slice by index
	OpSetIndex = instruction.OpSetIndex

	// Rotate the top three elements on the stack
	// Changes [a, b, c] to [b, c, a]
	OpRotate = instruction.OpRotate

	// Create a new slice
	OpNewSlice = instruction.OpNewSlice

	// Get the length of a slice or array
	OpLen = instruction.OpLen

	// Get an element from a slice or array by index
	OpGetElement = instruction.OpGetElement

	// Import a module
	OpImport = instruction.OpImport
)

// BinaryOp represents a binary operation
type BinaryOp = instruction.BinaryOp

const (
	OpAdd          BinaryOp = instruction.OpAdd
	OpSub          BinaryOp = instruction.OpSub
	OpMul          BinaryOp = instruction.OpMul
	OpDiv          BinaryOp = instruction.OpDiv
	OpMod          BinaryOp = instruction.OpMod
	OpEqual        BinaryOp = instruction.OpEqual
	OpNotEqual     BinaryOp = instruction.OpNotEqual
	OpLess         BinaryOp = instruction.OpLess
	OpLessEqual    BinaryOp = instruction.OpLessEqual
	OpGreater      BinaryOp = instruction.OpGreater
	OpGreaterEqual BinaryOp = instruction.OpGreaterEqual
	OpAnd          BinaryOp = instruction.OpAnd
	OpOr           BinaryOp = instruction.OpOr
)

// UnaryOp represents a unary operation
type UnaryOp = instruction.UnaryOp

const (
	OpNeg UnaryOp = instruction.OpNeg
	OpNot UnaryOp = instruction.OpNot
)

// Instruction represents a single VM instruction
type Instruction = instruction.Instruction

// NewInstruction creates a new instruction
func NewInstruction(op OpCode, arg interface{}, arg2 ...interface{}) *Instruction {
	return instruction.NewInstruction(op, arg, arg2...)
}

// VM represents the virtual machine
type VM struct {
	// Stack for expression evaluation
	stack []interface{}

	// Global variables
	globals map[string]interface{}

	// Local variables
	locals map[string]interface{}

	// Instructions
	instructions []*Instruction

	// Instruction pointer
	ip int

	// Return value
	retval interface{}

	// Function registry (for unified function calling)
	functionRegistry map[string]func(args ...interface{}) (interface{}, error)

	// Script function definitions (name -> startIP, endIP, paramCount)
	scriptFunctions map[string]*ScriptFunction

	// Type system
	typeSystem map[string]types.IType

	// Debug mode
	debug bool

	// Execution statistics
	executionCount int

	// Maximum instruction limit to prevent infinite loops
	maxInstructions int64

	// Module manager interface for on-demand module function access
	moduleManager ModuleManagerInterface
}

// ScriptFunction represents a script-defined function
type ScriptFunction struct {
	Name         string
	StartIP      int
	EndIP        int
	ParamCount   int
	ParamNames   []string // Store parameter names for use in function body
	ReceiverType string   // Store receiver type ("value" or "pointer")
}

// NewVM creates a new virtual machine
func NewVM() *VM {
	// Initialize with basic types
	typeSystem := make(map[string]types.IType)
	typeSystem["int"] = types.IntType.Clone()
	typeSystem["float64"] = types.Float64Type.Clone()
	typeSystem["string"] = types.StringType.Clone()
	typeSystem["bool"] = types.BoolType.Clone()

	return &VM{
		stack:            make([]interface{}, 0),
		globals:          make(map[string]interface{}),
		locals:           make(map[string]interface{}),
		instructions:     make([]*Instruction, 0),
		ip:               0,
		functionRegistry: make(map[string]func(args ...interface{}) (interface{}, error)),
		scriptFunctions:  make(map[string]*ScriptFunction),
		typeSystem:       typeSystem,
		debug:            false,
		executionCount:   0,
		maxInstructions:  1000000, // 默认最大指令数限制为100万条指令
		moduleManager:    nil,     // 初始化为空
	}
}

// RegisterFunction registers a function in the VM
func (vm *VM) RegisterFunction(name string, fn func(args ...interface{}) (interface{}, error)) {
	vm.functionRegistry[name] = fn
}

// RegisterScriptFunction registers a script-defined function
func (vm *VM) RegisterScriptFunction(name string, startIP, endIP, paramCount int) {
	vm.scriptFunctions[name] = &ScriptFunction{
		Name:       name,
		StartIP:    startIP,
		EndIP:      endIP,
		ParamCount: paramCount,
	}
}

// GetScriptFunction returns a script function by name
func (vm *VM) GetScriptFunction(name string) (*ScriptFunction, bool) {
	scriptFunc, exists := vm.scriptFunctions[name]
	return scriptFunc, exists
}

// Push pushes a value onto the stack
func (vm *VM) Push(value interface{}) {
	vm.stack = append(vm.stack, value)
}

// Pop pops a value from the stack
func (vm *VM) Pop() interface{} {
	if len(vm.stack) == 0 {
		return nil
	}

	value := vm.stack[len(vm.stack)-1]
	vm.stack = vm.stack[:len(vm.stack)-1]
	return value
}

// Peek returns the top value without removing it
func (vm *VM) Peek() interface{} {
	if len(vm.stack) == 0 {
		return nil
	}

	return vm.stack[len(vm.stack)-1]
}

// StackSize returns the current stack size
func (vm *VM) StackSize() int {
	return len(vm.stack)
}

// SetGlobal sets a global variable
func (vm *VM) SetGlobal(name string, value interface{}) {
	vm.globals[name] = value
}

// GetGlobal gets a global variable
func (vm *VM) GetGlobal(name string) (interface{}, bool) {
	value, ok := vm.globals[name]
	return value, ok
}

// SetLocal sets a local variable
func (vm *VM) SetLocal(name string, value interface{}) {
	if vm.locals == nil {
		vm.locals = make(map[string]interface{})
	}
	vm.locals[name] = value
}

// GetLocal gets a local variable
func (vm *VM) GetLocal(name string) (interface{}, bool) {
	if vm.locals == nil {
		return nil, false
	}
	value, ok := vm.locals[name]
	return value, ok
}

// AddInstruction adds an instruction to the VM
func (vm *VM) AddInstruction(instr *Instruction) {
	vm.instructions = append(vm.instructions, instr)
}

// GetInstructions returns all instructions
func (vm *VM) GetInstructions() []*Instruction {
	return vm.instructions
}

// Clear clears the VM state
func (vm *VM) Clear() {
	vm.stack = vm.stack[:0]
	vm.locals = make(map[string]interface{})
	vm.instructions = vm.instructions[:0]
	vm.ip = 0
	vm.retval = nil
}

// SetDebug enables or disables debug mode
func (vm *VM) SetDebug(debug bool) {
	vm.debug = debug
}

// GetExecutionCount returns the number of instructions executed
func (vm *VM) GetExecutionCount() int {
	return vm.executionCount
}

// SetMaxInstructions sets the maximum number of instructions that can be executed
func (vm *VM) SetMaxInstructions(maxInstructions int64) {
	vm.maxInstructions = maxInstructions
}

// GetMaxInstructions returns the maximum number of instructions that can be executed
func (vm *VM) GetMaxInstructions() int64 {
	return vm.maxInstructions
}

// Helper function to get a field from a struct, including embedded structs
func getFieldFromStruct(objMap map[string]interface{}, fieldName string) interface{} {
	// First try to find the field directly
	if field, exists := objMap[fieldName]; exists {
		return field
	}

	// If not found directly, try to find it in embedded structs
	// In a real implementation, we would need to recursively search embedded structs
	// For now, we'll just check if there's an embedded struct with the same name as the field
	// This is a simplified approach for demonstration purposes
	for key, value := range objMap {
		// Check if the key matches the field name (for embedded structs)
		if key == fieldName {
			return value
		}

		// If the value is a map (embedded struct), recursively search it
		if embeddedMap, ok := value.(map[string]interface{}); ok {
			if embeddedField := getFieldFromStruct(embeddedMap, fieldName); embeddedField != nil {
				return embeddedField
			}
		}
	}

	// If not found, return nil
	return nil
}

// Execute executes the instructions
func (vm *VM) Execute(ctx context.Context) (interface{}, error) {
	vm.ip = 0
	vm.executionCount = 0

	for vm.ip < len(vm.instructions) {
		// Check if maximum instruction limit is exceeded
		if vm.maxInstructions > 0 && int64(vm.executionCount) >= vm.maxInstructions {
			return nil, fmt.Errorf("maximum instruction limit exceeded: %d instructions executed", vm.executionCount)
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		instr := vm.instructions[vm.ip]

		// Debug output
		if vm.debug {
			fmt.Printf("IP: %d, Instruction: %s, Stack: %v\n", vm.ip, instr.String(), vm.stack)
		}

		switch instr.Op {
		case OpNop:
			// Do nothing
		case OpLoadConst:
			vm.Push(instr.Arg)
		case OpLoadName:
			name := instr.Arg.(string)
			if value, ok := vm.GetLocal(name); ok {
				vm.Push(value)
			} else if value, ok := vm.GetGlobal(name); ok {
				vm.Push(value)
			} else {
				return nil, fmt.Errorf("undefined variable: %s", name)
			}
		case OpStoreName:
			name := instr.Arg.(string)
			if len(vm.stack) == 0 {
				return nil, fmt.Errorf("stack underflow in STORE_NAME")
			}
			value := vm.Pop()
			vm.SetLocal(name, value)
		case OpPop:
			if len(vm.stack) == 0 {
				return nil, fmt.Errorf("stack underflow in POP")
			}
			vm.Pop()
		case OpCall:
			// Function call using unified function registry
			fnName := instr.Arg.(string)
			argCount := instr.Arg2.(int)

			// Check if we have enough arguments on the stack
			if len(vm.stack) < argCount {
				return nil, fmt.Errorf("stack underflow in CALL: expected %d arguments, got %d", argCount, len(vm.stack))
			}

			// Pop arguments from stack (in reverse order to maintain correct parameter order)
			args := make([]interface{}, argCount)
			for i := argCount - 1; i >= 0; i-- {
				args[i] = vm.Pop()
			}

			// Check if it's a module function call (format: moduleName.functionName)
			// Module functions should already be registered in the function registry
			// For script-defined functions, check if the receiver parameter needs to be handled specially
			if scriptFunc, exists := vm.scriptFunctions[fnName]; exists {
				// Handle receiver parameter based on receiver type
				if scriptFunc.ReceiverType == "value" && len(args) > 0 {
					// If it's a map (struct) and the receiver type is "value", create a copy
					if objMap, ok := args[0].(map[string]interface{}); ok {
						args[0] = deepCopyMap(objMap)
					}
				}
				// For pointer receivers, we pass the struct as-is (no copy needed)
			}

			// Execute function from registry
			result, err := vm.executeFunction(fnName, args...)
			if err != nil {
				return nil, err
			}

			vm.Push(result)
		case OpCallMethod:
			// Method call for struct methods
			fnName := instr.Arg.(string)
			argCount := instr.Arg2.(int)

			// Check if we have enough arguments on the stack
			if len(vm.stack) < argCount {
				return nil, fmt.Errorf("stack underflow in CALL_METHOD: expected %d arguments, got %d", argCount, len(vm.stack))
			}

			// Pop arguments from stack (in reverse order to maintain correct parameter order)
			args := make([]interface{}, argCount)
			for i := argCount - 1; i >= 0; i-- {
				args[i] = vm.Pop()
			}

			// Store the original receiver for pointer receiver methods
			var originalReceiver interface{}
			if scriptFunc, exists := vm.scriptFunctions[fnName]; exists && scriptFunc.ReceiverType == "pointer" && len(args) > 0 {
				originalReceiver = args[0]
			}

			// For script-defined methods, check if the receiver parameter needs to be handled specially
			if scriptFunc, exists := vm.scriptFunctions[fnName]; exists {
				// Handle receiver parameter based on receiver type
				if scriptFunc.ReceiverType == "value" && len(args) > 0 {
					// If it's a map (struct) and the receiver type is "value", create a copy
					if objMap, ok := args[0].(map[string]interface{}); ok {
						args[0] = deepCopyMap(objMap)
					}
				}
				// For pointer receivers, we pass the struct as-is (no copy needed)
			}

			// Execute function from registry
			result, err := vm.executeFunction(fnName, args...)
			if err != nil {
				return nil, err
			}

			// For method calls, handle the return value based on receiver type
			if scriptFunc, exists := vm.scriptFunctions[fnName]; exists {
				if scriptFunc.ReceiverType == "pointer" {
					// For pointer receiver methods, push the modified receiver back onto the stack
					if originalReceiver != nil {
						vm.Push(originalReceiver)
					} else {
						vm.Push(result)
					}
				} else {
					// For value receiver methods, push the result
					vm.Push(result)
				}
			} else {
				// For regular functions, push the result
				vm.Push(result)
			}
		case OpRegistFunction:
			// Register a script-defined function
			fnName := instr.Arg.(string)
			funcInfo := instr.Arg2.(*ScriptFunction)
			vm.scriptFunctions[fnName] = funcInfo
			fmt.Printf("Registered function: %s with %d params\n", fnName, funcInfo.ParamCount)
		case OpReturn:
			if len(vm.stack) == 0 {
				return nil, fmt.Errorf("stack underflow in RETURN")
			}
			vm.retval = vm.Pop()
			fmt.Printf("Return value: %v\n", vm.retval)
			return vm.retval, nil
		case OpJump:
			target := instr.Arg.(int)
			if target < 0 || target >= len(vm.instructions) {
				return nil, fmt.Errorf("invalid jump target: %d", target)
			}
			vm.ip = target
			// 不在这里增加 executionCount，而是在循环末尾统一增加
			continue
		case OpJumpIf:
			if len(vm.stack) == 0 {
				return nil, fmt.Errorf("stack underflow in JUMP_IF")
			}
			condition := vm.Pop()
			// Jump if condition is FALSE (negate the condition)
			if !isTruthy(condition) {
				target := instr.Arg.(int)
				if target < 0 || target >= len(vm.instructions) {
					return nil, fmt.Errorf("invalid jump target: %d", target)
				}
				vm.ip = target
				// 不在这里增加 executionCount，而是在循环末尾统一增加
				continue
			}
		case OpBinaryOp:
			if len(vm.stack) < 2 {
				return nil, fmt.Errorf("stack underflow in BINARY_OP: expected 2 values, got %d", len(vm.stack))
			}
			right := vm.Pop()
			left := vm.Pop()

			result, err := vm.executeBinaryOp(instr.Arg.(BinaryOp), left, right)
			if err != nil {
				return nil, err
			}

			vm.Push(result)
		case OpUnaryOp:
			if len(vm.stack) < 1 {
				return nil, fmt.Errorf("stack underflow in UNARY_OP: expected 1 value, got %d", len(vm.stack))
			}
			value := vm.Pop()

			result, err := vm.executeUnaryOp(instr.Arg.(UnaryOp), value)
			if err != nil {
				return nil, err
			}

			vm.Push(result)
		case OpNewStruct:
			// Create a new struct instance based on type definition
			var structInstance map[string]interface{}

			// If a type name is provided, try to create a struct with default values
			if typeName, ok := instr.Arg.(string); ok && typeName != "" {
				if structType, exists := vm.typeSystem[typeName]; exists {
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

			vm.Push(structInstance)
		case OpNewSlice:
			// Create a new slice
			// Arg should be the number of elements
			elementCount, ok := instr.Arg.(int)
			if !ok {
				elementCount = 0
			}

			// Check if we have enough elements on the stack
			if len(vm.stack) < elementCount {
				return nil, fmt.Errorf("stack underflow in NEW_SLICE: expected %d elements, got %d", elementCount, len(vm.stack))
			}

			// Create a slice with the elements
			slice := make([]interface{}, elementCount)
			// Pop elements from stack in reverse order to maintain correct order
			for i := elementCount - 1; i >= 0; i-- {
				slice[i] = vm.Pop()
			}

			vm.Push(slice)
		case OpGetField:
			if len(vm.stack) < 2 {
				return nil, fmt.Errorf("stack underflow in GET_FIELD: expected 2 values, got %d", len(vm.stack))
			}
			fieldName := vm.Pop()
			obj := vm.Pop()

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
				vm.Push(field)
			} else {
				// If obj is not a map, return nil
				vm.Push(nil)
			}
		case OpSetField:
			// Debug output
			if vm.debug {
				fmt.Printf("About to execute SET_FIELD, Stack size: %d, Stack: %v\n", len(vm.stack), vm.stack)
			}

			if len(vm.stack) < 3 {
				return nil, fmt.Errorf("stack underflow in SET_FIELD: expected 3 values, got %d", len(vm.stack))
			}
			// Pop in reverse order to get the correct values
			// Stack order should be: [object, fieldName, value]
			value := vm.Pop()
			fieldName := vm.Pop()
			obj := vm.Pop()

			// Convert fieldName to string if it's not already
			var fieldNameStr string
			if s, ok := fieldName.(string); ok {
				fieldNameStr = s
			} else {
				fieldNameStr = fmt.Sprintf("%v", fieldName)
			}

			// For now, we assume obj is a map
			if objMap, ok := obj.(map[string]interface{}); ok {
				// Debug output
				if vm.debug {
					fmt.Printf("Setting field '%s' to '%v' in object: %v\n", fieldNameStr, value, objMap)
				}
				objMap[fieldNameStr] = value
				// Debug output
				if vm.debug {
					fmt.Printf("Object after setting field: %v\n", objMap)
				}
				// Push the modified object back onto the stack
				vm.Push(objMap)
			} else {
				// Debug output
				if vm.debug {
					fmt.Printf("Object is not a map, type: %T, value: %v, pushing back unchanged\n", obj, obj)
				}
				// If obj is not a map, push it back unchanged
				vm.Push(obj)
			}
		case OpSetStructField:
			// Debug output
			if vm.debug {
				fmt.Printf("About to execute SET_STRUCT_FIELD, Stack size: %d, Stack: %v\n", len(vm.stack), vm.stack)
			}

			if len(vm.stack) < 3 {
				return nil, fmt.Errorf("stack underflow in SET_STRUCT_FIELD: expected 3 values, got %d", len(vm.stack))
			}
			// Pop in reverse order to get the correct values
			// Stack order should be: [object, fieldName, value]
			value := vm.Pop()
			fieldName := vm.Pop()
			obj := vm.Pop()

			// Convert fieldName to string if it's not already
			var fieldNameStr string
			if s, ok := fieldName.(string); ok {
				fieldNameStr = s
			} else {
				fieldNameStr = fmt.Sprintf("%v", fieldName)
			}

			// For now, we assume obj is a map
			if objMap, ok := obj.(map[string]interface{}); ok {
				// Debug output
				if vm.debug {
					fmt.Printf("Setting struct field '%s' to '%v' in object: %v\n", fieldNameStr, value, objMap)
				}
				objMap[fieldNameStr] = value
				// Debug output
				if vm.debug {
					fmt.Printf("Object after setting struct field: %v\n", objMap)
				}
				// Push the modified object back onto the stack
				vm.Push(objMap)
			} else {
				// Debug output
				if vm.debug {
					fmt.Printf("Object is not a map, type: %T, value: %v, pushing back unchanged\n", obj, obj)
				}
				// If obj is not a map, push it back unchanged
				vm.Push(obj)
			}
		case OpGetIndex:
			if len(vm.stack) < 2 {
				return nil, fmt.Errorf("stack underflow in GET_INDEX: expected 2 values, got %d", len(vm.stack))
			}
			index := vm.Pop()
			array := vm.Pop()

			// Convert index to int if it's not already
			var indexInt int
			if i, ok := index.(int); ok {
				indexInt = i
			} else {
				return nil, fmt.Errorf("index must be an integer, got %T", index)
			}

			// For now, we assume array is a slice of interface{}
			if arraySlice, ok := array.([]interface{}); ok {
				if indexInt >= 0 && indexInt < len(arraySlice) {
					vm.Push(arraySlice[indexInt])
				} else {
					vm.Push(nil)
				}
			} else {
				// If array is not a slice, return nil
				vm.Push(nil)
			}
		case OpSetIndex:
			if len(vm.stack) < 3 {
				return nil, fmt.Errorf("stack underflow in SET_INDEX: expected 3 values, got %d", len(vm.stack))
			}
			value := vm.Pop()
			index := vm.Pop()
			array := vm.Pop()

			// Convert index to int if it's not already
			var indexInt int
			if i, ok := index.(int); ok {
				indexInt = i
			} else {
				return nil, fmt.Errorf("index must be an integer, got %T", index)
			}

			// For now, we assume array is a slice of interface{}
			if arraySlice, ok := array.([]interface{}); ok {
				if indexInt >= 0 && indexInt < len(arraySlice) {
					arraySlice[indexInt] = value
					// Push the modified array back onto the stack
					vm.Push(arraySlice)
				} else {
					// Index out of bounds, push the array back unchanged
					vm.Push(arraySlice)
				}
			} else {
				// If array is not a slice, push it back unchanged
				vm.Push(array)
			}
		case OpRotate:
			if len(vm.stack) < 3 {
				return nil, fmt.Errorf("stack underflow in ROTATE: expected 3 values, got %d", len(vm.stack))
			}
			// Rotate the top three elements: [a, b, c] -> [b, c, a]
			top := vm.stack[len(vm.stack)-1]
			second := vm.stack[len(vm.stack)-2]
			third := vm.stack[len(vm.stack)-3]
			vm.stack[len(vm.stack)-1] = third
			vm.stack[len(vm.stack)-2] = top
			vm.stack[len(vm.stack)-3] = second
		case OpLen:
			if len(vm.stack) < 1 {
				return nil, fmt.Errorf("stack underflow in LEN: expected 1 value, got %d", len(vm.stack))
			}
			value := vm.Pop()

			// Get the length of the value if it's a slice or array
			switch v := value.(type) {
			case []interface{}:
				vm.Push(len(v))
			default:
				// For other types, return 0
				vm.Push(0)
			}
		case OpGetElement:
			if len(vm.stack) < 2 {
				return nil, fmt.Errorf("stack underflow in GET_ELEMENT: expected 2 values, got %d", len(vm.stack))
			}
			index := vm.Pop()
			array := vm.Pop()

			// Convert index to int if it's not already
			var indexInt int
			if i, ok := index.(int); ok {
				indexInt = i
			} else {
				return nil, fmt.Errorf("index must be an integer, got %T", index)
			}

			// Get the element at the index if array is a slice
			if arraySlice, ok := array.([]interface{}); ok {
				if indexInt >= 0 && indexInt < len(arraySlice) {
					vm.Push(arraySlice[indexInt])
				} else {
					vm.Push(nil)
				}
			} else {
				vm.Push(nil)
			}
		case OpImport:
			// Handle import instruction
			// Pop the module name from the stack
			if len(vm.stack) < 1 {
				return nil, fmt.Errorf("stack underflow in IMPORT: expected 1 value, got %d", len(vm.stack))
			}
			modulePath := vm.Pop()

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
			vm.SetGlobal(moduleName, moduleName)

			// Debug output
			if vm.debug {
				fmt.Printf("IMPORT: Imported module %s\n", moduleName)
			}
		}

		vm.ip++
		vm.executionCount++
	}

	return vm.retval, nil
}

// executeFunction executes a function from the registry
func (vm *VM) executeFunction(name string, args ...interface{}) (interface{}, error) {
	// Check if it's a script-defined function
	if scriptFunc, exists := vm.scriptFunctions[name]; exists {
		return vm.executeScriptFunction(scriptFunc, args...)
	}

	// Check if function exists in registry
	if fn, exists := vm.functionRegistry[name]; exists {
		return fn(args...)
	}

	// 检查是否为模块函数调用 (格式: moduleName.functionName)
	if parts := strings.Split(name, "."); len(parts) == 2 {
		moduleName := parts[0]
		functionName := parts[1]

		// 通过模块管理器按需调用模块函数
		if vm.moduleManager != nil {
			fmt.Printf("Calling module function: %s.%s with args: %v\n", moduleName, functionName, args)
			if result, err := vm.moduleManager.CallModuleFunction(moduleName, functionName, args...); err == nil {
				fmt.Printf("Successfully called module function: %s.%s, result: %v\n", moduleName, functionName, result)
				return result, nil
			} else {
				fmt.Printf("Failed to call module function: %s.%s, error: %v\n", moduleName, functionName, err)
				return nil, err
			}
		}
	}

	// If not found in registry, return error
	return nil, fmt.Errorf("undefined function: %s", name)
}

// executeScriptFunction executes a script-defined function
func (vm *VM) executeScriptFunction(scriptFunc *ScriptFunction, args ...interface{}) (interface{}, error) {
	// Check argument count
	if len(args) != scriptFunc.ParamCount {
		return nil, fmt.Errorf("function %s expects %d arguments, got %d", scriptFunc.Name, scriptFunc.ParamCount, len(args))
	}

	// Save current execution state
	currentIP := vm.ip
	currentLocals := vm.locals
	currentExecutionCount := vm.executionCount // 保存当前执行计数

	// Create new local scope for function execution
	vm.locals = make(map[string]interface{})

	// Set up function parameters as local variables
	// For script-defined functions, we need to handle parameters correctly
	// Push arguments in normal order (not reverse order) so that STORE_NAME instructions
	// can pop them in the correct order
	for i := 0; i < len(args); i++ {
		// For value receivers, we need to create a copy of the struct
		// For pointer receivers, we pass the struct as-is
		if i == 0 && scriptFunc.ParamCount > 0 && scriptFunc.ReceiverType != "" {
			// This is the receiver parameter for a method
			// Handle receiver parameter based on receiver type
			if objMap, ok := args[i].(map[string]interface{}); ok && scriptFunc.ReceiverType == "value" {
				// Create a deep copy of the map for value receiver
				copyMap := deepCopyMap(objMap)
				vm.Push(copyMap)
			} else if objMap, ok := args[i].(map[string]interface{}); ok && scriptFunc.ReceiverType == "pointer" {
				// For pointer receivers, pass the original struct (no copy)
				vm.Push(objMap)
			} else {
				vm.Push(args[i])
			}
		} else {
			vm.Push(args[i])
		}
	}

	// Execute function instructions
	startIP := scriptFunc.StartIP
	endIP := scriptFunc.EndIP

	// Execute instructions in the function
	for vm.ip = startIP; vm.ip < endIP && vm.ip < len(vm.instructions); vm.ip++ {
		// Check if maximum instruction limit is exceeded
		if vm.maxInstructions > 0 && int64(vm.executionCount) >= vm.maxInstructions {
			// Restore execution state before returning error
			vm.ip = currentIP
			vm.locals = currentLocals
			vm.executionCount = currentExecutionCount
			return nil, fmt.Errorf("maximum instruction limit exceeded: %d instructions executed", vm.executionCount)
		}

		instr := vm.instructions[vm.ip]

		// Debug output
		if vm.debug {
			fmt.Printf("Function IP: %d, Instruction: %s, Stack: %v\n", vm.ip, instr.String(), vm.stack)
		}

		switch instr.Op {
		case OpNop:
			// Do nothing
		case OpLoadConst:
			vm.Push(instr.Arg)
		case OpLoadName:
			name := instr.Arg.(string)
			if value, ok := vm.GetLocal(name); ok {
				vm.Push(value)
			} else if value, ok := vm.GetGlobal(name); ok {
				vm.Push(value)
			} else {
				// Restore execution state before returning error
				vm.ip = currentIP
				vm.locals = currentLocals
				vm.executionCount = currentExecutionCount
				return nil, fmt.Errorf("undefined variable: %s", name)
			}
		case OpStoreName:
			name := instr.Arg.(string)
			if len(vm.stack) == 0 {
				// Restore execution state before returning error
				vm.ip = currentIP
				vm.locals = currentLocals
				vm.executionCount = currentExecutionCount
				return nil, fmt.Errorf("stack underflow in STORE_NAME")
			}
			value := vm.Pop()
			vm.SetLocal(name, value)
		case OpCall:
			// Function call using unified function registry
			fnName := instr.Arg.(string)
			argCount := instr.Arg2.(int)

			// Check if we have enough arguments on the stack
			if len(vm.stack) < argCount {
				// Restore execution state before returning error
				vm.ip = currentIP
				vm.locals = currentLocals
				vm.executionCount = currentExecutionCount
				return nil, fmt.Errorf("stack underflow in CALL: expected %d arguments, got %d", argCount, len(vm.stack))
			}

			// Pop arguments from stack (in reverse order to maintain correct parameter order)
			callArgs := make([]interface{}, argCount)
			for i := argCount - 1; i >= 0; i-- {
				callArgs[i] = vm.Pop()
			}

			// For script-defined functions, check if the receiver parameter needs to be handled specially
			if scriptFunc, exists := vm.scriptFunctions[fnName]; exists {
				// Handle receiver parameter based on receiver type
				if scriptFunc.ReceiverType == "value" && len(callArgs) > 0 {
					// If it's a map (struct) and the receiver type is "value", create a copy
					if objMap, ok := callArgs[0].(map[string]interface{}); ok {
						callArgs[0] = deepCopyMap(objMap)
					}
				}
				// For pointer receivers, we pass the struct as-is (no copy needed)
			}

			// Execute function from registry
			result, err := vm.executeFunction(fnName, callArgs...)
			if err != nil {
				// Restore execution state before returning error
				vm.ip = currentIP
				vm.locals = currentLocals
				vm.executionCount = currentExecutionCount
				return nil, err
			}

			vm.Push(result)
		case OpReturn:
			if len(vm.stack) == 0 {
				// Restore execution state before returning error
				vm.ip = currentIP
				vm.locals = currentLocals
				vm.executionCount = currentExecutionCount
				return nil, fmt.Errorf("stack underflow in RETURN")
			}
			result := vm.Pop()

			// For pointer receiver methods, we need to return the modified receiver
			if scriptFunc.ReceiverType == "pointer" && len(args) > 0 {
				// Return the modified receiver (first argument)
				result = args[0]
			}

			// Restore execution state
			vm.ip = currentIP
			vm.locals = currentLocals

			return result, nil
		case OpJump:
			target := instr.Arg.(int)
			if target < 0 || target >= len(vm.instructions) {
				// Restore execution state before returning error
				vm.ip = currentIP
				vm.locals = currentLocals
				vm.executionCount = currentExecutionCount
				return nil, fmt.Errorf("invalid jump target: %d", target)
			}
			vm.ip = target - 1 // -1 because we increment at the end of the loop
		case OpJumpIf:
			if len(vm.stack) == 0 {
				// Restore execution state before returning error
				vm.ip = currentIP
				vm.locals = currentLocals
				vm.executionCount = currentExecutionCount
				return nil, fmt.Errorf("stack underflow in JUMP_IF")
			}
			condition := vm.Pop()
			// Jump if condition is FALSE (negate the condition)
			if !isTruthy(condition) {
				target := instr.Arg.(int)
				if target < 0 || target >= len(vm.instructions) {
					// Restore execution state before returning error
					vm.ip = currentIP
					vm.locals = currentLocals
					vm.executionCount = currentExecutionCount
					return nil, fmt.Errorf("invalid jump target: %d", target)
				}
				vm.ip = target - 1 // -1 because we increment at the end of the loop
			}
		case OpBinaryOp:
			if len(vm.stack) < 2 {
				// Restore execution state before returning error
				vm.ip = currentIP
				vm.locals = currentLocals
				vm.executionCount = currentExecutionCount
				return nil, fmt.Errorf("stack underflow in BINARY_OP: expected 2 values, got %d", len(vm.stack))
			}
			right := vm.Pop()
			left := vm.Pop()

			result, err := vm.executeBinaryOp(instr.Arg.(BinaryOp), left, right)
			if err != nil {
				// Restore execution state before returning error
				vm.ip = currentIP
				vm.locals = currentLocals
				vm.executionCount = currentExecutionCount
				return nil, err
			}

			vm.Push(result)
		case OpUnaryOp:
			if len(vm.stack) < 1 {
				// Restore execution state before returning error
				vm.ip = currentIP
				vm.locals = currentLocals
				vm.executionCount = currentExecutionCount
				return nil, fmt.Errorf("stack underflow in UNARY_OP: expected 1 value, got %d", len(vm.stack))
			}
			value := vm.Pop()

			result, err := vm.executeUnaryOp(instr.Arg.(UnaryOp), value)
			if err != nil {
				// Restore execution state before returning error
				vm.ip = currentIP
				vm.locals = currentLocals
				vm.executionCount = currentExecutionCount
				return nil, err
			}

			vm.Push(result)
		case OpNewStruct:
			// Create a new struct instance based on type definition
			var structInstance map[string]interface{}

			// If a type name is provided, try to create a struct with default values
			if typeName, ok := instr.Arg.(string); ok && typeName != "" {
				if structType, exists := vm.typeSystem[typeName]; exists {
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

			vm.Push(structInstance)
		case OpGetField:
			if len(vm.stack) < 2 {
				// Restore execution state before returning error
				vm.ip = currentIP
				vm.locals = currentLocals
				vm.executionCount = currentExecutionCount
				return nil, fmt.Errorf("stack underflow in GET_FIELD: expected 2 values, got %d", len(vm.stack))
			}
			fieldName := vm.Pop()
			obj := vm.Pop()

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
				vm.Push(field)
			} else {
				// If obj is not a map, return nil
				vm.Push(nil)
			}
		case OpSetField:
			// Debug output
			if vm.debug {
				fmt.Printf("Function IP: %d, About to execute SET_FIELD, Stack size: %d, Stack: %v\n", vm.ip, len(vm.stack), vm.stack)
			}

			if len(vm.stack) < 3 {
				// Restore execution state before returning error
				vm.ip = currentIP
				vm.locals = currentLocals
				vm.executionCount = currentExecutionCount
				return nil, fmt.Errorf("stack underflow in SET_FIELD: expected 3 values, got %d", len(vm.stack))
			}
			fieldName := vm.Pop()
			value := vm.Pop()
			obj := vm.Pop()

			// Convert fieldName to string if it's not already
			var fieldNameStr string
			if s, ok := fieldName.(string); ok {
				fieldNameStr = s
			} else {
				fieldNameStr = fmt.Sprintf("%v", fieldName)
			}

			// For now, we assume obj is a map
			if objMap, ok := obj.(map[string]interface{}); ok {
				// Debug output
				if vm.debug {
					fmt.Printf("Function IP: %d, Setting field '%s' to '%v' in object: %v\n", vm.ip, fieldNameStr, value, objMap)
				}
				objMap[fieldNameStr] = value
				// Debug output
				if vm.debug {
					fmt.Printf("Function IP: %d, Object after setting field: %v\n", vm.ip, objMap)
				}
				// Push the modified object back onto the stack
				vm.Push(objMap)
			} else {
				// Debug output
				if vm.debug {
					fmt.Printf("Function IP: %d, Object is not a map, type: %T, value: %v, pushing back unchanged\n", vm.ip, obj, obj)
				}
				// If obj is not a map, push it back unchanged
				vm.Push(obj)
			}
		case OpSetStructField:
			// Debug output
			if vm.debug {
				fmt.Printf("Function IP: %d, About to execute SET_STRUCT_FIELD, Stack size: %d, Stack: %v\n", vm.ip, len(vm.stack), vm.stack)
			}

			if len(vm.stack) < 3 {
				// Restore execution state before returning error
				vm.ip = currentIP
				vm.locals = currentLocals
				vm.executionCount = currentExecutionCount
				return nil, fmt.Errorf("stack underflow in SET_STRUCT_FIELD: expected 3 values, got %d", len(vm.stack))
			}
			// Pop in reverse order to get the correct values
			// Stack order should be: [object, fieldName, value]
			value := vm.Pop()
			fieldName := vm.Pop()
			obj := vm.Pop()

			// Convert fieldName to string if it's not already
			var fieldNameStr string
			if s, ok := fieldName.(string); ok {
				fieldNameStr = s
			} else {
				fieldNameStr = fmt.Sprintf("%v", fieldName)
			}

			// For now, we assume obj is a map
			if objMap, ok := obj.(map[string]interface{}); ok {
				// Debug output
				if vm.debug {
					fmt.Printf("Function IP: %d, Setting struct field '%s' to '%v' in object: %v\n", vm.ip, fieldNameStr, value, objMap)
				}
				objMap[fieldNameStr] = value
				// Debug output
				if vm.debug {
					fmt.Printf("Function IP: %d, Object after setting struct field: %v\n", vm.ip, objMap)
				}
				// Push the modified object back onto the stack
				vm.Push(objMap)
			} else {
				// Debug output
				if vm.debug {
					fmt.Printf("Function IP: %d, Object is not a map, type: %T, value: %v, pushing back unchanged\n", vm.ip, obj, obj)
				}
				// If obj is not a map, push it back unchanged
				vm.Push(obj)
			}
		case OpGetIndex:
			if len(vm.stack) < 2 {
				// Restore execution state before returning error
				vm.ip = currentIP
				vm.locals = currentLocals
				vm.executionCount = currentExecutionCount
				return nil, fmt.Errorf("stack underflow in GET_INDEX: expected 2 values, got %d", len(vm.stack))
			}
			index := vm.Pop()
			array := vm.Pop()

			// Convert index to int if it's not already
			var indexInt int
			if i, ok := index.(int); ok {
				indexInt = i
			} else {
				// Restore execution state before returning error
				vm.ip = currentIP
				vm.locals = currentLocals
				vm.executionCount = currentExecutionCount
				return nil, fmt.Errorf("index must be an integer, got %T", index)
			}

			// For now, we assume array is a slice of interface{}
			if arraySlice, ok := array.([]interface{}); ok {
				if indexInt >= 0 && indexInt < len(arraySlice) {
					vm.Push(arraySlice[indexInt])
				} else {
					vm.Push(nil)
				}
			} else {
				// If array is not a slice, return nil
				vm.Push(nil)
			}
		case OpSetIndex:
			if len(vm.stack) < 3 {
				// Restore execution state before returning error
				vm.ip = currentIP
				vm.locals = currentLocals
				vm.executionCount = currentExecutionCount
				return nil, fmt.Errorf("stack underflow in SET_INDEX: expected 3 values, got %d", len(vm.stack))
			}
			value := vm.Pop()
			index := vm.Pop()
			array := vm.Pop()

			// Convert index to int if it's not already
			var indexInt int
			if i, ok := index.(int); ok {
				indexInt = i
			} else {
				// Restore execution state before returning error
				vm.ip = currentIP
				vm.locals = currentLocals
				vm.executionCount = currentExecutionCount
				return nil, fmt.Errorf("index must be an integer, got %T", index)
			}

			// For now, we assume array is a slice of interface{}
			if arraySlice, ok := array.([]interface{}); ok {
				if indexInt >= 0 && indexInt < len(arraySlice) {
					arraySlice[indexInt] = value
					// Push the modified array back onto the stack
					vm.Push(arraySlice)
				} else {
					// Index out of bounds, push the array back unchanged
					vm.Push(arraySlice)
				}
			} else {
				// If array is not a slice, push it back unchanged
				vm.Push(array)
			}
		case OpRotate:
			if len(vm.stack) < 3 {
				// Restore execution state before returning error
				vm.ip = currentIP
				vm.locals = currentLocals
				vm.executionCount = currentExecutionCount
				return nil, fmt.Errorf("stack underflow in ROTATE: expected 3 values, got %d", len(vm.stack))
			}
			// Rotate the top three elements: [a, b, c] -> [b, c, a]
			top := vm.stack[len(vm.stack)-1]
			second := vm.stack[len(vm.stack)-2]
			third := vm.stack[len(vm.stack)-3]
			vm.stack[len(vm.stack)-1] = third
			vm.stack[len(vm.stack)-2] = top
			vm.stack[len(vm.stack)-3] = second
		}

		vm.executionCount++
	}

	// Restore execution state
	vm.ip = currentIP
	vm.locals = currentLocals

	// If we reach here, the function didn't return explicitly
	return nil, nil
}

// deepCopyMap creates a deep copy of a map
func deepCopyMap(original map[string]interface{}) map[string]interface{} {
	copyMap := make(map[string]interface{})
	for k, v := range original {
		// For nested maps, we need to deep copy them as well
		if nestedMap, ok := v.(map[string]interface{}); ok {
			copyMap[k] = deepCopyMap(nestedMap)
		} else {
			copyMap[k] = v
		}
	}
	return copyMap
}

// SetFunctionRegistry sets the function registry
func (vm *VM) SetFunctionRegistry(registry map[string]func(args ...interface{}) (interface{}, error)) {
	vm.functionRegistry = registry
}

// GetFunctionRegistry gets the function registry
func (vm *VM) GetFunctionRegistry() map[string]func(args ...interface{}) (interface{}, error) {
	return vm.functionRegistry
}

// RegisterType registers a type in the VM
func (vm *VM) RegisterType(name string, typ types.IType) {
	vm.typeSystem[name] = typ
}

// GetType gets a type from the VM
func (vm *VM) GetType(name string) (types.IType, bool) {
	typ, ok := vm.typeSystem[name]
	return typ, ok
}

// SetModuleManager sets the module manager interface for the VM
func (vm *VM) SetModuleManager(mm ModuleManagerInterface) {
	vm.moduleManager = mm
}

// executeBinaryOp executes a binary operation
func (vm *VM) executeBinaryOp(op BinaryOp, left, right interface{}) (interface{}, error) {
	switch op {
	case OpAdd:
		return vm.add(left, right)
	case OpSub:
		return vm.sub(left, right)
	case OpMul:
		return vm.mul(left, right)
	case OpDiv:
		return vm.div(left, right)
	case OpMod:
		return vm.mod(left, right)
	case OpEqual:
		return reflect.DeepEqual(left, right), nil
	case OpNotEqual:
		return !reflect.DeepEqual(left, right), nil
	case OpLess:
		return vm.less(left, right)
	case OpLessEqual:
		return vm.lessEqual(left, right)
	case OpGreater:
		return vm.greater(left, right)
	case OpGreaterEqual:
		return vm.greaterEqual(left, right)
	case OpAnd:
		return isTruthy(left) && isTruthy(right), nil
	case OpOr:
		return isTruthy(left) || isTruthy(right), nil
	default:
		return nil, fmt.Errorf("unknown binary operation: %d", op)
	}
}

// executeUnaryOp executes a unary operation
func (vm *VM) executeUnaryOp(op UnaryOp, value interface{}) (interface{}, error) {
	switch op {
	case OpNeg:
		switch v := value.(type) {
		case int:
			return -v, nil
		case float64:
			return -v, nil
		default:
			return nil, fmt.Errorf("unsupported type for negation: %T", value)
		}
	case OpNot:
		return !isTruthy(value), nil
	default:
		return nil, fmt.Errorf("unknown unary operation: %d", op)
	}
}

// Helper functions for arithmetic operations
func (vm *VM) add(left, right interface{}) (interface{}, error) {
	switch l := left.(type) {
	case int:
		switch r := right.(type) {
		case int:
			return l + r, nil
		case float64:
			return float64(l) + r, nil
		}
	case float64:
		switch r := right.(type) {
		case int:
			return l + float64(r), nil
		case float64:
			return l + r, nil
		}
	case string:
		if r, ok := right.(string); ok {
			return l + r, nil
		}
	}

	return nil, fmt.Errorf("unsupported types for addition: %T and %T", left, right)
}

func (vm *VM) sub(left, right interface{}) (interface{}, error) {
	switch l := left.(type) {
	case int:
		switch r := right.(type) {
		case int:
			return l - r, nil
		case float64:
			return float64(l) - r, nil
		}
	case float64:
		switch r := right.(type) {
		case int:
			return l - float64(r), nil
		case float64:
			return l - r, nil
		}
	}

	return nil, fmt.Errorf("unsupported types for subtraction: %T and %T", left, right)
}

func (vm *VM) mul(left, right interface{}) (interface{}, error) {
	switch l := left.(type) {
	case int:
		switch r := right.(type) {
		case int:
			return l * r, nil
		case float64:
			return float64(l) * r, nil
		}
	case float64:
		switch r := right.(type) {
		case int:
			return l * float64(r), nil
		case float64:
			return l * r, nil
		}
	}

	return nil, fmt.Errorf("unsupported types for multiplication: %T and %T", left, right)
}

func (vm *VM) div(left, right interface{}) (interface{}, error) {
	switch l := left.(type) {
	case int:
		switch r := right.(type) {
		case int:
			if r == 0 {
				return nil, fmt.Errorf("division by zero")
			}
			return l / r, nil
		case float64:
			if r == 0.0 {
				return nil, fmt.Errorf("division by zero")
			}
			return float64(l) / r, nil
		}
	case float64:
		switch r := right.(type) {
		case int:
			if r == 0 {
				return nil, fmt.Errorf("division by zero")
			}
			return l / float64(r), nil
		case float64:
			if r == 0.0 {
				return nil, fmt.Errorf("division by zero")
			}
			return l / r, nil
		}
	}

	return nil, fmt.Errorf("unsupported types for division: %T and %T", left, right)
}

func (vm *VM) mod(left, right interface{}) (interface{}, error) {
	if l, ok := left.(int); ok {
		if r, ok := right.(int); ok {
			if r == 0 {
				return nil, fmt.Errorf("modulo by zero")
			}
			return l % r, nil
		}
	}

	return nil, fmt.Errorf("unsupported types for modulo: %T and %T", left, right)
}

func (vm *VM) less(left, right interface{}) (interface{}, error) {
	switch l := left.(type) {
	case int:
		switch r := right.(type) {
		case int:
			return l < r, nil
		case float64:
			return float64(l) < r, nil
		}
	case float64:
		switch r := right.(type) {
		case int:
			return l < float64(r), nil
		case float64:
			return l < r, nil
		}
	case string:
		if r, ok := right.(string); ok {
			return l < r, nil
		}
	}

	return nil, fmt.Errorf("unsupported types for less than comparison: %T and %T", left, right)
}

func (vm *VM) lessEqual(left, right interface{}) (interface{}, error) {
	switch l := left.(type) {
	case int:
		switch r := right.(type) {
		case int:
			return l <= r, nil
		case float64:
			return float64(l) <= r, nil
		}
	case float64:
		switch r := right.(type) {
		case int:
			return l <= float64(r), nil
		case float64:
			return l <= r, nil
		}
	case string:
		if r, ok := right.(string); ok {
			return l <= r, nil
		}
	}

	return nil, fmt.Errorf("unsupported types for less than or equal comparison: %T and %T", left, right)
}

func (vm *VM) greater(left, right interface{}) (interface{}, error) {
	switch l := left.(type) {
	case int:
		switch r := right.(type) {
		case int:
			return l > r, nil
		case float64:
			return float64(l) > r, nil
		}
	case float64:
		switch r := right.(type) {
		case int:
			return l > float64(r), nil
		case float64:
			return l > r, nil
		}
	case string:
		if r, ok := right.(string); ok {
			return l > r, nil
		}
	}

	return nil, fmt.Errorf("unsupported types for greater than comparison: %T and %T", left, right)
}

func (vm *VM) greaterEqual(left, right interface{}) (interface{}, error) {
	switch l := left.(type) {
	case int:
		switch r := right.(type) {
		case int:
			return l >= r, nil
		case float64:
			return float64(l) >= r, nil
		}
	case float64:
		switch r := right.(type) {
		case int:
			return l >= float64(r), nil
		case float64:
			return l >= r, nil
		}
	case string:
		if r, ok := right.(string); ok {
			return l >= r, nil
		}
	}

	return nil, fmt.Errorf("unsupported types for greater than or equal comparison: %T and %T", left, right)
}

// isTruthy determines if a value is truthy
func isTruthy(value interface{}) bool {
	switch v := value.(type) {
	case nil:
		return false
	case bool:
		return v
	case int:
		return v != 0
	case float64:
		return v != 0.0
	case string:
		return v != ""
	default:
		return true
	}
}

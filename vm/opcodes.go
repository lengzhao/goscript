// Package vm provides the virtual machine implementation
package vm

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/lengzhao/goscript/context"
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

	// Swap the top two elements on the stack
	OpSwap = instruction.OpSwap

	// Create a new slice
	OpNewSlice = instruction.OpNewSlice

	// Get the length of a slice or array
	OpLen = instruction.OpLen

	// Get an element from a slice or array by index
	OpGetElement = instruction.OpGetElement

	// Import a module
	OpImport = instruction.OpImport

	// Enter a scope with a specific key
	OpEnterScopeWithKey = instruction.OpEnterScopeWithKey

	// Exit a scope with a specific key
	OpExitScopeWithKey = instruction.OpExitScopeWithKey

	// Create a new variable
	OpCreateVar = instruction.OpCreateVar
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

	// Create global context for new approach
	globalCtx := context.NewContext("global", nil)

	vm := &VM{
		stack:            NewStack(64, 10000), // 初始容量64，最大10000
		instructions:     make([]*Instruction, 0),
		ip:               0,
		functionRegistry: make(map[string]func(args ...interface{}) (interface{}, error)),
		scriptFunctions:  make(map[string]*ScriptFunction),
		typeSystem:       typeSystem,
		debug:            false,
		executionCount:   0,
		maxInstructions:  1000000, // 默认最大指令数限制为100万条指令
		moduleManager:    nil,     // 初始化为空
		handlers:         make(map[OpCode]OpHandler),
		globalCtx:        globalCtx,
		currentCtx:       globalCtx,
	}

	// Register all opcode handlers
	vm.registerHandlers()

	return vm
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
	if err := vm.stack.Push(value); err != nil {
		// In case of stack overflow, we could panic or handle gracefully
		// For now, we'll panic as this indicates a serious issue
		panic(fmt.Sprintf("Stack overflow: %v", err))
	}
}

// Pop pops a value from the stack
func (vm *VM) Pop() interface{} {
	value, err := vm.stack.Pop()
	if err != nil {
		// Return nil for stack underflow to maintain backward compatibility
		return nil
	}
	return value
}

// Peek returns the top value without removing it
func (vm *VM) Peek() interface{} {
	value, err := vm.stack.Peek()
	if err != nil {
		// Return nil for stack underflow to maintain backward compatibility
		return nil
	}
	return value
}

// StackSize returns the current stack size
func (vm *VM) StackSize() int {
	return vm.stack.Size()
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
	vm.stack.Clear()
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

	// Create a new context for this function execution with a unique path key
	// Use the function name as the path key to ensure uniqueness
	funcCtx := context.NewContext("function."+scriptFunc.Name, vm.currentCtx)

	// Add the new context as a child of the current context
	if vm.currentCtx != nil {
		vm.currentCtx.AddChild(funcCtx)
	}

	// Set the new context as current context
	vm.currentCtx = funcCtx

	// Push arguments onto the stack in the order they were passed
	// The function body will pop them in reverse order (last parameter first)
	for _, arg := range args {
		vm.Push(arg)
	}

	// Use the Run method to execute the function instructions
	// This ensures proper context management and error handling
	result, err := vm.Run(funcCtx, scriptFunc.StartIP, scriptFunc.EndIP)

	// Clean up: remove the function context from parent's children
	if vm.currentCtx != nil && vm.currentCtx.GetParent() != nil {
		parent := vm.currentCtx.GetParent()
		parent.RemoveChild(funcCtx.GetPathKey())
		// Restore parent context
		vm.currentCtx = parent
	}

	return result, err
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

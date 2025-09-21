// Package vm provides the virtual machine implementation
package vm

import (
	"fmt"
	"reflect"
)

// OpCode represents an operation code for the virtual machine
type OpCode byte

// String returns the string representation of an OpCode
func (op OpCode) String() string {
	switch op {
	case OpNop:
		return "OpNop"
	case OpLoadConst:
		return "OpLoadConst"
	case OpLoadName:
		return "OpLoadName"
	case OpStoreName:
		return "OpStoreName"
	case OpCall:
		return "OpCall"
	case OpRegistFunction:
		return "OpRegistFunction"
	case OpReturn:
		return "OpReturn"
	case OpJump:
		return "OpJump"
	case OpJumpIf:
		return "OpJumpIf"
	case OpBinaryOp:
		return "OpBinaryOp"
	case OpUnaryOp:
		return "OpUnaryOp"
	default:
		return fmt.Sprintf("OpCode(%d)", op)
	}
}

const (
	// No operation
	OpNop OpCode = iota

	// Load a constant onto the stack
	OpLoadConst

	// Load a variable by name
	OpLoadName

	// Store a value to a variable by name
	OpStoreName

	// Call a function
	OpCall

	// Register a script-defined function
	OpRegistFunction

	// Return from function
	OpReturn

	// Unconditional jump
	OpJump

	// Conditional jump
	OpJumpIf

	// Binary operation (add, sub, mul, div, etc.)
	OpBinaryOp

	// Unary operation (neg, not, etc.)
	OpUnaryOp
)

// BinaryOp represents a binary operation
type BinaryOp byte

const (
	OpAdd BinaryOp = iota
	OpSub
	OpMul
	OpDiv
	OpMod
	OpEqual
	OpNotEqual
	OpLess
	OpLessEqual
	OpGreater
	OpGreaterEqual
	OpAnd
	OpOr
)

// UnaryOp represents a unary operation
type UnaryOp byte

const (
	OpNeg UnaryOp = iota
	OpNot
)

// Instruction represents a single VM instruction
type Instruction struct {
	Op   OpCode
	Arg  interface{}
	Arg2 interface{}
}

// NewInstruction creates a new instruction
func NewInstruction(op OpCode, arg interface{}, arg2 ...interface{}) *Instruction {
	instr := &Instruction{
		Op:  op,
		Arg: arg,
	}

	if len(arg2) > 0 {
		instr.Arg2 = arg2[0]
	}

	return instr
}

// String returns the string representation of an instruction
func (i *Instruction) String() string {
	switch i.Op {
	case OpLoadConst:
		return fmt.Sprintf("LOAD_CONST %v", i.Arg)
	case OpLoadName:
		return fmt.Sprintf("LOAD_NAME %v", i.Arg)
	case OpStoreName:
		return fmt.Sprintf("STORE_NAME %v", i.Arg)
	case OpCall:
		return fmt.Sprintf("CALL %v %v", i.Arg, i.Arg2)
	case OpRegistFunction:
		return fmt.Sprintf("REGIST_FUNCTION %v %v", i.Arg, i.Arg2)
	case OpReturn:
		return "RETURN"
	case OpJump:
		return fmt.Sprintf("JUMP %v", i.Arg)
	case OpJumpIf:
		return fmt.Sprintf("JUMP_IF %v", i.Arg)
	case OpBinaryOp:
		return fmt.Sprintf("BINARY_OP %v", i.Arg)
	case OpUnaryOp:
		return fmt.Sprintf("UNARY_OP %v", i.Arg)
	default:
		return fmt.Sprintf("UNKNOWN(%d) %v %v", i.Op, i.Arg, i.Arg2)
	}
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

	// Debug mode
	debug bool

	// Execution statistics
	executionCount int
}

// ScriptFunction represents a script-defined function
type ScriptFunction struct {
	Name       string
	StartIP    int
	EndIP      int
	ParamCount int
	ParamNames []string // Store parameter names for use in function body
}

// NewVM creates a new virtual machine
func NewVM() *VM {
	return &VM{
		stack:            make([]interface{}, 0),
		globals:          make(map[string]interface{}),
		locals:           make(map[string]interface{}),
		instructions:     make([]*Instruction, 0),
		ip:               0,
		functionRegistry: make(map[string]func(args ...interface{}) (interface{}, error)),
		scriptFunctions:  make(map[string]*ScriptFunction),
		debug:            false,
		executionCount:   0,
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

// Execute executes the instructions
func (vm *VM) Execute() (interface{}, error) {
	vm.ip = 0
	vm.executionCount = 0

	for vm.ip < len(vm.instructions) {
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
		case OpCall:
			// Function call using unified function registry
			fnName := instr.Arg.(string)
			argCount := instr.Arg2.(int)

			// Check if we have enough arguments on the stack
			if len(vm.stack) < argCount {
				return nil, fmt.Errorf("stack underflow in CALL: expected %d arguments, got %d", argCount, len(vm.stack))
			}

			// Pop arguments from stack
			args := make([]interface{}, argCount)
			for i := argCount - 1; i >= 0; i-- {
				args[i] = vm.Pop()
			}

			// Execute function from registry
			result, err := vm.executeFunction(fnName, args...)
			if err != nil {
				return nil, err
			}

			vm.Push(result)
		case OpRegistFunction:
			// Register a script-defined function
			fnName := instr.Arg.(string)
			funcInfo := instr.Arg2.(*ScriptFunction)
			vm.scriptFunctions[fnName] = funcInfo
		case OpReturn:
			if len(vm.stack) == 0 {
				return nil, fmt.Errorf("stack underflow in RETURN")
			}
			vm.retval = vm.Pop()
			return vm.retval, nil
		case OpJump:
			target := instr.Arg.(int)
			if target < 0 || target >= len(vm.instructions) {
				return nil, fmt.Errorf("invalid jump target: %d", target)
			}
			vm.ip = target
			vm.executionCount++
			continue
		case OpJumpIf:
			if len(vm.stack) == 0 {
				return nil, fmt.Errorf("stack underflow in JUMP_IF")
			}
			condition := vm.Pop()
			if isTruthy(condition) {
				target := instr.Arg.(int)
				if target < 0 || target >= len(vm.instructions) {
					return nil, fmt.Errorf("invalid jump target: %d", target)
				}
				vm.ip = target
				vm.executionCount++
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

	// Create new local scope for function execution
	vm.locals = make(map[string]interface{})

	// Set up function parameters as local variables
	// Push parameters onto the stack in reverse order for the STORE_NAME instructions
	// The function body will pop them off in the correct order
	if len(scriptFunc.ParamNames) > 0 && len(scriptFunc.ParamNames) == len(args) {
		// Push parameters onto stack in reverse order
		// This is because STORE_NAME pops from the stack
		for i := len(args) - 1; i >= 0; i-- {
			vm.Push(args[i])
		}
	}

	// Execute function instructions
	startIP := scriptFunc.StartIP
	endIP := scriptFunc.EndIP

	// Execute instructions in the function
	for vm.ip = startIP; vm.ip < endIP && vm.ip < len(vm.instructions); vm.ip++ {
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
				return nil, fmt.Errorf("undefined variable: %s", name)
			}
		case OpStoreName:
			name := instr.Arg.(string)
			if len(vm.stack) == 0 {
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
				return nil, fmt.Errorf("stack underflow in CALL: expected %d arguments, got %d", argCount, len(vm.stack))
			}

			// Pop arguments from stack
			callArgs := make([]interface{}, argCount)
			for i := argCount - 1; i >= 0; i-- {
				callArgs[i] = vm.Pop()
			}

			// Execute function from registry
			result, err := vm.executeFunction(fnName, callArgs...)
			if err != nil {
				return nil, err
			}

			vm.Push(result)
		case OpReturn:
			if len(vm.stack) == 0 {
				return nil, fmt.Errorf("stack underflow in RETURN")
			}
			result := vm.Pop()

			// Restore execution state
			vm.ip = currentIP
			vm.locals = currentLocals

			return result, nil
		case OpJump:
			target := instr.Arg.(int)
			if target < 0 || target >= len(vm.instructions) {
				return nil, fmt.Errorf("invalid jump target: %d", target)
			}
			vm.ip = target - 1 // -1 because we increment at the end of the loop
		case OpJumpIf:
			if len(vm.stack) == 0 {
				return nil, fmt.Errorf("stack underflow in JUMP_IF")
			}
			condition := vm.Pop()
			if isTruthy(condition) {
				target := instr.Arg.(int)
				if target < 0 || target >= len(vm.instructions) {
					return nil, fmt.Errorf("invalid jump target: %d", target)
				}
				vm.ip = target - 1 // -1 because we increment at the end of the loop
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
		}
	}

	// Restore execution state
	vm.ip = currentIP
	vm.locals = currentLocals

	// If we reach here, the function didn't return explicitly
	return nil, nil
}

// SetFunctionRegistry sets the function registry
func (vm *VM) SetFunctionRegistry(registry map[string]func(args ...interface{}) (interface{}, error)) {
	vm.functionRegistry = registry
}

// GetFunctionRegistry gets the function registry
func (vm *VM) GetFunctionRegistry() map[string]func(args ...interface{}) (interface{}, error) {
	return vm.functionRegistry
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
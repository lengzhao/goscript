// Package vm implements the virtual machine for GoScript
package vm

import (
	"fmt"
	"strings"
	"sync"

	"github.com/lengzhao/goscript/context"
	"github.com/lengzhao/goscript/instruction"
	"github.com/lengzhao/goscript/types"
)

// VM represents the GoScript virtual machine
type VM struct {
	// Instructions organized by key (e.g., "main.main", "main.init")
	InstructionSets map[string][]*instruction.Instruction

	// Global context
	GlobalCtx *context.Context

	// Current context
	currentCtx *context.Context

	// All instructions (for compatibility with compiler tests)
	instructions []*instruction.Instruction

	// Registered functions that can be called from scripts
	functions map[string]ScriptFunction

	// Script function information for parameter names
	scriptFunctionInfos map[string]*ScriptFunctionInfo

	// Registered modules with simplified interface
	modules map[string]types.ModuleExecutor

	// Mutex for thread safety
	mu sync.RWMutex

	// Instruction counter for security limits
	instructionCount int64

	// Maximum number of instructions allowed (0 means no limit)
	maxInstructions int64

	// Debug mode
	debug bool
}

// ScriptFunction represents a function that can be called from scripts
type ScriptFunction func(args ...interface{}) (interface{}, error)

// ScriptFunctionInfo represents information about a script-defined function
type ScriptFunctionInfo struct {
	Name       string
	Key        string
	ParamCount int
	ParamNames []string // Add parameter names
}

// NewVM creates a new virtual machine
func NewVM() *VM {
	vm := &VM{
		InstructionSets:     make(map[string][]*instruction.Instruction),
		functions:           make(map[string]ScriptFunction),
		scriptFunctionInfos: make(map[string]*ScriptFunctionInfo),
		modules:             make(map[string]types.ModuleExecutor),
		instructions:        make([]*instruction.Instruction, 0),
		GlobalCtx:           context.NewContext("global", nil), // Global context with no parent
		maxInstructions:     10000,                             // Default limit of 10,000 instructions
	}
	return vm
}

// RegisterModule registers a module with a simplified interface
func (vm *VM) RegisterModule(name string, executor types.ModuleExecutor) {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	vm.modules[name] = executor
}

// GetModule retrieves a registered module by name
func (vm *VM) GetModule(name string) (types.ModuleExecutor, bool) {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	module, exists := vm.modules[name]
	return module, exists
}

// GetFunction retrieves a registered function by name
// This can be a standalone function or a module function (module.function)
func (vm *VM) GetFunction(name string) (ScriptFunction, bool) {
	vm.mu.RLock()
	defer vm.mu.RUnlock()

	// First check if it's a standalone function
	fn, exists := vm.functions[name]
	if exists {
		return fn, true
	}

	// Check if it's a module function (format: "module.function")
	if idx := strings.Index(name, "."); idx != -1 {
		moduleName := name[:idx]
		entrypoint := name[idx+1:]

		// Check if the module exists
		if module, moduleExists := vm.modules[moduleName]; moduleExists {
			// Create a wrapper function that calls the module executor
			wrapper := func(args ...interface{}) (interface{}, error) {
				return module(entrypoint, args...)
			}
			return wrapper, true
		}
	}

	return nil, false
}

// RegisterFunction registers a function that can be called from scripts
func (vm *VM) RegisterFunction(name string, fn ScriptFunction) {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	vm.functions[name] = fn
}

// RegisterScriptFunction registers a script-defined function
func (vm *VM) RegisterScriptFunction(name string, info *ScriptFunctionInfo) {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	// Store the function info for later use
	vm.scriptFunctionInfos[name] = info

	// Create a wrapper function that will execute the script function when called
	vm.functions[name] = func(args ...interface{}) (interface{}, error) {
		// Get the function instructions
		instructions, exists := vm.InstructionSets[info.Key]
		if !exists {
			return nil, fmt.Errorf("script function %s not found", info.Key)
		}

		functionCtx := context.NewContext(info.Key, vm.currentCtx)

		// Set function arguments as local variables using the actual parameter names
		paramNames := make([]string, len(args))

		// Use the actual parameter names from the function info if available
		if len(info.ParamNames) > 0 {
			// Use the actual parameter names from the function definition
			for i := 0; i < len(args) && i < len(info.ParamNames); i++ {
				paramNames[i] = info.ParamNames[i]
			}
			// Fill in any remaining parameters with default names
			for i := len(info.ParamNames); i < len(args); i++ {
				paramNames[i] = fmt.Sprintf("arg%d", i)
			}
		} else {
			// Fall back to default parameter names
			for i := 0; i < len(args); i++ {
				paramNames[i] = fmt.Sprintf("arg%d", i)
			}
		}

		// Set arguments as local variables with appropriate names
		for i, arg := range args {
			paramName := paramNames[i]
			// Create and set the variable with the actual argument value
			functionCtx.CreateVariableWithType(paramName, arg, "unknown")
		}

		// Save the current context
		previousCtx := vm.currentCtx

		// Set the current context for the function execution
		vm.currentCtx = functionCtx

		// Execute the function instructions using the executor
		executor := NewExecutor(vm)
		result, err := executor.executeInstructions(instructions)

		// Restore the previous context
		vm.currentCtx = previousCtx

		return result, err
	}
}

// GetAllScriptFunctions returns all registered script function information
func (vm *VM) GetAllScriptFunctions() map[string]*ScriptFunctionInfo {
	vm.mu.RLock()
	defer vm.mu.RUnlock()

	// Create a copy of the map to avoid race conditions
	result := make(map[string]*ScriptFunctionInfo)
	for name, info := range vm.scriptFunctionInfos {
		result[name] = info
	}

	return result
}

// GetInstructions returns all instructions (for compatibility with compiler tests)
func (vm *VM) GetInstructions() []*instruction.Instruction {
	vm.mu.RLock()
	defer vm.mu.RUnlock()

	// Collect all instructions from all instruction sets
	var allInstructions []*instruction.Instruction
	for key, instructions := range vm.InstructionSets {
		if vm.debug {
			fmt.Printf("Instructions for key %s:\n", key)
			for i, instr := range instructions {
				fmt.Printf("  %d: %s\n", i, instr.String())
			}
		}
		allInstructions = append(allInstructions, instructions...)
	}

	// Also include instructions added directly
	allInstructions = append(allInstructions, vm.instructions...)

	return allInstructions
}

// AddInstruction adds an instruction to the VM (for compatibility with compiler tests)
func (vm *VM) AddInstruction(instr *instruction.Instruction) {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	vm.instructions = append(vm.instructions, instr)
}

// SetMaxInstructions sets the maximum number of instructions allowed
func (vm *VM) SetMaxInstructions(max int64) {
	vm.maxInstructions = max
}

// GetInstructionCount returns the current instruction count
func (vm *VM) GetInstructionCount() int64 {
	return vm.instructionCount
}

// ResetInstructionCount resets the instruction counter
func (vm *VM) ResetInstructionCount() {
	vm.instructionCount = 0
}

// AddInstructionSet adds a set of instructions with a specific key
func (vm *VM) AddInstructionSet(key string, instructions []*instruction.Instruction) {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	vm.InstructionSets[key] = instructions
}

// GetInstructionSet retrieves instructions by key
func (vm *VM) GetInstructionSet(key string) ([]*instruction.Instruction, bool) {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	instructions, exists := vm.InstructionSets[key]
	return instructions, exists
}

// GetAllInstructionSets returns all instruction sets
func (vm *VM) GetAllInstructionSets() map[string][]*instruction.Instruction {
	vm.mu.RLock()
	defer vm.mu.RUnlock()

	// Create a copy of the map to avoid race conditions
	result := make(map[string][]*instruction.Instruction)
	for key, instructions := range vm.InstructionSets {
		result[key] = instructions
	}

	return result
}

// Execute runs the virtual machine with the given entry point
// If entryPoint is empty, it defaults to "main.main" or tries to find another main function
func (vm *VM) Execute(entryPoint string, args ...interface{}) (interface{}, error) {
	// Reset instruction count before execution
	vm.ResetInstructionCount()

	if entryPoint == "" {
		entryPoint = "main.main"
		// If main.main doesn't exist, try to find another main function
		if _, exists := vm.GetInstructionSet(entryPoint); !exists {
			// Try to find any function ending with ".main"
			for key := range vm.InstructionSets {
				if len(key) >= 5 && key[len(key)-5:] == ".main" {
					entryPoint = key
					break
				}
			}
		}
	}

	// Extract package name from entry point
	packageName := "main" // default
	if idx := len(entryPoint) - 5; idx > 0 {
		if entryPoint[idx:] == ".main" {
			packageName = entryPoint[:idx]
		}
	}

	// Create global context
	globalCtx := context.NewContext("global", nil)
	vm.GlobalCtx = globalCtx

	// Create package context (for main package)
	// The package context's parent is the global context
	packageCtx := context.NewContext(packageName, globalCtx)

	// First, execute package-level code (imports, global variable creation, etc.)
	// This would typically be in the package name itself
	if packageInstructions, exists := vm.GetInstructionSet(packageName); exists {
		vm.currentCtx = packageCtx
		executor := NewExecutor(vm)
		if _, err := executor.executeInstructions(packageInstructions); err != nil {
			return nil, fmt.Errorf("error executing package-level code: %w", err)
		}
	}

	if initInstructions, exists := vm.GetInstructionSet(packageName + ".init"); exists {
		vm.currentCtx = packageCtx
		executor := NewExecutor(vm)
		if _, err := executor.executeInstructions(initInstructions); err != nil {
			return nil, fmt.Errorf("error executing package init: %w", err)
		}
	}

	// Execute the entry point function
	instructions, exists := vm.GetInstructionSet(entryPoint)
	if !exists {
		return nil, fmt.Errorf("entry point %s not found", entryPoint)
	}

	// Create function context with package context as parent
	functionCtx := context.NewContext(entryPoint, packageCtx)
	vm.currentCtx = functionCtx

	// Set function arguments as local variables
	// Check if this is a script function with known parameter names
	paramNames := vm.getScriptFunctionParamNames(entryPoint, len(args))

	// Set arguments as local variables with appropriate names
	for i, arg := range args {
		paramName := paramNames[i]
		functionCtx.CreateVariableWithType(paramName, arg, "unknown")
	}

	// Execute the function using the executor
	executor := NewExecutor(vm)

	result, err := executor.executeInstructions(instructions)

	// Return result and error
	return result, err
}

// getScriptFunctionParamNames gets the parameter names for a script function
// If the function is not a registered script function, it falls back to generic names
func (vm *VM) getScriptFunctionParamNames(functionKey string, argCount int) []string {
	vm.mu.RLock()
	defer vm.mu.RUnlock()

	// Look for the function in script function infos
	for _, info := range vm.scriptFunctionInfos {
		if info.Key == functionKey && len(info.ParamNames) >= argCount {
			return info.ParamNames[:argCount]
		}
	}

	// Fall back to generic parameter names
	paramNames := make([]string, argCount)
	for i := 0; i < argCount; i++ {
		paramNames[i] = fmt.Sprintf("arg%d", i)
	}
	return paramNames
}

// SetDebug enables or disables debug mode
func (vm *VM) SetDebug(debug bool) {
	vm.debug = debug
}

// GetDebug returns the current debug mode
func (vm *VM) GetDebug() bool {
	return vm.debug
}

// executeBinaryOp executes a binary operation
func (vm *VM) executeBinaryOp(op instruction.BinaryOp, left, right interface{}) (interface{}, error) {
	// Debug information
	//fmt.Printf("Executing binary operation: %v with left=%v (type %T) and right=%v (type %T)\n", op, left, left, right, right)

	switch op {
	case instruction.OpAdd:
		// Handle different types of addition
		switch l := left.(type) {
		case int:
			if r, ok := right.(int); ok {
				return l + r, nil
			}
		case float64:
			if r, ok := right.(float64); ok {
				return l + r, nil
			}
		case string:
			if r, ok := right.(string); ok {
				return l + r, nil
			}
		}
		// Handle mixed types for addition
		// Convert int to float64 if one operand is float64
		if l, ok := left.(int); ok {
			if r, ok := right.(float64); ok {
				return float64(l) + r, nil
			}
		}
		if l, ok := left.(float64); ok {
			if r, ok := right.(int); ok {
				return l + float64(r), nil
			}
		}
		return nil, fmt.Errorf("unsupported types for addition: %T and %T", left, right)

	case instruction.OpSub:
		if l, ok := left.(int); ok {
			if r, ok := right.(int); ok {
				return l - r, nil
			}
		}
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l - r, nil
			}
		}
		// Handle mixed types
		if l, ok := left.(int); ok {
			if r, ok := right.(float64); ok {
				return float64(l) - r, nil
			}
		}
		if l, ok := left.(float64); ok {
			if r, ok := right.(int); ok {
				return l - float64(r), nil
			}
		}
		return nil, fmt.Errorf("unsupported types for subtraction: %T and %T", left, right)

	case instruction.OpMul:
		if l, ok := left.(int); ok {
			if r, ok := right.(int); ok {
				return l * r, nil
			}
		}
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l * r, nil
			}
		}
		// Handle mixed types
		if l, ok := left.(int); ok {
			if r, ok := right.(float64); ok {
				return float64(l) * r, nil
			}
		}
		if l, ok := left.(float64); ok {
			if r, ok := right.(int); ok {
				return l * float64(r), nil
			}
		}
		return nil, fmt.Errorf("unsupported types for multiplication: %T and %T", left, right)

	case instruction.OpDiv:
		if l, ok := left.(int); ok {
			if r, ok := right.(int); ok {
				if r == 0 {
					return nil, fmt.Errorf("division by zero")
				}
				return l / r, nil
			}
		}
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				if r == 0.0 {
					return nil, fmt.Errorf("division by zero")
				}
				return l / r, nil
			}
		}
		// Handle mixed types
		if l, ok := left.(int); ok {
			if r, ok := right.(float64); ok {
				if r == 0.0 {
					return nil, fmt.Errorf("division by zero")
				}
				return float64(l) / r, nil
			}
		}
		if l, ok := left.(float64); ok {
			if r, ok := right.(int); ok {
				if r == 0 {
					return nil, fmt.Errorf("division by zero")
				}
				return l / float64(r), nil
			}
		}
		return nil, fmt.Errorf("unsupported types for division: %T and %T", left, right)

	case instruction.OpMod:
		if l, ok := left.(int); ok {
			if r, ok := right.(int); ok {
				if r == 0 {
					return nil, fmt.Errorf("modulo by zero")
				}
				return l % r, nil
			}
		}
		return nil, fmt.Errorf("unsupported types for modulo: %T and %T", left, right)

	case instruction.OpEqual:
		return left == right, nil

	case instruction.OpNotEqual:
		return left != right, nil

	case instruction.OpLess:
		if l, ok := left.(int); ok {
			if r, ok := right.(int); ok {
				return l < r, nil
			}
		}
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l < r, nil
			}
		}
		// Handle mixed types
		if l, ok := left.(int); ok {
			if r, ok := right.(float64); ok {
				return float64(l) < r, nil
			}
		}
		if l, ok := left.(float64); ok {
			if r, ok := right.(int); ok {
				return l < float64(r), nil
			}
		}
		return nil, fmt.Errorf("unsupported types for less than comparison: %T and %T", left, right)

	case instruction.OpLessEqual:
		if l, ok := left.(int); ok {
			if r, ok := right.(int); ok {
				return l <= r, nil
			}
		}
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l <= r, nil
			}
		}
		// Handle mixed types
		if l, ok := left.(int); ok {
			if r, ok := right.(float64); ok {
				return float64(l) <= r, nil
			}
		}
		if l, ok := left.(float64); ok {
			if r, ok := right.(int); ok {
				return l <= float64(r), nil
			}
		}
		return nil, fmt.Errorf("unsupported types for less than or equal comparison: %T and %T", left, right)

	case instruction.OpGreater:
		if l, ok := left.(int); ok {
			if r, ok := right.(int); ok {
				return l > r, nil
			}
		}
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l > r, nil
			}
		}
		// Handle mixed types
		if l, ok := left.(int); ok {
			if r, ok := right.(float64); ok {
				return float64(l) > r, nil
			}
		}
		if l, ok := left.(float64); ok {
			if r, ok := right.(int); ok {
				return l > float64(r), nil
			}
		}
		return nil, fmt.Errorf("unsupported types for greater than comparison: %T and %T", left, right)

	case instruction.OpGreaterEqual:
		if l, ok := left.(int); ok {
			if r, ok := right.(int); ok {
				return l >= r, nil
			}
		}
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l >= r, nil
			}
		}
		// Handle mixed types
		if l, ok := left.(int); ok {
			if r, ok := right.(float64); ok {
				return float64(l) >= r, nil
			}
		}
		if l, ok := left.(float64); ok {
			if r, ok := right.(int); ok {
				return l >= float64(r), nil
			}
		}
		return nil, fmt.Errorf("unsupported types for greater than or equal comparison: %T and %T", left, right)

	case instruction.OpAnd:
		// Logical AND operation
		// In Go, && is short-circuit, but in our VM implementation, both operands are already evaluated
		// We just need to check if both are truthy
		return isTruthy(left) && isTruthy(right), nil

	case instruction.OpOr:
		// Logical OR operation
		// In Go, || is short-circuit, but in our VM implementation, both operands are already evaluated
		// We just need to check if either is truthy
		return isTruthy(left) || isTruthy(right), nil

	default:
		return nil, fmt.Errorf("unsupported binary operation: %d", op)
	}
}

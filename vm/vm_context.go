// Package vm provides the virtual machine implementation with context-based scope management
package vm

import (
	"fmt"

	"github.com/lengzhao/goscript/context"
	"github.com/lengzhao/goscript/types"
)

// VM represents the virtual machine with context-based scope management
type VM struct {
	// Optimized stack for expression evaluation
	stack *Stack

	// Global context (root of context hierarchy)
	globalCtx *context.Context

	// Current context
	currentCtx *context.Context

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

	// Opcode handlers dispatch table
	handlers map[OpCode]OpHandler
}

// NewVMWithContex creates a new virtual machine with context-based scope management
func NewVMWithContex() *VM {
	// Initialize with basic types
	typeSystem := make(map[string]types.IType)
	typeSystem["int"] = types.IntType.Clone()
	typeSystem["float64"] = types.Float64Type.Clone()
	typeSystem["string"] = types.StringType.Clone()
	typeSystem["bool"] = types.BoolType.Clone()

	// Create global context
	globalCtx := context.NewContext("global", nil)

	vm := &VM{
		stack:            NewStack(64, 10000), // Initial capacity 64, max 10000
		globalCtx:        globalCtx,
		currentCtx:       globalCtx,
		instructions:     make([]*Instruction, 0),
		ip:               0,
		functionRegistry: make(map[string]func(args ...interface{}) (interface{}, error)),
		scriptFunctions:  make(map[string]*ScriptFunction),
		typeSystem:       typeSystem,
		debug:            false,
		executionCount:   0,
		maxInstructions:  1000000, // Default maximum instruction limit of 1 million
		moduleManager:    nil,     // Initialize as nil
		handlers:         make(map[OpCode]OpHandler),
	}

	// Register all opcode handlers
	vm.registerHandlers()

	return vm
}

// EnterScope creates and enters a new scope with the given path key
func (vm *VM) EnterScope(pathKey string) *context.Context {
	// Create new context with current context as parent
	newCtx := context.NewContext(pathKey, vm.currentCtx)

	// Add the new context as a child of the current context
	if vm.currentCtx != nil {
		vm.currentCtx.AddChild(newCtx)
	}

	// Set as current context
	vm.currentCtx = newCtx

	// Debug output
	if vm.debug {
		fmt.Printf("Entered scope: %s\n", pathKey)
	}

	return newCtx
}

// ExitScope exits the current scope and returns to parent scope
func (vm *VM) ExitScope() *context.Context {
	if vm.currentCtx != nil && vm.currentCtx.GetParent() != nil {
		parent := vm.currentCtx.GetParent()

		// Remove the current context from parent's children
		if parent != nil {
			parent.RemoveChild(vm.currentCtx.GetPathKey())
		}

		// Set parent as current context
		vm.currentCtx = parent

		// Debug output
		if vm.debug {
			fmt.Printf("Exited scope, now in: %s\n", parent.GetPathKey())
		}

		return parent
	}

	// If we're already at the global context, stay there
	return vm.currentCtx
}

// SetVariable sets a variable, 对外接口，强制创建、覆盖
func (vm *VM) SetVariable(name string, value interface{}) error {
	if vm.currentCtx == nil {
		return fmt.Errorf("no current context")
	}

	vm.currentCtx.CreateVariableWithType(name, value, "unknown")

	// Debug output
	if vm.debug {
		fmt.Printf("Set variable in context %s: %s = %v\n", vm.currentCtx.GetPathKey(), name, value)
	}

	return nil
}

// GetVariable gets a variable, searching from current context up to root
func (vm *VM) GetVariable(name string) (interface{}, bool) {
	if vm.currentCtx == nil {
		return nil, false
	}

	return vm.currentCtx.GetVariable(name)
}

// DeleteVariable removes a variable from the current context
func (vm *VM) DeleteVariable(name string) {
	if vm.currentCtx == nil {
		return
	}

	vm.currentCtx.DeleteVariable(name)

	// Debug output
	if vm.debug {
		fmt.Printf("Deleted variable from context %s: %s\n", vm.currentCtx.GetPathKey(), name)
	}
}

// Run executes instructions within a specific context
func (vm *VM) Run(ctx *context.Context, startIP, endIP int) (interface{}, error) {
	// Validate input parameters
	if startIP < 0 || endIP > len(vm.instructions) || startIP > endIP {
		return nil, fmt.Errorf("invalid instruction range: [%d, %d)", startIP, endIP)
	}

	// Save current execution state
	prevCtx := vm.currentCtx
	prevIP := vm.ip
	prevRetval := vm.retval

	// Set new execution context
	vm.currentCtx = ctx
	vm.ip = startIP
	vm.retval = nil

	// Restore execution state when function returns
	defer func() {
		vm.currentCtx = prevCtx
		vm.ip = prevIP
		vm.retval = prevRetval
	}()

	// Execute instructions in the given range
	for vm.ip < endIP && vm.ip < len(vm.instructions) {
		// Check instruction limit
		if vm.maxInstructions > 0 && int64(vm.executionCount) >= vm.maxInstructions {
			return nil, fmt.Errorf("maximum instruction limit exceeded: %d instructions executed", vm.executionCount)
		}

		instr := vm.instructions[vm.ip]

		// Debug output
		if vm.debug {
			fmt.Printf("IP: %d, Instruction: %s, Stack: %v\n", vm.ip, instr.String(), vm.stack.GetSlice())
		}

		// Use dispatch table
		if handler, exists := vm.handlers[instr.Op]; exists {
			if err := handler(vm, instr); err != nil {
				return nil, fmt.Errorf("error executing instruction at IP %d: %w", vm.ip, err)
			}
		} else {
			return nil, fmt.Errorf("unknown opcode: %v at IP %d", instr.Op, vm.ip)
		}

		vm.ip++
		vm.executionCount++
	}

	return vm.retval, nil
}

// Execute executes the instructions using the new context-based approach
func (vm *VM) Execute(ctx *context.Context) (interface{}, error) {
	// Initialize with global context
	vm.currentCtx = vm.globalCtx
	vm.ip = 0
	vm.executionCount = 0
	vm.retval = nil

	// Run in the global context
	return vm.Run(vm.globalCtx, 0, len(vm.instructions))
}

// Package goscript provides the main interface for the GoScript engine
package goscript

import (
	"context"
	"fmt"
	"go/ast"
	"go/token"
	"time"

	"github.com/lengzhao/goscript/builtin"
	"github.com/lengzhao/goscript/compiler"
	"github.com/lengzhao/goscript/module"
	"github.com/lengzhao/goscript/parser"
	"github.com/lengzhao/goscript/runtime"
	"github.com/lengzhao/goscript/vm"
)

// Script represents a GoScript script
type Script struct {
	// Source code
	source []byte

	// Module manager
	moduleManager *module.ModuleManager

	// Virtual machine
	vm *vm.VM

	// Runtime
	runtime *runtime.Runtime

	// Debug mode
	debug bool

	// Execution statistics
	executionStats *ExecutionStats

	// Maximum number of instructions allowed (0 means no limit)
	maxInstructions int64
}

// ExecutionStats holds execution statistics
type ExecutionStats struct {
	ExecutionTime    time.Duration
	InstructionCount int
	ErrorCount       int
}

// NewScript creates a new script
func NewScript(source []byte) *Script {
	script := &Script{
		source:          source,
		moduleManager:   module.NewModuleManager(),
		vm:              vm.NewVM(),
		runtime:         runtime.NewRuntime(),
		debug:           false,
		executionStats:  &ExecutionStats{},
		maxInstructions: 10000, // Default limit of 10,000 instructions
	}

	// Register builtin functions with the VM
	for name, fn := range builtin.BuiltInFunctions {
		script.vm.RegisterFunction(name, func(f builtin.Function) func(args ...interface{}) (interface{}, error) {
			return func(args ...interface{}) (interface{}, error) {
				return f(args...)
			}
		}(fn))
	}

	return script
}

// SetMaxInstructions sets the maximum number of instructions allowed
func (s *Script) SetMaxInstructions(max int64) {
	s.maxInstructions = max
	s.vm.SetMaxInstructions(max)
}

// AddVariable adds a variable to the script
func (s *Script) AddVariable(name string, value interface{}) error {
	return s.vm.GlobalCtx.CreateVariableWithType(name, value, "unknow")
}

// GetVariable gets a variable from the script
func (s *Script) GetVariable(name string) (interface{}, bool) {
	return s.vm.GlobalCtx.GetVariable(name)
}

// SetVariable sets a variable in the script
func (s *Script) SetVariable(name string, value interface{}) error {
	return s.vm.GlobalCtx.SetVariable(name, value)
}

// AddFunction adds a function to the script
func (s *Script) AddFunction(name string, execFn vm.ScriptFunction) error {

	// Also register with the VM directly for immediate use
	s.vm.RegisterFunction(name, execFn)

	// Debug output
	if s.debug {
		fmt.Printf("Script: Added function %s\n", name)
	}

	return nil
}

// CallFunction calls a function in the script
func (s *Script) CallFunction(name string, args ...interface{}) (interface{}, error) {
	// Try to call the function in the current context
	return s.callFunctionInContext(name, args...)
}

// callFunctionInContext calls a function in the current context
func (s *Script) callFunctionInContext(name string, args ...interface{}) (interface{}, error) {
	// Debug output
	if s.debug {
		fmt.Printf("Script: Calling function %s with args %v\n", name, args)
	}

	// First check if it's a module function (format: moduleName.functionName)
	if len(name) > 0 && len(args) >= 0 {
		// Check if it's a module function call
		for i, char := range name {
			if char == '.' {
				moduleName := name[:i]
				functionName := name[i+1:]

				// Try to call the module function
				result, err := s.moduleManager.CallModuleFunction(moduleName, functionName, args...)
				if err == nil {
					return result, nil
				}

				// If it failed, continue to try other options
				break
			}
		}
	}

	// Try to call the function from the VM (functions registered via AddFunction)
	if fn, exists := s.vm.GetFunction(name); exists {
		result, err := fn(args...)
		if s.debug {
			if err != nil {
				fmt.Printf("Script: Error calling VM function %s: %v\n", name, err)
			} else {
				fmt.Printf("Script: Called VM function %s, result: %v\n", name, result)
			}
		}
		return result, err
	}

	// Try to call the function in the current module
	currentModule, exists := s.moduleManager.GetCurrentModule()
	if exists {
		if function, ok := currentModule.GetFunction(name); ok {
			// Call the function directly
			result, err := function(args...)
			if s.debug {
				if err != nil {
					fmt.Printf("Script: Error calling module function %s: %v\n", name, err)
				} else {
					fmt.Printf("Script: Called module function %s, result: %v\n", name, result)
				}
			}
			return result, err
		}
	}

	return nil, fmt.Errorf("function %s not found", name)
}

// Run executes the script
func (s *Script) Run() (interface{}, error) {
	return s.RunContext(context.Background())
}

// RunContext executes the script with a context
func (s *Script) RunContext(ctx context.Context) (interface{}, error) {
	fmt.Println("RunContext: Starting execution")
	startTime := time.Now()

	// Parse and compile the source code
	sourceStr := string(s.source)

	// Create a parser
	parser := parser.New()

	// Parse the source code into an AST
	astFile, err := parser.Parse("script.go", []byte(sourceStr), 0)
	if err != nil {
		return nil, fmt.Errorf("failed to parse source code: %w", err)
	}

	// Process import declarations before compilation
	for _, decl := range astFile.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.IMPORT {
			for _, spec := range genDecl.Specs {
				if importSpec, ok := spec.(*ast.ImportSpec); ok {
					// Get the import path (remove quotes)
					path := importSpec.Path.Value
					if len(path) > 2 {
						path = path[1 : len(path)-1] // Remove quotes
					}

					// Import the module
					err := s.ImportModule(path)
					if err != nil {
						return nil, fmt.Errorf("failed to import module %s: %w", path, err)
					}
				}
			}
		}
	}

	// Create a compiler instance
	compiler := compiler.NewCompiler(s.vm)

	// Compile the AST to bytecode
	err = compiler.Compile(astFile)
	if err != nil {
		return nil, fmt.Errorf("failed to compile AST: %w", err)
	}

	// Set max instructions in VM
	s.vm.SetMaxInstructions(s.maxInstructions)

	// Execute the VM
	fmt.Println("RunContext: Executing VM")
	result, err := s.vm.Execute("")
	fmt.Printf("RunContext: VM execution completed, result: %v, err: %v\n", result, err)

	// Update execution statistics
	s.executionStats.ExecutionTime = time.Since(startTime)
	if err != nil {
		s.executionStats.ErrorCount++
	}

	// Get instruction count from VM
	s.executionStats.InstructionCount = int(s.vm.GetInstructionCount())

	if err != nil {
		return nil, err
	}

	return result, nil
}

// ImportModule imports a module into the script
func (s *Script) ImportModule(moduleNames ...string) error {
	for _, moduleName := range moduleNames {
		err := s.moduleManager.ImportModule(moduleName)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetModuleManager returns the module manager
func (s *Script) GetModuleManager() *module.ModuleManager {
	return s.moduleManager
}

// String returns a string representation of the script
func (s *Script) String() string {
	return fmt.Sprintf("Script{source: %d bytes, modules: %d}",
		len(s.source), len(s.moduleManager.GetAllModules()))
}

// SetDebug enables or disables debug mode
func (s *Script) SetDebug(debug bool) {
	s.debug = debug
	s.vm.SetDebug(debug)
	s.moduleManager.SetDebug(debug)
	s.runtime.SetDebug(debug)
}

// GetExecutionStats returns execution statistics
func (s *Script) GetExecutionStats() *ExecutionStats {
	return s.executionStats
}

// GetVM returns the virtual machine
func (s *Script) GetVM() *vm.VM {
	return s.vm
}

// GetRuntime returns the runtime
func (s *Script) GetRuntime() *runtime.Runtime {
	return s.runtime
}

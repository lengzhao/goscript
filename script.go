// Package goscript provides the main interface for the GoScript engine
package goscript

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/lengzhao/goscript/builtin"
	"github.com/lengzhao/goscript/compiler"
	execContext "github.com/lengzhao/goscript/context"
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

	// Global execution context
	globalContext *execContext.ExecutionContext

	// Virtual machine
	vm *vm.VM

	// Runtime
	runtime *runtime.Runtime

	// Debug mode
	debug bool

	// Execution statistics
	executionStats *ExecutionStats
}

// ExecutionStats holds execution statistics
type ExecutionStats struct {
	ExecutionTime    time.Duration
	InstructionCount int
	ErrorCount       int
}

// SecurityContext defines security restrictions
type SecurityContext execContext.SecurityContext

// NewScript creates a new script
func NewScript(source []byte) *Script {
	script := &Script{
		source:         source,
		moduleManager:  module.NewModuleManager(),
		globalContext:  execContext.NewExecutionContext(),
		vm:             vm.NewVM(),
		runtime:        runtime.NewRuntime(),
		debug:          false,
		executionStats: &ExecutionStats{},
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

// SetSecurityContext sets the security context
func (s *Script) SetSecurityContext(securityCtx *execContext.SecurityContext) {
	s.globalContext.Security = securityCtx

	// Also set it in the module manager's global context
	s.moduleManager.GetGlobalContext().Security = securityCtx
}

// AddVariable adds a variable to the script
func (s *Script) AddVariable(name string, value interface{}) error {
	return s.globalContext.SetValue(name, value)
}

// GetVariable gets a variable from the script
func (s *Script) GetVariable(name string) (interface{}, bool) {
	return s.globalContext.GetValue(name)
}

// SetVariable sets a variable in the script
func (s *Script) SetVariable(name string, value interface{}) error {
	return s.globalContext.SetValue(name, value)
}

// AddFunction adds a function to the script
func (s *Script) AddFunction(name string, execFn execContext.Function) error {
	// Register the function in the global context
	err := s.globalContext.RegisterFunction(name, execFn)
	if err != nil {
		return err
	}

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

	// Try to call the function in the global context
	if fn, exists := s.globalContext.GetFunction(name); exists {
		result, err := fn(args...)
		if s.debug {
			if err != nil {
				fmt.Printf("Script: Error calling function %s: %v\n", name, err)
			} else {
				fmt.Printf("Script: Called function %s, result: %v\n", name, result)
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

	// Use the global context
	execCtx := s.globalContext

	// Apply security context timeout if set
	if execCtx.Security != nil && execCtx.Security.MaxExecutionTime > 0 {
		execCtx = execCtx.WithTimeout(execCtx.Security.MaxExecutionTime)
		// Also update the global context reference to use the new context with timeout
		s.globalContext = execCtx
	}

	// Create a module for the script
	scriptModule := module.NewModule("main")
	scriptModule.SetDebug(s.debug)

	// Register the module
	s.moduleManager.RegisterModule(scriptModule)

	// Set the current module
	s.moduleManager.SetCurrentModule("main")

	// Clear any existing instructions
	s.vm = vm.NewVM()
	s.vm.SetDebug(s.debug)
	s.vm.SetModuleManager(s.moduleManager) // 设置模块管理器引用

	// Set maximum instruction limit from security context
	if execCtx.Security != nil && execCtx.Security.MaxInstructions > 0 {
		s.vm.SetMaxInstructions(execCtx.Security.MaxInstructions)
	}

	// Register builtin functions with the VM
	for name, fn := range builtin.BuiltInFunctions {
		s.vm.RegisterFunction(name, func(f builtin.Function) func(args ...interface{}) (interface{}, error) {
			return func(args ...interface{}) (interface{}, error) {
				return f(args...)
			}
		}(fn))
	}

	// Generate bytecode based on source content
	// This is a simplified compilation process for demonstration purposes
	sourceStr := string(s.source)

	fmt.Printf("RunContext: Source code:\n%s\n", sourceStr)

	// Create a parser
	parser := parser.New()

	// Parse the source code into an AST
	astFile, err := parser.Parse("script.go", []byte(sourceStr), 0)
	if err != nil {
		return nil, fmt.Errorf("failed to parse source code: %w", err)
	}

	// Create a compiler instance
	compiler := compiler.NewCompiler(s.vm, s.globalContext)

	// Compile the AST to bytecode
	err = compiler.Compile(astFile)
	if err != nil {
		return nil, fmt.Errorf("failed to compile AST: %w", err)
	}

	// After compiling, we know which modules are needed
	// Import all needed modules
	neededModules := compiler.GetImports()
	fmt.Printf("Needed modules: %v\n", neededModules)

	// Import all needed modules
	for _, modulePath := range neededModules {
		// Extract module name from path
		parts := strings.Split(modulePath, "/")
		moduleName := parts[len(parts)-1]

		// Import the module
		err := s.moduleManager.ImportModule(moduleName)
		if err != nil {
			// If the module doesn't exist, continue to the next one
			fmt.Printf("Failed to import module %s: %v\n", moduleName, err)
			continue
		}
		fmt.Printf("Successfully imported module %s\n", moduleName)
	}

	// Register any functions that were added via AddFunction
	// We need to get all functions from the execution context
	if rootCtx := execCtx.GetRootContext(); rootCtx != nil {
		// Get all functions from the scope manager
		allFunctions := rootCtx.ScopeManager.GetAllFunctions()
		for name, fn := range allFunctions {
			// Register each function with the VM
			s.vm.RegisterFunction(name, fn)
		}
	}

	// Execute the VM
	fmt.Println("RunContext: Executing VM")
	result, err := s.vm.Execute(nil)
	fmt.Printf("RunContext: VM execution completed, result: %v, err: %v\n", result, err)

	// Update execution statistics
	s.executionStats.ExecutionTime = time.Since(startTime)
	s.executionStats.InstructionCount = s.vm.GetExecutionCount()
	if err != nil {
		s.executionStats.ErrorCount++
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

// compileSource generates bytecode instructions based on the source code content
// This implementation parses the AST and generates proper bytecode
func (s *Script) compileSource(sourceStr string) error {
	// Create a parser
	parser := parser.New()

	// Parse the source code into an AST
	astFile, err := parser.Parse("script.go", []byte(sourceStr), 0)
	if err != nil {
		return fmt.Errorf("failed to parse source code: %w", err)
	}

	// Create a compiler instance
	compiler := compiler.NewCompiler(s.vm, s.globalContext)

	// Compile the AST to bytecode
	err = compiler.Compile(astFile)
	if err != nil {
		return fmt.Errorf("failed to compile AST: %w", err)
	}

	return nil
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

// GetGlobalContext returns the global execution context
func (s *Script) GetGlobalContext() *execContext.ExecutionContext {
	return s.globalContext
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
	s.globalContext.SetDebug(debug)
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

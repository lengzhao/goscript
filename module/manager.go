// Package module provides module management for GoScript
package module

import (
	"fmt"

	"github.com/lengzhao/goscript/context"
	"github.com/lengzhao/goscript/symbol"
	"github.com/lengzhao/goscript/vm"
)

// Function represents a callable function in a module
type Function func(args ...interface{}) (interface{}, error)

// Module represents a module in GoScript
type Module struct {
	// Name is the module name
	Name string

	// Instructions are the bytecode instructions for the module
	Instructions []*vm.Instruction

	// SymbolTable is the symbol table for the module
	SymbolTable *symbol.SymbolTable

	// Context is the execution context for the module
	Context *context.ExecutionContext

	// Functions maps function names to functions
	Functions map[string]Function

	// Dependencies are module names this module depends on
	Dependencies []string

	// Debug mode
	Debug bool
}

// NewModule creates a new module
func NewModule(name string) *Module {
	return &Module{
		Name:         name,
		Instructions: make([]*vm.Instruction, 0),
		SymbolTable:  symbol.NewSymbolTable(),
		Context:      context.NewExecutionContext(),
		Functions:    make(map[string]Function),
		Dependencies: make([]string, 0),
		Debug:        false,
	}
}

// AddInstruction adds an instruction to the module
func (m *Module) AddInstruction(instruction *vm.Instruction) {
	m.Instructions = append(m.Instructions, instruction)

	// Debug output
	if m.Debug {
		fmt.Printf("Module %s: Added instruction %s\n", m.Name, instruction.String())
	}
}

// GetInstructions returns all instructions in the module
func (m *Module) GetInstructions() []*vm.Instruction {
	return m.Instructions
}

// AddSymbol adds a symbol to the module
func (m *Module) AddSymbol(sym *symbol.Symbol) {
	sym.Module = m.Name
	m.SymbolTable.Set(sym.Name, sym)

	// Debug output
	if m.Debug {
		fmt.Printf("Module %s: Added symbol %s\n", m.Name, sym.Name)
	}
}

// GetSymbol retrieves a symbol from the module
func (m *Module) GetSymbol(name string) *symbol.Symbol {
	return m.SymbolTable.Get(name)
}

// SetContext sets the execution context for the module
func (m *Module) SetContext(ctx *context.ExecutionContext) {
	m.Context = ctx
	m.Context.ModuleName = m.Name

	// Debug output
	if m.Debug {
		fmt.Printf("Module %s: Set context\n", m.Name)
	}
}

// AddFunction adds a function to the module
func (m *Module) AddFunction(name string, fn Function) error {
	if _, exists := m.Functions[name]; exists {
		return fmt.Errorf("function %s already exists in module %s", name, m.Name)
	}

	m.Functions[name] = fn

	// Also add the function to the symbol table
	sym := &symbol.Symbol{
		Module:     m.Name,
		ID:         name,
		Name:       name,
		Type:       nil, // Function type to be defined
		Address:    fn,
		ScopeLevel: 0,
		Mutable:    false,
	}

	m.SymbolTable.Set(name, sym)

	// Debug output
	if m.Debug {
		fmt.Printf("Module %s: Added function %s\n", m.Name, name)
	}

	return nil
}

// GetFunction retrieves a function from the module
func (m *Module) GetFunction(name string) (Function, bool) {
	fn, exists := m.Functions[name]
	return fn, exists
}

// GetAllFunctions returns all functions in the module
func (m *Module) GetAllFunctions() map[string]Function {
	return m.Functions
}

// AddDependency adds a dependency to this module
func (m *Module) AddDependency(moduleName string) {
	// Check if dependency already exists
	for _, dep := range m.Dependencies {
		if dep == moduleName {
			return // Already exists
		}
	}

	m.Dependencies = append(m.Dependencies, moduleName)

	// Debug output
	if m.Debug {
		fmt.Printf("Module %s: Added dependency %s\n", m.Name, moduleName)
	}
}

// GetDependencies returns all dependencies of this module
func (m *Module) GetDependencies() []string {
	return m.Dependencies
}

// SetDebug enables or disables debug mode
func (m *Module) SetDebug(debug bool) {
	m.Debug = debug
	m.Context.SetDebug(debug)
}

// String returns a string representation of the module
func (m *Module) String() string {
	return fmt.Sprintf("Module{name: %s, functions: %d, dependencies: %d}",
		m.Name, len(m.Functions), len(m.Dependencies))
}

// ModuleManager manages modules in GoScript
type ModuleManager struct {
	// modules maps module names to modules
	modules map[string]*Module

	// currentModule is the name of the current module
	currentModule string

	// globalContext is the global execution context
	globalContext *context.ExecutionContext

	// Debug mode
	debug bool
}

// NewModuleManager creates a new module manager
func NewModuleManager() *ModuleManager {
	return &ModuleManager{
		modules:       make(map[string]*Module),
		currentModule: "main",
		globalContext: context.NewExecutionContext(),
		debug:         false,
	}
}

// RegisterModule registers a module
func (mm *ModuleManager) RegisterModule(module *Module) error {
	if _, exists := mm.modules[module.Name]; exists {
		return fmt.Errorf("module %s already registered", module.Name)
	}

	mm.modules[module.Name] = module

	// Debug output
	if mm.debug {
		fmt.Printf("ModuleManager: Registered module %s\n", module.Name)
	}

	return nil
}

// GetModule retrieves a module by name
func (mm *ModuleManager) GetModule(name string) (*Module, bool) {
	module, exists := mm.modules[name]
	return module, exists
}

// SetCurrentModule sets the current module
func (mm *ModuleManager) SetCurrentModule(name string) error {
	if _, exists := mm.modules[name]; !exists {
		return fmt.Errorf("module %s not found", name)
	}

	mm.currentModule = name

	// Debug output
	if mm.debug {
		fmt.Printf("ModuleManager: Set current module to %s\n", name)
	}

	return nil
}

// GetCurrentModule returns the current module
func (mm *ModuleManager) GetCurrentModule() (*Module, bool) {
	return mm.GetModule(mm.currentModule)
}

// CallModuleFunction calls a function in a specific module
func (mm *ModuleManager) CallModuleFunction(moduleName, functionName string, args ...interface{}) (interface{}, error) {
	module, exists := mm.modules[moduleName]
	if !exists {
		return nil, fmt.Errorf("module %s not found", moduleName)
	}

	// Check if the function exists in the module
	fn, exists := module.GetFunction(functionName)
	if !exists {
		// Check if the function exists in the module's symbol table
		symbol := module.SymbolTable.Get(functionName)
		if symbol == nil {
			return nil, fmt.Errorf("function %s not found in module %s", functionName, moduleName)
		}

		// If it's in the symbol table but not in the functions map,
		// it might be a built-in function or a variable
		if fn, ok := symbol.Address.(Function); ok {
			// Call the function
			return fn(args...)
		}

		return nil, fmt.Errorf("function %s not found in module %s", functionName, moduleName)
	}

	// Call the function
	result, err := fn(args...)

	// Debug output
	if mm.debug {
		if err != nil {
			fmt.Printf("ModuleManager: Error calling function %s.%s: %v\n", moduleName, functionName, err)
		} else {
			fmt.Printf("ModuleManager: Called function %s.%s, result: %v\n", moduleName, functionName, result)
		}
	}

	return result, err
}

// ImportModule imports a module into the current context
func (mm *ModuleManager) ImportModule(moduleName string, ctx *context.ExecutionContext) error {
	module, exists := mm.modules[moduleName]
	if !exists {
		return fmt.Errorf("module %s not found", moduleName)
	}

	// Import all symbols from the module into the context
	symbols := module.SymbolTable.GetAllSymbols()
	for name, sym := range symbols {
		// For now, we'll just print a message as we don't have direct access to the symbol table
		// In a real implementation, we would add the symbol to the context's symbol table
		_ = name
		_ = sym

		// Debug output
		if mm.debug {
			fmt.Printf("ModuleManager: Imported symbol %s from module %s\n", name, moduleName)
		}
	}

	return nil
}

// GetGlobalContext returns the global execution context
func (mm *ModuleManager) GetGlobalContext() *context.ExecutionContext {
	return mm.globalContext
}

// GetAllModules returns all registered modules
func (mm *ModuleManager) GetAllModules() map[string]*Module {
	return mm.modules
}

// GetModuleNames returns the names of all registered modules
func (mm *ModuleManager) GetModuleNames() []string {
	names := make([]string, 0, len(mm.modules))
	for name := range mm.modules {
		names = append(names, name)
	}
	return names
}

// RemoveModule removes a module
func (mm *ModuleManager) RemoveModule(name string) error {
	if _, exists := mm.modules[name]; !exists {
		return fmt.Errorf("module %s not found", name)
	}

	// Check if this is the current module
	if mm.currentModule == name {
		return fmt.Errorf("cannot remove current module %s", name)
	}

	delete(mm.modules, name)

	// Debug output
	if mm.debug {
		fmt.Printf("ModuleManager: Removed module %s\n", name)
	}

	return nil
}

// SetDebug enables or disables debug mode
func (mm *ModuleManager) SetDebug(debug bool) {
	mm.debug = debug
	mm.globalContext.SetDebug(debug)

	// Set debug mode for all modules
	for _, module := range mm.modules {
		module.SetDebug(debug)
	}
}

// ResolveDependencies resolves module dependencies
func (mm *ModuleManager) ResolveDependencies() error {
	// This is a simple dependency resolution algorithm
	// In a real implementation, this would be more sophisticated

	// For now, we just check if all dependencies exist
	for moduleName, module := range mm.modules {
		for _, depName := range module.GetDependencies() {
			if _, exists := mm.modules[depName]; !exists {
				return fmt.Errorf("module %s depends on missing module %s", moduleName, depName)
			}
		}
	}

	return nil
}

// String returns a string representation of the module manager
func (mm *ModuleManager) String() string {
	result := "ModuleManager{modules: ["
	first := true
	for name := range mm.modules {
		if !first {
			result += ", "
		}
		result += name
		first = false
	}
	result += fmt.Sprintf("], current: %s}", mm.currentModule)
	return result
}

// DebugString returns a detailed string representation for debugging
func (mm *ModuleManager) DebugString() string {
	result := "ModuleManager{\n"
	result += fmt.Sprintf("  Current module: %s\n", mm.currentModule)
	result += "  Modules:\n"
	for name, module := range mm.modules {
		result += fmt.Sprintf("    %s: %s\n", name, module.String())
		// Print dependencies
		deps := module.GetDependencies()
		if len(deps) > 0 {
			result += "      Dependencies: "
			for i, dep := range deps {
				if i > 0 {
					result += ", "
				}
				result += dep
			}
			result += "\n"
		}
	}
	result += "}"
	return result
}

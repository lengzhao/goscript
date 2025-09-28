// Package module provides module management for GoScript
package module

import (
	"fmt"

	"github.com/lengzhao/goscript/builtin"
	"github.com/lengzhao/goscript/instruction"
	"github.com/lengzhao/goscript/symbol"
	"github.com/lengzhao/goscript/types"
)

// Function represents a callable function in a module
type Function = types.Function

// Module represents a module in GoScript
type Module struct {
	// Name is the module name
	Name string

	// Instructions are the bytecode instructions for the module
	Instructions []*instruction.Instruction

	// SymbolTable is the symbol table for the module
	SymbolTable *symbol.SymbolTable

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
		Instructions: make([]*instruction.Instruction, 0),
		SymbolTable:  symbol.NewSymbolTable(),
		Functions:    make(map[string]Function),
		Dependencies: make([]string, 0),
		Debug:        false,
	}
}

// AddInstruction adds an instruction to the module
func (m *Module) AddInstruction(instruction *instruction.Instruction) {
	m.Instructions = append(m.Instructions, instruction)

	// Debug output
	if m.Debug {
		fmt.Printf("Module %s: Added instruction %s\n", m.Name, instruction.String())
	}
}

// GetInstructions returns all instructions in the module
func (m *Module) GetInstructions() []*instruction.Instruction {
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

	// Debug mode
	debug bool
}

// NewModuleManager creates a new module manager
func NewModuleManager() *ModuleManager {
	mm := &ModuleManager{
		modules:       make(map[string]*Module),
		currentModule: "main",
		debug:         false,
	}

	// Register default modules
	mm.registerDefaultModules()

	return mm
}

// registerDefaultModules registers default builtin modules
func (mm *ModuleManager) registerDefaultModules() {
	// Register math module
	mathModule := NewModule("math")

	// Add math functions
	// For now, we'll add a simple abs function
	mathModule.AddFunction("abs", func(args ...interface{}) (interface{}, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("abs function requires 1 argument")
		}
		if num, ok := args[0].(int); ok {
			if num < 0 {
				return -num, nil
			}
			return num, nil
		}
		return nil, fmt.Errorf("abs function requires integer argument")
	})

	mm.RegisterModule(mathModule)
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
		// Check if it's a builtin module
		moduleFuncs, isBuiltin := builtin.GetModuleFunctions(moduleName)
		if !isBuiltin {
			return nil, fmt.Errorf("module %s not found", moduleName)
		}

		// Create a new module for the builtin functions
		module = NewModule(moduleName)

		// Add all functions from the builtin module
		for funcName, fn := range moduleFuncs {
			err := module.AddFunction(funcName, fn)
			if err != nil {
				return nil, fmt.Errorf("failed to add function %s to module %s: %w", funcName, moduleName, err)
			}
		}

		// Register the module
		err := mm.RegisterModule(module)
		if err != nil {
			return nil, fmt.Errorf("failed to register module %s: %w", moduleName, err)
		}
	} else {
		// Module exists, check if the function is already registered
		// If not, and it's a builtin module, add the function
		_, funcExists := module.GetFunction(functionName)
		if !funcExists {
			// Check if it's a builtin module
			moduleFuncs, isBuiltin := builtin.GetModuleFunctions(moduleName)
			if isBuiltin {
				// Check if the function exists in the builtin module
				if fn, ok := moduleFuncs[functionName]; ok {
					// Add the function to the module
					err := module.AddFunction(functionName, fn)
					if err != nil {
						return nil, fmt.Errorf("failed to add function %s to module %s: %w", functionName, moduleName, err)
					}
				}
			}
		}
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
func (mm *ModuleManager) ImportModule(moduleName string) error {
	// 只检查模块是否存在，不立即导入所有函数
	_, exists := mm.modules[moduleName]
	if !exists {
		// Check if it's a builtin module
		_, isBuiltin := builtin.GetModuleFunctions(moduleName)
		if !isBuiltin {
			return fmt.Errorf("module %s not found", moduleName)
		}

		// 创建模块但不立即注册所有函数
		module := NewModule(moduleName)
		err := mm.RegisterModule(module)
		if err != nil {
			return fmt.Errorf("failed to register module %s: %w", moduleName, err)
		}
	}

	// 只需记录模块已被导入，不需要实际导入所有函数
	return nil
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

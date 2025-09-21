// Package symbol provides symbol table management for GoScript
package symbol

import (
	"fmt"

	"github.com/lengzhao/goscript/types"
)

// Symbol represents a symbol in the symbol table
type Symbol struct {
	// Module is the module this symbol belongs to
	Module string

	// ID is the unique identifier for the symbol
	ID string

	// Name is the name of the symbol
	Name string

	// Type is the type of the symbol
	Type types.IType

	// Address is a pointer to the actual value storage
	Address interface{}

	// ScopeLevel is the scope level of the symbol
	ScopeLevel int

	// Mutable indicates whether the symbol can be modified
	Mutable bool
}

// SymbolTable manages symbols
type SymbolTable struct {
	// symbols maps symbol names to symbols
	symbols map[string]*Symbol

	// modules maps module names to lists of symbols
	modules map[string][]*Symbol

	// currentScope is the current scope level
	currentScope int

	// parent is the parent symbol table for scope nesting
	parent *SymbolTable
}

// NewSymbolTable creates a new symbol table
func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		symbols:      make(map[string]*Symbol),
		modules:      make(map[string][]*Symbol),
		currentScope: 0,
		parent:       nil,
	}
}

// NewChildSymbolTable creates a new symbol table with a parent
func NewChildSymbolTable(parent *SymbolTable) *SymbolTable {
	return &SymbolTable{
		symbols:      make(map[string]*Symbol),
		modules:      make(map[string][]*Symbol),
		currentScope: parent.currentScope + 1,
		parent:       parent,
	}
}

// Set adds or updates a symbol in the symbol table
func (st *SymbolTable) Set(name string, symbol *Symbol) {
	if symbol == nil {
		return
	}

	// Set the scope level of the symbol
	symbol.ScopeLevel = st.currentScope

	st.symbols[name] = symbol

	// Add to module if module is specified
	if symbol.Module != "" {
		st.modules[symbol.Module] = append(st.modules[symbol.Module], symbol)
	}
}

// Get retrieves a symbol by name
func (st *SymbolTable) Get(name string) *Symbol {
	// Look in current scope first
	if symbol, exists := st.symbols[name]; exists {
		return symbol
	}

	// If not found and we have a parent, look in parent
	if st.parent != nil {
		return st.parent.Get(name)
	}

	// Not found
	return nil
}

// GetLocal retrieves a symbol by name only in the current scope
func (st *SymbolTable) GetLocal(name string) *Symbol {
	// Look only in current scope
	if symbol, exists := st.symbols[name]; exists {
		return symbol
	}

	// Not found in current scope
	return nil
}

// GetCurrentScope returns the current scope level
func (st *SymbolTable) GetCurrentScope() int {
	return st.currentScope
}

// EnterScope increases the scope level
func (st *SymbolTable) EnterScope() {
	st.currentScope++
}

// ExitScope decreases the scope level
func (st *SymbolTable) ExitScope() {
	if st.currentScope > 0 {
		st.currentScope--
	}

	// Clean up symbols that belong to the exiting scope
	for name, symbol := range st.symbols {
		if symbol.ScopeLevel > st.currentScope {
			delete(st.symbols, name)

			// Also remove from module if applicable
			if symbol.Module != "" {
				if symbols, exists := st.modules[symbol.Module]; exists {
					// Find and remove the symbol from the module's list
					for i, sym := range symbols {
						if sym == symbol {
							// Remove the symbol from the slice
							st.modules[symbol.Module] = append(symbols[:i], symbols[i+1:]...)
							break
						}
					}
				}
			}
		}
	}
}

// Remove removes a symbol by name
func (st *SymbolTable) Remove(name string) {
	delete(st.symbols, name)
}

// GetAllSymbols returns all symbols in the current table
func (st *SymbolTable) GetAllSymbols() map[string]*Symbol {
	// Return a copy to prevent external modification
	result := make(map[string]*Symbol)
	for k, v := range st.symbols {
		result[k] = v
	}
	return result
}

// GetSymbolsByModule returns all symbols belonging to a module
func (st *SymbolTable) GetSymbolsByModule(moduleName string) ([]*Symbol, bool) {
	symbols, exists := st.modules[moduleName]
	return symbols, exists
}

// GetAllModules returns all module names
func (st *SymbolTable) GetAllModules() []string {
	modules := make([]string, 0, len(st.modules))
	for module := range st.modules {
		modules = append(modules, module)
	}
	return modules
}

// Contains checks if a symbol exists in the symbol table
func (st *SymbolTable) Contains(name string) bool {
	return st.Get(name) != nil
}

// UpdateSymbolType updates the type of an existing symbol
func (st *SymbolTable) UpdateSymbolType(name string, newType types.IType) error {
	symbol := st.Get(name)
	if symbol == nil {
		return fmt.Errorf("symbol %s not found", name)
	}

	symbol.Type = newType
	return nil
}

// String returns a string representation of the symbol table
func (st *SymbolTable) String() string {
	result := fmt.Sprintf("SymbolTable{scope: %d, symbols: [", st.currentScope)
	first := true
	for name, symbol := range st.symbols {
		if !first {
			result += ", "
		}
		result += fmt.Sprintf("%s:%s", name, symbol.Type.TypeName())
		first = false
	}
	result += "]}"
	return result
}

// DebugString returns a detailed string representation for debugging
func (st *SymbolTable) DebugString() string {
	result := fmt.Sprintf("SymbolTable{scope: %d, symbols: {\n", st.currentScope)
	for name, symbol := range st.symbols {
		result += fmt.Sprintf("  %s: {Module: %s, Type: %s, Scope: %d, Mutable: %t}\n",
			name, symbol.Module, symbol.Type.TypeName(), symbol.ScopeLevel, symbol.Mutable)
	}
	result += "}}"
	return result
}

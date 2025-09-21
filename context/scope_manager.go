// Package context provides scope management
package context

import (
	"fmt"
	"sync"
	"time"
)

// ScopeManager manages variable scopes
type ScopeManager struct {
	// Root scope
	root *Scope

	// Current scope
	current *Scope

	// Function registry
	functions map[string]Function

	// Cancel function for timeout
	cancelFunc func()

	// Deadline for timeout
	deadline time.Time

	// mutex for thread safety
	mu sync.RWMutex

	// Debug mode
	debug bool
}

// Scope represents a variable scope
type Scope struct {
	variables map[string]*Variable
	parent    *Scope
}

// Variable represents a variable in the scope
type Variable struct {
	Name  string
	Value interface{}
	Type  string
}

// Function represents a callable function
type Function interface {
	// Call executes the function with the given arguments
	Call(args ...interface{}) (interface{}, error)

	// Name returns the function name
	Name() string
}

// NewScopeManager creates a new scope manager
func NewScopeManager() *ScopeManager {
	root := &Scope{
		variables: make(map[string]*Variable),
		parent:    nil,
	}

	return &ScopeManager{
		root:      root,
		current:   root,
		functions: make(map[string]Function),
		debug:     false,
	}
}

// NewScopeManagerWithTimeout creates a new scope manager with timeout
func NewScopeManagerWithTimeout(timeout time.Duration) *ScopeManager {
	root := &Scope{
		variables: make(map[string]*Variable),
		parent:    nil,
	}

	deadline := time.Now().Add(timeout)

	return &ScopeManager{
		root:       root,
		current:    root,
		functions:  make(map[string]Function),
		deadline:   deadline,
		cancelFunc: func() {}, // Empty cancel function for now
		debug:      false,
	}
}

// NewScope creates a new scope (entering a new context)
func (sm *ScopeManager) NewScope() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Create a new scope with current scope as parent
	newScope := &Scope{
		variables: make(map[string]*Variable),
		parent:    sm.current,
	}

	sm.current = newScope

	// Debug output
	if sm.debug {
		scopeCount := 0
		scope := sm.current
		for scope != nil {
			scopeCount++
			scope = scope.parent
		}
		fmt.Printf("Entered new scope. Current scope level: %d\n", scopeCount)
	}
}

// ExitScope exits the current scope and returns to parent scope
func (sm *ScopeManager) ExitScope() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.current.parent != nil {
		sm.current = sm.current.parent

		// Debug output
		if sm.debug {
			scopeCount := 0
			scope := sm.current
			for scope != nil {
				scopeCount++
				scope = scope.parent
			}
			fmt.Printf("Exited scope. Current scope level: %d\n", scopeCount)
		}
	}
}

// SetVariable sets a variable in the current scope
func (sm *ScopeManager) SetVariable(name string, value interface{}, varType string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Set the variable in current scope
	sm.current.variables[name] = &Variable{
		Name:  name,
		Value: value,
		Type:  varType,
	}

	// Debug output
	if sm.debug {
		fmt.Printf("Set variable in scope: %s = %v (type: %s)\n", name, value, varType)
	}
}

// GetVariable gets a variable, searching from current scope up to root
func (sm *ScopeManager) GetVariable(name string) (*Variable, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	// Search from current scope up the chain
	scope := sm.current
	for scope != nil {
		if variable, exists := scope.variables[name]; exists {
			return variable, true
		}
		scope = scope.parent
	}

	return nil, false
}

// UpdateVariable updates an existing variable in the scope chain
func (sm *ScopeManager) UpdateVariable(name string, value interface{}) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Search from current scope up the chain to find the variable
	scope := sm.current
	for scope != nil {
		if variable, exists := scope.variables[name]; exists {
			variable.Value = value

			// Debug output
			if sm.debug {
				fmt.Printf("Updated variable: %s = %v\n", name, value)
			}

			return nil
		}
		scope = scope.parent
	}

	return fmt.Errorf("variable %s not found", name)
}

// DeleteVariable deletes a variable from the current scope
func (sm *ScopeManager) DeleteVariable(name string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Delete the variable from current scope
	delete(sm.current.variables, name)

	// Debug output
	if sm.debug {
		fmt.Printf("Deleted variable: %s\n", name)
	}
}

// GetCurrentScopeVariables returns all variables in the current scope
func (sm *ScopeManager) GetCurrentScopeVariables() map[string]*Variable {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	// Return a copy to prevent external modification
	result := make(map[string]*Variable)
	for k, v := range sm.current.variables {
		result[k] = v
	}
	return result
}

// GetAllVariables returns all variables in all scopes
func (sm *ScopeManager) GetAllVariables() map[string]*Variable {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	// Collect variables from all scopes
	result := make(map[string]*Variable)

	// Search from current scope up the chain
	scope := sm.current
	for scope != nil {
		for k, v := range scope.variables {
			// Only add if not already present (inner scopes override outer scopes)
			if _, exists := result[k]; !exists {
				result[k] = v
			}
		}
		scope = scope.parent
	}

	return result
}

// RegisterFunction registers a function in the scope manager
func (sm *ScopeManager) RegisterFunction(name string, fn Function) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.functions[name] = fn

	// Debug output
	if sm.debug {
		fmt.Printf("Registered function: %s\n", name)
	}
}

// GetFunction gets a function by name
func (sm *ScopeManager) GetFunction(name string) (Function, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	fn, exists := sm.functions[name]
	return fn, exists
}

// GetAllFunctions returns all registered functions
func (sm *ScopeManager) GetAllFunctions() map[string]Function {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	// Return a copy to prevent external modification
	result := make(map[string]Function)
	for k, v := range sm.functions {
		result[k] = v
	}
	return result
}

// IsCancelled checks if the scope has been cancelled (timeout)
func (sm *ScopeManager) IsCancelled() bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	// Check if we have a deadline and if it has passed
	if !sm.deadline.IsZero() && time.Now().After(sm.deadline) {
		return true
	}
	return false
}

// SetDeadline sets a deadline for the scope
func (sm *ScopeManager) SetDeadline(deadline time.Time) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.deadline = deadline
}

// GetDeadline returns the deadline of the scope
func (sm *ScopeManager) GetDeadline() (time.Time, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if sm.deadline.IsZero() {
		return time.Time{}, false
	}
	return sm.deadline, true
}

// SetDebug enables or disables debug mode
func (sm *ScopeManager) SetDebug(debug bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.debug = debug
}

// GetScopeLevel returns the current scope level
func (sm *ScopeManager) GetScopeLevel() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	level := 0
	scope := sm.current
	for scope != nil {
		level++
		scope = scope.parent
	}
	return level
}

// String returns a string representation of the scope manager
func (sm *ScopeManager) String() string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	scopeCount := 0
	scope := sm.current
	for scope != nil {
		scopeCount++
		scope = scope.parent
	}

	deadlineStr := "no deadline"
	if !sm.deadline.IsZero() {
		deadlineStr = sm.deadline.String()
	}

	return fmt.Sprintf("ScopeManager{scopes: %d, functions: %d, deadline: %s}", scopeCount, len(sm.functions), deadlineStr)
}

// DebugString returns a detailed string representation for debugging
func (sm *ScopeManager) DebugString() string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	result := "ScopeManager{\n"

	// Print scope hierarchy
	scope := sm.current
	level := 0
	for scope != nil {
		result += fmt.Sprintf("  Scope %d: %d variables\n", level, len(scope.variables))
		for name, variable := range scope.variables {
			result += fmt.Sprintf("    %s: %v (%s)\n", name, variable.Value, variable.Type)
		}
		scope = scope.parent
		level++
	}

	// Print functions
	result += fmt.Sprintf("  Functions: %d\n", len(sm.functions))
	for name := range sm.functions {
		result += fmt.Sprintf("    %s\n", name)
	}

	deadlineStr := "no deadline"
	if !sm.deadline.IsZero() {
		deadlineStr = sm.deadline.String()
	}
	result += fmt.Sprintf("  Deadline: %s\n", deadlineStr)
	result += "}"

	return result
}

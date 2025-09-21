// Package context provides execution context management
package context

import (
	"context"
	"fmt"
	"time"
)

// ExecutionContext represents an execution context with scope management
type ExecutionContext struct {
	// Go context for cancellation, timeout, and value storage
	Context context.Context

	// Cancel function to cancel the context
	Cancel context.CancelFunc

	// Scope manager for variable scope management
	ScopeManager *ScopeManager

	// Module name
	ModuleName string

	// Parent execution context
	Parent *ExecutionContext

	// Security context
	Security *SecurityContext

	// Debug mode
	Debug bool
}

// SecurityContext defines security restrictions
type SecurityContext struct {
	MaxExecutionTime  time.Duration
	MaxMemoryUsage    int64
	AllowedModules    []string
	ForbiddenKeywords []string
	AllowCrossModule  bool
}

// NewExecutionContext creates a new execution context
func NewExecutionContext() *ExecutionContext {
	ctx, cancel := context.WithCancel(context.Background())

	return &ExecutionContext{
		Context:      ctx,
		Cancel:       cancel,
		ScopeManager: NewScopeManager(),
		ModuleName:   "main",
		Parent:       nil,
		Security: &SecurityContext{
			MaxExecutionTime:  5 * time.Second,
			MaxMemoryUsage:    10 * 1024 * 1024, // 10MB
			AllowedModules:    []string{"fmt", "math"},
			ForbiddenKeywords: []string{"unsafe"},
			AllowCrossModule:  true,
		},
		Debug: false,
	}
}

// NewExecutionContextWithTimeout creates a new execution context with timeout
func NewExecutionContextWithTimeout(timeout time.Duration) *ExecutionContext {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	return &ExecutionContext{
		Context:      ctx,
		Cancel:       cancel,
		ScopeManager: NewScopeManagerWithTimeout(timeout),
		ModuleName:   "main",
		Parent:       nil,
		Security: &SecurityContext{
			MaxExecutionTime:  timeout,
			MaxMemoryUsage:    10 * 1024 * 1024, // 10MB
			AllowedModules:    []string{"fmt", "math"},
			ForbiddenKeywords: []string{"unsafe"},
			AllowCrossModule:  true,
		},
		Debug: false,
	}
}

// NewChildContext creates a new child execution context
func (ec *ExecutionContext) NewChildContext(moduleName string) *ExecutionContext {
	// Create a new context derived from parent
	childCtx, childCancel := context.WithCancel(ec.Context)

	// Create a new scope in the scope manager
	ec.ScopeManager.NewScope()

	return &ExecutionContext{
		Context:      childCtx,
		Cancel:       childCancel,
		ScopeManager: ec.ScopeManager, // Share the same scope manager
		ModuleName:   moduleName,
		Parent:       ec,
		Security:     ec.Security, // Share the same security context
		Debug:        ec.Debug,
	}
}

// ExitContext exits the current context and returns to parent
func (ec *ExecutionContext) ExitContext() *ExecutionContext {
	if ec.Parent != nil {
		// Exit the current scope
		ec.ScopeManager.ExitScope()
		return ec.Parent
	}
	return ec
}

// SetValue sets a value in the current scope
func (ec *ExecutionContext) SetValue(name string, value interface{}) error {
	// Determine the type of the value
	var varType string
	switch value.(type) {
	case int:
		varType = "int"
	case float64:
		varType = "float64"
	case string:
		varType = "string"
	case bool:
		varType = "bool"
	default:
		varType = "unknown"
	}

	// Set the variable in the scope manager
	ec.ScopeManager.SetVariable(name, value, varType)

	// Debug output
	if ec.Debug {
		fmt.Printf("Set variable: %s = %v (type: %s)\n", name, value, varType)
	}

	return nil
}

// GetValue gets a value from the scope chain
func (ec *ExecutionContext) GetValue(name string) (interface{}, bool) {
	if variable, exists := ec.ScopeManager.GetVariable(name); exists {
		// Debug output
		if ec.Debug {
			fmt.Printf("Get variable: %s = %v (type: %s)\n", name, variable.Value, variable.Type)
		}
		return variable.Value, true
	}
	return nil, false
}

// GetVariable gets a variable with type information
func (ec *ExecutionContext) GetVariable(name string) *Variable {
	variable, _ := ec.ScopeManager.GetVariable(name)
	return variable
}

// UpdateValue updates an existing value in the scope chain
func (ec *ExecutionContext) UpdateValue(name string, value interface{}) error {
	// Debug output
	if ec.Debug {
		fmt.Printf("Update variable: %s = %v\n", name, value)
	}
	return ec.ScopeManager.UpdateVariable(name, value)
}

// DeleteValue deletes a value from the current scope
func (ec *ExecutionContext) DeleteValue(name string) {
	// Debug output
	if ec.Debug {
		fmt.Printf("Delete variable: %s\n", name)
	}
	ec.ScopeManager.DeleteVariable(name)
}

// SetGlobalValue sets a global value in the root context
func (ec *ExecutionContext) SetGlobalValue(name string, value interface{}) error {
	// For global values, we set them in the root context
	rootEC := ec.GetRootContext()
	return rootEC.SetValue(name, value)
}

// GetRootContext gets the root execution context
func (ec *ExecutionContext) GetRootContext() *ExecutionContext {
	current := ec
	for current.Parent != nil {
		current = current.Parent
	}
	return current
}

// CheckModuleAccess checks if module access is allowed
func (ec *ExecutionContext) CheckModuleAccess(moduleName string) error {
	if ec.Security == nil {
		return nil
	}

	// If no allowed modules specified, allow all
	if len(ec.Security.AllowedModules) == 0 {
		return nil
	}

	// Check if module is in allowed list
	for _, allowed := range ec.Security.AllowedModules {
		if allowed == moduleName {
			return nil
		}
	}

	return fmt.Errorf("access to module %s is not allowed", moduleName)
}

// IsCancelled checks if the context has been cancelled
func (ec *ExecutionContext) IsCancelled() bool {
	select {
	case <-ec.Context.Done():
		return true
	default:
		return false
	}
}

// Deadline returns the deadline of the context
func (ec *ExecutionContext) Deadline() (time.Time, bool) {
	return ec.Context.Deadline()
}

// Err returns the error associated with the context
func (ec *ExecutionContext) Err() error {
	return ec.Context.Err()
}

// String returns a string representation of the execution context
func (ec *ExecutionContext) String() string {
	deadline, hasDeadline := ec.Deadline()
	if hasDeadline {
		return fmt.Sprintf("ExecutionContext{module: %s, deadline: %v}", ec.ModuleName, deadline)
	}
	return fmt.Sprintf("ExecutionContext{module: %s}", ec.ModuleName)
}

// RegisterFunction registers a function in the execution context
func (ec *ExecutionContext) RegisterFunction(name string, fn Function) error {
	ec.ScopeManager.RegisterFunction(name, fn)

	// Debug output
	if ec.Debug {
		fmt.Printf("Registered function: %s\n", name)
	}

	return nil
}

// GetFunction gets a function by name from the execution context
func (ec *ExecutionContext) GetFunction(name string) (Function, bool) {
	return ec.ScopeManager.GetFunction(name)
}

// GetAllFunctions returns all registered functions
func (ec *ExecutionContext) GetAllFunctions() map[string]Function {
	return ec.ScopeManager.GetAllFunctions()
}

// WithTimeout creates a new execution context with timeout
func (ec *ExecutionContext) WithTimeout(timeout time.Duration) *ExecutionContext {
	// Create a new context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	// Create a new execution context with the same scope manager
	return &ExecutionContext{
		Context:      ctx,
		Cancel:       cancel,
		ScopeManager: ec.ScopeManager, // Share the same scope manager
		ModuleName:   ec.ModuleName,
		Parent:       ec.Parent,
		Security:     ec.Security,
		Debug:        ec.Debug,
	}
}

// SetDebug enables or disables debug mode
func (ec *ExecutionContext) SetDebug(debug bool) {
	ec.Debug = debug
}

// GetScopeManager returns the scope manager
func (ec *ExecutionContext) GetScopeManager() *ScopeManager {
	return ec.ScopeManager
}

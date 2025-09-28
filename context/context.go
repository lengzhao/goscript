// Package vm provides the virtual machine implementation with context-based scope management
package context

import (
	"fmt"
)

// Context represents an execution context with hierarchical scope management
type Context struct {
	// Path key for identifying the context (e.g., "main.function.loop")
	pathKey string

	// Parent context reference
	parent *Context

	// Variables in this context
	variables map[string]interface{}

	// Variable types in this context
	types map[string]string

	// Child contexts
	children map[string]*Context
}

// NewContext creates a new context with the given path key and parent
func NewContext(pathKey string, parent *Context) *Context {
	return &Context{
		pathKey:   pathKey,
		parent:    parent,
		variables: make(map[string]interface{}),
		types:     make(map[string]string),
		children:  make(map[string]*Context),
	}
}

// GetPathKey returns the path key of this context
func (ctx *Context) GetPathKey() string {
	return ctx.pathKey
}

// GetParent returns the parent context
func (ctx *Context) GetParent() *Context {
	return ctx.parent
}

// GetVariable searches for a variable in the context hierarchy
// It first checks the current context, then recursively checks parent contexts
func (ctx *Context) GetVariable(name string) (interface{}, bool) {
	// First check current context
	if value, exists := ctx.variables[name]; exists {
		return value, true
	}

	// Then check parent context if exists
	if ctx.parent != nil {
		return ctx.parent.GetVariable(name)
	}

	return nil, false
}

// MustGetVariable gets a variable, panics if not found in the context hierarchy
func (ctx *Context) MustGetVariable(name string) interface{} {
	if value, exists := ctx.variables[name]; exists {
		return value
	}

	if ctx.parent != nil {
		return ctx.parent.MustGetVariable(name)
	}

	panic(fmt.Sprintf("variable %s not found in context hierarchy", name))
}

// SetVariable sets a variable in the current context
func (ctx *Context) SetVariable(name string, value interface{}) error {
	if _, exists := ctx.variables[name]; exists {
		ctx.variables[name] = value
		return nil
	}
	if ctx.parent == nil {
		return fmt.Errorf("variable %s not found in context hierarchy", name)
	}
	return ctx.parent.SetVariable(name, value)
}

// CreateVariableWithType sets a variable with its type in the current context
func (ctx *Context) CreateVariableWithType(name string, value interface{}, varType string) error {
	if _, exists := ctx.variables[name]; exists {
		return fmt.Errorf("variable %s already exists", name)
	}
	ctx.variables[name] = value
	ctx.types[name] = varType
	return nil
}

// GetVariableType gets the type of a variable in the context hierarchy
func (ctx *Context) GetVariableType(name string) (string, bool) {
	// First check current context
	if varType, exists := ctx.types[name]; exists {
		return varType, true
	}

	// Then check parent context if exists
	if ctx.parent != nil {
		return ctx.parent.GetVariableType(name)
	}

	return "", false
}

// HasVariable checks if a variable exists in the current context (not in hierarchy)
func (ctx *Context) HasVariable(name string) bool {
	_, exists := ctx.variables[name]
	return exists
}

// DeleteVariable removes a variable from the current context
func (ctx *Context) DeleteVariable(name string) {
	delete(ctx.variables, name)
	delete(ctx.types, name)
}

// GetAllVariables returns all variables in the current context
func (ctx *Context) GetAllVariables() map[string]interface{} {
	// Return a copy to prevent external modification
	result := make(map[string]interface{})
	for k, v := range ctx.variables {
		result[k] = v
	}
	return result
}

// GetAllVariablesWithTypes returns all variables with their types in the current context
func (ctx *Context) GetAllVariablesWithTypes() (map[string]interface{}, map[string]string) {
	// Return copies to prevent external modification
	vars := make(map[string]interface{})
	types := make(map[string]string)

	for k, v := range ctx.variables {
		vars[k] = v
	}

	for k, t := range ctx.types {
		types[k] = t
	}

	return vars, types
}

// AddChild adds a child context
func (ctx *Context) AddChild(child *Context) {
	ctx.children[child.pathKey] = child
}

// RemoveChild removes a child context
func (ctx *Context) RemoveChild(pathKey string) {
	delete(ctx.children, pathKey)
}

// GetChild gets a child context by path key
func (ctx *Context) GetChild(pathKey string) (*Context, bool) {
	child, exists := ctx.children[pathKey]
	return child, exists
}

// GetChildren returns all child contexts
func (ctx *Context) GetChildren() map[string]*Context {
	// Return a copy to prevent external modification
	result := make(map[string]*Context)
	for k, v := range ctx.children {
		result[k] = v
	}
	return result
}

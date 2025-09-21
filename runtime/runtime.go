// Package runtime provides runtime support for the GoScript engine
package runtime

import (
	"fmt"
	"reflect"

	"github.com/lengzhao/goscript/types"
)

// Runtime manages the execution environment for scripts
type Runtime struct {
	// Variables stores script variables
	Variables map[string]interface{}

	// Functions stores built-in and user-defined functions
	Functions map[string]Function

	// Types stores built-in and user-defined types
	Types map[string]reflect.Type

	// Modules stores imported modules
	Modules map[string]Module

	// Type system
	TypeSystem map[string]types.IType

	// Debug mode
	Debug bool
}

// NewRuntime creates a new Runtime instance
func NewRuntime() *Runtime {
	// Initialize with basic types
	typeSystem := make(map[string]types.IType)
	typeSystem["int"] = types.IntType.Clone()
	typeSystem["float64"] = types.Float64Type.Clone()
	typeSystem["string"] = types.StringType.Clone()
	typeSystem["bool"] = types.BoolType.Clone()

	return &Runtime{
		Variables:  make(map[string]interface{}),
		Functions:  make(map[string]Function),
		Types:      make(map[string]reflect.Type),
		Modules:    make(map[string]Module),
		TypeSystem: typeSystem,
		Debug:      false,
	}
}

// SetVariable sets a variable in the runtime
func (r *Runtime) SetVariable(name string, value interface{}) {
	r.Variables[name] = value

	// Debug output
	if r.Debug {
		fmt.Printf("Runtime: Set variable %s = %v\n", name, value)
	}
}

// GetVariable gets a variable from the runtime
func (r *Runtime) GetVariable(name string) (interface{}, bool) {
	value, ok := r.Variables[name]

	// Debug output
	if r.Debug && ok {
		fmt.Printf("Runtime: Get variable %s = %v\n", name, value)
	}

	return value, ok
}

// DeleteVariable deletes a variable from the runtime
func (r *Runtime) DeleteVariable(name string) {
	delete(r.Variables, name)

	// Debug output
	if r.Debug {
		fmt.Printf("Runtime: Delete variable %s\n", name)
	}
}

// RegisterFunction registers a function in the runtime
func (r *Runtime) RegisterFunction(name string, fn Function) {
	r.Functions[name] = fn

	// Debug output
	if r.Debug {
		fmt.Printf("Runtime: Register function %s\n", name)
	}
}

// GetFunction gets a function from the runtime
func (r *Runtime) GetFunction(name string) (Function, bool) {
	fn, ok := r.Functions[name]
	return fn, ok
}

// RegisterType registers a type in the runtime
func (r *Runtime) RegisterType(name string, typ reflect.Type) {
	r.Types[name] = typ

	// Debug output
	if r.Debug {
		fmt.Printf("Runtime: Register type %s\n", name)
	}
}

// GetType gets a type from the runtime
func (r *Runtime) GetType(name string) (reflect.Type, bool) {
	typ, ok := r.Types[name]
	return typ, ok
}

// ImportModule imports a module into the runtime
func (r *Runtime) ImportModule(module Module) error {
	name := module.Name()
	if _, exists := r.Modules[name]; exists {
		return fmt.Errorf("module %s already imported", name)
	}

	r.Modules[name] = module

	// Debug output
	if r.Debug {
		fmt.Printf("Runtime: Import module %s\n", name)
	}

	return nil
}

// GetModule gets a module from the runtime
func (r *Runtime) GetModule(name string) (Module, bool) {
	module, ok := r.Modules[name]
	return module, ok
}

// RegisterIType registers an IType in the runtime
func (r *Runtime) RegisterIType(name string, typ types.IType) {
	r.TypeSystem[name] = typ

	// Debug output
	if r.Debug {
		fmt.Printf("Runtime: Register IType %s\n", name)
	}
}

// GetIType gets an IType from the runtime
func (r *Runtime) GetIType(name string) (types.IType, bool) {
	typ, ok := r.TypeSystem[name]
	return typ, ok
}

// GetAllITypes returns all registered ITypes
func (r *Runtime) GetAllITypes() map[string]types.IType {
	return r.TypeSystem
}

// SetDebug enables or disables debug mode
func (r *Runtime) SetDebug(debug bool) {
	r.Debug = debug
}

// Function represents a callable function in the runtime
type Function func(args ...interface{}) (interface{}, error)

// Module represents a module in the runtime
type Module interface {
	// Name returns the module name
	Name() string

	// GetFunction returns a function by name
	GetFunction(name string) (Function, bool)

	// GetType returns a type by name
	GetType(name string) (reflect.Type, bool)
}

// Context holds the execution context for a script
type Context struct {
	Runtime *Runtime
	// Add other context information as needed
}

// NewContext creates a new execution context
func NewContext(runtime *Runtime) *Context {
	return &Context{
		Runtime: runtime,
	}
}

// String returns a string representation of the runtime
func (r *Runtime) String() string {
	return fmt.Sprintf("Runtime{variables: %d, functions: %d, types: %d, modules: %d}",
		len(r.Variables), len(r.Functions), len(r.Types), len(r.Modules))
}

// DebugString returns a detailed string representation for debugging
func (r *Runtime) DebugString() string {
	result := "Runtime{\n"
	result += fmt.Sprintf("  Variables: %d\n", len(r.Variables))
	for name, value := range r.Variables {
		result += fmt.Sprintf("    %s: %v\n", name, value)
	}
	result += fmt.Sprintf("  Functions: %d\n", len(r.Functions))
	for name := range r.Functions {
		result += fmt.Sprintf("    %s\n", name)
	}
	result += fmt.Sprintf("  Types: %d\n", len(r.Types))
	for name := range r.Types {
		result += fmt.Sprintf("    %s\n", name)
	}
	result += fmt.Sprintf("  Modules: %d\n", len(r.Modules))
	for name := range r.Modules {
		result += fmt.Sprintf("    %s\n", name)
	}
	result += fmt.Sprintf("  ITypes: %d\n", len(r.TypeSystem))
	for name := range r.TypeSystem {
		result += fmt.Sprintf("    %s\n", name)
	}
	result += "}"
	return result
}

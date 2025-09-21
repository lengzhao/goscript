// Package vm provides the virtual machine implementation
package vm

// ModuleManagerInterface defines the interface for module management
type ModuleManagerInterface interface {
	// CallModuleFunction calls a function in a specific module
	CallModuleFunction(moduleName, functionName string, args ...interface{}) (interface{}, error)
}
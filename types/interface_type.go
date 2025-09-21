// Package types defines the type system for GoScript
package types

import (
	"fmt"
	"reflect"
	"strings"
)

// InterfaceType represents an interface type
type InterfaceType struct {
	name     string
	methods  map[string]Method
	embedded []IType
}

// NewInterfaceType creates a new interface type
func NewInterfaceType(name string) *InterfaceType {
	return &InterfaceType{
		name:     name,
		methods:  make(map[string]Method),
		embedded: make([]IType, 0),
	}
}

// TypeName returns the name of the type
func (it *InterfaceType) TypeName() string {
	return it.name
}

// String returns the string representation of the type
func (it *InterfaceType) String() string {
	if len(it.methods) == 0 && len(it.embedded) == 0 {
		return "interface{}"
	}

	var parts []string

	// Add embedded interfaces
	for _, embedded := range it.embedded {
		parts = append(parts, embedded.TypeName())
	}

	// Add methods from embedded interfaces
	for _, embedded := range it.embedded {
		embeddedMethods := embedded.GetMethods()
		for _, method := range embeddedMethods {
			paramStrs := make([]string, len(method.Params))
			for i, param := range method.Params {
				paramStrs[i] = param.String()
			}

			returnStrs := make([]string, len(method.Returns))
			for i, ret := range method.Returns {
				returnStrs[i] = ret.String()
			}

			var methodStr string
			if len(returnStrs) == 0 {
				methodStr = fmt.Sprintf("%s(%s)", method.Name, strings.Join(paramStrs, ", "))
			} else if len(returnStrs) == 1 {
				methodStr = fmt.Sprintf("%s(%s) %s", method.Name, strings.Join(paramStrs, ", "), returnStrs[0])
			} else {
				methodStr = fmt.Sprintf("%s(%s) (%s)", method.Name, strings.Join(paramStrs, ", "), strings.Join(returnStrs, ", "))
			}

			parts = append(parts, methodStr)
		}
	}

	// Add methods
	for _, method := range it.methods {
		paramStrs := make([]string, len(method.Params))
		for i, param := range method.Params {
			paramStrs[i] = param.String()
		}

		returnStrs := make([]string, len(method.Returns))
		for i, ret := range method.Returns {
			returnStrs[i] = ret.String()
		}

		var methodStr string
		if len(returnStrs) == 0 {
			methodStr = fmt.Sprintf("%s(%s)", method.Name, strings.Join(paramStrs, ", "))
		} else if len(returnStrs) == 1 {
			methodStr = fmt.Sprintf("%s(%s) %s", method.Name, strings.Join(paramStrs, ", "), returnStrs[0])
		} else {
			methodStr = fmt.Sprintf("%s(%s) (%s)", method.Name, strings.Join(paramStrs, ", "), strings.Join(returnStrs, ", "))
		}

		parts = append(parts, methodStr)
	}

	return fmt.Sprintf("interface { %s }", strings.Join(parts, "; "))
}

// Equals compares this type with another type
func (it *InterfaceType) Equals(other IType) bool {
	if other == nil {
		return false
	}

	otherInterface, ok := other.(*InterfaceType)
	if !ok {
		return false
	}

	if it.name != otherInterface.name {
		return false
	}

	if len(it.methods) != len(otherInterface.methods) {
		return false
	}

	if len(it.embedded) != len(otherInterface.embedded) {
		return false
	}

	// Compare methods
	for name, method := range it.methods {
		otherMethod, exists := otherInterface.methods[name]
		if !exists {
			return false
		}

		if len(method.Params) != len(otherMethod.Params) {
			return false
		}

		if len(method.Returns) != len(otherMethod.Returns) {
			return false
		}

		for i, param := range method.Params {
			if !param.Equals(otherMethod.Params[i]) {
				return false
			}
		}

		for i, ret := range method.Returns {
			if !ret.Equals(otherMethod.Returns[i]) {
				return false
			}
		}
	}

	// Compare embedded interfaces
	for i, embedded := range it.embedded {
		if !embedded.Equals(otherInterface.embedded[i]) {
			return false
		}
	}

	return true
}

// Size returns the size of the type in bytes
func (it *InterfaceType) Size() int {
	// Interface types have a fixed size for method table pointers
	return 16
}

// Clone creates a copy of the type
func (it *InterfaceType) Clone() IType {
	clone := &InterfaceType{
		name:     it.name,
		methods:  make(map[string]Method),
		embedded: make([]IType, len(it.embedded)),
	}

	for name, method := range it.methods {
		paramClones := make([]IType, len(method.Params))
		for i, param := range method.Params {
			paramClones[i] = param.Clone()
		}

		returnClones := make([]IType, len(method.Returns))
		for i, ret := range method.Returns {
			returnClones[i] = ret.Clone()
		}

		clone.methods[name] = Method{
			Name:    method.Name,
			Params:  paramClones,
			Returns: returnClones,
		}
	}

	for i, embedded := range it.embedded {
		clone.embedded[i] = embedded.Clone()
	}

	return clone
}

// DefaultValue returns the default value for this type
func (it *InterfaceType) DefaultValue() interface{} {
	// Interface types have nil as default value
	return nil
}

// Kind returns the reflect.Kind of the type
func (it *InterfaceType) Kind() reflect.Kind {
	return reflect.Interface
}

// AddMethod adds a method to the interface
func (it *InterfaceType) AddMethod(name string, params []IType, returns []IType) {
	it.methods[name] = Method{
		Name:    name,
		Params:  params,
		Returns: returns,
	}
}

// AddEmbedded adds an embedded interface
func (it *InterfaceType) AddEmbedded(embedded IType) {
	it.embedded = append(it.embedded, embedded)
}

// GetEmbedded returns all embedded interfaces
func (it *InterfaceType) GetEmbedded() []IType {
	return it.embedded
}

// Implements checks if a type implements this interface
func (it *InterfaceType) Implements(typ IType) bool {
	// For now, we'll implement a simple check
	// In a real implementation, this would check if the type has all required methods
	switch t := typ.(type) {
	case *StructType:
		// Check if struct has all methods required by interface
		for methodName := range it.methods {
			// In a real implementation, we would check if the struct has a method
			// with the same name and signature
			_ = methodName
		}
		return true
	case *InterfaceType:
		// Check if interface embeds or has all methods of this interface
		for methodName, method := range it.methods {
			if otherMethod, exists := t.methods[methodName]; exists {
				// Check if method signatures match
				if len(method.Params) != len(otherMethod.Params) {
					return false
				}
				if len(method.Returns) != len(otherMethod.Returns) {
					return false
				}

				for i, param := range method.Params {
					if !param.Equals(otherMethod.Params[i]) {
						return false
					}
				}

				for i, ret := range method.Returns {
					if !ret.Equals(otherMethod.Returns[i]) {
						return false
					}
				}
			} else {
				// Check embedded interfaces
				found := false
				for _, embedded := range t.embedded {
					if embedded.HasMethod(methodName) {
						found = true
						break
					}
				}
				if !found {
					return false
				}
			}
		}
		return true
	default:
		return false
	}
}

// GetMethods returns all methods available on this type
func (it *InterfaceType) GetMethods() []Method {
	methods := make([]Method, 0, len(it.methods))

	// Add methods defined directly in this interface
	for _, method := range it.methods {
		methods = append(methods, method)
	}

	// Add methods from embedded interfaces
	for _, embedded := range it.embedded {
		embeddedMethods := embedded.GetMethods()
		methods = append(methods, embeddedMethods...)
	}

	return methods
}

// HasMethod checks if the type has a method with the given name
// This checks both directly defined methods and methods from embedded interfaces
func (it *InterfaceType) HasMethod(name string) bool {
	// Check directly defined methods
	if _, exists := it.methods[name]; exists {
		return true
	}

	// Check embedded interfaces
	for _, embedded := range it.embedded {
		if embedded.HasMethod(name) {
			return true
		}
	}

	return false
}

// GetMethod returns a method by name
// This checks both directly defined methods and methods from embedded interfaces
func (it *InterfaceType) GetMethod(name string) (Method, bool) {
	// Check directly defined methods
	if method, exists := it.methods[name]; exists {
		return method, true
	}

	// Check embedded interfaces
	for _, embedded := range it.embedded {
		if method, exists := embedded.GetMethod(name); exists {
			return method, true
		}
	}

	return Method{}, false
}

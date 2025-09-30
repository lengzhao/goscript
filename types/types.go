// Package types defines the type system for GoScript
package types

import (
	"reflect"
)

// Function represents a callable function in GoScript
type Function func(args ...interface{}) (interface{}, error)

// Method represents a method signature
type Method struct {
	Name    string
	Params  []IType
	Returns []IType
}

// IType is the interface that all types in GoScript must implement
type IType interface {
	// TypeName returns the name of the type
	TypeName() string

	// String returns the string representation of the type
	String() string

	// Equals compares this type with another type
	Equals(other IType) bool

	// Size returns the size of the type in bytes
	Size() int

	// Clone creates a copy of the type
	Clone() IType

	// DefaultValue returns the default value for this type
	DefaultValue() interface{}

	// Kind returns the reflect.Kind of the type
	Kind() reflect.Kind

	// GetMethods returns all methods available on this type
	GetMethods() []Method

	// HasMethod checks if the type has a method with the given name
	HasMethod(name string) bool

	// GetMethod returns a method by name
	GetMethod(name string) (Method, bool)
}

// Package types defines the type system for GoScript
package types

import (
	"fmt"
	"reflect"
)

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
}

// BaseType represents a basic type
type BaseType struct {
	name string
	size int
	kind reflect.Kind
}

// NewBaseType creates a new basic type
func NewBaseType(name string, size int, kind reflect.Kind) *BaseType {
	return &BaseType{
		name: name,
		size: size,
		kind: kind,
	}
}

// TypeName returns the name of the type
func (bt *BaseType) TypeName() string {
	return bt.name
}

// String returns the string representation of the type
func (bt *BaseType) String() string {
	return bt.name
}

// Equals compares this type with another type
func (bt *BaseType) Equals(other IType) bool {
	if other == nil {
		return false
	}
	return bt.name == other.TypeName()
}

// Size returns the size of the type in bytes
func (bt *BaseType) Size() int {
	return bt.size
}

// Clone creates a copy of the type
func (bt *BaseType) Clone() IType {
	return &BaseType{
		name: bt.name,
		size: bt.size,
		kind: bt.kind,
	}
}

// DefaultValue returns the default value for this type
func (bt *BaseType) DefaultValue() interface{} {
	switch bt.kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int(0)
	case reflect.Float32, reflect.Float64:
		return float64(0.0)
	case reflect.String:
		return ""
	case reflect.Bool:
		return false
	default:
		return nil
	}
}

// Kind returns the reflect.Kind of the type
func (bt *BaseType) Kind() reflect.Kind {
	return bt.kind
}

// IntType represents the int type
var IntType = NewBaseType("int", 8, reflect.Int)

// Float64Type represents the float64 type
var Float64Type = NewBaseType("float64", 8, reflect.Float64)

// StringType represents the string type
var StringType = NewBaseType("string", 16, reflect.String)

// BoolType represents the bool type
var BoolType = NewBaseType("bool", 1, reflect.Bool)

// GetTypeByName returns a type by its name
func GetTypeByName(name string) (IType, error) {
	switch name {
	case "int":
		return IntType.Clone(), nil
	case "float64":
		return Float64Type.Clone(), nil
	case "string":
		return StringType.Clone(), nil
	case "bool":
		return BoolType.Clone(), nil
	default:
		return nil, fmt.Errorf("unknown type: %s", name)
	}
}

// TypeFromReflect returns a IType from a reflect.Type
func TypeFromReflect(rt reflect.Type) IType {
	switch rt.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return IntType.Clone()
	case reflect.Float32, reflect.Float64:
		return Float64Type.Clone()
	case reflect.String:
		return StringType.Clone()
	case reflect.Bool:
		return BoolType.Clone()
	default:
		return NewBaseType(rt.Name(), int(rt.Size()), rt.Kind())
	}
}

// IsNumeric checks if a type is numeric
func IsNumeric(t IType) bool {
	if t == nil {
		return false
	}

	kind := t.Kind()
	return kind >= reflect.Int && kind <= reflect.Float64
}

// IsComparable checks if two types are comparable
func IsComparable(t1, t2 IType) bool {
	if t1 == nil || t2 == nil {
		return false
	}

	// Same types are always comparable
	if t1.TypeName() == t2.TypeName() {
		return true
	}

	// Numeric types are comparable
	if IsNumeric(t1) && IsNumeric(t2) {
		return true
	}

	return false
}

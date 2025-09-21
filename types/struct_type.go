// Package types defines the type system for GoScript
package types

import (
	"fmt"
	"reflect"
	"strings"
)

// StructType represents a struct type
type StructType struct {
	name   string
	fields map[string]IType
	tags   map[string]string
	// Store methods associated with this struct type
	methods map[string]Method
}

// NewStructType creates a new struct type
func NewStructType(name string) *StructType {
	return &StructType{
		name:    name,
		fields:  make(map[string]IType),
		tags:    make(map[string]string),
		methods: make(map[string]Method),
	}
}

// TypeName returns the name of the type
func (st *StructType) TypeName() string {
	return st.name
}

// String returns the string representation of the type
func (st *StructType) String() string {
	var fields []string
	for name, typ := range st.fields {
		tag := ""
		if t, exists := st.tags[name]; exists && t != "" {
			tag = fmt.Sprintf(" `%s`", t)
		}
		fields = append(fields, fmt.Sprintf("%s %s%s", name, typ.String(), tag))
	}
	return fmt.Sprintf("struct { %s }", strings.Join(fields, "; "))
}

// Equals compares this type with another type
func (st *StructType) Equals(other IType) bool {
	if other == nil {
		return false
	}

	otherStruct, ok := other.(*StructType)
	if !ok {
		return false
	}

	if st.name != otherStruct.name {
		return false
	}

	if len(st.fields) != len(otherStruct.fields) {
		return false
	}

	for name, typ := range st.fields {
		otherTyp, exists := otherStruct.fields[name]
		if !exists || !typ.Equals(otherTyp) {
			return false
		}
	}

	return true
}

// Size returns the size of the type in bytes
func (st *StructType) Size() int {
	// For simplicity, we return a fixed size for struct types
	// In a real implementation, this would be calculated based on fields
	return 24
}

// Clone creates a copy of the type
func (st *StructType) Clone() IType {
	clone := &StructType{
		name:    st.name,
		fields:  make(map[string]IType),
		tags:    make(map[string]string),
		methods: make(map[string]Method),
	}

	for name, typ := range st.fields {
		clone.fields[name] = typ.Clone()
	}

	for name, tag := range st.tags {
		clone.tags[name] = tag
	}

	for name, method := range st.methods {
		clone.methods[name] = method
	}

	return clone
}

// DefaultValue returns the default value for this type
func (st *StructType) DefaultValue() interface{} {
	// Return a map representing the struct with default values
	values := make(map[string]interface{})
	for name, typ := range st.fields {
		values[name] = typ.DefaultValue()
	}
	return values
}

// Kind returns the reflect.Kind of the type
func (st *StructType) Kind() reflect.Kind {
	return reflect.Struct
}

// AddField adds a field to the struct
func (st *StructType) AddField(name string, typ IType, tag string) {
	st.fields[name] = typ
	if tag != "" {
		st.tags[name] = tag
	}
}

// GetField returns the type of a field
func (st *StructType) GetField(name string) (IType, bool) {
	typ, exists := st.fields[name]
	return typ, exists
}

// HasField checks if the struct has a field with the given name
func (st *StructType) HasField(name string) bool {
	_, exists := st.fields[name]
	return exists
}

// GetFields returns all fields of the struct
func (st *StructType) GetFields() map[string]IType {
	return st.fields
}

// GetFieldNames returns the names of all fields
func (st *StructType) GetFieldNames() []string {
	names := make([]string, 0, len(st.fields))
	for name := range st.fields {
		names = append(names, name)
	}
	return names
}

// GetMethods returns all methods available on this type
func (st *StructType) GetMethods() []Method {
	methods := make([]Method, 0, len(st.methods))
	for _, method := range st.methods {
		methods = append(methods, method)
	}
	return methods
}

// HasMethod checks if the type has a method with the given name
func (st *StructType) HasMethod(name string) bool {
	_, exists := st.methods[name]
	return exists
}

// GetMethod returns a method by name
func (st *StructType) GetMethod(name string) (Method, bool) {
	method, exists := st.methods[name]
	return method, exists
}

// AddMethod adds a method to the struct type
func (st *StructType) AddMethod(name string, params []IType, returns []IType) {
	st.methods[name] = Method{
		Name:    name,
		Params:  params,
		Returns: returns,
	}
}

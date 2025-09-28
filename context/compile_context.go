// Package context provides the compile context implementation for organizing instructions during compilation
package context

import (
	"github.com/lengzhao/goscript/instruction"
)

// CompileContext represents a compile context for organizing instructions during compilation
type CompileContext struct {
	// Path key for identifying the context (e.g., "main.function.block1")
	pathKey string

	// Parent context reference
	parent *CompileContext

	// Instructions mapped by key
	instructions map[string][]*instruction.Instruction

	// Child contexts
	children map[string]*CompileContext
}

// NewCompileContext creates a new compile context with the given path key and parent
func NewCompileContext(pathKey string, parent *CompileContext) *CompileContext {
	return &CompileContext{
		pathKey:      pathKey,
		parent:       parent,
		instructions: make(map[string][]*instruction.Instruction),
		children:     make(map[string]*CompileContext),
	}
}

// GetPathKey returns the path key of this compile context
func (ctx *CompileContext) GetPathKey() string {
	return ctx.pathKey
}

// GetParent returns the parent compile context
func (ctx *CompileContext) GetParent() *CompileContext {
	return ctx.parent
}

// SetInstructions sets instructions for a specific key
func (ctx *CompileContext) SetInstructions(key string, instructions []*instruction.Instruction) {
	ctx.instructions[key] = instructions
}

// GetInstructions gets instructions for a specific key
func (ctx *CompileContext) GetInstructions(key string) ([]*instruction.Instruction, bool) {
	instructions, exists := ctx.instructions[key]
	return instructions, exists
}

// GetAllInstructions returns all instructions mapped by key
func (ctx *CompileContext) GetAllInstructions() map[string][]*instruction.Instruction {
	// Return a copy to prevent external modification
	result := make(map[string][]*instruction.Instruction)
	for k, v := range ctx.instructions {
		// Create a copy of the instruction slice
		instrCopy := make([]*instruction.Instruction, len(v))
		copy(instrCopy, v)
		result[k] = instrCopy
	}
	return result
}

// AddChild adds a child compile context
func (ctx *CompileContext) AddChild(child *CompileContext) {
	ctx.children[child.pathKey] = child
}

// RemoveChild removes a child compile context
func (ctx *CompileContext) RemoveChild(pathKey string) {
	delete(ctx.children, pathKey)
}

// GetChild gets a child compile context by path key
func (ctx *CompileContext) GetChild(pathKey string) (*CompileContext, bool) {
	child, exists := ctx.children[pathKey]
	return child, exists
}

// GetChildren returns all child compile contexts
func (ctx *CompileContext) GetChildren() map[string]*CompileContext {
	// Return a copy to prevent external modification
	result := make(map[string]*CompileContext)
	for k, v := range ctx.children {
		result[k] = v
	}
	return result
}
// Package vm provides optimized stack operations for the virtual machine
package vm

import (
	"fmt"
)

// Stack represents an optimized stack for the virtual machine
type Stack struct {
	// Pre-allocated slice for better performance
	data []interface{}

	// Current stack size (top index + 1)
	size int

	// Maximum stack size for bounds checking
	maxSize int

	// Pre-allocated capacity to reduce reallocations
	capacity int
}

// NewStack creates a new optimized stack
func NewStack(initialCapacity, maxSize int) *Stack {
	if initialCapacity <= 0 {
		initialCapacity = 64 // Default initial capacity
	}
	if maxSize <= 0 {
		maxSize = 10000 // Default maximum size
	}

	return &Stack{
		data:     make([]interface{}, initialCapacity),
		size:     0,
		maxSize:  maxSize,
		capacity: initialCapacity,
	}
}

// Push pushes a value onto the stack
func (s *Stack) Push(value interface{}) error {
	// Check if we need to expand the slice
	if s.size >= s.capacity {
		if s.size >= s.maxSize {
			return fmt.Errorf("stack overflow: maximum size %d exceeded", s.maxSize)
		}

		// Double the capacity, but don't exceed maxSize
		newCapacity := s.capacity * 2
		if newCapacity > s.maxSize {
			newCapacity = s.maxSize
		}

		// Reallocate with new capacity
		newData := make([]interface{}, newCapacity)
		copy(newData, s.data[:s.size])
		s.data = newData
		s.capacity = newCapacity
	}

	// Add the value
	s.data[s.size] = value
	s.size++
	return nil
}

// Pop pops a value from the stack
func (s *Stack) Pop() (interface{}, error) {
	if s.size == 0 {
		return nil, fmt.Errorf("stack underflow")
	}

	s.size--
	value := s.data[s.size]
	s.data[s.size] = nil // Help GC
	return value, nil
}

// Peek returns the top value without removing it
func (s *Stack) Peek() (interface{}, error) {
	if s.size == 0 {
		return nil, fmt.Errorf("stack underflow")
	}

	return s.data[s.size-1], nil
}

// PeekAt returns the value at the specified index (0 = top)
func (s *Stack) PeekAt(index int) (interface{}, error) {
	if index < 0 || index >= s.size {
		return nil, fmt.Errorf("stack index out of range: %d (size: %d)", index, s.size)
	}

	return s.data[s.size-1-index], nil
}

// Size returns the current stack size
func (s *Stack) Size() int {
	return s.size
}

// IsEmpty returns true if the stack is empty
func (s *Stack) IsEmpty() bool {
	return s.size == 0
}

// Clear clears the stack
func (s *Stack) Clear() {
	// Clear all elements to help GC
	for i := 0; i < s.size; i++ {
		s.data[i] = nil
	}
	s.size = 0
}

// GetSlice returns a slice of the current stack contents (for debugging)
func (s *Stack) GetSlice() []interface{} {
	if s.size == 0 {
		return nil
	}

	result := make([]interface{}, s.size)
	copy(result, s.data[:s.size])
	return result
}

// Rotate rotates the top n elements of the stack
// For n=3: [a, b, c] -> [b, c, a]
func (s *Stack) Rotate(n int) error {
	if n < 2 || n > s.size {
		return fmt.Errorf("invalid rotate count: %d (stack size: %d)", n, s.size)
	}

	// Rotate the top n elements
	// For [a, b, c] with n=3, we want [b, c, a]
	// This means moving the first element to the end
	start := s.size - n
	topElement := s.data[start]

	// Shift all elements one position to the left
	for i := 0; i < n-1; i++ {
		s.data[start+i] = s.data[start+i+1]
	}

	// Put the first element at the end
	s.data[start+n-1] = topElement

	return nil
}

// Swap swaps the top two elements
func (s *Stack) Swap() error {
	if s.size < 2 {
		return fmt.Errorf("stack underflow: need at least 2 elements for swap")
	}

	s.data[s.size-1], s.data[s.size-2] = s.data[s.size-2], s.data[s.size-1]
	return nil
}

// Dup duplicates the top element
func (s *Stack) Dup() error {
	if s.size == 0 {
		return fmt.Errorf("stack underflow: cannot duplicate empty stack")
	}

	value := s.data[s.size-1]
	return s.Push(value)
}

// GetCapacity returns the current capacity
func (s *Stack) GetCapacity() int {
	return s.capacity
}

// GetMaxSize returns the maximum allowed size
func (s *Stack) GetMaxSize() int {
	return s.maxSize
}

// String returns a string representation of the stack
func (s *Stack) String() string {
	if s.size == 0 {
		return "[]"
	}

	return fmt.Sprintf("%v", s.data[:s.size])
}

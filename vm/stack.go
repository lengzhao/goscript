// Package vm provides the virtual machine implementation with key-based instruction execution
package vm

// Stack represents a simple stack data structure for the VM
type Stack struct {
	data  []interface{}
	top   int // Index of the top element
	limit int // Maximum capacity to prevent unbounded growth
}

// NewStack creates a new stack
func NewStack() *Stack {
	capacity := 200
	return &Stack{
		data:  make([]interface{}, capacity),
		top:   -1, // -1 indicates empty stack
		limit: capacity,
	}
}

// Push adds an item to the top of the stack
func (s *Stack) Push(item interface{}) {
	// Check if we need to expand the stack
	if s.top+1 >= s.limit {
		// Double the capacity
		newLimit := s.limit * 2
		newData := make([]interface{}, newLimit)
		copy(newData, s.data)
		s.data = newData
		s.limit = newLimit
	}

	s.top++
	s.data[s.top] = item
}

// Pop removes and returns the top item from the stack
func (s *Stack) Pop() interface{} {
	if s.top < 0 {
		return nil
	}

	item := s.data[s.top]
	// Optionally clear the reference to help GC
	s.data[s.top] = nil
	s.top--

	return item
}

// Peek returns the top item from the stack without removing it
func (s *Stack) Peek() interface{} {
	if s.top < 0 {
		return nil
	}

	return s.data[s.top]
}

// Len returns the number of items in the stack
func (s *Stack) Len() int {
	return s.top + 1
}

// IsEmpty returns true if the stack is empty
func (s *Stack) IsEmpty() bool {
	return s.top < 0
}

// Clear removes all items from the stack
func (s *Stack) Clear() {
	// Clear references to help GC
	for i := 0; i <= s.top; i++ {
		s.data[i] = nil
	}
	s.top = -1
}

// Items returns a copy of the stack items for debugging
func (s *Stack) Items() []interface{} {
	if s.top < 0 {
		return []interface{}{}
	}
	items := make([]interface{}, s.top+1)
	copy(items, s.data[:s.top+1])
	return items
}

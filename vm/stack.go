// Package vm provides the virtual machine implementation with key-based instruction execution
package vm

// Stack represents a simple stack data structure for the VM
type Stack struct {
	data []interface{}
}

// NewStack creates a new stack
func NewStack() *Stack {
	return &Stack{
		data: make([]interface{}, 0),
	}
}

// Push adds an item to the top of the stack
func (s *Stack) Push(item interface{}) {
	s.data = append(s.data, item)
}

// Pop removes and returns the top item from the stack
func (s *Stack) Pop() interface{} {
	if len(s.data) == 0 {
		return nil
	}

	// Get the last item
	item := s.data[len(s.data)-1]

	// Remove the last item
	s.data = s.data[:len(s.data)-1]

	return item
}

// Peek returns the top item from the stack without removing it
func (s *Stack) Peek() interface{} {
	if len(s.data) == 0 {
		return nil
	}

	return s.data[len(s.data)-1]
}

// Len returns the number of items in the stack
func (s *Stack) Len() int {
	return len(s.data)
}

// IsEmpty returns true if the stack is empty
func (s *Stack) IsEmpty() bool {
	return len(s.data) == 0
}

// Clear removes all items from the stack
func (s *Stack) Clear() {
	s.data = s.data[:0]
}

// Items returns a copy of the stack items for debugging
func (s *Stack) Items() []interface{} {
	items := make([]interface{}, len(s.data))
	copy(items, s.data)
	return items
}

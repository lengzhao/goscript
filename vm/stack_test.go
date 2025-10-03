package vm

import (
	"testing"
)

func TestStackOperations(t *testing.T) {
	stack := NewStack()

	// Test empty stack
	if !stack.IsEmpty() {
		t.Error("New stack should be empty")
	}

	if stack.Len() != 0 {
		t.Error("New stack length should be 0")
	}

	if stack.Pop() != nil {
		t.Error("Pop from empty stack should return nil")
	}

	if stack.Peek() != nil {
		t.Error("Peek from empty stack should return nil")
	}

	// Test push and pop
	stack.Push(1)
	if stack.IsEmpty() {
		t.Error("Stack should not be empty after push")
	}

	if stack.Len() != 1 {
		t.Errorf("Stack length should be 1, got %d", stack.Len())
	}

	if stack.Peek() != 1 {
		t.Errorf("Peek should return 1, got %v", stack.Peek())
	}

	value := stack.Pop()
	if value != 1 {
		t.Errorf("Pop should return 1, got %v", value)
	}

	if !stack.IsEmpty() {
		t.Error("Stack should be empty after popping the only element")
	}

	// Test multiple operations
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	if stack.Len() != 3 {
		t.Errorf("Stack length should be 3, got %d", stack.Len())
	}

	if stack.Peek() != 3 {
		t.Errorf("Peek should return 3, got %v", stack.Peek())
	}

	// Pop all elements
	if stack.Pop() != 3 {
		t.Error("First pop should return 3")
	}

	if stack.Pop() != 2 {
		t.Error("Second pop should return 2")
	}

	if stack.Pop() != 1 {
		t.Error("Third pop should return 1")
	}

	if !stack.IsEmpty() {
		t.Error("Stack should be empty after popping all elements")
	}
}

func TestStackItems(t *testing.T) {
	stack := NewStack()

	// Test items on empty stack
	items := stack.Items()
	if len(items) != 0 {
		t.Error("Items should return empty slice for empty stack")
	}

	// Test items on populated stack
	stack.Push("a")
	stack.Push("b")
	stack.Push("c")

	items = stack.Items()
	if len(items) != 3 {
		t.Errorf("Items should return 3 elements, got %d", len(items))
	}

	if items[0] != "a" || items[1] != "b" || items[2] != "c" {
		t.Errorf("Items should return [a, b, c], got %v", items)
	}

	// Verify that Items returns a copy and doesn't affect the original stack
	items[0] = "modified"
	if stack.Peek() != "c" {
		t.Error("Modifying items copy should not affect original stack")
	}
}

func TestStackClear(t *testing.T) {
	stack := NewStack()
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	if stack.Len() != 3 {
		t.Errorf("Stack should have 3 elements, got %d", stack.Len())
	}

	stack.Clear()

	if !stack.IsEmpty() {
		t.Error("Stack should be empty after Clear")
	}

	if stack.Len() != 0 {
		t.Errorf("Stack length should be 0 after Clear, got %d", stack.Len())
	}
}

func BenchmarkStackPushPop(b *testing.B) {
	stack := NewStack()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		stack.Push(i)
		stack.Pop()
	}
}

func BenchmarkStackOperations(b *testing.B) {
	stack := NewStack()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		stack.Push(i)
	}

	for i := 0; i < b.N; i++ {
		stack.Pop()
	}
}
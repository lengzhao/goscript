package vm

import (
	"testing"
)

func TestStackRotate(t *testing.T) {
	stack := NewStack(10, 100)

	// Push three elements: a, b, c
	stack.Push("a")
	stack.Push("b")
	stack.Push("c")

	// Check initial state
	if stack.Size() != 3 {
		t.Errorf("Expected stack size 3, got %d", stack.Size())
	}

	// Rotate the top 3 elements: [a, b, c] -> [b, c, a]
	err := stack.Rotate(3)
	if err != nil {
		t.Errorf("Rotate failed: %v", err)
	}

	// Check final state
	if stack.Size() != 3 {
		t.Errorf("Expected stack size 3, got %d", stack.Size())
	}

	// Pop elements and check order
	c, _ := stack.Pop()
	b, _ := stack.Pop()
	a, _ := stack.Pop()

	if a != "b" || b != "c" || c != "a" {
		t.Errorf("Expected [b, c, a], got [%v, %v, %v]", a, b, c)
	}
}

func TestStackRotateFieldAssignment(t *testing.T) {
	stack := NewStack(10, 100)

	// Simulate field assignment case: [value, object, fieldName]
	stack.Push("value")
	stack.Push("object")
	stack.Push("fieldName")

	// Check initial state
	if stack.Size() != 3 {
		t.Errorf("Expected stack size 3, got %d", stack.Size())
	}

	// Rotate the top 3 elements: [value, object, fieldName] -> [object, fieldName, value]
	err := stack.Rotate(3)
	if err != nil {
		t.Errorf("Rotate failed: %v", err)
	}

	// Check final state
	if stack.Size() != 3 {
		t.Errorf("Expected stack size 3, got %d", stack.Size())
	}

	// Pop elements and check order
	value, _ := stack.Pop()
	fieldName, _ := stack.Pop()
	object, _ := stack.Pop()

	if object != "object" || fieldName != "fieldName" || value != "value" {
		t.Errorf("Expected [object, fieldName, value], got [%v, %v, %v]", object, fieldName, value)
	}
}

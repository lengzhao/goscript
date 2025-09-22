package vm

import (
	"testing"
)

// BenchmarkStackOperations benchmarks the optimized stack operations
func BenchmarkStackOperations(b *testing.B) {
	// Benchmark push operations
	b.Run("Push", func(b *testing.B) {
		vm := NewVM()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			vm.Push(i % 100) // Small limit to prevent overflow
		}
	})

	// Benchmark pop operations
	b.Run("Pop", func(b *testing.B) {
		vm := NewVM()
		// Pre-fill stack
		for i := 0; i < 100; i++ {
			vm.Push(i)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			vm.Pop()
		}
	})

	// Benchmark mixed operations
	b.Run("Mixed", func(b *testing.B) {
		vm := NewVM()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			vm.Push(i % 10) // Very small limit
			if i%2 == 0 {
				vm.Pop()
			}
		}
	})
}

// BenchmarkStackSize benchmarks stack size checking
func BenchmarkStackSize(b *testing.B) {
	vm := NewVM()

	// Pre-fill stack
	for i := 0; i < 100; i++ {
		vm.Push(i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = vm.StackSize()
	}
}

// BenchmarkStackRotate benchmarks the rotate operation
func BenchmarkStackRotate(b *testing.B) {
	vm := NewVM()

	// Pre-fill stack with 3 elements
	vm.Push(1)
	vm.Push(2)
	vm.Push(3)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		vm.stack.Rotate(3)
	}
}

// TestStackOptimization tests the stack optimization features
func TestStackOptimization(t *testing.T) {
	vm := NewVM()

	// Test initial capacity
	if vm.stack.GetCapacity() != 64 {
		t.Errorf("Expected initial capacity 64, got %d", vm.stack.GetCapacity())
	}

	// Test max size
	if vm.stack.GetMaxSize() != 10000 {
		t.Errorf("Expected max size 10000, got %d", vm.stack.GetMaxSize())
	}

	// Test stack operations
	vm.Push(1)
	vm.Push(2)
	vm.Push(3)

	if vm.StackSize() != 3 {
		t.Errorf("Expected stack size 3, got %d", vm.StackSize())
	}

	// Test peek
	value, err := vm.stack.Peek()
	if err != nil {
		t.Errorf("Peek failed: %v", err)
	}
	if value != 3 {
		t.Errorf("Expected peek value 3, got %v", value)
	}

	// Test pop
	value = vm.Pop()
	if value != 3 {
		t.Errorf("Expected pop value 3, got %v", value)
	}

	// Test rotate
	err = vm.stack.Rotate(2)
	if err != nil {
		t.Errorf("Rotate failed: %v", err)
	}

	// After rotate [1, 2] -> [2, 1]
	value = vm.Pop()
	if value != 1 {
		t.Errorf("Expected rotated value 1, got %v", value)
	}

	value = vm.Pop()
	if value != 2 {
		t.Errorf("Expected rotated value 2, got %v", value)
	}

	// Test stack underflow
	value = vm.Pop()
	if value != nil {
		t.Errorf("Expected nil for stack underflow, got %v", value)
	}

	// Test stack overflow protection
	for i := 0; i < 10001; i++ {
		err := vm.stack.Push(i)
		if err != nil {
			// Should get stack overflow error
			if i < 10000 {
				t.Errorf("Unexpected error at %d: %v", i, err)
			}
			break
		}
	}
}

// TestStackMemoryEfficiency tests memory efficiency
func TestStackMemoryEfficiency(t *testing.T) {
	vm := NewVM()

	// Test that stack grows efficiently
	initialCapacity := vm.stack.GetCapacity()

	// Push elements to trigger growth
	for i := 0; i < 100; i++ {
		vm.Push(i)
	}

	// Capacity should have grown
	if vm.stack.GetCapacity() <= initialCapacity {
		t.Errorf("Expected capacity to grow from %d, got %d", initialCapacity, vm.stack.GetCapacity())
	}

	// Clear and test that memory is freed
	vm.stack.Clear()
	if vm.StackSize() != 0 {
		t.Errorf("Expected empty stack after clear, got size %d", vm.StackSize())
	}
}

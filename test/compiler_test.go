package test

import (
	"testing"

	"github.com/lengzhao/goscript"
)

func TestCompilerBasic(t *testing.T) {
	// Test basic arithmetic
	script := goscript.NewScript([]byte(`
	package main

	func main() {
		return 1 + 2
	}
	`))

	result, err := script.Run()
	if err != nil {
		t.Fatalf("Failed to run script: %v", err)
	}

	if result != 3 {
		t.Errorf("Expected 3, got %v", result)
	}
}

func TestCompilerVariables(t *testing.T) {
	// Test variable assignment
	script := goscript.NewScript([]byte(`
	package main

	func main() {
		x := 10
		y := 20
		return x + y
	}
	`))

	result, err := script.Run()
	if err != nil {
		t.Fatalf("Failed to run script: %v", err)
	}

	if result != 30 {
		t.Errorf("Expected 30, got %v", result)
	}
}

func TestCompilerFunctionCall(t *testing.T) {
	// Test function call
	script := goscript.NewScript([]byte(`
	package main

	func main() {
		return add(5, 3)
	}
	`))

	// Add the add function
	script.AddFunction("add", func(args ...interface{}) (interface{}, error) {
		a := args[0].(int)
		b := args[1].(int)
		return a + b, nil
	})

	result, err := script.Run()
	if err != nil {
		t.Fatalf("Failed to run script: %v", err)
	}

	if result != 8 {
		t.Errorf("Expected 8, got %v", result)
	}
}

func TestCompilerNestedLoop(t *testing.T) {
	// Test nested loops
	script := goscript.NewScript([]byte(`
	package main

	func multiplyTables() int {
		result := 0
		for i := 1; i <= 2; i++ {
			for j := 1; j <= 2; j++ {
				result += i * j
			}
		}
		return result
	}

	func main() {
		result := multiplyTables()  // (1*1 + 1*2) + (2*1 + 2*2) = 1 + 2 + 2 + 4 = 9
		return result
	}
	`))

	result, err := script.Run()
	if err != nil {
		t.Fatalf("Failed to run script: %v", err)
	}

	if result != 9 {
		t.Errorf("Expected 9, got %v", result)
	}
}

func TestCompilerRecursiveFunction(t *testing.T) {
	// Test recursive function
	script := goscript.NewScript([]byte(`
	package main

	func factorial(n int) int {
		if n <= 1 {
			return 1
		}
		return n * factorial(n-1)
	}

	func main() {
		result := factorial(4)  // 4! = 24
		return result
	}
	`))

	result, err := script.Run()
	if err != nil {
		t.Fatalf("Failed to run script: %v", err)
	}

	if result != 24 {
		t.Errorf("Expected 24, got %v", result)
	}
}

func TestCompilerWhileLoop(t *testing.T) {
	// Test while loop (for loop with only condition)
	script := goscript.NewScript([]byte(`
	package main

	func sumDown(n int) int {
		result := 0
		for n > 0 {
			result += n
			n--
		}
		return result
	}

	func main() {
		result := sumDown(3)  // 3 + 2 + 1 = 6
		return result
	}
	`))

	result, err := script.Run()
	if err != nil {
		t.Fatalf("Failed to run script: %v", err)
	}

	if result != 6 {
		t.Errorf("Expected 6, got %v", result)
	}
}

func TestCompilerComplexCondition(t *testing.T) {
	// Test complex conditions
	script := goscript.NewScript([]byte(`
	package main

	func checkCondition(a, b int) int {
		if a > 0 && b > 0 {
			return a + b
		} else if a < 0 || b < 0 {
			return a * b
		} else {
			return 0
		}
	}

	func main() {
		result := checkCondition(3, 4)  // 3 > 0 && 4 > 0, so return 3 + 4 = 7
		return result
	}
	`))

	result, err := script.Run()
	if err != nil {
		t.Fatalf("Failed to run script: %v", err)
	}

	if result != 7 {
		t.Errorf("Expected 7, got %v", result)
	}
}

func TestCompilerCompoundAssignment(t *testing.T) {
	// Test compound assignment operators
	script := goscript.NewScript([]byte(`
	package main

	func compoundOps() int {
		x := 10
		x += 5   // x = 15
		x -= 3   // x = 12
		return x
	}

	func main() {
		result := compoundOps()  // Should return 12
		return result
	}
	`))

	result, err := script.Run()
	if err != nil {
		t.Fatalf("Failed to run script: %v", err)
	}

	if result != 12 {
		t.Errorf("Expected 12, got %v", result)
	}
}

func TestCompilerDecrementOperator(t *testing.T) {
	// Test decrement operator
	script := goscript.NewScript([]byte(`
	package main

	func decrementTest() int {
		x := 5
		x--
		return x
	}

	func main() {
		result := decrementTest()  // Should return 4
		return result
	}
	`))

	result, err := script.Run()
	if err != nil {
		t.Fatalf("Failed to run script: %v", err)
	}

	if result != 4 {
		t.Errorf("Expected 4, got %v", result)
	}
}

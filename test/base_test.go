package test

import (
	"testing"

	goscript "github.com/lengzhao/goscript"
)

func TestCallFunction(t *testing.T) {
	// Test basic arithmetic
	script := goscript.NewScript([]byte(`
	package test

	func add(a, b int) int {
		return a + b
	}

	func main() {
		return add(1,2)
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

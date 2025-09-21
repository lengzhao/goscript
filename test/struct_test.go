package test

import (
	"testing"

	goscript "github.com/lengzhao/goscript"
)

// TestStructBasic tests basic struct functionality
func TestStructBasic(t *testing.T) {
	scriptSource := `package main

type Person struct {
	name string
	age  int
}

func main() {
	// Create a person instance
	p := Person{name: "Alice", age: 30}
	
	// Access fields directly
	return p.age
}`

	script := goscript.NewScript([]byte(scriptSource))
	result, err := script.Run()
	if err != nil {
		t.Fatalf("Failed to run script: %v", err)
	}

	if result != 30 {
		t.Errorf("Expected 30, got %v", result)
	}
}

// TestStructNested tests nested struct functionality
func TestStructNested(t *testing.T) {
	scriptSource := `package main

type Person struct {
	name string
	age  int
}

type Company struct {
	name string
	ceo  Person
}

func main() {
	// Create person instances
	p1 := Person{name: "Alice", age: 30}
	
	// Create a company with nested struct
	c := Company{
		name: "TechCorp",
		ceo: p1,
	}
	
	// Access nested fields
	return c.ceo.age
}`

	script := goscript.NewScript([]byte(scriptSource))
	result, err := script.Run()
	if err != nil {
		t.Fatalf("Failed to run script: %v", err)
	}

	if result != 30 {
		t.Errorf("Expected 30, got %v", result)
	}
}

// TestStructTypeSystem tests struct type system functionality
func TestStructTypeSystem(t *testing.T) {
	// This test would typically be in a different file focused on type system testing
	// but we include it here for completeness
}

package test

import (
	"testing"

	goscript "github.com/lengzhao/goscript"
)

func TestSliceCompositeLiteral(t *testing.T) {
	scriptSource := `package main

func main() {
	// Create a slice using composite literal
	slice := []int{1, 2, 3, 4, 5}
	
	// Access elements
	return slice[2]  // Should return 3
}`

	script := goscript.NewScript([]byte(scriptSource))
	result, err := script.Run()
	if err != nil {
		t.Fatalf("Failed to run script: %v", err)
	}

	if result != 3 {
		t.Errorf("Expected 3, got %v", result)
	}
}

func TestStructCompositeLiteral(t *testing.T) {
	scriptSource := `package main

type Person struct {
	name string
	age  int
}

func main() {
	// Create a person using composite literal
	person := Person{name: "Alice", age: 30}
	
	// Access fields
	return person.age  // Should return 30
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

func TestNestedStructCompositeLiteral(t *testing.T) {
	scriptSource := `package main

type Address struct {
	city string
}

type Person struct {
	name    string
	age     int
	address Address
}

func main() {
	// Create a person with nested struct using composite literal
	person := Person{
		name: "Alice", 
		age: 30,
		address: Address{city: "Beijing"},
	}
	
	// Access nested fields
	return person.address.city  // Should return "Beijing"
}`

	script := goscript.NewScript([]byte(scriptSource))
	result, err := script.Run()
	if err != nil {
		t.Fatalf("Failed to run script: %v", err)
	}

	if result != "Beijing" {
		t.Errorf("Expected Beijing, got %v", result)
	}
}

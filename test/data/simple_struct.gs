// Simple struct test script
package main

type Person struct {
	name string
	age  int
}

func main() {
	// Create a person instance
	p := Person{name: "Alice", age: 30}
	
	// Access fields directly
	return p.age
}
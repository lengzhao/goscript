// Simple struct test script
package main

type Person struct {
	name string
	age  int
}

func main() {
	// Create a person instance
	p := Person{name: "Alice", age: 30}
	p2 := &Person{name: "Blob", age: 41}
	
	// Access fields directly
	return p.age+p2.age
}
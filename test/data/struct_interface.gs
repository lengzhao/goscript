// Struct and interface test script
package main

// Define a struct type
type Person struct {
	name string
	age  int
}

// Define an interface type
type Speaker interface {
	Speak() string
}

// Implement the interface for Person
func (p Person) Speak() string {
	return "Hello, my name is " + p.name
}

func main() {
	// Create a person instance
	p := Person{name: "Alice", age: 30}
	
	// Call the Speak method
	message := p.Speak()
	
	return message
}
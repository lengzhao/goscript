// Package main provides an example of using structs in GoScript
package main

import (
	"fmt"

	"github.com/lengzhao/goscript"
)

func main() {
	// Create a script that uses structs
	source := `
package main

type Person struct {
	name string
	age  int
}

func (p Person) GetName() string {
	return p.name
}

func (p *Person) SetAge(age int) {
	p.age = age
}

func (p Person) GetAge() int {
	return p.age
}

func main() {
	// Create a new person
	person := Person{name: "Alice", age: 30}
	
	// Use methods
	name := person.GetName()
	person.SetAge(31)
	age := person.GetAge()
	
	// Return a simple integer result instead of string concatenation
	return age
}
`

	script := goscript.NewScript([]byte(source))

	// Execute the script
	result, err := script.Run()
	if err != nil {
		fmt.Printf("Execution error: %v\n", err)
		return
	}

	fmt.Printf("Script result: %v\n", result)
}

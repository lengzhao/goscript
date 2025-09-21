// Complete struct and interface test script
package main

type Person struct {
	name string
	age  int
}

type Speaker interface {
	Speak() string
}

type Employee struct {
	Person
	company string
	salary  float64
}

func (e Employee) Speak() string {
	return "Hello, my name is " + e.name + " and I work at " + e.company
}

func (e Employee) GetSalary() float64 {
	return e.salary
}

func main() {
	// Create an employee
	emp := Employee{
		Person: Person{
			name: "Alice",
			age:  30,
		},
		company: "Acme Corp",
		salary:  50000.0,
	}
	
	// Test field access
	name := emp.name
	age := emp.age
	company := emp.company
	
	// Test method calls
	speech := emp.Speak()
	salary := emp.GetSalary()
	
	// Return a computed result
	return len(name) + age + len(company) + int(salary/1000)
}
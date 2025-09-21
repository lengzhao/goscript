package compiler

import (
	"testing"

	baseCtx "context"

	"github.com/lengzhao/goscript/context"
	"github.com/lengzhao/goscript/parser"
	"github.com/lengzhao/goscript/vm"
)

func TestFunctionRegistration(t *testing.T) {
	// Create a simple script with functions
	script := `
package main

func add(a, b int) int {
    return a + b
}

func main() {
    result := add(5, 3)
    return result
}
`

	// Parse the script
	parser := parser.New()
	astFile, err := parser.Parse("test.go", []byte(script), 0)
	if err != nil {
		t.Fatalf("Failed to parse script: %v", err)
	}

	// Create VM and context
	vmInstance := vm.NewVM()
	context := context.NewExecutionContext()

	// Create compiler and compile
	compiler := NewCompiler(vmInstance, context)
	err = compiler.Compile(astFile)
	if err != nil {
		t.Fatalf("Failed to compile script: %v", err)
	}

	// Check that we have instructions
	instructions := vmInstance.GetInstructions()
	if len(instructions) == 0 {
		t.Fatal("No instructions generated")
	}

	// Check that the first instruction is OpRegistFunction
	if instructions[0].Op != vm.OpRegistFunction {
		t.Errorf("Expected first instruction to be OpRegistFunction, got %s", instructions[0].Op.String())
	}

	// Execute the VM to register functions and run the script
	result, err := vmInstance.Execute(baseCtx.Background())
	if err != nil {
		t.Fatalf("Failed to execute VM: %v", err)
	}

	// Check that function was registered during execution
	scriptFunc, exists := vmInstance.GetScriptFunction("add")
	if !exists {
		t.Fatal("Function 'add' was not registered")
	}

	// Check function properties
	if scriptFunc.Name != "add" {
		t.Errorf("Expected function name 'add', got '%s'", scriptFunc.Name)
	}

	if scriptFunc.ParamCount != 2 {
		t.Errorf("Expected 2 parameters, got %d", scriptFunc.ParamCount)
	}

	expectedParams := []string{"a", "b"}
	for i, param := range expectedParams {
		if scriptFunc.ParamNames[i] != param {
			t.Errorf("Expected parameter %d to be '%s', got '%s'", i, param, scriptFunc.ParamNames[i])
		}
	}

	// Check result
	if result != 8 {
		t.Errorf("Expected result 8, got %v", result)
	}
}

func TestMultipleFunctionRegistration(t *testing.T) {
	// Create a script with multiple functions
	script := `
package main

func add(a, b int) int {
    return a + b
}

func multiply(a, b int) int {
    return a * b
}

func main() {
    result1 := add(3, 4)
    result2 := multiply(5, 6)
    return result1 + result2
}
`

	// Parse the script
	parser := parser.New()
	astFile, err := parser.Parse("test.go", []byte(script), 0)
	if err != nil {
		t.Fatalf("Failed to parse script: %v", err)
	}

	// Create VM and context
	vmInstance := vm.NewVM()
	context := context.NewExecutionContext()

	// Create compiler and compile
	compiler := NewCompiler(vmInstance, context)
	err = compiler.Compile(astFile)
	if err != nil {
		t.Fatalf("Failed to compile script: %v", err)
	}

	// Execute the VM to register functions and run the script
	result, err := vmInstance.Execute(baseCtx.Background())
	if err != nil {
		t.Fatalf("Failed to execute VM: %v", err)
	}

	// Check that functions were registered during execution
	functions := []string{"add", "multiply"}
	for _, funcName := range functions {
		_, exists := vmInstance.GetScriptFunction(funcName)
		if !exists {
			t.Errorf("Function '%s' was not registered", funcName)
		}
	}

	// Check result: add(3,4)=7, multiply(5,6)=30, total=37
	if result != 37 {
		t.Errorf("Expected result 37, got %v", result)
	}
}

// TestStructCompilation tests compilation of struct types and struct literals
func TestStructCompilation(t *testing.T) {
	// Create a script with struct type definition and struct literal
	script := `
package main

type Person struct {
	name string
	age  int
}

func main() {
	// Create a new person using struct literal
	person := Person{name: "Alice", age: 30}
	
	// Access fields
	return person.age
}
`

	// Parse the script
	parser := parser.New()
	astFile, err := parser.Parse("test.go", []byte(script), 0)
	if err != nil {
		t.Fatalf("Failed to parse script: %v", err)
	}

	// Create VM and context
	vmInstance := vm.NewVM()
	context := context.NewExecutionContext()

	// Create compiler and compile
	compiler := NewCompiler(vmInstance, context)
	err = compiler.Compile(astFile)
	if err != nil {
		t.Fatalf("Failed to compile script: %v", err)
	}

	// Execute the VM
	result, err := vmInstance.Execute(baseCtx.Background())
	if err != nil {
		t.Fatalf("Failed to execute VM: %v", err)
	}

	// Check result: person.age should be 30
	if result != 30 {
		t.Errorf("Expected result 30, got %v", result)
	}
}

// TestSliceCompilation tests compilation of slice operations
func TestSliceCompilation(t *testing.T) {
	// Create a script with slice operations
	script := `
package main

func main() {
	// Create a slice
	numbers := []int{1, 2, 3, 4, 5}
	
	// Access element by index
	value := numbers[2]
	
	// Modify element by index
	numbers[1] = 10
	
	// Return modified value
	return numbers[1]
}
`

	// Parse the script
	parser := parser.New()
	astFile, err := parser.Parse("test.go", []byte(script), 0)
	if err != nil {
		t.Fatalf("Failed to parse script: %v", err)
	}

	// Create VM and context
	vmInstance := vm.NewVM()
	context := context.NewExecutionContext()

	// Create compiler and compile
	compiler := NewCompiler(vmInstance, context)
	err = compiler.Compile(astFile)
	if err != nil {
		t.Fatalf("Failed to compile script: %v", err)
	}

	// Execute the VM
	result, err := vmInstance.Execute(baseCtx.Background())
	if err != nil {
		t.Fatalf("Failed to execute VM: %v", err)
	}

	// Check result: numbers[1] should be 10
	if result != 10 {
		t.Errorf("Expected result 10, got %v", result)
	}
}

// TestStructMethodCompilation tests compilation of struct methods
func TestStructMethodCompilation(t *testing.T) {
	// Create a script with struct methods
	script := `
package main

type Person struct {
	name string
	age  int
}

// Value receiver method
func (p Person) GetName() string {
	return p.name
}

// Pointer receiver method
func (p *Person) SetAge(age int) {
	p.age = age
}

// Value receiver method
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
	
	// Return age
	return age
}
`

	// Parse the script
	parser := parser.New()
	astFile, err := parser.Parse("test.go", []byte(script), 0)
	if err != nil {
		t.Fatalf("Failed to parse script: %v", err)
	}

	// Create VM and context
	vmInstance := vm.NewVM()
	context := context.NewExecutionContext()

	// Create compiler and compile
	compiler := NewCompiler(vmInstance, context)
	err = compiler.Compile(astFile)
	if err != nil {
		t.Fatalf("Failed to compile script: %v", err)
	}

	// Execute the VM
	result, err := vmInstance.Execute(baseCtx.Background())
	if err != nil {
		t.Fatalf("Failed to execute VM: %v", err)
	}

	// Check that methods were registered during compilation
	methods := []string{"Person.GetName", "Person.SetAge", "Person.GetAge"}
	for _, methodName := range methods {
		_, exists := vmInstance.GetScriptFunction(methodName)
		if !exists {
			t.Errorf("Method '%s' was not registered", methodName)
		}
	}

	// Check result: age should be 31 after SetAge(31)
	if result != 31 {
		t.Errorf("Expected result 31, got %v", result)
	}
}

// TestComplexStructAndSliceOperations tests more complex struct and slice operations
func TestComplexStructAndSliceOperations(t *testing.T) {
	// Create a script with complex struct and slice operations
	script := `
package main

type Person struct {
	name string
	age  int
}

func main() {
	// Create a slice of structs
	people := []Person{
		Person{name: "Alice", age: 30},
		Person{name: "Bob", age: 25},
		Person{name: "Charlie", age: 35},
	}
	
	// Access struct in slice and get field
	aliceAge := people[0].age
	
	// Modify struct in slice
	people[1].age = 26
	
	// Access modified value
	bobAge := people[1].age
	
	// Return sum
	return aliceAge + bobAge
}
`

	// Parse the script
	parser := parser.New()
	astFile, err := parser.Parse("test.go", []byte(script), 0)
	if err != nil {
		t.Fatalf("Failed to parse script: %v", err)
	}

	// Create VM and context
	vmInstance := vm.NewVM()
	context := context.NewExecutionContext()

	// Create compiler and compile
	compiler := NewCompiler(vmInstance, context)
	err = compiler.Compile(astFile)
	if err != nil {
		t.Fatalf("Failed to compile script: %v", err)
	}

	// Execute the VM
	result, err := vmInstance.Execute(baseCtx.Background())
	if err != nil {
		t.Fatalf("Failed to execute VM: %v", err)
	}

	// Check result: aliceAge (30) + bobAge (26) = 56
	if result != 56 {
		t.Errorf("Expected result 56, got %v", result)
	}
}

// TestNestedStructCompilation tests compilation of nested structs
func TestNestedStructCompilation(t *testing.T) {
	// Create a script with nested struct types
	script := `
package main

type Address struct {
	street string
	city   string
}

type Person struct {
	name    string
	age     int
	address Address
}

func main() {
	// Create a new person with nested address
	person := Person{
		name: "Alice",
		age: 30,
		address: Address{
			street: "123 Main St",
			city: "New York",
		},
	}
	
	// Access nested field
	return person.address.city
}
`

	// Parse the script
	parser := parser.New()
	astFile, err := parser.Parse("test.go", []byte(script), 0)
	if err != nil {
		t.Fatalf("Failed to parse script: %v", err)
	}

	// Create VM and context
	vmInstance := vm.NewVM()
	context := context.NewExecutionContext()

	// Create compiler and compile
	compiler := NewCompiler(vmInstance, context)
	err = compiler.Compile(astFile)
	if err != nil {
		t.Fatalf("Failed to compile script: %v", err)
	}

	// Execute the VM
	result, err := vmInstance.Execute(baseCtx.Background())
	if err != nil {
		t.Fatalf("Failed to execute VM: %v", err)
	}

	// Check result: person.address.city should be "New York"
	if result != "New York" {
		t.Errorf("Expected result 'New York', got %v", result)
	}
}

// TestSliceOfStructsCompilation tests compilation of slices of structs
func TestSliceOfStructsCompilation(t *testing.T) {
	// Create a script with slice of structs
	script := `
package main

type Person struct {
	name string
	age  int
}

func main() {
	// Create a slice of structs
	people := []Person{
		Person{name: "Alice", age: 30},
		Person{name: "Bob", age: 25},
		Person{name: "Charlie", age: 35},
	}
	
	// Access element in slice and get field
	return people[1].age
}
`

	// Parse the script
	parser := parser.New()
	astFile, err := parser.Parse("test.go", []byte(script), 0)
	if err != nil {
		t.Fatalf("Failed to parse script: %v", err)
	}

	// Create VM and context
	vmInstance := vm.NewVM()
	context := context.NewExecutionContext()

	// Create compiler and compile
	compiler := NewCompiler(vmInstance, context)
	err = compiler.Compile(astFile)
	if err != nil {
		t.Fatalf("Failed to compile script: %v", err)
	}

	// Execute the VM
	result, err := vmInstance.Execute(baseCtx.Background())
	if err != nil {
		t.Fatalf("Failed to execute VM: %v", err)
	}

	// Check result: people[1].age should be 25
	if result != 25 {
		t.Errorf("Expected result 25, got %v", result)
	}
}

// TestPointerReceiverMethod tests pointer receiver methods
func TestPointerReceiverMethod(t *testing.T) {
	// Create a script with pointer receiver methods
	script := `
package main

type Counter struct {
	value int
}

// Pointer receiver method
func (c *Counter) Increment() {
	c.value = c.value + 1
}

// Value receiver method
func (c Counter) GetValue() int {
	return c.value
}

func main() {
	// Create a new counter
	counter := Counter{value: 10}
	
	// Use methods
	counter.Increment()
	value := counter.GetValue()
	
	// Return value
	return value
}
`

	// Parse the script
	parser := parser.New()
	astFile, err := parser.Parse("test.go", []byte(script), 0)
	if err != nil {
		t.Fatalf("Failed to parse script: %v", err)
	}

	// Create VM and context
	vmInstance := vm.NewVM()
	context := context.NewExecutionContext()

	// Create compiler and compile
	compiler := NewCompiler(vmInstance, context)
	err = compiler.Compile(astFile)
	if err != nil {
		t.Fatalf("Failed to compile script: %v", err)
	}

	// Execute the VM
	result, err := vmInstance.Execute(baseCtx.Background())
	if err != nil {
		t.Fatalf("Failed to execute VM: %v", err)
	}

	// Check result: value should be 11 after increment
	if result != 11 {
		t.Errorf("Expected result 11, got %v", result)
	}
}

// TestMethodChaining tests method chaining
func TestMethodChaining(t *testing.T) {
	// Create a script with method chaining
	script := `
package main

type Calculator struct {
	value int
}

// Pointer receiver method that returns the receiver
func (c *Calculator) Add(x int) *Calculator {
	c.value = c.value + x
	return c
}

// Pointer receiver method that returns the receiver
func (c *Calculator) Multiply(x int) *Calculator {
	c.value = c.value * x
	return c
}

// Value receiver method to get the value
func (c Calculator) GetValue() int {
	return c.value
}

func main() {
	// Create a new calculator and chain methods
	calc := Calculator{value: 5}
	result := calc.Add(3).Multiply(2).GetValue()
	
	// Return result: (5 + 3) * 2 = 16
	return result
}
`

	// Parse the script
	parser := parser.New()
	astFile, err := parser.Parse("test.go", []byte(script), 0)
	if err != nil {
		t.Fatalf("Failed to parse script: %v", err)
	}

	// Create VM and context
	vmInstance := vm.NewVM()
	context := context.NewExecutionContext()

	// Create compiler and compile
	compiler := NewCompiler(vmInstance, context)
	err = compiler.Compile(astFile)
	if err != nil {
		t.Fatalf("Failed to compile script: %v", err)
	}

	// Execute the VM
	result, err := vmInstance.Execute(baseCtx.Background())
	if err != nil {
		t.Fatalf("Failed to execute VM: %v", err)
	}

	// Check result: (5 + 3) * 2 = 16
	if result != 16 {
		t.Errorf("Expected result 16, got %v", result)
	}
}

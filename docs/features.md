# GoScript Feature Details

## 1. Overview

GoScript is a script engine compatible with Go standard syntax that supports most of Go language's core features. This document details the features and capabilities currently supported by GoScript.

## 2. Supported Syntax Features

### 2.1 Basic Syntax Structure

#### Package Declaration
```go
package main
```

#### Import Declaration
```go
import "fmt"
import (
    "strings"
    "math"
)
```

#### Variable Declaration
```go
// var declaration
var x int
var y, z int = 10, 20
var (
    a = 1
    b = 2
)

// Short variable declaration
x := 10
name := "GoScript"
```

#### Constant Declaration
```go
const pi = 3.14159
const (
    a = 1
    b = 2
)
```

### 2.2 Data Types

#### Basic Types
- Integer types: int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64
- Floating-point types: float32, float64
- Boolean type: bool
- String type: string

#### Composite Types
- Array: [n]T
- Slice: []T
- Struct: struct
- Interface: interface{}

### 2.3 Control Structures

#### Conditional Statements
```go
if x > 0 {
    // do something
} else if x < 0 {
    // do something else
} else {
    // do another thing
}
```

#### Loop Statements
```go
// Traditional for loop
for i := 0; i < 10; i++ {
    // do something
}

// While loop form
for x > 0 {
    x--
}

// Infinite loop
for {
    // do something
    break
}
```

#### Range Statements
```go
// Iterate over slice
slice := []int{1, 2, 3}
for index, value := range slice {
    // index is the index, value is the value
}

// Iterate values only
for _, value := range slice {
    // only care about values
}

// Iterate indices only
for index := range slice {
    // only care about indices
}

// Iterate over strings
str := "hello"
for index, char := range str {
    // char is of rune type
}
```

### 2.4 Functions

#### Function Declaration
```go
func add(a, b int) int {
    return a + b
}

func greet(name string) string {
    return "Hello, " + name
}

// Multiple return values
func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}
```

#### Function Calls
```go
result := add(1, 2)
greeting := greet("World")
```

### 2.5 Structs and Methods

#### Struct Definition
```go
type Person struct {
    Name string
    Age  int
}

type Rectangle struct {
    Width, Height float64
}
```

#### Struct Initialization
```go
// Field name initialization
person := Person{Name: "Alice", Age: 30}

// Order initialization
person := Person{"Bob", 25}

// Composite literals
rect := Rectangle{Width: 10, Height: 5}
```

#### Method Definition
```go
// Value receiver method
func (p Person) GetName() string {
    return p.Name
}

// Pointer receiver method
func (p *Person) SetAge(age int) {
    p.Age = age
}

// Method with return value
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}
```

#### Method Calls
```go
person := Person{Name: "Alice", Age: 30}
name := person.GetName()
person.SetAge(31)
```

### 2.6 Operators

#### Arithmetic Operators
- Addition: +
- Subtraction: -
- Multiplication: *
- Division: /
- Modulus: %
- Increment: ++
- Decrement: --

#### Comparison Operators
- Equal: ==
- Not equal: !=
- Less than: <
- Less than or equal: <=
- Greater than: >
- Greater than or equal: >=

#### Logical Operators
- Logical AND: &&
- Logical OR: ||
- Logical NOT: !

#### Assignment Operators
- Simple assignment: =
- Add assignment: +=
- Subtract assignment: -=
- Multiply assignment: *=
- Divide assignment: /=
- Modulus assignment: %=

## 3. Built-in Functions

### 3.1 Basic Built-in Functions
- len(): Get the length of strings, arrays, slices, and maps
- int(): Convert value to integer
- float64(): Convert value to floating-point number
- string(): Convert value to string

## 4. Module System

### 4.1 Built-in Modules
GoScript provides the following built-in modules:
- math: Mathematical functions
- strings: String operations
- fmt: Formatted input/output
- json: JSON serialization and deserialization

### 4.2 Module Usage
```go
// Using module functions
result := math.Abs(-5.0)
upper := strings.ToUpper("hello")
```

## 5. Error Handling

### 5.1 Error Return
```go
func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}
```

### 5.2 Error Checking
```go
result, err := divide(10, 0)
if err != nil {
    // Handle error
    return
}
```

## 6. Limitations and Unsupported Features

### 6.1 Unsupported Syntax Features
- Goroutines and channels
- unsafe package
- Reflection (reflect package)
- Complete package management system
- Type assertions
- Concrete implementation of interfaces
- defer statements
- switch statements
- select statements

### 6.2 Type System Limitations
- No support for generics
- No support for complex usage of type aliases
- No support for struct tags

## 7. Performance Features

### 7.1 Compile and Execute
GoScript compiles source code into bytecode, then executes it in a virtual machine, which provides better performance compared to pure interpretation.

### 7.2 Scope Optimization
The key-based context management mechanism provides efficient scope lookup and variable management.

### 7.3 Memory Management
Object pooling and pre-allocation mechanisms reduce memory allocation and GC pressure.

## 8. Security Features

### 8.1 Resource Limitations
- Maximum execution time limit
- Maximum memory usage limit
- Maximum instruction count limit

### 8.2 Sandbox Environment
- Prohibition of dangerous system calls
- Restriction of file system access
- Restriction of network access

### 8.3 Module Access Control
- Configurable module access permissions
- Prohibited keyword list
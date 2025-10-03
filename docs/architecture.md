# GoScript Comprehensive Technical Documentation

## 1. Project Overview

GoScript is a script engine compatible with Go standard syntax that allows you to dynamically execute Go code within Go applications. This project aims to implement a safe, efficient, and extensible script execution environment.

### 1.1 Design Goals

1. **Syntax Compatibility**: As compatible as possible with Go standard syntax
2. **Security**: Provide a sandbox environment to limit dangerous operations
3. **Extensibility**: Support custom functions, types, and modules
4. **Performance**: Improve execution efficiency through compilation
5. **Usability**: Provide a simple API for Go applications to call

### 1.2 Core Features

- **Modular Design**: Separation of lexical analysis, syntax analysis, AST generation, and other components
- **Extensibility**: Support for custom functions and modules
- **Security**: Provides execution time and memory usage limits
- **Reuse of Go Native Modules**: Lexical analysis, syntax analysis, etc. directly reuse Go standard library
- **Key-based Context Management**: Manage scope and variables through unique identifiers
- **Complete Go Syntax Support**: Support for structs, methods, range statements, composite literals, etc.

## 2. Architecture Design

### 2.1 Overall Architecture

```
+-------------------+
|   Go Application  |
+-------------------+
          |
          v
+-------------------+
|  Script Interface | <- API layer, provided for Go application calls
+-------------------+
          |
          v
+-------------------+
|   Parser/Lexer    | <- Lexical and syntax analyzer (reusing Go standard library)
+-------------------+
          |
          v
+-------------------+
|   AST Generator   | <- Abstract syntax tree generator (reusing Go standard library)
+-------------------+
          |
          v
+-------------------+
|   Compiler        | <- Compiler, compiles AST to bytecode
+-------------------+
          |
          v
+-------------------+
|   Bytecode VM     | <- Bytecode virtual machine, executes compiled code
+-------------------+
          |
          v
+-------------------+
|  Runtime System   | <- Runtime system, manages objects, memory, etc.
+-------------------+
```

### 2.2 Core Components

GoScript's core components include:

1. **Script (script.go)**: Main API interface, providing methods like NewScript, Run, AddFunction, etc.
2. **Parser (parser/)**: Lexical and syntax analyzer, reusing Go standard library's `go/scanner` and `go/parser`
3. **Compiler (compiler/)**: Compiler that compiles AST into executable intermediate representation (bytecode)
4. **VM (vm/)**: Virtual machine that executes compiled bytecode
5. **Context (context/)**: Execution context management, managing variable scope and stack during script execution
6. **Instruction (instruction/)**: Instruction definitions, defining opcodes executable by the virtual machine
7. **Types (types/)**: Type system, defining type interfaces and module executors in GoScript
8. **Builtin (builtin/)**: Built-in functions and modules, providing standard library functionality such as math, strings, etc.

These components work together to provide a complete script execution environment.

### 2.3 Core Design Concepts

1. **Leveraging Go Standard Library**: Fully utilize Go standard library functions, especially the `context` package to manage execution context and variable scope
2. **Simplified Opcode Design**: Reduce virtual machine complexity and improve execution efficiency through simplified opcode design
3. **Modular Architecture**: Adopt modular design with clear component responsibilities for easy maintenance and extension

## 3. Type System

### 3.1 IType Interface

All types support the IType interface, providing a unified type operation interface:

```go
type IType interface {
    // TypeName returns the name of the type
    TypeName() string

    // String returns the string representation of the value
    String() string

    // Equals compares two values of the same type
    Equals(other IType) bool
}
```

### 3.2 Basic Type Implementation

- IntType: Integer type
- FloatType: Floating-point type
- StringType: String type
- BoolType: Boolean type

## 4. Execution Context and Scope Management

### 4.1 Context Structure

The current implementation uses a hierarchical context system based on the `context.Context` package:

```go
type Context struct {
    // Path key for identifying the context (e.g., "main.function.loop")
    pathKey string

    // Parent context reference
    parent *Context

    // Variables in this context
    variables map[string]interface{}

    // Variable types in this context
    types map[string]string

    // Child contexts
    children map[string]*Context
}
```

### 4.2 Scope Nesting

```
Global scope
└── Module scope
    └── Function scope
        └── Block scope
```

### 4.3 Variable Lookup and Isolation

The Context structure implements a natural variable lookup chain, searching from the current scope upward to the global scope. Variables in different scopes are naturally isolated to prevent variable pollution.

### 4.4 Key-based Context Management (New Design)

To better manage scopes, a key-based context management mechanism has been introduced:

1. **Unique Identifier**: Each scope has a unique key identifier
   - Global scope: `main`
   - Main function: `main.main`
   - Regular functions: `main.FunctionName`
   - Struct methods: `main.StructName.MethodName`
   - Other modules: `moduleName.FunctionName`
   - Code blocks: `main.main.block_1`

2. **Compile-time Tracking**: The compiler analyzes the current context's key when compiling BlockStmt and generates corresponding scope management instructions

3. **Runtime Management**: Runtime creates Context objects to manage variables and reference relationships, with each Context referencing its parent context

4. **Variable Lookup**: Variable lookup follows the scope chain, searching from the current context upward to the global context

5. **Scope Isolation**: Variables in different scopes are naturally isolated to prevent variable pollution

## 5. Virtual Machine and Opcodes

### 5.1 Simplified Opcodes

```go
const (
    OpNop        OpCode = iota // No operation
    OpLoadConst               // Load constant
    OpLoadName                // Load variable
    OpStoreName               // Store variable
    OpCall                    // Call function
    OpReturn                  // Return
    OpJump                    // Jump
    OpJumpIf                  // Conditional jump
    OpBinaryOp                // Binary operation
    OpUnaryOp                 // Unary operation
    OpEnterScope              // Enter scope
    OpExitScope               // Exit scope
    OpEnterScopeWithKey       // Enter scope with specified key
    OpExitScopeWithKey        // Exit scope with specified key
    OpCreateVar               // Create variable
    OpNewSlice                // Create slice
    OpNewStruct               // Create struct
    OpGetField                // Get struct field
    OpSetField                // Set struct field
    OpGetIndex                // Get indexed element
    OpSetIndex                // Set indexed element
    OpLen                     // Get length
    OpImport                  // Import module
)
```

### 5.2 Instruction Format

```go
type Instruction struct {
    Op   OpCode      // Opcode
    Arg  interface{} // Argument 1
    Arg2 interface{} // Argument 2
}
```

## 6. Function Registry Mechanism

### 6.1 Unified Function Calls

All functions (built-in functions, user-defined functions, module functions) are registered and called through the same mechanism.

### 6.2 Function Registration Process

1. Create function instance
2. Register function through VM.RegisterFunction
3. Function is stored in VM's function registry

### 6.3 Function Call Process

1. Check if it's a module function call (moduleName.functionName)
2. Look up function in VM's function registry
3. Execute function directly or through module executor

## 7. Module System

### 7.1 Module Structure

Modules are managed through the VM's module registry system:

```go
// Module functions are registered through the ModuleExecutor interface
type ModuleExecutor func(entrypoint string, args ...interface{}) (interface{}, error)
```

### 7.2 Module Management

Supports module definition and registration, function registration, inter-module calls, built-in module support, and module access control.

## 8. Syntax Support

### 8.1 Supported Go Syntax Features

1. **Basic Types**:
   - Boolean type (bool)
   - Numeric types (int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64)
   - String type (string)
   - Composite types ([]T, [n]T, map[K]T, struct, interface)

2. **Variable Declaration**:
   - var declaration
   - Short variable declaration (:=)
   - Constant declaration (const)

3. **Control Structures**:
   - Conditional statements (if/else)
   - Loop statements (for)
   - Range statements (for range)
   - Switch statements
   - Goto statements

4. **Functions**:
   - Function declaration
   - Function calls
   - Method declaration and calls
   - Multiple return values
   - Variadic parameters
   - Anonymous functions
   - Closures

5. **Data Structures**:
   - Arrays and slices
   - Maps
   - Structs
   - Pointers
   - Composite literals ([]int{1, 2, 3} or Person{name: "Alice"})

6. **Operators**:
   - Arithmetic operators (+, -, *, /, %)
   - Comparison operators (==, !=, <, <=, >, >=)
   - Logical operators (&&, ||, !)
   - Assignment operators (=, +=, -=, *=, /=, %=)
   - Increment/decrement operators (++, --)

7. **Error Handling**:
   - error type
   - panic/recover mechanism

### 8.2 Limited Syntax Features

For security and implementation simplification, the following Go features are restricted or not supported:

1. **Package Management**: No support for complete Go package system, functionality provided through module system
2. **Concurrency**: Restriction or no support for goroutines and channels, limited concurrency functionality provided through modules
3. **Low-level Operations**: No support for unsafe package, restriction of pointer operations
4. **Reflection**: Restriction or no support for reflect package
5. **System Calls**: Limited system functionality provided through modules

## 9. Security Mechanisms

### 9.1 Sandbox Environment

1. **Resource Limitations**:
   - Maximum execution time limit
   - Maximum memory usage limit
   - Maximum object allocation limit

2. **API Limitations**:
   - Prohibition of dangerous system calls
   - Restriction of file system access
   - Restriction of network access

3. **Syntax Limitations**:
   - Prohibition of certain keywords
   - Complexity restrictions

### 9.2 Security Context

```go
// Security is managed through VM instruction limits
// Script level security configuration
type Script struct {
    // Maximum number of instructions allowed (0 means no limit)
    maxInstructions int64
}

// Set security context
func (s *Script) SetMaxInstructions(max int64)
```

## 10. API Design

### 10.1 Script Interface

```go
type Script struct {
    // Script content
    source []byte
    
    // Virtual machine
    vm *vm.VM

    // Debug mode
    debug bool

    // Execution statistics
    executionStats *ExecutionStats

    // Maximum number of instructions allowed (0 means no limit)
    maxInstructions int64
}

// Create new script
func NewScript(source []byte) *Script

// Add function
func (s *Script) AddFunction(name string, fn vm.ScriptFunction) error

// Register module
func (s *Script) RegisterModule(moduleName string, executor types.ModuleExecutor)

// Compile and execute script
func (s *Script) Run() (interface{}, error)

// Execute script (with context)
func (s *Script) RunContext(ctx context.Context) (interface{}, error)
```

## 11. Usage Examples

### 11.1 Basic Usage

```go
package main

import (
    "fmt"
    "context"
    "time"
    "github.com/lengzhao/goscript"
)

func main() {
    // Create script
    source := `
package main

func main() {
    x := 10
    y := 20
    return x + y
}
`
    
    script := goscript.NewScript([]byte(source))
    script.SetDebug(true) // Enable debug mode
    
    // Execute script
    result, err := script.Run()
    if err != nil {
        fmt.Printf("Execution error: %v\n", err)
        return
    }
    
    fmt.Printf("Execution result: %v\n", result) // Output: 30
}
```

### 11.2 Custom Extensions

```go
// Custom function
func customFunction(args ...interface{}) (interface{}, error) {
    if len(args) != 2 {
        return nil, fmt.Errorf("requires 2 arguments")
    }
    
    a, ok1 := args[0].(int)
    b, ok2 := args[1].(int)
    if !ok1 || !ok2 {
        return nil, fmt.Errorf("arguments must be integers")
    }
    
    return a * b, nil
}

// Using custom function
func main() {
    source := `
package main

func main() {
    result := customMultiply(5, 6)
    return result
}
`
    
    script := goscript.NewScript([]byte(source))
    
    // Register custom function with the VM
    script.AddFunction("customMultiply", func(args ...interface{}) (interface{}, error) {
        if len(args) != 2 {
            return nil, fmt.Errorf("customMultiply function requires 2 arguments")
        }
        a, ok1 := args[0].(int)
        b, ok2 := args[1].(int)
        if !ok1 || !ok2 {
            return nil, fmt.Errorf("customMultiply function requires integer arguments")
        }
        return a * b, nil
    })
    
    result, err := script.Run()
    if err != nil {
        fmt.Printf("Execution error: %v\n", err)
        return
    }
    
    fmt.Printf("Execution result: %v\n", result) // Output: 30
}
```

## 12. Execution Flow

1. **Lexical Analysis**: Source code → Tokens (reusing Go standard library)
2. **Syntax Analysis**: Tokens → AST (reusing Go standard library)
3. **Compilation**: AST → Bytecode instructions and constant pool
4. **Execution**:
   - Create VM context
   - Load modules and functions
   - Execute bytecode instructions
   - Manage scope and variables

## 13. Performance Optimization

1. **Opcode Optimization**: Reduce virtual machine complexity and improve execution efficiency through simplified opcode design
2. **Scope Management Optimization**: Optimize scope management using hierarchical Context objects to reduce memory allocation and lookup time
3. **Standard Library Reuse**: Reuse Go standard library's lexical analysis, syntax analysis, and AST processing functions to ensure compatibility and performance

## 14. Extension Mechanisms

### 14.1 Custom Functions

Support for registering custom functions through Go code:

```go
// Register custom function
script.AddFunction("myFunc", func(args ...interface{}) (interface{}, error) {
    // Function implementation
    return result, nil
})
```

### 14.2 Custom Types

Support for implementing custom types through the IType interface:

```go
// Implement custom type by implementing the IType interface
type MyType struct {
    // Field definitions
}

func (m *MyType) TypeName() string {
    return "MyType"
}

func (m *MyType) String() string {
    // String representation
}

// Register custom type through custom functions or modules
```

### 14.3 Module System

Support for modular extensions:

```go
// Create and register module with ModuleExecutor
moduleExecutor := func(entrypoint string, args ...interface{}) (interface{}, error) {
    // Module implementation
    return result, nil
}

// Register module
script.RegisterModule("myModule", moduleExecutor)
```

## 15. Testing

Run all tests:

```
go test ./...
```

Run tests for specific packages:

```
go test ./lexer -v
go test ./parser -v
go test ./ast -v
go test ./compiler -v
go test ./runtime -v
go test ./vm -v
```

## 16. Summary

GoScript implements a concise, efficient, and secure script engine through the following approaches:

1. **Leveraging Go Standard Library**: Reusing mature Go standard library functions
2. **Simplified Design**: Reducing complexity through simplified opcode and component design
3. **Natural Scope Management**: Implementing natural scope management using hierarchical Context objects
4. **Modular Architecture**: Clear component responsibility division for easy maintenance and extension
5. **Built-in Security Mechanisms**: Providing instruction count limits for security controls

This design makes GoScript an easy-to-use, high-performance, and secure script engine suitable for various Go application dynamic execution needs.
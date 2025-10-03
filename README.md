# GoScript - Go-Compatible Scripting Engine

GoScript is a scripting engine compatible with Go standard syntax that allows you to dynamically execute Go code within Go applications.

[ZH 中文版](README_cn.md)

## Features

- **Syntax Compatibility**: As compatible as possible with Go standard syntax
- **Modular Design**: Separation of components such as lexical analysis, syntax analysis, and AST generation
- **Extensibility**: Support for custom functions and modules
- **Reuse of Go Native Modules**: Lexical analysis, syntax analysis, etc. directly reuse Go standard library

## Installation

```bash
go get github.com/lengzhao/goscript
```

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
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

### Custom Functions

```go
// Create custom function
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

// Using custom functions
func main() {
    source := `
package main

func main() {
    result := customMultiply(5, 6)
    return result
}
`
    
    script := goscript.NewScript([]byte(source))
    
    // Register the custom function
    err := script.AddFunction("customMultiply", func(args ...interface{}) (interface{}, error) {
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
    if err != nil {
        fmt.Printf("Failed to register custom function: %v\n", err)
        return
    }
    
    result, err := script.Run()
    if err != nil {
        fmt.Printf("Execution error: %v\n", err)
        return
    }
    
    fmt.Printf("Execution result: %v\n", result) // Output: 30
}
```

### Using Built-in Modules

GoScript provides several built-in modules, including `math`, `strings`, `fmt`, and `json`:

```go
func main() {
    // Using strings module
    source := `
package main

import "strings"

func main() {
    lowerStr := strings.ToLower("HELLO WORLD")
    hasWorld := strings.Contains("hello world", "world")
    
    return map[string]interface{}{
        "lower": lowerStr,
        "contains": hasWorld,
    }
}
`
    
    script := goscript.NewScript([]byte(source))
    
    // Register builtin modules
    modules := []string{"strings"}
    for _, moduleName := range modules {
        moduleExecutor, exists := builtin.GetModuleExecutor(moduleName)
        if exists {
            script.RegisterModule(moduleName, moduleExecutor)
        }
    }
    
    result, err := script.Run()
    if err != nil {
        fmt.Printf("Execution error: %v\n", err)
        return
    }
    
    fmt.Printf("Execution result: %v\n", result)
}
```

### Calling Functions Directly

```go
func main() {
    // Create a new script
    script := goscript.NewScript([]byte{})
    
    // Add a function using AddFunction method that uses arguments
    script.AddFunction("addFunc", func(args ...interface{}) (interface{}, error) {
        if len(args) != 2 {
            return nil, fmt.Errorf("addFunc requires 2 arguments")
        }
        a, ok1 := args[0].(int)
        b, ok2 := args[1].(int)
        if !ok1 || !ok2 {
            return nil, fmt.Errorf("addFunc requires integer arguments")
        }
        return a + b, nil
    })
    
    // Call the function directly
    result, err := script.CallFunction("addFunc", 5, 6)
    if err != nil {
        fmt.Printf("Failed to call function: %v\n", err)
        return
    }
    
    fmt.Printf("Function result: %v\n", result) // Output: 11
}
```

## Core Components

### Script

The main interface for the GoScript engine. Key methods include:

- `NewScript(source []byte) *Script` - Creates a new script
- `Run() (interface{}, error)` - Executes the script
- `AddFunction(name string, execFn vm.ScriptFunction) error` - Adds a custom function
- `CallFunction(name string, args ...interface{}) (interface{}, error)` - Calls a function directly
- `SetDebug(debug bool)` - Enables or disables debug mode
- `RegisterModule(moduleName string, executor types.ModuleExecutor)` - Registers a module
- `SetMaxInstructions(max int64)` - Sets the maximum number of instructions (default: 10000)

### Virtual Machine (VM)

The virtual machine is responsible for executing compiled bytecode:

```go
// Create virtual machine
vmInstance := vm.NewVM()

// Register function
vmInstance.RegisterFunction("multiply", func(args ...interface{}) (interface{}, error) {
    // Function implementation
    return result, nil
})

// Call function
result, err := vmInstance.Execute("main.main", arg1, arg2)

// Get function
fn, exists := vmInstance.GetFunction("functionName")
```

### Built-in Modules

GoScript provides several built-in modules:

1. **strings** - String manipulation functions
2. **math** - Mathematical functions
3. **fmt** - Formatting functions
4. **json** - JSON encoding/decoding functions

## Security Features

GoScript provides multiple security mechanisms to prevent script abuse of system resources:

### 1. Instruction Count Limit
Limit the maximum number of instructions a script can execute. The default limit is 10,000 instructions:

```go
script := goscript.NewScript(source)
// Default is 10,000 instructions
// script.SetMaxInstructions(10000)

// Set custom limit
script.SetMaxInstructions(5000) // Limit to 5,000 instructions

// Remove limit
script.SetMaxInstructions(0) // No limit
```

### 2. Module Access Control
Control which modules scripts can access by selectively registering modules.

## Testing

Run all tests:

```bash
go test ./...
```

Run tests for specific packages:

```bash
go test ./test -v
```

## Examples

Check the example programs in the `examples/` directory:

- `examples/basic/` - Basic usage examples
- `examples/custom_function/` - Custom function examples
- `examples/builtin_functions/` - Built-in function examples
- `examples/modules/` - Module usage examples
- `examples/interface_example/` - Interface examples
- `examples/struct_example/` - Struct examples

Run examples:

```bash
cd examples/custom_function
go run function_demo.go
```

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the project
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

## Contact

Project link: [https://github.com/lengzhao/goscript](https://github.com/lengzhao/goscript)
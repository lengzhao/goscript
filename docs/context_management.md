# GoScript Context Management Design Document

## 1. Overview

This document details the design of GoScript's key-based context management mechanism. This mechanism aims to provide a more precise and efficient scope management approach by assigning unique key identifiers to each scope, achieving isolation and lookup of variables, functions, and types.

## 2. Design Goals

1. **Unique Identification**: Each scope has a unique key identifier
2. **Hierarchical Management**: Support for nested and hierarchical scope relationships
3. **Efficient Lookup**: Rapid location and access to contexts through keys
4. **Runtime Support**: Creation of corresponding Context objects at runtime
5. **Compatibility**: Seamless integration with existing compilation and execution processes

## 3. Core Concepts

### 3.1 Scope Key Naming Convention

Each scope has a unique key following this naming convention:

- **Global scope**: `main`
- **Main function**: `main.main`
- **Regular functions**: `main.FunctionName`
- **Struct methods**: `main.StructName.MethodName`
- **Other modules**: `moduleName.FunctionName`

### 3.2 Compile-time Context (Compile Context)

During compilation, each scope has a corresponding compile-time context containing:
- Unique key identifier
- Reference to parent context
- Variable, function, and type information within the current scope

### 3.3 Runtime Context (Runtime Context)

At runtime, each scope corresponds to a Context object containing:
- Unique key identifier
- Reference to parent Context
- Local variables within the current scope

## 4. Implementation Details

### 4.1 Compiler Modifications

#### 4.1.1 Scope Path Tracking

The compiler needs to maintain the current scope path for generating unique keys:

```go
type Compiler struct {
    // ... existing fields ...
    currentScopePath []string // Track current scope path, e.g., ["main", "main"] for main function
}
```

#### 4.1.2 Scope Key Generation

Provide helper functions to generate unique keys for the current scope:

```go
func (c *Compiler) getCurrentScopeKey() string {
    if len(c.currentScopePath) == 0 {
        return "main"
    }
    return strings.Join(c.currentScopePath, ".")
}
```

#### 4.1.3 BlockStmt Compilation

When compiling BlockStmt, it's necessary to:
1. Analyze the current context's key
2. Generate scope management instructions
3. Compile statements within the block
4. Generate scope exit instructions

### 4.2 Virtual Machine Modifications

#### 4.2.1 Context Management

The virtual machine uses a hierarchical Context system:

```go
type VM struct {
    // ... existing fields ...
    GlobalCtx *context.Context  // Global context
    currentCtx *context.Context // Current context
}
```

#### 4.2.2 Context Structure

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

#### 4.2.3 Scope Management Instructions

Scope management is handled through context operations:

```go
const (
    // ... existing opcodes ...
    OpEnterScopeWithKey // Enter scope with specified key
    OpExitScopeWithKey  // Exit scope with specified key
)
```

### 4.3 Instruction Processing

#### 4.3.1 OpEnterScopeWithKey

Processing instructions for entering a scope with a specified key:
1. Get scope key from instruction parameters
2. Save current context
3. Look up or create corresponding context for the key
4. Set as current context

#### 4.3.2 OpExitScopeWithKey

Processing instructions for exiting a scope with a specified key:
1. Get scope key from instruction parameters
2. Verify that current context is correct
3. Restore parent context
4. Set as current context

## 5. Execution Flow

### 5.1 Compilation Phase

1. Compiler adds function name to scope path when compiling functions
2. Generate scope management instructions when compiling BlockStmt
3. Generate unique key identifiers for each scope

### 5.2 Runtime Phase

1. Virtual machine creates Context when executing OpEnterScopeWithKey instruction
2. Variable storage and lookup occur in the current Context
3. Virtual machine reverts to parent context when executing OpExitScopeWithKey instruction

## 6. Examples

### 6.1 Code Example

```go
package main

func main() {
    x := 10  // Variable stored in Context with key "main.main"
    
    func() {
        y := 20  // Variable stored in Context of anonymous function, parent is "main.main"
    }()
}
```

### 6.2 Generated Instruction Sequence

```
OpEnterScopeWithKey "main.main"  // Enter main function scope
OpLoadConst 10                   // Load constant 10
OpStoreName "x"                  // Store variable x
OpEnterScopeWithKey "main.main.func1"  // Enter anonymous function scope
OpLoadConst 20                   // Load constant 20
OpStoreName "y"                  // Store variable y
OpExitScopeWithKey "main.main.func1"   // Exit anonymous function scope
OpExitScopeWithKey "main.main"         // Exit main function scope
```

## 7. Advantages

1. **Precise Scope Management**: Each scope has a unique identifier, avoiding naming conflicts
2. **Efficient Variable Lookup**: Direct context location through keys reduces lookup time
3. **Clear Hierarchical Relationships**: Parent-child context references are clear, facilitating variable lookup and scope management
4. **Modular Support**: Natural support for isolation of different modules
5. **Debugging Friendly**: Keys clearly indicate the current scope

## 8. Considerations

1. **Performance Considerations**: Need to balance context management overhead and lookup efficiency
2. **Memory Management**: Timely cleanup of unused context objects
3. **Error Handling**: Ensure correct scope switching to avoid stack overflow or underflow
4. **Compatibility**: Ensure new mechanisms are compatible with existing code
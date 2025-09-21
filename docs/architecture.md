# GoScript 综合技术文档

## 1. 项目概述

GoScript是一个兼容Go标准语法的脚本引擎，它允许你在Go应用程序中动态执行Go代码。该项目旨在实现一个安全、高效且易于扩展的脚本执行环境。

### 1.1 设计目标

1. **语法兼容性**：尽可能兼容Go标准语法
2. **安全性**：提供沙箱环境，限制危险操作
3. **可扩展性**：支持自定义函数、类型和模块
4. **性能**：通过编译执行提高运行效率
5. **易用性**：提供简洁的API供Go程序调用

### 1.2 核心特性

- **模块化设计**：词法分析、语法分析、AST生成等组件分离
- **可扩展性**：支持自定义函数和模块
- **安全性**：提供执行时间和内存使用限制
- **复用Go原生模块**：词法分析、语法分析等直接复用Go标准库

## 2. 架构设计

### 2.1 整体架构

```
+-------------------+
|   Go Application  |
+-------------------+
          |
          v
+-------------------+
|  Script Interface | <- API层，提供给Go应用调用
+-------------------+
          |
          v
+-------------------+
|   Parser/Lexer    | <- 词法语法分析器（复用Go标准库）
+-------------------+
          |
          v
+-------------------+
|   AST Generator   | <- 抽象语法树生成器（复用Go标准库）
+-------------------+
          |
          v
+-------------------+
|   Compiler        | <- 编译器，将AST编译为字节码
+-------------------+
          |
          v
+-------------------+
|   Bytecode VM     | <- 字节码虚拟机，执行编译后的代码
+-------------------+
          |
          v
+-------------------+
|  Runtime System   | <- 运行时系统，管理对象、内存等
+-------------------+
```

### 2.2 核心组件

1. **词法分析器 (parser)**：使用Go标准库的`go/scanner`进行词法分析
2. **语法分析器 (parser)**：使用Go标准库的`go/parser`进行语法分析
3. **抽象语法树 (ast)**：使用Go标准库的`go/ast`处理AST节点
4. **编译器 (compiler)**：将AST编译为可执行的中间表示（字节码）
5. **运行时 (runtime)**：管理变量、函数、类型和模块
6. **虚拟机 (vm)**：执行编译后的字节码
7. **类型系统 (types)**：统一的类型系统，所有类型都实现IType接口
8. **符号表 (symbol)**：管理变量、函数、类型等符号信息
9. **模块管理 (module)**：管理不同模块及其指令集
10. **执行上下文 (context)**：管理脚本执行时的变量作用域和栈

### 2.3 核心设计理念

1. **利用Go标准库**：充分利用Go标准库的功能，特别是`context`包来管理执行上下文和变量作用域
2. **简化操作码设计**：通过简化操作码设计，减少虚拟机的复杂性，提高执行效率
3. **模块化架构**：采用模块化设计，各个组件职责清晰，便于维护和扩展

## 3. 类型系统

### 3.1 IType接口

所有类型都支持IType接口，提供统一的类型操作接口：

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

### 3.2 基本类型实现

- IntType: 整数类型
- FloatType: 浮点数类型
- StringType: 字符串类型
- BoolType: 布尔类型

## 4. 执行上下文与作用域管理

### 4.1 ExecutionContext结构

```go
type ExecutionContext struct {
    // Go context for cancellation, timeout, and value storage
    Context context.Context

    // Cancel function to cancel the context
    Cancel context.CancelFunc

    // Scope manager for variable scope management
    ScopeManager *ScopeManager

    // Module name
    ModuleName string

    // Parent execution context
    Parent *ExecutionContext

    // Security context
    Security *SecurityContext
}
```

### 4.2 作用域嵌套

```
全局作用域
└── 模块作用域
    └── 函数作用域
        └── 块作用域
```

### 4.3 变量查找与隔离

利用Go的`context`包实现自然的变量查找链，从当前作用域向上查找直到全局作用域。不同作用域之间的变量自然隔离，防止变量污染。

## 5. 虚拟机与操作码

### 5.1 简化后的操作码

```go
const (
    OpNop        OpCode = iota // 空操作
    OpLoadConst               // 加载常量
    OpLoadName                // 加载变量
    OpStoreName               // 存储变量
    OpCall                    // 调用函数
    OpReturn                  // 返回
    OpJump                    // 跳转
    OpJumpIf                  // 条件跳转
    OpBinaryOp                // 二元操作
    OpUnaryOp                 // 一元操作
)
```

### 5.2 指令格式

```go
type Instruction struct {
    Op   OpCode      // 操作码
    Arg  interface{} // 参数1
    Arg2 interface{} // 参数2
}
```

## 6. 函数注册表机制

### 6.1 统一函数调用

所有函数（内置函数、用户定义函数、模块函数）都通过相同的机制注册和调用。

### 6.2 函数注册流程

1. 创建函数实例
2. 通过ExecutionContext.RegisterFunction注册函数
3. 函数存储在ScopeManager的函数注册表中

### 6.3 函数调用流程

1. 检查是否为模块函数调用 (moduleName.functionName)
2. 在全局上下文中查找函数
3. 在当前模块中查找函数

## 7. 模块系统

### 7.1 Module结构

```go
type Module struct {
    Name         string                    // 模块名称
    Instructions []*vm.Instruction         // 指令集
    SymbolTable  *symbol.SymbolTable       // 符号表
    Context      *context.ExecutionContext  // 执行上下文
    Functions    map[string]Function       // 函数映射
}
```

### 7.2 模块管理

支持模块定义和注册、函数注册、模块间调用、内置模块支持和模块访问控制。

## 8. 语法支持

### 8.1 支持的Go语法特性

1. **基本类型**：
   - 布尔类型 (bool)
   - 数值类型 (int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64)
   - 字符串类型 (string)
   - 复合类型 ([]T, [n]T, map[K]T, struct, interface)

2. **变量声明**：
   - var 声明
   - 短变量声明 (:=)
   - 常量声明 (const)

3. **控制结构**：
   - 条件语句 (if/else)
   - 循环语句 (for)
   - switch语句
   - goto语句

4. **函数**：
   - 函数声明
   - 函数调用
   - 多返回值
   - 可变参数
   - 匿名函数
   - 闭包

5. **数据结构**：
   - 数组和切片
   - 映射
   - 结构体
   - 指针

6. **错误处理**：
   - error类型
   - panic/recover机制

### 8.2 限制的语法特性

为了安全性和简化实现，以下Go特性将被限制或不支持：

1. **包管理**：不支持完整的Go包系统，通过模块系统提供功能
2. **并发**：限制或不支持goroutine和channel，可通过模块提供受限的并发功能
3. **低级操作**：不支持unsafe包，限制指针操作
4. **反射**：限制或不支持reflect包
5. **系统调用**：通过模块提供受限的系统功能

## 9. 安全机制

### 9.1 沙箱环境

1. **资源限制**：
   - 最大执行时间限制
   - 最大内存使用限制
   - 最大对象分配限制

2. **API限制**：
   - 禁止危险系统调用
   - 限制文件系统访问
   - 限制网络访问

3. **语法限制**：
   - 禁止某些关键字
   - 限制复杂度

### 9.2 安全上下文

```go
type SecurityContext struct {
    // 最大执行时间（纳秒）
    MaxExecutionTime time.Duration
    
    // 最大内存使用（字节）
    MaxMemoryUsage int64
    
    // 允许的模块列表
    AllowedModules []string
    
    // 禁止的关键字
    ForbiddenKeywords []string
}

// 设置安全上下文
func (s *Script) SetSecurityContext(ctx *SecurityContext)
```

## 10. API设计

### 10.1 Script接口

```go
type Script struct {
    // 脚本内容
    source []byte
    
    // 变量映射
    variables map[string]interface{}
    
    // 模块映射
    modules map[string]Module
}

// 创建新脚本
func NewScript(source []byte) *Script

// 添加变量
func (s *Script) AddVariable(name string, value interface{}) error

// 添加模块
func (s *Script) AddModule(name string, module Module) error

// 编译脚本
func (s *Script) Compile() (*CompiledScript, error)

// 执行脚本
func (s *Script) Run() (interface{}, error)

// 执行脚本（带上下文）
func (s *Script) RunContext(ctx context.Context) (interface{}, error)
```

## 11. 使用示例

### 11.1 基本使用

```go
package main

import (
    "fmt"
    "context"
    "time"
    "github.com/lengzhao/goscript"
)

func main() {
    // 创建脚本
    source := `
package main

func main() {
    x := 10
    y := 20
    return x + y
}
`
    
    script := goscript.NewScript([]byte(source))
    
    // 设置安全上下文
    securityCtx := &goscript.SecurityContext{
        MaxExecutionTime: 5 * time.Second,
        MaxMemoryUsage:   10 * 1024 * 1024, // 10MB
    }
    script.SetSecurityContext(securityCtx)
    
    // 执行脚本
    result, err := script.RunContext(context.Background())
    if err != nil {
        fmt.Printf("执行错误: %v\n", err)
        return
    }
    
    fmt.Printf("执行结果: %v\n", result) // 输出: 30
}
```

### 11.2 自定义扩展

```go
// 自定义函数
func customFunction(args ...interface{}) (interface{}, error) {
    if len(args) != 2 {
        return nil, fmt.Errorf("需要2个参数")
    }
    
    a, ok1 := args[0].(int)
    b, ok2 := args[1].(int)
    if !ok1 || !ok2 {
        return nil, fmt.Errorf("参数必须是整数")
    }
    
    return a * b, nil
}

// 使用自定义函数
func main() {
    source := `
package main

func main() {
    result := customMultiply(5, 6)
    return result
}
`
    
    script := goscript.NewScript([]byte(source))
    script.AddFunction("customMultiply", runtime.NewBuiltInFunction("customMultiply", customFunction))
    
    result, err := script.Run()
    if err != nil {
        fmt.Printf("执行错误: %v\n", err)
        return
    }
    
    fmt.Printf("执行结果: %v\n", result) // 输出: 30
}
```

## 12. 执行流程

1. **词法分析**：源代码 → Tokens（复用Go标准库）
2. **语法分析**：Tokens → AST（复用Go标准库）
3. **编译**：AST → 字节码指令和常量池
4. **执行**：
   - 创建执行上下文
   - 加载模块和函数
   - 执行字节码指令
   - 管理作用域和变量

## 13. 性能优化

1. **操作码优化**：通过简化操作码设计，减少虚拟机的复杂性，提高执行效率
2. **作用域管理优化**：利用Go的`context`包优化作用域管理，减少内存分配和查找时间
3. **标准库复用**：复用Go标准库的词法分析、语法分析和AST处理功能，确保兼容性和性能

## 14. 扩展机制

### 14.1 自定义函数

支持通过Go代码注册自定义函数：

```go
// 注册自定义函数
script.AddFunction("myFunc", func(args ...interface{}) (interface{}, error) {
    // 函数实现
    return result, nil
})
```

### 14.2 自定义类型

支持实现自定义类型：

```go
// 实现自定义类型
type MyType struct {
    // 字段定义
}

func (m *MyType) TypeName() string {
    return "MyType"
}

func (m *MyType) String() string {
    // 字符串表示
}

// 注册自定义类型
script.AddType("MyType", &MyType{})
```

### 14.3 模块系统

支持模块化扩展：

```go
// 创建模块
module := NewModule("myModule")
module.AddFunction("func1", func1)
module.AddType("Type1", &Type1{})

// 注册模块
script.AddModule("myModule", module)
```

## 15. 测试

运行所有测试：

```
go test ./...
```

运行特定包的测试：

```
go test ./lexer -v
go test ./parser -v
go test ./ast -v
go test ./compiler -v
go test ./runtime -v
go test ./vm -v
```

## 16. 总结

GoScript通过以下方式实现了简洁、高效和安全的脚本引擎：

1. **利用Go标准库**：复用成熟的Go标准库功能
2. **简化设计**：通过简化操作码和组件设计降低复杂性
3. **自然的作用域管理**：利用Go的`context`包实现自然的作用域管理
4. **模块化架构**：清晰的组件职责划分便于维护和扩展
5. **内置安全机制**：提供多层次的安全控制

这种设计使得GoScript成为一个易于使用、高性能且安全的脚本引擎，适用于各种Go应用程序的动态执行需求。
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
- **基于Key的上下文管理**：通过唯一标识符管理作用域和变量
- **完整的Go语法支持**：支持结构体、方法、range语句、复合字面量等

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

GoScript的核心组件包括：

1. **Script (script.go)**：主要的API接口，提供NewScript、Run、AddFunction等方法
2. **Parser (parser/)**：词法和语法分析器，复用Go标准库的`go/scanner`和`go/parser`
3. **Compiler (compiler/)**：编译器，将AST编译为可执行的中间表示（字节码）
4. **VM (vm/)**：虚拟机，执行编译后的字节码
5. **Context (context/)**：执行上下文管理，管理脚本执行时的变量作用域和栈
6. **Instruction (instruction/)**：指令定义，定义虚拟机可执行的操作码
7. **Types (types/)**：类型系统，定义GoScript中的类型接口和模块执行器
8. **Builtin (builtin/)**：内置函数和模块，提供标准库功能如math、strings等

这些组件协同工作，提供了一个完整的脚本执行环境。

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

### 4.1 Context结构

当前实现使用基于`context.Context`包的分层上下文系统：

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

### 4.2 作用域嵌套

```
全局作用域
└── 模块作用域
    └── 函数作用域
        └── 块作用域
```

### 4.3 变量查找与隔离

利用Context结构实现自然的变量查找链，从当前作用域向上查找直到全局作用域。不同作用域之间的变量自然隔离，防止变量污染。

### 4.4 基于Key的上下文管理（新设计）

为了更好地管理作用域，引入了基于Key的上下文管理机制：

1. **唯一标识符**：每个作用域都有唯一的key标识符
   - 全局作用域：`main`
   - 主函数：`main.main`
   - 普通函数：`main.FunctionName`
   - 结构体方法：`main.StructName.MethodName`
   - 其他模块：`moduleName.FunctionName`
   - 代码块：`main.main.block_1`

2. **编译时跟踪**：编译器在编译BlockStmt时分析当前上下文的key，并生成相应的作用域管理指令

3. **运行时管理**：运行时创建Context对象来管理变量和引用关系，每个Context引用其父级上下文

4. **变量查找**：变量查找遵循作用域链，从当前上下文向上查找直到全局上下文

5. **作用域隔离**：不同作用域之间的变量自然隔离，防止变量污染

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
    OpEnterScope              // 进入作用域
    OpExitScope               // 退出作用域
    OpEnterScopeWithKey       // 进入指定key的作用域
    OpExitScopeWithKey        // 退出指定key的作用域
    OpCreateVar               // 创建变量
    OpNewSlice                // 创建切片
    OpNewStruct               // 创建结构体
    OpGetField                // 获取结构体字段
    OpSetField                // 设置结构体字段
    OpGetIndex                // 获取索引元素
    OpSetIndex                // 设置索引元素
    OpLen                     // 获取长度
    OpImport                  // 导入模块
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
2. 通过VM.RegisterFunction注册函数
3. 函数存储在VM的函数注册表中

### 6.3 函数调用流程

1. 检查是否为模块函数调用 (moduleName.functionName)
2. 在VM的函数注册表中查找函数
3. 直接执行函数或通过模块执行器执行

## 7. 模块系统

### 7.1 模块结构

模块通过VM的模块注册系统进行管理：

```go
// 模块函数通过ModuleExecutor接口注册
type ModuleExecutor func(entrypoint string, args ...interface{}) (interface{}, error)
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
   - range语句 (for range)
   - switch语句
   - goto语句

4. **函数**：
   - 函数声明
   - 函数调用
   - 方法声明和调用
   - 多返回值
   - 可变参数
   - 匿名函数
   - 闭包

5. **数据结构**：
   - 数组和切片
   - 映射
   - 结构体
   - 指针
   - 复合字面量 ([]int{1, 2, 3} 或 Person{name: "Alice"})

6. **操作符**：
   - 算术操作符 (+, -, *, /, %)
   - 比较操作符 (==, !=, <, <=, >, >=)
   - 逻辑操作符 (&&, ||, !)
   - 赋值操作符 (=, +=, -=, *=, /=, %=)
   - 自增自减操作符 (++, --)

7. **错误处理**：
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
// 安全性通过VM指令限制进行管理
// 脚本级别的安全配置
type Script struct {
    // 最大允许的指令数（0表示无限制）
    maxInstructions int64
}

// 设置安全上下文
func (s *Script) SetMaxInstructions(max int64)
```

## 10. API设计

### 10.1 Script接口

```go
type Script struct {
    // 脚本内容
    source []byte
    
    // 虚拟机
    vm *vm.VM

    // 调试模式
    debug bool

    // 执行统计信息
    executionStats *ExecutionStats

    // 最大允许的指令数（0表示无限制）
    maxInstructions int64
}

// 创建新脚本
func NewScript(source []byte) *Script

// 添加函数
func (s *Script) AddFunction(name string, fn vm.ScriptFunction) error

// 注册模块
func (s *Script) RegisterModule(moduleName string, executor types.ModuleExecutor)

// 编译和执行脚本
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
    script.SetDebug(true) // 启用调试模式
    
    // 执行脚本
    result, err := script.Run()
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
    
    // 通过VM注册自定义函数
    script.AddFunction("customMultiply", func(args ...interface{}) (interface{}, error) {
        if len(args) != 2 {
            return nil, fmt.Errorf("customMultiply函数需要2个参数")
        }
        a, ok1 := args[0].(int)
        b, ok2 := args[1].(int)
        if !ok1 || !ok2 {
            return nil, fmt.Errorf("customMultiply函数需要整数参数")
        }
        return a * b, nil
    })
    
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
   - 创建VM上下文
   - 加载模块和函数
   - 执行字节码指令
   - 管理作用域和变量

## 13. 性能优化

1. **操作码优化**：通过简化操作码设计，减少虚拟机的复杂性，提高执行效率
2. **作用域管理优化**：利用分层Context对象优化作用域管理，减少内存分配和查找时间
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

支持通过IType接口实现自定义类型：

```go
// 通过实现IType接口实现自定义类型
type MyType struct {
    // 字段定义
}

func (m *MyType) TypeName() string {
    return "MyType"
}

func (m *MyType) String() string {
    // 字符串表示
}

// 通过自定义函数或模块注册自定义类型
```

### 14.3 模块系统

支持模块化扩展：

```go
// 创建并通过ModuleExecutor注册模块
moduleExecutor := func(entrypoint string, args ...interface{}) (interface{}, error) {
    // 模块实现
    return result, nil
}

// 注册模块
script.RegisterModule("myModule", moduleExecutor)
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
3. **自然的作用域管理**：利用分层Context对象实现自然的作用域管理
4. **模块化架构**：清晰的组件职责划分便于维护和扩展
5. **内置安全机制**：提供指令数限制等安全控制

这种设计使得GoScript成为一个易于使用、高性能且安全的脚本引擎，适用于各种Go应用程序的动态执行需求。
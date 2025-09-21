# GoScript - Go兼容脚本引擎

GoScript是一个兼容Go标准语法的脚本引擎，它允许你在Go应用程序中动态执行Go代码。

## 架构设计

GoScript的架构设计详细信息请参见 [最终架构设计文档](docs/architecture_final.md)。

## 特性

- **语法兼容性**：尽可能兼容Go标准语法
- **模块化设计**：词法分析、语法分析、AST生成等组件分离
- **可扩展性**：支持自定义函数和模块
- **安全性**：提供执行时间和内存使用限制
- **复用Go原生模块**：词法分析、语法分析等直接复用Go标准库

## 架构

GoScript的架构包括以下核心组件：

1. **词法分析器 (lexer)**：使用Go标准库的`go/scanner`进行词法分析
2. **语法分析器 (parser)**：使用Go标准库的`go/parser`进行语法分析
3. **抽象语法树 (ast)**：使用Go标准库的`go/ast`处理AST节点
4. **编译器 (compiler)**：将AST编译为可执行的中间表示（字节码）
5. **运行时 (runtime)**：管理变量、函数、类型和模块
6. **虚拟机 (vm)**：执行编译后的字节码
7. **类型系统 (types)**：统一的类型系统，所有类型都实现IType接口
8. **符号表 (symbol)**：管理变量、函数、类型等符号信息
9. **模块管理 (module)**：管理不同模块及其指令集
10. **执行上下文 (context)**：管理脚本执行时的变量作用域和栈

## 安装

```bash
go get github.com/lengzhao/goscript
```

## 快速开始

### 基本用法

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

### 自定义函数

```go
// 创建自定义函数
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

## 核心组件

### 词法分析器

词法分析器将源代码分解为标记(tokens)：

```go
tokens, err := script.Lex()
if err != nil {
    // 处理错误
}

for _, token := range tokens {
    fmt.Printf("Token: %s, Value: %s\n", token.TokenType(), token.Value())
}
```

### 语法分析器

语法分析器将标记转换为抽象语法树(AST)：

```go
ast, err := script.Parse()
if err != nil {
    // 处理错误
}

// 遍历AST节点
ast.Inspect(ast, func(n ast.Node) bool {
    // 处理节点
    return true
})
```

### 编译器 (Compiler)

编译器将AST编译为字节码：

```go
// 创建编译器
compiler := compiler.NewCompiler(runtime)

// 编译AST
bytecode, constants, err := compiler.Compile(astFile)
if err != nil {
    // 处理错误
}

// bytecode包含编译后的指令
// constants包含常量池
```

### 运行时 (Runtime)

运行时管理脚本执行环境中的变量、函数、类型和模块：

```go
// 创建运行时
rt := runtime.NewRuntime()

// 注册自定义函数
customFunc := runtime.NewBuiltInFunction("multiply", func(args ...interface{}) (interface{}, error) {
    // 函数实现
    return result, nil
})
rt.RegisterFunction("multiply", customFunc)

// 设置变量
rt.SetVariable("testVar", 42)
```

### 虚拟机 (VM)

虚拟机负责执行编译后的字节码：

```go
// 创建虚拟机
vmInstance := vm.NewVM(runtime)

// 执行二进制操作
result, err := vmInstance.BinaryOperation(10, 20, "+")

// 调用函数
result, err := vmInstance.CallFunction("multiply", 5, 6)

// 栈操作
vmInstance.Push(42)
value := vmInstance.Pop()
```

## 安全机制

GoScript提供多种安全机制来限制脚本执行：

```go
securityCtx := &goscript.SecurityContext{
    MaxExecutionTime: 5 * time.Second,     // 最大执行时间
    MaxMemoryUsage:   10 * 1024 * 1024,    // 最大内存使用 (10MB)
    AllowedModules:   []string{"fmt"},     // 允许的模块
    ForbiddenKeywords: []string{"unsafe"}, // 禁止的关键字
}
script.SetSecurityContext(securityCtx)
```

## 测试

运行所有测试：

```bash
go test ./...
```

运行特定包的测试：

```bash
go test ./lexer -v
go test ./parser -v
go test ./ast -v
go test ./compiler -v
go test ./runtime -v
go test ./vm -v
```

## 示例

查看`examples/`目录中的示例程序：

- `examples/basic/` - 基本用法示例
- `examples/custom_functions/` - 自定义函数示例
- `examples/compiler/` - 编译器功能示例
- `examples/security/` - 安全机制示例
- `examples/custom_types/` - 自定义类型和模块示例
- `examples/error_handling/` - 错误处理示例
- `examples/comprehensive/` - 综合功能示例
- `examples/modules/` - 模块化示例

运行演示程序：

```bash
cd cmd/demo
go run main.go
```

运行特定示例：

```bash
cd examples/compiler
go run main.go
```

## 贡献

欢迎贡献代码！请遵循以下步骤：

1. Fork项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启Pull Request

## 许可证

本项目采用MIT许可证 - 查看[LICENSE](LICENSE)文件了解详情

## 联系

项目链接: [https://github.com/lengzhao/goscript](https://github.com/lengzhao/goscript)
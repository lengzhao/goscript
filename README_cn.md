# GoScript - Go兼容脚本引擎

GoScript是一个兼容Go标准语法的脚本引擎，它允许你在Go应用程序中动态执行Go代码。

## 特性

- **语法兼容性**：尽可能兼容Go标准语法
- **模块化设计**：词法分析、语法分析、AST生成等组件分离
- **可扩展性**：支持自定义函数和模块
- **复用Go原生模块**：词法分析、语法分析等直接复用Go标准库

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
    
    // 注册自定义函数
    err := script.AddFunction("customMultiply", func(args ...interface{}) (interface{}, error) {
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
    if err != nil {
        fmt.Printf("注册自定义函数失败: %v\n", err)
        return
    }
    
    result, err := script.Run()
    if err != nil {
        fmt.Printf("执行错误: %v\n", err)
        return
    }
    
    fmt.Printf("执行结果: %v\n", result) // 输出: 30
}
```

### 使用内置模块

GoScript提供了几个内置模块，包括`math`、`strings`、`fmt`和`json`：

```go
func main() {
    // 使用strings模块
    source := `
package main

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
    
    // 注册内置模块
    modules := []string{"strings"}
    for _, moduleName := range modules {
        moduleExecutor, exists := builtin.GetModuleExecutor(moduleName)
        if exists {
            script.RegisterModule(moduleName, moduleExecutor)
        }
    }
    
    result, err := script.Run()
    if err != nil {
        fmt.Printf("执行错误: %v\n", err)
        return
    }
    
    fmt.Printf("执行结果: %v\n", result)
}
```

### 直接调用函数

```go
func main() {
    // 创建新脚本
    script := goscript.NewScript([]byte{})
    
    // 使用AddFunction方法添加函数
    script.AddFunction("addFunc", func(args ...interface{}) (interface{}, error) {
        if len(args) != 2 {
            return nil, fmt.Errorf("addFunc需要2个参数")
        }
        a, ok1 := args[0].(int)
        b, ok2 := args[1].(int)
        if !ok1 || !ok2 {
            return nil, fmt.Errorf("addFunc需要整数参数")
        }
        return a + b, nil
    })
    
    // 直接调用函数
    result, err := script.CallFunction("addFunc", 5, 6)
    if err != nil {
        fmt.Printf("调用函数失败: %v\n", err)
        return
    }
    
    fmt.Printf("函数结果: %v\n", result) // 输出: 11
}
```

## 核心组件

### 脚本 (Script)

GoScript引擎的主要接口。关键方法包括：

- `NewScript(source []byte) *Script` - 创建新脚本
- `Run() (interface{}, error)` - 执行脚本
- `AddFunction(name string, execFn vm.ScriptFunction) error` - 添加自定义函数
- `CallFunction(name string, args ...interface{}) (interface{}, error)` - 直接调用函数
- `SetDebug(debug bool)` - 启用或禁用调试模式
- `RegisterModule(moduleName string, executor types.ModuleExecutor)` - 注册模块
- `SetMaxInstructions(max int64)` - 设置最大指令数（默认值：10000）

### 虚拟机 (VM)

虚拟机负责执行编译后的字节码：

```go
// 创建虚拟机
vmInstance := vm.NewVM()

// 注册函数
vmInstance.RegisterFunction("multiply", func(args ...interface{}) (interface{}, error) {
    // 函数实现
    return result, nil
})

// 调用函数
result, err := vmInstance.Execute("main.main", arg1, arg2)

// 获取函数
fn, exists := vmInstance.GetFunction("functionName")
```

### 内置模块

GoScript提供以下内置模块：

1. **strings** - 字符串操作函数
2. **math** - 数学函数
3. **fmt** - 格式化函数
4. **json** - JSON编码/解码函数

## 安全特性

GoScript提供了多种安全机制来防止脚本滥用系统资源：

### 1. 指令数限制
限制脚本可以执行的最大指令数。默认限制为10000条指令：

```go
script := goscript.NewScript(source)
// 默认为10000条指令
// script.SetMaxInstructions(10000)

// 设置自定义限制
script.SetMaxInstructions(5000) // 限制为5000条指令

// 移除限制
script.SetMaxInstructions(0) // 无限制
```

### 2. 模块访问控制
通过选择性注册模块来控制脚本可以访问的模块。

## 测试

运行所有测试：

```bash
go test ./...
```

运行特定包的测试：

```bash
go test ./test -v
```

## 示例

查看`examples/`目录中的示例程序：

- `examples/basic/` - 基本用法示例
- `examples/custom_function/` - 自定义函数示例
- `examples/builtin_functions/` - 内置函数示例
- `examples/modules/` - 模块使用示例
- `examples/interface_example/` - 接口示例
- `examples/struct_example/` - 结构体示例
- `test/data` - 各种功能的测试脚本

运行示例：

```bash
cd examples/custom_function
go run function_demo.go
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
# 函数注册表设计说明

## 1. 概述

本文档说明了GoScript中函数注册表的设计和实现。函数注册表提供了一个统一的机制来注册和调用各种类型的函数，包括内置函数、用户定义函数和模块函数。

## 2. 设计目标

1. **统一接口**：为所有类型的函数提供统一的注册和调用接口
2. **作用域管理**：支持函数的作用域管理，确保正确的函数查找
3. **性能优化**：提供高效的函数查找和调用机制
4. **可扩展性**：支持轻松添加新的函数类型

## 3. 核心组件

### 3.1 ExecutionContext
ExecutionContext是执行上下文，包含ScopeManager和函数注册表。

**主要方法**：
- `RegisterFunction(name string, fn Function) error` - 注册函数
- `GetFunction(name string) (Function, bool)` - 获取函数
- `WithTimeout(timeout time.Duration) *ExecutionContext` - 创建带超时的执行上下文

### 3.2 ScopeManager
ScopeManager管理变量作用域和函数注册表。

**主要方法**：
- `RegisterFunction(name string, fn Function)` - 注册函数
- `GetFunction(name string) (Function, bool)` - 获取函数
- `GetAllFunctions() map[string]Function` - 获取所有函数

### 3.3 Virtual Machine
虚拟机使用函数注册表来执行函数调用指令。

**主要方法**：
- `RegisterFunction(name string, fn func(args ...interface{}) (interface{}, error))` - 注册函数
- `SetFunctionRegistry(registry map[string]func(args ...interface{}) (interface{}, error))` - 设置函数注册表
- `GetFunctionRegistry() map[string]func(args ...interface{}) (interface{}, error)` - 获取函数注册表

## 4. 函数注册流程

```go
// 1. 创建脚本
script := NewScript(source)

// 2. 定义函数
myFunction := &SimpleFunction{
    name: "myFunc",
    fn: func(args ...interface{}) (interface{}, error) {
        // 函数实现
        return result, nil
    },
}

// 3. 注册函数
script.AddFunction("myFunc", myFunction)
```

## 5. 函数调用流程

```go
// 1. 调用函数
result, err := script.CallFunction("myFunc", arg1, arg2)

// 2. 函数查找顺序：
//    a. 检查是否为模块函数调用 (moduleName.functionName)
//    b. 在全局上下文中查找函数
//    c. 在当前模块中查找函数
```

## 6. 统一函数调用机制

所有函数（内置函数、用户定义函数、模块函数）都通过相同的机制注册和调用：

1. **内置函数**：在模块初始化时注册
2. **用户定义函数**：通过`AddFunction`方法注册
3. **模块函数**：在模块中注册并通过模块管理器调用

## 7. 示例代码

### 7.1 注册用户定义函数

```go
// 定义自定义函数
type CustomFunction struct {
    name string
    fn   func(args ...interface{}) (interface{}, error)
}

func (f *CustomFunction) Call(args ...interface{}) (interface{}, error) {
    return f.fn(args...)
}

func (f *CustomFunction) Name() string {
    return f.name
}

// 注册函数
script.AddFunction("customFunc", &CustomFunction{
    name: "customFunc",
    fn: func(args ...interface{}) (interface{}, error) {
        // 函数实现
        return "result", nil
    },
})
```

### 7.2 在虚拟机中调用函数

```go
// 在虚拟机中注册函数
vm.RegisterFunction("add", func(args ...interface{}) (interface{}, error) {
    if len(args) != 2 {
        return nil, fmt.Errorf("add function requires 2 arguments")
    }
    a, ok1 := args[0].(int)
    b, ok2 := args[1].(int)
    if !ok1 || !ok2 {
        return nil, fmt.Errorf("add function requires integer arguments")
    }
    return a + b, nil
})
```

## 8. 总结

函数注册表设计提供了一个灵活、统一的函数管理机制，支持各种类型的函数注册和调用。通过ExecutionContext和ScopeManager的配合，实现了作用域管理和函数查找，确保了函数调用的正确性和高效性。
# GoScript 精细化设计文档

## 1. 概述

本文档详细描述了GoScript的精细化设计，包括类型系统、符号表、操作码、执行上下文和模块化执行等核心组件。

## 2. 核心设计要点

### 2.1 IType接口定义
所有类型都支持IType接口，提供统一的类型操作接口。

### 2.2 符号表设计
符号表管理变量、函数、类型等符号信息，支持作用域嵌套和模块化管理。

### 2.3 操作码简化
操作码被简化为核心10个指令，提高虚拟机执行效率。

### 2.4 执行上下文优化
执行器引入ctx，所有的变量都存到ctx里，进入函数/子作用域，通过ScopeManager创建子作用域，从而实现变量隔离。

### 2.5 模块化执行
不同的模块对应不同的指令集，模块间调用通过模块名由执行器索引并调用。

## 3. 详细设计

### 3.1 类型系统 (types)

#### 3.1.1 IType接口
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

#### 3.1.2 基本类型实现
- IntType: 整数类型
- FloatType: 浮点数类型
- StringType: 字符串类型
- BoolType: 布尔类型

### 3.2 符号表 (symbol)

#### 3.2.1 Symbol结构
```go
type Symbol struct {
    Module     string  // 所属模块
    ID         string  // 唯一标识符
    Name       string  // 符号名称
    Type       IType   // 符号类型
    Address    interface{} // 内存地址或值
    ScopeLevel int     // 作用域层级
    Mutable    bool    // 是否可变
}
```

#### 3.2.2 SymbolTable结构
```go
type SymbolTable struct {
    symbols map[string]*Symbol
    module  string
}
```

### 3.3 操作码 (vm)

#### 3.3.1 简化后的操作码
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

#### 3.3.2 指令格式
```go
type Instruction struct {
    Op   OpCode      // 操作码
    Arg  interface{} // 参数1
    Arg2 interface{} // 参数2
}
```

### 3.4 执行上下文 (context)

#### 3.4.1 ExecutionContext结构
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

#### 3.4.2 ScopeManager结构
```go
type ScopeManager struct {
    // Root scope
    root *Scope

    // Current scope
    current *Scope

    // Function registry
    functions map[string]Function

    // Cancel function for timeout
    cancelFunc func()

    // Deadline for timeout
    deadline time.Time

    // mutex for thread safety
    mu sync.RWMutex
}
```

#### 3.4.3 函数注册支持
ExecutionContext支持函数注册和调用：
- `RegisterFunction(name string, fn Function) error` - 注册函数
- `GetFunction(name string) (Function, bool)` - 获取函数
- `WithTimeout(timeout time.Duration) *ExecutionContext` - 创建带超时的执行上下文

### 3.5 模块管理 (module)

#### 3.5.1 Module结构
```go
type Module struct {
    Name         string                    // 模块名称
    Instructions []*vm.Instruction         // 指令集
    SymbolTable  *symbol.SymbolTable       // 符号表
    Context      *context.ExecutionContext  // 执行上下文
    Functions    map[string]Function       // 函数映射
}
```

#### 3.5.2 ModuleManager结构
```go
type ModuleManager struct {
    modules       map[string]*Module        // 模块映射
    currentModule string                    // 当前模块
    globalContext *context.ExecutionContext  // 全局执行上下文
}
```

#### 3.5.3 模块间调用
模块间调用通过模块名索引并调用：
- `CallModuleFunction(moduleName, functionName string, args ...interface{})` - 调用模块函数

## 4. 执行流程

### 4.1 初始化阶段
1. 创建Script实例
2. 初始化ModuleManager
3. 创建全局ExecutionContext
4. 初始化虚拟机VM

### 4.2 编译阶段
1. 词法分析：源代码 → Tokens
2. 语法分析：Tokens → AST
3. 代码生成：AST → 字节码指令和常量池

### 4.3 执行阶段
1. 创建执行上下文
2. 加载模块和函数
3. 执行字节码指令
4. 管理作用域和变量

## 5. 作用域管理

### 5.1 作用域嵌套
```
全局作用域
└── 模块作用域
    └── 函数作用域
        └── 块作用域
```

### 5.2 变量查找
利用ScopeManager实现自然的变量查找链，从当前作用域向上查找直到全局作用域。

### 5.3 变量隔离
不同作用域之间的变量自然隔离，防止变量污染。

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

## 7. 安全机制

### 7.1 执行时间限制
通过ExecutionContext的WithTimeout方法实现执行时间限制。

### 7.2 内存使用限制
通过SecurityContext的MaxMemoryUsage字段控制内存使用。

### 7.3 模块访问控制
通过SecurityContext的AllowedModules字段控制可访问的模块。

## 8. 扩展机制

### 8.1 自定义函数
支持注册自定义Go函数到ExecutionContext。

### 8.2 自定义类型
通过实现IType接口定义新类型。

### 8.3 模块扩展
支持创建自定义模块并注册到ModuleManager。

## 9. 性能优化

### 9.1 操作码优化
通过简化操作码设计，减少虚拟机的复杂性，提高执行效率。

### 9.2 作用域管理优化
利用自定义ScopeManager优化作用域管理，减少内存分配和查找时间。

### 9.3 函数调用优化
通过统一的函数注册表机制优化函数查找和调用。

## 10. 总结

GoScript的精细化设计通过以下方式实现了简洁、高效和安全的脚本引擎：

1. **简化的操作码设计**：仅使用10个核心操作码，降低虚拟机复杂性
2. **优化的作用域管理**：通过ScopeManager实现高效的作用域管理
3. **统一的函数调用机制**：所有函数通过相同的机制注册和调用
4. **模块化架构**：清晰的组件职责划分便于维护和扩展
5. **内置安全机制**：提供多层次的安全控制

这种设计使得GoScript成为一个易于使用、高性能且安全的脚本引擎，适用于各种Go应用程序的动态执行需求。
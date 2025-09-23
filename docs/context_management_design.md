# GoScript 上下文管理设计文档

## 1. 概述

本文档详细描述了GoScript中基于key的上下文管理机制设计。该机制旨在提供一种更精确和高效的作用域管理方式，通过为每个作用域分配唯一的key标识符，实现变量、函数和类型的隔离与查找。

## 2. 设计目标

1. **唯一标识**：每个作用域都有唯一的key标识符
2. **层级管理**：支持作用域的嵌套和层级关系
3. **高效查找**：通过key快速定位和访问上下文
4. **运行时支持**：在运行时创建对应的runCtx对象
5. **兼容性**：与现有编译和执行流程无缝集成

## 3. 核心概念

### 3.1 作用域Key命名规范

每个作用域都有一个唯一的key，遵循以下命名规范：

- **全局作用域**：`main`
- **主函数**：`main.main`
- **普通函数**：`main.FunctionName`
- **结构体方法**：`main.StructName.MethodName`
- **其他模块**：`moduleName.FunctionName`

### 3.2 编译时上下文(Compile Context)

在编译阶段，每个作用域都有对应的编译时上下文，包含：
- 唯一的key标识符
- 父级上下文的引用
- 当前作用域内的变量、函数、类型信息

### 3.3 运行时上下文(Runtime Context)

在运行时，每个作用域对应一个runCtx对象，包含：
- 唯一的key标识符
- 父级runCtx的引用
- 当前作用域内的局部变量

## 4. 实现细节

### 4.1 编译器修改

#### 4.1.1 作用域路径跟踪

编译器需要维护当前的作用域路径，用于生成唯一的key：

```go
type Compiler struct {
    // ... existing fields ...
    currentScopePath []string // 跟踪当前作用域路径，如 ["main", "main"] 表示主函数
}
```

#### 4.1.2 作用域Key生成

提供辅助函数生成当前作用域的唯一key：

```go
func (c *Compiler) getCurrentScopeKey() string {
    if len(c.currentScopePath) == 0 {
        return "main"
    }
    return strings.Join(c.currentScopePath, ".")
}
```

#### 4.1.3 BlockStmt编译

在编译BlockStmt时，需要：
1. 分析当前上下文的key
2. 生成进入作用域的指令
3. 编译块内的语句
4. 生成退出作用域的指令

### 4.2 虚拟机修改

#### 4.2.1 上下文映射

虚拟机维护一个key到RuntimeContext的映射：

```go
type VM struct {
    // ... existing fields ...
    ctxMap   map[string]*RuntimeContext // key到上下文的映射
    ctxStack []*RuntimeContext          // 上下文栈，用于作用域嵌套
    currentCtx *RuntimeContext          // 当前上下文
}
```

#### 4.2.2 RuntimeContext结构

```go
type RuntimeContext struct {
    Key    string           // 唯一标识符
    Parent *RuntimeContext  // 父级上下文引用
    Locals map[string]interface{} // 局部变量映射
}
```

#### 4.2.3 作用域管理指令

新增两个操作码用于管理基于key的作用域：

```go
const (
    // ... existing opcodes ...
    OpEnterScopeWithKey // 进入指定key的作用域
    OpExitScopeWithKey  // 退出指定key的作用域
)
```

### 4.3 指令处理

#### 4.3.1 OpEnterScopeWithKey

处理进入指定key作用域的指令：
1. 从指令参数获取作用域key
2. 将当前上下文压入栈
3. 查找或创建对应key的上下文
4. 设置为当前上下文

#### 4.3.2 OpExitScopeWithKey

处理退出指定key作用域的指令：
1. 从指令参数获取作用域key
2. 验证当前是否在正确的上下文中
3. 从栈中弹出父级上下文
4. 设置为当前上下文

## 5. 执行流程

### 5.1 编译阶段

1. 编译器在编译函数时，将函数名添加到作用域路径
2. 编译BlockStmt时，生成作用域管理指令
3. 为每个作用域生成唯一的key标识符

### 5.2 运行时阶段

1. 虚拟机执行OpEnterScopeWithKey指令时创建runCtx
2. 变量存储和查找都在当前runCtx中进行
3. 虚拟机执行OpExitScopeWithKey指令时回退到父级上下文

## 6. 示例

### 6.1 代码示例

```go
package main

func main() {
    x := 10  // 变量存储在key为"main.main"的runCtx中
    
    func() {
        y := 20  // 变量存储在匿名函数的runCtx中，父级为"main.main"
    }()
}
```

### 6.2 生成的指令序列

```
OpEnterScopeWithKey "main.main"  // 进入主函数作用域
OpLoadConst 10                   // 加载常量10
OpStoreName "x"                  // 存储变量x
OpEnterScopeWithKey "main.main.func1"  // 进入匿名函数作用域
OpLoadConst 20                   // 加载常量20
OpStoreName "y"                  // 存储变量y
OpExitScopeWithKey "main.main.func1"   // 退出匿名函数作用域
OpExitScopeWithKey "main.main"         // 退出主函数作用域
```

## 7. 优势

1. **精确的作用域管理**：每个作用域都有唯一标识，避免命名冲突
2. **高效的变量查找**：通过key直接定位上下文，减少查找时间
3. **清晰的层级关系**：父子上下文引用明确，便于变量查找和作用域管理
4. **模块化支持**：天然支持不同模块的隔离
5. **调试友好**：通过key可以清楚地知道当前所处的作用域

## 8. 注意事项

1. **性能考虑**：需要平衡上下文管理的开销和查找效率
2. **内存管理**：及时清理不再使用的上下文对象
3. **错误处理**：确保作用域切换的正确性，避免栈溢出或下溢
4. **兼容性**：确保新机制与现有代码兼容
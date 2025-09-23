# VM 重构方案

## 1. 当前架构分析

当前的 VM 实现有以下特点：
- 使用全局变量和局部变量映射来管理变量
- 通过 ExecutionContext 和 ScopeManager 管理作用域
- 编译器生成指令，VM 执行指令
- 作用域管理分散在多个组件中

## 2. 重构目标

根据架构师的要求，重构目标是：
1. 内部实现一个 run 接口，用于执行作用域里的指令
2. 进入子作用域（函数、代码块等），调用一层 run()
3. run 入口，基于当前的作用域，创建 ctx，ctx.parent 会指向上级的 ctx（基于 pathKey 找到上级 ctx）
4. 每个包都会有在 VM 里保留全局的 ctx
5. 每个 ctx 都有对应的 pathKey
6. 变量的赋值、取值，先从当前 ctx 里面查找，如果没有，则向上查找
7. 所有变量都需要预先创建，如果赋值时，找不到对应变量，则属于异常
8. 退出作用域，销毁对应 ctx

## 3. 详细重构方案

### 3.1 核心数据结构设计

#### Context 结构体：
```go
type Context struct {
    pathKey     string           // 作用域路径键
    parent      *Context         // 父级上下文
    variables   map[string]interface{} // 变量映射
    types       map[string]string      // 变量类型映射
    children    map[string]*Context    // 子上下文映射
}
```

#### 增强的 VM 结构体：
```go
type VM struct {
    stack          *Stack
    globalCtx      *Context        // 全局上下文
    currentCtx     *Context        // 当前上下文
    instructions   []*Instruction
    ip             int
    retval         interface{}
    functionRegistry map[string]func(args ...interface{}) (interface{}, error)
    scriptFunctions map[string]*ScriptFunction
    typeSystem     map[string]types.IType
    debug          bool
    executionCount int
    maxInstructions int64
    moduleManager  ModuleManagerInterface
    handlers       map[OpCode]OpHandler
}
```

### 3.2 核心接口设计

#### run 接口：
```go
// Run executes instructions within a specific context
func (vm *VM) Run(ctx *Context, startIP, endIP int) (interface{}, error)
```

#### 作用域管理接口：
```go
// EnterScope creates and enters a new scope
func (vm *VM) EnterScope(pathKey string) *Context

// ExitScope exits the current scope and returns to parent
func (vm *VM) ExitScope() *Context

// GetCurrentContext returns the current execution context
func (vm *VM) GetCurrentContext() *Context
```

#### 变量操作接口：
```go
// SetVariable sets a variable in the current context
func (vm *VM) SetVariable(name string, value interface{}) error

// GetVariable gets a variable, searching from current context up to root
func (vm *VM) GetVariable(name string) (interface{}, bool)

// MustGetVariable gets a variable, panics if not found
func (vm *VM) MustGetVariable(name string) interface{}
```

### 3.3 Context 实现细节

#### Context 创建和管理：
```go
// NewContext creates a new context
func NewContext(pathKey string, parent *Context) *Context {
    return &Context{
        pathKey:   pathKey,
        parent:    parent,
        variables: make(map[string]interface{}),
        types:     make(map[string]string),
        children:  make(map[string]*Context),
    }
}

// GetVariable searches for a variable in the context hierarchy
func (ctx *Context) GetVariable(name string) (interface{}, bool) {
    // First check current context
    if value, exists := ctx.variables[name]; exists {
        return value, true
    }
    
    // Then check parent context if exists
    if ctx.parent != nil {
        return ctx.parent.GetVariable(name)
    }
    
    return nil, false
}

// SetVariable sets a variable in the current context
func (ctx *Context) SetVariable(name string, value interface{}) {
    ctx.variables[name] = value
}

// MustGetVariable gets a variable, panics if not found
func (ctx *Context) MustGetVariable(name string) interface{} {
    if value, exists := ctx.variables[name]; exists {
        return value
    }
    
    if ctx.parent != nil {
        return ctx.parent.MustGetVariable(name)
    }
    
    panic(fmt.Sprintf("variable %s not found in context hierarchy", name))
}
```

### 3.4 指令处理重构

#### 作用域相关指令处理：
```go
// handleEnterScopeWithKey handles entering a scope with a specific key
func (vm *VM) handleEnterScopeWithKey(v *VM, instr *Instruction) error

// handleExitScopeWithKey handles exiting a scope with a specific key
func (vm *VM) handleExitScopeWithKey(v *VM, instr *Instruction) error
```

#### 变量操作指令处理：
```go
// handleLoadName searches for variable in context hierarchy
func (vm *VM) handleLoadName(v *VM, instr *Instruction) error

// handleStoreName stores variable in current context
func (vm *VM) handleStoreName(v *VM, instr *Instruction) error
```

### 3.5 VM 执行流程重构

#### 主执行函数：
```go
// Execute executes the instructions using the new context-based approach
func (vm *VM) Execute(ctx context.Context) (interface{}, error) {
    // Initialize with global context
    vm.currentCtx = vm.globalCtx
    vm.ip = 0
    vm.executionCount = 0

    // Run in the global context
    return vm.Run(vm.globalCtx, 0, len(vm.instructions))
}
```

#### Run 函数实现：
```go
// Run executes instructions within a specific context
func (vm *VM) Run(ctx *Context, startIP, endIP int) (interface{}, error) {
    // Save current execution state
    prevCtx := vm.currentCtx
    prevIP := vm.ip
    
    // Set new execution context
    vm.currentCtx = ctx
    vm.ip = startIP
    
    defer func() {
        // Restore execution state
        vm.currentCtx = prevCtx
        vm.ip = prevIP
    }()
    
    // Execute instructions in the given range
    for vm.ip < endIP && vm.ip < len(vm.instructions) {
        // Check instruction limit
        if vm.maxInstructions > 0 && int64(vm.executionCount) >= vm.maxInstructions {
            return nil, fmt.Errorf("maximum instruction limit exceeded: %d instructions executed", vm.executionCount)
        }
        
        // Handle context cancellation
        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        default:
        }
        
        instr := vm.instructions[vm.ip]
        
        // Debug output
        if vm.debug {
            fmt.Printf("IP: %d, Instruction: %s, Stack: %v\n", vm.ip, instr.String(), vm.stack.GetSlice())
        }
        
        // Use dispatch table
        if handler, exists := vm.handlers[instr.Op]; exists {
            if err := handler(vm, instr); err != nil {
                return nil, err
            }
        } else {
            return nil, fmt.Errorf("unknown opcode: %v at IP %d", instr.Op, vm.ip)
        }
        
        vm.ip++
        vm.executionCount++
    }
    
    return vm.retval, nil
}
```

### 3.6 函数调用重构

#### 脚本函数执行：
```go
// executeScriptFunction executes a script-defined function with context-based scopes
func (vm *VM) executeScriptFunction(scriptFunc *ScriptFunction, args ...interface{}) (interface{}, error) {
    // Check argument count
    if len(args) != scriptFunc.ParamCount {
        return nil, fmt.Errorf("function %s expects %d arguments, got %d", scriptFunc.Name, scriptFunc.ParamCount, len(args))
    }

    // Create new context for function execution
    funcCtx := vm.EnterScope(fmt.Sprintf("function.%s", scriptFunc.Name))
    defer vm.ExitScope()

    // Set up function parameters as local variables in the new context
    for i, paramName := range scriptFunc.ParamNames {
        if i < len(args) {
            funcCtx.SetVariable(paramName, args[i])
        }
    }

    // Execute function instructions in the new context
    return vm.Run(funcCtx, scriptFunc.StartIP, scriptFunc.EndIP)
}
```

### 3.7 编译器集成

#### 编译器适配：
- 编译器需要生成作用域管理指令（OpEnterScopeWithKey, OpExitScopeWithKey）
- 编译器需要生成变量创建指令（OpCreateVar）
- 编译器需要维护作用域路径信息

#### 作用域路径管理：
- 使用点分隔的路径作为作用域键（如 "main.function.loop"）
- 父级作用域可以通过路径推导（去掉最后一个部分）

## 4. 实施步骤

1. **第一阶段：核心数据结构实现**
   - 实现 Context 结构体和相关方法
   - 修改 VM 结构体，添加上下文管理字段
   - 实现基本的上下文操作接口

2. **第二阶段：指令处理重构**
   - 重构现有的指令处理函数，使用新的上下文机制
   - 实现作用域管理指令处理
   - 实现变量操作指令处理

3. **第三阶段：执行流程重构**
   - 实现 Run 接口
   - 重构 Execute 函数
   - 重构函数调用机制

4. **第四阶段：编译器集成**
   - 修改编译器生成作用域管理指令
   - 确保编译器与新的 VM 上下文机制兼容

5. **第五阶段：测试和验证**
   - 运行现有测试确保功能正常
   - 添加新的测试用例验证作用域管理
   - 性能测试确保重构没有引入性能问题

## 5. 预期收益

1. **更清晰的作用域管理**：通过上下文层次结构明确管理变量作用域
2. **更好的可维护性**：将作用域管理逻辑集中到 VM 中
3. **更强的一致性**：所有变量操作都遵循相同的作用域查找规则
4. **更好的错误处理**：变量未声明时能及时发现并报错
5. **更灵活的执行模型**：支持在不同作用域中执行指令
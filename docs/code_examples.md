# GoScript改进代码示例

## 1. VM操作码优化示例

### 当前实现问题
```go
// 当前VM实现 - 大量重复代码
func (vm *VM) Execute(ctx context.Context) (interface{}, error) {
    for vm.ip < len(vm.instructions) {
        instr := vm.instructions[vm.ip]
        
        switch instr.Op {
        case OpAdd:
            if len(vm.stack) < 2 {
                return nil, fmt.Errorf("stack underflow")
            }
            right := vm.Pop()
            left := vm.Pop()
            result, err := vm.add(left, right)
            if err != nil {
                return nil, err
            }
            vm.Push(result)
        case OpSub:
            if len(vm.stack) < 2 {
                return nil, fmt.Errorf("stack underflow")
            }
            right := vm.Pop()
            left := vm.Pop()
            result, err := vm.sub(left, right)
            if err != nil {
                return nil, err
            }
            vm.Push(result)
        // ... 大量重复代码
        }
        vm.ip++
    }
}
```

### 改进后的实现
```go
// 改进后的VM实现 - 使用操作码分发表
type OpHandler func(vm *VM, instr *Instruction) error

type VM struct {
    stack      []Value
    stackTop   int
    stackSize  int
    code       []Instruction
    ip         int
    handlers   map[OpCode]OpHandler
    valuePool  *ValuePool
}

func NewVM() *VM {
    vm := &VM{
        stack:     make([]Value, 1024),
        stackSize: 1024,
        handlers:  make(map[OpCode]OpHandler),
        valuePool: NewValuePool(),
    }
    
    // 注册操作码处理器
    vm.registerHandlers()
    return vm
}

func (vm *VM) registerHandlers() {
    vm.handlers[OpLoadConst] = vm.handleLoadConst
    vm.handlers[OpLoadName] = vm.handleLoadName
    vm.handlers[OpStoreName] = vm.handleStoreName
    vm.handlers[OpAdd] = vm.handleAdd
    vm.handlers[OpSub] = vm.handleSub
    vm.handlers[OpMul] = vm.handleMul
    vm.handlers[OpDiv] = vm.handleDiv
    // ... 其他操作码
}

func (vm *VM) Execute(ctx context.Context) (Value, error) {
    for vm.ip < len(vm.code) {
        instr := vm.code[vm.ip]
        
        if handler, exists := vm.handlers[instr.Op]; exists {
            if err := handler(vm, instr); err != nil {
                return nil, err
            }
        } else {
            return nil, fmt.Errorf("unknown opcode: %v", instr.Op)
        }
        
        vm.ip++
    }
    
    if vm.stackTop > 0 {
        return vm.stack[vm.stackTop-1], nil
    }
    return nil, nil
}

// 具体的操作码处理器
func (vm *VM) handleAdd(vm *VM, instr *Instruction) error {
    if vm.stackTop < 2 {
        return &VmError{
            Type:    ErrorTypeStackUnderflow,
            Message: "expected 2 values for addition",
            IP:      vm.ip,
        }
    }
    
    right := vm.pop()
    left := vm.pop()
    result := left.Add(right)
    vm.push(result)
    return nil
}

func (vm *VM) handleSub(vm *VM, instr *Instruction) error {
    if vm.stackTop < 2 {
        return &VmError{
            Type:    ErrorTypeStackUnderflow,
            Message: "expected 2 values for subtraction",
            IP:      vm.ip,
        }
    }
    
    right := vm.pop()
    left := vm.pop()
    result := left.Sub(right)
    vm.push(result)
    return nil
}
```

## 2. 值类型系统示例

### 当前实现问题
```go
// 当前实现 - 使用interface{}，类型检查在运行时
type VM struct {
    stack []interface{}
}

func (vm *VM) add(left, right interface{}) (interface{}, error) {
    switch l := left.(type) {
    case int:
        switch r := right.(type) {
        case int:
            return l + r, nil
        case float64:
            return float64(l) + r, nil
        }
    case float64:
        switch r := right.(type) {
        case int:
            return l + float64(r), nil
        case float64:
            return l + r, nil
        }
    }
    return nil, fmt.Errorf("unsupported types")
}
```

### 改进后的实现
```go
// 改进后的实现 - 使用值类型系统
type Value interface {
    Type() ValueType
    Int() int64
    Float() float64
    String() string
    Bool() bool
    Add(other Value) Value
    Sub(other Value) Value
    Mul(other Value) Value
    Div(other Value) Value
    Clone() Value
    Reset()
}

type ValueType int

const (
    TypeInt ValueType = iota
    TypeFloat
    TypeString
    TypeBool
    TypeNil
)

type IntValue struct {
    value int64
}

func (v *IntValue) Type() ValueType { return TypeInt }
func (v *IntValue) Int() int64 { return v.value }
func (v *IntValue) Float() float64 { return float64(v.value) }
func (v *IntValue) String() string { return strconv.FormatInt(v.value, 10) }
func (v *IntValue) Bool() bool { return v.value != 0 }

func (v *IntValue) Add(other Value) Value {
    switch other.Type() {
    case TypeInt:
        return &IntValue{value: v.value + other.Int()}
    case TypeFloat:
        return &FloatValue{value: float64(v.value) + other.Float()}
    default:
        panic("unsupported type for addition")
    }
}

func (v *IntValue) Sub(other Value) Value {
    switch other.Type() {
    case TypeInt:
        return &IntValue{value: v.value - other.Int()}
    case TypeFloat:
        return &FloatValue{value: float64(v.value) - other.Float()}
    default:
        panic("unsupported type for subtraction")
    }
}

func (v *IntValue) Clone() Value {
    return &IntValue{value: v.value}
}

func (v *IntValue) Reset() {
    v.value = 0
}

// 值池实现
type ValuePool struct {
    intPool    sync.Pool
    floatPool  sync.Pool
    stringPool sync.Pool
    boolPool   sync.Pool
}

func NewValuePool() *ValuePool {
    return &ValuePool{
        intPool: sync.Pool{
            New: func() interface{} { return &IntValue{} },
        },
        floatPool: sync.Pool{
            New: func() interface{} { return &FloatValue{} },
        },
        stringPool: sync.Pool{
            New: func() interface{} { return &StringValue{} },
        },
        boolPool: sync.Pool{
            New: func() interface{} { return &BoolValue{} },
        },
    }
}

func (vp *ValuePool) GetInt() *IntValue {
    return vp.intPool.Get().(*IntValue)
}

func (vp *ValuePool) PutInt(v *IntValue) {
    v.Reset()
    vp.intPool.Put(v)
}
```

## 3. 类型检查器示例

### 当前实现问题
```go
// 当前实现 - 缺乏类型检查
func (c *Compiler) compileBinaryExpr(expr *ast.BinaryExpr) error {
    // 编译左操作数
    err := c.compileExpr(expr.X)
    if err != nil {
        return err
    }
    
    // 编译右操作数
    err = c.compileExpr(expr.Y)
    if err != nil {
        return err
    }
    
    // 直接生成操作码，没有类型检查
    switch expr.Op {
    case token.ADD:
        c.emitInstruction(vm.NewInstruction(vm.OpAdd, nil, nil))
    case token.SUB:
        c.emitInstruction(vm.NewInstruction(vm.OpSub, nil, nil))
    }
    
    return nil
}
```

### 改进后的实现
```go
// 改进后的实现 - 带类型检查的编译器
type TypeChecker struct {
    symbols map[string]Type
    errors  []TypeError
}

type Type interface {
    Name() string
    Size() int
    AssignableTo(other Type) bool
    ConvertibleTo(other Type) bool
}

type IntType struct{}

func (t *IntType) Name() string { return "int" }
func (t *IntType) Size() int { return 8 }
func (t *IntType) AssignableTo(other Type) bool {
    return other.Name() == "int" || other.Name() == "float64"
}
func (t *IntType) ConvertibleTo(other Type) bool {
    return other.Name() == "float64" || other.Name() == "string"
}

type Compiler struct {
    typeChecker *TypeChecker
    vm          *VM
    ip          int
}

func (c *Compiler) compileBinaryExpr(expr *ast.BinaryExpr) error {
    // 1. 类型检查
    leftType, err := c.typeChecker.CheckExpr(expr.X)
    if err != nil {
        return err
    }
    
    rightType, err := c.typeChecker.CheckExpr(expr.Y)
    if err != nil {
        return err
    }
    
    // 2. 类型兼容性检查
    if !c.typeChecker.IsCompatible(leftType, rightType, expr.Op) {
        return &TypeError{
            Message: fmt.Sprintf("incompatible types: %s and %s", leftType.Name(), rightType.Name()),
            Line:    expr.Pos().Line,
            Column:  expr.Pos().Column,
        }
    }
    
    // 3. 编译左操作数
    err = c.compileExpr(expr.X)
    if err != nil {
        return err
    }
    
    // 4. 编译右操作数
    err = c.compileExpr(expr.Y)
    if err != nil {
        return err
    }
    
    // 5. 生成类型特化的操作码
    resultType := c.typeChecker.GetResultType(leftType, rightType, expr.Op)
    opCode := c.getTypedOpCode(expr.Op, resultType)
    c.emitInstruction(vm.NewInstruction(opCode, nil, nil))
    
    return nil
}

func (c *Compiler) getTypedOpCode(op token.Token, resultType Type) OpCode {
    switch op {
    case token.ADD:
        switch resultType.Name() {
        case "int":
            return OpAddInt
        case "float64":
            return OpAddFloat
        case "string":
            return OpAddString
        }
    case token.SUB:
        switch resultType.Name() {
        case "int":
            return OpSubInt
        case "float64":
            return OpSubFloat
        }
    }
    return OpNop
}
```

## 4. JIT编译器示例

### 当前实现问题
```go
// 当前实现 - 解释执行，性能差
func (vm *VM) executeAdd(left, right Value) Value {
    // 每次都要进行类型检查和转换
    switch l := left.(type) {
    case int:
        switch r := right.(type) {
        case int:
            return l + r
        case float64:
            return float64(l) + r
        }
    case float64:
        switch r := right.(type) {
        case int:
            return l + float64(r)
        case float64:
            return l + r
        }
    }
    return nil
}
```

### 改进后的实现
```go
// 改进后的实现 - JIT编译
type JITCompiler struct {
    targetArch string
    codeGen    CodeGenerator
    optimizer  Optimizer
}

type CodeGenerator interface {
    Generate(ir *IR) []byte
}

type IR struct {
    Instructions []IRInstruction
    Constants    []Value
    Functions    []IRFunction
}

type IRInstruction struct {
    Op    IROpCode
    Args  []int
    Type  ValueType
}

type IROpCode int

const (
    IROpLoadConst IROpCode = iota
    IROpAdd
    IROpSub
    IROpMul
    IROpDiv
    IROpCall
    IROpReturn
)

func (jit *JITCompiler) Compile(ir *IR) ([]byte, error) {
    // 1. 优化IR
    optimized := jit.optimizer.Optimize(ir)
    
    // 2. 生成机器码
    machineCode := jit.codeGen.Generate(optimized)
    
    // 3. 验证和修复
    if err := jit.validateCode(machineCode); err != nil {
        return nil, err
    }
    
    return machineCode, nil
}

// 类型特化的代码生成器
type TypedCodeGenerator struct {
    arch string
}

func (tcg *TypedCodeGenerator) Generate(ir *IR) []byte {
    var code []byte
    
    for _, instr := range ir.Instructions {
        switch instr.Op {
        case IROpAdd:
            switch instr.Type {
            case TypeInt:
                code = append(code, tcg.generateIntAdd()...)
            case TypeFloat:
                code = append(code, tcg.generateFloatAdd()...)
            }
        case IROpSub:
            switch instr.Type {
            case TypeInt:
                code = append(code, tcg.generateIntSub()...)
            case TypeFloat:
                code = append(code, tcg.generateFloatSub()...)
            }
        }
    }
    
    return code
}

func (tcg *TypedCodeGenerator) generateIntAdd() []byte {
    // 生成整数加法的机器码
    // 这里简化处理，实际实现需要根据目标架构生成
    return []byte{
        0x48, 0x01, 0xf8, // add %rdi, %rax
        0xc3,             // ret
    }
}

func (tcg *TypedCodeGenerator) generateFloatAdd() []byte {
    // 生成浮点数加法的机器码
    return []byte{
        0xf2, 0x0f, 0x58, 0xc1, // addsd %xmm1, %xmm0
        0xc3,                   // ret
    }
}
```

## 5. 并发执行示例

### 当前实现问题
```go
// 当前实现 - 单线程执行
func (vm *VM) Execute(ctx context.Context) (interface{}, error) {
    for vm.ip < len(vm.instructions) {
        // 单线程执行，无法利用多核
        instr := vm.instructions[vm.ip]
        // ... 执行指令
        vm.ip++
    }
}
```

### 改进后的实现
```go
// 改进后的实现 - 支持并发执行
type ConcurrentVM struct {
    coroutines []*Coroutine
    scheduler  *Scheduler
    channels   map[string]*Channel
}

type Coroutine struct {
    id       int
    vm       *VM
    state    CoroutineState
    channel  chan Value
    parent   *Coroutine
    children []*Coroutine
}

type Scheduler struct {
    readyQueue []*Coroutine
    running    *Coroutine
    mutex      sync.Mutex
}

func (vm *ConcurrentVM) SpawnCoroutine(code []Instruction) *Coroutine {
    coroutine := &Coroutine{
        id:    vm.nextCoroutineID(),
        vm:    vm.clone(),
        state: StateReady,
    }
    
    vm.coroutines = append(vm.coroutines, coroutine)
    vm.scheduler.Schedule(coroutine)
    
    return coroutine
}

func (vm *ConcurrentVM) Execute(ctx context.Context) (Value, error) {
    for {
        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        default:
            if vm.scheduler.HasReady() {
                coroutine := vm.scheduler.GetNext()
                if coroutine != nil {
                    vm.scheduler.Run(coroutine)
                }
            } else {
                break
            }
        }
    }
    
    return vm.getResult(), nil
}

// 通道实现
type Channel struct {
    buffer    []Value
    capacity  int
    senders   []*Coroutine
    receivers []*Coroutine
    mutex     sync.Mutex
}

func (ch *Channel) Send(value Value) error {
    ch.mutex.Lock()
    defer ch.mutex.Unlock()
    
    if len(ch.buffer) >= ch.capacity {
        return ErrChannelFull
    }
    
    ch.buffer = append(ch.buffer, value)
    
    // 唤醒等待的接收者
    if len(ch.receivers) > 0 {
        receiver := ch.receivers[0]
        ch.receivers = ch.receivers[1:]
        receiver.Resume()
    }
    
    return nil
}

func (ch *Channel) Receive() (Value, error) {
    ch.mutex.Lock()
    defer ch.mutex.Unlock()
    
    if len(ch.buffer) > 0 {
        value := ch.buffer[0]
        ch.buffer = ch.buffer[1:]
        return value, nil
    }
    
    // 没有数据，等待发送者
    return nil, ErrChannelEmpty
}
```

## 6. 性能测试示例

```go
// 性能测试和基准测试
func BenchmarkVMExecution(b *testing.B) {
    vm := NewVM()
    code := generateTestCode(1000) // 生成1000条指令的测试代码
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := vm.Execute(context.Background())
        if err != nil {
            b.Fatal(err)
        }
    }
}

func BenchmarkJITExecution(b *testing.B) {
    vm := NewJITVM()
    code := generateTestCode(1000)
    
    // 预热JIT编译器
    for i := 0; i < 100; i++ {
        vm.Execute(context.Background())
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := vm.Execute(context.Background())
        if err != nil {
            b.Fatal(err)
        }
    }
}

func TestPerformanceImprovement(t *testing.T) {
    // 测试性能改进效果
    vm := NewVM()
    jitVM := NewJITVM()
    
    code := generateComplexCode()
    
    // 测试解释执行
    start := time.Now()
    _, err := vm.Execute(context.Background())
    if err != nil {
        t.Fatal(err)
    }
    interpretTime := time.Since(start)
    
    // 测试JIT执行
    start = time.Now()
    _, err = jitVM.Execute(context.Background())
    if err != nil {
        t.Fatal(err)
    }
    jitTime := time.Since(start)
    
    improvement := float64(interpretTime) / float64(jitTime)
    t.Logf("JIT improvement: %.2fx", improvement)
    
    if improvement < 2.0 {
        t.Errorf("Expected at least 2x improvement, got %.2fx", improvement)
    }
}
```

这些代码示例展示了如何具体实施改进方案，从VM优化到JIT编译，从类型检查到并发执行。每个示例都提供了当前实现的问题和改进后的解决方案，可以作为实际开发的参考。

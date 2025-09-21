# GoScript 方法调用统一化改造方案

## 目标

改造 GoScript 的函数调用机制，使结构体方法调用与普通函数调用统一：

1. `struct.Function`，比如有 struct Rectangle，对应的方法为 Area，则自动注册方法 `Rectangle.Area(r Rectangle)`
2. 如果是可以修改 struct 的，指针传递的，则注册的是 `Rectangle.Set(r *Rectangle, value int)`
3. 自动将 struct 对象转成函数参数的第一个参数
4. 这样的方式，就能够统一普通函数调用和 struct 的方法调用

## 当前实现分析

### 方法注册

当前实现中，方法在编译时被注册为带有接收者类型前缀的函数名：
- 值接收者方法：`Rectangle.Area`
- 指针接收者方法：`Rectangle.SetHeight`

### 方法调用

在方法调用时，编译器会：
1. 将接收者对象压入栈中
2. 将方法参数压入栈中
3. 发出 OpCall 指令，参数数量为参数个数+1（包括接收者）

### 虚拟机执行

在虚拟机执行 OpCall 指令时：
1. 从栈中弹出参数（包括接收者）
2. 调用 executeScriptFunction 执行函数
3. 在 executeScriptFunction 中处理接收者参数的特殊逻辑

## 改造方案

### 1. 编译器改造

#### 方法注册
修改 [compileMethod](../compiler/compiler.go#L377-L452) 函数，使方法注册时符合新的规范：

```go
// 值接收者方法 Rectangle.Area() 注册为 Rectangle.Area(r Rectangle)
// 指针接收者方法 Rectangle.SetHeight() 注册为 Rectangle.SetHeight(r *Rectangle)
```

#### 方法调用编译
修改 [compileCallExpr](../compiler/compiler.go#L1177-L1223) 函数，确保方法调用时正确处理参数顺序：

```go
// 对于方法调用 obj.Method(arg1, arg2)
// 编译后栈中顺序应为: [obj, arg1, arg2]
// OpCall 指令参数数量为: 3
```

### 2. 虚拟机改造

#### 函数执行
修改 [executeScriptFunction](../vm/opcodes.go#L811-L997) 函数，确保参数正确处理：

```go
// 对于值接收者方法，需要复制接收者对象
// 对于指针接收者方法，直接使用接收者对象
```

## 实施步骤

1. 修改编译器中的方法注册逻辑
2. 修改编译器中的方法调用编译逻辑
3. 修改虚拟机中的函数执行逻辑
4. 测试验证改造效果

## 预期效果

改造后，以下两种调用方式应该等价：

```go
// 方法调用
rect := Rectangle{width: 10, height: 5}
area := rect.Area()

// 等价的函数调用
rect := Rectangle{width: 10, height: 5}
area := Rectangle.Area(rect)
```

这样就可以统一普通函数调用和结构体方法调用的处理逻辑，简化虚拟机实现。
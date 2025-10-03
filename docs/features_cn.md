# GoScript 功能特性详解

## 1. 概述

GoScript是一个兼容Go标准语法的脚本引擎，支持大部分Go语言的核心特性。本文档详细说明GoScript当前支持的功能和特性。

## 2. 支持的语法特性

### 2.1 基本语法结构

#### 包声明
```go
package main
```

#### 导入声明
```go
import "fmt"
import (
    "strings"
    "math"
)
```

#### 变量声明
```go
// var声明
var x int
var y, z int = 10, 20
var (
    a = 1
    b = 2
)

// 短变量声明
x := 10
name := "GoScript"
```

#### 常量声明
```go
const pi = 3.14159
const (
    a = 1
    b = 2
)
```

### 2.2 数据类型

#### 基本类型
- 整数类型：int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64
- 浮点类型：float32, float64
- 布尔类型：bool
- 字符串类型：string

#### 复合类型
- 数组：[n]T
- 切片：[]T
- 结构体：struct
- 接口：interface{}

### 2.3 控制结构

#### 条件语句
```go
if x > 0 {
    // do something
} else if x < 0 {
    // do something else
} else {
    // do another thing
}
```

#### 循环语句
```go
// 传统for循环
for i := 0; i < 10; i++ {
    // do something
}

// while循环形式
for x > 0 {
    x--
}

// 无限循环
for {
    // do something
    break
}
```

#### Range语句
```go
// 遍历切片
slice := []int{1, 2, 3}
for index, value := range slice {
    // index是索引，value是值
}

// 只遍历值
for _, value := range slice {
    // 只关心值
}

// 只遍历索引
for index := range slice {
    // 只关心索引
}

// 遍历字符串
str := "hello"
for index, char := range str {
    // char是rune类型
}
```

### 2.4 函数

#### 函数声明
```go
func add(a, b int) int {
    return a + b
}

func greet(name string) string {
    return "Hello, " + name
}

// 多返回值
func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}
```

#### 函数调用
```go
result := add(1, 2)
greeting := greet("World")
```

### 2.5 结构体和方法

#### 结构体定义
```go
type Person struct {
    Name string
    Age  int
}

type Rectangle struct {
    Width, Height float64
}
```

#### 结构体初始化
```go
// 字段名初始化
person := Person{Name: "Alice", Age: 30}

// 顺序初始化
person := Person{"Bob", 25}

// 复合字面量
rect := Rectangle{Width: 10, Height: 5}
```

#### 方法定义
```go
// 值接收者方法
func (p Person) GetName() string {
    return p.Name
}

// 指针接收者方法
func (p *Person) SetAge(age int) {
    p.Age = age
}

// 带返回值的方法
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}
```

#### 方法调用
```go
person := Person{Name: "Alice", Age: 30}
name := person.GetName()
person.SetAge(31)
```

### 2.6 操作符

#### 算术操作符
- 加法：+
- 减法：-
- 乘法：*
- 除法：/
- 取模：%
- 自增：++
- 自减：--

#### 比较操作符
- 等于：==
- 不等于：!=
- 小于：< 
- 小于等于：<=
- 大于：>
- 大于等于：>=

#### 逻辑操作符
- 逻辑与：&&
- 逻辑或：||
- 逻辑非：!

#### 赋值操作符
- 简单赋值：=
- 加赋值：+=
- 减赋值：-=
- 乘赋值：*=
- 除赋值：/=
- 模赋值：%=

## 3. 内置函数

### 3.1 基本内置函数
- len()：获取字符串、数组、切片、映射的长度
- int()：将值转换为整数
- float64()：将值转换为浮点数
- string()：将值转换为字符串

## 4. 模块系统

### 4.1 内置模块
GoScript提供以下内置模块：
- math：数学函数
- strings：字符串操作
- fmt：格式化输入输出
- json：JSON序列化和反序列化

### 4.2 模块使用
```go
// 使用模块函数
result := math.Abs(-5.0)
upper := strings.ToUpper("hello")
```

## 5. 错误处理

### 5.1 错误返回
```go
func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}
```

### 5.2 错误检查
```go
result, err := divide(10, 0)
if err != nil {
    // 处理错误
    return
}
```

## 6. 限制和不支持的特性

### 6.1 不支持的语法特性
- goroutine和channel
- unsafe包
- 反射(reflect包)
- 完整的包管理系统
- 类型断言
- 接口的具体实现
- defer语句
- switch语句
- select语句

### 6.2 类型系统限制
- 不支持泛型
- 不支持类型别名的复杂用法
- 不支持结构体标签(tag)

## 7. 性能特性

### 7.1 编译执行
GoScript将源代码编译为字节码，然后在虚拟机中执行，相比纯解释执行有更好的性能。

### 7.2 作用域优化
基于Key的上下文管理机制提供了高效的作用域查找和变量管理。

### 7.3 内存管理
通过对象池和预分配机制减少内存分配和GC压力。

## 8. 安全特性

### 8.1 资源限制
- 最大执行时间限制
- 最大内存使用限制
- 最大指令数限制

### 8.2 沙箱环境
- 禁止危险系统调用
- 限制文件系统访问
- 限制网络访问

### 8.3 模块访问控制
- 可配置的模块访问权限
- 禁止关键字列表
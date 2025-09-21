# 当前问题分析

## 问题描述

在测试用例 `TestScriptsInDataFolder` 中，`struct2.gs` 脚本执行结果不正确：
- 期望返回：60
- 实际返回：50

## 脚本分析

```go
type Rectangle struct {
	width  int
	height int
}

func (r Rectangle) Area() int {
	return r.width * r.height
}

// 没带指针的方法，修改应该是无效的
func (r Rectangle)SetWidth(width int) {
	r.width = width
}

func (r *Rectangle)SetHeight(height int) {
	r.height = height
}

func main() {
	// Create shapes
	rect := Rectangle{width: 10, height: 5}

	// 修改应该是无效的
	rect.SetWidth(6)
	// 修改应该是有效的
	rect.SetHeight(6)
	
	// Calculate areas
	rectArea := rect.Area()
	
	return rectArea
}
```

按照 Go 语言的语义：
1. `rect.SetWidth(6)` 应该是无效的，因为是值接收者，修改的是副本
2. `rect.SetHeight(6)` 应该是有效的，因为是指针接收者，修改的是原对象
3. 最终 `rect.Area()` 应该返回 `10 * 6 = 60`

但实际返回 50，说明 `SetHeight` 方法没有正确修改原始对象。

## 问题定位

通过调试输出发现，在执行 `Rectangle.SetWidth` 方法时，栈的内容顺序不正确：

```
Function IP: 44, About to execute SET_FIELD, Stack size: 3, Stack: [6 map[height:5 width:10] width]
Function IP: 44, Object is not a map, type: int, value: 6, pushing back unchanged
```

这表明在执行 `SET_FIELD` 操作时，对象和值在栈中的位置颠倒了。

## 根本原因

问题出在编译器生成 `SET_FIELD` 指令时栈中元素的顺序不正确。在函数内部执行 `SET_FIELD` 时，正确的栈顺序应该是：
- 栈顶：字段名
- 栈顶-1：值
- 栈顶-2：对象

但在实际执行中，顺序变成了：
- 栈顶：字段名
- 栈顶-1：对象
- 栈顶-2：值

## 解决方案

需要修改编译器中的 `compileCompositeLit` 函数，通过添加 `ROTATE` 指令来调整栈中元素的顺序，使结构体字面量初始化能够正确工作。

同时需要检查方法调用时参数的处理逻辑，确保在 `executeScriptFunction` 方法中正确处理接收者参数。
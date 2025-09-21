# GoScript 测试说明

## 测试结构

```
test/
├── script_test.go     # 测试主文件
├── data/              # 测试脚本数据
│   ├── hello.gs       # Hello World 脚本
│   ├── add.gs         # 加法运算脚本
│   ├── loop.gs        # 循环脚本
│   ├── conditional.gs # 条件语句脚本
│   └── function_call.gs # 函数调用脚本
```

## 运行测试

```bash
# 运行所有测试
go test ./test -v

# 运行特定测试
go test ./test -run TestScriptsInDataFolder -v
```

## 测试内容

### TestScriptsInDataFolder
自动加载 `test/data` 目录下的所有 `.gs` 脚本文件并执行验证。

### TestScriptWithCustomFunction
测试自定义函数的注册和调用。

### TestScriptWithVariables
测试变量的定义和使用。

## 脚本文件说明

### hello.gs
简单的 Hello World 脚本，返回字符串 "Hello, World!"。

### add.gs
包含加法函数的脚本，测试函数定义和调用。

### loop.gs
包含循环的脚本，计算 1 到 10 的和。

### conditional.gs
包含条件语句的脚本，返回两个数中的最大值。

### function_call.gs
包含嵌套函数调用的脚本，测试函数间的调用关系。

## 注意事项

1. 所有脚本文件必须使用 `.gs` 作为后缀
2. 脚本文件应包含 `main()` 函数作为入口点
3. 测试会自动验证脚本的执行结果
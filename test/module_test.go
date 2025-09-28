package test

import (
	"context"
	"testing"

	goscript "github.com/lengzhao/goscript"
)

func TestModuleImport(t *testing.T) {
	// 测试模块导入功能
	scriptSource := `
package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
	return 42
}
`

	// 创建脚本
	script := goscript.NewScript([]byte(scriptSource))

	// 设置指令数限制
	script.SetMaxInstructions(1000)

	// 导入模块
	err := script.ImportModule("fmt")
	if err != nil {
		t.Fatalf("Failed to import module: %v", err)
	}

	// 执行脚本
	result, err := script.RunContext(context.Background())
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	// 验证结果
	if result != 42 {
		t.Errorf("Expected 42, but got %v", result)
	}
}

func TestMultipleModuleImports(t *testing.T) {
	// 测试多个模块导入
	scriptSource := `
package main

import (
	"fmt"
	"math"
)

func main() {
	result := math.Abs(5.5)
	fmt.Printf("Absolute value: %f\n", result)
	return int(result)
}
`

	// 创建脚本
	script := goscript.NewScript([]byte(scriptSource))

	// 设置指令数限制
	script.SetMaxInstructions(1000)

	// 导入多个模块
	err := script.ImportModule("fmt", "math")
	if err != nil {
		t.Fatalf("Failed to import modules: %v", err)
	}

	// 执行脚本
	result, err := script.RunContext(context.Background())
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	// 验证结果
	if result != 5 {
		t.Errorf("Expected 5, but got %v", result)
	}
}

func TestModuleFunctionalityInScript(t *testing.T) {
	// Test using module functions through the script API
	script := goscript.NewScript([]byte(""))

	// Import modules explicitly
	err := script.ImportModule("strings", "json")
	if err != nil {
		t.Fatalf("Failed to import modules: %v", err)
	}

	// Test calling strings module functions through the script
	result1, err := script.CallFunction("strings.ToLower", "HELLO WORLD")
	if err != nil {
		t.Fatalf("Failed to call strings.ToLower function: %v", err)
	}
	if result1 != "hello world" {
		t.Errorf("Expected strings.ToLower('HELLO WORLD') to be 'hello world', got %v", result1)
	}

	result2, err := script.CallFunction("strings.Contains", "hello world", "world")
	if err != nil {
		t.Fatalf("Failed to call strings.Contains function: %v", err)
	}
	if result2 != true {
		t.Errorf("Expected strings.Contains('hello world', 'world') to be true, got %v", result2)
	}

	t.Log("Module functionality in script test passed")
}

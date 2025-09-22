package test

import (
	"context"
	"strings"
	"testing"
	"time"

	goscript "github.com/lengzhao/goscript"
	execContext "github.com/lengzhao/goscript/context"
)

func TestInstructionLimit(t *testing.T) {
	// 测试无限循环脚本是否会因为指令数限制而被终止
	scriptSource := `
package test

func main() {
	i := 0
	for {
		i = i + 1
	}
	return i
}
`

	// 创建脚本
	script := goscript.NewScript([]byte(scriptSource))

	// 设置较小的指令数限制用于测试
	securityCtx := &execContext.SecurityContext{
		MaxExecutionTime:  5 * time.Second,
		MaxMemoryUsage:    10 * 1024 * 1024, // 10MB
		AllowedModules:    []string{"fmt", "math"},
		ForbiddenKeywords: []string{"unsafe"},
		AllowCrossModule:  true,
		MaxInstructions:   1000, // 设置较小的指令数限制
	}
	script.SetSecurityContext((*execContext.SecurityContext)(securityCtx))

	// 执行脚本
	_, err := script.RunContext(context.Background())

	// 验证是否因为指令数限制而返回错误
	if err == nil {
		t.Error("Expected error due to instruction limit, but got nil")
		return
	}

	// 验证错误信息是否包含指令数限制相关的内容
	if err.Error() != "maximum instruction limit exceeded: 1000 instructions executed" {
		t.Errorf("Expected instruction limit error, but got: %v", err)
	}

	// 验证执行统计信息
	stats := script.GetExecutionStats()
	if stats.InstructionCount < 1000 {
		t.Errorf("Expected at least 1000 instructions executed, but got: %d", stats.InstructionCount)
	}
}

func TestInstructionLimitInFunction(t *testing.T) {
	// 测试在函数中执行无限循环是否会因为指令数限制而被终止
	scriptSource := `
package test

func infiniteLoop() {
	i := 0
	for {
		i = i + 1
	}
}

func main() {
	infiniteLoop()
	return 0
}
`

	// 创建脚本
	script := goscript.NewScript([]byte(scriptSource))

	// 设置较小的指令数限制用于测试
	securityCtx := &execContext.SecurityContext{
		MaxExecutionTime:  5 * time.Second,
		MaxMemoryUsage:    10 * 1024 * 1024, // 10MB
		AllowedModules:    []string{"fmt", "math"},
		ForbiddenKeywords: []string{"unsafe"},
		AllowCrossModule:  true,
		MaxInstructions:   1000, // 设置较小的指令数限制
	}
	script.SetSecurityContext((*execContext.SecurityContext)(securityCtx))

	// 执行脚本
	_, err := script.RunContext(context.Background())

	// 验证是否因为指令数限制而返回错误
	if err == nil {
		t.Error("Expected error due to instruction limit, but got nil")
		return
	}

	// 验证错误信息是否包含指令数限制相关的内容
	if !strings.Contains(err.Error(), "maximum instruction limit exceeded") {
		t.Errorf("Expected instruction limit error, but got: %v", err)
	}

	// 验证执行统计信息
	stats := script.GetExecutionStats()
	if stats.InstructionCount < 1 {
		t.Errorf("Expected at least 1 instructions executed, but got: %d", stats.InstructionCount)
	}
}

func TestNormalExecutionWithinLimit(t *testing.T) {
	// 测试正常脚本在指令数限制内能够正常执行
	scriptSource := `
package test

func main() {
	return 42
}
`

	// 创建脚本
	script := goscript.NewScript([]byte(scriptSource))

	// 设置较大的指令数限制
	securityCtx := &execContext.SecurityContext{
		MaxExecutionTime:  5 * time.Second,
		MaxMemoryUsage:    10 * 1024 * 1024, // 10MB
		AllowedModules:    []string{"fmt", "math"},
		ForbiddenKeywords: []string{"unsafe"},
		AllowCrossModule:  true,
		MaxInstructions:   1000, // 设置较小的指令数限制用于简单脚本
	}
	script.SetSecurityContext((*execContext.SecurityContext)(securityCtx))

	// 执行脚本
	result, err := script.RunContext(context.Background())

	// 验证执行是否成功
	if err != nil {
		stats := script.GetExecutionStats()
		t.Logf("Execution failed with %d instructions executed", stats.InstructionCount)
		t.Errorf("Expected successful execution, but got error: %v", err)
		return
	}

	// 验证返回结果是否正确
	if result != 42 {
		t.Errorf("Expected result 42, but got: %v", result)
	}

	// 验证执行统计信息
	stats := script.GetExecutionStats()
	if stats.InstructionCount >= 1000 {
		t.Errorf("Expected less than 1000 instructions executed, but got: %d", stats.InstructionCount)
	}
}

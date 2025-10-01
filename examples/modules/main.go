package main

import (
	"fmt"
	"log"

	"github.com/lengzhao/goscript/builtin"
	"github.com/lengzhao/goscript/vm"
)

func main() {
	// Create a new VM
	vmInstance := vm.NewVM()

	// Register builtin modules using the new ModuleExecutor interface
	modules := []string{"strings", "math"}
	for _, moduleName := range modules {
		moduleExecutor, exists := builtin.GetModuleExecutor(moduleName)
		if exists {
			vmInstance.RegisterModule(moduleName, moduleExecutor)
			fmt.Printf("Registered builtin module: %s\n", moduleName)
		}
	}

	// Test strings module function
	toUpperFn, exists := vmInstance.GetFunction("strings.ToUpper")
	if !exists {
		log.Fatal("Function strings.ToUpper not found")
	}

	result, err := toUpperFn("hello")
	if err != nil {
		log.Fatalf("Failed to call strings.ToUpper: %v", err)
	}

	fmt.Printf("strings.ToUpper('hello') = %v\n", result)

	// Test math module function
	absFn, exists := vmInstance.GetFunction("math.Abs")
	if !exists {
		log.Fatal("Function math.Abs not found")
	}

	result, err = absFn(-5)
	if err != nil {
		log.Fatalf("Failed to call math.Abs: %v", err)
	}

	fmt.Printf("math.Abs(-5) = %v\n", result)
}

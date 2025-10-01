package test

import (
	"testing"

	_ "embed"

	"github.com/lengzhao/goscript"
	"github.com/lengzhao/goscript/builtin"
)

//go:embed module_test/math.gs
var mathCode string

//go:embed module_test/main.gs
var mainCode string

func TestModuleGSFileCompilation(t *testing.T) {
	// Create a VM for the module
	moduleVM := goscript.NewScript([]byte(mathCode))

	moduleVM.Build()

	// Create a script with the main code
	script := goscript.NewScript([]byte(mainCode))
	script.RegisterModule("math", moduleVM.CallFunction)

	// Run the script
	result, err := script.Run()
	if err != nil {
		t.Fatalf("Failed to run script: %v", err)
	}

	// Check the result
	// The main.gs file calculates a result, let's just check that it's not nil
	if result == nil {
		t.Errorf("Expected result to be non-nil, got %v", result)
	}
}

func TestModuleFunctionCallDirect(t *testing.T) {
	// Create a script
	script := goscript.NewScript([]byte{})

	// Get the VM from the script
	vmInstance := script.GetVM()

	// Register builtin modules
	modules := []string{"strings", "math"}
	for _, moduleName := range modules {
		moduleExecutor, exists := builtin.GetModuleExecutor(moduleName)
		if exists {
			vmInstance.RegisterModule(moduleName, moduleExecutor)
		}
	}

	// Try to get the math.Max function directly from the VM
	fn, exists := vmInstance.GetFunction("math.Max")
	if !exists {
		t.Fatalf("Function math.Max not found")
	}

	// Call the function
	result, err := fn(5, 10)
	if err != nil {
		t.Fatalf("Failed to call math.Max: %v", err)
	}
	if result != 10 {
		t.Errorf("Expected math.Max(5, 10) to be 10, got %v", result)
	}
}

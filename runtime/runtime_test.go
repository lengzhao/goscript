package runtime

import (
	"fmt"
	"reflect"
	"testing"
)

func TestRuntime(t *testing.T) {
	rt := NewRuntime()

	// Test variable management
	rt.SetVariable("testVar", 42)
	value, ok := rt.GetVariable("testVar")
	if !ok {
		t.Error("Expected to find testVar")
	}
	if value != 42 {
		t.Errorf("Expected value 42, got %v", value)
	}

	rt.DeleteVariable("testVar")
	_, ok = rt.GetVariable("testVar")
	if ok {
		t.Error("Expected testVar to be deleted")
	}

	// Test function registration
	testFn := NewBuiltInFunction("testFunc", func(args ...interface{}) (interface{}, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("expected 2 arguments")
		}
		a, ok1 := args[0].(int)
		b, ok2 := args[1].(int)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("arguments must be integers")
		}
		return a + b, nil
	})

	rt.RegisterFunction("testFunc", testFn)
	fn, ok := rt.GetFunction("testFunc")
	if !ok {
		t.Error("Expected to find testFunc")
	}
	if fn.Name() != "testFunc" {
		t.Errorf("Expected function name 'testFunc', got '%s'", fn.Name())
	}

	// Test function execution
	result, err := fn.Call(10, 20)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != 30 {
		t.Errorf("Expected result 30, got %v", result)
	}

	// Test type registration
	rt.RegisterType("int", reflect.TypeOf(0))
	typ, ok := rt.GetType("int")
	if !ok {
		t.Error("Expected to find type 'int'")
	}
	if typ.Kind() != reflect.Int {
		t.Errorf("Expected int type, got %v", typ.Kind())
	}
}

func TestModule(t *testing.T) {
	rt := NewRuntime()

	// Create a mock module
	mockModule := &mockModule{
		name: "mock",
		functions: map[string]Function{
			"mockFunc": NewBuiltInFunction("mockFunc", func(args ...interface{}) (interface{}, error) {
				return "mock result", nil
			}),
		},
		types: map[string]reflect.Type{
			"mockType": reflect.TypeOf(""),
		},
	}

	// Test module import
	err := rt.ImportModule(mockModule)
	if err != nil {
		t.Errorf("Expected no error importing module, got %v", err)
	}

	// Try to import the same module again
	err = rt.ImportModule(mockModule)
	if err == nil {
		t.Error("Expected error when importing duplicate module")
	}

	// Test module retrieval
	module, ok := rt.GetModule("mock")
	if !ok {
		t.Error("Expected to find mock module")
	}
	if module.Name() != "mock" {
		t.Errorf("Expected module name 'mock', got '%s'", module.Name())
	}

	// Test function from module
	fn, ok := module.GetFunction("mockFunc")
	if !ok {
		t.Error("Expected to find mockFunc in module")
	}
	result, err := fn.Call()
	if err != nil {
		t.Errorf("Expected no error calling mockFunc, got %v", err)
	}
	if result != "mock result" {
		t.Errorf("Expected 'mock result', got %v", result)
	}

	// Test type from module
	typ, ok := module.GetType("mockType")
	if !ok {
		t.Error("Expected to find mockType in module")
	}
	if typ.Kind() != reflect.String {
		t.Errorf("Expected string type, got %v", typ.Kind())
	}
}

// Mock module for testing
type mockModule struct {
	name      string
	functions map[string]Function
	types     map[string]reflect.Type
}

func (m *mockModule) Name() string {
	return m.name
}

func (m *mockModule) GetFunction(name string) (Function, bool) {
	fn, ok := m.functions[name]
	return fn, ok
}

func (m *mockModule) GetType(name string) (reflect.Type, bool) {
	typ, ok := m.types[name]
	return typ, ok
}

// Package builtin provides built-in functions and modules for GoScript
package builtin

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"github.com/lengzhao/goscript/types"
)

// Strings module functions
var StringsModule = map[string]types.Function{
	"Contains": func(args ...interface{}) (interface{}, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("contains function requires 2 arguments")
		}
		s, ok1 := args[0].(string)
		substr, ok2 := args[1].(string)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("contains function requires string arguments")
		}
		return strings.Contains(s, substr), nil
	},
	"HasPrefix": func(args ...interface{}) (interface{}, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("hasPrefix function requires 2 arguments")
		}
		s, ok1 := args[0].(string)
		prefix, ok2 := args[1].(string)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("hasPrefix function requires string arguments")
		}
		return strings.HasPrefix(s, prefix), nil
	},
	"HasSuffix": func(args ...interface{}) (interface{}, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("hasSuffix function requires 2 arguments")
		}
		s, ok1 := args[0].(string)
		suffix, ok2 := args[1].(string)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("hasSuffix function requires string arguments")
		}
		return strings.HasSuffix(s, suffix), nil
	},
	"ToLower": func(args ...interface{}) (interface{}, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("toLower function requires 1 argument")
		}
		s, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("toLower function requires string argument")
		}
		return strings.ToLower(s), nil
	},
	"ToUpper": func(args ...interface{}) (interface{}, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("toUpper function requires 1 argument")
		}
		s, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("toUpper function requires string argument")
		}
		return strings.ToUpper(s), nil
	},
	"Trim": func(args ...interface{}) (interface{}, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("trim function requires 2 arguments")
		}
		s, ok1 := args[0].(string)
		cutset, ok2 := args[1].(string)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("trim function requires string arguments")
		}
		return strings.Trim(s, cutset), nil
	},
	"TrimSpace": func(args ...interface{}) (interface{}, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("trimSpace function requires 1 argument")
		}
		s, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("trimSpace function requires string argument")
		}
		return strings.TrimSpace(s), nil
	},
	"Split": func(args ...interface{}) (interface{}, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("split function requires 2 arguments")
		}
		s, ok1 := args[0].(string)
		sep, ok2 := args[1].(string)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("split function requires string arguments")
		}
		result := strings.Split(s, sep)
		// Convert []string to []interface{}
		interfaceResult := make([]interface{}, len(result))
		for i, v := range result {
			interfaceResult[i] = v
		}
		return interfaceResult, nil
	},
	"Join": func(args ...interface{}) (interface{}, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("join function requires 2 arguments")
		}
		slice, ok1 := args[0].([]interface{})
		sep, ok2 := args[1].(string)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("join function requires []interface{} and string arguments")
		}
		// Convert []interface{} to []string
		stringsSlice := make([]string, len(slice))
		for i, v := range slice {
			if str, ok := v.(string); ok {
				stringsSlice[i] = str
			} else {
				stringsSlice[i] = fmt.Sprintf("%v", v)
			}
		}
		return strings.Join(stringsSlice, sep), nil
	},
}

// Fmt module functions
var FmtModule = map[string]types.Function{
	"Printf": func(args ...interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("printf function requires at least 1 argument")
		}
		format, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("first argument to printf must be a string")
		}
		// For simplicity, we'll just return the formatted string
		// In a real implementation, this would actually print to stdout
		if len(args) == 1 {
			return format, nil
		}
		return fmt.Sprintf(format, args[1:]...), nil
	},
	"Println": func(args ...interface{}) (interface{}, error) {
		// Print all arguments with spaces between them and a newline at the end
		fmt.Println(args...)
		// Return nil as Println doesn't return a value
		return nil, nil
	},
	"Sprintf": func(args ...interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("sprintf function requires at least 1 argument")
		}
		format, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("first argument to sprintf must be a string")
		}
		return fmt.Sprintf(format, args[1:]...), nil
	},
	"Sprint": func(args ...interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("sprint function requires at least 1 argument")
		}
		return fmt.Sprint(args...), nil
	},
}

// Math module functions
var MathModule = map[string]types.Function{
	"Abs": func(args ...interface{}) (interface{}, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("abs function requires 1 argument")
		}
		switch v := args[0].(type) {
		case int:
			if v < 0 {
				return -v, nil
			}
			return v, nil
		case float64:
			return math.Abs(v), nil
		default:
			return nil, fmt.Errorf("abs function requires numeric argument, got %T", v)
		}
	},
	"Max": func(args ...interface{}) (interface{}, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("max function requires 2 arguments")
		}
		switch a := args[0].(type) {
		case int:
			if b, ok := args[1].(int); ok {
				if a > b {
					return a, nil
				}
				return b, nil
			}
		case float64:
			if b, ok := args[1].(float64); ok {
				return math.Max(a, b), nil
			}
		}
		return nil, fmt.Errorf("max function requires numeric arguments of the same type")
	},
	"Min": func(args ...interface{}) (interface{}, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("min function requires 2 arguments")
		}
		switch a := args[0].(type) {
		case int:
			if b, ok := args[1].(int); ok {
				if a < b {
					return a, nil
				}
				return b, nil
			}
		case float64:
			if b, ok := args[1].(float64); ok {
				return math.Min(a, b), nil
			}
		}
		return nil, fmt.Errorf("min function requires numeric arguments of the same type")
	},
	"Sqrt": func(args ...interface{}) (interface{}, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("sqrt function requires 1 argument")
		}
		if v, ok := args[0].(float64); ok {
			return math.Sqrt(v), nil
		}
		return nil, fmt.Errorf("sqrt function requires float64 argument")
	},
}

// JSON module functions
var JSONModule = map[string]types.Function{
	"Marshal": func(args ...interface{}) (interface{}, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("marshal function requires 1 argument")
		}
		// Convert Go value to JSON
		jsonData, err := json.Marshal(args[0])
		if err != nil {
			return nil, fmt.Errorf("failed to marshal to JSON: %w", err)
		}
		return string(jsonData), nil
	},
	"Unmarshal": func(args ...interface{}) (interface{}, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("unmarshal function requires 1 argument")
		}
		jsonStr, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("unmarshal function requires string argument")
		}
		// Convert JSON string to Go value
		var result interface{}
		err := json.Unmarshal([]byte(jsonStr), &result)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
		}
		return result, nil
	},
}

// GetModuleFunctions returns the functions for a given module
func GetModuleFunctions(moduleName string) (map[string]types.Function, bool) {
	switch moduleName {
	case "strings":
		return StringsModule, true
	case "fmt":
		return FmtModule, true
	case "math":
		return MathModule, true
	case "json":
		return JSONModule, true
	default:
		return nil, false
	}
}

// GetModuleExecutor returns a ModuleExecutor for a given module
func GetModuleExecutor(moduleName string) (types.ModuleExecutor, bool) {
	// Get the module functions
	moduleFuncs, exists := GetModuleFunctions(moduleName)
	if !exists {
		return nil, false
	}

	// Create a ModuleExecutor that delegates to the module functions
	moduleExecutor := func(entrypoint string, args ...interface{}) (interface{}, error) {
		// Look up the function in the module
		if fn, exists := moduleFuncs[entrypoint]; exists {
			// Call the function with the provided arguments
			return fn(args...)
		}
		return nil, fmt.Errorf("function %s not found in module %s", entrypoint, moduleName)
	}

	return moduleExecutor, true
}

func ListAllModules() []string {
	return []string{"strings", "fmt", "math", "json"}
}

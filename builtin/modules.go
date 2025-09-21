// Package builtin provides built-in functions and modules for GoScript
package builtin

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Strings module functions
var StringsModule = map[string]Function{
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
var FmtModule = map[string]Function{
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

// JSON module functions
var JSONModule = map[string]Function{
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
func GetModuleFunctions(moduleName string) (map[string]Function, bool) {
	switch moduleName {
	case "strings":
		return StringsModule, true
	case "fmt":
		return FmtModule, true
	case "json":
		return JSONModule, true
	default:
		return nil, false
	}
}

func ListAllModules() []string {
	return []string{"strings", "fmt", "json"}
}

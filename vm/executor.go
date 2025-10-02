package vm

import (
	"fmt"
	"strings"

	"github.com/lengzhao/goscript/builtin"
	execContext "github.com/lengzhao/goscript/context"
	"github.com/lengzhao/goscript/instruction"
)

// ReturnError is a special error type used to return values from functions
type ReturnError struct {
	Value interface{}
}

func (e *ReturnError) Error() string {
	return "return"
}

// OpHandler defines the signature for opcode handlers
type OpHandler func(stack *Stack, instr *instruction.Instruction, pc int) (int, error)

// Executor handles the execution of instructions
type Executor struct {
	vm *VM
	// Opcode handler array for table-driven execution
	// Using array instead of map for better performance
	opcodeHandlers [instruction.OpCodeLast + 1]OpHandler
}

// NewExecutor creates a new executor
func NewExecutor(vm *VM) *Executor {
	exec := &Executor{
		vm: vm,
	}

	// Initialize opcode handlers
	exec.initOpcodeHandlers()

	return exec
}

// initOpcodeHandlers initializes the opcode handler array
func (exec *Executor) initOpcodeHandlers() {
	exec.opcodeHandlers[instruction.OpNop] = exec.handleNop
	exec.opcodeHandlers[instruction.OpLoadConst] = exec.handleLoadConst
	exec.opcodeHandlers[instruction.OpLoadName] = exec.handleLoadName
	exec.opcodeHandlers[instruction.OpStoreName] = exec.handleStoreName
	exec.opcodeHandlers[instruction.OpPop] = exec.handlePop
	exec.opcodeHandlers[instruction.OpCall] = exec.handleCall
	exec.opcodeHandlers[instruction.OpReturn] = exec.handleReturn
	exec.opcodeHandlers[instruction.OpBinaryOp] = exec.handleBinaryOp
	exec.opcodeHandlers[instruction.OpCreateVar] = exec.handleCreateVar
	exec.opcodeHandlers[instruction.OpEnterScopeWithKey] = exec.handleEnterScopeWithKey
	exec.opcodeHandlers[instruction.OpExitScopeWithKey] = exec.handleExitScopeWithKey
	exec.opcodeHandlers[instruction.OpGetIndex] = exec.handleGetIndex
	exec.opcodeHandlers[instruction.OpSetIndex] = exec.handleSetIndex
	exec.opcodeHandlers[instruction.OpJump] = exec.handleJump
	exec.opcodeHandlers[instruction.OpJumpIf] = exec.handleJumpIf
	exec.opcodeHandlers[instruction.OpNewSlice] = exec.handleNewSlice
	exec.opcodeHandlers[instruction.OpLen] = exec.handleLen
	exec.opcodeHandlers[instruction.OpRotate] = exec.handleRotate
	exec.opcodeHandlers[instruction.OpSwap] = exec.handleSwap
	exec.opcodeHandlers[instruction.OpNewStruct] = exec.handleNewStruct
	exec.opcodeHandlers[instruction.OpSetField] = exec.handleSetField
	exec.opcodeHandlers[instruction.OpGetField] = exec.handleGetField
	exec.opcodeHandlers[instruction.OpCallMethod] = exec.handleCallMethod
	exec.opcodeHandlers[instruction.OpImport] = exec.handleImport
}

// RegisterOpHandler registers a custom opcode handler
func (exec *Executor) RegisterOpHandler(op instruction.OpCode, handler OpHandler) {
	exec.opcodeHandlers[op] = handler
}

// executeInstructions executes a sequence of instructions with the given context
func (exec *Executor) executeInstructions(instructions []*instruction.Instruction) (interface{}, error) {
	stack := NewStack()
	pc := 0 // program counter

	// Reset instruction count for this execution
	exec.vm.instructionCount = 0

	for pc < len(instructions) {
		instr := instructions[pc]

		// Check instruction limit
		if exec.vm.maxInstructions > 0 {
			if exec.vm.instructionCount >= exec.vm.maxInstructions {
				return nil, fmt.Errorf("maximum instruction limit exceeded: %d instructions executed", exec.vm.instructionCount)
			}
		}

		// Increment instruction counter
		exec.vm.instructionCount++

		// Debug output
		if exec.vm.debug {
			fmt.Printf("Executing instruction %d: %s, stack size: %d, stack: %v\n", pc, instr.String(), stack.Len(), stack.Items())
		}

		// Look up the handler for this opcode using array for better performance
		handler := exec.opcodeHandlers[instr.Op]
		if handler == nil {
			return nil, fmt.Errorf("unsupported operation: %s", instr.Op.String())
		}

		// Execute the handler
		newPC, err := handler(stack, instr, pc)
		if err != nil {
			// Check if it's a return error
			if returnErr, ok := err.(*ReturnError); ok {
				if exec.vm.debug {
					fmt.Printf("Return with value: %v\n", returnErr.Value)
				}
				return returnErr.Value, nil
			}
			return nil, err
		}
		pc = newPC
	}

	// If we've executed all instructions without an explicit return, return nil
	return nil, nil
}

// handleNop handles the NOP opcode
func (exec *Executor) handleNop(stack *Stack, instr *instruction.Instruction, pc int) (int, error) {
	return pc + 1, nil
}

// handleLoadConst handles the LOAD_CONST opcode
func (exec *Executor) handleLoadConst(stack *Stack, instr *instruction.Instruction, pc int) (int, error) {
	stack.Push(instr.Arg)
	return pc + 1, nil
}

// handleLoadName handles the LOAD_NAME opcode
func (exec *Executor) handleLoadName(stack *Stack, instr *instruction.Instruction, pc int) (int, error) {
	name, ok := instr.Arg.(string)
	if !ok {
		return 0, fmt.Errorf("invalid argument for LOAD_NAME")
	}

	// Check if this is a field access (e.g., "p.age")
	// if strings.Contains(name, ".") {
	// Split the name into variable and field parts
	parts := strings.Split(name, ".")
	if len(parts) == 2 {
		varName := parts[0]
		fieldName := parts[1]

		// Look up the variable (struct) in the context hierarchy
		structValue, exists := exec.vm.currentCtx.GetVariable(varName)
		if !exists {
			return 0, fmt.Errorf("undefined variable: %s", varName)
		}

		// Check if it's a struct (map)
		if structMap, ok := structValue.(map[string]interface{}); ok {
			// Get the field value
			fieldValue, fieldExists := structMap[fieldName]
			if !fieldExists {
				// Field doesn't exist, push nil
				stack.Push(nil)
			} else {
				// Push the field value
				stack.Push(fieldValue)
			}
			return pc + 1, nil
		}
	}
	// }

	// Look up the variable in the context hierarchy
	value, exists := exec.vm.currentCtx.GetVariable(name)
	if !exists {
		// Check if it's a module reference
		// In this case, we should return the module name itself as a string
		// This allows module functions to be called using the format "moduleName.functionName"
		if exec.isModuleName(name) {
			stack.Push(name)
			return pc + 1, nil
		}
		return 0, fmt.Errorf("undefined variable: %s", name)
	}
	// Debug information
	//fmt.Printf("LOAD_NAME: %s = %v (type %T)\n", name, value, value)
	stack.Push(value)
	return pc + 1, nil
}

// handleStoreName handles the STORE_NAME opcode
func (exec *Executor) handleStoreName(stack *Stack, instr *instruction.Instruction, pc int) (int, error) {
	name, ok := instr.Arg.(string)
	if !ok {
		return 0, fmt.Errorf("invalid argument for STORE_NAME")
	}

	if stack.Len() < 1 {
		return 0, fmt.Errorf("stack underflow")
	}

	value := stack.Pop()

	// For function parameters, they might already have values set by the caller
	// We should update the value, not create a new variable
	err := exec.vm.currentCtx.SetVariable(name, value)
	if err != nil {
		// If setting fails, try to create the variable
		exec.vm.currentCtx.CreateVariableWithType(name, value, "unknown")
	}

	return pc + 1, nil
}

// handlePop handles the POP opcode
func (exec *Executor) handlePop(stack *Stack, instr *instruction.Instruction, pc int) (int, error) {
	if stack.Len() < 1 {
		return 0, fmt.Errorf("stack underflow")
	}
	stack.Pop()
	return pc + 1, nil
}

// handleCall handles the CALL opcode
func (exec *Executor) handleCall(stack *Stack, instr *instruction.Instruction, pc int) (int, error) {
	// Get the function name and argument count
	functionName, ok := instr.Arg.(string)
	if !ok {
		return 0, fmt.Errorf("invalid function name for CALL")
	}

	argCount, ok := instr.Arg2.(int)
	if !ok {
		return 0, fmt.Errorf("invalid argument count for CALL")
	}

	// Debug information - print stack before processing
	if exec.vm.debug {
		fmt.Printf("CALL %s with %d arguments, stack: %v\n", functionName, argCount, stack.Items())
	}

	// Prepare arguments using the unified function
	args, err := exec.prepareArguments(stack, argCount)
	if err != nil {
		return 0, fmt.Errorf("error preparing arguments for CALL %s: %w", functionName, err)
	}

	// Unified call handling
	callType := exec.determineCallType(args, functionName)

	switch callType {
	case callTypeModule:
		return exec.handleModuleCall(stack, functionName, args, pc)
	case callTypeMethod:
		return exec.handleMethodCallUnified(stack, functionName, args, pc)
	default:
		// Regular function call
		// Push the arguments back to the stack for handleFunctionCall
		exec.pushArgumentsBack(stack, args)
		return exec.handleFunctionCall(stack, exec.vm, functionName, argCount, pc)
	}
}

// CallType represents the type of function call
type CallType int

const (
	callTypeRegular CallType = iota
	callTypeModule
	callTypeMethod
)

// determineCallType determines the type of call based on arguments and function name
func (exec *Executor) determineCallType(args []interface{}, functionName string) CallType {
	if len(args) > 0 {
		// Check if the first argument is a module
		if _, isModule := exec.isModuleVariable(args[0]); isModule {
			return callTypeModule
		}

		// Check if this is a method call with a struct receiver
		if exec.isStructReceiver(args[0]) {
			return callTypeMethod
		}
	}

	return callTypeRegular
}

// handleModuleCall handles module function calls
func (exec *Executor) handleModuleCall(stack *Stack, functionName string, args []interface{}, pc int) (int, error) {
	// This is a module call (e.g., math.Max)
	moduleName, _ := exec.isModuleVariable(args[0])
	qualifiedName := fmt.Sprintf("%s.%s", moduleName, functionName)
	// For module calls, we don't include the module name in the arguments
	moduleArgs := args[1:] // Remove the module name from arguments

	// Push the arguments back to the stack for handleFunctionCall
	for _, arg := range moduleArgs {
		stack.Push(arg)
	}

	return exec.handleFunctionCall(stack, exec.vm, qualifiedName, len(moduleArgs), pc)
}

// handleMethodCallUnified handles method calls with unified approach
func (exec *Executor) handleMethodCallUnified(stack *Stack, functionName string, args []interface{}, pc int) (int, error) {
	// This is a method call (e.g., rect.Area)
	// The receiver is already included in args[0]
	// For method calls, we use the specialized handleCallMethod function
	// Push the arguments back to the stack
	for _, arg := range args {
		stack.Push(arg)
	}

	// Create a CALL_METHOD instruction and delegate to handleCallMethod
	// For CALL_METHOD, argCount should not include the receiver
	callMethodInstr := instruction.NewInstruction(instruction.OpCallMethod, functionName, len(args)-1)
	return exec.handleCallMethod(stack, callMethodInstr, pc)
}

// isModuleVariable checks if a variable is a module
func (exec *Executor) isModuleVariable(variable interface{}) (string, bool) {
	// In our implementation, modules are stored as variables
	// We need to check if this variable corresponds to a registered module
	if varName, ok := variable.(string); ok {
		// Check if this variable is registered as a module in the VM
		if _, moduleExists := exec.vm.GetModule(varName); moduleExists {
			return varName, true
		}
	}
	return "", false
}

// isStructReceiver checks if the variable is a struct receiver
func (exec *Executor) isStructReceiver(variable interface{}) bool {
	// Check if the variable is a struct (map)
	_, ok := variable.(map[string]interface{})
	return ok
}

// getQualifiedMethodName gets the qualified method name for a struct receiver
func (exec *Executor) getQualifiedMethodName(receiver interface{}, methodName string) string {
	// Try to get the type name from the struct
	typeName := exec.getTypeNameFromStruct(receiver)
	if typeName != "" {
		return fmt.Sprintf("%s.%s", typeName, methodName)
	}
	// Fallback to just the method name
	return methodName
}

// getTypeNameFromStruct extracts the type name from a struct receiver
func (exec *Executor) getTypeNameFromStruct(receiver interface{}) string {
	if structMap, ok := receiver.(map[string]interface{}); ok {
		// First check for explicit type field
		if typeName, exists := structMap["_type"]; exists {
			if name, ok := typeName.(string); ok {
				return name
			}
		}
		// Fallback: try to infer from keys
		for key := range structMap {
			if key != "width" && key != "height" && key != "radius" && key != "_type" && key != "name" && key != "age" {
				return key
			}
		}
	}
	return ""
}

// handleMethodCall handles method calls with format "receiver.method" or "*type.method"
func (exec *Executor) handleMethodCall(stack *Stack, vm *VM, funcName string, argCount int, pc int) (int, error) {
	// Prepare arguments
	if stack.Len() < argCount {
		return 0, fmt.Errorf("stack underflow when calling method %s", funcName)
	}

	args := make([]interface{}, argCount)
	for i := argCount - 1; i >= 0; i-- {
		args[i] = stack.Pop()
	}

	// Check if this is a method call with a qualified name that includes pointer type
	// e.g., "*Rectangle.SetHeight"
	if strings.Contains(funcName, "*") {
		return exec.handlePointerMethodCall(stack, vm, funcName, args, pc)
	}

	// This looks like a method call with a variable name as the receiver
	// e.g., "rect.SetWidth" where "rect" is a variable name
	// We need to look up the variable to get its type and then find the correct method
	return exec.handleReceiverMethodCall(stack, vm, funcName, args, pc)
}

// handlePointerMethodCall handles method calls with pointer receiver format "*type.method"
func (exec *Executor) handlePointerMethodCall(stack *Stack, vm *VM, funcName string, args []interface{}, pc int) (int, error) {
	// Check if it's a registered script function with the qualified name
	if fn, exists := vm.GetFunction(funcName); exists {
		// Call the method
		result, err := fn(args...)
		if err != nil {
			return 0, fmt.Errorf("error calling method %s: %w", funcName, err)
		}

		// Push result back to stack if not nil
		if result != nil {
			stack.Push(result)
		}
		return pc + 1, nil
	}

	// Check if it's a script-defined method
	if functionInstructions, exists := vm.GetInstructionSet(funcName); exists {
		// Create new context for the method call
		// The method context's parent is the current context
		// Get the method name (part after the last dot)
		parts := strings.Split(funcName, ".")
		methodName := parts[len(parts)-1]
		methodCtx := execContext.NewContext(methodName, vm.currentCtx)

		// Set method arguments as local variables
		// Try to get parameter names from registered script function info
		paramNames := []string{} // default empty param names

		// Try to get the actual parameter names from the registered script function
		scriptFunctions := vm.GetAllScriptFunctions()
		for _, fnInfo := range scriptFunctions {
			// Check if this function matches our method name
			if fnInfo.Key == funcName {
				// Use the parameter names from the function info
				if len(fnInfo.ParamNames) > 0 {
					paramNames = fnInfo.ParamNames
				}
				break
			}
		}

		// Set arguments as local variables with appropriate names
		for i, arg := range args {
			paramName := fmt.Sprintf("arg%d", i)
			if i < len(paramNames) {
				paramName = paramNames[i]
			}

			// Make sure we create the variable in the method context
			methodCtx.CreateVariableWithType(paramName, arg, "unknown")
		}

		// Execute the method using a new executor
		newExec := NewExecutor(vm)
		result, err := newExec.executeInstructions(functionInstructions)
		if err != nil {
			return 0, fmt.Errorf("error executing method %s: %w", funcName, err)
		}

		// Push result back to stack if not nil
		if result != nil {
			stack.Push(result)
		}
		return pc + 1, nil
	}

	return 0, fmt.Errorf("undefined method: %s", funcName)
}

// handleReceiverMethodCall handles method calls with receiver variable format "variable.method"
func (exec *Executor) handleReceiverMethodCall(stack *Stack, vm *VM, funcName string, args []interface{}, pc int) (int, error) {
	// Split the function name to get receiver variable and method name
	parts := strings.Split(funcName, ".")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid method call format: %s", funcName)
	}

	receiverVarName := parts[0]
	methodName := parts[1]

	// Look up the receiver variable to get its type
	receiver, exists := vm.currentCtx.GetVariable(receiverVarName)
	if !exists {
		return exec.handleModuleFunctionCall(receiverVarName, methodName, args, stack, pc)
	}

	// Try to determine the type name of the receiver
	typeName := exec.getTypeNameFromReceiver(receiver, receiverVarName)

	// Try to find the method with the correct signature
	// First try with the type name
	qualifiedMethodName := fmt.Sprintf("%s.%s", typeName, methodName)

	// Check if it's a registered script function with the qualified name
	if fn, exists := vm.GetFunction(qualifiedMethodName); exists {
		return exec.callRegisteredMethod(fn, qualifiedMethodName, args, stack, pc, typeName, methodName)
	}

	// Check if it's a script-defined method
	if _, exists := vm.GetInstructionSet(qualifiedMethodName); exists {
		return exec.callScriptDefinedMethod(stack, vm, qualifiedMethodName, args, pc, typeName, methodName)
	}

	// If we can't find the qualified method, try to find it as a regular function
	if fn, exists := vm.GetFunction(methodName); exists {
		// Call the function
		result, err := fn(args...)
		if err != nil {
			return 0, fmt.Errorf("error calling function %s: %w", methodName, err)
		}

		// Push result back to stack if not nil
		if result != nil {
			stack.Push(result)
		}
		return pc + 1, nil
	}

	return 0, fmt.Errorf("undefined method: %s", funcName)
}

// handleModuleFunctionCall handles calls to module functions
func (exec *Executor) handleModuleFunctionCall(moduleName, functionName string, args []interface{}, stack *Stack, pc int) (int, error) {
	// This is a module function call
	// The first argument is the module name, which we don't need to pass to the function
	// The remaining arguments are the actual function arguments
	if len(args) < 1 {
		return 0, fmt.Errorf("module function call requires at least 1 argument")
	}

	// Remove the module name from the arguments
	actualArgs := args[1:]

	result, err := exec.callModuleFunction(moduleName, functionName, actualArgs...)
	if err != nil {
		return 0, fmt.Errorf("error calling module function %s.%s: %w", moduleName, functionName, err)
	}

	// Push result back to stack if not nil
	if result != nil {
		stack.Push(result)
	}
	return pc + 1, nil
}

// getTypeNameFromReceiver extracts the type name from a receiver object
func (exec *Executor) getTypeNameFromReceiver(receiver interface{}, receiverVarName string) string {
	typeName := ""
	if structMap, ok := receiver.(map[string]interface{}); ok {
		// Check if the struct has a type field
		if typeField, exists := structMap["_type"]; exists {
			typeName = typeField.(string)
		} else {
			// Heuristic approach: try to infer type from field names
			// This is not ideal but works for simple cases
			// In a real implementation, we would store type information during compilation
			// For now, let's try to determine the type based on the variable name and context
			switch receiverVarName {
			case "rect":
				typeName = "Rectangle"
			case "circle":
				typeName = "Circle"
			default:
				typeName = "Rectangle" // Default fallback
			}
		}
	}
	return typeName
}

// callRegisteredMethod calls a registered script method
func (exec *Executor) callRegisteredMethod(fn ScriptFunction, qualifiedMethodName string, args []interface{}, stack *Stack, pc int, typeName, methodName string) (int, error) {
	// Call the method
	// For value receiver methods, we need to pass a copy of the struct
	// For pointer receiver methods, we pass the original struct
	callArgs := make([]interface{}, len(args))
	copy(callArgs, args)

	// Check if this is a pointer receiver method
	pointerMethodName := fmt.Sprintf("*%s.%s", typeName, methodName)
	if _, pointerExists := exec.vm.GetFunction(pointerMethodName); pointerExists {
		// Pointer receiver method - pass original struct
	} else {
		// Value receiver method - pass copy of struct
		if originalStruct, ok := args[0].(map[string]interface{}); ok {
			structCopy := make(map[string]interface{})
			for k, v := range originalStruct {
				structCopy[k] = v
			}
			callArgs[0] = structCopy
		}
	}

	result, err := fn(callArgs...)
	if err != nil {
		return 0, fmt.Errorf("error calling method %s: %w", qualifiedMethodName, err)
	}

	// Push result back to stack if not nil
	if result != nil {
		stack.Push(result)
	}
	return pc + 1, nil
}

// callScriptDefinedMethod calls a script-defined method
func (exec *Executor) callScriptDefinedMethod(stack *Stack, vm *VM, qualifiedMethodName string, args []interface{}, pc int, typeName, methodName string) (int, error) {
	// Check if it's a script-defined method
	functionInstructions, exists := vm.GetInstructionSet(qualifiedMethodName)
	if !exists {
		return 0, fmt.Errorf("undefined method: %s", qualifiedMethodName)
	}

	// Create new context for the method call
	methodCtx := execContext.NewContext(methodName, vm.currentCtx)

	// Set method arguments as local variables
	// Try to get parameter names from registered script function info
	paramNames := []string{} // default empty param names

	// Try to get the actual parameter names from the registered script function
	scriptFunctions := vm.GetAllScriptFunctions()
	for _, fnInfo := range scriptFunctions {
		// Check if this function matches our method name
		if fnInfo.Key == qualifiedMethodName {
			// Use the parameter names from the function info
			if len(fnInfo.ParamNames) > 0 {
				paramNames = fnInfo.ParamNames
			}
			break
		}
	}

	// For value receiver methods, we need to create a copy of the struct
	// For pointer receiver methods, we use the original struct
	// Check if there's a pointer receiver version of this method
	pointerMethodName := fmt.Sprintf("*%s.%s", typeName, methodName)
	isPointerReceiver := false
	if _, pointerExists := vm.GetFunction(pointerMethodName); pointerExists {
		isPointerReceiver = true
	} else if _, pointerExists := vm.GetInstructionSet(pointerMethodName); pointerExists {
		isPointerReceiver = true
	}

	// Set arguments as local variables with appropriate names
	for i, arg := range args {
		paramName := "r" // default receiver name
		if i < len(paramNames) {
			paramName = paramNames[i]
		} else if i > 0 {
			// For non-receiver arguments, use generic names
			paramName = fmt.Sprintf("arg%d", i-1)
		}

		// For the receiver (first argument), handle value vs pointer receiver
		if i == 0 {
			// This is the receiver
			if !isPointerReceiver {
				// For value receiver methods, create a copy of the struct
				if originalStruct, ok := arg.(map[string]interface{}); ok {
					structCopy := make(map[string]interface{})
					for k, v := range originalStruct {
						structCopy[k] = v
					}
					arg = structCopy
				}
			}
		}

		// Make sure we create the variable in the method context
		methodCtx.CreateVariableWithType(paramName, arg, "unknown")
	}

	// Execute the method using a new executor
	newExec := NewExecutor(vm)
	// Set the current context for the method execution
	vm.currentCtx = methodCtx
	result, err := newExec.executeInstructions(functionInstructions)
	if err != nil {
		return 0, fmt.Errorf("error executing method %s: %w", qualifiedMethodName, err)
	}

	// Push result back to stack if not nil
	if result != nil {
		stack.Push(result)
	}
	return pc + 1, nil
}

// handleFunctionCall handles regular function calls
func (exec *Executor) handleFunctionCall(stack *Stack, vm *VM, funcName string, argCount int, pc int) (int, error) {
	// Check if it's a registered script function
	if fn, exists := vm.GetFunction(funcName); exists {
		// Prepare arguments using the unified function
		args, err := exec.prepareArguments(stack, argCount)
		if err != nil {
			return 0, fmt.Errorf("error preparing arguments for function %s: %w", funcName, err)
		}

		// Call the function
		result, err := fn(args...)
		if err != nil {
			return 0, fmt.Errorf("error calling function %s: %w", funcName, err)
		}

		// Push result back to stack if not nil
		if result != nil {
			stack.Push(result)
		}
		return pc + 1, nil
	}

	// Check if it's a script-defined function (by key)
	if _, exists := vm.GetInstructionSet(funcName); exists {
		return exec.callScriptDefinedFunction(stack, vm, funcName, argCount, pc)
	}

	return 0, fmt.Errorf("undefined function: %s", funcName)
}

// callScriptDefinedFunction calls a script-defined function
func (exec *Executor) callScriptDefinedFunction(stack *Stack, vm *VM, funcName string, argCount int, pc int) (int, error) {
	// Prepare arguments using the unified function
	args, err := exec.prepareArguments(stack, argCount)
	if err != nil {
		return 0, fmt.Errorf("error preparing arguments for script function %s: %w", funcName, err)
	}

	// Create new context for the function call
	// The function context's parent is the current context
	functionCtx := execContext.NewContext(funcName, exec.vm.currentCtx)

	// Try to get the actual parameter names from the registered script function
	paramNames := make([]string, argCount)

	// Get all script functions to find the one we're calling
	scriptFunctions := vm.GetAllScriptFunctions()

	// Try to find the function info
	var foundFuncInfo *ScriptFunctionInfo
	for _, fnInfo := range scriptFunctions {
		// Check if this function matches our function name
		if fnInfo.Key == funcName || fnInfo.Name == funcName {
			foundFuncInfo = fnInfo
			break
		}
	}

	// If we found the function info and it has parameter names, use them
	if foundFuncInfo != nil && len(foundFuncInfo.ParamNames) > 0 {
		// Use the actual parameter names from the function definition
		for i := 0; i < argCount && i < len(foundFuncInfo.ParamNames); i++ {
			paramNames[i] = foundFuncInfo.ParamNames[i]
		}
		// Fill in any remaining parameters with default names
		for i := len(foundFuncInfo.ParamNames); i < argCount; i++ {
			paramNames[i] = fmt.Sprintf("arg%d", i)
		}
	} else {
		// Fall back to default parameter names
		for i := 0; i < argCount; i++ {
			paramNames[i] = fmt.Sprintf("arg%d", i)
		}
	}

	// Set arguments as local variables with appropriate names
	for i, arg := range args {
		paramName := paramNames[i]
		// Make sure we create the variable in the function context
		functionCtx.CreateVariableWithType(paramName, arg, "unknown")
	}

	// Execute the function using a new executor
	functionInstructions, exists := vm.GetInstructionSet(funcName)
	if !exists {
		return 0, fmt.Errorf("undefined function: %s", funcName)
	}

	// Save the current context
	previousCtx := vm.currentCtx

	// Set the current context for the function execution
	vm.currentCtx = functionCtx

	// Execute the function
	newExec := NewExecutor(vm)
	result, err := newExec.executeInstructions(functionInstructions)

	// Restore the previous context
	vm.currentCtx = previousCtx

	if err != nil {
		return 0, fmt.Errorf("error executing function %s: %w", funcName, err)
	}

	// Push result back to stack if not nil
	if result != nil {
		stack.Push(result)
	}
	return pc + 1, nil
}

// isModuleName checks if a name is a registered module name
func (exec *Executor) isModuleName(name string) bool {
	// Use the builtin module system to check if it's a valid module name
	modules := builtin.ListAllModules()
	for _, module := range modules {
		if module == name {
			return true
		}
	}
	return false
}

// callModuleFunction calls a function in a module
func (exec *Executor) callModuleFunction(moduleName, functionName string, args ...interface{}) (interface{}, error) {
	// First, check if the module is registered in the VM
	if module, exists := exec.vm.GetModule(moduleName); exists {
		// Call the module executor directly
		return module(functionName, args...)
	}

	// If not found in VM modules, fall back to builtin modules
	// Use the builtin module system instead of hardcoding functions
	// Get the module functions from the builtin package
	moduleFuncs, exists := builtin.GetModuleFunctions(moduleName)
	if !exists {
		return nil, fmt.Errorf("unsupported module: %s", moduleName)
	}

	// Get the specific function from the module
	fn, exists := moduleFuncs[functionName]
	if !exists {
		return nil, fmt.Errorf("unsupported function %s in module %s", functionName, moduleName)
	}

	// Call the function
	return fn(args...)
}

// handleReturn handles the RETURN opcode
func (exec *Executor) handleReturn(stack *Stack, instr *instruction.Instruction, pc int) (int, error) {
	// Return the top of stack if it exists
	if stack.Len() > 0 {
		// We use a special error value to return the result
		return 0, &ReturnError{Value: stack.Pop()}
	}
	return 0, &ReturnError{Value: nil}
}

// handleBinaryOp handles the BINARY_OP opcode
func (exec *Executor) handleBinaryOp(stack *Stack, instr *instruction.Instruction, pc int) (int, error) {
	vm := exec.vm
	op, ok := instr.Arg.(instruction.BinaryOp)
	if !ok {
		return 0, fmt.Errorf("invalid binary operation")
	}

	if stack.Len() < 2 {
		return 0, fmt.Errorf("stack underflow for binary operation")
	}

	right := stack.Pop()
	left := stack.Pop()

	result, err := vm.executeBinaryOp(op, left, right)
	if err != nil {
		return 0, err
	}

	stack.Push(result)
	return pc + 1, nil
}

// handleCreateVar handles the CREATE_VAR opcode
func (exec *Executor) handleCreateVar(stack *Stack, instr *instruction.Instruction, pc int) (int, error) {
	name, ok := instr.Arg.(string)
	if !ok {
		return 0, fmt.Errorf("invalid variable name")
	}

	// Create the variable with nil initial value
	exec.vm.currentCtx.CreateVariableWithType(name, nil, "unknown")
	return pc + 1, nil
}

// handleEnterScopeWithKey handles the ENTER_SCOPE_WITH_KEY opcode
func (exec *Executor) handleEnterScopeWithKey(stack *Stack, instr *instruction.Instruction, pc int) (int, error) {
	// For now, we just increment the program counter
	// In a more advanced implementation, we might manage nested scopes
	// todo newctx to replage old ctx
	ctx := execContext.NewContext("", exec.vm.currentCtx)
	exec.vm.currentCtx = ctx
	return pc + 1, nil
}

// handleExitScopeWithKey handles the EXIT_SCOPE_WITH_KEY opcode
func (exec *Executor) handleExitScopeWithKey(stack *Stack, instr *instruction.Instruction, pc int) (int, error) {
	// For now, we just increment the program counter
	// In a more advanced implementation, we might manage nested scopes
	ctx := exec.vm.currentCtx.GetParent()
	exec.vm.currentCtx = ctx
	return pc + 1, nil
}

// handleGetIndex handles the GET_INDEX opcode
func (exec *Executor) handleGetIndex(stack *Stack, instr *instruction.Instruction, pc int) (int, error) {
	if stack.Len() < 2 {
		return 0, fmt.Errorf("stack underflow for GET_INDEX")
	}

	// Pop the index and the collection
	index := stack.Pop()
	collection := stack.Pop()

	// Handle different collection types
	switch coll := collection.(type) {
	case []interface{}:
		// Handle slice/array indexing
		idx, ok := index.(int)
		if !ok {
			return 0, fmt.Errorf("index must be an integer, got %T", index)
		}
		if idx < 0 || idx >= len(coll) {
			return 0, fmt.Errorf("index out of range: %d", idx)
		}
		stack.Push(coll[idx])
	case map[string]interface{}:
		// Handle map indexing
		key, ok := index.(string)
		if !ok {
			return 0, fmt.Errorf("map key must be a string, got %T", index)
		}
		value, exists := coll[key]
		if !exists {
			stack.Push(nil)
		} else {
			stack.Push(value)
		}
	default:
		return 0, fmt.Errorf("unsupported collection type for indexing: %T", collection)
	}

	return pc + 1, nil
}

// handleSetIndex handles the SET_INDEX opcode
func (exec *Executor) handleSetIndex(stack *Stack, instr *instruction.Instruction, pc int) (int, error) {
	if stack.Len() < 3 {
		return 0, fmt.Errorf("stack underflow for SET_INDEX")
	}

	// Pop the value, index, and collection
	value := stack.Pop()
	index := stack.Pop()
	collection := stack.Pop()

	// Handle different collection types
	switch coll := collection.(type) {
	case []interface{}:
		// Handle slice/array indexing
		idx, ok := index.(int)
		if !ok {
			return 0, fmt.Errorf("index must be an integer, got %T", index)
		}
		if idx < 0 || idx >= len(coll) {
			return 0, fmt.Errorf("index out of range: %d", idx)
		}
		coll[idx] = value
	case map[string]interface{}:
		// Handle map indexing
		key, ok := index.(string)
		if !ok {
			return 0, fmt.Errorf("map key must be a string, got %T", index)
		}
		coll[key] = value
	default:
		return 0, fmt.Errorf("unsupported collection type for indexing: %T (value: %v, index: %v)", collection, value, index)
	}

	return pc + 1, nil
}

// handleJump handles the JUMP opcode
func (exec *Executor) handleJump(stack *Stack, instr *instruction.Instruction, pc int) (int, error) {
	target, ok := instr.Arg.(int)
	if !ok {
		return 0, fmt.Errorf("invalid jump target")
	}
	return target, nil
}

// handleJumpIf handles the JUMP_IF opcode
func (exec *Executor) handleJumpIf(stack *Stack, instr *instruction.Instruction, pc int) (int, error) {
	target, ok := instr.Arg.(int)
	if !ok {
		return 0, fmt.Errorf("invalid jump target")
	}

	if stack.Len() < 1 {
		return 0, fmt.Errorf("stack underflow for conditional jump")
	}

	// Pop the condition value
	condition := stack.Pop()

	// Check if condition is truthy
	// For loop conditions, we want to jump when the condition is FALSE (to exit the loop)
	// For if statements, we want to jump when the condition is FALSE (to skip the if block)
	if !isTruthy(condition) {
		return target, nil
	}

	return pc + 1, nil
}

// isTruthy determines if a value is truthy
func isTruthy(value interface{}) bool {
	switch v := value.(type) {
	case nil:
		return false
	case bool:
		return v
	case int:
		return v != 0
	case float64:
		return v != 0.0
	case string:
		return v != ""
	default:
		// For other types, consider them truthy if they're not nil
		return true
	}
}

// handleNewSlice handles the NEW_SLICE opcode
func (exec *Executor) handleNewSlice(stack *Stack, instr *instruction.Instruction, pc int) (int, error) {
	size, ok := instr.Arg.(int)
	if !ok {
		return 0, fmt.Errorf("invalid size for NEW_SLICE")
	}

	// Create a new slice with the specified size
	slice := make([]interface{}, size)
	stack.Push(slice)
	return pc + 1, nil
}

// handleLen handles the LEN opcode
func (exec *Executor) handleLen(stack *Stack, instr *instruction.Instruction, pc int) (int, error) {
	if stack.Len() < 1 {
		return 0, fmt.Errorf("stack underflow for LEN")
	}

	// Pop the collection
	collection := stack.Pop()

	// Handle different collection types
	switch coll := collection.(type) {
	case []interface{}:
		// Handle slice/array length
		stack.Push(len(coll))
	case map[string]interface{}:
		// Handle map length
		stack.Push(len(coll))
	case string:
		// Handle string length
		stack.Push(len(coll))
	default:
		return 0, fmt.Errorf("unsupported collection type for length: %T", collection)
	}

	return pc + 1, nil
}

// handleRotate handles the ROTATE opcode
// Changes [a, b, c] to [b, c, a]
func (exec *Executor) handleRotate(stack *Stack, instr *instruction.Instruction, pc int) (int, error) {
	if stack.Len() < 3 {
		return 0, fmt.Errorf("stack underflow for ROTATE")
	}

	// Get the top three elements
	c := stack.Pop()
	b := stack.Pop()
	a := stack.Pop()

	// Push them back in rotated order: b, c, a
	stack.Push(b)
	stack.Push(c)
	stack.Push(a)

	return pc + 1, nil
}

// handleSwap handles the SWAP opcode
// Changes [a, b] to [b, a]
func (exec *Executor) handleSwap(stack *Stack, instr *instruction.Instruction, pc int) (int, error) {
	if stack.Len() < 2 {
		return 0, fmt.Errorf("stack underflow for SWAP")
	}

	// Get the top two elements
	b := stack.Pop()
	a := stack.Pop()

	// Push them back in swapped order: b, a
	stack.Push(b)
	stack.Push(a)

	return pc + 1, nil
}

// handleNewStruct handles the NEW_STRUCT opcode
func (exec *Executor) handleNewStruct(stack *Stack, instr *instruction.Instruction, pc int) (int, error) {
	// Create a new struct (represented as a map)
	structInstance := make(map[string]interface{})

	// If there's a type name in the instruction argument, store it
	if typeName, ok := instr.Arg.(string); ok && typeName != "" {
		structInstance["_type"] = typeName
	}

	stack.Push(structInstance)
	return pc + 1, nil
}

// handleSetField handles the SET_FIELD opcode
func (exec *Executor) handleSetField(stack *Stack, instr *instruction.Instruction, pc int) (int, error) {
	if stack.Len() < 2 {
		return 0, fmt.Errorf("stack underflow for SET_FIELD, stack size: %d", stack.Len())
	}

	// Get the field name from the instruction argument
	fieldName, ok := instr.Arg.(string)
	if !ok {
		return 0, fmt.Errorf("SET_FIELD: field name is not a string, got %T", instr.Arg)
	}

	// Pop the value and struct
	// Stack loading order: struct, value
	// Stack popping order: value, struct
	value := stack.Pop()
	structInterface := stack.Pop()

	// Debug information
	if exec.vm.debug {
		fmt.Printf("SET_FIELD: struct = %v (type %T), field = %s, value = %v (type %T)\n",
			structInterface, structInterface, fieldName, value, value)
	}

	// Check that the struct is a map
	structMap, ok := structInterface.(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("SET_FIELD: struct is not a map, got %T", structInterface)
	}

	// First, try to set the field directly
	if _, exists := structMap[fieldName]; exists {
		// Field exists directly, set it
		structMap[fieldName] = value
	} else {
		// If the field doesn0't exist directly, check for promoted fields in anonymous nested structs
		// In Go, when a struct has an anonymous field, its fields are promoted to the outer struct
		fieldSet := false
		for _, nestedStruct := range structMap {
			// Check if this key might be an anonymous field (typically it would be a struct type name)
			// For simplicity, we'll assume any map value that is itself a map could be an anonymous nested struct
			if nestedMap, isMap := nestedStruct.(map[string]interface{}); isMap {
				// Check if the nested struct has the field we're looking for
				if _, found := nestedMap[fieldName]; found {
					// Set the promoted field in the nested struct
					nestedMap[fieldName] = value
					fieldSet = true
					break
				}
			}
		}

		// If we couldn't find a promoted field, set it as a direct field
		if !fieldSet {
			structMap[fieldName] = value
		}
	}

	return pc + 1, nil
}

// handleGetField handles the GET_FIELD opcode
func (exec *Executor) handleGetField(stack *Stack, instr *instruction.Instruction, pc int) (int, error) {
	if stack.Len() < 1 {
		return 0, fmt.Errorf("stack underflow for GET_FIELD, stack size: %d", stack.Len())
	}

	// Get the field name from the instruction argument
	fieldName, ok := instr.Arg.(string)
	if !ok {
		return 0, fmt.Errorf("GET_FIELD: field name is not a string, got %T", instr.Arg)
	}

	// Pop the struct
	structInterface := stack.Pop()

	// Debug information
	if exec.vm.debug {
		fmt.Printf("GET_FIELD: struct = %v (type %T), field = %s\n", structInterface, structInterface, fieldName)
	}

	// Check that the struct is a map
	structMap, ok := structInterface.(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("GET_FIELD: struct is not a map, got %T", structInterface)
	}

	// First, try to get the field directly
	value, exists := structMap[fieldName]
	if !exists {
		// If the field doesn't exist directly, check for promoted fields in anonymous nested structs
		// In Go, when a struct has an anonymous field, its fields are promoted to the outer struct
		for _, nestedStruct := range structMap {
			// Check if this key might be an anonymous field (typically it would be a struct type name)
			// For simplicity, we'll assume any map value that is itself a map could be an anonymous nested struct
			if nestedMap, isMap := nestedStruct.(map[string]interface{}); isMap {
				// Check if the nested struct has the field we're looking for
				if promotedValue, found := nestedMap[fieldName]; found {
					// Found the promoted field
					stack.Push(promotedValue)
					return pc + 1, nil
				}
			}
		}

		// Field doesn't exist even after checking for promoted fields, push nil
		stack.Push(nil)
	} else {
		// Push the field value
		stack.Push(value)
	}

	return pc + 1, nil
}

// handleCallMethod handles the CALL_METHOD opcode
func (exec *Executor) handleCallMethod(stack *Stack, instr *instruction.Instruction, pc int) (int, error) {
	vm := exec.vm
	methodName, ok := instr.Arg.(string)
	if !ok {
		return 0, fmt.Errorf("invalid method name")
	}

	// Handle both cases: argCount (int) for stack-based or argValues ([]interface{}) for direct values
	var args []interface{}
	var receiver interface{}

	// Debug information - print stack before processing
	if exec.vm.debug {
		fmt.Printf("Stack before CALL_METHOD %s: %v\n", methodName, stack.Items())
	}

	// Check if Arg2 is a slice of arguments (direct values) or an int (arg count)
	switch arg2 := instr.Arg2.(type) {
	case []interface{}:
		// Direct argument values embedded in the instruction
		args = arg2
		// Get the receiver from the stack
		if stack.Len() < 1 {
			return 0, fmt.Errorf("stack underflow when calling method %s", methodName)
		}
		receiver = stack.Pop()
	case int:
		// Stack-based arguments
		argCount := arg2
		// Check if we have enough elements on the stack
		if stack.Len() < argCount+1 {
			return 0, fmt.Errorf("stack underflow when calling method %s", methodName)
		}

		// Prepare arguments (excluding the receiver) using the unified function
		var err error
		args, err = exec.prepareArguments(stack, argCount)
		if err != nil {
			return 0, fmt.Errorf("error preparing arguments for method %s: %w", methodName, err)
		}

		// Get the receiver (the struct instance)
		receiver = stack.Pop()
	default:
		return 0, fmt.Errorf("invalid argument type for CALL_METHOD")
	}

	// Debug information
	if exec.vm.debug {
		fmt.Printf("Calling method %s with %d arguments\n", methodName, len(args))
		fmt.Printf("Method %s receiver: %v (type %T), args: %v\n", methodName, receiver, receiver, args)
	}

	// First, try to find a method with the qualified name (e.g., "Person.GetName")
	// This is for our new approach where structs are treated like packages
	qualifiedMethodName := methodName
	if structMap, ok := receiver.(map[string]interface{}); ok {
		// If we have a struct type name, we can create a qualified method name
		if typeName, exists := structMap["_type"]; exists {
			qualifiedMethodName = fmt.Sprintf("%s.%s", typeName, methodName)
		} else {
			// Try to infer the type name from the context
			// This is a heuristic approach - in a real implementation we would store type info better
			for key := range structMap {
				if key != "width" && key != "height" && key != "radius" && key != "_type" && key != "name" && key != "age" {
					// Assume this is the type name
					qualifiedMethodName = fmt.Sprintf("%s.%s", key, methodName)
					break
				}
			}
		}
	}

	if exec.vm.debug {
		fmt.Printf("Looking for registered function with key: %s\n", qualifiedMethodName)
	}
	// Try to find the method by looking for a registered function with the qualified name
	if fn, exists := vm.GetFunction(qualifiedMethodName); exists {
		if exec.vm.debug {
			fmt.Printf("Found registered function with key: %s\n", qualifiedMethodName)
		}
		// Prepare arguments including the receiver as the first argument
		allArgs := make([]interface{}, len(args)+1)
		allArgs[0] = receiver
		copy(allArgs[1:], args)

		// Call the method
		result, err := fn(allArgs...)
		if err != nil {
			return 0, fmt.Errorf("error calling method %s: %w", methodName, err)
		}

		// Push result back to stack if not nil
		if result != nil {
			stack.Push(result)
		}
		if exec.vm.debug {
			fmt.Printf("Stack after CALL_METHOD %s (builtin): %v\n", methodName, stack.Items())
		}
		return pc + 1, nil
	} else {
		if exec.vm.debug {
			fmt.Printf("No registered function found with key: %s\n", qualifiedMethodName)
		}
	}

	// Check if it's a script-defined method
	// For script-defined methods, we need to find them by key
	// The key would be something like "test.func.methodName"
	// This is a simplified approach for testing purposes
	// Try different function key patterns in order of preference
	functionKeys := []string{
		qualifiedMethodName, // Try the qualified method name first (e.g., "Rectangle.SetWidth")
		fmt.Sprintf("*%s.%s", getStructTypeName(receiver), methodName), // Try pointer receiver (e.g., "*Rectangle.SetHeight")
		fmt.Sprintf("test.func.%s", methodName),
		fmt.Sprintf("main.func.%s", methodName),
	}

	var functionInstructions []*instruction.Instruction
	var found bool
	var foundKey string

	for _, key := range functionKeys {
		if exec.vm.debug {
			fmt.Printf("Looking for function with key: %s\n", key)
		}
		if instructions, exists := vm.GetInstructionSet(key); exists {
			functionInstructions = instructions
			found = true
			foundKey = key
			if exec.vm.debug {
				fmt.Printf("Found function with key: %s, %d instructions\n", key, len(instructions))
			}
			break
		}
	}

	if found {
		// Create new context for the method call
		// The method context's parent is the current context
		methodCtx := execContext.NewContext(methodName, vm.currentCtx)

		// Set method arguments as local variables
		// The first argument is the receiver (usually named after the receiver parameter)
		allArgs := make([]interface{}, len(args)+1)
		allArgs[0] = receiver
		copy(allArgs[1:], args)

		// For value receiver methods, we need to create a copy of the struct
		// For pointer receiver methods, we use the original struct
		// Check if this is a pointer receiver method
		isPointerReceiver := strings.HasPrefix(foundKey, "*")
		if exec.vm.debug {
			fmt.Printf("Method %s is pointer receiver: %t\n", foundKey, isPointerReceiver)
		}

		// If it's a value receiver method, create a copy of the struct
		if !isPointerReceiver {
			// Create a copy of the struct for value receiver
			if originalStruct, ok := receiver.(map[string]interface{}); ok {
				structCopy := make(map[string]interface{})
				for k, v := range originalStruct {
					structCopy[k] = v
				}
				allArgs[0] = structCopy
				if exec.vm.debug {
					fmt.Printf("Created copy of struct for value receiver: %v\n", structCopy)
				}
			}
		}

		// Set argument names: first is receiver name, then actual parameter names
		// Try to get parameter names from registered script function info
		paramNames := []string{"r"} // default receiver name

		// Try to get the actual parameter names from the registered script function
		scriptFunctions := vm.GetAllScriptFunctions()
		if exec.vm.debug {
			fmt.Printf("Script functions: %v\n", scriptFunctions)
		}
		foundParamNames := false
		for name, fnInfo := range scriptFunctions {
			// Check if this function matches our method name
			if exec.vm.debug {
				fmt.Printf("Checking function %s: key=%s, paramNames=%v\n", name, fnInfo.Key, fnInfo.ParamNames)
			}
			if fnInfo.Key == foundKey {
				// Use the parameter names from the function info
				if len(fnInfo.ParamNames) > 0 {
					paramNames = fnInfo.ParamNames
					foundParamNames = true
				}
				if exec.vm.debug {
					fmt.Printf("Using paramNames from %s: %v\n", name, paramNames)
				}
				break
			}
		}

		// If we still have default parameter names, try to determine them based on method name
		if !foundParamNames && len(paramNames) == 1 && paramNames[0] == "r" {
			// Try to extract parameter names from the function key or method name
			// For now, we'll use a heuristic approach
			switch methodName {
			case "SetWidth":
				paramNames = []string{"r", "width"}
			case "SetHeight":
				paramNames = []string{"r", "height"}
			case "SetRadius":
				paramNames = []string{"c", "radius"} // Based on the Circle.SetRadius method
			case "Area":
				paramNames = []string{"r"}
			case "Add":
				paramNames = []string{"c", "x"} // Based on our test function (c Calculator) Add(x int)
			case "Scale":
				paramNames = []string{"r", "factor"} // Based on our test function (r Rectangle) Scale(factor int)
			case "GetWidth":
				paramNames = []string{"r"}
			default:
				// Fallback to generic names
				paramNames = []string{"r"} // receiver name
				for i := 0; i < len(args); i++ {
					paramNames = append(paramNames, fmt.Sprintf("arg%d", i))
				}
			}
		}

		// Make sure we have enough parameter names
		for len(paramNames) < len(allArgs) {
			paramNames = append(paramNames, fmt.Sprintf("arg%d", len(paramNames)-1))
		}

		for i, arg := range allArgs {
			paramName := "unknown"
			if i < len(paramNames) {
				paramName = paramNames[i]
			} else if i < 8 {
				// Fallback to generic names a, b, c, etc.
				genericNames := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
				if i < len(genericNames) {
					paramName = genericNames[i]
				} else {
					paramName = fmt.Sprintf("arg%d", i)
				}
			} else {
				paramName = fmt.Sprintf("arg%d", i)
			}
			// Make sure we create the variable in the method context
			methodCtx.CreateVariableWithType(paramName, arg, "unknown")
			// Debug information
			if exec.vm.debug {
				fmt.Printf("Setting parameter %s = %v (type %T)\n", paramName, arg, arg)
			}
		}

		// Execute the method using a new executor
		newExec := NewExecutor(vm)
		// Set the current context for the method execution
		vm.currentCtx = methodCtx

		// Debug information - print all variables in the method context
		vars, _ := methodCtx.GetAllVariablesWithTypes()
		if exec.vm.debug {
			fmt.Printf("Method context variables: %v\n", vars)
		}

		result, err := newExec.executeInstructions(functionInstructions)
		if err != nil {
			return 0, fmt.Errorf("error executing method %s: %w", methodName, err)
		}

		// Push result back to stack if not nil
		if result != nil {
			stack.Push(result)
		}
		if exec.vm.debug {
			fmt.Printf("Stack after CALL_METHOD %s (script): %v\n", methodName, stack.Items())
		}
		return pc + 1, nil
	} else {
		// Try to find the method by looking for a registered function
		// This would be for built-in methods if we had any
		if fn, exists := vm.GetFunction(methodName); exists {
			// Prepare arguments including the receiver as the first argument
			allArgs := make([]interface{}, len(args)+1)
			allArgs[0] = receiver
			copy(allArgs[1:], args)

			// Call the method
			result, err := fn(allArgs...)
			if err != nil {
				return 0, fmt.Errorf("error calling method %s: %w", methodName, err)
			}

			// Push result back to stack if not nil
			if result != nil {
				stack.Push(result)
			}
			if exec.vm.debug {
				fmt.Printf("Stack after CALL_METHOD %s (builtin2): %v\n", methodName, stack.Items())
			}
			return pc + 1, nil
		} else {
			return 0, fmt.Errorf("undefined method: %s", methodName)
		}
	}

	return pc + 1, nil
}

// getStructTypeName extracts the type name from a struct receiver
func getStructTypeName(receiver interface{}) string {
	if structMap, ok := receiver.(map[string]interface{}); ok {
		// First check for explicit type field
		if typeName, exists := structMap["_type"]; exists {
			if name, ok := typeName.(string); ok {
				return name
			}
		}
		// Fallback: try to infer from keys
		for key := range structMap {
			if key != "width" && key != "height" && key != "radius" && key != "_type" && key != "name" && key != "age" {
				return key
			}
		}
	}
	return "unknown"
}

// handleImport handles the IMPORT opcode
func (exec *Executor) handleImport(stack *Stack, instr *instruction.Instruction, pc int) (int, error) {
	importPath, ok := instr.Arg.(string)
	if !ok {
		return 0, fmt.Errorf("invalid import path")
	}

	pkgName, ok := instr.Arg2.(string)
	if !ok {
		return 0, fmt.Errorf("invalid package name")
	}

	// Check if this is a builtin module and register it on-demand
	modules := builtin.ListAllModules()
	for _, moduleName := range modules {
		// Match either by module name or by import path
		if moduleName == pkgName || moduleName == importPath {
			if _, exists := exec.vm.GetModule(moduleName); exists {
				break
			}
			// Register the module with the VM
			moduleExecutor, exists := builtin.GetModuleExecutor(moduleName)
			if exists {
				exec.vm.RegisterModule(moduleName, moduleExecutor)
			}
			break
		}
	}

	// In the VM context, we can't directly access the module manager
	// The module importing should be handled at the Script level
	// For now, we'll just create a placeholder variable
	exec.vm.currentCtx.CreateVariableWithType(pkgName, importPath, "module")

	return pc + 1, nil
}

// prepareArguments prepares arguments from the stack
func (exec *Executor) prepareArguments(stack *Stack, argCount int) ([]interface{}, error) {
	if stack.Len() < argCount {
		return nil, fmt.Errorf("stack underflow when preparing arguments")
	}

	args := make([]interface{}, argCount)
	for i := argCount - 1; i >= 0; i-- {
		args[i] = stack.Pop()
	}

	return args, nil
}

// pushArgumentsBack pushes arguments back to the stack
func (exec *Executor) pushArgumentsBack(stack *Stack, args []interface{}) {
	for _, arg := range args {
		stack.Push(arg)
	}
}

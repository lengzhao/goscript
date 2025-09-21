// Package compiler implements the GoScript compiler
// It compiles AST nodes to bytecode instructions
package compiler

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"github.com/lengzhao/goscript/context"
	"github.com/lengzhao/goscript/types"
	"github.com/lengzhao/goscript/vm"
)

// FunctionInfo holds information about a compiled function
type FunctionInfo struct {
	Name         string
	StartIP      int
	EndIP        int
	ParamCount   int
	ParamNames   []string // Store parameter names for use in function body
	ReceiverType string   // Store receiver type ("value" or "pointer")
	// Reference to the ScriptFunction that will be registered at runtime
	ScriptFunction *vm.ScriptFunction
}

// Compiler compiles AST nodes to bytecode
type Compiler struct {
	// Virtual machine to generate instructions for
	vm *vm.VM

	// Execution context for variable and function lookups
	context *context.ExecutionContext

	// Instruction pointer for jumps
	ip int

	// Function definitions to compile after main
	functions []*ast.FuncDecl

	// Compiled functions map
	functionMap map[string]*FunctionInfo

	// Variable type mapping to track variable types during compilation
	variableTypes map[string]string

	// Expression type mapping to track expression types during compilation
	// This is used for method chaining to determine the type of the receiver
	expressionTypes map[ast.Expr]string
}

// NewCompiler creates a new compiler
func NewCompiler(vm *vm.VM, context *context.ExecutionContext) *Compiler {
	return &Compiler{
		vm:              vm,
		context:         context,
		ip:              0,
		functions:       make([]*ast.FuncDecl, 0),
		functionMap:     make(map[string]*FunctionInfo),
		variableTypes:   make(map[string]string),
		expressionTypes: make(map[ast.Expr]string),
	}
}

// Compile compiles an AST file to bytecode
func (c *Compiler) Compile(file *ast.File) error {
	// Collect all function declarations and type declarations
	var typeDecls []*ast.GenDecl
	var mainFunc *ast.FuncDecl
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			if d.Recv != nil {
				// This is a method declaration
				c.functions = append(c.functions, d)
				fmt.Printf("Found method: %s with receiver\n", d.Name.Name)
			} else if d.Name.Name == "main" {
				mainFunc = d
				fmt.Printf("Found main function\n")
			} else {
				c.functions = append(c.functions, d)
				fmt.Printf("Found function: %s\n", d.Name.Name)
			}
		case *ast.GenDecl:
			// Handle type declarations
			if d.Tok == token.TYPE {
				typeDecls = append(typeDecls, d)
				fmt.Printf("Found type declaration\n")
			}
		}
	}

	// Process type declarations first
	for _, decl := range typeDecls {
		err := c.compileTypeDecl(decl)
		if err != nil {
			return err
		}
	}

	// Create function info for all functions first
	for _, fn := range c.functions {
		var funcName string
		var receiverType string // "value" or "pointer"
		if fn.Recv != nil {
			// This is a method declaration
			// Get receiver type name to create a unique method name
			var receiverTypeName string
			if len(fn.Recv.List) > 0 {
				receiver := fn.Recv.List[0]
				// Extract type name from receiver
				switch t := receiver.Type.(type) {
				case *ast.Ident:
					receiverTypeName = t.Name
					receiverType = "value"
				case *ast.StarExpr:
					if ident, ok := t.X.(*ast.Ident); ok {
						receiverTypeName = ident.Name
					}
					receiverType = "pointer"
				}

				// Create unique method name
				funcName = receiverTypeName + "." + fn.Name.Name
				fmt.Printf("Creating method info: %s (receiver type: %s)\n", funcName, receiverType)
			} else {
				funcName = fn.Name.Name
			}
		} else {
			// This is a regular function
			funcName = fn.Name.Name
		}

		// Create function info
		funcInfo := &FunctionInfo{
			Name:         funcName,
			ParamNames:   make([]string, 0),
			ReceiverType: receiverType,
		}

		// Count parameters and collect parameter names
		// For methods, we need to add the receiver as the first parameter
		if fn.Recv != nil && len(fn.Recv.List) > 0 {
			// Add receiver as first parameter
			receiver := fn.Recv.List[0]
			if len(receiver.Names) > 0 {
				funcInfo.ParamCount++
				funcInfo.ParamNames = append(funcInfo.ParamNames, receiver.Names[0].Name)
			} else {
				// Anonymous receiver, use a default name
				funcInfo.ParamCount++
				funcInfo.ParamNames = append(funcInfo.ParamNames, "_receiver")
			}
		}

		if fn.Type.Params != nil {
			// Count all parameter names (handling multiple names per field)
			for _, param := range fn.Type.Params.List {
				for _, name := range param.Names {
					funcInfo.ParamCount++
					funcInfo.ParamNames = append(funcInfo.ParamNames, name.Name)
				}
			}
		}

		// Create ScriptFunction that will be registered at runtime
		funcInfo.ScriptFunction = &vm.ScriptFunction{
			Name:         funcName, // Use the unique function name
			ParamCount:   funcInfo.ParamCount,
			ParamNames:   funcInfo.ParamNames,
			ReceiverType: receiverType,
		}

		// Store function info with the unique function name
		c.functionMap[funcName] = funcInfo
		fmt.Printf("Stored function info: %s\n", funcName)
	}

	// Compile function definitions first (except main)
	// Generate OpRegistFunction instructions for each function
	for _, fn := range c.functions {
		err := c.compileFunctionRegistration(fn)
		if err != nil {
			return err
		}
	}

	// Compile main function
	if mainFunc != nil {
		err := c.compileBlockStmt(mainFunc.Body)
		if err != nil {
			return err
		}
	}

	// Now compile function bodies
	for _, fn := range c.functions {
		err := c.compileFunctionBody(fn)
		if err != nil {
			return err
		}
	}

	return nil
}

// compileTypeDecl compiles a type declaration
func (c *Compiler) compileTypeDecl(decl *ast.GenDecl) error {
	// Process each specification in the declaration
	for _, spec := range decl.Specs {
		if typeSpec, ok := spec.(*ast.TypeSpec); ok {
			// Get the type name
			typeName := typeSpec.Name.Name

			// Process the type based on its structure
			switch t := typeSpec.Type.(type) {
			case *ast.StructType:
				// Handle struct type declaration
				err := c.compileStructType(typeName, t)
				if err != nil {
					return err
				}
			case *ast.InterfaceType:
				// Handle interface type declaration
				err := c.compileInterfaceType(typeName, t)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// compileStructType compiles a struct type declaration
func (c *Compiler) compileStructType(name string, structType *ast.StructType) error {
	// Create a new StructType
	structTypeDef := types.NewStructType(name)

	// Process fields
	if structType.Fields != nil {
		for _, field := range structType.Fields.List {
			// Get field names
			for _, fieldName := range field.Names {
				// Get field type
				fieldTypeName := getTypeName(field.Type)
				fieldType, err := types.GetTypeByName(fieldTypeName)
				if err != nil {
					// For now, we'll just print a warning and use a default type
					fmt.Printf("Warning: Unknown field type %s, using interface{} as default\n", fieldTypeName)
					fieldType = types.NewInterfaceType("")
				}

				// Get field tag if present
				var tag string
				if field.Tag != nil {
					tag = field.Tag.Value
				}

				// Add field to struct type
				structTypeDef.AddField(fieldName.Name, fieldType, tag)
			}
		}
	}

	// Register the struct type in the context
	fmt.Printf("Compiling struct type: %s with fields %v\n", name, structTypeDef.GetFieldNames())

	// Register the struct type in the VM's type system
	c.vm.RegisterType(name, structTypeDef)

	return nil
}

// compileInterfaceType compiles an interface type declaration
func (c *Compiler) compileInterfaceType(name string, interfaceType *ast.InterfaceType) error {
	// Create a new InterfaceType
	interfaceTypeDef := types.NewInterfaceType(name)

	// Process methods and embedded interfaces
	if interfaceType.Methods != nil {
		for _, method := range interfaceType.Methods.List {
			switch {
			case method.Names != nil:
				// This is a method declaration
				for _, methodName := range method.Names {
					// Get method signature
					methodType := method.Type.(*ast.FuncType)

					// Process parameters
					var params []types.IType
					if methodType.Params != nil {
						for _, param := range methodType.Params.List {
							for _, _ = range param.Names {
								// Get parameter type
								paramTypeName := getTypeName(param.Type)
								paramType, err := types.GetTypeByName(paramTypeName)
								if err != nil {
									// For now, we'll just print a warning and use a default type
									fmt.Printf("Warning: Unknown parameter type %s, using interface{} as default\n", paramTypeName)
									paramType = types.NewInterfaceType("")
								}
								params = append(params, paramType)
							}
						}
					}

					// Process return values
					var returns []types.IType
					if methodType.Results != nil {
						for _, result := range methodType.Results.List {
							for _, _ = range result.Names {
								// Get return type
								returnTypeName := getTypeName(result.Type)
								returnType, err := types.GetTypeByName(returnTypeName)
								if err != nil {
									// For now, we'll just print a warning and use a default type
									fmt.Printf("Warning: Unknown return type %s, using interface{} as default\n", returnTypeName)
									returnType = types.NewInterfaceType("")
								}
								returns = append(returns, returnType)
							}
						}
					}

					// Add method to interface type
					interfaceTypeDef.AddMethod(methodName.Name, params, returns)
				}
			default:
				// This is an embedded interface
				embeddedTypeName := getTypeName(method.Type)
				embeddedType, err := types.GetTypeByName(embeddedTypeName)
				if err != nil {
					// For now, we'll just print a warning and use a default type
					fmt.Printf("Warning: Unknown embedded interface type %s, using interface{} as default\n", embeddedTypeName)
					embeddedType = types.NewInterfaceType("")
				}
				interfaceTypeDef.AddEmbedded(embeddedType)
			}
		}
	}

	// Register the interface type in the context
	fmt.Printf("Compiling interface type: %s with methods %v\n", name, getInterfaceMethodNames(interfaceTypeDef))

	// TODO: Register the interface type in the runtime context

	return nil
}

// Helper function to get type name from ast.Expr
func getTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		// For now, we'll just return the selector name
		return t.Sel.Name
	case *ast.StarExpr:
		// Pointer type
		return "*" + getTypeName(t.X)
	case *ast.ArrayType:
		if t.Len == nil {
			// Slice type
			return "[]" + getTypeName(t.Elt)
		} else {
			// Array type
			return "[" + getTypeName(t.Len) + "]" + getTypeName(t.Elt)
		}
	default:
		return "interface{}"
	}
}

// Helper function to get method names from an interface type
func getInterfaceMethodNames(interfaceType *types.InterfaceType) []string {
	methods := make([]string, 0)
	// We can't directly access the methods map, so we'll return a placeholder
	// In a real implementation, we would have a method to get method names
	methods = append(methods, "[methods]")
	return methods
}

// compileMethod compiles a method declaration
func (c *Compiler) compileMethod(fn *ast.FuncDecl) error {
	// Get method name
	methodName := fn.Name.Name

	// Get receiver type name to create a unique method name
	var receiverTypeName string
	var receiverType string // "value" or "pointer"
	if fn.Recv != nil && len(fn.Recv.List) > 0 {
		receiver := fn.Recv.List[0]
		// Extract type name from receiver
		switch t := receiver.Type.(type) {
		case *ast.Ident:
			receiverTypeName = t.Name
			receiverType = "value"
		case *ast.StarExpr:
			if ident, ok := t.X.(*ast.Ident); ok {
				receiverTypeName = ident.Name
			}
			receiverType = "pointer"
		}
	}

	// Create unique method name
	uniqueMethodName := methodName
	if receiverTypeName != "" {
		uniqueMethodName = receiverTypeName + "." + methodName
	}

	fmt.Printf("Compiling method: %s (unique name: %s, receiver type: %s)\n", methodName, uniqueMethodName, receiverType)

	// Create function info
	funcInfo := &FunctionInfo{
		Name:         uniqueMethodName,
		ParamNames:   make([]string, 0),
		ReceiverType: receiverType,
	}

	// Add receiver as first parameter
	if fn.Recv != nil && len(fn.Recv.List) > 0 {
		receiver := fn.Recv.List[0]
		if len(receiver.Names) > 0 {
			// Add receiver name to parameter names
			funcInfo.ParamNames = append(funcInfo.ParamNames, receiver.Names[0].Name)
			funcInfo.ParamCount++
		}
	}

	// Count parameters and collect parameter names
	if fn.Type.Params != nil {
		// Count all parameter names (handling multiple names per field)
		for _, param := range fn.Type.Params.List {
			for _, name := range param.Names {
				funcInfo.ParamCount++
				funcInfo.ParamNames = append(funcInfo.ParamNames, name.Name)
			}
		}
	}

	// Create ScriptFunction that will be registered at runtime
	funcInfo.ScriptFunction = &vm.ScriptFunction{
		Name:         uniqueMethodName,
		ParamCount:   funcInfo.ParamCount,
		ParamNames:   funcInfo.ParamNames,
		ReceiverType: receiverType,
	}

	fmt.Printf("Registering script function: %s with receiver type: %s\n", uniqueMethodName, receiverType)

	// Store function info
	c.functionMap[uniqueMethodName] = funcInfo

	// Generate OpRegistFunction instruction
	c.emitInstruction(vm.NewInstruction(vm.OpRegistFunction, uniqueMethodName, funcInfo.ScriptFunction))

	return nil
}

// compileFunctionRegistration generates OpRegistFunction instruction for a function
func (c *Compiler) compileFunctionRegistration(fn *ast.FuncDecl) error {
	// Get function name
	var funcName string
	if fn.Recv != nil {
		// This is a method declaration
		// Get receiver type name to create a unique method name
		var receiverTypeName string
		if len(fn.Recv.List) > 0 {
			receiver := fn.Recv.List[0]
			// Extract type name from receiver
			switch t := receiver.Type.(type) {
			case *ast.Ident:
				receiverTypeName = t.Name
			case *ast.StarExpr:
				if ident, ok := t.X.(*ast.Ident); ok {
					receiverTypeName = ident.Name
				}
			}

			// Create unique method name
			funcName = receiverTypeName + "." + fn.Name.Name
		} else {
			funcName = fn.Name.Name
		}
	} else {
		// This is a regular function
		funcName = fn.Name.Name
	}

	fmt.Printf("Compiling function registration for: %s\n", funcName)

	// Get function info
	funcInfo, exists := c.functionMap[funcName]
	if !exists {
		return fmt.Errorf("function %s not found in function map", funcName)
	}

	// Generate OpRegistFunction instruction
	c.emitInstruction(vm.NewInstruction(vm.OpRegistFunction, funcName, funcInfo.ScriptFunction))
	fmt.Printf("Emitted OpRegistFunction for: %s\n", funcName)

	return nil
}

// compileFunctionBody compiles a function body
func (c *Compiler) compileFunctionBody(fn *ast.FuncDecl) error {
	// Get function name
	var funcName string
	if fn.Recv != nil {
		// This is a method
		// Get receiver type name to create a unique method name
		var receiverTypeName string
		if len(fn.Recv.List) > 0 {
			receiver := fn.Recv.List[0]
			// Extract type name from receiver
			switch t := receiver.Type.(type) {
			case *ast.Ident:
				receiverTypeName = t.Name
			case *ast.StarExpr:
				if ident, ok := t.X.(*ast.Ident); ok {
					receiverTypeName = ident.Name
				}
			}

			// Create unique method name
			funcName = receiverTypeName + "." + fn.Name.Name
		} else {
			funcName = fn.Name.Name
		}
	} else {
		// This is a regular function
		funcName = fn.Name.Name
	}

	// Get existing function info
	funcInfo, exists := c.functionMap[funcName]
	if !exists {
		return fmt.Errorf("function %s not found in function map", funcName)
	}

	// Set the start IP for the function
	funcInfo.ScriptFunction.StartIP = c.ip

	// Generate code to store parameters as local variables
	// Parameters are pushed onto the stack in reverse order during function call
	// We need to pop them and store them as local variables with their original names
	for i := len(funcInfo.ParamNames) - 1; i >= 0; i-- {
		paramName := funcInfo.ParamNames[i]
		// Pop parameter from stack and store as local variable
		c.emitInstruction(vm.NewInstruction(vm.OpStoreName, paramName, nil))
	}

	// Compile the function body
	err := c.compileBlockStmt(fn.Body)
	if err != nil {
		return err
	}

	// Set the end IP for the function
	funcInfo.ScriptFunction.EndIP = c.ip

	// If function doesn't have an explicit return, add one
	c.emitInstruction(vm.NewInstruction(vm.OpLoadConst, nil, nil))
	c.emitInstruction(vm.NewInstruction(vm.OpReturn, nil, nil))

	return nil
}

// compileBlockStmt compiles a block statement
func (c *Compiler) compileBlockStmt(block *ast.BlockStmt) error {
	// Compile each statement in the block
	for _, stmt := range block.List {
		err := c.compileStmt(stmt)
		if err != nil {
			return err
		}
	}

	return nil
}

// compileStmt compiles a statement
func (c *Compiler) compileStmt(stmt ast.Stmt) error {
	switch s := stmt.(type) {
	case *ast.ExprStmt:
		return c.compileExprStmt(s)
	case *ast.AssignStmt:
		return c.compileAssignStmt(s)
	case *ast.ReturnStmt:
		return c.compileReturnStmt(s)
	case *ast.IfStmt:
		return c.compileIfStmt(s)
	case *ast.ForStmt:
		return c.compileForStmt(s)
	case *ast.RangeStmt:
		return c.compileRangeStmt(s)
	case *ast.IncDecStmt:
		return c.compileIncDecStmt(s)
	default:
		return fmt.Errorf("unsupported statement type: %T", stmt)
	}
}

// compileRangeStmt compiles a range statement
func (c *Compiler) compileRangeStmt(stmt *ast.RangeStmt) error {
	// More complete implementation of range statement
	// This will iterate over the elements of a slice

	// Compile the expression being ranged over (the slice)
	err := c.compileExpr(stmt.X)
	if err != nil {
		return err
	}

	// Store the slice in a temporary variable
	c.emitInstruction(vm.NewInstruction(vm.OpStoreName, "_range_slice", nil))

	// Load the slice again
	c.emitInstruction(vm.NewInstruction(vm.OpLoadName, "_range_slice", nil))

	// Get the length of the slice
	c.emitInstruction(vm.NewInstruction(vm.OpLen, nil, nil))

	// Store the length in a temporary variable
	c.emitInstruction(vm.NewInstruction(vm.OpStoreName, "_range_length", nil))

	// Initialize counter
	c.emitInstruction(vm.NewInstruction(vm.OpLoadConst, 0, nil)) // counter = 0
	c.emitInstruction(vm.NewInstruction(vm.OpStoreName, "_range_counter", nil))

	// Start of loop
	startIP := c.ip

	// Check counter < length
	c.emitInstruction(vm.NewInstruction(vm.OpLoadName, "_range_counter", nil))
	c.emitInstruction(vm.NewInstruction(vm.OpLoadName, "_range_length", nil))
	c.emitInstruction(vm.NewInstruction(vm.OpBinaryOp, vm.OpLess, nil))

	// Jump if condition is false (counter >= length)
	jumpIfInstr := vm.NewInstruction(vm.OpJumpIf, 0, nil) // Placeholder target
	c.emitInstruction(jumpIfInstr)

	// If there's a key variable, assign the current counter value to it
	// But skip if the key is the blank identifier "_"
	if stmt.Key != nil {
		if ident, ok := stmt.Key.(*ast.Ident); ok && ident.Name != "_" {
			// Load the current counter value
			c.emitInstruction(vm.NewInstruction(vm.OpLoadName, "_range_counter", nil))

			// Store it in the key variable
			c.emitInstruction(vm.NewInstruction(vm.OpStoreName, ident.Name, nil))
		}
	}

	// If there's a value variable, assign the current element value to it
	if stmt.Value != nil {
		// Load the slice
		c.emitInstruction(vm.NewInstruction(vm.OpLoadName, "_range_slice", nil))

		// Load the current counter (index)
		c.emitInstruction(vm.NewInstruction(vm.OpLoadName, "_range_counter", nil))

		// Get the element at the index
		c.emitInstruction(vm.NewInstruction(vm.OpGetElement, nil, nil))

		// Store it in the value variable
		if ident, ok := stmt.Value.(*ast.Ident); ok {
			c.emitInstruction(vm.NewInstruction(vm.OpStoreName, ident.Name, nil))
		}
	}

	// Compile the loop body
	err = c.compileBlockStmt(stmt.Body)
	if err != nil {
		return err
	}

	// Increment counter
	c.emitInstruction(vm.NewInstruction(vm.OpLoadName, "_range_counter", nil))
	c.emitInstruction(vm.NewInstruction(vm.OpLoadConst, 1, nil))
	c.emitInstruction(vm.NewInstruction(vm.OpBinaryOp, vm.OpAdd, nil))
	c.emitInstruction(vm.NewInstruction(vm.OpStoreName, "_range_counter", nil))

	// Jump back to start
	c.emitInstruction(vm.NewInstruction(vm.OpJump, startIP, nil))

	// Update jump target to after the loop
	jumpIfInstr.Arg = c.ip

	// Clean up temporary variables
	c.emitInstruction(vm.NewInstruction(vm.OpLoadConst, nil, nil))
	c.emitInstruction(vm.NewInstruction(vm.OpStoreName, "_range_counter", nil))
	c.emitInstruction(vm.NewInstruction(vm.OpLoadConst, nil, nil))
	c.emitInstruction(vm.NewInstruction(vm.OpStoreName, "_range_length", nil))
	c.emitInstruction(vm.NewInstruction(vm.OpLoadConst, nil, nil))
	c.emitInstruction(vm.NewInstruction(vm.OpStoreName, "_range_slice", nil))

	return nil
}

// compileExprStmt compiles an expression statement
func (c *Compiler) compileExprStmt(stmt *ast.ExprStmt) error {
	return c.compileExpr(stmt.X)
}

// compileAssignStmt compiles an assignment statement
func (c *Compiler) compileAssignStmt(stmt *ast.AssignStmt) error {
	// Handle assignment operators
	switch stmt.Tok {
	case token.ASSIGN, token.DEFINE:
		// Compile the right-hand side expression
		err := c.compileExpr(stmt.Rhs[0])
		if err != nil {
			return err
		}

		// Handle the left-hand side
		switch lhs := stmt.Lhs[0].(type) {
		case *ast.Ident:
			// Store the result in the variable
			c.emitInstruction(vm.NewInstruction(vm.OpStoreName, lhs.Name, nil))

			// Track variable type if we can determine it from the RHS
			if compositeLit, ok := stmt.Rhs[0].(*ast.CompositeLit); ok {
				if typeIdent, ok := compositeLit.Type.(*ast.Ident); ok {
					c.variableTypes[lhs.Name] = typeIdent.Name
				}
			}
		case *ast.SelectorExpr:
			// Handle field assignment (e.g., obj.field = value)
			// Compile the object expression first
			err := c.compileExpr(lhs.X)
			if err != nil {
				return err
			}

			// Push the field name as a constant
			c.emitInstruction(vm.NewInstruction(vm.OpLoadConst, lhs.Sel.Name, nil))

			// The value to assign is already on the stack from RHS compilation
			// Emit instruction to set the field (arguments are in reverse order on stack)
			// Stack: [value, object, fieldName] -> OpSetStructField takes: [object, fieldName, value]
			// So we need to rotate the stack
			c.emitInstruction(vm.NewInstruction(vm.OpRotate, nil, nil))
			c.emitInstruction(vm.NewInstruction(vm.OpSetStructField, nil, nil))
		case *ast.IndexExpr:
			// Handle index assignment (e.g., array[index] = value)
			// Compile the array/slice expression first
			err := c.compileExpr(lhs.X)
			if err != nil {
				return err
			}

			// Compile the index expression
			err = c.compileExpr(lhs.Index)
			if err != nil {
				return err
			}

			// The value to assign is already on the stack from RHS compilation
			// At this point, the stack should have: [value, array, index]
			// But OpSetIndex expects: [array, index, value]
			// So we need to rearrange the stack

			// Emit instruction to rotate the top three elements
			// This will change [value, array, index] to [array, index, value]
			c.emitInstruction(vm.NewInstruction(vm.OpRotate, nil, nil))

			// Emit instruction to set the element at the index
			c.emitInstruction(vm.NewInstruction(vm.OpSetIndex, nil, nil))
		default:
			return fmt.Errorf("unsupported assignment target: %T", lhs)
		}
	case token.ADD_ASSIGN:
		// Handle += operator
		// First load the current value of the variable
		switch lhs := stmt.Lhs[0].(type) {
		case *ast.Ident:
			c.emitInstruction(vm.NewInstruction(vm.OpLoadName, lhs.Name, nil))
		case *ast.SelectorExpr:
			// Handle field access (e.g., obj.field)
			// Compile the object expression first
			err := c.compileExpr(lhs.X)
			if err != nil {
				return err
			}

			// Push the field name as a constant
			c.emitInstruction(vm.NewInstruction(vm.OpLoadConst, lhs.Sel.Name, nil))

			// Emit instruction to get the field
			c.emitInstruction(vm.NewInstruction(vm.OpGetField, nil, nil))
		default:
			return fmt.Errorf("unsupported assignment target for +=: %T", lhs)
		}

		// Compile the right-hand side expression
		err := c.compileExpr(stmt.Rhs[0])
		if err != nil {
			return err
		}

		// Add the values
		c.emitInstruction(vm.NewInstruction(vm.OpBinaryOp, vm.OpAdd, nil))

		// Store the result back
		switch lhs := stmt.Lhs[0].(type) {
		case *ast.Ident:
			c.emitInstruction(vm.NewInstruction(vm.OpStoreName, lhs.Name, nil))
		case *ast.SelectorExpr:
			// Handle field assignment (e.g., obj.field += value)
			// Compile the object expression first
			err := c.compileExpr(lhs.X)
			if err != nil {
				return err
			}

			// Push the field name as a constant
			c.emitInstruction(vm.NewInstruction(vm.OpLoadConst, lhs.Sel.Name, nil))

			// The value to assign is already on the stack from the addition operation
			// Emit instruction to set the field
			c.emitInstruction(vm.NewInstruction(vm.OpSetStructField, nil, nil))
		}
	case token.SUB_ASSIGN:
		// Handle -= operator
		// First load the current value of the variable
		switch lhs := stmt.Lhs[0].(type) {
		case *ast.Ident:
			c.emitInstruction(vm.NewInstruction(vm.OpLoadName, lhs.Name, nil))
		case *ast.SelectorExpr:
			// Handle field access (e.g., obj.field)
			// Compile the object expression first
			err := c.compileExpr(lhs.X)
			if err != nil {
				return err
			}

			// Push the field name as a constant
			c.emitInstruction(vm.NewInstruction(vm.OpLoadConst, lhs.Sel.Name, nil))

			// Emit instruction to get the field
			c.emitInstruction(vm.NewInstruction(vm.OpGetField, nil, nil))
		default:
			return fmt.Errorf("unsupported assignment target for -=: %T", lhs)
		}

		// Compile the right-hand side expression
		err := c.compileExpr(stmt.Rhs[0])
		if err != nil {
			return err
		}

		// Subtract the values
		c.emitInstruction(vm.NewInstruction(vm.OpBinaryOp, vm.OpSub, nil))

		// Store the result back
		switch lhs := stmt.Lhs[0].(type) {
		case *ast.Ident:
			c.emitInstruction(vm.NewInstruction(vm.OpStoreName, lhs.Name, nil))
		case *ast.SelectorExpr:
			// Handle field assignment (e.g., obj.field -= value)
			// Compile the object expression first
			err := c.compileExpr(lhs.X)
			if err != nil {
				return err
			}

			// Push the field name as a constant
			c.emitInstruction(vm.NewInstruction(vm.OpLoadConst, lhs.Sel.Name, nil))

			// The value to assign is already on the stack from the subtraction operation
			// Emit instruction to set the field
			c.emitInstruction(vm.NewInstruction(vm.OpSetStructField, nil, nil))
		}
	case token.MUL_ASSIGN:
		// Handle *= operator
		// First load the current value of the variable
		switch lhs := stmt.Lhs[0].(type) {
		case *ast.Ident:
			c.emitInstruction(vm.NewInstruction(vm.OpLoadName, lhs.Name, nil))
		case *ast.SelectorExpr:
			// Handle field access (e.g., obj.field)
			// Compile the object expression first
			err := c.compileExpr(lhs.X)
			if err != nil {
				return err
			}

			// Push the field name as a constant
			c.emitInstruction(vm.NewInstruction(vm.OpLoadConst, lhs.Sel.Name, nil))

			// Emit instruction to get the field
			c.emitInstruction(vm.NewInstruction(vm.OpGetField, nil, nil))
		default:
			return fmt.Errorf("unsupported assignment target for *=: %T", lhs)
		}

		// Compile the right-hand side expression
		err := c.compileExpr(stmt.Rhs[0])
		if err != nil {
			return err
		}

		// Multiply the values
		c.emitInstruction(vm.NewInstruction(vm.OpBinaryOp, vm.OpMul, nil))

		// Store the result back
		switch lhs := stmt.Lhs[0].(type) {
		case *ast.Ident:
			c.emitInstruction(vm.NewInstruction(vm.OpStoreName, lhs.Name, nil))
		case *ast.SelectorExpr:
			// Handle field assignment (e.g., obj.field *= value)
			// Compile the object expression first
			err := c.compileExpr(lhs.X)
			if err != nil {
				return err
			}

			// Push the field name as a constant
			c.emitInstruction(vm.NewInstruction(vm.OpLoadConst, lhs.Sel.Name, nil))

			// The value to assign is already on the stack from the multiplication operation
			// Emit instruction to set the field
			c.emitInstruction(vm.NewInstruction(vm.OpSetStructField, nil, nil))
		}
	case token.QUO_ASSIGN:
		// Handle /= operator
		// First load the current value of the variable
		switch lhs := stmt.Lhs[0].(type) {
		case *ast.Ident:
			c.emitInstruction(vm.NewInstruction(vm.OpLoadName, lhs.Name, nil))
		case *ast.SelectorExpr:
			// Handle field access (e.g., obj.field)
			// Compile the object expression first
			err := c.compileExpr(lhs.X)
			if err != nil {
				return err
			}

			// Push the field name as a constant
			c.emitInstruction(vm.NewInstruction(vm.OpLoadConst, lhs.Sel.Name, nil))

			// Emit instruction to get the field
			c.emitInstruction(vm.NewInstruction(vm.OpGetField, nil, nil))
		default:
			return fmt.Errorf("unsupported assignment target for /=: %T", lhs)
		}

		// Compile the right-hand side expression
		err := c.compileExpr(stmt.Rhs[0])
		if err != nil {
			return err
		}

		// Divide the values
		c.emitInstruction(vm.NewInstruction(vm.OpBinaryOp, vm.OpDiv, nil))

		// Store the result back
		switch lhs := stmt.Lhs[0].(type) {
		case *ast.Ident:
			c.emitInstruction(vm.NewInstruction(vm.OpStoreName, lhs.Name, nil))
		case *ast.SelectorExpr:
			// Handle field assignment (e.g., obj.field /= value)
			// Compile the object expression first
			err := c.compileExpr(lhs.X)
			if err != nil {
				return err
			}

			// Push the field name as a constant
			c.emitInstruction(vm.NewInstruction(vm.OpLoadConst, lhs.Sel.Name, nil))

			// The value to assign is already on the stack from the division operation
			// Emit instruction to set the field
			c.emitInstruction(vm.NewInstruction(vm.OpSetStructField, nil, nil))
		}
	default:
		return fmt.Errorf("unsupported assignment operator: %s", stmt.Tok)
	}

	return nil
}

// compileReturnStmt compiles a return statement
func (c *Compiler) compileReturnStmt(stmt *ast.ReturnStmt) error {
	// Compile the return expression if it exists
	if len(stmt.Results) > 0 {
		err := c.compileExpr(stmt.Results[0])
		if err != nil {
			return err
		}
	} else {
		// If no return value, return nil
		c.emitInstruction(vm.NewInstruction(vm.OpLoadConst, nil, nil))
	}

	// Emit return instruction
	c.emitInstruction(vm.NewInstruction(vm.OpReturn, nil, nil))

	return nil
}

// compileIfStmt compiles an if statement
func (c *Compiler) compileIfStmt(stmt *ast.IfStmt) error {
	// Compile the condition
	err := c.compileExpr(stmt.Cond)
	if err != nil {
		return err
	}

	// Emit a conditional jump instruction (placeholder target)
	// Jump if condition is FALSE (skip the if body)
	jumpIfInstr := vm.NewInstruction(vm.OpJumpIf, 0, nil) // Placeholder target
	c.emitInstruction(jumpIfInstr)

	// Compile the if body
	err = c.compileBlockStmt(stmt.Body)
	if err != nil {
		return err
	}

	// If there's an else block, we need to jump over it at the end of the if block
	var elseJumpInstr *vm.Instruction
	if stmt.Else != nil {
		// Emit an unconditional jump to skip the else block
		elseJumpInstr = vm.NewInstruction(vm.OpJump, 0, nil) // Placeholder target
		c.emitInstruction(elseJumpInstr)
	}

	// Update the conditional jump target to after the if body
	jumpIfInstr.Arg = c.ip

	// Compile the else block if it exists
	if stmt.Else != nil {
		switch elseStmt := stmt.Else.(type) {
		case *ast.BlockStmt:
			err = c.compileBlockStmt(elseStmt)
			if err != nil {
				return err
			}
		case *ast.IfStmt:
			// Handle else if as a nested if statement
			err = c.compileIfStmt(elseStmt)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported else statement type: %T", elseStmt)
		}

		// Update the else jump target to after the else block
		elseJumpInstr.Arg = c.ip
	}

	return nil
}

// compileForStmt compiles a for statement
func (c *Compiler) compileForStmt(stmt *ast.ForStmt) error {
	// Compile the init statement if it exists
	if stmt.Init != nil {
		err := c.compileStmt(stmt.Init)
		if err != nil {
			return err
		}
	}

	// Save the start IP for looping
	startIP := c.ip

	// Compile the condition if it exists
	if stmt.Cond != nil {
		err := c.compileExpr(stmt.Cond)
		if err != nil {
			return err
		}

		// Emit a conditional jump to exit the loop
		jumpIfInstr := vm.NewInstruction(vm.OpJumpIf, 0, nil) // Placeholder target
		c.emitInstruction(jumpIfInstr)

		// Compile the loop body
		err = c.compileBlockStmt(stmt.Body)
		if err != nil {
			return err
		}

		// Compile the post statement if it exists
		if stmt.Post != nil {
			err = c.compileStmt(stmt.Post)
			if err != nil {
				return err
			}
		}

		// Emit an unconditional jump back to the start
		c.emitInstruction(vm.NewInstruction(vm.OpJump, startIP, nil))

		// Update the conditional jump target to after the loop
		jumpIfInstr.Arg = c.ip
	} else {
		// Infinite loop - compile the body
		err := c.compileBlockStmt(stmt.Body)
		if err != nil {
			return err
		}

		// Compile the post statement if it exists
		if stmt.Post != nil {
			err = c.compileStmt(stmt.Post)
			if err != nil {
				return err
			}
		}

		// Emit an unconditional jump back to the start
		c.emitInstruction(vm.NewInstruction(vm.OpJump, startIP, nil))
	}

	return nil
}

// compileExpr compiles an expression
func (c *Compiler) compileExpr(expr ast.Expr) error {
	switch e := expr.(type) {
	case *ast.BasicLit:
		return c.compileBasicLit(e)
	case *ast.BinaryExpr:
		return c.compileBinaryExpr(e)
	case *ast.CallExpr:
		return c.compileCallExpr(e)
	case *ast.Ident:
		return c.compileIdent(e)
	case *ast.ParenExpr:
		return c.compileExpr(e.X)
	case *ast.CompositeLit:
		return c.compileCompositeLit(e)
	case *ast.SelectorExpr:
		return c.compileSelectorExpr(e)
	case *ast.IndexExpr:
		return c.compileIndexExpr(e)
	case *ast.UnaryExpr:
		return c.compileUnaryExpr(e)
	default:
		return fmt.Errorf("unsupported expression type: %T", expr)
	}
}

// compileUnaryExpr compiles a unary expression
func (c *Compiler) compileUnaryExpr(expr *ast.UnaryExpr) error {
	// Compile the operand
	err := c.compileExpr(expr.X)
	if err != nil {
		return err
	}

	// Handle different unary operators
	switch expr.Op {
	case token.AND: // & operator (address of)
		// For now, we'll just leave the operand on the stack as-is
		// In a more complete implementation, we would need to handle pointers properly
		// For our simple case, we'll treat &Person{} as just Person{}
		return nil
	case token.SUB: // - operator (negation)
		c.emitInstruction(vm.NewInstruction(vm.OpUnaryOp, vm.OpNeg, nil))
	case token.NOT: // ! operator (logical not)
		c.emitInstruction(vm.NewInstruction(vm.OpUnaryOp, vm.OpNot, nil))
	default:
		return fmt.Errorf("unsupported unary operator: %s", expr.Op)
	}

	return nil
}

// compileBasicLit compiles a basic literal
func (c *Compiler) compileBasicLit(lit *ast.BasicLit) error {
	switch lit.Kind {
	case token.INT:
		// Parse the integer value
		value, err := strconv.Atoi(lit.Value)
		if err != nil {
			return err
		}
		c.emitInstruction(vm.NewInstruction(vm.OpLoadConst, value, nil))
	case token.FLOAT:
		// Parse the float value
		value, err := strconv.ParseFloat(lit.Value, 64)
		if err != nil {
			return err
		}
		c.emitInstruction(vm.NewInstruction(vm.OpLoadConst, value, nil))
	case token.STRING:
		// Remove quotes from string literal
		value := lit.Value[1 : len(lit.Value)-1]
		c.emitInstruction(vm.NewInstruction(vm.OpLoadConst, value, nil))
	default:
		return fmt.Errorf("unsupported literal kind: %s", lit.Kind)
	}

	return nil
}

// compileBinaryExpr compiles a binary expression
func (c *Compiler) compileBinaryExpr(expr *ast.BinaryExpr) error {
	// Compile left operand
	err := c.compileExpr(expr.X)
	if err != nil {
		return err
	}

	// Compile right operand
	err = c.compileExpr(expr.Y)
	if err != nil {
		return err
	}

	// Emit the appropriate binary operation
	switch expr.Op {
	case token.ADD:
		c.emitInstruction(vm.NewInstruction(vm.OpBinaryOp, vm.OpAdd, nil))
	case token.SUB:
		c.emitInstruction(vm.NewInstruction(vm.OpBinaryOp, vm.OpSub, nil))
	case token.MUL:
		c.emitInstruction(vm.NewInstruction(vm.OpBinaryOp, vm.OpMul, nil))
	case token.QUO:
		c.emitInstruction(vm.NewInstruction(vm.OpBinaryOp, vm.OpDiv, nil))
	case token.REM:
		c.emitInstruction(vm.NewInstruction(vm.OpBinaryOp, vm.OpMod, nil))
	case token.EQL:
		c.emitInstruction(vm.NewInstruction(vm.OpBinaryOp, vm.OpEqual, nil))
	case token.NEQ:
		c.emitInstruction(vm.NewInstruction(vm.OpBinaryOp, vm.OpNotEqual, nil))
	case token.LSS:
		c.emitInstruction(vm.NewInstruction(vm.OpBinaryOp, vm.OpLess, nil))
	case token.LEQ:
		c.emitInstruction(vm.NewInstruction(vm.OpBinaryOp, vm.OpLessEqual, nil))
	case token.GTR:
		c.emitInstruction(vm.NewInstruction(vm.OpBinaryOp, vm.OpGreater, nil))
	case token.GEQ:
		c.emitInstruction(vm.NewInstruction(vm.OpBinaryOp, vm.OpGreaterEqual, nil))
	case token.LAND:
		c.emitInstruction(vm.NewInstruction(vm.OpBinaryOp, vm.OpAnd, nil))
	case token.LOR:
		c.emitInstruction(vm.NewInstruction(vm.OpBinaryOp, vm.OpOr, nil))
	default:
		return fmt.Errorf("unsupported binary operator: %s", expr.Op)
	}

	return nil
}

// compileCallExpr compiles a function call expression
func (c *Compiler) compileCallExpr(expr *ast.CallExpr) error {
	// Handle method calls (e.g., obj.Method())
	if selExpr, ok := expr.Fun.(*ast.SelectorExpr); ok {
		// Get the method name
		methodName := selExpr.Sel.Name

		// Determine the type of the receiver
		var typeName string
		var varName string
		if ident, ok := selExpr.X.(*ast.Ident); ok {
			// Get the type from our variable type mapping
			varName = ident.Name
			typeName = c.variableTypes[varName]
		} else if callExpr, ok := selExpr.X.(*ast.CallExpr); ok {
			// This is a method chain, e.g., obj.Method1().Method2()
			// We need to determine the return type of the previous method call
			// Check if we have tracked the type of the previous expression
			if trackedType, exists := c.expressionTypes[callExpr]; exists {
				typeName = trackedType
			} else {
				// Try to infer the return type from the method signature
				if prevSelExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
					prevMethodName := prevSelExpr.Sel.Name
					// Try to find a method that matches the name and get its return type
					// We'll look for methods that return a pointer to a struct (for chaining)
					for funcName := range c.functionMap {
						if strings.HasSuffix(funcName, "."+prevMethodName) {
							// Extract the type name from the function name
							parts := strings.Split(funcName, ".")
							if len(parts) == 2 {
								// Check if this method returns a pointer to the same type (for chaining)
								// This is a heuristic - in a real implementation we would have proper type information
								typeName = parts[0]
								break
							}
						}
					}
				}
			}
		}

		// Create unique method name
		uniqueMethodName := methodName
		if typeName != "" {
			uniqueMethodName = typeName + "." + methodName
		}

		// Get function info to determine receiver type
		funcInfo, exists := c.functionMap[uniqueMethodName]
		isPointerReceiver := exists && funcInfo != nil && funcInfo.ReceiverType == "pointer"

		// For method calls, we need to compile the receiver first
		// This ensures the correct order when arguments are popped during function call
		err := c.compileExpr(selExpr.X)
		if err != nil {
			return err
		}

		// Compile all arguments
		argCount := len(expr.Args)
		for _, arg := range expr.Args {
			err := c.compileExpr(arg)
			if err != nil {
				return err
			}
		}

		// For method calls, we pass the receiver as the first argument
		// The stack order should be: [receiver, arg1, ..., argN] so that when popped,
		// we get [receiver, arg1, ..., argN]
		c.emitInstruction(vm.NewInstruction(vm.OpCallMethod, uniqueMethodName, argCount+1)) // +1 for receiver

		// For pointer receivers, we need to update the original variable with the modified struct
		// For value receivers, the original variable should remain unchanged
		if isPointerReceiver && varName != "" {
			// If it's a pointer receiver, we need to store the modified struct back to the variable
			c.emitInstruction(vm.NewInstruction(vm.OpStoreName, varName, nil))
		}

		// Track the type of this method call expression for potential method chaining
		// If this method returns a pointer to a struct (for chaining), track that type
		if funcInfo != nil && funcInfo.ReceiverType == "pointer" {
			// Heuristic: assume pointer receiver methods return the same type for chaining
			// Extract type name from the function name (e.g., "Calculator.Add" -> "Calculator")
			parts := strings.Split(uniqueMethodName, ".")
			if len(parts) >= 2 {
				c.expressionTypes[expr] = parts[0]
			}
		}

		return nil
	}

	// Handle regular function calls
	ident, ok := expr.Fun.(*ast.Ident)
	if !ok {
		return fmt.Errorf("unsupported function call type: %T", expr.Fun)
	}

	// Compile all arguments
	argCount := len(expr.Args)
	for _, arg := range expr.Args {
		err := c.compileExpr(arg)
		if err != nil {
			return err
		}
	}

	// Emit the regular function call instruction for external functions
	c.emitInstruction(vm.NewInstruction(vm.OpCall, ident.Name, argCount))

	return nil
}

// compileIdent compiles an identifier
func (c *Compiler) compileIdent(ident *ast.Ident) error {
	// Emit a load name instruction
	c.emitInstruction(vm.NewInstruction(vm.OpLoadName, ident.Name, nil))

	return nil
}

// compileCompositeLit compiles a composite literal (e.g., struct literal, slice literal)
func (c *Compiler) compileCompositeLit(lit *ast.CompositeLit) error {
	// Check if this is a slice/array literal
	if _, ok := lit.Type.(*ast.ArrayType); ok {
		// This is a slice or array literal
		// Compile all elements first
		for _, elt := range lit.Elts {
			err := c.compileExpr(elt)
			if err != nil {
				return err
			}
		}

		// Emit instruction to create a new slice with the compiled elements
		// The elements are on the stack, and we pass the count as an argument
		c.emitInstruction(vm.NewInstruction(vm.OpNewSlice, len(lit.Elts), nil))

		return nil
	}

	// Handle struct literals
	// Get the type name if available
	var typeName string
	if ident, ok := lit.Type.(*ast.Ident); ok {
		typeName = ident.Name
	}

	// Emit instruction to create a new struct
	// Pass the type name as an argument
	c.emitInstruction(vm.NewInstruction(vm.OpNewStruct, typeName, nil))

	// Process each key-value pair in the composite literal
	for _, elt := range lit.Elts {
		switch kv := elt.(type) {
		case *ast.KeyValueExpr:
			// Get the key (field name) first
			if keyIdent, ok := kv.Key.(*ast.Ident); ok {
				// Push the field name
				c.emitInstruction(vm.NewInstruction(vm.OpLoadConst, keyIdent.Name, nil))
			}

			// Compile the value
			err := c.compileExpr(kv.Value)
			if err != nil {
				return err
			}

			// Emit instruction to set the field
			// Stack at this point: [struct, fieldName, value]
			// OpSetStructField expects: [struct, fieldName, value]
			// So the order is already correct, no need to rotate
			c.emitInstruction(vm.NewInstruction(vm.OpSetStructField, nil, nil))
		}
	}

	return nil
}

// compileSelectorExpr compiles a selector expression (e.g., obj.field)
func (c *Compiler) compileSelectorExpr(expr *ast.SelectorExpr) error {
	// Compile the expression on the left side of the selector
	err := c.compileExpr(expr.X)
	if err != nil {
		return err
	}

	// Get the field name
	fieldName := expr.Sel.Name

	// Push the field name onto the stack
	c.emitInstruction(vm.NewInstruction(vm.OpLoadConst, fieldName, nil))

	// Emit instruction to get the field
	c.emitInstruction(vm.NewInstruction(vm.OpGetField, nil, nil))

	return nil
}

// compileIndexExpr compiles an index expression (e.g., array[index])
func (c *Compiler) compileIndexExpr(expr *ast.IndexExpr) error {
	// Compile the array/slice expression
	err := c.compileExpr(expr.X)
	if err != nil {
		return err
	}

	// Compile the index expression
	err = c.compileExpr(expr.Index)
	if err != nil {
		return err
	}

	// Emit instruction to get the element at the index
	c.emitInstruction(vm.NewInstruction(vm.OpGetIndex, nil, nil))

	return nil
}

// compileIncDecStmt compiles an increment/decrement statement
func (c *Compiler) compileIncDecStmt(stmt *ast.IncDecStmt) error {
	// Load the current value of the variable
	switch x := stmt.X.(type) {
	case *ast.Ident:
		c.emitInstruction(vm.NewInstruction(vm.OpLoadName, x.Name, nil))
	default:
		return fmt.Errorf("unsupported increment/decrement target: %T", x)
	}

	// Load constant 1
	c.emitInstruction(vm.NewInstruction(vm.OpLoadConst, 1, nil))

	// Emit the appropriate binary operation
	switch stmt.Tok {
	case token.INC:
		c.emitInstruction(vm.NewInstruction(vm.OpBinaryOp, vm.OpAdd, nil))
	case token.DEC:
		c.emitInstruction(vm.NewInstruction(vm.OpBinaryOp, vm.OpSub, nil))
	default:
		return fmt.Errorf("unsupported increment/decrement operator: %s", stmt.Tok)
	}

	// Store the result back to the variable
	switch x := stmt.X.(type) {
	case *ast.Ident:
		c.emitInstruction(vm.NewInstruction(vm.OpStoreName, x.Name, nil))
	default:
		return fmt.Errorf("unsupported increment/decrement target: %T", x)
	}

	return nil
}

// emitInstruction adds an instruction to the VM and updates the IP
func (c *Compiler) emitInstruction(instr *vm.Instruction) {
	c.vm.AddInstruction(instr)
	c.ip++
}

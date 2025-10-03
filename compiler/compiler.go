// Package compiler implements the GoScript compiler
// It compiles AST nodes to bytecode instructions with key-based instruction management
package compiler

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"github.com/lengzhao/goscript/context"
	"github.com/lengzhao/goscript/instruction"
	"github.com/lengzhao/goscript/vm"
)

// Compiler compiles AST nodes to bytecode with key-based instruction management
type Compiler struct {
	// Virtual machine to generate instructions for
	vm *vm.VM

	// Compile context for organizing instructions during compilation
	compileContext *context.CompileContext

	// Current scope key for generating unique keys
	currentScopeKey string

	// Package name from the AST
	packageName string

	// Counter for generating unique keys
	keyCounter int

	// Instructions for the current scope
	currentInstructions []*instruction.Instruction

	// Imported modules map (package name -> import path)
	importedModules map[string]string

	// Label positions map (label name -> instruction index)
	labelPositions map[string]int
}

// NewCompiler creates a new compiler with key-based instruction management
func NewCompiler(vmInstance *vm.VM) *Compiler {
	// Create a temporary compile context, will be updated when we know the package name
	compileCtx := context.NewCompileContext("main", nil)
	return &Compiler{
		vm:                  vmInstance,
		compileContext:      compileCtx,
		currentScopeKey:     "main",
		packageName:         "main",
		keyCounter:          0,
		currentInstructions: make([]*instruction.Instruction, 0),
		importedModules:     make(map[string]string),
		labelPositions:      make(map[string]int),
	}
}

// Compile compiles an AST file to bytecode with key-based instruction management
func (c *Compiler) Compile(file *ast.File) error {
	// Get package name from AST
	if file.Name != nil {
		c.packageName = file.Name.Name
	}

	// Update compile context with proper package name
	c.compileContext = context.NewCompileContext(c.packageName, nil)
	c.currentScopeKey = c.packageName
	c.currentInstructions = make([]*instruction.Instruction, 0)

	// Process import declarations first
	for _, decl := range file.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.IMPORT {
			if err := c.compileGenDecl(genDecl); err != nil {
				return err
			}
		}
	}

	// Store package-level instructions if any
	if len(c.currentInstructions) > 0 {
		c.compileContext.SetInstructions(c.packageName, c.currentInstructions)
	}

	// Process function declarations
	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			if err := c.compileFunction(fn); err != nil {
				return err
			}
		}
	}

	// Transfer all compiled instructions to the VM
	return c.transferInstructions()
}

// compileGenDecl compiles general declarations (variables, types, etc.)
func (c *Compiler) compileGenDecl(decl *ast.GenDecl) error {
	switch decl.Tok {
	case token.IMPORT:
		// Handle import declarations
		return c.compileImportDecl(decl)
	case token.VAR:
		// Handle variable declarations
		return c.compileVarDecl(decl)
	case token.TYPE:
		// Handle type declarations (structs, etc.)
		return c.compileTypeDecl(decl)
	}
	return nil
}

// compileImportDecl compiles import declarations
func (c *Compiler) compileImportDecl(decl *ast.GenDecl) error {
	for _, spec := range decl.Specs {
		if importSpec, ok := spec.(*ast.ImportSpec); ok {
			// Get the import path (remove quotes)
			path := importSpec.Path.Value
			if len(path) > 2 {
				path = path[1 : len(path)-1] // Remove quotes
			}

			// Get the package name (either explicit or inferred from path)
			var pkgName string
			if importSpec.Name != nil {
				pkgName = importSpec.Name.Name
			} else {
				// Infer package name from path (simplified approach)
				parts := strings.Split(path, "/")
				pkgName = parts[len(parts)-1]
			}

			// Store the imported module
			c.importedModules[pkgName] = path

			// Emit the import instruction
			c.emitInstruction(instruction.NewInstruction(instruction.OpImport, path, pkgName))

			// Also create a variable for the module with "module" type
			// This will allow us to handle module calls uniformly with method calls
			c.emitInstruction(instruction.NewInstruction(instruction.OpCreateVar, pkgName, "module"))
		}
	}
	return nil
}

// compileVarDecl compiles variable declarations
func (c *Compiler) compileVarDecl(decl *ast.GenDecl) error {
	for _, spec := range decl.Specs {
		if valueSpec, ok := spec.(*ast.ValueSpec); ok {
			// Handle each variable in the declaration
			for i, name := range valueSpec.Names {
				// Create the variable
				c.emitInstruction(instruction.NewInstruction(instruction.OpCreateVar, name.Name, nil))

				// If there's an initial value, compile it and assign it
				if i < len(valueSpec.Values) && valueSpec.Values[i] != nil {
					if err := c.compileExpr(valueSpec.Values[i]); err != nil {
						return err
					}
					c.emitInstruction(instruction.NewInstruction(instruction.OpStoreName, name.Name, nil))
				} else {
					// Initialize with nil if no initial value
					c.emitInstruction(instruction.NewInstruction(instruction.OpLoadConst, nil, nil))
					c.emitInstruction(instruction.NewInstruction(instruction.OpStoreName, name.Name, nil))
				}
			}
		}
	}
	return nil
}

// compileTypeDecl compiles type declarations
func (c *Compiler) compileTypeDecl(decl *ast.GenDecl) error {
	// For now, we'll just acknowledge type declarations
	// In a more complete implementation, we would process struct definitions, etc.
	for _, spec := range decl.Specs {
		if typeSpec, ok := spec.(*ast.TypeSpec); ok {
			fmt.Printf("Compiling type declaration: %s\n", typeSpec.Name.Name)
			// TODO: Process struct types and other complex types
		}
	}
	return nil
}

// compileFunction compiles a function declaration
func (c *Compiler) compileFunction(fn *ast.FuncDecl) error {
	// Generate function key
	funcKey := c.generateFunctionKey(fn)

	// Save current state
	prevScopeKey := c.currentScopeKey
	prevInstructions := c.currentInstructions

	// Set new scope key
	c.currentScopeKey = funcKey
	c.currentInstructions = make([]*instruction.Instruction, 0)

	// Collect parameter names
	var paramNames []string

	// Compile receiver parameter if this is a method
	if fn.Recv != nil && len(fn.Recv.List) > 0 {
		// This is a method, compile the receiver parameter
		for _, param := range fn.Recv.List {
			for _, name := range param.Names {
				c.emitInstruction(instruction.NewInstruction(instruction.OpCreateVar, name.Name, nil))
				// Note: We don't load parameter values here because they will be set by VM when calling the function
				// The VM will map the actual arguments to these parameter names
				paramNames = append(paramNames, name.Name)
			}
		}
	}

	// Compile function parameters as local variables
	if fn.Type.Params != nil {
		for _, param := range fn.Type.Params.List {
			// Handle parameters with explicit names
			if len(param.Names) > 0 {
				for _, name := range param.Names {
					c.emitInstruction(instruction.NewInstruction(instruction.OpCreateVar, name.Name, nil))
					// Note: We don't load parameter values here because they will be set by VM when calling the function
					// The VM will map the actual arguments to these parameter names
					paramNames = append(paramNames, name.Name)
				}
			} else {
				// Handle parameters without explicit names (e.g., in simplified syntax where name is in the type field)
				// In GoScript's simplified syntax, the parameter name might be stored in the type field
				if ident, ok := param.Type.(*ast.Ident); ok {
					// The parameter name is stored in the type field
					paramName := ident.Name
					c.emitInstruction(instruction.NewInstruction(instruction.OpCreateVar, paramName, nil))
					paramNames = append(paramNames, paramName)
				}
			}
		}
	}

	// Compile function body
	if err := c.compileBlockStmt(fn.Body); err != nil {
		// Restore previous state
		c.currentScopeKey = prevScopeKey
		c.currentInstructions = prevInstructions
		return err
	}

	// Store instructions in compile context with function key
	if len(c.currentInstructions) > 0 {
		c.compileContext.SetInstructions(funcKey, c.currentInstructions)
	}

	// Restore previous state
	c.currentScopeKey = prevScopeKey
	c.currentInstructions = prevInstructions

	// Register function with VM
	scriptFunc := &vm.ScriptFunctionInfo{
		Name:       fn.Name.Name,
		Key:        funcKey,
		ParamCount: c.getParamCount(fn),
		ParamNames: paramNames,
	}
	c.vm.RegisterScriptFunction(fn.Name.Name, scriptFunc)

	return nil
}

// generateFunctionKey generates a unique key for a function
func (c *Compiler) generateFunctionKey(fn *ast.FuncDecl) string {
	// Check if this is a method (has receiver)
	if fn.Recv != nil && len(fn.Recv.List) > 0 {
		// This is a method, generate key in format "struct.method"
		// Get the receiver type name
		if len(fn.Recv.List) > 0 {
			receiver := fn.Recv.List[0]
			if len(receiver.Names) > 0 {
				// Get receiver type name
				typeName := c.getTypeNameWithPointer(receiver.Type)
				if typeName != "" {
					return fmt.Sprintf("%s.%s", typeName, fn.Name.Name)
				}
			}
		}
	}

	if fn.Name.Name == "main" {
		return fmt.Sprintf("%s.main", c.packageName)
	}
	return fmt.Sprintf("%s.func.%s", c.packageName, fn.Name.Name)
}

// getTypeName extracts the type name from an AST expression
func (c *Compiler) getTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		// Pointer type, get the underlying type
		return c.getTypeName(t.X)
	case *ast.SelectorExpr:
		// Qualified type, e.g., pkg.Type
		if ident, ok := t.X.(*ast.Ident); ok {
			return fmt.Sprintf("%s.%s", ident.Name, t.Sel.Name)
		}
	}
	return ""
}

// getTypeNameWithPointer extracts the type name from an AST expression, including pointer information
func (c *Compiler) getTypeNameWithPointer(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		// Pointer type, get the underlying type and prefix with "*"
		underlyingType := c.getTypeName(t.X)
		return fmt.Sprintf("*%s", underlyingType)
	case *ast.SelectorExpr:
		// Qualified type, e.g., pkg.Type
		if ident, ok := t.X.(*ast.Ident); ok {
			return fmt.Sprintf("%s.%s", ident.Name, t.Sel.Name)
		}
	}
	return ""
}

// getParamCount gets the number of parameters for a function
func (c *Compiler) getParamCount(fn *ast.FuncDecl) int {
	if fn.Type.Params == nil {
		return 0
	}
	count := 0
	for _, param := range fn.Type.Params.List {
		count += len(param.Names)
	}
	return count
}

// compileBlockStmt compiles a block statement with key-based scope management
func (c *Compiler) compileBlockStmt(block *ast.BlockStmt) error {
	// Generate a unique scope key for this block
	scopeKey := c.generateKey("block")

	// Emit instruction to enter the block scope
	c.emitInstruction(instruction.NewInstruction(instruction.OpEnterScopeWithKey, scopeKey, nil))

	// Compile each statement in the block
	for _, stmt := range block.List {
		if err := c.compileStmt(stmt); err != nil {
			return err
		}
	}

	// Emit instruction to exit the block scope
	c.emitInstruction(instruction.NewInstruction(instruction.OpExitScopeWithKey, scopeKey, nil))

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
	case *ast.BlockStmt:
		return c.compileBlockStmt(s)
	case *ast.IncDecStmt:
		return c.compileIncDecStmt(s)
	case *ast.DeclStmt:
		// Handle declaration statements (local variables)
		if genDecl, ok := s.Decl.(*ast.GenDecl); ok {
			return c.compileGenDecl(genDecl)
		}
	case *ast.SwitchStmt:
		return c.compileSwitchStmt(s)
	case *ast.LabeledStmt:
		// Handle labeled statements
		return c.compileLabeledStmt(s)
	case *ast.BranchStmt:
		// Handle branch statements (goto, break, continue, fallthrough)
		return c.compileBranchStmt(s)
	default:
		return fmt.Errorf("unsupported statement type: %T", stmt)
	}
	return nil
}

// compileRangeStmt compiles a range statement
func (c *Compiler) compileRangeStmt(stmt *ast.RangeStmt) error {
	// Generate unique names for loop variables
	rangeVarName := c.generateKey("range_var")
	counterVarName := c.generateKey("range_counter")
	lengthVarName := c.generateKey("range_length")

	// Compile the expression being ranged over
	if err := c.compileExpr(stmt.X); err != nil {
		return err
	}

	// Store the collection in a temporary variable
	c.emitInstruction(instruction.NewInstruction(instruction.OpStoreName, rangeVarName, nil))

	// Get the length of the collection and store it
	c.emitInstruction(instruction.NewInstruction(instruction.OpLoadName, rangeVarName, nil))
	c.emitInstruction(instruction.NewInstruction(instruction.OpLen, nil, nil))
	c.emitInstruction(instruction.NewInstruction(instruction.OpStoreName, lengthVarName, nil))

	// Create loop counter variable (initialized to 0)
	c.emitInstruction(instruction.NewInstruction(instruction.OpCreateVar, counterVarName, nil))
	c.emitInstruction(instruction.NewInstruction(instruction.OpLoadConst, 0, nil))
	c.emitInstruction(instruction.NewInstruction(instruction.OpStoreName, counterVarName, nil))

	// Save the start IP for looping
	startIP := len(c.currentInstructions)

	// Check loop condition: counter < length
	c.emitInstruction(instruction.NewInstruction(instruction.OpLoadName, counterVarName, nil))
	c.emitInstruction(instruction.NewInstruction(instruction.OpLoadName, lengthVarName, nil))
	c.emitInstruction(instruction.NewInstruction(instruction.OpBinaryOp, instruction.OpLess, nil))

	// Emit a conditional jump to exit the loop (when condition is false)
	jumpIfInstr := instruction.NewInstruction(instruction.OpJumpIf, 0, nil) // Placeholder target
	c.emitInstruction(jumpIfInstr)

	// Set up loop variables if needed
	if stmt.Key != nil {
		// For range with key (index)
		if keyIdent, ok := stmt.Key.(*ast.Ident); ok {
			// Set the key variable to the current counter value
			c.emitInstruction(instruction.NewInstruction(instruction.OpCreateVar, keyIdent.Name, nil))
			c.emitInstruction(instruction.NewInstruction(instruction.OpLoadName, counterVarName, nil))
			c.emitInstruction(instruction.NewInstruction(instruction.OpStoreName, keyIdent.Name, nil))
		}
	}

	if stmt.Value != nil {
		// For range with value
		if valueIdent, ok := stmt.Value.(*ast.Ident); ok {
			// Get the value from the collection at the current index
			c.emitInstruction(instruction.NewInstruction(instruction.OpCreateVar, valueIdent.Name, nil))
			c.emitInstruction(instruction.NewInstruction(instruction.OpLoadName, rangeVarName, nil))
			c.emitInstruction(instruction.NewInstruction(instruction.OpLoadName, counterVarName, nil))
			c.emitInstruction(instruction.NewInstruction(instruction.OpGetIndex, nil, nil))
			c.emitInstruction(instruction.NewInstruction(instruction.OpStoreName, valueIdent.Name, nil))
		}
	}

	// Compile the loop body with its own scope
	if err := c.compileBlockStmt(stmt.Body); err != nil {
		return err
	}

	// Increment the counter
	c.emitInstruction(instruction.NewInstruction(instruction.OpLoadName, counterVarName, nil))
	c.emitInstruction(instruction.NewInstruction(instruction.OpLoadConst, 1, nil))
	c.emitInstruction(instruction.NewInstruction(instruction.OpBinaryOp, instruction.OpAdd, nil))
	c.emitInstruction(instruction.NewInstruction(instruction.OpStoreName, counterVarName, nil))

	// Emit an unconditional jump back to the start
	c.emitInstruction(instruction.NewInstruction(instruction.OpJump, startIP, nil))

	// Update the conditional jump target to after the loop
	jumpIfInstr.Arg = len(c.currentInstructions)

	return nil
}

// compileExprStmt compiles an expression statement
func (c *Compiler) compileExprStmt(stmt *ast.ExprStmt) error {
	return c.compileExpr(stmt.X)
}

// compileAssignStmt compiles an assignment statement
func (c *Compiler) compileAssignStmt(stmt *ast.AssignStmt) error {
	// Handle the left-hand side first for index expressions and selector expressions
	switch lhs := stmt.Lhs[0].(type) {
	case *ast.IndexExpr:
		// Handle index assignment (e.g., array[index] = value)
		// For index assignment, we need to compile in a specific order:
		// 1. Compile the collection (e.g., array)
		// 2. Compile the index (e.g., index)
		// 3. Compile the value to assign
		// 4. Emit SET_INDEX instruction

		// Compile the expression being indexed (e.g., array)
		if err := c.compileExpr(lhs.X); err != nil {
			return err
		}

		// Compile the index expression (e.g., index)
		if err := c.compileExpr(lhs.Index); err != nil {
			return err
		}

		// Handle compound assignment operators for index expressions
		if stmt.Tok != token.ASSIGN { // Not a simple assignment
			// For compound assignment, we need to load the current value first
			// Emit GET_INDEX to get the current value
			c.emitInstruction(instruction.NewInstruction(instruction.OpGetIndex, nil, nil))

			// Compile the right-hand side expression
			if err := c.compileExpr(stmt.Rhs[0]); err != nil {
				return err
			}

			// Apply the binary operation
			switch stmt.Tok {
			case token.ADD_ASSIGN:
				c.emitInstruction(instruction.NewInstruction(instruction.OpBinaryOp, instruction.OpAdd, nil))
			case token.SUB_ASSIGN:
				c.emitInstruction(instruction.NewInstruction(instruction.OpBinaryOp, instruction.OpSub, nil))
			case token.MUL_ASSIGN:
				c.emitInstruction(instruction.NewInstruction(instruction.OpBinaryOp, instruction.OpMul, nil))
			case token.QUO_ASSIGN:
				c.emitInstruction(instruction.NewInstruction(instruction.OpBinaryOp, instruction.OpDiv, nil))
			case token.REM_ASSIGN:
				c.emitInstruction(instruction.NewInstruction(instruction.OpBinaryOp, instruction.OpMod, nil))
			default:
				return fmt.Errorf("unsupported compound assignment operator: %s", stmt.Tok)
			}
		} else {
			// Simple assignment
			// Compile the right-hand side expression (the value to assign)
			err := c.compileExpr(stmt.Rhs[0])
			if err != nil {
				return err
			}
		}

		// Emit the SET_INDEX instruction
		c.emitInstruction(instruction.NewInstruction(instruction.OpSetIndex, nil, nil))
		return nil
	case *ast.SelectorExpr:
		// Handle selector assignment (e.g., struct.field = value)
		// For selector assignment, we need to compile in a specific order:
		// 1. Compile the expression being selected (e.g., struct)
		// 2. Compile the value to assign
		// 3. Emit SET_FIELD instruction with field name as argument

		// Handle compound assignment operators for selector expressions
		if stmt.Tok != token.ASSIGN { // Not a simple assignment
			// For compound assignment, we need to:
			// 1. Load the struct
			if err := c.compileExpr(lhs.X); err != nil {
				return err
			}
			// 2. Get the current value
			c.emitInstruction(instruction.NewInstruction(instruction.OpGetField, lhs.Sel.Name, nil))

			// Compile the right-hand side expression
			if err := c.compileExpr(stmt.Rhs[0]); err != nil {
				return err
			}

			// Apply the binary operation
			switch stmt.Tok {
			case token.ADD_ASSIGN:
				c.emitInstruction(instruction.NewInstruction(instruction.OpBinaryOp, instruction.OpAdd, nil))
			case token.SUB_ASSIGN:
				c.emitInstruction(instruction.NewInstruction(instruction.OpBinaryOp, instruction.OpSub, nil))
			case token.MUL_ASSIGN:
				c.emitInstruction(instruction.NewInstruction(instruction.OpBinaryOp, instruction.OpMul, nil))
			case token.QUO_ASSIGN:
				c.emitInstruction(instruction.NewInstruction(instruction.OpBinaryOp, instruction.OpDiv, nil))
			case token.REM_ASSIGN:
				c.emitInstruction(instruction.NewInstruction(instruction.OpBinaryOp, instruction.OpMod, nil))
			default:
				return fmt.Errorf("unsupported compound assignment operator: %s", stmt.Tok)
			}

			// For compound assignment, we need to load the struct again for SET_FIELD
			// The stack at this point is: [new_value]
			// We need to get: [struct, new_value]
			if err := c.compileExpr(lhs.X); err != nil {
				return err
			}
			// Stack is now: [new_value, struct]
			// We need to swap to get: [struct, new_value]
			c.emitInstruction(instruction.NewInstruction(instruction.OpSwap, nil, nil))
		} else {
			// Simple assignment
			// Compile the expression being selected (e.g., struct)
			if err := c.compileExpr(lhs.X); err != nil {
				return err
			}

			// Compile the right-hand side expression (the value to assign)
			err := c.compileExpr(stmt.Rhs[0])
			if err != nil {
				return err
			}
			// The stack order is already correct: [struct, value]
			// No need to swap
		}

		// Emit the SET_FIELD instruction with field name as argument
		c.emitInstruction(instruction.NewInstruction(instruction.OpSetField, lhs.Sel.Name, nil))
		return nil
	}

	// Handle compound assignment operators for regular variables
	if stmt.Tok != token.ASSIGN && stmt.Tok != token.DEFINE { // Not a simple assignment or declaration
		// For compound assignment, we need to load the current value first
		switch lhs := stmt.Lhs[0].(type) {
		case *ast.Ident:
			c.emitInstruction(instruction.NewInstruction(instruction.OpLoadName, lhs.Name, nil))
		default:
			return fmt.Errorf("unsupported assignment target for compound assignment: %T", lhs)
		}

		// Compile the right-hand side expression
		if err := c.compileExpr(stmt.Rhs[0]); err != nil {
			return err
		}

		// Apply the binary operation
		switch stmt.Tok {
		case token.ADD_ASSIGN:
			c.emitInstruction(instruction.NewInstruction(instruction.OpBinaryOp, instruction.OpAdd, nil))
		case token.SUB_ASSIGN:
			c.emitInstruction(instruction.NewInstruction(instruction.OpBinaryOp, instruction.OpSub, nil))
		case token.MUL_ASSIGN:
			c.emitInstruction(instruction.NewInstruction(instruction.OpBinaryOp, instruction.OpMul, nil))
		case token.QUO_ASSIGN:
			c.emitInstruction(instruction.NewInstruction(instruction.OpBinaryOp, instruction.OpDiv, nil))
		case token.REM_ASSIGN:
			c.emitInstruction(instruction.NewInstruction(instruction.OpBinaryOp, instruction.OpMod, nil))
		default:
			return fmt.Errorf("unsupported compound assignment operator: %s", stmt.Tok)
		}
	} else {
		// For regular assignments, compile the right-hand side first
		err := c.compileExpr(stmt.Rhs[0])
		if err != nil {
			return err
		}
	}

	// Handle the left-hand side
	switch lhs := stmt.Lhs[0].(type) {
	case *ast.Ident:
		// For short variable declaration (:=), create the variable first
		if stmt.Tok == token.DEFINE {
			c.emitInstruction(instruction.NewInstruction(instruction.OpCreateVar, lhs.Name, nil))
		}
		// Store the result in the variable
		c.emitInstruction(instruction.NewInstruction(instruction.OpStoreName, lhs.Name, nil))
	default:
		return fmt.Errorf("unsupported assignment target: %T", lhs)
	}
	return nil
}

// compileReturnStmt compiles a return statement
func (c *Compiler) compileReturnStmt(stmt *ast.ReturnStmt) error {
	// If there are return values, compile them
	if len(stmt.Results) > 0 {
		if err := c.compileExpr(stmt.Results[0]); err != nil {
			return err
		}
	} else {
		// If no return value, return nil
		c.emitInstruction(instruction.NewInstruction(instruction.OpLoadConst, nil, nil))
	}

	// Emit return instruction
	c.emitInstruction(instruction.NewInstruction(instruction.OpReturn, nil, nil))
	return nil
}

// compileIfStmt compiles an if statement using goto-based approach
func (c *Compiler) compileIfStmt(stmt *ast.IfStmt) error {
	// Compile the condition
	if err := c.compileExpr(stmt.Cond); err != nil {
		return err
	}

	// Generate labels for true and false branches
	trueLabel := c.generateKey("if_true")
	falseLabel := c.generateKey("if_false")
	endLabel := c.generateKey("if_end")

	// Emit a conditional jump to the false branch if condition is false
	// JUMP_IF jumps when the condition is FALSE, so if condition is false, jump to falseLabel
	c.emitInstruction(instruction.NewInstruction(instruction.OpJumpIf, falseLabel, nil))

	// If we reach here, condition was TRUE, so jump to true branch
	c.emitInstruction(instruction.NewInstruction(instruction.OpJump, trueLabel, nil))

	// False branch (else part)
	c.emitInstruction(instruction.NewInstruction(instruction.OpLabel, falseLabel, nil))
	if stmt.Else != nil {
		// Compile the else block
		if elseStmt, ok := stmt.Else.(*ast.BlockStmt); ok {
			if err := c.compileBlockStmt(elseStmt); err != nil {
				return err
			}
		}
	}
	// Jump to end after executing else block
	c.emitInstruction(instruction.NewInstruction(instruction.OpJump, endLabel, nil))

	// True branch (if body)
	c.emitInstruction(instruction.NewInstruction(instruction.OpLabel, trueLabel, nil))
	if err := c.compileBlockStmt(stmt.Body); err != nil {
		return err
	}
	// Jump to end after executing if body
	c.emitInstruction(instruction.NewInstruction(instruction.OpJump, endLabel, nil))

	// End of if statement
	c.emitInstruction(instruction.NewInstruction(instruction.OpLabel, endLabel, nil))

	return nil
}

// compileForStmt compiles a for statement with key-based block management
func (c *Compiler) compileForStmt(stmt *ast.ForStmt) error {
	// Compile the init statement if it exists
	if stmt.Init != nil {
		if err := c.compileStmt(stmt.Init); err != nil {
			return err
		}
	}

	// Save the start IP for looping
	startIP := len(c.currentInstructions)

	// Compile the condition if it exists
	if stmt.Cond != nil {
		if err := c.compileExpr(stmt.Cond); err != nil {
			return err
		}

		// Emit a conditional jump to exit the loop (when condition is false)
		// We need to jump to the instruction after the loop body when condition is false
		jumpIfInstr := instruction.NewInstruction(instruction.OpJumpIf, 0, nil) // Placeholder target
		c.emitInstruction(jumpIfInstr)

		// Compile the loop body with its own scope
		if err := c.compileBlockStmt(stmt.Body); err != nil {
			return err
		}

		// Compile the post statement if it exists
		if stmt.Post != nil {
			if err := c.compileStmt(stmt.Post); err != nil {
				return err
			}
		}

		// Emit an unconditional jump back to the start
		c.emitInstruction(instruction.NewInstruction(instruction.OpJump, startIP, nil))

		// Update the conditional jump target to after the loop
		// This is where we exit the loop when condition is false
		jumpIfInstr.Arg = len(c.currentInstructions)
	} else {
		// Infinite loop - compile the body with its own scope
		if err := c.compileBlockStmt(stmt.Body); err != nil {
			return err
		}

		// Compile the post statement if it exists
		if stmt.Post != nil {
			if err := c.compileStmt(stmt.Post); err != nil {
				return err
			}
		}

		// Emit an unconditional jump back to the start
		c.emitInstruction(instruction.NewInstruction(instruction.OpJump, startIP, nil))
	}

	return nil
}

// compileIncDecStmt compiles an increment or decrement statement
func (c *Compiler) compileIncDecStmt(stmt *ast.IncDecStmt) error {
	// Load the current value of the variable
	switch x := stmt.X.(type) {
	case *ast.Ident:
		c.emitInstruction(instruction.NewInstruction(instruction.OpLoadName, x.Name, nil))
	default:
		return fmt.Errorf("unsupported increment/decrement target: %T", x)
	}

	c.emitInstruction(instruction.NewInstruction(instruction.OpLoadConst, 1, nil))
	// Emit the appropriate instruction
	if stmt.Tok == token.INC {
		c.emitInstruction(instruction.NewInstruction(instruction.OpBinaryOp, instruction.OpAdd, nil))
	} else if stmt.Tok == token.DEC {
		c.emitInstruction(instruction.NewInstruction(instruction.OpBinaryOp, instruction.OpSub, nil))
	}

	// Store the result back
	switch x := stmt.X.(type) {
	case *ast.Ident:
		c.emitInstruction(instruction.NewInstruction(instruction.OpStoreName, x.Name, nil))
	default:
		return fmt.Errorf("unsupported increment/decrement target: %T", x)
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
	case *ast.IndexExpr:
		return c.compileIndexExpr(e)
	case *ast.CompositeLit:
		return c.compileCompositeLit(e)
	case *ast.KeyValueExpr:
		// For KeyValueExpr, we just need to compile the value
		// The key is handled by the parent CompositeLit
		return c.compileExpr(e.Value)
	case *ast.SelectorExpr:
		return c.compileSelectorExpr(e)
	case *ast.UnaryExpr:
		return c.compileUnaryExpr(e)
	default:
		return fmt.Errorf("unsupported expression type: %T", expr)
	}
}

// compileUnaryExpr compiles a unary expression
func (c *Compiler) compileUnaryExpr(expr *ast.UnaryExpr) error {
	// For now, we only handle the address operator (&)
	if expr.Op == token.AND {
		// For & expression, we just compile the operand
		// In a more complete implementation, we would need to handle pointers properly
		return c.compileExpr(expr.X)
	}

	return fmt.Errorf("unsupported unary operator: %s", expr.Op)
}

// compileCompositeLit compiles a composite literal (e.g., []int{1, 2, 3} or Person{name: "Alice"})
func (c *Compiler) compileCompositeLit(lit *ast.CompositeLit) error {
	// Check if this is a slice literal (no key specified for elements)
	isSlice := len(lit.Elts) > 0
	if isSlice {
		// Check if the first element is not a KeyValueExpr, which indicates a slice
		_, isKeyValue := lit.Elts[0].(*ast.KeyValueExpr)
		isSlice = !isKeyValue
	}

	if isSlice {
		// Handle slice literals like []int{1, 2, 3}
		// Create a new slice with the appropriate size
		c.emitInstruction(instruction.NewInstruction(instruction.OpNewSlice, len(lit.Elts), nil))

		// Store the slice in a temporary variable so we can reference it multiple times
		tempVarName := c.generateKey("slice_lit")
		c.emitInstruction(instruction.NewInstruction(instruction.OpStoreName, tempVarName, nil))

		// Compile each element and add it to the slice
		for i, elem := range lit.Elts {
			// Load the slice reference
			c.emitInstruction(instruction.NewInstruction(instruction.OpLoadName, tempVarName, nil))

			// Compile the index
			c.emitInstruction(instruction.NewInstruction(instruction.OpLoadConst, i, nil))

			// Compile the element value
			if err := c.compileExpr(elem); err != nil {
				return err
			}

			// Set the element in the slice
			// Stack should be: [..., slice, index, value]
			c.emitInstruction(instruction.NewInstruction(instruction.OpSetIndex, nil, nil))
		}

		// Load the final slice onto the stack
		c.emitInstruction(instruction.NewInstruction(instruction.OpLoadName, tempVarName, nil))
	} else {
		// Handle struct literals like Person{name: "Alice"}
		// Create a new struct with type information if available
		var structType string
		if lit.Type != nil {
			// Try to extract type name from the composite literal type
			structType = c.getTypeName(lit.Type)
		}
		c.emitInstruction(instruction.NewInstruction(instruction.OpNewStruct, structType, nil))

		// Store the struct in a temporary variable so we can reference it multiple times
		tempVarName := c.generateKey("composite_lit")
		c.emitInstruction(instruction.NewInstruction(instruction.OpStoreName, tempVarName, nil))

		// Compile each element and add it to the struct
		for _, elem := range lit.Elts {
			// Handle KeyValueExpr (for struct fields) or regular expressions (for slice elements)
			switch e := elem.(type) {
			case *ast.KeyValueExpr:
				// This is a struct field assignment: key: value
				// Load the struct reference
				c.emitInstruction(instruction.NewInstruction(instruction.OpLoadName, tempVarName, nil))

				// Compile the key (field name)
				var fieldName string
				if keyIdent, ok := e.Key.(*ast.Ident); ok {
					fieldName = keyIdent.Name
				} else {
					return fmt.Errorf("unsupported key type in KeyValueExpr: %T", e.Key)
				}

				// Compile the value
				if err := c.compileExpr(e.Value); err != nil {
					return err
				}

				// Set the field in the struct
				// Stack should be: [..., struct, value]
				c.emitInstruction(instruction.NewInstruction(instruction.OpSetField, fieldName, nil))
			default:
				// For slice elements or other types, we would need different handling
				// But for now, let's focus on struct support
				return fmt.Errorf("unsupported composite literal element type: %T", elem)
			}
		}

		// Load the final struct onto the stack
		c.emitInstruction(instruction.NewInstruction(instruction.OpLoadName, tempVarName, nil))
	}

	return nil
}

// compileIndexExpr compiles an index expression (e.g., array[index])
func (c *Compiler) compileIndexExpr(expr *ast.IndexExpr) error {
	// Compile the expression being indexed (e.g., array)
	if err := c.compileExpr(expr.X); err != nil {
		return err
	}

	// Compile the index expression (e.g., index)
	if err := c.compileExpr(expr.Index); err != nil {
		return err
	}

	// Emit the GET_INDEX instruction
	c.emitInstruction(instruction.NewInstruction(instruction.OpGetIndex, nil, nil))

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
		c.emitInstruction(instruction.NewInstruction(instruction.OpLoadConst, value, nil))
	case token.FLOAT:
		// Parse the float value
		value, err := strconv.ParseFloat(lit.Value, 64)
		if err != nil {
			return err
		}
		c.emitInstruction(instruction.NewInstruction(instruction.OpLoadConst, value, nil))
	case token.STRING:
		// Remove quotes from string literal
		value := lit.Value[1 : len(lit.Value)-1]
		c.emitInstruction(instruction.NewInstruction(instruction.OpLoadConst, value, nil))
	default:
		return fmt.Errorf("unsupported literal kind: %s", lit.Kind)
	}
	return nil
}

// compileBinaryExpr compiles a binary expression
func (c *Compiler) compileBinaryExpr(expr *ast.BinaryExpr) error {
	// Compile left operand
	if err := c.compileExpr(expr.X); err != nil {
		return err
	}

	// Compile right operand
	if err := c.compileExpr(expr.Y); err != nil {
		return err
	}

	// Emit the appropriate binary operation
	switch expr.Op {
	case token.ADD:
		c.emitInstruction(instruction.NewInstruction(instruction.OpBinaryOp, instruction.OpAdd, nil))
	case token.SUB:
		c.emitInstruction(instruction.NewInstruction(instruction.OpBinaryOp, instruction.OpSub, nil))
	case token.MUL:
		c.emitInstruction(instruction.NewInstruction(instruction.OpBinaryOp, instruction.OpMul, nil))
	case token.QUO:
		c.emitInstruction(instruction.NewInstruction(instruction.OpBinaryOp, instruction.OpDiv, nil))
	case token.REM:
		c.emitInstruction(instruction.NewInstruction(instruction.OpBinaryOp, instruction.OpMod, nil))
	case token.EQL:
		c.emitInstruction(instruction.NewInstruction(instruction.OpBinaryOp, instruction.OpEqual, nil))
	case token.NEQ:
		c.emitInstruction(instruction.NewInstruction(instruction.OpBinaryOp, instruction.OpNotEqual, nil))
	case token.LSS:
		c.emitInstruction(instruction.NewInstruction(instruction.OpBinaryOp, instruction.OpLess, nil))
	case token.LEQ:
		c.emitInstruction(instruction.NewInstruction(instruction.OpBinaryOp, instruction.OpLessEqual, nil))
	case token.GTR:
		c.emitInstruction(instruction.NewInstruction(instruction.OpBinaryOp, instruction.OpGreater, nil))
	case token.GEQ:
		c.emitInstruction(instruction.NewInstruction(instruction.OpBinaryOp, instruction.OpGreaterEqual, nil))
	case token.LAND: // Logical AND (&&)
		c.emitInstruction(instruction.NewInstruction(instruction.OpBinaryOp, instruction.OpAnd, nil))
	case token.LOR: // Logical OR (||)
		c.emitInstruction(instruction.NewInstruction(instruction.OpBinaryOp, instruction.OpOr, nil))
	default:
		return fmt.Errorf("unsupported binary operator: %s", expr.Op)
	}

	return nil
}

// compileCallExpr compiles a function call expression with key-based calling
func (c *Compiler) compileCallExpr(expr *ast.CallExpr) error {
	// Handle different types of function calls
	switch fun := expr.Fun.(type) {
	case *ast.Ident:
		// Regular function calls (e.g., add(1, 2))
		// Compile all arguments
		argCount := len(expr.Args)
		for _, arg := range expr.Args {
			if err := c.compileExpr(arg); err != nil {
				return err
			}
		}

		// Emit the function call instruction with key-based calling
		c.emitInstruction(instruction.NewInstruction(instruction.OpCall, fun.Name, argCount))
	case *ast.SelectorExpr:
		// Method calls (e.g., p.SetWidth(20)) or module calls (e.g., math.Max(1, 2))
		// For unified handling, we'll compile the receiver and then use OpCall
		// First, compile the receiver (e.g., p or math)
		if err := c.compileExpr(fun.X); err != nil {
			return err
		}

		// Compile all arguments
		argCount := len(expr.Args)
		for _, arg := range expr.Args {
			if err := c.compileExpr(arg); err != nil {
				return err
			}
		}

		// For unified handling, we use the format "receiver.functionName"
		// The receiver will be on the stack as the first argument
		functionName := fun.Sel.Name
		// Emit the function call instruction with the function name only
		// The receiver is already on the stack as the first argument
		c.emitInstruction(instruction.NewInstruction(instruction.OpCall, functionName, argCount+1))
	default:
		return fmt.Errorf("unsupported function call type: %T", expr.Fun)
	}

	return nil
}

// compileIdent compiles an identifier
func (c *Compiler) compileIdent(ident *ast.Ident) error {
	// Emit a load name instruction
	c.emitInstruction(instruction.NewInstruction(instruction.OpLoadName, ident.Name, nil))
	return nil
}

// compileSelectorExpr compiles a selector expression (e.g., p.age)
func (c *Compiler) compileSelectorExpr(expr *ast.SelectorExpr) error {
	// For field access, we need to:
	// 1. Compile the expression being selected (e.g., p)
	// 2. Emit the GET_FIELD instruction with the field name as argument
	// But according to the new architecture, we should emit OpLoadName with the qualified name

	// Compile the expression being selected (e.g., p)
	if err := c.compileExpr(expr.X); err != nil {
		return err
	}

	// Emit the GET_FIELD instruction with the field name as argument
	c.emitInstruction(instruction.NewInstruction(instruction.OpGetField, expr.Sel.Name, nil))

	return nil
}

// generateKey generates a unique key for a code block
func (c *Compiler) generateKey(prefix string) string {
	c.keyCounter++
	return fmt.Sprintf("%s.%s_%d", c.currentScopeKey, prefix, c.keyCounter)
}

// emitInstruction adds an instruction to the current scope
func (c *Compiler) emitInstruction(instr *instruction.Instruction) {
	c.currentInstructions = append(c.currentInstructions, instr)
}

// transferInstructions transfers all compiled instructions from the compile context to the VM
func (c *Compiler) transferInstructions() error {
	// First, resolve label positions for goto instructions
	c.resolveLabelPositions()

	// Transfer instructions from the compile context
	instructions := c.compileContext.GetAllInstructions()

	// Transfer each set of instructions with their keys
	for key, instrs := range instructions {
		fmt.Printf("Transferring instructions for key: %s, count: %d\n", key, len(instrs))

		// Add instruction set with key to the VM
		c.vm.AddInstructionSet(key, instrs)
	}

	return nil
}

// resolveLabelPositions resolves label positions for goto instructions
func (c *Compiler) resolveLabelPositions() {
	// Get all instruction sets
	allInstructions := c.compileContext.GetAllInstructions()

	// Process each instruction set
	for key, instructions := range allInstructions {
		// Create a map of label names to their positions within this instruction set
		labelMap := make(map[string]int)

		// First pass: collect all label positions
		for i, instr := range instructions {
			if instr.Op == instruction.OpLabel {
				if labelName, ok := instr.Arg.(string); ok {
					labelMap[labelName] = i
				}
			}
		}

		// Second pass: resolve goto and jumpif instructions
		for _, instr := range instructions {
			if instr.Op == instruction.OpJump || instr.Op == instruction.OpJumpIf {
				if labelName, ok := instr.Arg.(string); ok {
					if targetPos, exists := labelMap[labelName]; exists {
						// Update the instruction with the actual target position
						instr.Arg = targetPos
					} else {
						// Label not found in current scope, check if it's a forward reference
						// For now, we'll leave it as is and let the VM handle it
						fmt.Printf("Warning: Label '%s' not found in scope '%s'\n", labelName, key)
					}
				}
			}
		}
	}
}

// compileSwitchStmt compiles a switch statement using goto-based approach
func (c *Compiler) compileSwitchStmt(stmt *ast.SwitchStmt) error {
	// Generate a unique scope key for this switch statement
	scopeKey := c.generateKey("switch")

	// Emit instruction to enter the switch scope
	c.emitInstruction(instruction.NewInstruction(instruction.OpEnterScopeWithKey, scopeKey, nil))

	// Compile the switch tag (expression to switch on) and store it in a variable
	var tagVarName string
	if stmt.Tag != nil {
		// Compile the tag expression
		if err := c.compileExpr(stmt.Tag); err != nil {
			return err
		}

		// Store the tag value in a temporary variable
		tagVarName = c.generateKey("switch_tag")
		c.emitInstruction(instruction.NewInstruction(instruction.OpCreateVar, tagVarName, nil))
		c.emitInstruction(instruction.NewInstruction(instruction.OpStoreName, tagVarName, nil))
	}

	// Generate labels for cases
	caseLabels := make([]string, len(stmt.Body.List))
	defaultLabel := ""
	endLabel := c.generateKey("end_switch")

	// First pass: generate labels
	for i, clause := range stmt.Body.List {
		caseClause, ok := clause.(*ast.CaseClause)
		if !ok {
			return fmt.Errorf("unexpected clause type in switch: %T", clause)
		}

		if len(caseClause.List) == 0 {
			// Default case
			defaultLabel = c.generateKey("default")
			caseLabels[i] = defaultLabel
		} else {
			// Regular case
			caseLabels[i] = c.generateKey("case")
		}
	}

	// If no default case, use endLabel as default
	if defaultLabel == "" {
		defaultLabel = endLabel
	}

	// Generate condition checks and jumps
	for i, clause := range stmt.Body.List {
		caseClause, ok := clause.(*ast.CaseClause)
		if !ok {
			return fmt.Errorf("unexpected clause type in switch: %T", clause)
		}

		if len(caseClause.List) == 0 {
			// Default case - no condition check needed
			continue
		} else {
			// Regular case with conditions
			// For each expression in the case list, check if it matches the tag
			for _, expr := range caseClause.List {
				// Load the tag value
				if tagVarName != "" {
					c.emitInstruction(instruction.NewInstruction(instruction.OpLoadName, tagVarName, nil))
				}

				// Compile the case expression
				if err := c.compileExpr(expr); err != nil {
					return err
				}

				// Emit a binary equality operation
				c.emitInstruction(instruction.NewInstruction(instruction.OpBinaryOp, instruction.OpEqual, nil))

				// Emit a conditional jump to the case body if the condition is true
				// Since JUMP_IF jumps when the condition is FALSE, we need to invert our logic.
				// We want to jump to the case body when the condition is TRUE.
				// So we use a trick: we put the jump to the case body AFTER the JUMP_IF instruction.
				// If the condition is TRUE, JUMP_IF won't jump and we'll continue to the GOTO.
				// If the condition is FALSE, JUMP_IF will jump to skip the GOTO.

				// Create a label to skip the goto instruction
				skipGotoLabel := c.generateKey("skip_goto")

				// Jump to skipGotoLabel if condition is FALSE
				c.emitInstruction(instruction.NewInstruction(instruction.OpJumpIf, skipGotoLabel, nil))

				// If we reach here, condition was TRUE, so jump to case body
				c.emitInstruction(instruction.NewInstruction(instruction.OpJump, caseLabels[i], nil))

				// Label to skip the goto instruction
				c.emitInstruction(instruction.NewInstruction(instruction.OpLabel, skipGotoLabel, nil))
			}
		}
	}

	// Jump to default case if no conditions matched
	c.emitInstruction(instruction.NewInstruction(instruction.OpJump, defaultLabel, nil))

	// Process each case clause body
	for i, clause := range stmt.Body.List {
		caseClause, ok := clause.(*ast.CaseClause)
		if !ok {
			return fmt.Errorf("unexpected clause type in switch: %T", clause)
		}

		// Emit label for this case
		c.emitInstruction(instruction.NewInstruction(instruction.OpLabel, caseLabels[i], nil))

		// Compile each statement in the case body
		for _, caseStmt := range caseClause.Body {
			if err := c.compileStmt(caseStmt); err != nil {
				return err
			}
		}

		// Jump to end of switch after executing the case body
		c.emitInstruction(instruction.NewInstruction(instruction.OpJump, endLabel, nil))
	}

	// Emit label for end of switch (this is also the default label if no default case exists)
	c.emitInstruction(instruction.NewInstruction(instruction.OpLabel, endLabel, nil))

	// Emit instruction to exit the switch scope
	c.emitInstruction(instruction.NewInstruction(instruction.OpExitScopeWithKey, scopeKey, nil))

	return nil
}

// compileLabeledStmt compiles a labeled statement
func (c *Compiler) compileLabeledStmt(stmt *ast.LabeledStmt) error {
	// Record the position of this label
	labelName := stmt.Label.Name
	c.labelPositions[labelName] = len(c.currentInstructions)

	// Emit a label instruction
	c.emitInstruction(instruction.NewInstruction(instruction.OpLabel, labelName, nil))

	// Compile the statement that follows the label
	return c.compileStmt(stmt.Stmt)
}

// compileBranchStmt compiles a branch statement (goto, break, continue, fallthrough)
func (c *Compiler) compileBranchStmt(stmt *ast.BranchStmt) error {
	switch stmt.Tok {
	case token.GOTO:
		// Handle goto statement
		if stmt.Label != nil {
			// Emit a goto instruction with the label name
			// The actual target position will be resolved later during linking
			c.emitInstruction(instruction.NewInstruction(instruction.OpJump, stmt.Label.Name, nil))
		} else {
			return fmt.Errorf("goto statement must have a label")
		}
	case token.BREAK:
		// Handle break statement
		c.emitInstruction(instruction.NewInstruction(instruction.OpBreak, nil, nil))
	case token.CONTINUE:
		// For now, we don't support continue, but we could add it later
		return fmt.Errorf("continue statement not yet supported")
	case token.FALLTHROUGH:
		// For now, we don't support fallthrough, but we could add it later
		return fmt.Errorf("fallthrough statement not yet supported")
	default:
		return fmt.Errorf("unsupported branch statement: %s", stmt.Tok)
	}
	return nil
}

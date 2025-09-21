// Package compiler implements the GoScript compiler
// It compiles AST nodes to bytecode instructions
package compiler

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"

	"github.com/lengzhao/goscript/context"
	"github.com/lengzhao/goscript/vm"
)

// FunctionInfo holds information about a compiled function
type FunctionInfo struct {
	Name       string
	StartIP    int
	EndIP      int
	ParamCount int
	ParamNames []string // Store parameter names for use in function body
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
}

// NewCompiler creates a new compiler
func NewCompiler(vm *vm.VM, context *context.ExecutionContext) *Compiler {
	return &Compiler{
		vm:          vm,
		context:     context,
		ip:          0,
		functions:   make([]*ast.FuncDecl, 0),
		functionMap: make(map[string]*FunctionInfo),
	}
}

// Compile compiles an AST file to bytecode
func (c *Compiler) Compile(file *ast.File) error {
	// Collect all function declarations first
	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			c.functions = append(c.functions, fn)
		}
	}

	// Create function info for all functions first
	for _, fn := range c.functions {
		if fn.Name.Name != "main" {
			// Create function info
			funcInfo := &FunctionInfo{
				Name:       fn.Name.Name,
				ParamNames: make([]string, 0),
			}

			// Count parameters and collect parameter names
			if fn.Type.Params != nil {
				funcInfo.ParamCount = 0
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
				Name:       fn.Name.Name,
				ParamCount: funcInfo.ParamCount,
				ParamNames: funcInfo.ParamNames,
			}

			// Store function info
			c.functionMap[fn.Name.Name] = funcInfo
		}
	}

	// Compile function definitions first (except main)
	// Generate OpRegistFunction instructions for each function
	for _, fn := range c.functions {
		if fn.Name.Name != "main" {
			err := c.compileFunctionRegistration(fn)
			if err != nil {
				return err
			}
		}
	}

	// Compile main function
	var mainFunc *ast.FuncDecl
	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok && fn.Name.Name == "main" {
			mainFunc = fn
			break
		}
	}

	if mainFunc != nil {
		err := c.compileBlockStmt(mainFunc.Body)
		if err != nil {
			return err
		}
	}

	// Now compile function bodies
	for _, fn := range c.functions {
		if fn.Name.Name != "main" {
			err := c.compileFunctionBody(fn)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// compileFunctionRegistration generates OpRegistFunction instruction for a function
func (c *Compiler) compileFunctionRegistration(fn *ast.FuncDecl) error {
	// Get function info
	funcInfo, exists := c.functionMap[fn.Name.Name]
	if !exists {
		return fmt.Errorf("function %s not found in function map", fn.Name.Name)
	}

	// Generate OpRegistFunction instruction
	c.emitInstruction(vm.NewInstruction(vm.OpRegistFunction, fn.Name.Name, funcInfo.ScriptFunction))

	return nil
}

// compileFunctionBody compiles a function body
func (c *Compiler) compileFunctionBody(fn *ast.FuncDecl) error {
	// Get existing function info
	funcInfo, exists := c.functionMap[fn.Name.Name]
	if !exists {
		return fmt.Errorf("function %s not found in function map", fn.Name.Name)
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
	case *ast.IncDecStmt:
		return c.compileIncDecStmt(s)
	default:
		return fmt.Errorf("unsupported statement type: %T", stmt)
	}
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
		default:
			return fmt.Errorf("unsupported assignment target: %T", lhs)
		}
	case token.ADD_ASSIGN:
		// Handle += operator
		// First load the current value of the variable
		switch lhs := stmt.Lhs[0].(type) {
		case *ast.Ident:
			c.emitInstruction(vm.NewInstruction(vm.OpLoadName, lhs.Name, nil))
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
		}
	case token.SUB_ASSIGN:
		// Handle -= operator
		// First load the current value of the variable
		switch lhs := stmt.Lhs[0].(type) {
		case *ast.Ident:
			c.emitInstruction(vm.NewInstruction(vm.OpLoadName, lhs.Name, nil))
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
		}
	case token.MUL_ASSIGN:
		// Handle *= operator
		// First load the current value of the variable
		switch lhs := stmt.Lhs[0].(type) {
		case *ast.Ident:
			c.emitInstruction(vm.NewInstruction(vm.OpLoadName, lhs.Name, nil))
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
		}
	case token.QUO_ASSIGN:
		// Handle /= operator
		// First load the current value of the variable
		switch lhs := stmt.Lhs[0].(type) {
		case *ast.Ident:
			c.emitInstruction(vm.NewInstruction(vm.OpLoadName, lhs.Name, nil))
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
	default:
		return fmt.Errorf("unsupported expression type: %T", expr)
	}
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
	// Get the function name
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

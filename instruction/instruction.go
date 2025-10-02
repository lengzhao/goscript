// Package instruction provides instruction definitions for GoScript
package instruction

import (
	"fmt"
)

// OpCode represents an operation code for the virtual machine
type OpCode byte

const (
	// No operation
	OpNop OpCode = iota

	// Load a constant onto the stack
	OpLoadConst

	// Load a variable by name
	OpLoadName

	// Store a value to a variable by name
	OpStoreName

	// Pop a value from the stack (discard it)
	OpPop

	// Call a function
	OpCall

	// Call a struct method
	OpCallMethod

	// Register a script-defined function
	OpRegistFunction

	// Return from function
	OpReturn

	// Unconditional jump
	OpJump

	// Conditional jump
	OpJumpIf

	// Binary operation (add, sub, mul, div, etc.)
	OpBinaryOp

	// Unary operation (neg, not, etc.)
	OpUnaryOp

	// Create a new struct instance
	OpNewStruct

	// Access a field of a struct
	OpGetField

	// Set a field of a struct
	OpSetField

	// Set a field of a struct with explicit stack order
	OpSetStructField

	// Access an element of an array/slice by index
	OpGetIndex

	// Set an element of an array/slice by index
	OpSetIndex

	// Rotate the top three elements on the stack
	// Changes [a, b, c] to [b, c, a]
	OpRotate

	// Swap the top two elements on the stack
	OpSwap

	// Create a new slice
	OpNewSlice

	// Get the length of a slice or array
	OpLen

	// Get an element from a slice or array by index
	OpGetElement

	// Import a module
	OpImport

	// Enter a scope with a specific key
	OpEnterScopeWithKey

	// Exit a scope with a specific key
	OpExitScopeWithKey

	// Create a new variable
	OpCreateVar

	// Break from loop
	OpBreak

	OpCodeLast
)

// String returns the string representation of an OpCode
func (op OpCode) String() string {
	switch op {
	case OpNop:
		return "OpNop"
	case OpLoadConst:
		return "OpLoadConst"
	case OpLoadName:
		return "OpLoadName"
	case OpStoreName:
		return "OpStoreName"
	case OpPop:
		return "OpPop"
	case OpCall:
		return "OpCall"
	case OpCallMethod:
		return "OpCallMethod"
	case OpRegistFunction:
		return "OpRegistFunction"
	case OpReturn:
		return "OpReturn"
	case OpJump:
		return "OpJump"
	case OpJumpIf:
		return "OpJumpIf"
	case OpBinaryOp:
		return "OpBinaryOp"
	case OpUnaryOp:
		return "OpUnaryOp"
	case OpNewStruct:
		return "OpNewStruct"
	case OpGetField:
		return "OpGetField"
	case OpSetField:
		return "OpSetField"
	case OpSetStructField:
		return "OpSetStructField"
	case OpGetIndex:
		return "OpGetIndex"
	case OpSetIndex:
		return "OpSetIndex"
	case OpRotate:
		return "OpRotate"
	case OpSwap:
		return "OpSwap"
	case OpNewSlice:
		return "OpNewSlice"
	case OpLen:
		return "OpLen"
	case OpGetElement:
		return "OpGetElement"
	case OpImport:
		return "OpImport"
	case OpEnterScopeWithKey:
		return "OpEnterScopeWithKey"
	case OpExitScopeWithKey:
		return "OpExitScopeWithKey"
	case OpCreateVar:
		return "OpCreateVar"
	case OpBreak:
		return "OpBreak"
	default:
		return fmt.Sprintf("OpCode(%d)", op)
	}
}

// BinaryOp represents a binary operation
type BinaryOp byte

const (
	OpAdd BinaryOp = iota
	OpSub
	OpMul
	OpDiv
	OpMod
	OpEqual
	OpNotEqual
	OpLess
	OpLessEqual
	OpGreater
	OpGreaterEqual
	OpAnd
	OpOr
)

// UnaryOp represents a unary operation
type UnaryOp byte

const (
	OpNeg UnaryOp = iota
	OpNot
)

// Instruction represents a single VM instruction
type Instruction struct {
	Op   OpCode
	Arg  interface{}
	Arg2 interface{}
}

// NewInstruction creates a new instruction
func NewInstruction(op OpCode, arg interface{}, arg2 ...interface{}) *Instruction {
	instr := &Instruction{
		Op:  op,
		Arg: arg,
	}

	if len(arg2) > 0 {
		instr.Arg2 = arg2[0]
	}

	return instr
}

// String returns the string representation of an instruction
func (i *Instruction) String() string {
	switch i.Op {
	case OpNop:
		return "NOP"
	case OpLoadConst:
		return fmt.Sprintf("LOAD_CONST %v", i.Arg)
	case OpLoadName:
		return fmt.Sprintf("LOAD_NAME %v", i.Arg)
	case OpStoreName:
		return fmt.Sprintf("STORE_NAME %v", i.Arg)
	case OpPop:
		return "POP"
	case OpCall:
		return fmt.Sprintf("CALL %v %v", i.Arg, i.Arg2)
	case OpCallMethod:
		return fmt.Sprintf("CALL_METHOD %v %v", i.Arg, i.Arg2)
	case OpRegistFunction:
		return fmt.Sprintf("REGIST_FUNCTION %v %v", i.Arg, i.Arg2)
	case OpReturn:
		return "RETURN"
	case OpJump:
		return fmt.Sprintf("JUMP %v", i.Arg)
	case OpJumpIf:
		return fmt.Sprintf("JUMP_IF %v", i.Arg)
	case OpBinaryOp:
		return fmt.Sprintf("BINARY_OP %v", i.Arg)
	case OpUnaryOp:
		return fmt.Sprintf("UNARY_OP %v", i.Arg)
	case OpNewStruct:
		return fmt.Sprintf("NEW_STRUCT %v", i.Arg)
	case OpGetField:
		return fmt.Sprintf("GET_FIELD %v", i.Arg)
	case OpSetField:
		return fmt.Sprintf("SET_FIELD %v", i.Arg)
	case OpSetStructField:
		return fmt.Sprintf("SET_STRUCT_FIELD %v", i.Arg)
	case OpGetIndex:
		return fmt.Sprintf("GET_INDEX %v", i.Arg)
	case OpSetIndex:
		return fmt.Sprintf("SET_INDEX %v", i.Arg)
	case OpRotate:
		return fmt.Sprintf("ROTATE %v", i.Arg)
	case OpSwap:
		return "SWAP"
	case OpNewSlice:
		return fmt.Sprintf("NEW_SLICE %v", i.Arg)
	case OpLen:
		return fmt.Sprintf("LEN %v", i.Arg)
	case OpGetElement:
		return fmt.Sprintf("GET_ELEMENT %v", i.Arg)
	case OpImport:
		return fmt.Sprintf("IMPORT %v", i.Arg)
	case OpEnterScopeWithKey:
		return fmt.Sprintf("ENTER_SCOPE_WITH_KEY %v", i.Arg)
	case OpExitScopeWithKey:
		return fmt.Sprintf("EXIT_SCOPE_WITH_KEY %v", i.Arg)
	case OpCreateVar:
		return fmt.Sprintf("CREATE_VAR %v", i.Arg)
	case OpBreak:
		return "BREAK"
	default:
		return fmt.Sprintf("UNKNOWN(%d) %v %v", i.Op, i.Arg, i.Arg2)
	}
}

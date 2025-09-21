// Package ast provides functionality for working with Go AST nodes
// It wraps the standard Go ast package and provides additional utilities
package ast

import (
	"go/ast"
	"go/token"
)

// Node is a wrapper around ast.Node with additional functionality
type Node struct {
	ast.Node
}

// NewNode creates a new Node wrapper
func NewNode(n ast.Node) *Node {
	return &Node{Node: n}
}

// NodeType returns the type of the node as a string
func (n *Node) NodeType() string {
	switch n.Node.(type) {
	case *ast.File:
		return "File"
	case *ast.FuncDecl:
		return "FuncDecl"
	case *ast.GenDecl:
		return "GenDecl"
	case *ast.AssignStmt:
		return "AssignStmt"
	case *ast.ExprStmt:
		return "ExprStmt"
	case *ast.CallExpr:
		return "CallExpr"
	case *ast.BinaryExpr:
		return "BinaryExpr"
	case *ast.Ident:
		return "Ident"
	case *ast.BasicLit:
		return "BasicLit"
	default:
		return "Unknown"
	}
}

// Position returns the position of the node in the source code
func (n *Node) Position(fset *token.FileSet) token.Position {
	return fset.Position(n.Node.Pos())
}

// Walk traverses the AST rooted at node n in depth-first order
func (n *Node) Walk(fn func(*Node) bool) {
	ast.Inspect(n.Node, func(node ast.Node) bool {
		if node != nil {
			wrapped := NewNode(node)
			return fn(wrapped)
		}
		return true
	})
}

// File represents a parsed Go source file
type File struct {
	*ast.File
}

// NewFile creates a new File wrapper
func NewFile(f *ast.File) *File {
	return &File{File: f}
}

// FunctionDecl represents a function declaration
type FunctionDecl struct {
	*ast.FuncDecl
}

// NewFunctionDecl creates a new FunctionDecl wrapper
func NewFunctionDecl(fd *ast.FuncDecl) *FunctionDecl {
	return &FunctionDecl{FuncDecl: fd}
}

// Name returns the name of the function
func (f *FunctionDecl) Name() string {
	return f.FuncDecl.Name.Name
}

// Body returns the body of the function
func (f *FunctionDecl) Body() *ast.BlockStmt {
	return f.FuncDecl.Body
}

// Expression represents an expression node
type Expression struct {
	ast.Expr
}

// NewExpression creates a new Expression wrapper
func NewExpression(e ast.Expr) *Expression {
	return &Expression{Expr: e}
}

// BinaryExpression represents a binary expression
type BinaryExpression struct {
	*ast.BinaryExpr
}

// NewBinaryExpression creates a new BinaryExpression wrapper
func NewBinaryExpression(be *ast.BinaryExpr) *BinaryExpression {
	return &BinaryExpression{BinaryExpr: be}
}

// Operator returns the operator of the binary expression
func (be *BinaryExpression) Operator() string {
	return be.BinaryExpr.Op.String()
}

// LeftOperand returns the left operand of the binary expression
func (be *BinaryExpression) LeftOperand() ast.Expr {
	return be.BinaryExpr.X
}

// RightOperand returns the right operand of the binary expression
func (be *BinaryExpression) RightOperand() ast.Expr {
	return be.BinaryExpr.Y
}

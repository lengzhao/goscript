package ast

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestNode(t *testing.T) {
	// Parse a simple Go source
	src := `package main

func add(a, b int) int {
	return a + b
}`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		t.Fatalf("Failed to parse source: %v", err)
	}

	// Create a wrapped node
	node := NewNode(file)

	// Test node type
	if node.NodeType() != "File" {
		t.Errorf("Expected node type 'File', got '%s'", node.NodeType())
	}

	// Test position
	pos := node.Position(fset)
	if pos.Line != 1 {
		t.Errorf("Expected position line 1, got %d", pos.Line)
	}

	// Test walk function
	count := 0
	node.Walk(func(n *Node) bool {
		count++
		return true
	})

	if count < 5 {
		t.Errorf("Expected at least 5 nodes, got %d", count)
	}
}

func TestFile(t *testing.T) {
	src := `package main

func add(a, b int) int {
	return a + b
}`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		t.Fatalf("Failed to parse source: %v", err)
	}

	wrappedFile := NewFile(file)

	if wrappedFile.Name.Name != "main" {
		t.Errorf("Expected package name 'main', got '%s'", wrappedFile.Name.Name)
	}
}

func TestFunctionDecl(t *testing.T) {
	src := `package main

func add(a, b int) int {
	return a + b
}`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		t.Fatalf("Failed to parse source: %v", err)
	}

	// Find the function declaration
	var funcDecl *ast.FuncDecl
	ast.Inspect(file, func(n ast.Node) bool {
		if fd, ok := n.(*ast.FuncDecl); ok {
			funcDecl = fd
			return false
		}
		return true
	})

	if funcDecl == nil {
		t.Fatalf("Failed to find function declaration")
	}

	wrappedFunc := NewFunctionDecl(funcDecl)

	if wrappedFunc.Name() != "add" {
		t.Errorf("Expected function name 'add', got '%s'", wrappedFunc.Name())
	}

	if wrappedFunc.Body() == nil {
		t.Error("Expected function body, got nil")
	}
}

func TestBinaryExpression(t *testing.T) {
	src := "x + y"

	expr, err := parser.ParseExpr(src)
	if err != nil {
		t.Fatalf("Failed to parse expression: %v", err)
	}

	binExpr, ok := expr.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("Expected BinaryExpr, got %T", expr)
	}

	wrappedBinExpr := NewBinaryExpression(binExpr)

	if wrappedBinExpr.Operator() != "+" {
		t.Errorf("Expected operator '+', got '%s'", wrappedBinExpr.Operator())
	}

	left, ok := wrappedBinExpr.LeftOperand().(*ast.Ident)
	if !ok {
		t.Fatalf("Expected left operand to be Ident, got %T", wrappedBinExpr.LeftOperand())
	}

	if left.Name != "x" {
		t.Errorf("Expected left operand name 'x', got '%s'", left.Name)
	}

	right, ok := wrappedBinExpr.RightOperand().(*ast.Ident)
	if !ok {
		t.Fatalf("Expected right operand to be Ident, got %T", wrappedBinExpr.RightOperand())
	}

	if right.Name != "y" {
		t.Errorf("Expected right operand name 'y', got '%s'", right.Name)
	}
}

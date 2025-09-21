package parser

import (
	"go/ast"
	"go/parser"
	"testing"
)

func TestParser(t *testing.T) {
	input := `package main

import "fmt"

func main() {
	x := 10
	y := 20
	fmt.Println(x + y)
}`

	p := New()
	file, err := p.Parse("test.go", []byte(input), parser.ParseComments)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if file == nil {
		t.Fatalf("Parse returned nil file")
	}

	if file.Name.Name != "main" {
		t.Errorf("Expected package name 'main', got '%s'", file.Name.Name)
	}

	// Check that we have one function declaration
	var funcDecl *ast.FuncDecl
	ast.Inspect(file, func(n ast.Node) bool {
		if fd, ok := n.(*ast.FuncDecl); ok && fd.Name.Name == "main" {
			funcDecl = fd
			return false // stop inspection
		}
		return true
	})

	if funcDecl == nil {
		t.Fatalf("Expected to find main function, but didn't")
	}

	// Check function body
	if funcDecl.Body == nil {
		t.Fatalf("Function body is nil")
	}

	// Check that we have statements in the function body
	// We expect 3 statements: 2 assignments and 1 function call
	if len(funcDecl.Body.List) != 3 {
		t.Errorf("Expected 3 statements, got %d", len(funcDecl.Body.List))
		for i, stmt := range funcDecl.Body.List {
			t.Logf("Statement %d: %T", i, stmt)
		}
	}
}

func TestParseExpr(t *testing.T) {
	input := "x + y * z"
	p := New()
	expr, err := p.ParseExpr([]byte(input))
	if err != nil {
		t.Fatalf("ParseExpr failed: %v", err)
	}

	if expr == nil {
		t.Fatalf("ParseExpr returned nil")
	}

	// Check that it's a binary expression
	be, ok := expr.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("Expected BinaryExpr, got %T", expr)
	}

	// Check operator
	if be.Op.String() != "+" {
		t.Errorf("Expected operator '+', got '%s'", be.Op.String())
	}

	// Check left operand is an identifier
	if _, ok := be.X.(*ast.Ident); !ok {
		t.Errorf("Expected left operand to be identifier, got %T", be.X)
	}

	// Check right operand is a binary expression (y * z)
	right, ok := be.Y.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("Expected right operand to be BinaryExpr, got %T", be.Y)
	}

	if right.Op.String() != "*" {
		t.Errorf("Expected operator '*' in right operand, got '%s'", right.Op.String())
	}
}

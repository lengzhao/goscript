// Package parser implements a parser for Go source code
// It wraps the standard Go parser to generate AST
package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// Parser wraps the Go standard library parser
type Parser struct {
	fset *token.FileSet
}

// New creates a new Parser
func New() *Parser {
	return &Parser{
		fset: token.NewFileSet(),
	}
}

// Parse parses the source code and returns the AST
func (p *Parser) Parse(filename string, src []byte, mode parser.Mode) (*ast.File, error) {
	return parser.ParseFile(p.fset, filename, src, mode)
}

// ParseExpr parses a single expression and returns the AST
func (p *Parser) ParseExpr(src []byte) (ast.Expr, error) {
	return parser.ParseExpr(string(src))
}

// FileSet returns the file set used by the parser
func (p *Parser) FileSet() *token.FileSet {
	return p.fset
}

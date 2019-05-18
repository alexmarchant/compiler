package parser

import (
	"github.com/alexmarchant/compiler/lexer"
)

// ValueType ...
type ValueType int

// IntType...
const (
	IntType ValueType = iota
	VoidType
)

type sequence struct {
	elements []sequenceElement
}

type sequenceElement struct {
	tokenType lexer.TokenType
	required  bool
}

func trimLeftNewLines(tokens *[]lexer.Token) {
	tokensVal := *tokens
	for tokensVal[0].Type == lexer.LineBreak {
		tokensVal = tokensVal[1:]
	}
	*tokens = tokensVal
}

// Parse returns an AST of a whole program
func Parse(tokens []lexer.Token) (*Program, error) {
	return parseProgram(&tokens)
}

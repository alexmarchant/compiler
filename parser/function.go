package parser

import (
	"fmt"
	"errors"

	"github.com/alexmarchant/compiler/lexer"
)

// Function is a function
type Function struct {
	Name       string
	ReturnType ValueType
	Statements []*Statement
}

func parseFunction(tokens *[]lexer.Token) (*Function, error) {
	function := Function{}

	tokensVal := *tokens

	if tokensVal[0].Type != lexer.KeywordFunc ||
		tokensVal[1].Type != lexer.Identifier ||
		tokensVal[2].Type != lexer.OpeningParenthesis ||
		tokensVal[3].Type != lexer.ClosingParenthesis {
		return nil, errors.New("Invalid function declaration")
	}

	function.Name = tokensVal[1].Source

	tokensVal = tokensVal[4:]
	if tokensVal[0].Type == lexer.KeywordInt {
		function.ReturnType = IntType
	} else {
		function.ReturnType = VoidType
	}

	tokensVal = tokensVal[1:]
	if tokensVal[0].Type != lexer.OpeningCurlyBrace {
		return nil, errors.New("Invalid function declaration")
	}

	tokensVal = tokensVal[1:]
	*tokens = tokensVal
	statement, err := parseStatement(tokens)
	if err != nil {
		return nil, err
	}
	function.Statements = append(function.Statements, statement)

	if tokensVal[0].Type != lexer.ClosingCurlyBrace {
		fmt.Printf("end: %+v\n", tokensVal)
		return nil, errors.New("Missing closing curly brace")
	}

	*tokens = tokensVal[1:]

	return &function, nil
}
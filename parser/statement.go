package parser

import (
	"errors"
	"github.com/alexmarchant/compiler/lexer"
)

// Statement is chunk of code
type Statement struct {
	Expression Expression
}

func parseStatement(tokens []lexer.Token) (Statement, error) {
	tokens = trimLeftNewLines(tokens)
	statement, err := parseReturnStatement(tokens)
	if err == nil {
		return statement, nil
	}
	return Statement{}, errors.New("Invalid statement")
}

func parseReturnStatement(tokens []lexer.Token) (Statement, error) {
	if tokens[0].TokenType != lexer.KeywordReturn {
		return Statement{}, errors.New("Missing 'return' keyword")
	}
	expression, err := parseExpression(tokens[1:])
	if err != nil {
		return Statement{}, err
	}
	return Statement{Expression: expression}, nil
}
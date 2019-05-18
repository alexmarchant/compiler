package parser

import (
	"fmt"
	"errors"

	"github.com/alexmarchant/compiler/lexer"
)

// Statement is chunk of code
type Statement struct {
	Type                         StatementType
	ReturnStatement              *ReturnStatement
	ExpressionStatement          *ExpressionStatement
	VariableDeclarationStatement *VariableDeclarationStatement
}

// StatementType is a type of statement
type StatementType int

// ReturnStatement et all is an enum
const (
	ReturnStatementType StatementType = iota
	ExpressionStatementType
)

// ReturnStatement ...
type ReturnStatement struct {
	Expression *Expression
}

// ExpressionStatement ...
type ExpressionStatement struct {
	Expression *Expression
}

// VariableDeclarationStatement ...
type VariableDeclarationStatement struct {
}

func parseStatement(tokens *[]lexer.Token) (*Statement, error) {
	trimLeftNewLines(tokens)
	statement, _ := parseReturnStatement(tokens)
	if statement != nil {
		return statement, nil
	}

	statement, _ = parseExpressionStatement(tokens)
	if statement != nil {
		return statement, nil
	}

	return nil, errors.New("Invalid statement")
}

func parseReturnStatement(tokens *[]lexer.Token) (*Statement, error) {
	tokensVal := *tokens
	if tokensVal[0].Type != lexer.KeywordReturn {
		return nil, errors.New("Missing 'return' keyword")
	}

	tokensVal = tokensVal[1:]
	*tokens = tokensVal
	expression, err := parseExpression(tokens)
	if err != nil {
		return nil, err
	}

	returnStatement := ReturnStatement{Expression: expression}
	return &Statement{
		Type:            ReturnStatementType,
		ReturnStatement: &returnStatement,
	}, nil
}

func parseExpressionStatement(tokens *[]lexer.Token) (*Statement, error) {
	expression, err := parseExpression(tokens)
	if err != nil {
		return nil, err
	}

	expressionStatement := ExpressionStatement{Expression: expression}
	return &Statement{
		Type:                ExpressionStatementType,
		ExpressionStatement: &expressionStatement,
	}, nil
}

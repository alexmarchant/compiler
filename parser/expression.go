package parser

import (
	"errors"
	"strconv"

	"github.com/alexmarchant/compiler/lexer"
)

// ExpressionType is an enum
type ExpressionType string

// Addition et all are expression types
const (
	AdditionExpressionType ExpressionType = "AdditionExpressionType"
	LiteralExpressionType  ExpressionType = "LiteralExpressionType"
)

// Expression ...
type Expression interface {
	Type() ExpressionType
}

// AdditionExpression ...
type AdditionExpression struct {
	Left  int
	Right int
}

// Type ...
func (e *AdditionExpression) Type() ExpressionType {
	return AdditionExpressionType
}

// LiteralExpressionValueType ...
type LiteralExpressionValueType string

// IntegerLiteralType ...
const (
	IntegerLiteralType LiteralExpressionValueType = "IntegerLiteralType"
)

// LiteralExpression ...
type LiteralExpression struct {
	Type           LiteralExpressionValueType
	IntegerLiteral int
}

func parseExpression(tokens *[]lexer.Token) (*Expression, error) {
	expression, _ := parseAdditionExpression(tokens)
	if expression != nil {
		return expression, nil
	}

	expression, _ = parseLiteralExpression(tokens)
	if expression != nil {
		return expression, nil
	}

	return nil, errors.New("Invalid expression")
}

func parseAdditionExpression(tokens *[]lexer.Token) (*Expression, error) {
	tokensVal := *tokens
	if tokensVal[0].Type != lexer.IntegerLiteral ||
		tokensVal[1].Type != lexer.PlusSign ||
		tokensVal[2].Type != lexer.IntegerLiteral {
		return nil, errors.New("Invalid addition expression")
	}

	left, _ := strconv.Atoi(tokensVal[0].Source)
	right, _ := strconv.Atoi(tokensVal[2].Source)
	expression := AdditionExpression{Left: left, Right: right}

	tokensVal = tokensVal[3:]
	*tokens = tokensVal

	return &Expression{
		Type:               AdditionExpressionType,
		AdditionExpression: &expression,
	}, nil
}

func parseLiteralExpression(tokens *[]lexer.Token) (*Expression, error) {
	tokensVal := *tokens
	if tokensVal[0].Type != lexer.IntegerLiteral {
		return nil, errors.New("Invalid literal expression")
	}

	value, _ := strconv.Atoi(tokensVal[0].Source)
	literalExpression := LiteralExpression{
		Type:           IntegerLiteralType,
		IntegerLiteral: value,
	}

	tokensVal = tokensVal[1:]
	*tokens = tokensVal

	expression := Expression{
		Type:              LiteralExpressionType,
		LiteralExpression: &literalExpression,
	}
	return &expression, nil
}

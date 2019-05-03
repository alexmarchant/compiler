package parser

import (
	"errors"
	"github.com/alexmarchant/compiler/lexer"
)

// ExpressionType is an enum
type ExpressionType int

// Addition et all are expression types
const (
	Addition ExpressionType = iota
)

// Expression is chunk of code
type Expression struct {
	Type ExpressionType
}

func parseExpression(tokens []lexer.Token) (Expression, error) {
	return parseAdditionExpression(tokens)
}

func parseAdditionExpression(tokens []lexer.Token) (Expression, error) {
	elements := []sequenceElement{
		sequenceElement{tokenType: lexer.IntegerLiteral, required: true},
		sequenceElement{tokenType: lexer.PlusSign, required: true},
		sequenceElement{tokenType: lexer.IntegerLiteral, required: true},
	}
	sequence := sequence{elements: elements}
	match, _ := sequence.match(tokens)
	if !match {
		return Expression{}, errors.New("Invalid addition expression")
	}
	return Expression{Type: Addition}, nil
}
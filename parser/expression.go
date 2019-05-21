package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alexmarchant/compiler/lexer"
)

// BinaryOperator ...
type BinaryOperator string

// BinaryOperatorPlus ...
const (
	BinaryOperatorPlus           BinaryOperator = "BinaryOperatorPlus"
	BinaryOperatorMinus          BinaryOperator = "BinaryOperatorMinus"
	BinaryOperatorMultiplication BinaryOperator = "BinaryOperatorMultiplication"
	BinaryOperatorDivision       BinaryOperator = "BinaryOperatorDivision"
)

func (b BinaryOperator) precedence() int {
	switch b {
	case BinaryOperatorPlus:
		return 20
	case BinaryOperatorMinus:
		return 20
	case BinaryOperatorMultiplication:
		return 40
	case BinaryOperatorDivision:
		return 40
	default:
		panic("Fallthrough")
	}
}

// ExpressionType ...
type ExpressionType string

// ExpressionTypeInt ...
const (
	ExpressionTypeInt                ExpressionType = "ExpressionTypeInt"
	ExpressionTypeString             ExpressionType = "ExpressionTypeString"
	ExpressionTypeArray              ExpressionType = "ExpressionTypeArray"
	ExpressionTypeReturn             ExpressionType = "ExpressionTypeReturn"
	ExpressionTypeBinary             ExpressionType = "ExpressionTypeBinary"
	ExpressionTypeCall               ExpressionType = "ExpressionTypeCall"
	ExpressionTypeParen              ExpressionType = "ExpressionTypeParen"
	ExpressionTypeVariableAssignment ExpressionType = "ExpressionTypeVariableAssignment"
	ExpressionTypeVariable           ExpressionType = "ExpressionTypeVariable"
)

// Expression ...
type Expression interface {
	ExpressionType() ExpressionType
}

// IntExpression ...
type IntExpression struct {
	Value int
}

// ExpressionType ...
func (e *IntExpression) ExpressionType() ExpressionType {
	return ExpressionTypeInt
}

// StringExpression ...
type StringExpression struct {
	Value string
}

// ExpressionType ...
func (e *StringExpression) ExpressionType() ExpressionType {
	return ExpressionTypeString
}

// ArrayExpression ...
type ArrayExpression struct {
	Elements []Expression
}

// ExpressionType ...
func (e *ArrayExpression) ExpressionType() ExpressionType {
	return ExpressionTypeArray
}

// ReturnExpression ...
type ReturnExpression struct {
	Expression Expression
}

// ExpressionType ...
func (e *ReturnExpression) ExpressionType() ExpressionType {
	return ExpressionTypeReturn
}

// BinaryExpression ...
type BinaryExpression struct {
	Op  BinaryOperator
	LHS Expression
	RHS Expression
}

// ExpressionType ...
func (e *BinaryExpression) ExpressionType() ExpressionType {
	return ExpressionTypeBinary
}

// CallExpression ...
type CallExpression struct {
	Callee string
	Params []Expression
}

// ExpressionType ...
func (e *CallExpression) ExpressionType() ExpressionType {
	return ExpressionTypeCall
}

// ParenExpression ...
type ParenExpression struct {
	Expression Expression
}

// ExpressionType ...
func (e *ParenExpression) ExpressionType() ExpressionType {
	return ExpressionTypeParen
}

// AssignmentExpression ...
type VariableExpression struct {
	Name string
}

// ExpressionType ...
func (e *VariableExpression) ExpressionType() ExpressionType {
	return ExpressionTypeVariable
}

// AssignmentExpression ...
type VariableAssignmentExpression struct {
	Name       string
	Type       ValueType
	Expression Expression
}

// ExpressionType ...
func (e *VariableAssignmentExpression) ExpressionType() ExpressionType {
	return ExpressionTypeVariableAssignment
}

func parseExpression() Expression {
	lhs := parsePrimaryExpression()
	return parseBinaryOperatorRHS(0, lhs)
}

func parsePrimaryExpression() Expression {
	switch {
	case tokens[index].Type == lexer.IntegerLiteral:
		return parseIntLiteralExpression()
	case tokens[index].Type == lexer.StringLiteral:
		return parseStringLiteralExpression()
	case tokens[index].Type == lexer.OpeningBracket:
		return parseArrayLiteralExpression()
	case tokens[index].Type == lexer.KeywordReturn:
		return parseReturnExpression()
	case tokens[index].Type == lexer.KeywordVar:
		return parseAssignmentExpression()
	case tokens[index].Type == lexer.OpeningParen:
		return parseParenExpression()
	case tokens[index].Type == lexer.Identifier:
		return parseIdentifierExpression()
	default:
		msg := fmt.Sprintf("Invalid token: %v", tokens[index])
		panic(msg)
	}
}

func parseBinaryOperatorRHS(expressionPrecendence int, lhs Expression) Expression {
	for {
		if tokens[index].Type == lexer.LineBreak || !isBinaryOperator() {
			return lhs
		}
		binOp := parseBinaryOperator()
		tokenPrecedence := binOp.precedence()
		if tokenPrecedence < expressionPrecendence {
			return lhs
		}

		rhs := parsePrimaryExpression()

		if !isBinaryOperator() {
			return &BinaryExpression{
				Op:  binOp,
				LHS: lhs,
				RHS: rhs,
			}
		}

		binOp = parseBinaryOperator()
		nextPrecedence := binOp.precedence()
		if expressionPrecendence < nextPrecedence {
			rhs = parseBinaryOperatorRHS(tokenPrecedence+1, rhs)
		}
		lhs = &BinaryExpression{
			Op:  binOp,
			LHS: lhs,
			RHS: rhs,
		}
	}
}

func isBinaryOperator() bool {
	token := tokens[index]
	switch token.Type {
	case lexer.PlusSign:
		return true
	case lexer.MinusSign:
		return true
	case lexer.MultiplicationSign:
		return true
	case lexer.DivisionSign:
		return true
	default:
		return false
	}
}

func parseBinaryOperator() BinaryOperator {
	token := tokens[index]
	index++

	switch token.Type {
	case lexer.PlusSign:
		return BinaryOperatorPlus
	case lexer.MinusSign:
		return BinaryOperatorMinus
	case lexer.MultiplicationSign:
		return BinaryOperatorMultiplication
	case lexer.DivisionSign:
		return BinaryOperatorDivision
	default:
		panic("Fallthrough")
	}
}

func parseIntLiteralExpression() *IntExpression {
	value, err := strconv.Atoi(tokens[index].Source)
	if err != nil {
		panic("Invalid int")
	}
	index++
	return &IntExpression{
		Value: value,
	}
}

func parseStringLiteralExpression() *StringExpression {
	token := tokens[index]
	value := strings.Trim(token.Source, "\"")
	index++
	return &StringExpression{
		Value: value,
	}
}

func parseArrayLiteralExpression() *ArrayExpression {
	index++
	expressions := []Expression{}
	for {
		if tokens[index].Type == lexer.ClosingBracket {
			index++
			break
		}
		if tokens[index].Type == lexer.Comma {
			index++
			continue
		}

		expressions = append(expressions, parseExpression())
	}
	return &ArrayExpression{
		Elements: expressions,
	}
}

func parseReturnExpression() *ReturnExpression {
	if tokens[index].Type != lexer.KeywordReturn {
		panic("Invalid return expression")
	}
	index++
	return &ReturnExpression{Expression: parseExpression()}
}

func parseAssignmentExpression() *VariableAssignmentExpression {
	exp := &VariableAssignmentExpression{}

	if tokens[index].Type != lexer.KeywordVar {
		panic("Invalid assignment expression")
	}
	index++

	if tokens[index].Type != lexer.Identifier {
		panic("Invalid assignment expression")
	}
	exp.Name = tokens[index].Source
	index++

	if tokens[index].Type != lexer.Colon {
		panic("Invalid assignment expression")
	}
	index++

	expType := parseValueType()
	if expType == nil {
		panic("Invalid assignment expression")
	}
	exp.Type = *expType

	if tokens[index].Type != lexer.Equals {
		panic("Invalid assignment expression")
	}
	index++

	exp.Expression = parseExpression()

	return exp
}

func parseParenExpression() Expression {
	if tokens[index].Type != lexer.OpeningParen {
		panic("Invalid paren expression")
	}
	index++
	expression := parseExpression()
	if tokens[index].Type != lexer.ClosingParen {
		panic("Invalid paren expression")
	}
	index++
	return &ParenExpression{
		Expression: expression,
	}
}

func parseIdentifierExpression() Expression {
	if tokens[index+1].Type == lexer.OpeningParen {
		return parseCallExpression()
	} else {
		name := tokens[index].Source
		index++
		return &VariableExpression{
			Name: name,
		}
	}
}

func parseCallExpression() *CallExpression {
	name := tokens[index].Source
	index++

	// (
	if tokens[index].Type != lexer.OpeningParen {
		panic("Invalid identifier expression")
	}
	index++

	// Parse params
	expressions := []Expression{}
	for {
		// check for closing paren )
		if tokens[index].Type == lexer.ClosingParen {
			index++
			break
		}
		// check for comma
		if tokens[index].Type == lexer.Comma {
			index++
			continue
		}
		// parse expression
		expressions = append(
			expressions,
			parseExpression())
	}

	return &CallExpression{
		Callee: name,
		Params: expressions,
	}
}

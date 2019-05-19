package parser

import (
	"strconv"

	"github.com/alexmarchant/compiler/lexer"
)

// NodeType ...
type NodeType string

// NodeTypeFunction ...
const (
	NodeTypeFunction   NodeType = "NodeTypeFunction"
	NodeTypeExpression NodeType = "NodeTypeExpression"
)

// Node ...
type Node interface {
	NodeType() NodeType
}

// ValueType ...
type ValueType string

// ValueTypeInt ...
const (
	ValueTypeInt ValueType = "ValueTypeInt"
)

// Function ...
type Function struct {
	Name        string
	ReturnType  *ValueType
	Expressions []Expression
}

// NodeType ...
func (f *Function) NodeType() NodeType {
	return NodeTypeFunction
}

// ExpressionType ...
type ExpressionType string

// ExpressionTypeInt ...
const (
	ExpressionTypeInt    ExpressionType = "ExpressionTypeInt"
	ExpressionTypeReturn ExpressionType = "ExpressionTypeReturn"
)

// Expression ...
type Expression interface {
	NodeType() NodeType
	ExpressionType() ExpressionType
}

// ExpressionNode ...
type ExpressionNode struct{}

// NodeType ...
func (e *ExpressionNode) NodeType() NodeType {
	return NodeTypeExpression
}

// IntExpression ...
type IntExpression struct {
	ExpressionNode
	Value int
}

// ExpressionType ...
func (e *IntExpression) ExpressionType() ExpressionType {
	return ExpressionTypeInt
}

// ReturnExpression ...
type ReturnExpression struct {
	ExpressionNode
	Expression Expression
}

// ExpressionType ...
func (e *ReturnExpression) ExpressionType() ExpressionType {
	return ExpressionTypeReturn
}

var tokens []lexer.Token
var index = 0

// Parse returns an AST of a whole program
func Parse(parseTokens []lexer.Token) []Node {
	nodes := []Node{}
	tokens = parseTokens

	for {
		token := tokens[index]

		switch {
		case token.Type == lexer.EOF:
			return nodes
		case token.Type == lexer.KeywordFunc:
			nodes = append(
				nodes,
				parseFunction())
		default:
			nodes = append(
				nodes,
				parseExpression())
		}
	}
}

func parseFunction() *Function {
	function := Function{}

	if tokens[index].Type != lexer.KeywordFunc {
		panic("Function declaration missing func keyword")
	}
	index++

	if tokens[index].Type != lexer.Identifier {
		panic("Function declaration missing name")
	}
	function.Name = tokens[index].Source
	index++

	if tokens[index].Type != lexer.OpeningParen {
		panic("Function declaration missing opening paren")
	}
	index++

	if tokens[index].Type != lexer.ClosingParen {
		panic("Function declaration missing closing paren")
	}
	index++

	function.ReturnType = parseValueType()

	if tokens[index].Type != lexer.OpeningCurlyBrace {
		panic("Function declaration missing opening curly brace")
	}
	index++

	for {
		// Parse expressions until we hit closing brace
		if tokens[index].Type == lexer.ClosingCurlyBrace {
			index++
			break
		}

		// Skip line breaks
		if tokens[index].Type == lexer.LineBreak {
			index++
			continue
		}

		function.Expressions = append(
			function.Expressions,
			parseExpression())
	}

	return &function
}

func parseValueType() *ValueType {
	token := tokens[index]

	switch token.Type {
	case lexer.KeywordInt:
		index++
		value := ValueTypeInt
		return &value
	default:
		return nil
	}
}

func parseExpression() Expression {
	switch {
	case tokens[index].Type == lexer.IntegerLiteral:
		return parseIntExpression()
	case tokens[index].Type == lexer.KeywordReturn:
		return parseReturnExpression()
	default:
		panic("Invalid expression")
	}
}

func parseIntExpression() *IntExpression {
	value, err := strconv.Atoi(tokens[index].Source)
	if err != nil {
		panic("Invalid int")
	}
	index++
	return &IntExpression{
		Value: value,
	}
}

func parseReturnExpression() *ReturnExpression {
	if tokens[index].Type != lexer.KeywordReturn {
		panic("Invalid return expression")
	}
	index++
	return &ReturnExpression{Expression: parseExpression()}
}

package parser

import "github.com/alexmarchant/compiler/lexer"

// Prototype ...
type Prototype struct {
	Name       string
	ReturnType *ValueType
}

// Function ...
type Function struct {
	Prototype   *Prototype
	Expressions []Expression
}

// NodeType ...
func (f *Function) NodeType() NodeType {
	return NodeTypeFunction
}

func parsePrototype() *Prototype {
	prototype := &Prototype{}

	if tokens[index].Type != lexer.KeywordFn {
		panic("Function declaration missing func keyword")
	}
	index++

	if tokens[index].Type != lexer.Identifier {
		panic("Function declaration missing name")
	}
	prototype.Name = tokens[index].Source
	index++

	if tokens[index].Type != lexer.OpeningParen {
		panic("Function declaration missing opening paren")
	}
	index++

	if tokens[index].Type != lexer.ClosingParen {
		panic("Function declaration missing closing paren")
	}
	index++

	prototype.ReturnType = parseValueType()

	return prototype
}

func parseFunction() *Function {
	function := &Function{}
	function.Prototype = parsePrototype()

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

	return function
}

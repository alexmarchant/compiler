package parser

import (
	"errors"
	"github.com/alexmarchant/compiler/lexer"
)

// Function is a function
type Function struct {
	Name       string
	Statements []Statement
}

func parseFunction(tokens []lexer.Token) (Function, error) {
	function := Function{}

	elements := []sequenceElement{
		sequenceElement{tokenType: lexer.KeywordFunc, required: true},
		sequenceElement{tokenType: lexer.Identifier, required: true},
		sequenceElement{tokenType: lexer.OpeningParenthesis, required: true},
		sequenceElement{tokenType: lexer.ClosingParenthesis, required: true},
		sequenceElement{tokenType: lexer.LineBreak, required: false},
		sequenceElement{tokenType: lexer.OpeningCurlyBrace, required: true},
	}
	functionDeclaration := sequence{elements: elements}
	match, endIndex := functionDeclaration.match(tokens)
	if !match {
		return function, errors.New("Invalid function declaration")
	}

	function.Name = tokens[1].Source
	tokens = tokens[endIndex:]
	
	statement, err := parseStatement(tokens)
	if err != nil {
		return function, err
	}
	function.Statements = append(function.Statements, statement)

	return function, nil
}

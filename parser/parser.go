package parser

import (
	"errors"
	"github.com/alexmarchant/compiler/lexer"
)

type statement struct {
	tokens []lexer.Token
}

type tokenSequence struct {
	elements []tokenSequenceElement
}

type tokenSequenceElement struct {
	tokenType lexer.TokenType
	required  bool
}

func (t tokenSequence) match(tokens []lexer.Token) bool {
	currentIndex := 0

	for _, el := range t.elements {
		if el.tokenType != tokens[currentIndex].TokenType && el.required {
			return false
		}
		currentIndex++
	}

	return true
}

// Parse returns an AST of a whole program
func Parse(tokens []lexer.Token) (Program, error) {
	return parseProgram(tokens)
}

func parseFunction(tokens []lexer.Token) (Function, []lexer.Token, error) {
	thisFunction := Function{}

	elements := []tokenSequenceElement{
		tokenSequenceElement{tokenType: lexer.KeywordFunc, required: true},
		tokenSequenceElement{tokenType: lexer.Identifier, required: true},
		tokenSequenceElement{tokenType: lexer.OpeningParenthesis, required: true},
		tokenSequenceElement{tokenType: lexer.ClosingParenthesis, required: true},
		tokenSequenceElement{tokenType: lexer.OpeningCurlyBrace, required: true},
	}
	functionDeclaration := tokenSequence{elements: elements}
	if !functionDeclaration.match(tokens) {
		return thisFunction, tokens, errors.New("Invalid function declaration")
	}

	thisFunction.Name = tokens[1].Source
	tokens = tokens[4:]
	thisStatement := statement{}

	for _, token := range tokens {
		if token.TokenType == lexer.ClosingCurlyBrace {
			break
		}
		thisStatement.tokens = append(thisStatement.tokens, token)
	}

	thisFunction.Statements = append(thisFunction.Statements, thisStatement)
	return thisFunction, tokens, nil
}

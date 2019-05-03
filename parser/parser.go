package parser

import (
	"github.com/alexmarchant/compiler/lexer"
)

type sequence struct {
	elements []sequenceElement
}

type sequenceElement struct {
	tokenType lexer.TokenType
	required  bool
}

func (t sequence) match(tokens []lexer.Token) (bool, int) {
	currentIndex := 0

	for _, el := range t.elements {
		if el.tokenType != tokens[currentIndex].TokenType {
			if el.required {
				return false, 0
			} else {
				continue
			}
		}
		currentIndex++
	}

	return true, currentIndex
}

func trimLeftNewLines(tokens []lexer.Token) []lexer.Token {
	for tokens[0].TokenType == lexer.LineBreak {
		tokens = tokens[1:]
	}
	return tokens
}

// Parse returns an AST of a whole program
func Parse(tokens []lexer.Token) (Program, error) {
	return parseProgram(tokens)
}

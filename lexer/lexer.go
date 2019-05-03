package lexer

import (
	"regexp"
	"strings"
)

// TokenType is an enum
type TokenType int

// KeywordFun et all are TokenTypes
const (
	KeywordFunc TokenType = iota
	KeywordReturn
	Identifier
	IntegerLiteral
	OpeningParenthesis
	ClosingParenthesis
	OpeningCurlyBrace
	ClosingCurlyBrace
	LineBreak
	PlusSign
)

func (t TokenType) tokenTypeRegex() string {
	switch t {
	case KeywordFunc:
		return "^func"
	case KeywordReturn:
		return "^return"
	case Identifier:
		return "^[a-zA-Z]+"
	case IntegerLiteral:
		return "^\\d+"
	case OpeningParenthesis:
		return "^\\("
	case ClosingParenthesis:
		return "^\\)"
	case OpeningCurlyBrace:
		return "^{"
	case ClosingCurlyBrace:
		return "^}"
	case LineBreak:
		return "^\n"
	case PlusSign:
		return "^\\+"
	default:
		panic("Switch fallthrough")
	}
}

// Token is a token
type Token struct {
	TokenType TokenType
	Source    string
}

func match(regexString string, source string) (bool, []int) {
	r, _ := regexp.Compile(regexString)
	match := r.MatchString(source)
	if match {
		indexes := r.FindStringIndex(source)
		return true, indexes
	} else {
		return false, []int{}
	}
}

func parseTokens(source string) []Token {
	source = strings.Trim(source, " ")
	tokens := []Token{}
	tokenTypes := []TokenType{
		KeywordFunc,
		KeywordReturn,
		OpeningParenthesis,
		ClosingParenthesis,
		OpeningCurlyBrace,
		ClosingCurlyBrace,
		LineBreak,
		PlusSign,
		IntegerLiteral,
		Identifier,
	}

	for len(source) > 0 {
		source = strings.Trim(source, " ")
		found := false

		for _, tokenType := range tokenTypes {
			regex := tokenType.tokenTypeRegex()
			if match, indexes := match(regex, source); match {
				tokens = append(tokens, Token{
					TokenType: tokenType,
					Source:    source[indexes[0]:indexes[1]],
				})
				source = source[indexes[1]:]
				found = true
				break
			}
		}

		if !found {
			panic("Can't find next token")
		}
	}

	return tokens
}

// Lex returns tokens
func Lex(source string) []Token {
	return parseTokens(source)
}

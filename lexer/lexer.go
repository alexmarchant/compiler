package lexer

import (
	"regexp"
	"strings"
)

// TokenType is an enum
type TokenType string

// KeywordFun et all are TokenTypes
const (
	KeywordFunc       TokenType = "KeywordFunc"
	KeywordReturn     TokenType = "KeywordReturn"
	KeywordInt        TokenType = "KeywordInt"
	Identifier        TokenType = "Identifier"
	IntegerLiteral    TokenType = "IntegerLiteral"
	OpeningParen      TokenType = "OpeningParen"
	ClosingParen      TokenType = "ClosingParen"
	OpeningCurlyBrace TokenType = "OpeningCurlyBrace"
	ClosingCurlyBrace TokenType = "ClosingCurlyBrace"
	LineBreak         TokenType = "LineBreak"
	PlusSign          TokenType = "PlusSign"
	EOF               TokenType = "EOF"
)

func (t TokenType) tokenTypeRegex() string {
	switch t {
	case KeywordFunc:
		return "^func"
	case KeywordReturn:
		return "^return"
	case KeywordInt:
		return "^int"
	case Identifier:
		return "^[a-zA-Z]+"
	case IntegerLiteral:
		return "^\\d+"
	case OpeningParen:
		return "^\\("
	case ClosingParen:
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
	Type   TokenType
	Source string
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
		KeywordInt,
		OpeningParen,
		ClosingParen,
		OpeningCurlyBrace,
		ClosingCurlyBrace,
		LineBreak,
		PlusSign,
		IntegerLiteral,
		Identifier,
		EOF,
	}

	for len(source) > 0 {
		source = strings.Trim(source, " ")
		found := false

		for _, tokenType := range tokenTypes {
			regex := tokenType.tokenTypeRegex()
			if match, indexes := match(regex, source); match {
				tokens = append(tokens, Token{
					Type:   tokenType,
					Source: source[indexes[0]:indexes[1]],
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

	tokens = append(tokens, Token{
		Type:   EOF,
		Source: "",
	})

	return tokens
}

// Lex returns tokens
func Lex(source string) []Token {
	return parseTokens(source)
}

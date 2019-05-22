package lexer

import (
	"fmt"
	"regexp"
	"strings"
)

// TokenType is an enum
type TokenType string

// KeywordFun et all are TokenTypes
const (
	KeywordFn          TokenType = "KeywordFn"
	KeywordReturn      TokenType = "KeywordReturn"
	KeywordVar         TokenType = "KeywordVar"
	KeywordStruct      TokenType = "KeywordStruct"
	KeywordInt         TokenType = "KeywordInt"
	KeywordIntArray    TokenType = "KeywordIntArray"
	KeywordString      TokenType = "KeywordString"
	Identifier         TokenType = "Identifier"
	IntegerLiteral     TokenType = "IntegerLiteral"
	StringLiteral      TokenType = "StringLiteral"
	Period             TokenType = "Period"
	Colon              TokenType = "Colon"
	Equals             TokenType = "Equals"
	OpeningParen       TokenType = "OpeningParen"
	ClosingParen       TokenType = "ClosingParen"
	OpeningCurlyBrace  TokenType = "OpeningCurlyBrace"
	ClosingCurlyBrace  TokenType = "ClosingCurlyBrace"
	LineBreak          TokenType = "LineBreak"
	PlusSign           TokenType = "PlusSign"
	MinusSign          TokenType = "MinusSign"
	MultiplicationSign TokenType = "MultiplicationSign"
	DivisionSign       TokenType = "DivisionSign"
	EOF                TokenType = "EOF"
	Comma              TokenType = "Comma"
	OpeningBracket     TokenType = "OpeningBracket"
	ClosingBracket     TokenType = "ClosingBracket"
)

func (t TokenType) tokenTypeRegex() string {
	switch t {
	case KeywordFn:
		return "^fn"
	case KeywordReturn:
		return "^return"
	case KeywordVar:
		return "^var"
	case KeywordStruct:
		return "^struct"
	case KeywordInt:
		return "^Int"
	case KeywordIntArray:
		return "^IntArray"
	case KeywordString:
		return "^String"
	case Identifier:
		return "^[a-zA-Z_]+"
	case IntegerLiteral:
		return "^\\d+"
	case StringLiteral:
		return "^\"(.*?)\""
	case Period:
		return "^\\."
	case Colon:
		return "^:"
	case Equals:
		return "^="
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
	case MinusSign:
		return "^\\-"
	case MultiplicationSign:
		return "^\\*"
	case DivisionSign:
		return "^\\/"
	case Comma:
		return "^,"
	case OpeningBracket:
		return "^\\["
	case ClosingBracket:
		return "^\\]"
	default:
		msg := fmt.Sprintf("Unrecognized token: %s", t)
		panic(msg)
	}
}

// Token is a token
type Token struct {
	Type   TokenType
	Source string
}

// Lex returns tokens
func Lex(source string) []Token {
	return parseTokens(source)
}

func parseTokens(source string) []Token {
	source = strings.Trim(source, " ")
	tokens := []Token{}
	tokenTypes := []TokenType{
		KeywordFn,
		KeywordReturn,
		KeywordVar,
		KeywordStruct,
		KeywordIntArray,
		KeywordInt,
		KeywordString,
		Period,
		Colon,
		Equals,
		OpeningParen,
		ClosingParen,
		OpeningCurlyBrace,
		ClosingCurlyBrace,
		LineBreak,
		PlusSign,
		MinusSign,
		MultiplicationSign,
		DivisionSign,
		IntegerLiteral,
		StringLiteral,
		Identifier,
		Comma,
		OpeningBracket,
		ClosingBracket,
	}

	for len(source) > 0 {
		source = strings.Trim(source, " \t")
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
			msg := fmt.Sprintf("Next token not recognized: %s", source)
			panic(msg)
		}
	}

	tokens = append(tokens, Token{
		Type:   EOF,
		Source: "",
	})

	return tokens
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

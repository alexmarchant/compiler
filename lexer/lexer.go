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
)

// Token is a token
type Token struct {
	TokenType TokenType
	Source    string
}

func parseTokens(source string) []Token {
	source = strings.Trim(source, " ")
	tokens := []Token{}

	for len(source) > 0 {
		source = strings.Trim(source, " ")

		if match, _ := regexp.MatchString("^\n", source); match {
			tokens = append(tokens, Token{
				TokenType: LineBreak,
				Source:    "\\n",
			})
			source = source[1:]
			continue
		}

		if match, _ := regexp.MatchString("^func", source); match {
			tokens = append(tokens, Token{
				TokenType: KeywordFunc,
				Source:    "func",
			})
			source = source[4:]
			continue
		}

		if match, _ := regexp.MatchString("^func", source); match {
			tokens = append(tokens, Token{
				TokenType: KeywordFunc,
				Source:    "func",
			})
			source = source[4:]
			continue
		}

		if match, _ := regexp.MatchString("^return", source); match {
			tokens = append(tokens, Token{
				TokenType: KeywordReturn,
				Source:    "return",
			})
			source = source[6:]
			continue
		}

		if match, _ := regexp.MatchString("^\\(", source); match {
			tokens = append(tokens, Token{
				TokenType: OpeningParenthesis,
				Source:    "(",
			})
			source = source[1:]
			continue
		}

		if match, _ := regexp.MatchString("^\\)", source); match {
			tokens = append(tokens, Token{
				TokenType: ClosingParenthesis,
				Source:    ")",
			})
			source = source[1:]
			continue
		}

		if match, _ := regexp.MatchString("^{", source); match {
			tokens = append(tokens, Token{
				TokenType: OpeningCurlyBrace,
				Source:    "{",
			})
			source = source[1:]
			continue
		}

		if match, _ := regexp.MatchString("^}", source); match {
			tokens = append(tokens, Token{
				TokenType: ClosingCurlyBrace,
				Source:    "}",
			})
			source = source[1:]
			continue
		}

		r, _ := regexp.Compile("^\\d+")
		if r.MatchString(source) {
			indexes := r.FindStringIndex(source)
			tokens = append(tokens, Token{
				TokenType: IntegerLiteral,
				Source:    source[indexes[0]:indexes[1]],
			})
			source = source[indexes[1]:]
			continue
		}

		r, _ = regexp.Compile("^[a-zA-Z]+")
		if r.MatchString(source) {
			indexes := r.FindStringIndex(source)
			tokens = append(tokens, Token{
				TokenType: Identifier,
				Source:    source[indexes[0]:indexes[1]],
			})
			source = source[indexes[1]:]
			continue
		}

		panic("Can't find next token")
	}

	return tokens
}

// Lex returns tokens
func Lex(source string) []Token {
	return parseTokens(source)
}

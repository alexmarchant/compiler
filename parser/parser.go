package parser

import (
	"fmt"
	"errors"

	"github.com/alexmarchant/compiler/lexer"
)

// NodeType ...
type NodeType string

// NodeTypeFunction ...
const (
	NodeTypeFunction NodeType = "NodeTypeFunction"
	NodeTypeStruct   NodeType = "NodeTypeStruct"
)

// Node ...
type Node interface {
	NodeType() NodeType
}

// Prop ...
type Prop struct {
	Name string
	Type string
}

var tokens []lexer.Token
var index int

// Parse returns an AST of a whole program
func Parse(someTokens []lexer.Token) []Node {
	nodes := []Node{}
	tokens = someTokens
	index = 0

	for {
		token := tokens[index]

		switch {
		case token.Type == lexer.EOF:
			return nodes
		case token.Type == lexer.LineBreak:
			index++
			continue
		case token.Type == lexer.KeywordFn:
			nodes = append(
				nodes,
				parseFunction())
		case token.Type == lexer.KeywordStruct:
			nodes = append(
				nodes,
				parseStruct())
		default:
			msg := fmt.Sprintf("Don't know how to parse: %v", token)
			panic(msg)
		}
	}
}

func parseValueType() (string, error) {
	token := tokens[index]

	switch token.Type {
	case lexer.KeywordInt:
		index++
		return "int", nil
	case lexer.KeywordString:
		index++
		return "String*", nil
	case lexer.Identifier:
		index++
		return token.Source, nil
	default:
		return "", errors.New("Invalid value type")
	}
}

func parseProp() *Prop {
	prop := &Prop{}
	if tokens[index].Type != lexer.Identifier {
		panic("Struct prop missing name")
	}
	prop.Name = tokens[index].Source
	index++

	if tokens[index].Type != lexer.Colon {
		panic("Struct prop missing colon")
	}
	index++

	valueType, err := parseValueType()
	if err != nil {
		panic("Struct prop has invalid type")
	}
	prop.Type = valueType

	return prop
}
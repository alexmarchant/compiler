package parser

import (
	"github.com/alexmarchant/compiler/lexer"
)

// NodeType ...
type NodeType string

// NodeTypeFunction ...
const (
	NodeTypeFunction NodeType = "NodeTypeFunction"
)

// Node ...
type Node interface {
	NodeType() NodeType
}

// ValueType ...
type ValueType string

// ValueTypeInt ...
const (
	ValueTypeInt      ValueType = "ValueTypeInt"
	ValueTypeString   ValueType = "ValueTypeString"
	ValueTypeIntArray ValueType = "ValueTypeIntArray"
)

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
		default:
			panic("Fallthrough")
		}
	}
}

func parseValueType() *ValueType {
	token := tokens[index]

	switch token.Type {
	case lexer.KeywordInt:
		index++
		value := ValueTypeInt
		return &value
	case lexer.KeywordString:
		index++
		value := ValueTypeString
		return &value
	case lexer.KeywordIntArray:
		index++
		value := ValueTypeIntArray
		return &value
	default:
		return nil
	}
}

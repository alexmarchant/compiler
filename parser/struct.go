package parser

import (
	"github.com/alexmarchant/compiler/lexer"
)

// Struct ...
type Struct struct {
	Name      string
	Props     []*StructProp
	Functions []*Function
}

// NodeType ...
func (f *Struct) NodeType() NodeType {
	return NodeTypeStruct
}

// StructProp ...
type StructProp struct {
	Name string
	Type ValueType
}

func parseStruct() *Struct {
	str := &Struct{}

	if tokens[index].Type != lexer.KeywordStruct {
		panic("Struct missing struct keyword")
	}
	index++

	if tokens[index].Type != lexer.Identifier {
		panic("Struct missing name")
	}
	str.Name = tokens[index].Source
	index++

	if tokens[index].Type != lexer.OpeningCurlyBrace {
		panic("Struct missing opening curly Brace")
	}
	index++

	for {
		if tokens[index].Type == lexer.LineBreak {
			index++
			continue
		}

		if tokens[index].Type == lexer.ClosingCurlyBrace {
			index++
			break
		}

		switch tokens[index].Type {
		case lexer.KeywordFn:
			str.Functions = append(
				str.Functions,
				parseFunction())
		case lexer.Identifier:
			str.Props = append(
				str.Props,
				parseStructProp())
		default:
			panic("Invalid struct")
		}
	}

	return str
}

func parseStructProp() *StructProp {
	prop := &StructProp{}
	if tokens[index].Type != lexer.Identifier {
		panic("Struct prop missing name")
	}
	prop.Name = tokens[index].Source
	index++

	if tokens[index].Type != lexer.Colon {
		panic("Struct prop missing colon")
	}
	index++

	valueType := parseValueType()
	if valueType == nil {
		panic("Struct prop has invalid type")
	}
	prop.Type = *valueType

	return prop
}

package generator

import (
	"fmt"
	"strings"

	"github.com/alexmarchant/compiler/parser"
)

func isBuiltin(callee string) bool {
	switch callee {
	case "println":
		return true
	default:
		return false
	}
}

func generateBuiltin(exp *parser.CallExpression, functionVariables *map[string]string) string {
	switch exp.Callee {
	case "println":
		return generatePrintln(exp, functionVariables)
	default:
		panic("Unkown builtin")
	}
}

func generatePrintln(exp *parser.CallExpression, functionVariables *map[string]string) string {
	format := []string{}
	args := []string{}
	for _, param := range exp.Params {
		switch param.ExpressionType() {
		case parser.ExpressionTypeInt:
			format = append(format, "%d")
			args = append(args, generateExpression(param, functionVariables))
		case parser.ExpressionTypeString:
			format = append(format, "%s")
			arg := fmt.Sprintf("%s->value", generateExpression(param, functionVariables))
			args = append(args, arg)
		case parser.ExpressionTypeVariable:
			exp := param.(*parser.VariableExpression)
			val, ok := (*functionVariables)[exp.Name]
			if !ok {
				panic("Undeclared variable")
			}
			val = strings.Trim(val, "*")
			switch {
			case val == "int":
				format = append(format, "%d")
				args = append(args, generateExpression(param, functionVariables))
			case val == "String":
				format = append(format, "%s")
				arg := fmt.Sprintf("%s->value", exp.Name)
				args = append(args, arg)
			case isTypeStruct(val):
				format = append(format, "%s")
				arg := fmt.Sprintf("%s__toString(%s)->value", val, exp.Name)
				args = append(args, arg)
			default:
				msg := fmt.Sprintf("Unknown type: %s", val)
				panic(msg)
			}
		default:
			panic("Not handling this case yet")
		}
	}
	formatString := strings.Join(format, " ")
	argsString := strings.Join(args, ", ")
	return fmt.Sprintf("printf(\"%s\\n\", %s)", formatString, argsString)
}

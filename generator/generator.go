package generator

import (
	"fmt"
	"strings"

	"github.com/alexmarchant/compiler/parser"
)

// Generate ...
func Generate(nodes []parser.Node) string {
	code := "#include <stdio.h>\n"
	code += "#include \"runtime/runtime.h\"\n"
	code += "\n"

	for _, node := range nodes {
		code += generateNode(node)
	}

	return code
}

func generateNode(node parser.Node) string {
	switch node.NodeType() {
	case parser.NodeTypeFunction:
		function := node.(*parser.Function)
		return generateFunction(function)
	case parser.NodeTypeStruct:
		str := node.(*parser.Struct)
		return generateStruct(str)
	default:
		panic("Invalid NodeType")
	}
}

func generateFunction(function *parser.Function) string {
	code := ""

	if function.Prototype.ReturnType == nil {
		if function.Prototype.Name == "main" {
			code += "int "
		} else {
			code += "void "
		}
	} else {
		code += fmt.Sprintf("%s ", cValueType(*function.Prototype.ReturnType))
	}

	var funcName string
	if function.Prototype.Name == "main" {
		code += fmt.Sprintf("%s", function.Prototype.Name)
	} else {
		code += fmt.Sprintf("func_%s", function.Prototype.Name)
	}
	code += fmt.Sprintf("%s() {\n", funcName)

	header := ""
	body := ""
	varCount := 0

	for _, expression := range function.Expressions {
		body += "\t"
		body += generateExpression(expression, &header, &varCount)
		body += ";\n"
	}

	if function.Prototype.Name == "main" && function.Prototype.ReturnType == nil {
		body += "\treturn 0;\n"
	}

	code += header
	code += "\n"
	code += body
	code += "}\n\n"

	return code
}

func generateStruct(str *parser.Struct) string {
	code := fmt.Sprintf("typedef struct _%s {\n", str.Name)
	for _, prop := range str.Props {
		code += fmt.Sprintf("\t%s %s;\n", cValueType(prop.Type), prop.Name)
	}
	code += fmt.Sprintf("} %s;\n\n", str.Name)
	return code
}

func generateExpression(expression parser.Expression, functionHeader *string, functionVarCount *int) string {
	switch expression.ExpressionType() {
	case parser.ExpressionTypeInt:
		exp := expression.(*parser.IntExpression)
		return fmt.Sprintf("%d", exp.Value)
	case parser.ExpressionTypeString:
		exp := expression.(*parser.StringExpression)
		*functionVarCount++
		arrayVarName := fmt.Sprintf("var_%d", *functionVarCount)
		*functionHeader += fmt.Sprintf("\tCharArray* %s = char_array_make();\n", arrayVarName)
		*functionHeader += fmt.Sprintf(
			"\tchar_array_add(%s, \"%s\");\n",
			arrayVarName,
			exp.Value)
		return arrayVarName
	case parser.ExpressionTypeArray:
		exp := expression.(*parser.ArrayExpression)
		*functionVarCount++
		arrayVarName := fmt.Sprintf("var_%d", *functionVarCount)
		*functionHeader += fmt.Sprintf("\tIntArray* %s = int_array_make();\n", arrayVarName)
		for _, elExp := range exp.Elements {
			*functionHeader += fmt.Sprintf(
				"\tint_array_push(%s, %s);\n",
				arrayVarName,
				generateExpression(elExp, functionHeader, functionVarCount))
		}
		return arrayVarName
	case parser.ExpressionTypeReturn:
		exp := expression.(*parser.ReturnExpression)
		code := "return "
		code += generateExpression(exp.Expression, functionHeader, functionVarCount)
		return code
	case parser.ExpressionTypeBinary:
		exp := expression.(*parser.BinaryExpression)
		code := generateExpression(exp.LHS, functionHeader, functionVarCount)
		code += fmt.Sprintf(" %s ", cBinaryOperator(exp.Op))
		code += generateExpression(exp.RHS, functionHeader, functionVarCount)
		return code
	case parser.ExpressionTypeCall:
		exp := expression.(*parser.CallExpression)

		// Call c functions directly
		if exp.Callee == "callCFunc" {
			functionNameExpression := exp.Params[0].(*parser.StringExpression)
			code := fmt.Sprintf("%s(", functionNameExpression.Value)
			paramCode := []string{}
			for _, exp := range exp.Params[1:] {
				stringExp := exp.(*parser.StringExpression)
				paramCode = append(
					paramCode,
					fmt.Sprintf("\"%s\"", stringExp.Value))
			}
			code += strings.Join(paramCode, ", ")
			code += ")"
			return code
		}

		// Native functions
		code := fmt.Sprintf("func_%s(", exp.Callee)
		paramCode := []string{}
		for _, exp := range exp.Params {
			paramCode = append(
				paramCode,
				generateExpression(exp, functionHeader, functionVarCount))
		}
		code += strings.Join(paramCode, ", ")
		code += ")"
		return code
	case parser.ExpressionTypeParen:
		exp := expression.(*parser.ParenExpression)
		code := "("
		code += generateExpression(exp.Expression, functionHeader, functionVarCount)
		code += ")"
		return code
	case parser.ExpressionTypeVariableAssignment:
		exp := expression.(*parser.VariableAssignmentExpression)
		code := cValueType(exp.Type)
		code += fmt.Sprintf(" %s = ", exp.Name)
		code += generateExpression(exp.Expression, functionHeader, functionVarCount)
		return code
	case parser.ExpressionTypeVariable:
		exp := expression.(*parser.VariableExpression)
		return exp.Name
	default:
		panic("Fallthrough")
	}
}

func cBinaryOperator(binOp parser.BinaryOperator) string {
	switch binOp {
	case parser.BinaryOperatorPlus:
		return "+"
	case parser.BinaryOperatorMinus:
		return "-"
	case parser.BinaryOperatorMultiplication:
		return "*"
	case parser.BinaryOperatorDivision:
		return "/"
	default:
		panic("Fallthrough")
	}
}

func cValueType(valueType parser.ValueType) string {
	switch valueType {
	case parser.ValueTypeInt:
		return "int"
	case parser.ValueTypeString:
		return "CharArray*"
	case parser.ValueTypeIntArray:
		return "IntArray*"
	default:
		panic("Unrecognized value type")
	}
}

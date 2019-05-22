package generator

import (
	"fmt"
	"strings"
	"io/ioutil"
	"os/exec"
	"bytes"

	"github.com/alexmarchant/compiler/parser"
)

var customTypes = map[string]string{}

// GenerateC ...
func GenerateC(nodes []parser.Node) string {
	code := "#include <stdio.h>\n"
	code += "#include <stdlib.h>\n"
	code += "#include \"runtime/runtime.h\"\n"
	code += "\n"

	for _, node := range nodes {
		code += generateNode(node)
	}

	ioutil.WriteFile("./out.c", []byte(code), 0644)
	return code
}

// CompileC ...
func CompileC() {
	cmd := exec.Command("sh", "-c", "clang out.c runtime/*.c -o out")
	var errLog bytes.Buffer
	cmd.Stderr = &errLog
	err := cmd.Run()
	if len(errLog.String()) > 0 {
		fmt.Println("\n--COMPILATION--")
		fmt.Printf(errLog.String())
	}
	if err != nil {
		panic(err)
	}
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
	variables := map[string]string{}

	code := ""
	code += fmt.Sprintf("%s ", function.Prototype.ReturnType)
	code += fmt.Sprintf("%s(", function.Prototype.Name)

	props := []string{}
	for _, prop := range function.Prototype.Props {
		variables[prop.Name] = prop.Type
		prop := fmt.Sprintf("%s %s", prop.Type, prop.Name)
		props = append(props, prop)
	}
	code += strings.Join(props, ", ")
	code += ") {\n"

	body := ""

	for _, expression := range function.Expressions {
		body += fmt.Sprintf("\t%s;\n", generateExpression(expression, &variables))
	}

	code += body
	code += "}\n\n"

	return code
}

func generateStruct(str *parser.Struct) string {
	customTypes[str.Name] = "struct"

	// Struct def
	code := fmt.Sprintf("typedef struct _%s {\n", str.Name)
	for _, prop := range str.Props {
		code += fmt.Sprintf("\t%s %s;\n", prop.Type, prop.Name)
	}
	code += fmt.Sprintf("} %s;\n\n", str.Name)

	// Make struct
	code += fmt.Sprintf("%s* %s__make() {\n", str.Name, str.Name)
	code += fmt.Sprintf("\t%s* val = malloc(sizeof(%s));\n", str.Name, str.Name)
	code += `	if (!val) {
		printf("Error allocating memory");
		exit(1);
	}
`
	code += "\treturn val;\n"
	code += "}\n\n"

	// Struct functions
	for _, function := range str.Functions {
		code += generateStructFunction(str, function)
	}

	return code
}

func generateStructFunction(str *parser.Struct, function *parser.Function) string {
	copy := *function
	copy.Prototype.Name = fmt.Sprintf("%s__%s", str.Name, function.Prototype.Name)
	structProp := &parser.Prop{
		Name: "self",
		Type: fmt.Sprintf("%s*", str.Name),
	}
	copy.Prototype.Props = append(
		[]*parser.Prop{structProp},
		copy.Prototype.Props...)
	
	return generateFunction(&copy)
}

func generateExpression(expression parser.Expression, functionVariables *map[string]string) string {
	switch expression.ExpressionType() {
	case parser.ExpressionTypeInt:
		exp := expression.(*parser.IntExpression)
		return fmt.Sprintf("%d", exp.Value)
	case parser.ExpressionTypeString:
		exp := expression.(*parser.StringExpression)
		return fmt.Sprintf("String__make(\"%s\")", exp.Value)
	case parser.ExpressionTypeReturn:
		exp := expression.(*parser.ReturnExpression)
		code := "return "
		code += generateExpression(exp.Expression, functionVariables)
		return code
	case parser.ExpressionTypeBinary:
		exp := expression.(*parser.BinaryExpression)
		code := generateExpression(exp.LHS, functionVariables)
		code += fmt.Sprintf(" %s ", cBinaryOperator(exp.Op))
		code += generateExpression(exp.RHS, functionVariables)
		return code
	case parser.ExpressionTypeCall:
		exp := expression.(*parser.CallExpression)

		if isBuiltin(exp.Callee) {
			return generateBuiltin(exp, functionVariables)
		}

		// Struct declaration looks like a call expression,
		// TODO: should catch this in the parser
		if isTypeStruct(exp.Callee) {
			return fmt.Sprintf("%s__make()", exp.Callee)
		}

		// Tradition call
		code := fmt.Sprintf("%s(", exp.Callee)
		paramCode := []string{}
		for _, exp := range exp.Params {
			paramCode = append(
				paramCode,
				generateExpression(exp, functionVariables))
		}
		code += strings.Join(paramCode, ", ")
		code += ")"
		return code
	case parser.ExpressionTypeParen:
		exp := expression.(*parser.ParenExpression)
		code := "("
		code += generateExpression(exp.Expression, functionVariables)
		code += ")"
		return code
	case parser.ExpressionTypeVariableDeclaration:
		exp := expression.(*parser.VariableDeclarationExpression)
		expType := exp.Type
		if isTypeStruct(exp.Type) {
			expType += "*"
		}
		(*functionVariables)[exp.Name] = expType
		code := fmt.Sprintf("%s %s = ", expType, exp.Name)
		code += generateExpression(exp.Expression, functionVariables)
		return code
	case parser.ExpressionTypeVariableAssignment:
		exp := expression.(*parser.VariableAssignmentExpression)
		code := fmt.Sprintf("%s = ", exp.Name)
		code += generateExpression(exp.Expression, functionVariables)
		return code
	case parser.ExpressionTypeVariable:
		exp := expression.(*parser.VariableExpression)
		return exp.Name
	case parser.ExpressionTypeAccessor:
		exp := expression.(*parser.AccessorExpression)
		code := ""

		targetType, ok := (*functionVariables)[exp.Target]
		if !ok {
			msg := fmt.Sprintf("Calling undeclared variable: %s", exp.Target)
			panic(msg)
		}

		if exp.Expression.ExpressionType() == parser.ExpressionTypeCall {
			callExp := exp.Expression.(*parser.CallExpression)
			callExp.Params = append(
				[]parser.Expression{&parser.VariableExpression{Name: exp.Target}},
				callExp.Params...)
			code += fmt.Sprintf(
				"%s__%s",
				strings.Trim(targetType, "*"),
				generateExpression(exp.Expression, functionVariables))
		} else {
			code += fmt.Sprintf(
				"%s->%s",
				exp.Target,
				generateExpression(exp.Expression, functionVariables))
		}

		return code
	default:
		msg := fmt.Sprintf("Unhandled expression type: %v", expression.ExpressionType())
		panic(msg)
	}
}

func isTypeStruct(valueType string) bool {
	val, ok := customTypes[valueType]
	return ok && val == "struct"
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

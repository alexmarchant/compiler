package generator

import (
	"fmt"

	"github.com/alexmarchant/compiler/parser"
)

var varCount = 0

// Generate ...
func Generate(program *parser.Program) string {
	ir := ""
	generateProgram(program, &ir)
	return ir
}

func generateProgram(program *parser.Program, ir *string) {
	generateFunction(program.MainFunction, ir)
}

func generateFunction(function *parser.Function, ir *string) {
	*ir += "define "
	switch function.ReturnType {
	case parser.IntType:
		*ir += "i64 "
	case parser.VoidType:
		*ir += "void "
	default:
		panic("Switch fallthrough")
	}

	*ir += fmt.Sprintf("@%s", function.Name)
	*ir += "() {\n"
	*ir += "entry:\n"

	for _, statement := range function.Statements {
		generateStatement(statement, ir)
	}

	*ir += "}\n"
}

func generateStatement(statement *parser.Statement, ir *string) {
	*ir += "  "

	switch statement.Type {
	case parser.ReturnStatementType:
		generateReturnStatement(statement.ReturnStatement, ir)
	case parser.ExpressionStatementType:
		generateExpressionStatement(statement.ExpressionStatement, ir)
	default:
		panic("Switch fallthrough")
	}
}

func generateReturnStatement(returnStatement *parser.ReturnStatement, ir *string) {
	id := generateExpression(returnStatement.Expression, ir)
	*ir += fmt.Sprintf("ret %s", id)
}

func generateExpressionStatement(expressionStatement *parser.ExpressionStatement, ir *string) string {
	return generateExpression(expressionStatement.Expression, ir)
}

func generateExpression(expression *parser.Expression, ir *string) string {
	switch expression.Type {
	case parser.AdditionExpressionType:
		return generateAdditionExpression(expression.AdditionExpression, ir)
	case parser.LiteralExpressionType:
		return generateLiteralExpression(expression.LiteralExpression, ir)
	default:
		panic("Switch fallthrough")
	}
}

func generateAdditionExpression(additionExpression *parser.AdditionExpression, ir *string) string {
	id := nextID()
	*ir += fmt.Sprintf("%s = add i64 %d %d", id, additionExpression.Left, additionExpression.Right)
	return id
}

func generateLiteralExpression(literalExpression *parser.LiteralExpression, ir *string) string {
	id := nextID()
	switch literalExpression.Type {
	case parser.IntegerLiteralType:
		*ir += fmt.Sprintf("%s = %d", id, literalExpression.IntegerLiteral)
	default:
		panic("Switch fallthrough")
	}
	return id
}

func nextID() string {
	id := fmt.Sprintf("\\%%d", varCount)
	varCount++
	return id
}
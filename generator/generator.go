package generator

import (
	"fmt"

	"github.com/alexmarchant/compiler/parser"
)

var asm = ""

// Generate ...
func Generate(nodes []parser.Node) string {
	for _, node := range nodes {
		generateNode(node)
	}

	return asm
}

func generateNode(node parser.Node) {
	switch node.NodeType() {
	case parser.NodeTypeFunction:
		function := node.(*parser.Function)
		generateFunction(function)
	case parser.NodeTypeExpression:
		expression := node.(parser.Expression)
		generateExpression(expression)
	default:
		panic("Invalid NodeType")
	}
}

func generateFunction(function *parser.Function) {
	asm += fmt.Sprintf("global\t%s\n\n", function.Name)
	asm += fmt.Sprintf("%s:\n", function.Name)

	// Function prologue (start new stack frame)
	asm += "\tpush\t%rbp\n"
	asm += "\tmov\t\t%rsp, %rbp\n"

	for _, expression := range function.Expressions {
		generateExpression(expression)
	}

	if function.ReturnType == nil {
		asm += "\tmov\t\t0, %rax\n"
	}

	// Function epilogue
	asm += "\tmov\t\t%rbp, %rsp\n"
	asm += "\tpop\t\t%rbp\n"
	asm += "\tret\n"
}

func generateExpression(expression parser.Expression) {
	switch expression.ExpressionType() {
	case parser.ExpressionTypeInt:
		intExpression := expression.(*parser.IntExpression)
		asm += fmt.Sprintf("\tmov\t\t%d, %%rax\n", intExpression.Value)
	case parser.ExpressionTypeReturn:
		retExpression := expression.(*parser.ReturnExpression)
		generateExpression(retExpression.Expression)
	}
}

package parser

import "github.com/alexmarchant/compiler/lexer"

// Program is a whole program
type Program struct {
	Functions []Function
}

func parseProgram(tokens []lexer.Token) (Program, error) {
	program := Program{}
	for len(tokens) > 0 {
		function, leftoverTokens, err := parseFunction(tokens)
		if err != nil {
			return program, err
		}
		program.Functions = append(program.Functions, function)
		tokens = leftoverTokens
	}
	return program, nil
}

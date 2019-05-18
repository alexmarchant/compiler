package parser

import "github.com/alexmarchant/compiler/lexer"

// Program is a whole program
type Program struct {
	MainFunction *Function
}

func parseProgram(tokens *[]lexer.Token) (*Program, error) {
	program := Program{}
	function, err := parseFunction(tokens)
	if err != nil {
		return nil, err
	}
	program.MainFunction = function
	return &program, nil
}

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"github.com/alexmarchant/compiler/lexer"
	"github.com/alexmarchant/compiler/parser"
)

func main() {
	if len(os.Args) < 2 {
		panic("File arg required")
	}
	filepath := os.Args[1]

	dat, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	source := string(dat)
	tokens := lexer.Lex(source)
	fmt.Printf("%+v\n", tokens)

	ast, err := parser.Parse(tokens)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", ast)
}

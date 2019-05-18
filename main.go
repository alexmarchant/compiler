package main

import (
	"io/ioutil"
	"os"
	"fmt"

	"github.com/alexmarchant/compiler/lexer"
	"github.com/alexmarchant/compiler/parser"
	"github.com/alexmarchant/compiler/generator"
	"github.com/sanity-io/litter"
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
	// litter.Dump(tokens)

	program, err := parser.Parse(tokens)
	if err != nil {
		panic(err)
	}
	litter.Dump(program)

	ir := generator.Generate(program)
	fmt.Print(ir)
}

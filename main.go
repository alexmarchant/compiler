package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/alexmarchant/compiler/generator"
	"github.com/alexmarchant/compiler/lexer"
	"github.com/alexmarchant/compiler/parser"
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
	litter.Dump(tokens)

	fmt.Println("--AST--")
	nodes := parser.Parse(tokens)
	litter.Dump(nodes)

	fmt.Println("\n--CODE--")
	code := generator.Generate(nodes)
	fmt.Print(code)

	// Generate binary
	ioutil.WriteFile("./out.c", []byte(code), 0644)
	cmd := exec.Command("sh", "-c", "clang out.c runtime/*.c -o out")
	var errLog bytes.Buffer
	cmd.Stderr = &errLog
	err = cmd.Run()
	if len(errLog.String()) > 0 {
		fmt.Println("\n--COMPILATION--")
		fmt.Printf(errLog.String())
	}
	if err != nil {
		panic(err)
	}
	os.Remove("out.c")
}

package main

import (
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
	// litter.Dump(tokens)

	fmt.Println("--AST--")
	nodes := parser.Parse(tokens)
	litter.Dump(nodes)

	fmt.Println("\n--ASM--")
	asm := generator.Generate(nodes)
	fmt.Print(asm)

	// Generate binary
	ioutil.WriteFile("./out.asm", []byte(asm), 0644)
	cmd := exec.Command("sh", "-c", "nasm -f macho64 out.asm")
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
	cmd = exec.Command("sh", "-c", "ld -macosx_version_min 10.7.0 -lSystem -o out out.o")
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}

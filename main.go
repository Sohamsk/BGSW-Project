package main

import (
	"bosch/listener"
	"bosch/parser"
	"log"
	"os"

	"github.com/antlr4-go/antlr/v4"
)

func main() {
	input, _ := antlr.NewFileStream(os.Args[1])

	lexer := parser.NewVisualBasic6Lexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.NewVisualBasic6Parser(stream)
	p.BuildParseTrees = true
	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))

	tree := p.StartRule()
	antlr.ParseTreeWalkerDefault.Walk(listener.NewTreeShapeListener(), tree)

	f, err := os.Create("op.txt")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	f.Seek(0, 0)
	f.WriteString(tree.GetText())

	f.Close()
}

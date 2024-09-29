package main

import (
	"bosch/listener"
	"bosch/parser"
	"bufio"
	"log"
	"os"

	"github.com/antlr4-go/antlr/v4"
)

func main() {
	input, err := antlr.NewFileStream(os.Args[1])

	if err != nil {
		log.Panic("File error")
	}

	lexer := parser.NewVisualBasic6Lexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.NewVisualBasic6Parser(stream)
	p.BuildParseTrees = true
	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	tree := p.StartRule()

	f, err := os.Create("op.json")
	if err != nil {
		log.Panic(err)
	}
	f.Seek(0, 0)
	w := bufio.NewWriter(f)
	w.WriteString("{\"Name\":\"whatever\", \"ModuleType\": \"something\"")
	antlr.ParseTreeWalkerDefault.Walk(listener.NewTreeShapeListener(w), tree)
	w.WriteString("}")
	w.Flush()
	f.Close()
}

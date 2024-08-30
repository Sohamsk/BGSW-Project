package main

import (
	"bosch/parser"
	"fmt"
	"log"
	"os"

	"github.com/antlr4-go/antlr/v4"
)

type TreeShapeListener struct {
	*parser.BaseVisualBasic6ParserListener
}

func NewTreeShapeListener() *TreeShapeListener {
	return new(TreeShapeListener)
}

func (treeShapeListener *TreeShapeListener) EnterEveryRule(ctx antlr.ParserRuleContext) {
	fmt.Println(ctx.GetText())
}

func main() {
	input, _ := antlr.NewFileStream(os.Args[1])

	lexer := parser.NewVisualBasic6Lexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.NewVisualBasic6Parser(stream)
	p.BuildParseTrees = true
	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))

	tree := p.StartRule()
	antlr.ParseTreeWalkerDefault.Walk(NewTreeShapeListener(), tree)

	f, err := os.Create("op.txt")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	f.Seek(0, 0)
	f.WriteString(tree.GetText())

	f.Close()
}

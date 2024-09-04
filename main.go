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

func (s *TreeShapeListener) VisitTerminal(node antlr.TerminalNode) {
	fmt.Println(node.GetSymbol())
}

func (s *TreeShapeListener) ExitStartRule(ctx *parser.StartRuleContext) {
	fmt.Println("Start Rule")
}

func (s *TreeShapeListener) EnterVariableStmt(ctx *parser.VariableStmtContext) {
	fmt.Println("DIM")
}

// func (treeShapeListener *TreeShapeListener) EnterWhileWendStmt(ctx parser.VariableStmtContext) {
// 	fmt.Println("Entered while")
// }

// func (s *TreeShapeListener) EnterEveryRule(ctx antlr.ParserRuleContext) {
// 	fmt.Printf("NEWLINE: %s\n", ctx.)
// }

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

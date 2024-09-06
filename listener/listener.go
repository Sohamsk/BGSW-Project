package listener

import (
	"bosch/parser"
	"bufio"
	"fmt"

	"github.com/antlr4-go/antlr/v4"
)

type TreeShapeListener struct {
	*parser.BaseVisualBasic6ParserListener
	writer *bufio.Writer
}

func NewTreeShapeListener(writer *bufio.Writer) *TreeShapeListener {
	l := new(TreeShapeListener)
	l.writer = writer
	return l
}

func (s *TreeShapeListener) VisitTerminal(node antlr.TerminalNode) {
	s.writer.WriteString(node.GetSymbol().String())
	fmt.Println(node.GetSymbol())
}

func (s *TreeShapeListener) EnterStartRule(ctx *parser.StartRuleContext) {
	fmt.Println("Enter Start Rule")
}

func (s *TreeShapeListener) ExitStartRule(ctx *parser.StartRuleContext) {
	fmt.Println("Exit Start Rule")
}

func (s *TreeShapeListener) EnterDoLoopStmt(ctx *parser.DoLoopStmtContext) {
	fmt.Println("Enter do statement")
}

func (s *TreeShapeListener) ExitDoLoopStmt(ctx *parser.DoLoopStmtContext) {
	fmt.Println("Exit do statement")
}

// func (treeShapeListener *TreeShapeListener) EnterWhileWendStmt(ctx parser.VariableStmtContext) {
// 	fmt.Println("Entered while")
// }

// func (s *TreeShapeListener) EnterEveryRule(ctx antlr.ParserRuleContext) {
// 	fmt.Printf("NEWLINE: %s\n", ctx.)
// }

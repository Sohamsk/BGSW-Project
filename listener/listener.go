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
	// 218 token type is variable name
	if node.GetSymbol().GetTokenType() == 218 {
		fmt.Println(node.GetText())
		// figure something out to place the variable name and their types somewhere
		// i will figure something out to make the it so that the varibale declaration statement state is stored
		// i am thinking of using a stack to save the state to use later while constructing the tree
	}
	// this lists out all the tokens in the code
	s.writer.WriteString(node.GetSymbol().String() + "\n")
}

func (s *TreeShapeListener) EnterStartRule(ctx *parser.StartRuleContext) {
	fmt.Println("Enter Start Rule")
}

func (s *TreeShapeListener) ExitStartRule(ctx *parser.StartRuleContext) {
	fmt.Println("Exit Start Rule")
}

func (s *TreeShapeListener) EnterVariableStmt(ctx *parser.VariableStmtContext) {
	// for _, v := range ctx.GetChildren() {
	// 	fmt.Println(v)
	// }
	fmt.Println(ctx.GetStart().GetTokenType())
}

func (s *TreeShapeListener) ExitVariableStmt(ctx *parser.VariableStmtContext) {
	fmt.Println("Exit DIM")
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

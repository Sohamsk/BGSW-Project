package listener

import (
	"bosch/parser"
	"bufio"
	"fmt"
	"bosch/stack"

	"github.com/antlr4-go/antlr/v4"
)

type TreeShapeListener struct {
	*parser.BaseVisualBasic6ParserListener
	writer *bufio.Writer
	stack stack.Stack
}

func NewTreeShapeListener(writer *bufio.Writer) *TreeShapeListener {
	l := new(TreeShapeListener)
	l.writer = writer
	l.stack = *stack.InitStack()
	return l
}

func (s *TreeShapeListener) VisitTerminal(node antlr.TerminalNode) {
	// 218 token type is variable name
//	if node.GetSymbol().GetTokenType() == 218 {
//		fmt.Println(node.GetText())
		// figure something out to place the variable name and their types somewhere
		// i will figure something out to make the it so that the varibale declaration statement state is stored
		// i am thinking of using a stack to save the state to use later while constructing the tree

//	}
	// this lists out all the tokens in the code
//	s.writer.WriteString(node.GetSymbol().String())
	fmt.Println(node.GetSymbol().String())
}

func (s *TreeShapeListener) EnterStartRule(ctx *parser.StartRuleContext) {
	s.writer.WriteString("{")
}

func (s *TreeShapeListener) ExitStartRule(ctx *parser.StartRuleContext) {
	s.writer.WriteString("}")
}

func (s *TreeShapeListener) EnterVariableSubStmt(ctx *parser.VariableSubStmtContext) {
	nodes := ctx.GetChildren()
	s.writer.WriteString("\"DeclareVariable\": {")
	s.writer.WriteString("\"Identifier\": \"" + nodes[0].(antlr.ParseTree).GetText() + "\",")
	s.writer.WriteString("\"Type\": {")
	if (len(nodes) == 3) {
		s.writer.WriteString("\"Static\": true,")
		s.writer.WriteString("\"type\": \"" + nodes[2].GetChild(2).(antlr.RuleNode).GetText() + "\"")
	} else {
		s.writer.WriteString("\"Static\": false")
	}
	s.writer.WriteString("}}")
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

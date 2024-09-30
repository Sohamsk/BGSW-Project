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
	fmt.Println(node.GetSymbol().String())
}

func (s *TreeShapeListener) EnterStartRule(ctx *parser.StartRuleContext) {
	s.writer.WriteString("\"body\": [")
}

func (s *TreeShapeListener) ExitStartRule(ctx *parser.StartRuleContext) {
	s.writer.Flush()
	s.writer.WriteString("]")
}

func (s *TreeShapeListener) EnterVariableSubStmt(ctx *parser.VariableSubStmtContext) {
	nodes := ctx.GetChildren()
	s.writer.WriteString("{\"RuleType\": \"DeclareVariable\",")
	s.writer.WriteString("\"Identifier\": \"" + nodes[0].(antlr.ParseTree).GetText() + "\",")
	if (len(nodes) == 3) {
		s.writer.WriteString("\"Type\": \"" + nodes[2].GetChild(2).(antlr.RuleNode).GetText() + "\"")
	} else {
		s.writer.WriteString("\"Type\": \"VARIANT\"")
	}
	s.writer.WriteString("},")
}


func (s *TreeShapeListener) EnterLetStmt(ctx *parser.LetStmtContext) {
	nodes := ctx.GetChildren()
	flip := false
	var lhs, rhs string
	s.writer.WriteString("{\"RuleType\": \"LetStatement\",")
	for _, node := range(nodes) {
		switch node.(type) {		// we WILL have to handle function calls some how
		case antlr.RuleNode:
			if (!flip) {
				lhs += node.(antlr.RuleNode).GetText()
			} else {
				rhs += node.(antlr.RuleNode).GetText()
			}
		case antlr.TerminalNode:
			sym := node.(antlr.TerminalNode).GetText()
			if (sym != " ") {
				if (sym == "=" || sym == "+=" || sym == "-=" ) {
					flip = !flip
					s.writer.WriteString("\"Operation\":\""+ sym +"\",")
				} else {
					if (!flip) {
						lhs += sym
					} else {
						rhs += sym
					}
				}
			}
		}
	}
	s.writer.WriteString("\"left\":\"" + lhs + "\",")
	s.writer.WriteString("\"Right\":\""+ rhs + "\"")
	s.writer.WriteString("},")
}

func (s *TreeShapeListener) EnterDoLoopStmt(ctx *parser.DoLoopStmtContext) {
	fmt.Println("Enter do statement")
}

func (s *TreeShapeListener) ExitDoLoopStmt(ctx *parser.DoLoopStmtContext) {
	fmt.Println("Exit do statement")
}


package listener

import (
	"bosch/parser"
	"bosch/stack"
	"bufio"
	"fmt"

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
	//fmt.Println(node.GetSymbol().String())
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

func handleLetExpression(nodes []antlr.Tree) {
	if (len(nodes) == 0) {
		return
	}
	fmt.Println("\"Left\": {")
	outer:
	for i, node := range(nodes) {
		switch node.(type) {
		case antlr.TerminalNode:
			if (node.(antlr.TerminalNode).GetText() != " ") {
				fmt.Println(node.(antlr.TerminalNode).GetText())
				fmt.Println("}\"Right\": {")
				fmt.Println("}")
				handleLetExpression(nodes[i + 1:])
				break outer
			}
		case antlr.RuleNode:
			if (node.(antlr.RuleNode).GetChildCount() == 1) {
				fmt.Println(node.(antlr.RuleNode).GetText())
			} else {
				handleLetExpression(node.(antlr.RuleNode).GetChildren())
			}
		}
	}
	fmt.Println("}")
}
func (s *TreeShapeListener) EnterLetStmt(ctx *parser.LetStmtContext) {
	//fmt.Println(parser.VisualBasic6ParserParserStaticData.RuleNames[ctx.GetRuleIndex()])
	nodes := ctx.GetChildren()
	fmt.Println("{")
	handleLetExpression(nodes)
	fmt.Println("__________")
}

func (s *TreeShapeListener) EnterDoLoopStmt(ctx *parser.DoLoopStmtContext) {
	fmt.Println("Enter do statement")
}

func (s *TreeShapeListener) ExitDoLoopStmt(ctx *parser.DoLoopStmtContext) {
	fmt.Println("Exit do statement")
}


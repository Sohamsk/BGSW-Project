package listener

import (
	"bosch/parser"
	"github.com/antlr4-go/antlr/v4"
)

func (s *TreeShapeListener) EnterForNextStmt(ctx *parser.ForNextStmtContext) {
	nodes := ctx.GetChildren()
	var initialization, condition, from string

	// Start the ForLoopStatement object
	s.writer.WriteString("{\"RuleType\":\"ForLoopStatement\",")

	for i, node := range nodes {
		switch n := node.(type) {
		case antlr.TerminalNode:
			// Capture initialization
			if n.GetText() == "For" && i+2 < len(nodes) {
				initialization = nodes[i+2].(antlr.ParseTree).GetText()
			} else if n.GetText() == "To" && i+2 < len(nodes) {
				condition = nodes[i+2].(antlr.ParseTree).GetText()
			} else if n.GetText() == "=" && i+1 < len(nodes) {
				from = nodes[i+1].(antlr.ParseTree).GetText()
			}
		}
	}
	s.writer.WriteString("\"Initialization\": \"" + initialization + "\",")
	s.writer.WriteString("\"Start\": \"" + from + "\",")
	s.writer.WriteString("\"End\": \"" + condition + "\",\"Body\": [")
}

func (s *TreeShapeListener) ExitForNextStmt(ctx *parser.ForNextStmtContext) {
	s.exitContext()
	s.writer.WriteString("}}")
}

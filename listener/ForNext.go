package listener

import (
	"bosch/parser"

	"github.com/antlr4-go/antlr/v4"
)

func (s *TreeShapeListener) EnterForNextStmt(ctx *parser.ForNextStmtContext) {
	nodes := ctx.GetChildren()
	var initialization, condition, from, step string

	// Start the ForLoopStatement object
	s.writer.WriteString("{\"RuleType\":\"ForNextStmt\",")

	for i, node := range nodes {
		switch n := node.(type) {
		case antlr.TerminalNode:
			// Capture initialization
			if n.GetText() == "For" && i+2 < len(nodes) {
				initialization = nodes[i+2].(antlr.ParseTree).GetText()
			} else if n.GetText() == "To" && i+2 < len(nodes) {
				condition = nodes[i+2].(antlr.ParseTree).GetText()
			} else if n.GetText() == "Step" && i+1 < len(nodes) {
				step = nodes[i+2].(antlr.ParseTree).GetText()
			} else if n.GetText() == "=" && i+1 < len(nodes) {
				from = nodes[i+2].(antlr.ParseTree).GetText()
			}
		}
	}
	s.writer.WriteString("\"Initialization\": \"" + initialization + "\",")
	s.writer.WriteString("\"Start\": \"" + from + "\",")
	s.writer.WriteString("\"End\": \"" + condition + "\",")
	s.writer.WriteString("\"Step\": \"" + step + "\",")
	s.writer.WriteString("\"Body\": [")
}

func (s *TreeShapeListener) ExitForNextStmt(ctx *parser.ForNextStmtContext) {
	s.exitContext()
	s.writer.WriteString("}")
}

package listener

import (
	"bosch/parser"
	"github.com/antlr4-go/antlr/v4"
)

func (s *TreeShapeListener) EnterPrintStmt(ctx *parser.PrintStmtContext) {
	nodes := ctx.GetChildren()
	//    s.writer.WriteString("{\"CommentedUI\":\"" + ctx.GetText() + "\"},")
	for _, node := range nodes {
		switch node := node.(type) {
		case antlr.RuleNode:
			text := node.GetText()
			s.writer.WriteString("{\"RuleType\":\"PrintStmt\", \"Data\": \"" + text + "\"},")
		}
	}
}

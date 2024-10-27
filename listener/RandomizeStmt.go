package listener

import (
	"bosch/parser"
	"github.com/antlr4-go/antlr/v4"
)

func (s *TreeShapeListener) EnterRandomizeStmt(ctx *parser.RandomizeStmtContext) {
	nodes := ctx.GetChildren()
	var seed string
	s.writer.WriteString("{\"RuleType\":\"RandomizeStatement\", ")

	if len(nodes) > 1 {
		seed = nodes[2].(antlr.ParseTree).GetText()
		s.writer.WriteString("\"Seed\": \"" + seed + "\"}")
	} else {
		s.writer.WriteString("\"Seed\": \"None\"}")
	}
	s.writer.WriteString("")
}

func (s *TreeShapeListener) ExitRandomizeStmt(ctx *parser.RandomizeStmtContext) {
	s.writer.WriteString("}")
}

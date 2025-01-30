package listener

import (
	"bosch/parser"
	"github.com/antlr4-go/antlr/v4"
)

func (s *TreeShapeListener) EnterVariableSubStmt(ctx *parser.VariableSubStmtContext) {
	nodes := ctx.GetChildren()
	s.writer.WriteString("{\"RuleType\": \"DeclareVariable\",")
	s.writer.WriteString("\"Identifier\": \"" + nodes[0].(antlr.ParseTree).GetText() + "\",")
	if len(nodes) == 3 {
		s.writer.WriteString("\"Type\": \"" + nodes[2].GetChild(2).(antlr.RuleNode).GetText() + "\"")
		s.SymTab[nodes[0].(antlr.ParseTree).GetText()] = nodes[2].GetChild(2).(antlr.RuleNode).GetText()
	} else {
		s.writer.WriteString("\"Type\": \"VARIANT\"")
		s.SymTab[nodes[0].(antlr.ParseTree).GetText()] = "VARIANT"
	}
	s.writer.WriteString("},")
}

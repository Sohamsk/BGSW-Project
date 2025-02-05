package listener

import (
	"bosch/parser"
	"strings"

	"github.com/antlr4-go/antlr/v4"
)

func (s *TreeShapeListener) EnterForEachStmt(ctx *parser.ForEachStmtContext) {
	nodes := ctx.GetChildren()
	var variable, collection string

	s.writer.WriteString("{\"RuleType\":\"ForEachStatement\",")

	for i, node := range nodes {
		switch n := node.(type) {
		case antlr.TerminalNode:
			if strings.ToLower(n.GetText()) == "each" && i+1 < len(nodes) {
				variable = nodes[i+2].(antlr.ParseTree).GetText()
			} else if strings.ToLower(n.GetText()) == "in" && i+1 < len(nodes) {
				collection = nodes[i+2].(antlr.ParseTree).GetText()
			}
		}
	}

	// Write JSON structure
	s.writer.WriteString("\"Element\": \"" + variable + "\",")
	s.writer.WriteString("\"Collection\": \"" + collection + "\",")
	s.writer.WriteString("\"Body\": [")
}

func (s *TreeShapeListener) ExitForEachStmt(ctx *parser.ForEachStmtContext) {
	s.exitContext()
	s.writer.WriteString("}")
}

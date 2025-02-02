package listener

import (
	"bosch/parser"
	"bufio"
	"bytes"
	"strings"
)

// TODO:Refactor this and IfThenElse statement later with a single handler

func (s *TreeShapeListener) EnterMacroIfBlockStmt(ctx *parser.MacroIfBlockStmtContext) {
	// These are actually Compiler Directives
	nodes := ctx.GetChildren()
	s.writer.WriteString("{\"RuleType\":\"MacroIfBlock\",")
	for _, node := range nodes {
		switch node := node.(type) {
		case *parser.IfConditionStmtContext:
			s.writer.WriteString("\"IsBlock\":true,")
			var buffer bytes.Buffer
			writer := bufio.NewWriter(&buffer)
			handleLetExpression(node.GetChild(0).GetChildren(), writer)
			writer.Flush()
			s.writer.WriteString("\"Condition\": [" + strings.Trim(buffer.String(), ",") + "],")
			s.writer.WriteString("\"IfBlock\": [")
		}
	}
}

func (s *TreeShapeListener) ExitMacroIfBlockStmt(ctx *parser.MacroIfBlockStmtContext) {
	s.exitContext()
	s.writer.WriteString("},")
}

func (s *TreeShapeListener) EnterMacroElseIfBlockStmt(ctx *parser.MacroElseIfBlockStmtContext) {
	s.writer.WriteString("{\"RuleType\":\"MacroElseIf\",")
	nodes := ctx.GetChildren()
	for _, node := range nodes {
		switch node := node.(type) {
		case *parser.IfConditionStmtContext:
			var buffer bytes.Buffer
			writer := bufio.NewWriter(&buffer)
			handleLetExpression(node.GetChild(0).GetChildren(), writer)
			writer.Flush()
			s.writer.WriteString("\"Condition\": [" + strings.Trim(buffer.String(), ",") + "],")
			s.writer.WriteString("\"ElseIfBlock\": [")
		}
	}
}

func (s *TreeShapeListener) ExitMacroElseIfBlockStmt(ctx *parser.MacroElseIfBlockStmtContext) {
	s.exitContext()
	s.writer.WriteString("},")
}

func (s *TreeShapeListener) EnterMacroElseBlockStmt(ctx *parser.MacroElseBlockStmtContext) {
	s.writer.WriteString("{\"RuleType\":\"MacroElseBlock\",")
	s.writer.WriteString("\"Body\": [")

}

func (s *TreeShapeListener) ExitMacroElseBlockStmt(ctx *parser.MacroElseBlockStmtContext) {
	s.exitContext()
	s.writer.WriteString("},")
}

func (s *TreeShapeListener) EnterMacroEndIfStmt(c *parser.MacroEndIfStmtContext) {
	s.writer.WriteString("{\"RuleType\":\"EndIf\"}")
}

package listener

import (
	"bosch/parser"
	"bufio"
	"bytes"
	"strings"
)

func (s *TreeShapeListener) EnterIfBlockStmt(ctx *parser.IfBlockStmtContext) {
	nodes := ctx.GetChildren()
	s.writer.WriteString("{\"RuleType\":\"IfThenElse\",")
	for _, node := range nodes {
		switch node := node.(type) {
		case *parser.IfConditionStmtContext:
			s.writer.WriteString("\"IsBlock\":\"True\",")
			var buffer bytes.Buffer
			writer := bufio.NewWriter(&buffer)
			handleLetExpression(node.GetChild(0).GetChildren(), writer)
			writer.Flush()
			s.writer.WriteString("\"Condition\": [" + strings.Trim(buffer.String(), ",") + "],")
			s.writer.WriteString("\"IfBlock\": [")
		}
	}
}

func (s *TreeShapeListener) ExitIfBlockStmt(ctx *parser.IfBlockStmtContext) {
	s.exitContext()
	s.writer.WriteString("},")
}

func (s *TreeShapeListener) EnterIfElseIfBlockStmt(ctx *parser.IfElseIfBlockStmtContext) {
	s.writer.WriteString("{\"RuleType\":\"ElseIf\",")
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

func (s *TreeShapeListener) ExitIfElseIfBlockStmt(ctx *parser.IfElseIfBlockStmtContext) {
	s.exitContext()
	s.writer.WriteString("},")
}

func (s *TreeShapeListener) EnterIfElseBlockStmt(ctx *parser.IfElseBlockStmtContext) {
	s.writer.WriteString("{\"RuleType\":\"ElseBlock\",")
	s.writer.WriteString("\"Body\": [")

}

func (s *TreeShapeListener) ExitIfElseBlockStmt(ctx *parser.IfElseBlockStmtContext) {
	s.exitContext()
	s.writer.WriteString("},")
}

func (s *TreeShapeListener) ExitBlockIfThenElseStmt(ctx *parser.BlockIfThenElseContext) {
	s.exitContext()
	s.writer.WriteString("},")
}

package listener

import (
	"bosch/parser"
)

func (s *TreeShapeListener) EnterIfBlockStmt(ctx *parser.IfBlockStmtContext) {
	nodes := ctx.GetChildren()
	s.writer.WriteString("{\"RuleType\":\"IfThenElse\",")
	for _, node := range nodes {
		switch node := node.(type) {
		case *parser.IfConditionStmtContext:
			s.writer.WriteString("\"IsBlock\":\"True\",")
			IfCondition := node.GetText()
			s.writer.WriteString("\"Condition\":\"{" + IfCondition + "}\",")
			s.writer.WriteString("\"IfBlock\": ")
		}
	}
}

func (s *TreeShapeListener) ExitIfBlockStmt(ctx *parser.IfBlockStmtContext) {
	s.writer.WriteString("},")
}

func (s *TreeShapeListener) EnterIfElseIfBlockStmt(ctx *parser.IfElseIfBlockStmtContext) {
	s.writer.WriteString("{\"RuleType\":\"ElseIf\",")
	nodes := ctx.GetChildren()
	for _, node := range nodes {
		switch node := node.(type) {
		case *parser.IfConditionStmtContext:
			IfCondition := node.GetText()
			s.writer.WriteString("\"Condition\":\"{" + IfCondition + "}\",")
			s.writer.WriteString("\"ElseIfBlock\": ")
		}
	}
}

func (s *TreeShapeListener) ExitIfElseIfBlockStmt(ctx *parser.IfElseIfBlockStmtContext) {
	s.writer.WriteString("},")
}

func (s *TreeShapeListener) EnterIfElseBlockStmt(ctx *parser.IfElseBlockStmtContext) {
	s.writer.WriteString("{\"ElseBlock\": ")

}

func (s *TreeShapeListener) ExitIfElseBlockStmt(ctx *parser.IfElseBlockStmtContext) {
	s.writer.WriteString("}")
}

func (s *TreeShapeListener) ExitBlockIfThenElseStmt(ctx *parser.BlockIfThenElseContext) {
	s.writer.WriteString("}")
}

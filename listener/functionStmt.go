package listener

import (
	"bosch/parser"
)

// EnterFunctionStmt is called when production functionStmt is entered.
func (s *TreeShapeListener) EnterFunctionStmt(ctx *parser.FunctionStmtContext) {
	s.handleFuncLikeDecl(ctx, "FuncStatement", false, "")
}

// ExitFunctionStmt is called when production functionStmt is exited.
func (s *TreeShapeListener) ExitFunctionStmt(ctx *parser.FunctionStmtContext) {
	s.exitContext()
	s.writer.WriteString("},")
}

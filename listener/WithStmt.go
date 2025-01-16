package listener

import (
	"bosch/converter"
	"bosch/parser"
	"encoding/json"

	"github.com/antlr4-go/antlr/v4"
)

func (s *TreeShapeListener) EnterWithStmt(ctx *parser.WithStmtContext) {
	with := converter.WithStmt{}
	with.RuleType = "WithStatement"
	r := ctx.GetChild(2)
	with.Object = r.(antlr.RuleContext).GetText()
	jsonData, err := json.Marshal(with)
	if err != nil {
		panic(err)
	}
	s.writer.WriteString(string(jsonData)[:len(jsonData)-5] + "[")
}

func (s *TreeShapeListener) ExitWithStmt(ctx *parser.WithStmtContext) {
	s.exitContext()
	s.writer.WriteString("},")
}

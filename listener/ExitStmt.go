package listener

import (
	"bosch/converter"
	"bosch/parser"
	"encoding/json"

	"github.com/antlr4-go/antlr/v4"
)

func (s *TreeShapeListener) EnterExit_Function(ctx *parser.Exit_FunctionContext) {
	rules := parser.VisualBasic6ParserParserStaticData.RuleNames
	var parent antlr.Tree
	parent = ctx
	for parent != nil && rules[parent.(antlr.RuleContext).GetRuleIndex()] != "functionStmt" {
		parent = parent.GetParent()
	}
	if parent == nil {
		panic("Syntax error")
	}
	ir := converter.ReturnStmt{}
	ir.RuleType = "ReturnStatement"
	ir.ReturnVariableName = parent.GetChild(2).(*parser.AmbiguousIdentifierContext).GetText()

	jsonData, err := json.Marshal(ir)
	if err != nil {
		panic(err)
	}
	s.writer.WriteString(string(jsonData) + ",")
}

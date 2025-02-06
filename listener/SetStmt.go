package listener

import (
	"bosch/converter/models"
	"bosch/parser"
	"encoding/json"

	"github.com/antlr4-go/antlr/v4"
)

func (s *TreeShapeListener) EnterSetStmt(ctx *parser.SetStmtContext) {
	rules := parser.VisualBasic6ParserParserStaticData.RuleNames
	children := ctx.GetChildren()
	var id string
	var class string
	var isNew bool

	for _, node := range children {
		node_val, ok := node.(antlr.RuleContext)
		if !ok {
			continue
		}
		switch rules[node_val.GetRuleIndex()] {
		case "implicitCallStmt_InStmt":
			id = node_val.GetText()
		case "valueStmt":
			if node.GetChildCount() == 1 {
				isNew = false
				class = node_val.GetText()
			} else {
				isNew = true
				child := node_val.GetChild(2)
				class = child.(antlr.RuleContext).GetText()
			}
		}
	}

	var set models.SetStmt
	set.RuleType = "SetStatement"
	set.IsNew = isNew
	set.Identifier = id
	set.Class, _ = json.Marshal(class)

	jsonData, _ := json.Marshal(set)
	s.writer.WriteString(string(jsonData) + ",")
	if isNew {
		s.SymTab[id] = class
	}
}

package listener

import (
	"bosch/converter/models"
	"bosch/parser"
	"encoding/json"

	"github.com/antlr4-go/antlr/v4"
)

// EnterSetStmt is called when production setStmt is entered.
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
				// TODO: handle function calls that deliver the object
				class = child.(antlr.RuleContext).GetText()
			}
		}
	}
	print(string(id + " " + class))
	if isNew {
		println(" true")
	} else {
		println(" false")
	}
	var set models.SetStmt
	set.RuleType = "SetStatement"
	set.IsNew = isNew
	set.Identifier = id
	set.Class, _ = json.Marshal(class)

	jsonData, _ := json.Marshal(set)
	s.writer.WriteString(string(jsonData) + ",")

}

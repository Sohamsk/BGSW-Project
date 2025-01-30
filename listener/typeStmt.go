package listener

import (
	"bosch/converter/models"
	"bosch/parser"
	"encoding/json"

	"github.com/antlr4-go/antlr/v4"
)

func (s *TreeShapeListener) EnterTypeStmt(ctx *parser.TypeStmtContext) {
	nodes := ctx.GetChildren()
	typestmt := models.TypeStmt{}
	typestmt.Rule.RuleType = "TypeStmtRule"
	for _, node := range nodes {
		switch node.(type) {
		case *parser.TypeStmt_ElementContext:
			tempDim := models.Dim{}
			nodes := node.GetChildren()
			for _, child := range nodes {
				switch child.(type) {
				case *parser.AmbiguousIdentifierContext:
					tempDim.Identifier = child.(antlr.ParseTree).GetText()
				case *parser.AsTypeClauseContext:
					tempDim.Type = child.GetChild(2).(antlr.ParseTree).GetText()

				}
			}
			if tempDim.Type == "" {
				tempDim.Type = "variant"
			}
			typeElementJSON, ok := json.Marshal(tempDim)
			if ok != nil {
				panic("Error in typeStmt")
			}
			typestmt.TypeElements = append(typestmt.TypeElements, typeElementJSON)
		case *parser.PublicPrivateVisibilityContext:
			typestmt.Visibility = node.(antlr.ParseTree).GetText()
		case *parser.AmbiguousIdentifierContext:
			typestmt.Name = node.(antlr.ParseTree).GetText()
		}
	}

	jsonData, err := json.Marshal(typestmt)
	if err != nil {
		panic(err)
	}
	s.writer.WriteString(string(jsonData) + ",")

}

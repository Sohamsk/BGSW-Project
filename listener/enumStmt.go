package listener

import (
	"bosch/converter/models"
	"bosch/parser"
	"encoding/json"
	"strings"

	"github.com/antlr4-go/antlr/v4"
)

func (s *TreeShapeListener) EnterEnumerationStmt(ctx *parser.EnumerationStmtContext) {
	nodes := ctx.GetChildren()
	enum := models.EnumStmt{}
	enum.Rule.RuleType = "EnumerationRule"
	for _, node := range nodes {
		switch node.(type) {
		case *parser.EnumerationStmt_ConstantContext:
			enumConstant := strings.TrimSpace(node.(antlr.ParseTree).GetText())
			enumConstant = strings.TrimSuffix(enumConstant, "\n\t")
			enumConstant = strings.TrimSuffix(enumConstant, "\n")
			enum.EnumValues = append(enum.EnumValues, enumConstant)
		case *parser.PublicPrivateVisibilityContext:
			enum.Visibility = node.(antlr.ParseTree).GetText()
		case *parser.AmbiguousIdentifierContext:
			enum.Name = node.(antlr.ParseTree).GetText()
		}
	}

	jsonData, err := json.Marshal(enum)
	if err != nil {
		panic(err)
	}
	s.writer.WriteString(string(jsonData) + ",")

}

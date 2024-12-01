package listener

import (
	"bosch/parser"
	"fmt"

	"github.com/antlr4-go/antlr/v4"
)

func (s *TreeShapeListener) EnterPrintStmt(ctx *parser.PrintStmtContext) {
	var outputList []string

	for _, child := range ctx.GetChildren() {

		if ruleNode, ok := child.(antlr.RuleNode); ok {

			traverseAndCollectIdentifiers(ruleNode, &outputList)
		}
	}
	for _, identifier := range outputList {
		if identifier != "" {
			s.writer.WriteString(fmt.Sprintf("{\"RuleType\":\"PrintStmt\", \"Data\": \"%s\"},", identifier))
		}
	}
}

// Helper function to traverse and collect identifiers
func traverseAndCollectIdentifiers(node antlr.Tree, outputList *[]string) {
	switch t := node.(type) {
	case antlr.TerminalNode:
		text := t.GetText()
		if text != "Print" && text != "," && text != " " {
			*outputList = append(*outputList, text)
		}
	case antlr.RuleNode:
		for _, child := range t.GetChildren() {
			traverseAndCollectIdentifiers(child, outputList)
		}
	}
}

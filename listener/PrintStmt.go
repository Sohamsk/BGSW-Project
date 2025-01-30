package listener

import (
	"bosch/parser"
	"encoding/json"
	"strings"

	"github.com/antlr4-go/antlr/v4"
)

func (s *TreeShapeListener) EnterPrintStmt(ctx *parser.PrintStmtContext) {
	var outputList []string

	// Traverse children to collect arguments
	for _, child := range ctx.GetChildren() {
		if ruleNode, ok := child.(antlr.RuleNode); ok {
			traverseAndCollectIdentifiers(ruleNode, &outputList)
		} else if terminalNode, ok := child.(antlr.TerminalNode); ok {
			text := strings.TrimSpace(terminalNode.GetText())
			if text != "Print" && text != "," && text != "" {
				outputList = append(outputList, text)
			}
		}
	}

	// Ensure proper JSON formatting using json.Marshal
	jsonData, err := json.Marshal(map[string]interface{}{
		"RuleType": "PrintStmt",
		"Data":     outputList,
	})
	if err != nil {
		panic("Error: Failed to marshal PrintStmt JSON")
	}

	s.writer.WriteString(string(jsonData))
}

// Helper function to traverse and collect identifiers
func traverseAndCollectIdentifiers(node antlr.Tree, outputList *[]string) {
	switch t := node.(type) {
	case antlr.TerminalNode:
		text := strings.TrimSpace(t.GetText()) // Trim whitespace
		if text != "Print" && text != "," && text != "" {
			*outputList = append(*outputList, text)
		}
	case antlr.RuleNode:
		for _, child := range t.GetChildren() {
			traverseAndCollectIdentifiers(child, outputList)
		}
	}
}

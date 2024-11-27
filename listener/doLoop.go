package listener

import (
	"bosch/parser"
	"bufio"
	"bytes"
	"github.com/antlr4-go/antlr/v4"
	"strings"
)

func (s *TreeShapeListener) EnterDoLoopStmt(ctx *parser.DoLoopStmtContext) {
	rules := parser.VisualBasic6ParserParserStaticData.RuleNames
	s.writer.WriteString("{\"RuleType\":\"DoLoopStatement\",")
	nodes := ctx.GetChildren()
	var loop bool
	for _, node := range nodes {
		_, ok := node.(antlr.TerminalNode)
		if ok {
			if strings.ToLower(node.(antlr.TerminalNode).GetText()) == "while" || strings.ToLower(node.(antlr.TerminalNode).GetText()) == "until" {
				s.writer.WriteString("\"Kind\":\"" + strings.ToLower(node.(antlr.TerminalNode).GetText()) + "\",")
			} else if strings.ToLower(node.(antlr.TerminalNode).GetText()) == "loop" {
				loop = true
			}
			continue
		} else if rules[node.(antlr.RuleContext).GetRuleIndex()] == "valueStmt" {
			var buffer bytes.Buffer
			if loop {
				s.writer.WriteString("\"BeforeLoop\":false,")
			} else {
				s.writer.WriteString("\"BeforeLoop\":true,")
			}
			writer := bufio.NewWriter(&buffer)
			handleLetExpression(node.(antlr.RuleContext).GetChildren(), writer)
			writer.Flush()
			s.writer.WriteString("\"Condition\": [" + strings.Trim(buffer.String(), ",") + "],")
		}

	}
	s.writer.WriteString("\"Body\": [")
}
func (s *TreeShapeListener) ExitDoLoopStmt(ctx *parser.DoLoopStmtContext) {
	s.exitContext()
	s.writer.WriteString("},") // Close the DoLoopStatement object
}

package listener

import (
	"bosch/parser"
	"bufio"
	"bytes"
	"strings"

	"github.com/antlr4-go/antlr/v4"
)

func handleFuncCalls(single antlr.Tree) string {
	rules := parser.VisualBasic6ParserParserStaticData.RuleNames
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	w.WriteString("{\"Type\": \"functioncall\",")
	for _, child := range single.(antlr.RuleNode).GetChildren() {
		switch val := child.(type) {
		case antlr.TerminalNode:
			continue
		case antlr.RuleNode:
			if rules[val.(antlr.RuleContext).GetRuleIndex()] == "argsCall" {
				for _, node := range val.GetChildren() {
					switch node := node.(type) {
					case antlr.TerminalNode:
						break
					case antlr.RuleNode:
						var some antlr.Tree
						var parent antlr.Tree
						proc := false
						some = node.GetChild(0)
					outer:
						for !proc {
							some = some.GetChild(0)
							switch some := some.(type) {
							case antlr.TerminalNode:
								parent = some.GetParent()
								break outer
							case antlr.RuleContext:
								if rules[some.GetRuleIndex()] == "iCS_S_ProcedureOrArrayCall" {
									proc = true
								}
							}
						}
						if proc {
							w.WriteString(handleFuncCalls(node.GetChild(0).GetChild(0).GetChild(0)) + ",")
						} else {
							w.WriteString("{\"type\":\"" + rules[parent.(antlr.RuleContext).GetRuleIndex()] + "\", \"sym\": \"" + strings.Trim(node.GetText(), "\"") + "\"},")
						}
					}
				}
			} else {
				w.WriteString("\"Identifier\": \"" + val.GetText() + "\", \"Arguments\": [")
			}
		}
	}
	w.Flush()
	str := buf.String()
	str = strings.TrimRight(str, ",") + "]}"
	return str
}

func (s *TreeShapeListener) EnterICS_B_ProcedureCall(ctx *parser.ICS_B_ProcedureCallContext) {
	s.writer.WriteString(handleFuncCalls(ctx) + ",")
}

func (s *TreeShapeListener) EnterECS_ProcedureCall(ctx *parser.ECS_ProcedureCallContext) {
	s.writer.WriteString(handleFuncCalls(ctx) + ",")
}

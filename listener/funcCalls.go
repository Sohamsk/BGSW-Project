package listener

import (
	"bosch/parser"
	"bufio"
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/antlr4-go/antlr/v4"
)

func handleFuncCalls(single antlr.Tree) string {
	rules := parser.VisualBasic6ParserParserStaticData.RuleNames
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	// w.WriteString("{\"Type\": \"functioncall\",")
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
						proc := false
						some = node.GetChild(0)
					outer:
						for !proc {
							some = some.GetChild(0)
							switch some := some.(type) {
							case antlr.TerminalNode:
								break outer
							case antlr.RuleContext:
								if rules[some.GetRuleIndex()] == "iCS_S_ProcedureOrArrayCall" {
									proc = true
								}
							}
						}
						if proc {
							w.WriteString("{\"Type\": \"FunctionCall\"," + handleFuncCalls(node.GetChild(0).GetChild(0).GetChild(0)) + "},")
						} else {
							var buf1 bytes.Buffer
							wr := bufio.NewWriter(&buf1)
							w.WriteString("{\"Type\":\"Expression\", \"body\":[")
							handleLetExpression(node.GetChildren(), wr)
							wr.Flush()
							w.WriteString(strings.Trim(string(buf1.Bytes()), ",") + "]},")
						}
					}
				}
			} else {
				sym := val.GetText()
				if strings.HasPrefix(sym, ".") {

					obj, exists := st.Peek()
					if exists {
						sym = obj + sym
					}
				}
				w.WriteString("\"Identifier\": \"" + sym + "\", \"Arguments\": [")
			}
		}
	}
	w.Flush()
	str := buf.String()
	str = strings.TrimRight(str, ",") + "]"
	return str
}

func handleMethodCalls(single antlr.Tree) string {
	rules := parser.VisualBasic6ParserParserStaticData.RuleNames
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	var id string
	args := false
	for _, child := range single.(antlr.RuleNode).GetChildren() {
		switch val := child.(type) {
		case antlr.TerminalNode:
			if val.GetText() == "." {
				id += "."
			}
			continue
		case antlr.RuleNode:
			if rules[val.(antlr.RuleContext).GetRuleIndex()] == "argsCall" {
				args = true
				w.WriteString("\"Identifier\": \"" + id + "\", \"Arguments\": [")

				for _, node := range val.GetChildren() {
					switch node := node.(type) {
					case antlr.TerminalNode:

						break
					case antlr.RuleNode:
						var some antlr.Tree
						proc := false
						some = node.GetChild(0)
					outer:
						for !proc {
							some = some.GetChild(0)
							switch some := some.(type) {
							case antlr.TerminalNode:
								break outer
							case antlr.RuleContext:
								if rules[some.GetRuleIndex()] == "iCS_S_ProcedureOrArrayCall" {
									proc = true
								}
							}
						}
						if proc {
							// figure out if this is func or method and use accordingly
							w.WriteString("{\"Type\": \"FunctionCall\"," + handleFuncCalls(node.GetChild(0).GetChild(0).GetChild(0)) + "},")
						} else {
							var buf1 bytes.Buffer
							wr := bufio.NewWriter(&buf1)
							w.WriteString("{\"Type\":\"Expression\", \"body\":[")
							handleLetExpression(node.GetChildren(), wr)
							wr.Flush()
							w.WriteString(strings.Trim(string(buf1.Bytes()), ",") + "]},")
						}
					}
				}
			} else {
				sym := val.GetText()
				id += sym
				// make it iterative such that till it gets the proper name without leading . it keeps trying
				var local []string
				for strings.HasPrefix(id, ".") {
					obj, exists := st.Peek()
					local = append(local, obj)
					st.Pop()
					if exists {
						id = obj + id
					} else {
						log.Println("Error: There may be a syntax error")
					}
				}
				for i := len(local) - 1; i >= 0; i-- {
					fmt.Println(i)
					st.Push(local[i])
				}
			}
		}
	}
	if !args {
		w.WriteString("\"Identifier\": \"" + strings.Trim(id, ".") + "\", \"Arguments\": [")
	}
	w.Flush()
	str := buf.String()
	str = strings.TrimRight(str, ",") + "]"
	return str
}

func (s *TreeShapeListener) EnterICS_B_ProcedureCall(ctx *parser.ICS_B_ProcedureCallContext) {
	s.writer.WriteString("{\"RuleType\": \"FunctionCall\",")
	s.writer.WriteString(handleFuncCalls(ctx) + "},")
}

func (s *TreeShapeListener) EnterECS_ProcedureCall(ctx *parser.ECS_ProcedureCallContext) {
	s.writer.WriteString("{\"RuleType\": \"FunctionCall\",")
	s.writer.WriteString(handleFuncCalls(ctx) + "},")
}

func (s *TreeShapeListener) EnterICS_B_MemberProcedureCall(ctx *parser.ICS_B_MemberProcedureCallContext) {
	s.writer.WriteString("{\"RuleType\": \"FunctionCall\",")
	s.writer.WriteString(handleMethodCalls(ctx) + "},")
}

package listener

import (
	"bosch/converter/models"
	"bosch/parser"
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"strings"

	"github.com/antlr4-go/antlr/v4"
)

// EnterEventStmt is called when production eventStmt is entered.
func (s *TreeShapeListener) EnterEventStmt(ctx *parser.EventStmtContext) {
	var event models.EventStatement
	event.RuleType = "EventStatement"

	children := ctx.GetChildren()

	for _, child := range children {
		switch child_t := child.(type) {
		case *parser.AmbiguousIdentifierContext:
			event.Identifier = child_t.GetText()
		case *parser.VisibilityContext:
			event.Visibility = child_t.GetText()
		case *parser.ArgListContext:
			for _, grandChild := range child_t.GetChildren() {
				_, ok := grandChild.(*parser.ArgContext)
				if ok {
					decl := models.DeclArg{}
					for _, greatGrandchild := range grandChild.GetChildren() {
						switch greatGrandchild.(type) {
						case antlr.TerminalNode:
							if greatGrandchild.(antlr.ParseTree).GetText() == "ByVal" {
								decl.IsPassedByRef = false
							}
						case *parser.AmbiguousIdentifierContext:
							decl.ArgumentName = greatGrandchild.(antlr.ParseTree).GetText()
						case *parser.AsTypeClauseContext:
							decl.ArgumentType = greatGrandchild.GetChild(2).(antlr.ParseTree).GetText()
						case *parser.TypeHintContext:
							argType, err := determineTypeFromHint(byte(greatGrandchild.(antlr.RuleContext).GetText()[0]))
							if err != nil {
								log.Println(err)
								argType = "Variant"
							}
							decl.ArgumentType = argType
						}
					}
					event.Arguments = append(event.Arguments, decl)
				}
			}
		}
	}

	obj, err := json.Marshal(event)
	if err != nil {
		log.Println("Error: ", err)
	}

	s.writer.WriteString(string(obj) + ",")
}

// EnterRaiseEventStmt is called when production raiseEventStmt is entered.
func (s *TreeShapeListener) EnterRaiseEventStmt(ctx *parser.RaiseEventStmtContext) {
	rules := parser.VisualBasic6ParserParserStaticData.RuleNames
	var w bytes.Buffer
	for _, child := range ctx.GetChildren() {
		child_rule, ok := child.(antlr.RuleNode)
		if ok {
			if rules[child_rule.(antlr.RuleContext).GetRuleIndex()] == "argsCall" {
				for _, node := range child_rule.GetChildren() {
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
				w.WriteString("\"Identifier\": \"" + child_rule.GetText() + "?.Invoke" + "\", \"Arguments\": [")
			}
		}
	}

	str := w.String()
	str = strings.TrimRight(str, ",") + "]"
	s.writer.WriteString("{\"RuleType\": \"FunctionCall\"," + str + "},")
}

// ExitRaiseEventStmt is called when production raiseEventStmt is exited.
func (s *TreeShapeListener) ExitRaiseEventStmt(ctx *parser.RaiseEventStmtContext) {}

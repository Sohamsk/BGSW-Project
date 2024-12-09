package listener

import (
	"bosch/parser"
	"bufio"
	"bytes"
	//	"fmt"
	"strings"

	"github.com/antlr4-go/antlr/v4"
)

// the arg call may also be a function so we may need to something there too
func handleLetExpression(nodes []antlr.Tree, w *bufio.Writer) {
	if len(nodes) == 0 {
		return
	}
	for _, node := range nodes {
		switch node := node.(type) {
		case antlr.TerminalNode:
			sym := node.GetText()
			if sym == " " {
				continue
			} else if sym == "(" || sym == ")" {
				w.WriteString("{\"Type\": \"Parenthesis\",")
			} else {
				w.WriteString("{\"Type\":\"Operator\",")
			}
			w.WriteString("\"Symbol\": \"" + sym + "\"},")
		case antlr.RuleNode:
			if node.GetChildCount() == 1 {
				single := node.GetChild(0).GetChild(0)
				_, holds := single.(antlr.TerminalNode)
				if holds {
					w.WriteString("{\"Type\": \"" + fetchParentOfTerminal(node) + "\",\"Symbol\": \"" + strings.Trim(single.(antlr.TerminalNode).GetText(), "\"") + "\"},")
				} else if parser.VisualBasic6ParserParserStaticData.RuleNames[single.(antlr.RuleContext).GetRuleIndex()] == "iCS_S_ProcedureOrArrayCall" {
					w.WriteString("{\"Type\": \"FunctionCall\",")
					w.WriteString(handleFuncCalls(single) + "},")
				} else {
					// find type of a node that is not a func call or and expression
					//					fmt.Println("fetch " + fetchParentOfTerminal(node.GetChild(0)))
					w.WriteString("{\"Type\":\"" + fetchParentOfTerminal(node) + "\", \"Symbol\":\"" + node.GetText() + "\"},")
				}
			} else {
				//				fmt.Println("nest")
				handleLetExpression(node.GetChildren(), w)
				//				fmt.Println("nested")
			}
		}
	}
}

// to do ternary operators, *Functions and procedures ,
func (s *TreeShapeListener) EnterLetStmt(ctx *parser.LetStmtContext) {
	// fmt.Println(parser.VisualBasic6ParserParserStaticData.RuleNames[ctx.GetRuleIndex()])
	nodes := ctx.GetChildren()
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	//fmt.Println(ctx.GetText())
	s.writer.WriteString("{\"RuleType\": \"expression\", \"Body\": ")
	writer.WriteString("[")
	handleLetExpression(nodes, writer)
	writer.Flush()
	str := buffer.String()
	str = strings.TrimRight(str, ",")
	s.writer.WriteString(str + "]},")
}

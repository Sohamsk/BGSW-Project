package listener

import (
	"bosch/parser"
	"bosch/stack"
	"bufio"
	"bytes"
	"reflect"
	"strings"

	"github.com/antlr4-go/antlr/v4"
)

type TreeShapeListener struct {
	*parser.BaseVisualBasic6ParserListener
	buf    *bytes.Buffer
	writer *bufio.Writer
	stack  stack.Stack
}

func NewTreeShapeListener(writer *bufio.Writer, buf *bytes.Buffer) *TreeShapeListener {
	l := new(TreeShapeListener)
	l.writer = writer
	l.buf = buf
	l.stack = *stack.InitStack()
	return l
}

//	func (s *TreeShapeListener) VisitTerminal(node antlr.TerminalNode) {
//		// 218 token type is variable name
//
// //	if node.GetSymbol().GetTokenType() == 218 {
// //		fmt.Println(node.GetText())
//
//	// figure something out to place the variable name and their types somewhere
//	// i will figure something out to make the it so that the varibale declaration statement state is stored
//	// i am thinking of using a stack to save the state to use later while constructing the tree
//
// //	}
//
//		//fmt.Println(node.GetSymbol().String())
//	}
func (s *TreeShapeListener) exitContext() {
	s.writer.Flush()
	string := s.buf.String()
	string = strings.Trim(string, ",") + "]"
	// fmt.Println(string)
	s.buf.Reset()
	s.writer.WriteString(string)
}

func (s *TreeShapeListener) EnterStartRule(ctx *parser.StartRuleContext) {
	s.writer.WriteString("\"body\": [")
}

func (s *TreeShapeListener) ExitStartRule(ctx *parser.StartRuleContext) {
	s.exitContext()
}

func fetchParentOfTerminal(someTree antlr.Tree) string {
	rules := parser.VisualBasic6ParserParserStaticData.RuleNames
	if reflect.TypeOf(someTree) == reflect.TypeOf(new(antlr.TerminalNode)) {
		return rules[someTree.GetParent().(antlr.RuleContext).GetRuleIndex()]
	} else {
		if someTree.GetChildCount() > 1 {
			fetchParentOfTerminal(someTree.GetChild(0))

		}
	}

	//	for _, node := range someTree.GetChildren() {
	//		switch node := node.(type) {
	//		case antlr.TerminalNode:
	//			break
	//		case antlr.RuleNode:
	//			var some antlr.Tree
	//			flag := false
	//			some = node.GetChild(0)
	//		outer:
	//			for !flag {
	//				some = some.GetChild(0)
	//				switch some := some.(type) {
	//				case antlr.TerminalNode:
	//					parent = some.GetParent()
	//					break outer
	//				case antlr.RuleContext:
	//					fmt.Println(rules[some.GetRuleIndex()])
	//
	//				}
	//			}
	//			if !flag {
	//			} else {
	//				return ""
	//			}
	//		}
	//	}
	return "kahitri gandlay "
}

func (s *TreeShapeListener) EnterICS_B_ProcedureCall(ctx *parser.ICS_B_ProcedureCallContext) {
	s.writer.WriteString(handleFuncCalls(ctx) + ",")
}

func (s *TreeShapeListener) EnterECS_ProcedureCall(ctx *parser.ECS_ProcedureCallContext) {
	s.writer.WriteString(handleFuncCalls(ctx) + ",")
}

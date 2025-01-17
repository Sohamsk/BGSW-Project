package listener

import (
	"bosch/parser"
	"bufio"
	"bytes"
	"strings"

	"github.com/antlr4-go/antlr/v4"
)

type TreeShapeListener struct {
	*parser.BaseVisualBasic6ParserListener
	buf    *bytes.Buffer
	writer *bufio.Writer
}

func NewTreeShapeListener(writer *bufio.Writer, buf *bytes.Buffer) *TreeShapeListener {
	l := new(TreeShapeListener)
	l.writer = writer
	l.buf = buf
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
	_, holds := someTree.(antlr.TerminalNode)
	if holds {
		return rules[someTree.GetParent().(antlr.RuleContext).GetRuleIndex()]
	} else {
		return fetchParentOfTerminal(someTree.GetChild(0))
	}
}

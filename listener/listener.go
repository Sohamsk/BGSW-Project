package listener

import (
	//	"bosch/converter/models"
	"bosch/parser"
	"bufio"
	"bytes"

	//	"encoding/json"
	//	"fmt"
	"strings"

	"github.com/antlr4-go/antlr/v4"
)

type TreeShapeListener struct {
	*parser.BaseVisualBasic6ParserListener
	buf    *bytes.Buffer
	writer *bufio.Writer
	SymTab map[string]string
}

func NewTreeShapeListener(writer *bufio.Writer, buf *bytes.Buffer) *TreeShapeListener {
	l := new(TreeShapeListener)
	l.writer = writer
	l.buf = buf
	l.SymTab = make(map[string]string)
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

var RuleMap = map[string]bool{ // map for handled rules
	"startRule":                     true,
	"ifThenElseStmt":                true,
	"letStmt":                       true,
	"forEachStmt":                   true,
	"forNextStmt":                   true,
	"blockStmt":                     true,
	"doLoopStmt":                    true,
	"selectCaseStmt":                true,
	"functionStmt":                  true,
	"withStmt":                      true,
	"variableStmt":                  true,
	"deftypeStmt":                   true,
	"printStmt":                     true,
	"comment":                       true,
	"subStmt":                       true,
	"exitStmt":                      true,
	"ifBlockStmt":                   true,
	"ifConditionStmt":               true,
	"ifElseIfBlockStmt":             true,
	"ifElseBlockStmt":               true,
	"iCS_S_VariableOrProcedureCall": true,
	"implicitCallStmt_InStmt":       true,
	"iCS_S_ProcedureOrArrayCall":    true,
	"argCall":                       true,
	"argsCall":                      true,
	"literal":                       true,
	"integerLiteral":                true,
	"ambiguousIdentifier":           true,
	"variableListStmt":              true,
	"variableSubStmt":               true,
	"asTypeClause":                  true,
	"arg":                           true,
	"argList":                       true,
	"type_":                         true,
	"baseType":                      true,
	"module":                        true,
	"moduleBody":                    true,
	"moduleBodyElement":             true,
	"moduleBlock":                   true,
	"block":                         true,
	"valueStmt":                     true,
	"octalLiteral":                  true,
	"doubleLiteral":                 true,
}

func (s *TreeShapeListener) EnterEveryRule(ctx antlr.ParserRuleContext) {
	// rules := parser.VisualBasic6ParserParserStaticData.RuleNames
	// //fmt.Println(rules[ctx.GetRuleIndex()])
	// _, ok := RuleMap[rules[ctx.GetRuleIndex()]]
	// //
	// if !ok { // This means rule is not handled so we'll just send it to MultiLineComment
	//
	//	UnhandledRule := models.MultiLineComment{}
	//	UnhandledRule.RuleType = "UnhandledRule"
	//	UnhandledRule.MultiLineComment = ctx.GetText()
	//	jsonData, err := json.Marshal(UnhandledRule)
	//	if err != nil {
	//		panic(err)
	//	}
	//	s.writer.WriteString(string(jsonData) + ",")
	//
	// }
}

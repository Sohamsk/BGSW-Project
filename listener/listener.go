package listener

import (
	"bosch/converter/models"
	"bosch/parser"
	"bufio"
	"bytes"
	"fmt"
	"log"

	"encoding/json"
	// "fmt"
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
	"ambiguousIdentifier":           true,
	"ambiguousKeyword":              true,
	"argCall":                       true,
	"argList":                       true,
	"argsCall":                      true,
	"arg":                           true,
	"asTypeClause":                  true,
	"baseType":                      true,
	"blockStmt":                     true,
	"block":                         true,
	"certainIdentifier":             true,
	"comment":                       true,
	"complexType":                   true,
	"deftypeStmt":                   true,
	"doLoopStmt":                    true,
	"doubleLiteral":                 true,
	"eCS_ProcedureCall":             true,
	"enumerationStmt":               true,
	"exitStmt":                      true,
	"forEachStmt":                   true,
	"forNextStmt":                   true,
	"functionStmt":                  true,
	"iCS_S_ProcedureOrArrayCall":    true,
	"iCS_S_VariableOrProcedureCall": true,
	"iCS_B_MemberOrProcedureCall":   true,
	"iCS_B_ProcedureCall":           true,
	"ifBlockStmt":                   true,
	"ifConditionStmt":               true,
	"ifElseBlockStmt":               true,
	"ifElseIfBlockStmt":             true,
	"ifThenElseStmt":                true,
	"implicitCallStmt_InStmt":       true,
	"implicitCallStmt_InBlock":      true,
	"integerLiteral":                true,
	"letStmt":                       true,
	"literal":                       true,
	"moduleBlock":                   true,
	"moduleBodyElement":             true,
	"moduleBody":                    true,
	"module":                        true,
	"octalLiteral":                  true,
	"printStmt":                     true,
	"propertyGetStmt":               true,
	"propertySetStmt":               true,
	"propertyLetStmt":               true,
	"selectCaseStmt":                true,
	"startRule":                     true,
	"subStmt":                       true,
	"setStmt":                       true,
	"type_":                         true,
	"typeHint":                      true,
	"typeStmt":                      true,
	"valueStmt":                     true,
	"variableListStmt":              true,
	"variableStmt":                  true,
	"variableSubStmt":               true,
	"visibility":                    true,
	"withStmt":                      true,
}

func (s *TreeShapeListener) EnterEveryRule(ctx antlr.ParserRuleContext) {
	rules := parser.VisualBasic6ParserParserStaticData.RuleNames
	//fmt.Println(rules[ctx.GetRuleIndex()])
	context_Type := rules[ctx.GetRuleIndex()]
	_, ok := RuleMap[context_Type]
	//
	if !ok { // This means rule is not handled so we'll just send it to MultiLineComment

		UnhandledRule := models.MultiLineComment{}
		UnhandledRule.RuleType = "UnhandledRule"
		if context_Type != "moduleBodyElement" || context_Type != "blockStmt" {
			UnhandledRule.MultiLineComment = ctx.GetText()
		}

		jsonData, err := json.Marshal(UnhandledRule)
		if err != nil {
			panic(err)
		}
		s.writer.WriteString(string(jsonData) + ",")
		log.Println("Warning: Not converting the rule:", context_Type)

	}
}

func determineTypeFromHint(hint byte) (string, error) {
	switch hint {
	case '%':
		return "Integer", nil
	case '&':
		return "Long", nil
	case '!':
		return "Single", nil
	case '#':
		return "Double", nil
	case '@':
		return "Currency", nil
	case '$':
		return "String", nil
	}
	return "", fmt.Errorf("Error: Unknown typehint")
}

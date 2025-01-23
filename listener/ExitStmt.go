package listener

import (
	"bosch/converter/models"
	"bosch/parser"
	"encoding/json"
	"log"

	"github.com/antlr4-go/antlr/v4"
)

func printReturn(returnName string) string {
	ir := models.ReturnStmt{}
	ir.RuleType = "ReturnStatement"
	ir.ReturnVariableName = returnName
	jsonData, err := json.Marshal(ir)
	if err != nil {
		panic(err)
	}
	return string(jsonData)
}

func (s *TreeShapeListener) EnterExit_Function(ctx *parser.Exit_FunctionContext) {
	rules := parser.VisualBasic6ParserParserStaticData.RuleNames
	var parent antlr.Tree
	parent = ctx
	for parent != nil && rules[parent.(antlr.RuleContext).GetRuleIndex()] != "functionStmt" {
		parent = parent.GetParent()
	}
	if parent == nil {
		panic("Syntax error")
	}
	varName := parent.GetChild(2).(*parser.AmbiguousIdentifierContext).GetText()

	jsonData := printReturn(varName)
	s.writer.WriteString(string(jsonData) + ",")
	log.Println("ReturnStatement converted to IR")
}

func printBreak(s *TreeShapeListener) {
	brk := models.BreakStmt{}
	brk.RuleType = "BreakStatement"
	jsonData, err := json.Marshal(brk)
	if err != nil {
		panic(err)
	}
	s.writer.WriteString(string(jsonData) + ",")
}

func (s *TreeShapeListener) EnterExit_Do(ctx *parser.Exit_DoContext) {
	printBreak(s)
}

func (s *TreeShapeListener) EnterExit_For(ctx *parser.Exit_ForContext) {
	printBreak(s)
}

func (s *TreeShapeListener) EnterExit_Sub(ctx *parser.Exit_SubContext) {
	jsonData := printReturn("")
	s.writer.WriteString(string(jsonData) + ",")
	log.Println("ReturnStatement converted to IR")
}

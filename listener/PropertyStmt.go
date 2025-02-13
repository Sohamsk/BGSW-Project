package listener

import (
	"bosch/converter/models"
	"bosch/parser"
	"encoding/json"
	"log"
	"reflect"

	"github.com/antlr4-go/antlr/v4"
)

func (s *TreeShapeListener) handleFuncLikeDecl(ctx antlr.ParserRuleContext, ruleType string, isProp bool, propType string) {
	nodes := ctx.GetChildren()
	prop := models.PropertyStatement{}
	if ruleType == "SubStatement" {
		prop.ReturnType = ""
	} else {
		prop.ReturnType = "Variant"
	}
	prop.RuleType = ruleType
	prop.Visibility = ""
	for _, child := range nodes {
		switch child.(type) {
		case *parser.VisibilityContext:
			prop.Visibility = child.(antlr.ParseTree).GetText()
		case *parser.AmbiguousIdentifierContext:
			prop.Identifier = child.(antlr.ParseTree).GetText()
		case *parser.AsTypeClauseContext:
			prop.ReturnType = child.GetChild(2).(antlr.RuleNode).GetText()
		case *parser.ArgListContext:
			for _, grandchild := range child.GetChildren() {
				if reflect.TypeOf(grandchild) == reflect.TypeOf(new(parser.ArgContext)) {
					decl := models.DeclArg{}
					for _, greatGrandchild := range grandchild.GetChildren() {
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
					prop.Arguments = append(prop.Arguments, decl)
				}
			}
		}
	}
	if isProp {
		s.SymTab["prop:"+propType+":"+prop.Identifier] = ""
	} else {
		s.SymTab["func:"+prop.Identifier] = prop.ReturnType
	}
	str, _ := json.Marshal(prop)
	s.writer.WriteString(string(str)[:len(str)-5] + "[")
}

func (s *TreeShapeListener) EnterPropertyGetStmt(ctx *parser.PropertyGetStmtContext) {
	s.handleFuncLikeDecl(ctx, "PropertyGetStatement", true, "get")
}

func (s *TreeShapeListener) ExitPropertyGetStmt(ctx *parser.PropertyGetStmtContext) {
	s.exitContext()
	s.writer.WriteString("},")
}

func (s *TreeShapeListener) EnterPropertySetStmt(ctx *parser.PropertySetStmtContext) {
	s.handleFuncLikeDecl(ctx, "PropertySetStatement", true, "set")
}

func (s *TreeShapeListener) ExitPropertySetStmt(ctx *parser.PropertySetStmtContext) {
	s.exitContext()
	s.writer.WriteString("},")
}

func (s *TreeShapeListener) EnterPropertyLetStmt(ctx *parser.PropertyLetStmtContext) {
	s.handleFuncLikeDecl(ctx, "PropertyLetStatement", true, "set")
}

func (s *TreeShapeListener) ExitPropertyLetStmt(ctx *parser.PropertyLetStmtContext) {
	s.exitContext()
	s.writer.WriteString("},")
}

package listener

import (
	"bosch/converter/models"
	"bosch/parser"
	"encoding/json"
	"fmt"
	"log"

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

// ExitEventStmt is called when production eventStmt is exited.
func (s *TreeShapeListener) ExitEventStmt(ctx *parser.EventStmtContext) {}

package listener

import (
	"bosch/parser"
	"fmt"
	"github.com/antlr4-go/antlr/v4"
	"reflect"
	"strings"
)

// EnterFunctionStmt is called when production functionStmt is entered.
func (s *TreeShapeListener) EnterFunctionStmt(ctx *parser.FunctionStmtContext) {
	nodes := ctx.GetChildren()
	returnType := "Variant"
	s.writer.WriteString("{\"RuleType\": \"FuncStatement\",")
	// handling arguments of a Sub
	index := 1
	Visibility := ""
	var arguments []string
	for _, child := range nodes {
		switch child.(type) {
		case *parser.VisibilityContext:
			Visibility = child.(antlr.ParseTree).GetText()
		case *parser.AmbiguousIdentifierContext:
			s.writer.WriteString("\"Identifier\": \"" + child.(antlr.ParseTree).GetText() + "\",")
			s.writer.WriteString("\"Visibility\": \"" + Visibility + "\",")
			s.writer.WriteString("\"Arguments\": [")
		case *parser.AsTypeClauseContext:
			returnType = child.GetChild(2).(antlr.RuleNode).GetText()
		case *parser.ArgListContext:
			passedByRef := true
			for _, grandchild := range child.GetChildren() {
				if reflect.TypeOf(grandchild) == reflect.TypeOf(new(parser.ArgContext)) {
					for _, greatGrandchild := range grandchild.GetChildren() {
						switch greatGrandchild.(type) {
						case antlr.TerminalNode:
							if greatGrandchild.(antlr.ParseTree).GetText() == "ByVal" {
								passedByRef = false
							}
						case *parser.AmbiguousIdentifierContext:
							arguments = append(arguments, fmt.Sprintf("{\"ArgumentName\":\"%s\",", greatGrandchild.(antlr.ParseTree).GetText()))
							index++
						case *parser.AsTypeClauseContext:
							argType := greatGrandchild.GetChild(2).(antlr.ParseTree).GetText()
							arguments[len(arguments)-1] += fmt.Sprintf("\"ArgumentType\": \"%s\",", argType)
							arguments[len(arguments)-1] += fmt.Sprintf("\"IsPassedByRef\": %t}", passedByRef)
						case *parser.TypeHintContext:
							arguments[len(arguments)-1] += fmt.Sprintf("\"ArgumentTypeHint\": \"%s\",", greatGrandchild.(antlr.ParseTree).GetText())
							arguments[len(arguments)-1] += fmt.Sprintf("\"IsPassedByRef\": %t}", passedByRef)
						}

					}
				}
			}
		}
		//		fmt.Printf("arguments is: %T\n", child)
	}
	if len(arguments) > 0 {
		s.writer.WriteString(strings.Join(arguments, ","))
	}
	s.writer.WriteString("],")
	s.writer.WriteString("\"ReturnType\":\"" + returnType + "\",")
	s.writer.WriteString("\"Body\": [")
}

// ExitFunctionStmt is called when production functionStmt is exited.
func (s *TreeShapeListener) ExitFunctionStmt(ctx *parser.FunctionStmtContext) {
	s.exitContext()
	s.writer.WriteString("},")
}

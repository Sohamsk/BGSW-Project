package listener

import (
	"bosch/parser"
	//"fmt"
	//	"github.com/antlr4-go/antlr/v4"
	//	"reflect"
	//	"strings"
)

//func (s *TreeShapeListener) EnterSubStmt(ctx *parser.SubStmtContext) {
//	nodes := ctx.GetChildren()
//	s.writer.WriteString("{\"RuleType\":\"SubStatement\",")
//	// handling arguments of a Sub
//	index := 1
//	Visibility := ""
//	var arguments []string
//	for _, child := range nodes {
//		switch child.(type) {
//		case *parser.VisibilityContext:
//			Visibility = child.(antlr.ParseTree).GetText()
//		case *parser.AmbiguousIdentifierContext:
//			s.writer.WriteString("\"Identifier\": \"" + child.(antlr.ParseTree).GetText() + "\",")
//			s.writer.WriteString("\"Visibility\": \"" + Visibility + "\",")
//			s.writer.WriteString("\"Arguments\": [")
//		case *parser.ArgListContext:
//			passedByRef := true
//			for _, grandchild := range child.GetChildren() {
//				//				fmt.Printf("arguments is: %T\n", grandchild)
//				if reflect.TypeOf(grandchild) == reflect.TypeOf(new(parser.ArgContext)) {
//					for _, greatGrandchild := range grandchild.GetChildren() {
//						switch greatGrandchild.(type) {
//						case antlr.TerminalNode:
//							if greatGrandchild.(antlr.ParseTree).GetText() == "ByVal" {
//								passedByRef = false
//							}
//
//						case *parser.AmbiguousIdentifierContext:
//							arguments = append(arguments, fmt.Sprintf("{\"ArgumentName\":\"%s\",", greatGrandchild.(antlr.ParseTree).GetText()))
//							index++
//						case *parser.AsTypeClauseContext:
//							argType := greatGrandchild.GetChild(2).(antlr.ParseTree).GetText()
//							arguments[len(arguments)-1] += fmt.Sprintf("\"ArgumentType\": \"%s\",", argType)
//							arguments[len(arguments)-1] += fmt.Sprintf("\"IsPassedByRef\": %t}", passedByRef)
//						case *parser.TypeHintContext:
//							arguments[len(arguments)-1] += fmt.Sprintf("\"ArgumentTypeHint\": \"%s\",", greatGrandchild.(antlr.ParseTree).GetText())
//							arguments[len(arguments)-1] += fmt.Sprintf("\"IsPassedByRef\": \"%t\"}", passedByRef)
//						}
//
//					}
//				}
//			}
//		}
//		//		fmt.Printf("arguments is: %T\n", child)
//	}
//	if len(arguments) > 0 {
//		s.writer.WriteString(strings.Join(arguments, ","))
//	}
//	s.writer.WriteString("],")
//	s.writer.WriteString("\"SubBody\": [")
//}

func (s *TreeShapeListener) EnterSubStmt(ctx *parser.SubStmtContext) {
	s.handleFuncLikeDecl(ctx, "SubStatement", false, "")
}

func (s *TreeShapeListener) ExitSubStmt(ctx *parser.SubStmtContext) {
	s.exitContext()
	s.writer.WriteString("},")
}

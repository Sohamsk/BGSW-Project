package listener

import (
	"bosch/parser"
	"bosch/stack"
	"bufio"
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"github.com/antlr4-go/antlr/v4"
)

type TreeShapeListener struct {
	*parser.BaseVisualBasic6ParserListener
	writer *bufio.Writer
	stack  stack.Stack
}

func NewTreeShapeListener(writer *bufio.Writer) *TreeShapeListener {
	l := new(TreeShapeListener)
	l.writer = writer
	l.stack = *stack.InitStack()
	return l
}

//func (s *TreeShapeListener) VisitTerminal(node antlr.TerminalNode) {
//	// 218 token type is variable name
////	if node.GetSymbol().GetTokenType() == 218 {
////		fmt.Println(node.GetText())
//		// figure something out to place the variable name and their types somewhere
//		// i will figure something out to make the it so that the varibale declaration statement state is stored
//		// i am thinking of using a stack to save the state to use later while constructing the tree
////	}
//	//fmt.Println(node.GetSymbol().String())
//}

func (s *TreeShapeListener) EnterStartRule(ctx *parser.StartRuleContext) {
	s.writer.WriteString("\"body\": [")
}

func (s *TreeShapeListener) ExitStartRule(ctx *parser.StartRuleContext) {
	s.writer.Flush()
	s.writer.WriteString("]")
}

func (s *TreeShapeListener) EnterVariableSubStmt(ctx *parser.VariableSubStmtContext) {
	nodes := ctx.GetChildren()
	s.writer.WriteString("{\"RuleType\": \"DeclareVariable\",")
	s.writer.WriteString("\"Identifier\": \"" + nodes[0].(antlr.ParseTree).GetText() + "\",")
	if len(nodes) == 3 {
		s.writer.WriteString("\"Type\": \"" + nodes[2].GetChild(2).(antlr.RuleNode).GetText() + "\"")
	} else {
		s.writer.WriteString("\"Type\": \"VARIANT\"")
	}
	s.writer.WriteString("},")
}

func fetchParentOfTerminal(someTree antlr.RuleNode) string {
	rules := parser.VisualBasic6ParserParserStaticData.RuleNames
	if reflect.TypeOf(someTree) == reflect.TypeOf(new(antlr.TerminalNode)) {
		return rules[someTree.GetParent().(antlr.RuleContext).GetRuleIndex()]
	} else {
		if someTree.GetChildCount() > 1 {
			//			fetchParentOfTerminal(someTree.GetChild(0))
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

// the arg call may also be a function so we may need to something there too
func handleFuncCalls(single antlr.Tree) string {
	rules := parser.VisualBasic6ParserParserStaticData.RuleNames
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	w.WriteString("{\"Type\": \"functioncall\",")
	for _, child := range single.(antlr.RuleNode).GetChildren() {
		switch val := child.(type) {
		case antlr.TerminalNode:
			continue
		case antlr.RuleNode:
			if rules[val.(antlr.RuleContext).GetRuleIndex()] == "argsCall" {
				for _, node := range val.GetChildren() {
					switch node := node.(type) {
					case antlr.TerminalNode:
						break
					case antlr.RuleNode:
						var some antlr.Tree
						var parent antlr.Tree
						proc := false
						some = node.GetChild(0)
					outer:
						for !proc {
							some = some.GetChild(0)
							switch some := some.(type) {
							case antlr.TerminalNode:
								parent = some.GetParent()
								break outer
							case antlr.RuleContext:
								if rules[some.GetRuleIndex()] == "iCS_S_ProcedureOrArrayCall" {
									proc = true
								}
							}
						}
						if proc {
							w.WriteString(handleFuncCalls(node.GetChild(0).GetChild(0).GetChild(0)) + ",")
						} else {
							w.WriteString("{\"type\":\"" + rules[parent.(antlr.RuleContext).GetRuleIndex()] + "\", \"sym\": \"" + strings.Trim(node.GetText(), "\"") + "\"},")
						}
					}
				}
			} else {
				w.WriteString("\"Identifier\": \"" + val.GetText() + "\", \"Arguments\": [")
			}
		}
	}
	w.Flush()
	str := buf.String()
	str = strings.TrimRight(str, ",") + "]}"
	return str
}

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
				//	fmt.Println(parser.VisualBasic6ParserParserStaticData.RuleNames[single.(antlr.RuleContext).GetRuleIndex()])
				if parser.VisualBasic6ParserParserStaticData.RuleNames[single.(antlr.RuleContext).GetRuleIndex()] == "iCS_S_ProcedureOrArrayCall" {
					w.WriteString(handleFuncCalls(single) + ",")
				} else {
					// find type of a node that is not a func call or and expression
					fmt.Println(fetchParentOfTerminal(node))
					w.WriteString("{\"Identifier\": \"" + node.GetText() + "\"},")
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

	s.writer.WriteString("{\"RuleType\": \"expression\", \"Body\": ")
	writer.WriteString("[")
	handleLetExpression(nodes, writer)
	writer.Flush()
	str := buffer.String()
	str = strings.TrimRight(str, ",")
	s.writer.WriteString(str + "]},")
}

func (s *TreeShapeListener) EnterSubStmt(ctx *parser.SubStmtContext) {
	nodes := ctx.GetChildren()
	s.writer.WriteString("{\"SubStatement\": {")
	// handling arguments of a Sub
	index := 1
	Visibility := "Public"
	var arguments []string
	for _, child := range nodes {
		switch child.(type) {
		case *parser.VisibilityContext:
			Visibility = child.(antlr.ParseTree).GetText()
		case *parser.AmbiguousIdentifierContext:
			s.writer.WriteString("\"SubName\": \"" + child.(antlr.ParseTree).GetText() + "\",")
			s.writer.WriteString("\"Visibility\": \"" + Visibility + "\",")
			s.writer.WriteString("\"arguments\": [")
		case *parser.ArgListContext:
			passedByRef := true
			for _, grandchild := range child.GetChildren() {
				fmt.Printf("arguments is: %T\n", grandchild)
				if reflect.TypeOf(grandchild) == reflect.TypeOf(new(parser.ArgContext)) {
					for _, greatGrandchild := range grandchild.GetChildren() {
						switch greatGrandchild.(type) {
						case antlr.TerminalNode:
							if greatGrandchild.(antlr.ParseTree).GetText() == "ByVal" {
								passedByRef = false
							}

						case *parser.AmbiguousIdentifierContext:
							arguments = append(arguments, fmt.Sprintf("{\"ArgumentName%d\":\"%s\",", index, greatGrandchild.(antlr.ParseTree).GetText()))
							index++
						case *parser.AsTypeClauseContext:
							argType := greatGrandchild.GetChild(2).(antlr.ParseTree).GetText()
							arguments[len(arguments)-1] += fmt.Sprintf("\"ArgumentType\": \"%s\",", argType)
							arguments[len(arguments)-1] += fmt.Sprintf("\"IsPassedByRef\": \"%t\"}", passedByRef)
						case *parser.TypeHintContext:
							arguments[len(arguments)-1] += fmt.Sprintf("\"ArgumentTypeHint\": \"%s\",", greatGrandchild.(antlr.ParseTree).GetText())
							arguments[len(arguments)-1] += fmt.Sprintf("\"IsPassedByRef\": \"%t\"}", passedByRef)
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
	s.writer.WriteString("\"SubBody\": [")
	//	block := nodes[5].GetChildren() // discuss the array accessing
	//	handleSubBody(block)
}
func (s *TreeShapeListener) ExitSubStmt(ctx *parser.SubStmtContext) {
	s.writer.WriteString("]}} ")
}

func (s *TreeShapeListener) EnterICS_B_ProcedureCall(ctx *parser.ICS_B_ProcedureCallContext) {
	s.writer.WriteString(handleFuncCalls(ctx) + ",")
}

func (s *TreeShapeListener) EnterECS_ProcedureCall(ctx *parser.ECS_ProcedureCallContext) {
	s.writer.WriteString(handleFuncCalls(ctx) + ",")
}

//func handleSubBody(blockTree []antlr.Tree) {
//	nodes := blockTree
//	for _, child := range nodes {
//		if reflect.TypeOf(child) == reflect.TypeOf(new(parser.BlockStmtContext)) {
//			//			fmt.Println(child.(antlr.ParseTree).GetText())
//		}
//		//		fmt.Printf("arguments is: %T\n", child)
//
//	}
//
//}

func (s *TreeShapeListener) EnterDoLoopStmt(ctx *parser.DoLoopStmtContext) {
	fmt.Println("Enter do statement")

}

func (s *TreeShapeListener) ExitDoLoopStmt(ctx *parser.DoLoopStmtContext) {
	fmt.Println("Exit do statement")
}

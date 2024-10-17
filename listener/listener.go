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
    buf *bytes.Buffer
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
}
func (s *TreeShapeListener) ExitSubStmt(ctx *parser.SubStmtContext) {
    s.exitContext()
	s.writer.WriteString("}} ")
}

func (s *TreeShapeListener) EnterICS_B_ProcedureCall(ctx *parser.ICS_B_ProcedureCallContext) {
	s.writer.WriteString(handleFuncCalls(ctx) + ",")
}

func (s *TreeShapeListener) EnterECS_ProcedureCall(ctx *parser.ECS_ProcedureCallContext) {
	s.writer.WriteString(handleFuncCalls(ctx) + ",")
}

func (s *TreeShapeListener) EnterDoLoopStmt(ctx *parser.DoLoopStmtContext) {
	rules := parser.VisualBasic6ParserParserStaticData.RuleNames
    s.writer.WriteString("{\"RuleType\":\"DoLoopStatement\",")
    nodes := ctx.GetChildren()
    for _, node := range nodes {
        _, ok := node.(antlr.TerminalNode)
        if (ok) {
            if (strings.ToLower(node.(antlr.TerminalNode).GetText()) == "while" || strings.ToLower(node.(antlr.TerminalNode).GetText()) == "until") {
                s.writer.WriteString("\"Kind\":\"" + strings.ToLower(node.(antlr.TerminalNode).GetText()) + "\",")
            }
            continue
        } else if (!ok) {
            if (rules[node.(antlr.RuleContext).GetRuleIndex()] == "valueStmt") {
                var buffer bytes.Buffer
                writer := bufio.NewWriter(&buffer)
                handleLetExpression(node.(antlr.RuleContext).GetChildren(), writer)
                writer.Flush()
                s.writer.WriteString("\"Condition\": [" + strings.Trim(buffer.String(), ",") + "],")
            }
        }

    }
    s.writer.WriteString("\"Body\": [")
}
func (s *TreeShapeListener) ExitDoLoopStmt(ctx *parser.DoLoopStmtContext) {
    s.exitContext()
    s.writer.WriteString("},") // Close the DoLoopStatement object
}

func (s *TreeShapeListener) EnterPrintStmt(ctx *parser.PrintStmtContext) {
	nodes := ctx.GetChildren()
	//    s.writer.WriteString("{\"CommentedUI\":\"" + ctx.GetText() + "\"},")
	for _, node := range nodes {
		switch node := node.(type) {
		case antlr.RuleNode:
			text := node.GetText()
			s.writer.WriteString("{\"RuleType\":\"PrintStmt\", \"Data\": \"" + text + "\"},")
		}
	}
}

func (s *TreeShapeListener) EnterForNextStmt(ctx *parser.ForNextStmtContext) {
	nodes := ctx.GetChildren()
	var initialization, condition, from string

	// Start the ForLoopStatement object
	s.writer.WriteString("{\"RuleType\":\"ForLoopStatement\",\"Body\":[")

	for i, node := range nodes {
		switch n := node.(type) {
		case antlr.TerminalNode:
			// Capture initialization
			if n.GetText() == "For" && i+1 < len(nodes) {
				initialization = nodes[i+2].(antlr.ParseTree).GetText()
			} else if n.GetText() == "To" && i+1 < len(nodes) {
				condition = nodes[i+2].(antlr.ParseTree).GetText()
			} else if n.GetText() == "=" && i+1 < len(nodes) {
				from = nodes[i+2].(antlr.ParseTree).GetText()
			}
		}
	}
	s.writer.WriteString("{\"Initialization\": \"" + initialization + "\",")
	s.writer.WriteString("\"Start\": \"" + from + "\",")
	s.writer.WriteString("\"End\": \"" + condition + "\"},")
}

func (s *TreeShapeListener) ExitForNextStmt(ctx *parser.ForNextStmtContext) {
	s.writer.WriteString("]}")
	s.writer.WriteString("}")
}
func (s *TreeShapeListener) EnterDeftypeStmt(ctx *parser.DeftypeStmtContext) {
    nodes := ctx.GetChildren()
    s.writer.WriteString("{\"RuleType\":\"DefType\",")
    for _, node := range(nodes) {
        if (reflect.TypeOf(node) == reflect.TypeOf(new (parser.LetterrangeContext))) {
            s.writer.WriteString("\"LetterRange\":\"")
            first := true
            for _, child := range(node.GetChildren()) {
                if (reflect.TypeOf(child) == reflect.TypeOf(new(parser.CertainIdentifierContext))) {
                    if (first) {
                        s.writer.WriteString(child.(antlr.RuleNode).GetText())
                        first = !first
                    } else {
                        s.writer.WriteString("-" + child.(antlr.RuleNode).GetText())
                    }
                }
            }
            s.writer.WriteString("\"")
        } else if (node.(antlr.TerminalNode).GetText() != " ") {
            dataType := strings.ToLower(node.(antlr.TerminalNode).GetText())

            switch dataType {
            case "defbool":
                s.writer.WriteString("\"DataType\":\"Boolean\",")
            case "defbyte":
                s.writer.WriteString("\"DataType\":\"Byte\",")
            case "defint":
                s.writer.WriteString("\"DataType\":\"Integer\",")
            case "deflng":
                s.writer.WriteString("\"DataType\":\"Long\",")
            case "deflnglng":
                s.writer.WriteString("\"DataType\":\"LongLong (valid on 64-bit platforms only)\",")
            case "deflngptr":
                s.writer.WriteString("\"DataType\":\"LongPtr\",")
            case "defcur":
                s.writer.WriteString("\"DataType\":\"Currency\",")
            case "defsng":
                s.writer.WriteString("\"DataType\":\"Single\",")
            case "defdbl":
                s.writer.WriteString("\"DataType\":\"Double\",")
            case "defDec":
                s.writer.WriteString("\"DataType\":\"Decimal (not currently supported)\",")
            case "defdate":
                s.writer.WriteString("\"DataType\":\"Date\",")
            case "defstr":
                s.writer.WriteString("\"DataType\":\"String\",")
            case "defobj":
                s.writer.WriteString("\"DataType\":\"Object\",")
            case "defvar":
                s.writer.WriteString("\"DataType\":\"Variant\",")
            }
        } else {
            continue
        }
    }
    s.writer.WriteString("},")
}

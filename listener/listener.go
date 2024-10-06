package listener

import (
	"bosch/parser"
	"bosch/stack"
	"bufio"
	"bytes"
	"fmt"
	"reflect"
	"strconv"
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

func fetchParentOfTerminal(someRuleNode antlr.RuleNode, someRuleNodeName string) string {
	rules := parser.VisualBasic6ParserParserStaticData.RuleNames

	for _, node := range someRuleNode.GetChildren() {
		switch node := node.(type) {
		case antlr.TerminalNode:
			break
		case antlr.RuleNode:
			var some antlr.Tree
			var parent antlr.Tree
			flag := false
			some = node.GetChild(0)
		outer:
			for !flag {
				some = some.GetChild(0)
				switch some := some.(type) {
				case antlr.TerminalNode:
					parent = some.GetParent()
					break outer
				case antlr.RuleContext:
					fmt.Println(rules[some.GetRuleIndex()])
					if rules[some.GetRuleIndex()] == someRuleNodeName {
						flag = true
					}
				}
			}
			if !flag {
				return rules[parent.(antlr.RuleContext).GetRuleIndex()]
			} else {
				return ""
			}
		}
	}
	return ""
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
                //fmt.Println("YO TYPE IS :" + fetchParentOfTerminal(val, "iCS_S_ProcedureOrArrayCall"))
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
                                fmt.Println(rules[some.GetRuleIndex()])
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

func handleLetExpression(nodes []antlr.Tree, w *bufio.Writer, first bool) {
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
					w.WriteString("{\"Identifier\": \"" + node.GetText() + "\"},")
				}
			} else {
				//				fmt.Println("nest")
				handleLetExpression(node.GetChildren(), w, false)
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
	handleLetExpression(nodes, writer, true)
	writer.Flush()
	str := buffer.String()
	str = strings.TrimRight(str, ",")
	s.writer.WriteString(str + "]},")
}

func (s *TreeShapeListener) EnterSubStmt(ctx *parser.SubStmtContext) {
	nodes := ctx.GetChildren()
	s.writer.WriteString("{\"SubStatement\": {")
	s.writer.WriteString("\"SubName\": \"" + nodes[2].(antlr.ParseTree).GetText() + "\",")
	s.writer.WriteString("\"arguments\": [")
	// handling arguments of a Sub
	index := 1
	for _, child := range nodes[3].GetChildren() {
		passedByRef := true
		if reflect.TypeOf(child) == reflect.TypeOf(new(parser.ArgContext)) {
			for _, grandchild := range child.GetChildren() {
				switch grandchild.(type) {
				case antlr.TerminalNode:
					if grandchild.(antlr.ParseTree).GetText() == "ByVal" {
						passedByRef = false
					}

				case *parser.AmbiguousIdentifierContext:
					s.writer.WriteString("{")
					s.writer.WriteString("\"ArgumentName" + strconv.Itoa(index) + "\": \"" + grandchild.(antlr.ParseTree).GetText() + "\",")
					index += 1
				case *parser.AsTypeClauseContext:
					s.writer.WriteString("\"ArgumentType\": \"" + grandchild.GetChild(2).(antlr.ParseTree).GetText() + "\",")
					s.writer.WriteString("\"IsPassedByRef\": \"" + strconv.FormatBool(passedByRef) + "\"")
				}
				//				fmt.Printf("arguments is: %T\n", grandchild)
			}
			s.writer.WriteString("},") // Figure out a way to avoid the trailing comma
		}
		//		fmt.Printf("arguments is: %T\n", child)
	}
	s.writer.WriteString("],")
	s.writer.WriteString("\"SubBody\": [")
	//	block := nodes[5].GetChildren() // discuss the array accessing
	//	handleSubBody(block)
}
func (s *TreeShapeListener) ExitSubStmt(ctx *parser.SubStmtContext) {
	s.writer.WriteString("]}} ")
}

//func (s *TreeShapeListener)  f

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

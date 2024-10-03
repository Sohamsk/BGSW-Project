package listener

import (
	"bosch/parser"
	"bosch/stack"
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"reflect"
	"strconv"

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

func handleLetExpression(nodes []antlr.Tree, w *bufio.Writer, first bool) {
	if (len(nodes) == 0) {
		return
	}
	if (!first && len(nodes) == 1) {
		fmt.Print(",")
	}
	for _, node := range(nodes) {
		switch node.(type) {
		case antlr.TerminalNode:
			sym := node.(antlr.TerminalNode).GetText()
			if (sym == " ") {
				continue
			} else if (sym == "(" || sym == ")") {
				w.WriteString("{\"Type\": \"Parenthesis\",")
			} else {
				w.WriteString("{\"Type\":\"Operator\",")
			}
			w.WriteString("\"Symbol\": \"" + sym + "\"},")
//			fmt.Println(sym)
		case antlr.RuleNode:
			if (node.(antlr.RuleNode).GetChildCount() == 1) {
				w.WriteString("{\"Identifier\": \"" + node.(antlr.RuleNode).GetText() + "\"},")
				if (parser.VisualBasic6ParserParserStaticData.RuleNames[node.GetChild(0).GetChild(0).(antlr.RuleContext).GetRuleIndex()] == "iCS_S_ProcedureOrArrayCall") {
					fmt.Println("this is either an array or a function call")
				}
				fmt.Println(node.(antlr.RuleNode).GetText())
			} else {
//				fmt.Println("nest")
				handleLetExpression(node.(antlr.RuleNode).GetChildren(), w, false)
//				fmt.Println("nested")
			}
		}
	}
}

// to do ternary operators, *Functions and procedures ,
func (s *TreeShapeListener) EnterLetStmt(ctx *parser.LetStmtContext) {
	fmt.Println(parser.VisualBasic6ParserParserStaticData.RuleNames[ctx.GetRuleIndex()])
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
		if reflect.TypeOf(child) == reflect.TypeOf(new(parser.ArgContext)) {
			for _, grandchild := range child.GetChildren() {
				switch grandchild.(type) {
				case *parser.AmbiguousIdentifierContext:
					s.writer.WriteString("{")
					s.writer.WriteString("\"ArgumentName" + strconv.Itoa(index) + "\": \"" + grandchild.(antlr.ParseTree).GetText() + "\",")
					index += 1
				case *parser.AsTypeClauseContext:
					s.writer.WriteString("\"ArgumentType\": \"" + grandchild.GetChild(2).(antlr.ParseTree).GetText() + "\"")
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

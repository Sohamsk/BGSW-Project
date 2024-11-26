package listener

import (
	"bosch/parser"
	"bufio"
	"bytes"
	//	"fmt"
	"strings"

	"github.com/antlr4-go/antlr/v4"
)

func (s *TreeShapeListener) EnterSelectCaseStmt(ctx *parser.SelectCaseStmtContext) {
	nodes := ctx.GetChildren()
	s.writer.WriteString("{\"RuleType\":\"SelectCaseStmt\",")
	for _, node := range nodes {
		switch node := node.(type) {
		case *parser.VsICSContext: //it is the identifier the selectCaseStmt works on
			s.writer.WriteString("\"SelectCaseIdentifier\": \"" + node.GetText() + "\",")
			s.writer.WriteString("\"Cases\": [")
		}
	}
}

func (s *TreeShapeListener) ExitSelectCaseStmt(ctx *parser.SelectCaseStmtContext) {
	s.exitContext()
	s.writer.WriteString("},")
}

func (s *TreeShapeListener) EnterSC_Case(ctx *parser.SC_CaseContext) {
	nodes := ctx.GetChildren()
	for _, node := range nodes {
		_, fits := node.(*parser.CaseCondExprContext) // sC_Cond == CaseCondExprContext
		//		fmt.Printf("%T", node, "\n")
		if fits {
			var buffer bytes.Buffer
			writer := bufio.NewWriter(&buffer)
			//handleLetExpression(node.GetChildren(), writer)

			childrenOfsC_Cond := node.GetChildren()
			for _, child := range childrenOfsC_Cond {

				//				fmt.Printf("%T", child, "\n")
				_, isTerminal := child.(antlr.TerminalNode)
				_, isCondExprIS := child.(*parser.CaseCondExprIsContext) // will have to implement separate logic for this
				if !isTerminal && !isCondExprIS {
					handleLetExpression(child.GetChildren(), writer)
				}
				// This is marked as EDGE CASE FOR NOW
				// TODO: SOLVE THIS EDGE CASE USING CUSTOM LOGIC
				//	if isCondExprIS {
				//		handleLetExpression(child.GetChildren(), writer)
				//	}
			}
			writer.Flush()
			s.writer.WriteString("{\"Case\": {\"Condition\":[" + strings.Trim(buffer.String(), ",") + "],")
		}
		_, isDefault := node.(*parser.CaseCondElseContext)
		if isDefault {
			var buffer bytes.Buffer
			writer := bufio.NewWriter(&buffer)

			handleLetExpression(node.GetChildren(), writer)
			writer.Flush()
			s.writer.WriteString("{\"DefaultCase\": {")
		}
	}
	s.writer.WriteString("\"CaseBody\": [")

}

func (s *TreeShapeListener) ExitSC_Case(ctx *parser.SC_CaseContext) {
	s.exitContext()
	s.writer.WriteString("}},")
}

//func (s *TreeShapeListener) EnterCaseCondExpr(ctx *parser.CaseCondExprContext) {
//	nodes := ctx.GetChildren()
//	for _, node := range nodes {
//		//fmt.Printf("%T", node, "\n")
//		_, isTerminal := node.(antlr.TerminalNode)
//		if isTerminal {
//			//	fmt.Println(node.(antlr.TerminalNode).GetText())
//			if node.(antlr.TerminalNode).GetText() == "," {
//			}
//		}
//
//	}
//}

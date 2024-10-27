package listener

import (
	"bosch/parser"
	"github.com/antlr4-go/antlr/v4"
	"reflect"
	"strings"
)

func (s *TreeShapeListener) EnterDeftypeStmt(ctx *parser.DeftypeStmtContext) {
	nodes := ctx.GetChildren()
	s.writer.WriteString("{\"RuleType\":\"DefType\",")
	for _, node := range nodes {
		if reflect.TypeOf(node) == reflect.TypeOf(new(parser.LetterrangeContext)) {
			s.writer.WriteString("\"LetterRange\":\"")
			first := true
			for _, child := range node.GetChildren() {
				if reflect.TypeOf(child) == reflect.TypeOf(new(parser.CertainIdentifierContext)) {
					if first {
						s.writer.WriteString(child.(antlr.RuleNode).GetText())
						first = !first
					} else {
						s.writer.WriteString("-" + child.(antlr.RuleNode).GetText())
					}
				}
			}
			s.writer.WriteString("\"")
		} else if node.(antlr.TerminalNode).GetText() != " " {
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

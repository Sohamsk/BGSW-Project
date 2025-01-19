package listener

import (
	"bosch/parser"
	"bosch/stack"
	"fmt"
	"strings"
)

type WithStmtListener struct {
	*parser.BaseVisualBasic6ParserListener
	stack stack.Stack
	buf   *strings.Builder
}

func NewWithStmtListener() *WithStmtListener {
	return &WithStmtListener{
		stack: *stack.InitStack(),
		buf:   &strings.Builder{},
	}
}

func (w *WithStmtListener) EnterWithStmt(ctx *parser.WithStmtContext) {
	object := ctx.GetChild(1).GetChildren()
	w.stack.Push(object)
	w.buf.WriteString(fmt.Sprintf("\"with\": {\"object\": \"%s\", \"operations\": [", object))
}

func (w *WithStmtListener) ExitWithStmt(ctx *parser.WithStmtContext) {

	w.buf.WriteString("]}")

	w.stack.Pop()
	fmt.Println(w.buf.String())
	w.buf.Reset()
}

func (w *WithStmtListener) EnterImplicitCallStmt_InStmt(ctx *parser.ImplicitCallStmt_InStmtContext) {
	currentObject, ok := w.stack.Peek().(string)
	if !ok {
		fmt.Println("Error: Unable to fetch the current object from the stack.")
		return
	}

	methodOrProperty := ctx.GetChild(0).GetChildren()
	w.buf.WriteString(fmt.Sprintf("{\"object\": \"%s\", \"action\": \"%s\"},", currentObject, methodOrProperty))
}

func (w *WithStmtListener) ExitImplicitCallStmt_InStmt(ctx *parser.ImplicitCallStmt_InStmtContext) {

}

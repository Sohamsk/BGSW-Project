package models

import "encoding/json"

type FileContext struct {
	FileName string
	FileType string
	Body     []json.RawMessage `json:"body"`
}

type Rule struct {
	RuleType string
}

type Dim struct {
	Identifier string `json:"identifier"`
	Type       string `json:"type"`
	Scope      string `json:"scope,omitempty"`
	WithEvents bool   `json:"withEvents,omitempty"`
}
type ArgType struct {
	Type string
}

type Literal struct {
	ArgType
	Symbol string
}

type FuncArg struct {
	ArgType
	Identifier string
	Arguments  []json.RawMessage
}

type FuncRule struct {
	Rule
	Identifier string
	Arguments  []json.RawMessage
}

type ExpressionRule struct {
	Rule
	Body []json.RawMessage
}

type DeclArg struct {
	ArgumentName  string
	ArgumentType  string
	IsPassedByRef bool
}

type SubStmt struct {
	Rule
	Identifier string
	Visibility string
	Arguments  []DeclArg
	SubBody    []json.RawMessage
}

type DoloopStmt struct {
	Rule
	Kind       string
	BeforeLoop bool
	Condition  []json.RawMessage
	Body       []json.RawMessage
}

type FuncDecl struct {
	Rule
	Identifier string
	Visibility string
	Arguments  []DeclArg
	ReturnType string
	Body       []json.RawMessage
}
type ElseIfRule struct {
	Rule
	Condition   []json.RawMessage
	ElseIfBlock []json.RawMessage
}

type ElseRule struct {
	Rule
	Body []json.RawMessage
}

type IfThenElseStmtRule struct {
	Rule
	IsBlock   bool
	Condition []json.RawMessage
	IfBlock   []json.RawMessage
}
type ForNext struct {
	Rule
	IdentifierName string            `json:"IdentifierName"`
	Start          string            `json:"Start"`
	End            string            `json:"End"`
	Step           string            `json:"Step"`
	Body           []json.RawMessage `json:"Body"`
}

type ReturnStmt struct {
	Rule
	ReturnVariableName string
}

type BreakStmt struct {
	Rule
}

type Comment struct {
	Rule
	CommentText string `json:"CommentText"`
}

type WithStmt struct {
	Rule
	Object string
	Body   []ExpressionRule
}

type MultiLineComment struct { // This struct is for adding a MultiLineComment in C# code for part that's not converted
	Rule
	MultiLineComment string `json:"MultiLineComment"`
}

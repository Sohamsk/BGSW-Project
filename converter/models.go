package converter

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
	Rule
	Identifier string
	Type       string
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

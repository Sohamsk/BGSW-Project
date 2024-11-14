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
	Symbol string `json:"sym"`
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

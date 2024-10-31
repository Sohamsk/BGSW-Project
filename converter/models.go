package converter

import "encoding/json"

type FileContext struct {
    FileName string
    FileType string
    Body []json.RawMessage `json:"body"`
}

type Rule struct {
    RuleType string
}

type Dim struct {
    Rule
    Identifier string
    Type string
}

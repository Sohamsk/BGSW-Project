package converter

import (
	"encoding/json"
	"fmt"
)

func DeclareVariableRule(content json.RawMessage) string {
    dim := Dim{}
    err := json.Unmarshal(content, &dim)
    if err != nil {
        panic("Error: Incorrect node")
    }
    fmt.Println(dim.Identifier)
    return string(content)
}

var funcMap = map[string]func(json.RawMessage)string{
    "DeclareVariable": DeclareVariableRule,
}

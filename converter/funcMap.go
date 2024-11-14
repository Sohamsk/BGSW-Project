package converter

import (
	"encoding/json"
	"fmt"
	"strings"
)

func incorrectNode() {
	panic("Error: Incorrect node")
}

func incorrectArg() {
	panic("incorrect argument")
}

func DeclareVariableRule(content json.RawMessage) string {
	dim := Dim{}
	err := json.Unmarshal(content, &dim)
	if err != nil {
		incorrectNode()
	}
	fmt.Println(dim.Identifier)
	return string(content)
}

func FuncCallRule(content json.RawMessage) string {
	fun := FuncRule{}
	err := json.Unmarshal(content, &fun)
	if err != nil {
		incorrectNode()
	}
	var sb strings.Builder
	sb.WriteString(fun.Identifier + "(")
	for _, raw := range fun.Arguments {
		arg := ArgType{}
		err := json.Unmarshal(raw, &arg)
		if err != nil {
			incorrectArg()
		}
		if arg.Type == "FunctionCall" {
			sb.WriteString(FuncCallArg(raw))
		} else {
			arg := Literal{}
			json.Unmarshal(raw, &arg)
			if arg.Type == "literal" {
				sb.WriteString("\"" + arg.Symbol + "\",")
			} else {
				sb.WriteString(arg.Symbol + ",")
			}
		}
	}
	return strings.Trim(sb.String(), ",") + ");"
}
func FuncCallArg(content json.RawMessage) string {
	fun := FuncArg{}
	err := json.Unmarshal(content, &fun)
	if err != nil {
		incorrectNode()
	}
	var sb strings.Builder
	sb.WriteString(fun.Identifier + "(")
	for _, raw := range fun.Arguments {
		arg := ArgType{}
		err := json.Unmarshal(raw, &arg)
		if err != nil {
			incorrectArg()
		}
		if arg.Type == "FunctionCall" {
			sb.WriteString(FuncCallArg(raw))
		} else {
			arg := Literal{}
			json.Unmarshal(raw, &arg)
			if arg.Type == "literal" {
				sb.WriteString("\"" + arg.Symbol + "\",")
			} else {
				sb.WriteString(arg.Symbol)
			}
		}
	}
	return strings.Trim(sb.String(), ",") + ")"
}

var funcMap = map[string]func(json.RawMessage) string{
	"DeclareVariable": DeclareVariableRule,
	"FunctionCall":    FuncCallRule,
}

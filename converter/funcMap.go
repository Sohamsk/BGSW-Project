package converter

import (
	"encoding/json"
	"fmt"
	"strings"
)

var funcMap map[string]func(json.RawMessage) string

func init() {
	funcMap = map[string]func(json.RawMessage) string{
		"DeclareVariable": DeclareVariableRule,
		"FunctionCall":    FuncCallRule,
		"expression":      ExpressionRuleHandler,
		"SubStatement":    SubStmtHandler,
		"DoLoopStatement": DoLoopStmtHandler,
	}
}

var vb_cs_types = map[string]string{
	"Boolean":   "bool",
	"Byte":      "byte",
	"Currency":  "decimal",
	"Date":      "DateTime",
	"Double":    "double",
	"Integer":   "short",
	"Long":      "int",
	"Object":    "object",
	"Single":    "float",
	"String":    "string",
	"Variant":   "object",  // Variant usually maps to object
	"Byte()":    "byte[]",  // Byte array
	"Integer()": "short[]", // Integer array
	"Long()":    "int[]",   // Long array
}

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

func ExpressionRuleHandler(content json.RawMessage) string {
	expr := ExpressionRule{}
	err := json.Unmarshal(content, &expr)
	if err != nil {
		incorrectNode()
	}
	var sb strings.Builder
	for _, raw := range expr.Body {
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
			sb.WriteString(arg.Symbol)
		}
	}

	return sb.String() + ";"
}

func SubStmtHandler(content json.RawMessage) string {
	sub := SubStmt{}
	err := json.Unmarshal(content, &sub)
	if err != nil {
		incorrectNode()
	}
	var sb strings.Builder
	if sub.Visibility == "Public" {
		sb.WriteString("public ")
	}
	sb.WriteString("void " + sub.Identifier + "(")
	for _, arg := range sub.Arguments {
		sb.WriteString(vb_cs_types[arg.ArgumentType] + " " + arg.ArgumentName + ",")
	}
	str := sb.String()
	sb.Reset()
	sb.WriteString(strings.Trim(str, ","))
	sb.WriteString(") {")
	sb.WriteString(handleBody(sub.SubBody))
	sb.WriteString("}")
	return sb.String()
}

func ProcessCondition(parts []json.RawMessage) string {
	var sb strings.Builder
	for _, raw := range parts {
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
			sb.WriteString(arg.Symbol)
		}
	}
	return sb.String()
}

func DoLoopStmtHandler(content json.RawMessage) string {
	loop := DoloopStmt{}
	err := json.Unmarshal(content, &loop)
	if err != nil {
		incorrectNode()
	}
	var sb strings.Builder
	if loop.BeforeLoop {
		sb.WriteString("while(")
		if loop.Kind == "until" {
			sb.WriteString("!(")
			sb.WriteString(ProcessCondition(loop.Condition) + "))")
		} else {
			sb.WriteString(ProcessCondition(loop.Condition) + ")")
		}
	} else {
		sb.WriteString("do")
	}
	sb.WriteString("{")
	sb.WriteString(handleBody(loop.Body))
	sb.WriteString("}")
	if !loop.BeforeLoop {
		sb.WriteString("while(")
		if loop.Kind == "until" {
			sb.WriteString("!(")
			sb.WriteString(ProcessCondition(loop.Condition) + "));")
		} else {
			sb.WriteString(ProcessCondition(loop.Condition) + ");")
		}
	}
	return sb.String()
}

package converter

import (
	"encoding/json"
	"fmt"
	"strconv"
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
		"FuncStatement":   FunctionHandler,
		"IfThenElse":      IfThenElseStmtHandler,
		"ElseIf":          ElseIfHandler,
		"ElseBlock":       ElseHandler,
		"ForNextStmt":     ForNextRule,
		"ReturnStatement": ReturnStmtHandler,
		"CommentRule":     CommentHandler,
		"WithStatement":   WithStmtHandler,
	}
}

var vb_cs_types = map[string]string{
	"boolean":   "bool",
	"byte":      "byte",
	"currency":  "decimal",
	"date":      "DateTime",
	"double":    "double",
	"integer":   "short",
	"long":      "int",
	"object":    "object",
	"single":    "float",
	"string":    "string",
	"variant":   "object",  // Variant usually maps to object
	"byte()":    "byte[]",  // Byte array
	"integer()": "short[]", // Integer array
	"long()":    "int[]",   // Long array
}

func incorrectNode() {
	panic("Error: Incorrect node")
}

func incorrectArg() {
	panic("incorrect argument")
}

func DeclareVariableRule(content json.RawMessage) string {
	// Unmarshal the JSON content into a Dim struct
	dim := Dim{}
	err := json.Unmarshal(content, &dim)
	if err != nil {
		// Handle incorrect JSON by calling a predefined error handler
		incorrectNode()
		return ""
	}

	// Convert the type using the mapping
	csType, exists := vb_cs_types[strings.ToLower(dim.Type)]
	if !exists {
		// If type not found, use the original type (might be a custom type)
		csType = dim.Type
	}

	// Construct the C# variable declaration
	var declaration string

	// Handle scope (optional)
	if dim.Scope != "" {
		switch strings.ToLower(dim.Scope) {
		case "public":
			declaration += "public "
		case "private":
			declaration += "private "
		case "protected":
			declaration += "protected "
		case "friend":
			// In C#, friend (internal) accessibility
			declaration += "internal "
		}
	}

	// Handle WithEvents (specific to VB6)
	if dim.WithEvents {
		// In C#, events are typically handled differently
		// This is a placeholder - actual implementation might vary
		declaration += "event "
	}

	// Finalize the declaration
	declaration += fmt.Sprintf("%s %s;", csType, dim.Identifier)

	return declaration
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
				sb.WriteString(arg.Symbol + ",")
			}
		}
	}
	return strings.Trim(sb.String(), ",") + ")"
}

func processExpressions(expr ExpressionRule) string {
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
			if arg.Type == "literal" {
				sb.WriteString("\"" + arg.Symbol + "\"")
			} else {
				sb.WriteString(arg.Symbol)
			}
		}
	}
	return sb.String() + ";"
}

func ExpressionRuleHandler(content json.RawMessage) string {
	expr := ExpressionRule{}
	err := json.Unmarshal(content, &expr)
	if err != nil {
		incorrectNode()
	}
	return processExpressions(expr)
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
		sb.WriteString(vb_cs_types[strings.ToLower(arg.ArgumentType)] + " " + arg.ArgumentName + ",")
	}
	str := sb.String()
	sb.Reset()
	sb.WriteString(strings.Trim(str, ","))
	sb.WriteString(") {")
	sb.WriteString(handleBody(sub.SubBody))
	sb.WriteString("}")
	return sb.String()
}

func handleBodyFunc(rules []json.RawMessage, name, returnType string) string {
	fmt.Println("in func")
	var result string
	result += vb_cs_types[strings.ToLower(returnType)] + " " + name + ";"
	for _, rule := range rules {
		result += ConvertRule(rule)
	}
	result += "return " + name + ";"
	return result
}

func FunctionHandler(content json.RawMessage) string {
	funct := FuncDecl{}
	err := json.Unmarshal(content, &funct)
	if err != nil {
		incorrectNode()
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s %s %s (", strings.ToLower(funct.Visibility), vb_cs_types[strings.ToLower(funct.ReturnType)], funct.Identifier))
	for _, arg := range funct.Arguments {
		sb.WriteString(vb_cs_types[strings.ToLower(arg.ArgumentType)] + " " + arg.ArgumentName + ",")
	}
	str := sb.String()
	sb.Reset()
	sb.WriteString(strings.Trim(str, ",") + "){" + handleBodyFunc(funct.Body, funct.Identifier, funct.ReturnType) + "}") // need a seperate body handler to handle functions returning values as there is no return keyword in vb6
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
			if arg.Symbol == "=" {
				sb.WriteString("==")
			} else {
				sb.WriteString(arg.Symbol)
			}
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
func IfThenElseStmtHandler(content json.RawMessage) string {
	var ifStmt IfThenElseStmtRule
	err := json.Unmarshal(content, &ifStmt)
	if err != nil {
		incorrectNode()
	}

	var sb strings.Builder

	// Handle the `IfThenElseStmtRule`
	sb.WriteString("if (")
	sb.WriteString(ProcessCondition(ifStmt.Condition)) // Using handleBody for Condition
	sb.WriteString(") {\n")
	sb.WriteString(handleBody(ifStmt.IfBlock)) // Using handleBody for IfBlock
	sb.WriteString("\n}")

	return sb.String()
}

func ElseIfHandler(content json.RawMessage) string {
	var elseIfStmt ElseIfRule
	err := json.Unmarshal(content, &elseIfStmt)
	if err != nil {
		incorrectNode()
	}

	var sb strings.Builder

	// Handle the `ElseIfRule`
	sb.WriteString("else if (")
	sb.WriteString(ProcessCondition(elseIfStmt.Condition)) // Using handleBody for Condition
	sb.WriteString(") {\n")
	sb.WriteString(handleBody(elseIfStmt.ElseIfBlock)) // Using handleBody for ElseIfBlock
	sb.WriteString("\n}")

	return sb.String()
}

func ElseHandler(content json.RawMessage) string {
	var elseStmt ElseRule
	err := json.Unmarshal(content, &elseStmt)
	if err != nil {
		incorrectNode()
	}

	var sb strings.Builder

	// Handle the `ElseRule`
	sb.WriteString("else {\n")
	sb.WriteString(handleBody(elseStmt.Body)) // Using handleBody for ElseBlock
	sb.WriteString("\n}")
	return sb.String()
}
func ForNextRule(content json.RawMessage) string {
	forNext := ForNext{}
	err := json.Unmarshal(content, &forNext)
	if err != nil {
		panic("Error: Incorrect node")
	}

	// Convert start, end, and step to integers
	start, _ := strconv.Atoi(forNext.Start)
	end, _ := strconv.Atoi(forNext.End)
	step, _ := strconv.Atoi(forNext.Step)
	// Generate the loop representation
	// Note: C# uses <= for the condition, similar to VB6
	loop := fmt.Sprintf("for(int %s = %d; %s <= %d; %s += %d)",
		forNext.IdentifierName,
		start,
		forNext.IdentifierName,
		end,
		forNext.IdentifierName,
		step,
	)

	var sb strings.Builder
	sb.WriteString(loop)
	sb.WriteString(" {\n")
	sb.WriteString(handleBody(forNext.Body))
	sb.WriteString("\n}")

	return sb.String()
}

func ReturnStmtHandler(content json.RawMessage) string {
	ret := ReturnStmt{}
	err := json.Unmarshal(content, &ret)
	if err != nil {
		incorrectNode()
	}

	return "return " + ret.ReturnVariableName + ";"
}

func CommentHandler(content json.RawMessage) string {
	comment := Comment{}
	err := json.Unmarshal(content, &comment)
	if err != nil {
		panic("Error: Incorrect node")
	}
	return fmt.Sprintf("// %s\n", comment.CommentText)
}

func handleBodyWith(expressions []ExpressionRule, objectName string) string {
	var result string
	for _, expression := range expressions {
		result += objectName + processExpressions(expression)
	}
	return result
}

func WithStmtHandler(content json.RawMessage) string {
	withStmt := WithStmt{}
	err := json.Unmarshal(content, &withStmt)
	if err != nil {
		incorrectNode()
	}
	return handleBodyWith(withStmt.Body, withStmt.Object)
}

package converter

import (
	"bosch/converter/models"
	"encoding/json"
	"fmt"
	"log"
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
		"SetStatement":    SetStatementHandler,
		"UnhandledRule":   MultiLineCommentHandler,
		"EnumerationRule": EnumsHandler,
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
	log.Println("Error: Incorrect node")
}

func incorrectArg() {
	log.Println("incorrect argument")
}

func DeclareVariableRule(content json.RawMessage) string {
	// Unmarshal the JSON content into a Dim struct
	dim := models.Dim{}
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
	fun := models.FuncRule{}
	err := json.Unmarshal(content, &fun)
	if err != nil {
		incorrectNode()
		return ""
	}
	var sb strings.Builder
	sb.WriteString(fun.Identifier + "(")
	for _, raw := range fun.Arguments {
		arg := models.ArgType{}
		err := json.Unmarshal(raw, &arg)
		if err != nil {
			incorrectArg()
		}
		if arg.Type == "FunctionCall" {
			sb.WriteString(FuncCallArg(raw))
		} else {
			arg := models.Literal{}
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
	fun := models.FuncArg{}
	err := json.Unmarshal(content, &fun)
	if err != nil {
		incorrectNode()
		return ""
	}
	var sb strings.Builder
	sb.WriteString(fun.Identifier + "(")
	for _, raw := range fun.Arguments {
		arg := models.ArgType{}
		err := json.Unmarshal(raw, &arg)
		if err != nil {
			incorrectArg()
		}
		if arg.Type == "FunctionCall" {
			sb.WriteString(FuncCallArg(raw))
		} else {
			arg := models.Literal{}
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

func processExpressions(expr models.ExpressionRule) string {
	var sb strings.Builder
	for _, raw := range expr.Body {
		arg := models.ArgType{}
		err := json.Unmarshal(raw, &arg)
		if err != nil {
			incorrectArg()
		}
		if arg.Type == "FunctionCall" {
			sb.WriteString(FuncCallArg(raw))
		} else {
			arg := models.Literal{}
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
	expr := models.ExpressionRule{}
	err := json.Unmarshal(content, &expr)
	if err != nil {
		incorrectNode()
		return ""
	}
	return processExpressions(expr)
}

func getVisibility(visibility string, sb *strings.Builder) {
	if strings.ToLower(visibility) == "friend" {
		sb.WriteString("internal ")
	} else if visibility == "" {
		if strings.ToLower(global.FileType) == "frm" {
			sb.WriteString("private ")
		} else {
			sb.WriteString("public ")
		}
	} else {
		sb.WriteString(strings.ToLower(visibility) + " ")
	}
}

func SubStmtHandler(content json.RawMessage) string {
	sub := models.SubStmt{}
	err := json.Unmarshal(content, &sub)
	if err != nil {
		incorrectNode()
		return ""
	}
	var sb strings.Builder
	getVisibility(sub.Visibility, &sb)
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
		inter, err := ConvertRule(rule)
		if err == nil {
			result += inter
		}
	}
	result += "return " + name + ";"
	return result
}

func FunctionHandler(content json.RawMessage) string {
	funct := models.FuncDecl{}
	err := json.Unmarshal(content, &funct)
	if err != nil {
		incorrectNode()
		return ""
	}
	var sb strings.Builder
	getVisibility(funct.Visibility, &sb)
	sb.WriteString(fmt.Sprintf("%s %s (", vb_cs_types[strings.ToLower(funct.ReturnType)], funct.Identifier))
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
		arg := models.ArgType{}
		err := json.Unmarshal(raw, &arg)
		if err != nil {
			incorrectArg()
		}
		if arg.Type == "FunctionCall" {
			sb.WriteString(FuncCallArg(raw))
		} else {
			arg := models.Literal{}
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
	loop := models.DoloopStmt{}
	err := json.Unmarshal(content, &loop)
	if err != nil {
		incorrectNode()
		return ""
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
	var ifStmt models.IfThenElseStmtRule
	err := json.Unmarshal(content, &ifStmt)
	if err != nil {
		incorrectNode()
		return ""
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
	var elseIfStmt models.ElseIfRule
	err := json.Unmarshal(content, &elseIfStmt)
	if err != nil {
		incorrectNode()
		return ""
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
	var elseStmt models.ElseRule
	err := json.Unmarshal(content, &elseStmt)
	if err != nil {
		incorrectNode()
		return ""
	}

	var sb strings.Builder

	// Handle the `ElseRule`
	sb.WriteString("else {\n")
	sb.WriteString(handleBody(elseStmt.Body)) // Using handleBody for ElseBlock
	sb.WriteString("\n}")
	return sb.String()
}
func ForNextRule(content json.RawMessage) string {
	forNext := models.ForNext{}
	err := json.Unmarshal(content, &forNext)
	if err != nil {
		incorrectNode()
		return ""
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
	ret := models.ReturnStmt{}
	err := json.Unmarshal(content, &ret)
	if err != nil {
		incorrectNode()
		return ""
	}

	return "return " + ret.ReturnVariableName + ";"
}

func CommentHandler(content json.RawMessage) string {
	comment := models.Comment{}
	err := json.Unmarshal(content, &comment)
	if err != nil {
		incorrectNode()
		return ""
	}
	return fmt.Sprintf("// %s\n", comment.CommentText)
}

func MultiLineCommentHandler(content json.RawMessage) string {
	MultiLineComment := models.MultiLineComment{}
	err := json.Unmarshal(content, &MultiLineComment)
	if err != nil {
		panic("Error: Incorrect node")
	}
	return fmt.Sprintf("/* \n" + MultiLineComment.MultiLineComment + "\n*/\n")
}
func handleBodyWith(expressions []models.ExpressionRule, objectName string) string {
	var result string
	for _, expression := range expressions {
		result += objectName + processExpressions(expression)
	}
	return result
}

func WithStmtHandler(content json.RawMessage) string {
	withStmt := models.WithStmt{}
	err := json.Unmarshal(content, &withStmt)
	if err != nil {
		incorrectNode()
		return ""
	}
	return handleBodyWith(withStmt.Body, withStmt.Object)
}

func SetStatementHandler(content json.RawMessage) string {
	set := models.SetStmt{}
	err := json.Unmarshal(content, &set)
	if err != nil {
		incorrectNode()
		return ""
	}

	var class string
	if set.IsNew {
		class = fmt.Sprintf("new %s()", strings.Trim(string(set.Class), "\""))
	} else {
		class = strings.Trim(string(set.Class), "\"")
	}

	return fmt.Sprintf("%s = %s;", set.Identifier, class)
}

func EnumsHandler(content json.RawMessage) string {
	enumStmt := models.EnumStmt{}
	err := json.Unmarshal(content, &enumStmt)
	if err != nil {
		panic("Error: Incorrect node")
	}

	// Start building the enum string with the enum name
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("public enum %s//change visibility as per requirement\n{\n", enumStmt.Name))

	// Process each enum value
	for i, value := range enumStmt.EnumValues {
		// Convert the value to a string since it's currently type any
		valueStr := fmt.Sprintf("%v", value)

		// Add comma for all elements except the last one
		if i < len(enumStmt.EnumValues)-1 {
			builder.WriteString(fmt.Sprintf("    %s,\n", valueStr))
		} else {
			builder.WriteString(fmt.Sprintf("    %s\n", valueStr))
		}
	}

	builder.WriteString("}\n")
	return builder.String()
}

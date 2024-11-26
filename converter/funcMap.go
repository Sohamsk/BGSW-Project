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

func ForNextRule(content json.RawMessage) string {
	forNext := ForNext{}
	err := json.Unmarshal(content, &forNext)
	if err != nil {
		panic("Error: Incorrect node")
	}
	// Generate the loop representation
	loop := fmt.Sprintf("for(int %s=%d; %s<=%d; %s+=%d)",
		forNext.IdentifierName,
		forNext.Start,
		forNext.IdentifierName,
		forNext.End,
		forNext.IdentifierName,
		forNext.Step)

	// Print the loop representation
	fmt.Println(loop)
	return loop
}

var funcMap = map[string]func(json.RawMessage) string{
	"DeclareVariable": DeclareVariableRule,
	"ForNext":         ForNextRule,
}

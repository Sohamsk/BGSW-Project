package converter

import (
	"encoding/json"
	"fmt"
)

type RawItem struct {
	Type string `json:"ruletype"`
}

func handleBody(rules []json.RawMessage) string {
	var result string
	for _, rule := range rules {
		result += ConvertRule(rule)
	}
	return result
}

func Convert(raw string) string {
	context := FileContext{}
	err := json.Unmarshal([]byte(raw), &context)
	if err != nil {
		panic("Error: Error unmarshalling json")
	}
	fmt.Println("Name: " + context.FileName + " Type: " + context.FileType)
	fmt.Println(handleBody(context.Body))
	return raw
}

// the converter should take the json string and the project context which we'll get on parsing the vbp file
func ConvertRule(rawMsg json.RawMessage) string {
	raw := RawItem{}
	err := json.Unmarshal([]byte(rawMsg), &raw)
	if err != nil {
		panic("Error: Error unmarshalling json")
	}

	action, ok := funcMap[raw.Type]
	if !ok {
		panic(raw.Type)
	}
	return action(rawMsg)
}

package converter

import (
	"bosch/converter/models"
	"encoding/json"
	"fmt"
)

func handleBody(rules []json.RawMessage) string {
	var result string
	for _, rule := range rules {
		result += ConvertRule(rule)
	}
	return result
}

func Convert(raw string) string {
	context := models.FileContext{}
	err := json.Unmarshal([]byte(raw), &context)
	if err != nil {
		panic("Error: Error unmarshalling json")
	}
	converted := fmt.Sprintf("class %s {%s}", context.FileName, handleBody(context.Body))
	return converted
}

// the converter should take the json string and the project context which we'll get on parsing the vbp file
func ConvertRule(rawMsg json.RawMessage) string {
	raw := models.Rule{}
	err := json.Unmarshal([]byte(rawMsg), &raw)
	if err != nil {
		panic("Error: Error unmarshalling json")
	}

	action, ok := funcMap[raw.RuleType]
	if !ok {
		panic(raw.RuleType + " is unknown")
	}
	return action(rawMsg)
}

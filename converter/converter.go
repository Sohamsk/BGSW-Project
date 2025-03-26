package converter

import (
	"bosch/converter/models"
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

type state struct {
	FileType string
	Symtab   map[string]string
}

var global state

func handleBody(rules []json.RawMessage) string {
	var result string
	for _, rule := range rules {
		raw := models.Rule{}
		json.Unmarshal(rule, &raw)
		
		// Skip functions in Execute()
		if raw.RuleType == "FuncStatement" {
			continue
		}

		inter, err := ConvertRule(rule)
		if err == nil {
			result += inter + "\n"
		}
	}
	return result
}



func Convert(raw string, symtab map[string]string) (string, error) {
	context := models.FileContext{}
	err := json.Unmarshal([]byte(raw), &context)
	global.FileType = context.FileType
	global.Symtab = symtab
	if err != nil {
		return "", fmt.Errorf("Error: %s", err)
	}

	// Separate functions and statements
	functions := handleFunctions(context.Body)  
	otherStatements := handleBody(context.Body) 

	// Build the C# class with functions outside execute
	converted := fmt.Sprintf(`using System;
class %s { 
%s 

public static void Execute() { 
%s 
} 

public static void Main(string[] args) { 
Execute(); 
}
}`, context.FileName, functions, otherStatements)

	return converted, nil
}

func handleFunctions(rules []json.RawMessage) string {
	var result string
	for _, rule := range rules {
		raw := models.Rule{}
		json.Unmarshal(rule, &raw)
		
		// Check if this rule is a function (based on heuristic)
		if raw.RuleType == "FuncStatement" {
			inter, err := ConvertRule(rule)
			if err == nil {
				result += inter + "\n" // Add the function separately
			}
		}
	}
	return result
}

// the converter should take the json string and the project context which we'll get on parsing the vbp file
func ConvertRule(rawMsg json.RawMessage) (string, error) {
	raw := models.Rule{}
	err := json.Unmarshal([]byte(rawMsg), &raw)
	if err != nil {
		error := errors.New("Error: Error unmarshalling json")
		log.Println(error)
		return "", error
	}

	action, ok := funcMap[raw.RuleType]
	if !ok {
		error := errors.New("Error:" + raw.RuleType + " is unknown")
		log.Println(error)
		return "", error
	}
	return action(rawMsg), nil
}

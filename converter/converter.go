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

func handleBody(rules []json.RawMessage, fromSub bool) string {
	var result string
	for _, rule := range rules {
		raw := models.Rule{}
		json.Unmarshal(rule, &raw)
		
		// Skip functions in Execute()
		if raw.RuleType == "FuncStatement" || raw.RuleType=="SubStatement" || raw.RuleType=="EnumerationRule"{
			continue
		}
		if !fromSub && raw.RuleType == "DeclareVariable" {
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
		return "", err
	}

	// Separate functions and statements
	functions := handleoutsidestmts(context.Body)  
	otherStatements := handleBody(context.Body, false) 

	// Build the C# class with functions outside execute
	converted := fmt.Sprintf(`using System;

class %s { 
%s 

public static void Execute() { 
%s 
} 

public static void Main(string[] args) { 

}
}`, context.FileName, functions, otherStatements)

	return converted, nil
}


func handleoutsidestmts(rules []json.RawMessage) string {
	var result string
	for _, rule := range rules {
		raw := models.Rule{}
		json.Unmarshal(rule, &raw)
		
		// Check if this rule is a function (based on heuristic)
		if raw.RuleType == "FuncStatement" || raw.RuleType=="SubStatement" || raw.RuleType == "DeclareVariable" || raw.RuleType=="EnumerationRule"{
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
		error := errors.New("error: Error unmarshalling json")
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

package converter

import (
	"encoding/json"
	"fmt"
)

type Item interface{}

type RawItem struct {
	Type string          `json:"ruletype"`
	Raw  json.RawMessage `json:"-"`
}

// the converter should take the json string and the project context which we'll get on parsing the vbp file
func ConvertRules() {
    rawString := `{"RuleType": "DeclareVariable","Identifier": "x","Type": "INTEGER"}`
    raw := RawItem{}
    json.Unmarshal([]byte(rawString), &raw)

    switch raw.Type {
        case "DeclareVariable":
            fmt.Println("inside a DeclareVariable")
        default:
            fmt.Println("Error unknown rule")
    }
    fmt.Println(raw.Type)
}

package converter

type Rule struct {
    RuleType string
}

type Dim struct {
    Rule
    Identifier string
    Type string
}

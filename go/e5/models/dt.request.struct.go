package models

type DTRequestStruct struct {
	DTRules []DTRule `json:"dt_rules,omitempty"`
}

type DTRule struct {
	RuleID      string    `json:"rule_id,omitempty"`
	Description string    `json:"description,omitempty"`
	Severity    string    `json:"severity,omitempty"`
	Action      string    `json:"action,omitempty"`
	Condition   Condition `json:"condition,omitempty"`
}

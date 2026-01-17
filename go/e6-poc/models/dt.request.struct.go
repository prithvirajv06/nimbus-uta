package models

type DTRequestStruct struct {
	DTRules []DTRule `json:"dt_rules,omitempty"`
}

type DTRule struct {
	RuleID      string      `json:"rule_id,omitempty"`
	Description string      `json:"description,omitempty"`
	Severity    string      `json:"severity,omitempty"`
	Action      []Action    `json:"action,omitempty"`
	Condition   []Condition `json:"condition,omitempty"`
}

type Condition struct {
	Operator string    `json:"operator,omitempty"`
	Operands Variables `json:"operands,omitempty"`
}

type Action struct {
	Type   string    `json:"type,omitempty"`
	Target Variables `json:"target,omitempty"`
	Value  string    `json:"value,omitempty"`
}

func NewDTRequestStruct() *DTRequestStruct {
	return &DTRequestStruct{
		DTRules: []DTRule{},
	}
}

func (dtr *DTRequestStruct) AddDTRule(rule DTRule) {
	dtr.DTRules = append(dtr.DTRules, rule)
}

type Variables struct {
	VarKey             string      `bson:"var_key" json:"var_key"`
	ContextVarKey      string      `bson:"context_var_key" json:"context_var_key"`
	ContextVarToCreate string      `bson:"context_var_to_create" json:"context_var_to_create"`
	Label              string      `bson:"label" json:"label"`
	Type               string      `bson:"type" json:"type"`
	IsRequired         bool        `bson:"is_required" json:"is_required"`
	Value              interface{} `bson:"value" json:"value"`
	Children           []Variables `bson:"children,omitempty" json:"children,omitempty"`
	IsClickable        bool        `bson:"_" json:"_"`
}

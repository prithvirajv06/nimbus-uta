package models

// WorkflowStep represents a node in the workflow graph
type WorkflowStep struct {
	StepID         string  `json:"step_id"`
	Type           string  `json:"type"` // use interface{} if < Go 1.18
	Label          string  `json:"label"`
	Icon           string  `json:"icon"`
	Target         string  `json:"target"`
	Value          *string `json:"value,omitempty"`
	Statement      *string `json:"statement,omitempty"`
	StatementLabel *string `json:"statementLabel,omitempty"`

	Children      []WorkflowStep `json:"children,omitempty"`
	TrueChildren  []WorkflowStep `json:"true_children,omitempty"`  // IF branch
	FalseChildren []WorkflowStep `json:"false_children,omitempty"` // ELSE branch

	IsOpen     *bool   `json:"isOpen,omitempty"`
	ContextVar *string `json:"context_var,omitempty"`

	ConditionConfig []ConditionConfig `json:"condition_config,omitempty"`
}

type LogicalOperator string

const (
	AND LogicalOperator = "AND"
	OR  LogicalOperator = "OR"
)

type ConditionConfig struct {
	LeftVar         Variables       `json:"left_var"`
	Operator        string          `json:"operator"`
	RightValue      any             `json:"right_value"`
	PreceedingLogic LogicalOperator `json:"preceeding_logic"`
}

func (w *WorkflowStep) IsConditional() bool {
	return len(w.TrueChildren) > 0 || len(w.FalseChildren) > 0
}

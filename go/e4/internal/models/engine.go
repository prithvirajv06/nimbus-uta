package models

// WorkflowStep represents a node in the workflow graph
type WorkflowStep struct {
	StepID         string     `json:"step_id,omitempty" bson:"step_id,omitempty"`
	Type           string     `json:"type,omitempty" bson:"type,omitempty"` // use interface{} if < Go 1.18
	Label          string     `json:"label,omitempty" bson:"label,omitempty"`
	Icon           string     `json:"icon,omitempty" bson:"icon,omitempty"`
	Target         *Variables `json:"target,omitempty" bson:"target,omitempty"`
	Value          *string    `json:"value,omitempty" bson:"value,omitempty"`
	Statement      *string    `json:"statement,omitempty" bson:"statement,omitempty"`
	StatementLabel *string    `json:"statement_label,omitempty" bson:"statement_label,omitempty"`

	Children      []WorkflowStep  `json:"children,omitempty" bson:"children,omitempty"`
	TrueChildren  *[]WorkflowStep `json:"true_children,omitempty" bson:"true_children,omitempty"`   // IF branch
	FalseChildren *[]WorkflowStep `json:"false_children,omitempty" bson:"false_children,omitempty"` // ELSE branch
	IsOpen        *bool           `json:"is_open,omitempty" bson:"is_open,omitempty"`
	ContextVar    *string         `json:"context_var,omitempty" bson:"context_var,omitempty"`

	ConditionConfig *[]ConditionConfig `json:"condition_config,omitempty" bson:"condition_config,omitempty"`
}

type LogicalOperator string

const (
	AND LogicalOperator = "AND"
	OR  LogicalOperator = "OR"
)

type ConditionConfig struct {
	LeftVar         Variables       `json:"left_var,omitempty" bson:"left_var,omitempty"`
	Operator        string          `json:"operator,omitempty" bson:"operator,omitempty"`
	RightValue      any             `json:"right_value,omitempty" bson:"right_value,omitempty"`
	PreceedingLogic LogicalOperator `json:"preceeding_logic,omitempty" bson:"preceeding_logic,omitempty"`
}

func (w *WorkflowStep) IsConditional() bool {
	return len(*w.TrueChildren) > 0 || len(*w.FalseChildren) > 0
}

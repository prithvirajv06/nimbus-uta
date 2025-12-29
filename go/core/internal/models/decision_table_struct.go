package models

type DecisionTable struct {
	NIMB_ID         string          `bson:"nimb_id" json:"nimb_id"`
	Description     string          `json:"description" gorm:"type:text" bson:"description"`
	NoOfRows        int             `json:"no_of_rows" gorm:"not null" bson:"no_of_rows"` //
	NoOfInputs      int             `json:"no_of_inputs" gorm:"not null" bson:"no_of_inputs"`
	NoOfOutputs     int             `json:"no_of_outputs" gorm:"not null" bson:"no_of_outputs"` //
	Name            string          `json:"name" gorm:"type:varchar(255);not null" bson:"name"`
	HitPolicy       string          `json:"hit_policy" gorm:"type:varchar(50)" bson:"hit_policy"` // UNIQUE, FIRST, PRIORITY, COLLECT, COLLECT_SUM, COLLECT_COUNT
	Inputs          []TableInput    `json:"inputs" gorm:"type:jsonb" bson:"inputs"`
	Outputs         []TableOutput   `json:"outputs" gorm:"type:jsonb" bson:"outputs"`
	VariablePackage VariablePackage `json:"variable_package" gorm:"type:jsonb" bson:"variable_package"`
	Rules           [][]string      `json:"rules" gorm:"type:jsonb" bson:"rules"` // Rows of input criteria followed by output values
}

type TableInput struct {
	Variable string `json:"variable" gorm:"type:varchar(255);not null" bson:"variable"`
	Label    string `json:"label" gorm:"type:varchar(255);not null" bson:"label"`
}

type TableOutput struct {
	Variable      string   `json:"variable" gorm:"type:varchar(255);not null" bson:"variable"`
	Label         string   `json:"label" gorm:"type:varchar(255);not null" bson:"label"`
	AllowedValues []string `json:"allowed_values,omitempty" bson:"allowed_values,omitempty"`
	IsPriority    bool     `json:"is_priority" gorm:"default:false" bson:"is_priority"`
}

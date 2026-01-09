package models

type DecisionTable struct {
	NIMB_ID         string          `bson:"nimb_id" json:"nimb_id"`
	Description     string          `json:"description" gorm:"type:text" bson:"description"`
	NoOfRows        int             `json:"no_of_rows" gorm:"not null" bson:"no_of_rows"` //
	NoOfInputs      int             `json:"no_of_inputs" gorm:"not null" bson:"no_of_inputs"`
	NoOfOutputs     int             `json:"no_of_outputs" gorm:"not null" bson:"no_of_outputs"` //
	Name            string          `json:"name" gorm:"type:varchar(255);not null" bson:"name"`
	HitPolicy       string          `json:"hit_policy" gorm:"type:varchar(50)" bson:"hit_policy"` // UNIQUE, FIRST, PRIORITY, COLLECT, COLLECT_SUM, COLLECT_COUNT
	InputsColumns   []Variables     `json:"input_columns" gorm:"type:jsonb" bson:"input_columns"`
	OutputsColumns  []Variables     `json:"output_columns" gorm:"type:jsonb" bson:"output_columns"`
	VariablePackage VariablePackage `json:"variable_package" gorm:"type:jsonb" bson:"variable_package"`
	Rules           [][]Variables   `json:"rules" gorm:"type:jsonb" bson:"rules"` // Rows of input criteria followed by output values
	Audit           Audit           `json:"audit" gorm:"type:jsonb" bson:"audit"`
}

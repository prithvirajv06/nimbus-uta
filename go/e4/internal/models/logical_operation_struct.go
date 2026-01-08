package models

type LogicFlow struct {
	NIMB_ID         string          `bson:"nimb_id" json:"nimb_id"`
	Name            string          `bson:"name" json:"name"`
	Description     string          `bson:"description" json:"description"`
	Active          bool            `bson:"active" json:"active"`
	NoOfBranches    int             `bson:"no_of_branches" json:"no_of_branches"`
	VariablePackage VariablePackage `bson:"variable_package" json:"variable_package"`
	Steps           []LogicalStep   `bson:"steps" json:"steps"`
	Audit           Audit           `bson:"audit" json:"audit"`
}

type LogicalStep struct {
	OperationName    string      `json:"operation_name" gorm:"type:varchar(255)" bson:"operation_name"`
	Condition        Condition   `json:"condition" gorm:"type:jsonb" bson:"condition"`
	Operator         string      `json:"operator" gorm:"type:varchar(50)" bson:"operator"`
	OperationIfTrue  []Operation `json:"operation_if_true" gorm:"type:jsonb" bson:"operation_if_true"`
	OperationIfFalse []Operation `json:"operation_if_false" gorm:"type:jsonb" bson:"operation_if_false"`
}

type Condition struct {
	Operator     string        `json:"operator" gorm:"type:varchar(50)" bson:"operator"`
	Variable     Variables     `json:"variable" gorm:"type:jsonb" bson:"variable"`
	Logical      string        `json:"logical" gorm:"type:varchar(50)" bson:"logical"`
	OpValue      interface{}   `json:"op_value" gorm:"type:jsonb" bson:"op_value"`
	ArrayFilters []ArrayFilter `json:"array_filters,omitempty" gorm:"type:jsonb" bson:"array_filters,omitempty"`

	Conditions []Condition `json:"conditions" gorm:"type:jsonb" bson:"conditions"`
}

type Operation struct {
	Variable     Variables     `json:"variable" gorm:"type:varchar(255)" bson:"variable"`
	Operation    string        `json:"operation" gorm:"type:varchar(50)" bson:"operation"`
	OpValue      interface{}   `json:"op_value" gorm:"type:jsonb" bson:"op_value"`
	ValueIsPath  bool          `json:"value_is_path" gorm:"default:false" bson:"value_is_path"`
	ArrayFilters []ArrayFilter `json:"array_filters,omitempty" gorm:"type:jsonb" bson:"array_filters,omitempty"`
}

type ArrayFilter struct {
	ArrayName string      `json:"array_name" gorm:"type:varchar(255)" bson:"array_name"`
	Property  string      `json:"property" gorm:"type:varchar(255)" bson:"property"`
	Logical   string      `json:"logical" gorm:"type:varchar(50)" bson:"logical"`
	OpValue   interface{} `json:"op_value" gorm:"type:jsonb" bson:"op_value"`
}

package models

type LogicFlow struct {
	ID              string          `bson:"_id,omitempty" json:"id"`
	Name            string          `bson:"name" json:"name"`
	Description     string          `bson:"description" json:"description"`
	Active          bool            `bson:"active" json:"active"`
	NoOfBranches    int             `bson:"no_of_branches" json:"no_of_branches"`
	VariablePackage VariablePackage `bson:"variable_package" json:"variable_package"`
	LogicalSteps    []LogicalStep   `bson:"logical_steps" json:"logical_steps"`
	Audit           Audit           `bson:"audit" json:"audit"`
}

type LogicalStep struct {
	Name             string      `json:"name" gorm:"type:varchar(255)" bson:"name"`
	Condition        Condition   `json:"condition" gorm:"type:jsonb" bson:"condition"`
	Variable         string      `json:"variable" gorm:"type:varchar(255)" bson:"variable"`
	Logical          string      `json:"logical" gorm:"type:varchar(50)" bson:"logical"`
	Value            interface{} `json:"value" gorm:"type:jsonb" bson:"value"`
	OperationIfTrue  []Operation `json:"operation_if_true" gorm:"type:jsonb" bson:"operation_if_true"`
	OperationIfFalse []Operation `json:"operation_if_false" gorm:"type:jsonb" bson:"operation_if_false"`
}

type Condition struct {
	Operator   string      `json:"operator" gorm:"type:varchar(50)" bson:"operator"`
	Conditions []Condition `json:"conditions" gorm:"type:jsonb" bson:"conditions"`
	Variable   string      `json:"variable" gorm:"type:varchar(255)" bson:"variable"`
	Logical    string      `json:"logical" gorm:"type:varchar(50)" bson:"logical"`
	Value      interface{} `json:"value" gorm:"type:jsonb" bson:"value"`
}

type Operation struct {
	Variable     string        `json:"variable" gorm:"type:varchar(255)" bson:"variable"`
	Operation    string        `json:"operation" gorm:"type:varchar(50)" bson:"operation"`
	Value        interface{}   `json:"value" gorm:"type:jsonb" bson:"value"`
	ValueIsPath  bool          `json:"value_is_path" gorm:"default:false" bson:"value_is_path"`
	Type         string        `json:"type" gorm:"type:varchar(50)" bson:"type"` // e.g., 'NUMBER' | 'STRING' | 'BOOLEAN' | 'JSON' | 'DATE'
	ArrayFilters []ArrayFilter `json:"array_filters,omitempty" gorm:"type:jsonb" bson:"array_filters,omitempty"`
}

type ArrayFilter struct {
	ArrayName string      `json:"array_name" gorm:"type:varchar(255)" bson:"array_name"`
	Property  string      `json:"property" gorm:"type:varchar(255)" bson:"property"`
	Logical   string      `json:"logical" gorm:"type:varchar(50)" bson:"logical"`
	Value     interface{} `json:"value" gorm:"type:jsonb" bson:"value"`
}

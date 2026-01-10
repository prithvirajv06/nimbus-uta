package models

type LogicFlow struct {
	NIMB_ID         string          `bson:"nimb_id" json:"nimb_id"`
	Name            string          `bson:"name" json:"name"`
	Description     string          `bson:"description" json:"description"`
	Active          bool            `bson:"active" json:"active"`
	NoOfBranches    int             `bson:"no_of_branches" json:"no_of_branches"`
	VariablePackage VariablePackage `bson:"variable_package" json:"variable_package"`
	LogicalSteps    []WorkflowDef   `bson:"logical_steps" json:"logical_steps"`
	Audit           Audit           `bson:"audit" json:"audit"`
}

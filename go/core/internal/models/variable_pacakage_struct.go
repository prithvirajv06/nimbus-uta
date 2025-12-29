package models

type VariablePackage struct {
	ID          string      `bson:"_id,omitempty" json:"id"`
	PackageName string      `bson:"package_name" json:"package_name"`
	Description string      `bson:"description" json:"description"`
	Variables   []Variables `bson:"variables" json:"variables"`
}

type Variables struct {
	VarKey     string `bson:"var_key" json:"var_key"`
	Label      string `bson:"label" json:"label"`
	Type       string `bson:"type" json:"type"`
	IsRequired bool   `bson:"is_required" json:"is_required"`
}

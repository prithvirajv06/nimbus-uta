package models

type VariablePackageRequet struct {
	PackageName string `bson:"package_name" json:"package_name"`
	Description string `bson:"description" json:"description"`
	JSONStr     string `bson:"json_str" json:"json_str"`
}

type VariablePackage struct {
	NIMB_ID     string      `bson:"nimb_id" json:"nimb_id"`
	PackageName string      `bson:"package_name" json:"package_name"`
	Description string      `bson:"description" json:"description"`
	Variables   []Variables `bson:"variables" json:"variables"`
	Audit       Audit       `bson:"audit" json:"audit"`
}

type Variables struct {
	VarKey     string `bson:"var_key" json:"var_key"`
	Label      string `bson:"label" json:"label"`
	Type       string `bson:"type" json:"type"`
	IsRequired bool   `bson:"is_required" json:"is_required"`
}

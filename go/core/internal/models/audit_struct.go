package models

type Audit struct {
	CreatedAt       int64  `bson:"created_at" json:"created_at"`
	CreatedBy       string `bson:"created_by" json:"created_by"`
	ModifiedAt      int64  `bson:"modified_at" json:"modified_at"`
	ModifiedBy      string `bson:"modified_by" json:"modified_by"`
	IsProdCandidate bool   `bson:"is_prod_candidate" json:"is_prod_candidate"`

	//Status fields
	Active     string `bson:"active" json:"active"`
	IsArchived bool   `bson:"is_archived" json:"is_archived,omitempty"`

	//Versioning fields
	Version      int `bson:"version" json:"version,omitempty"`
	MinorVersion int `bson:"minor_version" json:"minor_version,omitempty"`
}

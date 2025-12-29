package models

type Audit struct {
	CreatedAt       int64  `bson:"created_at" json:"created_at"`
	CreatedBy       string `bson:"created_by" json:"created_by"`
	ModifiedAt      int64  `bson:"modified_at" json:"modified_at"`
	ModifiedBy      string `bson:"modified_by" json:"modified_by"`
	IsProdCandicate bool   `bson:"is_prod_candidate" json:"is_prod_candidate"`
}

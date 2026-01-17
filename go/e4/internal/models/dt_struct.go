package models

type HitPolicy string

const (
	First   HitPolicy = "FIRST"   // Execute only the first match
	Collect HitPolicy = "COLLECT" // Execute all matches
	Unique  HitPolicy = "UNIQUE"  // Error if more than one match
)

type DTAction struct {
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

type DTRule struct {
	NIMB_ID    string            `json:"nimb_id"`
	ID         string            `json:"id"`
	Priority   int               `json:"priority"` // Used for sorting or selection
	Conditions map[string]string `json:"conditions"`
	Actions    []DTAction        `json:"actions"`
	Audit      Audit             `bson:"audit" json:"audit"`
}

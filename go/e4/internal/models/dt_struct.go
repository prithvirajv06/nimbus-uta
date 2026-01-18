package models

type HitPolicy string

const (
	First   HitPolicy = "FIRST"   // Execute only the first match
	Collect HitPolicy = "COLLECT" // Execute all matches
	Unique  HitPolicy = "UNIQUE"  // Error if more than one match
)

type DTRuleSet struct {
	InputColumns  []DTColumnMetadata       `json:"columns"`
	OutputColumns []DTColumnMetadata       `json:"output_columns"`
	Rules         []map[string]interface{} `json:"rules"`
	HitPolicy     HitPolicy                `json:"hit_policy"`
}

type DTAction struct {
	Target Variables   `json:"target"`
	Value  interface{} `json:"value"`
}

type DTColumnMetadata struct {
	ID         string       `json:"id"`
	Priority   int          `json:"priority"` // Used for sorting or selection
	Conditions WorkflowStep `json:"conditions"`
	Actions    []DTAction   `json:"actions"`
}

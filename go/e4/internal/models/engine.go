package models

import "encoding/json"

type WorkflowDef struct {
	NIMB_ID  string         `bson:"nimb_id" json:"nimb_id"`
	Engine   string         `json:"engine_name"`
	Pipeline []PipelineStep `json:"pipeline"`
	Metadata []VariableMeta `json:"variable_metadata"`
	Audit    Audit          `json:"audit" bson:"audit"`
}

type VariableMeta struct {
	Key   string      `json:"var_key"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type PipelineStep struct {
	Type string `json:"type"`

	// Condition fields
	Statement string         `json:"statement,omitempty"`
	Children  []PipelineStep `json:"children,omitempty"`

	// Assignment/Array fields
	Target string          `json:"target,omitempty"`
	Value  json.RawMessage `json:"value,omitempty"` // Keep raw to handle "Strings" vs Objects {}
	// Network fields
	URL    string `json:"url,omitempty"`
	Method string `json:"method,omitempty"`
	//Local Variable to store result
	ContextVar string `json:"context_var,omitempty"`
}

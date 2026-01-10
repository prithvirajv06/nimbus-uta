package engine

import (
	"encoding/json"
	"fmt"

	"github.com/prithvirajv06/nimbus-uta/go/engine/internal/models"
)

func TestGenerateScript() {
	// Sample metadata
	metadata := []models.VariableMeta{
		{Key: "metaVar", Value: "42"},
	}

	// Sample pipeline steps
	// Define steps as JSON
	stepsJSON := `
[
	{
		"type": "assignment",
		"target": "amount",
		"value": "100"
	},
	{
		"type": "push_array",
		"target": "items",
		"value": "{\"id\":1}"
	},
	{
		"type": "condition",
		"statement": "amount > 50",
		"children": [
			{
				"type": "assignment",
				"target": "status",
				"value": "approved"
			}
		]
	},
	{
		"type": "for_each",
		"target": "items",
		"context_var": "item",
		"children": [
			{
				"type": "assignment",
				"target": "processed",
				"value": "true"
			}
		]
	}
]
`
	var steps []models.PipelineStep
	if err := json.Unmarshal([]byte(stepsJSON), &steps); err != nil {
		fmt.Printf("Failed to unmarshal steps: %v", err)
		return
	}

	wf := models.WorkflowDef{
		NIMB_ID:  "wf1",
		Metadata: metadata,
		Pipeline: steps,
	}

	script := GenerateScript(wf)
	// Print the generated script in file
	fmt.Println("Generated Script:\n", script)
}

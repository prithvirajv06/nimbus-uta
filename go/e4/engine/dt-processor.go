package engine

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/prithvirajv06/nimbus-uta/go/engine/internal/models"
)

type DTProcessor struct {
	Rules       []models.DTRule
	ArrayColumn string
	Policy      models.HitPolicy
}

func (p *DTProcessor) Process(input map[string]interface{}) error {
	arrayData := GetNestedValue(input, p.ArrayColumn)

	if items, ok := arrayData.([]interface{}); ok {
		for _, item := range items {
			if err := p.evaluate(input, item); err != nil {
				return err
			}
		}
	} else {
		return p.evaluate(input, nil)
	}
	return nil
}

func (p *DTProcessor) evaluate(global map[string]interface{}, local interface{}) error {
	var matches []models.DTRule

	for _, rule := range p.Rules {
		match := true
		for path, expected := range rule.Conditions {
			actual := GetValue(global, local, path, p.ArrayColumn)
			if fmt.Sprintf("%v", actual) != expected {
				match = false
				break
			}
		}

		if match {
			matches = append(matches, rule)
			// Optimization: If policy is FIRST, we can stop looking immediately
			if p.Policy == models.First {
				break
			}
		}
	}

	// Validate UNIQUE policy
	if p.Policy == models.Unique && len(matches) > 1 {
		return fmt.Errorf("hit policy violation: UNIQUE expected 1 match, found %d", len(matches))
	}

	// Apply actions based on matches found
	for _, rule := range matches {
		for _, action := range rule.Actions {
			applyAction(global, local, action, p.ArrayColumn)
		}
	}

	return nil
}

// Helper to apply the action logic
func applyAction(global map[string]interface{}, local interface{}, action models.DTAction, arrayCol string) {
	// 1. Handle Wildcard Broadcast: "cart[*].discount"
	wildcardPrefix := arrayCol + "[*]."
	if strings.HasPrefix(action.Path, wildcardPrefix) {
		subPath := strings.TrimPrefix(action.Path, wildcardPrefix)

		// Get the array from the global context
		if arrayData, ok := GetNestedValue(global, arrayCol).([]interface{}); ok {
			for _, item := range arrayData {
				SetNestedValue(item, subPath, action.Value)
			}
		}
		return
	}

	// 2. Handle Specific Item Update: "cart.discount" (updates current item in loop)
	if local != nil && strings.HasPrefix(action.Path, arrayCol+".") {
		subPath := strings.TrimPrefix(action.Path, arrayCol+".")
		SetNestedValue(local, subPath, action.Value)
		return
	}

	// 3. Handle Global Update: "total_discount"
	SetNestedValue(global, action.Path, action.Value)
}

// ... (GetValue, GetNestedValue, and SetNestedValue functions from previous example remain the same) ...
// SetNestedValue writes a value to a map based on dot notation
func SetNestedValue(data interface{}, path string, value interface{}) {
	parts := strings.Split(path, ".")
	current := data

	for i, part := range parts {
		m, ok := current.(map[string]interface{})
		if !ok {
			return
		}

		if i == len(parts)-1 {
			m[part] = value
		} else {
			// If nested map doesn't exist, create it
			if m[part] == nil {
				m[part] = make(map[string]interface{})
			}
			current = m[part]
		}
	}
}

// GetValue (same logic as before)
func GetValue(globalCtx map[string]interface{}, localItem interface{}, path string, arrayCol string) interface{} {
	if localItem != nil && strings.HasPrefix(path, arrayCol+".") {
		subPath := strings.TrimPrefix(path, arrayCol+".")
		return GetNestedValue(localItem, subPath)
	}
	return GetNestedValue(globalCtx, path)
}

// GetNestedValue retrieves a value from a map using dot notation (e.g., "user.profile.age")
func GetNestedValue(data interface{}, path string) interface{} {
	parts := strings.Split(path, ".")
	var current interface{} = data

	for _, part := range parts {
		if m, ok := current.(map[string]interface{}); ok {
			current = m[part]
		} else {
			return nil
		}
	}
	return current
}

func main() {
	rawJSON := `{
		"user_segment": "gold",
		"cart": [
			{"item": "laptop", "price": 1000}
		]
	}`

	var data map[string]interface{}
	json.Unmarshal([]byte(rawJSON), &data)

	processor := DTProcessor{
		ArrayColumn: "cart",
		Policy:      models.Unique, // Try changing this to First or Collect
		Rules: []models.DTRule{
			{
				ID:         "GOLD_DISCOUNT",
				Conditions: map[string]string{"user_segment": "gold"},
				Actions:    []models.DTAction{{Path: "cart.discount", Value: 0.1}},
			},
		},
	}

	err := processor.Process(data)
	if err != nil {
		fmt.Println("Error:", err)
	}

	output, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(output))
}

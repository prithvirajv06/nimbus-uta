package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"unicode"
)

type Variables struct {
	Variable           string      `json:"variable"`
	SampleValue        string      `json:"sample_value"`
	Description        string      `json:"description"`
	ContextVar         string      `json:"context_var"`
	ContextVarToCreate string      `json:"context_var_to_create"`
	Type               string      `json:"type"`
	Required           bool        `json:"required"`
	Children           []Variables `json:"children,omitempty"`
}

func main() {
	tesJson, err := os.ReadFile("test.json")
	if err != nil {
		panic(err)
	}
	var variable Variables
	processVariable(string(tesJson), "data", &variable)

}

func processVariable(jsonStr string, contextPar string, variable *Variables) ([]Variables, error) {
	var data interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	var variables []Variables
	extractVariablesV2(data, "", &variables)
	return variables, nil
}

func extractVariablesV2(data any, prefix string, variables *[]Variables) {
	switch v := data.(type) {
	case map[string]any:
		for key, value := range v {
			fullKey := key
			contextVarKey := strings.ReplaceAll(prefix, ".", "_")
			contextVarKey = strings.ReplaceAll(contextVarKey, "[*]", "Array")
			if prefix != "" {
				fullKey = prefix + "." + key
			}
			valueType := getType(value)
			switch valueType {
			case "object":
				// Recurse into the object
				var children []models.Variables
				extractVariablesV2(value, fullKey, &children)
				*variables = append(*variables, models.Variables{
					VarKey:        fullKey,
					ContextVarKey: contextVarKey,
					Label:         formatLabel(key),
					Type:          valueType,
					IsRequired:    true,
					Value:         toJSONString(value),
					Children:      children,
				})
			case "array":
				// Check array elements
				var children []models.Variables
				if arr, ok := value.([]any); ok && len(arr) > 0 {
					// Analyze first element to understand array structure
					firstElemType := getType(arr[0])
					if firstElemType == "object" {
						extractVariablesV2(arr[0], fullKey+"[*]", &children)
						*variables = append(*variables, models.Variables{
							VarKey:        fullKey,
							ContextVarKey: contextVarKey,
							Label:         formatLabel(key),
							Type:          valueType,
							IsRequired:    true,
							Value:         toJSONString(value),
							Children:      children,
						})
					}

				}
			default:
				// It's a primitive type (string, number, boolean, null)
				*variables = append(*variables, models.Variables{
					VarKey:        fullKey,
					ContextVarKey: contextVarKey,
					Label:         formatLabel(key),
					Type:          valueType,
					IsRequired:    true,
					Value:         toJSONString(value),
				})
			}
		}
	}
}

func toJSONString(value any) string {
	bytes, err := json.Marshal(value)
	if err != nil {
		return value.(string)
	}
	return string(bytes)
}

func getType(value interface{}) string {
	if value == nil {
		return "null"
	}

	switch value.(type) {
	case string:
		return "string"
	case float64, int, int64, int32, float32:
		return "number"
	case bool:
		return "boolean"
	case []interface{}:
		return "array"
	case map[string]interface{}:
		return "object"
	default:
		return "unknown"
	}
}

func formatLabel(key string) string {
	// Handle empty key
	if key == "" {
		return ""
	}

	// Remove array notation if present
	key = strings.TrimSuffix(key, "[]")

	// Convert camelCase or PascalCase to words
	var result strings.Builder
	var prev rune

	for i, r := range key {
		// Handle underscores and hyphens
		if r == '_' || r == '-' {
			result.WriteRune(' ')
			continue
		}

		// Add space before uppercase letter if:
		// - Not first character
		// - Previous was lowercase or digit
		if i > 0 && unicode.IsUpper(r) && (unicode.IsLower(prev) || unicode.IsDigit(prev)) {
			result.WriteRune(' ')
		}

		// First letter or letter after space should be uppercase
		if i == 0 || prev == ' ' {
			result.WriteRune(unicode.ToUpper(r))
		} else {
			result.WriteRune(r)
		}

		prev = r
	}

	return result.String()
}

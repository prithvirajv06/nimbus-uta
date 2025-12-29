package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode"

	"github.com/prithvirajv06/nimbus-uta/go/core/internal/models"
)

func extractVariablesFromJSON(jsonStr string) ([]models.Variables, error) {
	var data interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	var variables []models.Variables
	extractVariables(data, "", &variables)
	return variables, nil
}

func extractVariables(data interface{}, prefix string, variables *[]models.Variables) {
	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			fullKey := key
			if prefix != "" {
				fullKey = prefix + "." + key
			}

			valueType := getType(value)

			// If it's an object, recurse into it
			if valueType == "object" {
				extractVariables(value, fullKey, variables)
			} else if valueType == "array" {
				// Add the array itself
				*variables = append(*variables, models.Variables{
					VarKey:     fullKey,
					Label:      formatLabel(key),
					Type:       valueType,
					IsRequired: true,
				})

				// Check array elements
				if arr, ok := value.([]interface{}); ok && len(arr) > 0 {
					// Analyze first element to understand array structure
					firstElemType := getType(arr[0])
					if firstElemType == "object" {
						extractVariables(arr[0], fullKey+"[]", variables)
					}
				}
			} else {
				// It's a primitive type (string, number, boolean, null)
				*variables = append(*variables, models.Variables{
					VarKey:     fullKey,
					Label:      formatLabel(key),
					Type:       valueType,
					IsRequired: true,
				})
			}
		}
	}
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

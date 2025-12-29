package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/prithvirajv06/nimbus-uta/go/core/config"
	"github.com/prithvirajv06/nimbus-uta/go/core/internal/models"
	"github.com/prithvirajv06/nimbus-uta/go/core/internal/repository"
	"github.com/prithvirajv06/nimbus-uta/go/core/internal/utils"
	"github.com/prithvirajv06/nimbus-uta/go/core/pkg/database"
	"github.com/prithvirajv06/nimbus-uta/go/core/pkg/messaging"
)

type VariablePackageService struct {
	mongo    *database.MongoDB
	rabbitMQ *messaging.RabbitMQ
	cfg      *config.Config
}

func (s *VariablePackageService) CreateNewVariablePackageFromJSON(c *gin.Context) {
	var payload models.VariablePackageRequet
	err := c.ShouldBindJSON(&payload)
	if HandleError(c, err, "Unable to unmarshel payload") {
		return
	}
	variables, err := extractVariablesFromJSON(payload.JSONStr)
	if HandleError(c, err, "Failed to extract variables from JSON") {
		return
	}
	var newNimId = utils.GenerateNIMBID("VAR_PKG")
	varPackage := models.VariablePackage{
		NIMB_ID:     newNimId,
		PackageName: payload.PackageName,
		Description: payload.Description,
		Variables:   variables,
	}
	varPackage.Audit.Version, _ = GetNextVersionNumber(c, s.mongo.Database, newNimId)
	varPackage.Audit.MinorVersion = 1
	repo := repository.NewGenericRepository[models.VariablePackage](c.Request.Context(), s.mongo.Database, "variable_packages")
	_, err = repo.InsertOne(varPackage)
	if HandleError(c, err, "Failed to create variable package") {
		return
	}
	RespondJSON(c, 201, "success", "Variable pacakage created", varPackage)
}

func (s *VariablePackageService) UpdateVariablePackage(c *gin.Context) {
	var payload models.VariablePackage
	err := c.ShouldBindJSON(&payload)
	if HandleError(c, err, "Unable to unmarshel payload") {
		return
	}
	// Archive Old Version and Create New Version
	repo := repository.NewGenericRepository[models.VariablePackage](c.Request.Context(), s.mongo.Database, "variable_packages")
	err = ArchiveEntity(repo, payload.NIMB_ID)
	if HandleError(c, err, "Failed to archive old version of variable package") {
		return
	}
	payload.Audit.MinorVersion += 1
	_, err = repo.InsertOne(payload)
	if HandleError(c, err, "Failed to create new version of variable package") {
		return
	}
	RespondJSON(c, 200, "success", "Variable package updated successfully", payload)
}

func (s *VariablePackageService) GetVariablePackageByNIMBID(c *gin.Context) {
	nimbID := c.Param("nimb_id")
	repo := repository.NewGenericRepository[models.VariablePackage](c.Request.Context(), s.mongo.Database, "variable_packages")
	varPackage, err := repo.FindOne(map[string]interface{}{"nimb_id": nimbID, "audit.is_archived": false})
	if HandleError(c, err, "Failed to retrieve variable package") {
		return
	}
	RespondJSON(c, 200, "success", "Variable package retrieved successfully", varPackage)
}

func (s *VariablePackageService) ArchiveVariablePackage(c *gin.Context) {
	repo := repository.NewGenericRepository[models.VariablePackage](c.Request.Context(), s.mongo.Database, "variable_packages")
	err := ArchiveEntity(repo, c.Param("nimb_id"))
	if HandleError(c, err, "Failed to archive variable package") {
		return
	}
	RespondJSON(c, 200, "success", "Variable package archived successfully", nil)
}

func (s *VariablePackageService) ListAllVariablePackages(c *gin.Context) {
	var filter models.VariablePackage
	err := c.ShouldBindQuery(&filter)
	if HandleError(c, err, "Invalid query parameters") {
		return
	}
	repo := repository.NewGenericRepository[models.VariablePackage](c.Request.Context(), s.mongo.Database, "variable_packages")
	varPackages, err := repo.FindMany(filter)
	HandleError(c, err, "Failed to retrieve variable packages")
	RespondJSON(c, 200, "success", "Variable packages retrieved successfully", varPackages)
}

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

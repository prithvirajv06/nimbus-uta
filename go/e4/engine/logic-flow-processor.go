package engine

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/Knetic/govaluate"
	"github.com/gorilla/websocket"
	"github.com/prithvirajv06/nimbus-uta/go/engine/internal/models"
)

// RuleEngine handles rule execution with context and operations
type RuleEngine struct {
	GlobalVars map[string]interface{}
	Input      map[string]interface{}
	Context    map[string]interface{} // Current loop context
	DB         *sql.DB
	HTTPClient *http.Client
	WSConn     *websocket.Conn
	Log        []string
}

func (re *RuleEngine) GetLog() []string {
	return re.Log
}

func (re *RuleEngine) AddLog(entry string) {
	re.Log = append(re.Log, entry)
}

// NewRuleEngine creates a new rule engine instance
func NewRuleEngine(input map[string]interface{}) *RuleEngine {
	return &RuleEngine{
		GlobalVars: make(map[string]interface{}),
		Input:      input,
		Context:    make(map[string]interface{}),
		HTTPClient: &http.Client{},
	}
}

// ExecuteWorkflow executes the entire workflow
func (re *RuleEngine) ExecuteWorkflow(config models.LogicFlow) error {
	re.AddLog("Starting workflow execution")
	for _, step := range config.LogicalSteps {
		re.AddLog(fmt.Sprintf("Executing step: %s", step.Type))
		if err := re.executeStep(step); err != nil {
			re.AddLog(fmt.Sprintf("Error executing step %s: %v", step.Type, err))
			return err
		}
	}
	re.AddLog("Workflow execution completed")
	return nil
}

// executeStep executes a single logical step
func (re *RuleEngine) executeStep(step models.WorkflowStep) error {
	re.AddLog(fmt.Sprintf("Processing step type: %s", step.Type))
	switch step.Type {
	case "start":
		return re.executeChildren(step.Children)
	case "for_each":
		return re.executeForEach(step)
	case "condition":
		return re.executeCondition(step)
	case "assignment":
		return re.executeAssignment(step)
	default:
		re.AddLog(fmt.Sprintf("Unknown step type: %s", step.Type))
		return fmt.Errorf("unknown step type: %s", step.Type)
	}
}

// executeForEach handles for_each loops
func (re *RuleEngine) executeForEach(step models.WorkflowStep) error {
	// Resolve the array path
	arrayPath := step.Target.VarKey
	re.AddLog(fmt.Sprintf("Executing for_each on array path: %s", arrayPath))
	data, err := re.GetValue(arrayPath)
	if err != nil {
		re.AddLog(fmt.Sprintf("Failed to get array at %s: %v", arrayPath, err))
		return fmt.Errorf("failed to get array at %s: %v", arrayPath, err)
	}
	slice := reflect.ValueOf(data)
	if slice.Kind() != reflect.Slice {
		re.AddLog(fmt.Sprintf("for_each target is not a slice: %s", arrayPath))
		return fmt.Errorf("for_each target is not a slice: %s", arrayPath)
	}
	originalContext := make(map[string]interface{})
	for k, v := range re.Context {
		originalContext[k] = v
	}
	for i := 0; i < slice.Len(); i++ {
		item := slice.Index(i).Interface()
		if step.ContextVar != "" {
			re.Context[step.ContextVar] = item
			re.Context[step.ContextVar+"_index"] = i
		}
		re.AddLog(fmt.Sprintf("for_each iteration %d, contextVar: %s", i, step.ContextVar))
		for _, child := range step.Children {
			if err := re.executeStep(child); err != nil {
				re.AddLog(fmt.Sprintf("Error in for_each child step: %v", err))
				return err
			}
		}
	}
	re.Context = originalContext
	re.AddLog("for_each completed")
	return nil
}

// executeCondition handles conditional logic
func (re *RuleEngine) executeCondition(step models.WorkflowStep) error {
	if step.Statement == "" {
		re.AddLog("Condition statement is empty")
		return fmt.Errorf("condition statement is empty")
	}
	if re.isExpression(step.Statement) {
		for _, cond := range step.ConditionConfig {
			re.AddLog(fmt.Sprintf("Condition config: LeftVar=%v, Operator=%s, RightValue=%v", cond.LeftVar, cond.Operator, cond.RightValue))
			processConditionConfig(re, cond)
		}
	}
	re.AddLog(fmt.Sprintf("Evaluating condition: %s", step.Statement))
	result, err := re.evaluateExpression(step.Statement, step.ContextMap)
	if err != nil {
		re.AddLog(fmt.Sprintf("Failed to evaluate condition '%s': %v", step.Statement, err))
		return fmt.Errorf("failed to evaluate condition '%s': %v", step.Statement, err)
	}
	boolResult := false
	switch v := result.(type) {
	case bool:
		boolResult = v
	case float64:
		boolResult = v != 0
	case string:
		boolResult = v != ""
	default:
		boolResult = v != nil
	}
	re.AddLog(fmt.Sprintf("Condition result: %v", boolResult))
	if boolResult {
		re.AddLog("Executing TrueChildren for condition")
		return re.executeChildren(step.TrueChildren)
	} else {
		re.AddLog("Executing FalseChildren for condition")
		return re.executeChildren(step.FalseChildren)
	}
}

// executeAssignment handles value assignments
func (re *RuleEngine) executeAssignment(step models.WorkflowStep) error {
	targetPath := step.Target.VarKey
	re.AddLog(fmt.Sprintf("Assignment step: targetPath=%s", targetPath))
	var value interface{}
	var err error
	if valueStr, ok := step.Value.(string); ok {
		if re.isExpression(valueStr) {
			re.AddLog(fmt.Sprintf("Evaluating assignment expression: %s", valueStr))
			value, err = re.evaluateExpression(valueStr, step.ContextMap)
			if err != nil {
				re.AddLog(fmt.Sprintf("Failed to evaluate assignment expression '%s': %v", valueStr, err))
				return fmt.Errorf("failed to evaluate expression '%s': %v", valueStr, err)
			}
		} else {
			value = valueStr
		}
	} else {
		value = step.Value
	}
	targetPath = re.resolveContextPath(targetPath, step.ContextMap)
	re.AddLog(fmt.Sprintf("Setting value at path: %s", targetPath))
	//set value with type
	if step.Target.Type != "" {
		switch step.Target.Type {
		case "string":
			value = fmt.Sprintf("%v", value)
		case "number":
			switch v := value.(type) {
			case float64:
				value = int(v)
			case string:
				intVal, err := strconv.Atoi(v)
				if err != nil {
					re.AddLog(fmt.Sprintf("Failed to convert value to int: %v", err))
					return fmt.Errorf("failed to convert value to int: %v", err)
				}
				value = intVal
			}
		case "float":
			switch v := value.(type) {
			case int:
				value = float64(v)
			case string:
				floatVal, err := strconv.ParseFloat(v, 64)
				if err != nil {
					re.AddLog(fmt.Sprintf("Failed to convert value to float: %v", err))
					return fmt.Errorf("failed to convert value to float: %v", err)
				}
				value = floatVal
			}
		case "boolean":
			switch v := value.(type) {
			case string:
				boolVal, err := strconv.ParseBool(v)
				if err != nil {
					re.AddLog(fmt.Sprintf("Failed to convert value to bool: %v", err))
					return fmt.Errorf("failed to convert value to bool: %v", err)
				}
				value = boolVal
			}
		}
	}
	err = re.SetValue(targetPath, value)
	if err != nil {
		re.AddLog(fmt.Sprintf("Error setting value at %s: %v", targetPath, err))
	} else {
		re.AddLog(fmt.Sprintf("Value set at %s", targetPath))
	}
	return err
}

// executeChildren executes a list of child steps
func (re *RuleEngine) executeChildren(children []models.WorkflowStep) error {
	re.AddLog(fmt.Sprintf("Executing %d child steps", len(children)))
	for _, child := range children {
		if err := re.executeStep(child); err != nil {
			re.AddLog(fmt.Sprintf("Error executing child step: %v", err))
			return err
		}
	}
	re.AddLog("All child steps executed")
	return nil
}

// evaluateExpression evaluates a string expression
func (re *RuleEngine) evaluateExpression(expr string, contextMap []models.LoopContextMap) (interface{}, error) {
	// Replace array wildcards with context values
	re.AddLog(fmt.Sprintf("Evaluating expression: %s", expr))
	resolvedExpr := re.resolveContextExpression(expr, contextMap)
	parameters := make(map[string]interface{})
	varMapping := make(map[string]string)
	re.extractAndResolveVars(resolvedExpr, parameters, varMapping)
	sanitizedExpr := resolvedExpr
	for original, sanitized := range varMapping {
		sanitizedExpr = strings.ReplaceAll(sanitizedExpr, original, sanitized)
	}
	expression, err := govaluate.NewEvaluableExpression(sanitizedExpr)
	if err != nil {
		re.AddLog(fmt.Sprintf("Expression parse error: %v (expr: %s)", err, sanitizedExpr))
		return nil, fmt.Errorf("expression parse error: %v (expr: %s)", err, sanitizedExpr)
	}
	result, err := expression.Evaluate(parameters)
	if err != nil {
		re.AddLog(fmt.Sprintf("Evaluation error: %v (params: %+v)", err, parameters))
		return nil, fmt.Errorf("evaluation error: %v (params: %+v)", err, parameters)
	}
	re.AddLog(fmt.Sprintf("Expression result: %v", result))
	return result, nil
}

// resolveContextExpression replaces wildcards with actual context indices
func (re *RuleEngine) resolveContextExpression(expr string, contextMap []models.LoopContextMap) string {
	if len(contextMap) == 0 {
		return expr
	}

	resolved := expr
	for _, cm := range contextMap {
		// Get the index from context
		indexKey := cm.ContextKey + "_index"
		if idx, ok := re.Context[indexKey]; ok {
			// Replace data.applicants[*] with data.applicants[0]
			wildcardPattern := regexp.MustCompile(regexp.QuoteMeta(cm.VarKey))
			replacement := strings.Replace(cm.VarKey, "[*]", fmt.Sprintf("[%d]", idx), -1)
			resolved = wildcardPattern.ReplaceAllString(resolved, replacement)
		}
	}

	return resolved
}

// resolveContextPath resolves array wildcards in paths
func (re *RuleEngine) resolveContextPath(path string, contextMap []models.LoopContextMap) string {
	if len(contextMap) == 0 {
		return path
	}

	resolved := path
	for _, cm := range contextMap {
		indexKey := cm.ContextKey + "_index"
		if idx, ok := re.Context[indexKey]; ok {
			wildcardPattern := regexp.MustCompile(regexp.QuoteMeta(cm.VarKey))
			replacement := strings.Replace(cm.VarKey, "[*]", fmt.Sprintf("[%d]", idx), -1)
			resolved = wildcardPattern.ReplaceAllString(resolved, replacement)
		}
	}

	return resolved
}

// isExpression checks if a string is an expression
func (re *RuleEngine) isExpression(s string) bool {
	// Contains operators or function calls
	return strings.Contains(s, "+") || strings.Contains(s, "-") ||
		strings.Contains(s, "*") || strings.Contains(s, "/") ||
		strings.Contains(s, ">") || strings.Contains(s, "<") ||
		strings.Contains(s, "==") || strings.Contains(s, "!=") ||
		strings.Contains(s, "?") || strings.Contains(s, ".") && !strings.HasPrefix(s, "\"")
}

// extractAndResolveVars extracts variables from expression and resolves their values
func (re *RuleEngine) extractAndResolveVars(expr string, parameters map[string]interface{}, varMapping map[string]string) {
	// Find all variable patterns like data.applicants[0].creditScore
	varPattern := regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*(\.[a-zA-Z_][a-zA-Z0-9_]*|\[[0-9]+\])*`)
	matches := varPattern.FindAllString(expr, -1)

	for _, match := range matches {
		// Skip operators and keywords
		if match == "true" || match == "false" || match == "null" {
			continue
		}

		// Get the actual value
		value, err := re.GetValue(match)
		if err != nil {
			// Variable doesn't exist, skip
			continue
		}

		// Create sanitized variable name (replace dots and brackets)
		sanitized := strings.ReplaceAll(match, ".", "_")
		sanitized = strings.ReplaceAll(sanitized, "[", "_")
		sanitized = strings.ReplaceAll(sanitized, "]", "_")

		parameters[sanitized] = value
		varMapping[match] = sanitized
	}
}

// flattenMap flattens nested maps into dot notation
func (re *RuleEngine) flattenMap(prefix string, data interface{}, result map[string]interface{}) {
	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			newPrefix := key
			if prefix != "" {
				newPrefix = prefix + "." + key
			}
			re.flattenMap(newPrefix, value, result)
		}
	case []interface{}:
		for i, item := range v {
			newPrefix := fmt.Sprintf("%s[%d]", prefix, i)
			re.flattenMap(newPrefix, item, result)
		}
	default:
		if prefix != "" {
			result[prefix] = v
		}
	}
}

// GetValue retrieves value from input using dot notation
func (re *RuleEngine) GetValue(path string) (interface{}, error) {
	re.AddLog(fmt.Sprintf("GetValue called for path: %s", path))
	if strings.HasPrefix(path, "\"") && strings.HasSuffix(path, "\"") {
		return strings.Trim(path, "\""), nil
	}
	if strings.HasPrefix(path, "'") && strings.HasSuffix(path, "'") {
		return strings.Trim(path, "'"), nil
	}
	parts := re.parsePath(path)
	var current interface{} = re.Input
	for i, part := range parts {
		if part.isArray {
			if i == 0 || !parts[i-1].isArray {
				switch v := current.(type) {
				case map[string]interface{}:
					current = v[part.name]
					if current == nil {
						re.AddLog(fmt.Sprintf("Path not found: %s", part.name))
						return nil, fmt.Errorf("path not found: %s", part.name)
					}
				default:
					re.AddLog(fmt.Sprintf("Cannot access property %s on non-object", part.name))
					return nil, fmt.Errorf("cannot access property %s on non-object", part.name)
				}
			}
			slice := reflect.ValueOf(current)
			if slice.Kind() != reflect.Slice {
				re.AddLog(fmt.Sprintf("Not a slice: %s (got %T)", part.name, current))
				return nil, fmt.Errorf("not a slice: %s (got %T)", part.name, current)
			}
			if part.index >= slice.Len() || part.index < 0 {
				re.AddLog(fmt.Sprintf("Index out of bounds: %d (len: %d)", part.index, slice.Len()))
				return nil, fmt.Errorf("index out of bounds: %d (len: %d)", part.index, slice.Len())
			}
			current = slice.Index(part.index).Interface()
		} else {
			switch v := current.(type) {
			case map[string]interface{}:
				val, ok := v[part.name]
				if !ok {
					re.AddLog(fmt.Sprintf("Property not found: %s", part.name))
					return nil, fmt.Errorf("property not found: %s", part.name)
				}
				current = val
			default:
				re.AddLog(fmt.Sprintf("Cannot traverse path at %s: not an object (got %T)", part.name, current))
				return nil, fmt.Errorf("cannot traverse path at %s: not an object (got %T)", part.name, current)
			}
		}
	}
	re.AddLog(fmt.Sprintf("GetValue result for path %s: %v", path, current))
	return current, nil
}

// SetValue sets value in input using dot notation
func (re *RuleEngine) SetValue(path string, value interface{}) error {
	re.AddLog(fmt.Sprintf("SetValue called for path: %s, value: %v", path, value))
	parts := re.parsePath(path)
	var current interface{} = re.Input
	for i := 0; i < len(parts)-1; i++ {
		part := parts[i]
		if part.isArray {
			switch v := current.(type) {
			case map[string]interface{}:
				current = v[part.name]
			}
			slice := reflect.ValueOf(current)
			if slice.Kind() != reflect.Slice {
				re.AddLog(fmt.Sprintf("Not a slice: %s", part.name))
				return fmt.Errorf("not a slice: %s", part.name)
			}
			current = slice.Index(part.index).Interface()
		} else {
			switch v := current.(type) {
			case map[string]interface{}:
				if v[part.name] == nil {
					v[part.name] = make(map[string]interface{})
				}
				current = v[part.name]
			default:
				re.AddLog(fmt.Sprintf("Cannot traverse path: %s", path))
				return fmt.Errorf("cannot traverse path: %s", path)
			}
		}
	}
	lastPart := parts[len(parts)-1]
	if lastPart.isArray {
		if m, ok := current.(map[string]interface{}); ok {
			arr := m[lastPart.name]
			slice := reflect.ValueOf(arr)
			if slice.Kind() == reflect.Slice {
				sliceVal := arr.([]interface{})
				if lastPart.index < len(sliceVal) {
					if _, ok := sliceVal[lastPart.index].(map[string]interface{}); ok {
						re.AddLog("Path incomplete for array element")
						return fmt.Errorf("path incomplete for array element")
					}
					sliceVal[lastPart.index] = value
					re.AddLog(fmt.Sprintf("Set value in array at %s[%d]", lastPart.name, lastPart.index))
				}
			}
		}
	} else {
		if m, ok := current.(map[string]interface{}); ok {
			m[lastPart.name] = value
			re.AddLog(fmt.Sprintf("Set value for key %s", lastPart.name))
			return nil
		}
	}
	return nil
}

// PathPart represents a part of a path
type PathPart struct {
	name    string
	isArray bool
	index   int
}

// parsePath parses a path string into parts
func (re *RuleEngine) parsePath(path string) []PathPart {
	parts := []PathPart{}
	segments := strings.Split(path, ".")

	for _, seg := range segments {
		if strings.Contains(seg, "[") && strings.Contains(seg, "]") {
			name := seg[:strings.Index(seg, "[")]
			indexStr := seg[strings.Index(seg, "[")+1 : strings.Index(seg, "]")]
			index, _ := strconv.Atoi(indexStr)
			parts = append(parts, PathPart{name: name, isArray: true, index: index})
		} else {
			parts = append(parts, PathPart{name: seg, isArray: false})
		}
	}

	return parts
}

// HTTPGet performs HTTP GET
func (re *RuleEngine) HTTPGet(url string) (map[string]interface{}, error) {
	re.AddLog(fmt.Sprintf("HTTP GET request to URL: %s", url))
	resp, err := re.HTTPClient.Get(url)
	if err != nil {
		re.AddLog(fmt.Sprintf("HTTP GET error: %v", err))
		return nil, err
	}
	defer resp.Body.Close()
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		re.AddLog(fmt.Sprintf("HTTP GET decode error: %v", err))
		return nil, err
	}
	re.AddLog(fmt.Sprintf("HTTP GET response: %v", result))
	return result, nil
}

// HTTPPost performs HTTP POST
func (re *RuleEngine) HTTPPost(url string, data interface{}) (map[string]interface{}, error) {
	re.AddLog(fmt.Sprintf("HTTP POST request to URL: %s, data: %v", url, data))
	jsonData, err := json.Marshal(data)
	if err != nil {
		re.AddLog(fmt.Sprintf("HTTP POST marshal error: %v", err))
		return nil, err
	}
	resp, err := re.HTTPClient.Post(url, "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		re.AddLog(fmt.Sprintf("HTTP POST error: %v", err))
		return nil, err
	}
	defer resp.Body.Close()
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		re.AddLog(fmt.Sprintf("HTTP POST decode error: %v", err))
		return nil, err
	}
	re.AddLog(fmt.Sprintf("HTTP POST response: %v", result))
	return result, nil
}

// QueryDB executes database query
func (re *RuleEngine) QueryDB(query string, args ...interface{}) ([]map[string]interface{}, error) {
	re.AddLog(fmt.Sprintf("QueryDB called: %s, args: %v", query, args))
	if re.DB == nil {
		re.AddLog("Database not connected")
		return nil, fmt.Errorf("database not connected")
	}
	rows, err := re.DB.Query(query, args...)
	if err != nil {
		re.AddLog(fmt.Sprintf("QueryDB error: %v", err))
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		re.AddLog(fmt.Sprintf("QueryDB columns error: %v", err))
		return nil, err
	}
	results := []map[string]interface{}{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}
		if err := rows.Scan(valuePtrs...); err != nil {
			re.AddLog(fmt.Sprintf("QueryDB scan error: %v", err))
			return nil, err
		}
		row := make(map[string]interface{})
		for i, col := range columns {
			row[col] = values[i]
		}
		results = append(results, row)
	}
	re.AddLog(fmt.Sprintf("QueryDB results: %v", results))
	return results, nil
}

// SendWebSocketMessage sends WebSocket message
func (re *RuleEngine) SendWebSocketMessage(message interface{}) error {
	re.AddLog(fmt.Sprintf("SendWebSocketMessage called: %v", message))
	if re.WSConn == nil {
		re.AddLog("WebSocket not connected")
		return fmt.Errorf("websocket not connected")
	}
	err := re.WSConn.WriteJSON(message)
	if err != nil {
		re.AddLog(fmt.Sprintf("WebSocket send error: %v", err))
	} else {
		re.AddLog("WebSocket message sent successfully")
	}
	return err
}

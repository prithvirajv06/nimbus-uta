package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/prithvirajv06/nimbus-uta/go/engine/internal/models"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// --- Engine ---

type Engine struct {
	regexCache sync.Map // Thread-safe cache for Regex compilation
}

// --- LogStack Helper ---

func addInfoLog(logStack *[]models.LogStackEntry, message string) {
	*logStack = append(*logStack, models.LogStackEntry{
		Timestamp: time.Now().UTC(),
		Type:      "INFO",
		Message:   message,
	})
}
func addErrorLog(logStack *[]models.LogStackEntry, message string) {
	*logStack = append(*logStack, models.LogStackEntry{
		Timestamp: time.Now().UTC(),
		Type:      "ERROR",
		Message:   message,
	})
}

func addResultLog(logStack *[]models.LogStackEntry, message string) {
	*logStack = append(*logStack, models.LogStackEntry{
		Timestamp: time.Now().UTC(),
		Type:      "RESULT",
		Message:   message,
	})
}

func addLogs(logStack *[]models.LogStackEntry, logs []models.LogStackEntry) {
	*logStack = append(*logStack, logs...)
}

func NewDTEngine() *Engine {
	return &Engine{}
}

// --- Decision Table Execution ---

func (e *Engine) ProcessDecisionTable(ctx context.Context, table models.DecisionTable, data []byte) ([]byte, []models.LogStackEntry, error) {
	logStack := []models.LogStackEntry{}

	// Helper: recursively process arrays in data for all input columns with [*]
	var processRecursive func(ctx context.Context, table models.DecisionTable, data []byte, prefix string) ([]byte, error)
	processRecursive = func(ctx context.Context, table models.DecisionTable, data []byte, prefix string) ([]byte, error) {
		// Find all array input columns at this level
		arrayInputs := []string{}
		for _, inputDef := range table.InputsColumns {
			if strings.Contains(inputDef.VarKey, "[*]") {
				arrKey := inputDef.VarKey
				if prefix != "" && strings.HasPrefix(arrKey, prefix+".") {
					arrKey = arrKey[len(prefix)+1:]
				}
				arrKey = strings.Split(arrKey, "[*]")[0]
				if arrKey != "" && !contains(arrayInputs, arrKey) {
					arrayInputs = append(arrayInputs, arrKey)
				}
			}
		}
		if len(arrayInputs) == 0 {
			// No array input at this level, process as single object
			var matchedRows [][]models.Variables
			for _, ruleRow := range table.Rules {
				if ctx.Err() != nil {
					addErrorLog(&logStack, "Context cancelled during decision table processing")
					return nil, ctx.Err()
				}
				isMatch := true
				for i, inputDef := range table.InputsColumns {
					addInfoLog(&logStack, fmt.Sprintf("Evaluating Input Column: %+v", inputDef))
					actualVal := gjson.GetBytes(data, inputDef.VarKey)
					if cond, logs := e.evaluateCell(ruleRow[i].Value.(string), actualVal, &logStack); !cond {
						addLogs(&logStack, logs)
						isMatch = false
						break
					}
				}
				if isMatch {
					addInfoLog(&logStack, fmt.Sprintf("Rule Matched: %v", ruleRow))
					matchedRows = append(matchedRows, ruleRow)
					if table.HitPolicy == "FIRST" {
						addInfoLog(&logStack, "Hit Policy FIRST - stopping after first match")
						break
					}
				}
			}
			finalValues, logs, err := e.applyHitPolicy(table, matchedRows, &logStack)
			addLogs(&logStack, logs)
			if err != nil {
				addErrorLog(&logStack, fmt.Sprintf("Error applying hit policy: %v", err))
				return nil, err
			}
			// Set output values
			output := data
			for variable, value := range finalValues {
				output, err = sjson.SetBytes(output, variable, value)
				if err != nil {
					return nil, err
				}
			}
			return output, nil
		}

		// If there are array inputs, process each array recursively
		output := data
		for _, arrKey := range arrayInputs {
			arr := gjson.GetBytes(output, arrKey)
			if arr.IsArray() {
				var newArr []interface{}
				for _, item := range arr.Array() {
					itemBytes := []byte(item.Raw)
					// Recursively process nested arrays in this item
					processedItem, err := processRecursive(ctx, table, itemBytes, arrKey)
					if err != nil {
						return nil, err
					}
					// Now process DT for this object (with all nested arrays already processed)
					var matchedRows [][]models.Variables
					for _, ruleRow := range table.Rules {
						if ctx.Err() != nil {
							addErrorLog(&logStack, "Context cancelled during decision table processing")
							return nil, ctx.Err()
						}
						isMatch := true
						for i, inputDef := range table.InputsColumns {
							// Only match columns that are for this array level
							if strings.HasPrefix(inputDef.VarKey, arrKey+"[*].") {
								subKey := strings.TrimPrefix(inputDef.VarKey, arrKey+"[*].")
								actualVal := gjson.GetBytes(processedItem, subKey)
								addInfoLog(&logStack, fmt.Sprintf("Evaluating Input Column: %+v", inputDef))
								if cond, logs := e.evaluateCell(ruleRow[i].Value.(string), actualVal, &logStack); !cond {
									addLogs(&logStack, logs)
									isMatch = false
									break
								}
							}
						}
						if isMatch {
							addInfoLog(&logStack, fmt.Sprintf("Rule Matched: %v", ruleRow))
							matchedRows = append(matchedRows, ruleRow)
							if table.HitPolicy == "FIRST" {
								addInfoLog(&logStack, "Hit Policy FIRST - stopping after first match")
								break
							}
						}
					}
					finalValues, logs, err := e.applyHitPolicy(table, matchedRows, &logStack)
					addLogs(&logStack, logs)
					if err != nil {
						addErrorLog(&logStack, fmt.Sprintf("Error applying hit policy: %v", err))
						return nil, err
					}
					// Set values for this object
					var obj map[string]interface{}
					if err := json.Unmarshal([]byte(item.Raw), &obj); err != nil {
						addErrorLog(&logStack, fmt.Sprintf("Failed to unmarshal array item: %v", err))
						return nil, err
					}
					for variable, value := range finalValues {
						if strings.HasPrefix(variable, arrKey+"[*].") {
							obj[strings.TrimPrefix(variable, arrKey+"[*].")] = value
						}
					}
					newArr = append(newArr, obj)
				}
				// Set the processed array back to the output
				var err error
				output, err = sjson.SetBytes(output, arrKey, newArr)
				if err != nil {
					return nil, err
				}
			}
		}
		return output, nil
	}

	output, err := processRecursive(ctx, table, data, "")
	return output, logStack, err
}

// Helper: check if a string is in a slice
func contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

func (e *Engine) applyHitPolicy(table models.DecisionTable, matchedRows [][]models.Variables, logStack *[]models.LogStackEntry) (map[string]interface{}, []models.LogStackEntry, error) {
	results := make(map[string]interface{})
	localLogs := []models.LogStackEntry{}
	addInfoLog(logStack, fmt.Sprintf("Applying Hit Policy: %s", table.HitPolicy))
	if len(matchedRows) == 0 {
		addInfoLog(logStack, "No matching rules found")
		return results, localLogs, nil
	}

	policy := strings.ToUpper(table.HitPolicy)
	addInfoLog(logStack, fmt.Sprintf("Total Matched Rows: %d", len(matchedRows)))

	switch policy {
	case "FIRST":
		addInfoLog(logStack, "Hit Policy FIRST - returning first matched row")
		firstRowVals, logs := e.extractRowValues(table, matchedRows[0], logStack)
		addLogs(logStack, logs)
		return firstRowVals, localLogs, nil

	case "ANY":
		addInfoLog(logStack, "Hit Policy ANY - validating all matched rows are identical")
		firstRowVals, logs := e.extractRowValues(table, matchedRows[0], logStack)
		addLogs(logStack, logs)
		for i := 1; i < len(matchedRows); i++ {
			addInfoLog(logStack, fmt.Sprintf("Comparing matched row %d to first row", i+1))
			rowVals, logs := e.extractRowValues(table, matchedRows[i], logStack)
			addLogs(logStack, logs)
			addInfoLog(logStack, fmt.Sprintf("First Row Values: %v, Current Row Values: %v", firstRowVals, rowVals))
			if !isEqual(firstRowVals, rowVals) {
				addErrorLog(logStack, "Hit Policy ANY violated: conflicting outputs found")
				return nil, localLogs, fmt.Errorf("hit policy ANY violated: conflicting outputs found")
			}
		}
		return firstRowVals, localLogs, nil

	case "PRIORITY":
		addInfoLog(logStack, "Hit Policy PRIORITY - selecting highest priority row")
		bestRow := matchedRows[0]
		for i := 1; i < len(matchedRows); i++ {
			addInfoLog(logStack, fmt.Sprintf("Comparing matched row %d to current best row", i+1))
			if e.isHigherPriority(table, matchedRows[i], bestRow) {
				addInfoLog(logStack, "Found new best priority row")
				bestRow = matchedRows[i]
			}
		}
		bestRowVals, logs := e.extractRowValues(table, bestRow, logStack)
		addLogs(logStack, logs)
		return bestRowVals, localLogs, nil
	}

	aggregated := make(map[string][]interface{})
	if policy == "OUTPUT ORDER" {
		sort.SliceStable(matchedRows, func(i, j int) bool {
			return e.isHigherPriority(table, matchedRows[i], matchedRows[j])
		})
	}

	for _, row := range matchedRows {
		vals, logs := e.extractRowValues(table, row, logStack)
		addLogs(logStack, logs)
		for k, v := range vals {
			aggregated[k] = append(aggregated[k], v)
		}
	}

	for _, outDef := range table.OutputsColumns {
		list, exists := aggregated[outDef.VarKey]
		if !exists {
			continue
		}
		switch policy {
		case "SUM":
			sum := 0.0
			for _, v := range list {
				sum += toFloat(v)
			}
			results[outDef.VarKey] = sum
		case "MIN":
			minVal := math.MaxFloat64
			for _, v := range list {
				f := toFloat(v)
				if f < minVal {
					minVal = f
				}
			}
			results[outDef.VarKey] = minVal
		case "MAX":
			maxVal := -math.MaxFloat64
			for _, v := range list {
				f := toFloat(v)
				if f > maxVal {
					maxVal = f
				}
			}
			results[outDef.VarKey] = maxVal
		case "COUNT":
			results[outDef.VarKey] = len(list)
		case "COLLECT", "ALL", "RULE ORDER", "OUTPUT ORDER":
			results[outDef.VarKey] = list
		default:
			return nil, localLogs, fmt.Errorf("unsupported hit policy: %s", policy)
		}
	}
	return results, localLogs, nil
}

// extractRowValues converts a raw rule row into a map of variable->value
func (e *Engine) extractRowValues(table models.DecisionTable, row []models.Variables, logStack *[]models.LogStackEntry) (map[string]interface{}, []models.LogStackEntry) {
	localLogs := []models.LogStackEntry{}
	addInfoLog(logStack, "Extracting row values")
	res := make(map[string]interface{})
	for i, outDef := range table.OutputsColumns {
		addInfoLog(logStack, fmt.Sprintf("Processing output column: %s", outDef.VarKey))
		valIndex := len(table.InputsColumns) + i
		if valIndex < len(row) {
			valStr := row[valIndex].Value.(string)
			if num, err := strconv.ParseFloat(valStr, 64); err == nil {
				res[outDef.VarKey] = num
			} else {
				if strings.ToLower(valStr) == "true" {
					res[outDef.VarKey] = true
				} else if strings.ToLower(valStr) == "false" {
					res[outDef.VarKey] = false
				} else {
					res[outDef.VarKey] = valStr
				}
			}
		}
	}
	addInfoLog(logStack, fmt.Sprintf("Extracted Values: %v", res))
	return res, localLogs
}

// isHigherPriority checks if rowA has a higher priority than rowB.
// Returns true if rowA > rowB.
// isHigherPriority finds the "winning" row based on the explicit priority column
func (e *Engine) isHigherPriority(table models.DecisionTable, rowA, rowB []models.Variables) bool {
	for i, outDef := range table.OutputsColumns {
		if outDef.IsRequired {
			idx := len(table.InputsColumns) + i
			pA, _ := strconv.ParseFloat(rowA[idx].Value.(string), 64)
			pB, _ := strconv.ParseFloat(rowB[idx].Value.(string), 64)
			if pA != pB {
				return pA > pB
			}
		}
	}
	return false
}

func indexOf(slice []string, val string) int {
	for i, item := range slice {
		if item == val {
			return i
		}
	}
	return -1
}

func isEqual(m1, m2 map[string]interface{}) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k, v1 := range m1 {
		if v2, ok := m2[k]; !ok || fmt.Sprintf("%v", v1) != fmt.Sprintf("%v", v2) {
			return false
		}
	}
	return true
}

// applyAggregation handles COLLECT, COLLECT_SUM, COLLECT_COUNT
func (e *Engine) applyAggregation(data []byte, path string, valStr string, policy string) ([]byte, error) {
	current := gjson.GetBytes(data, path)
	var val interface{} = valStr
	if num, err := strconv.ParseFloat(valStr, 64); err == nil {
		val = num
	}
	switch policy {
	case "COLLECT_SUM":
		numVal, ok := val.(float64)
		if !ok {
			return nil, fmt.Errorf("cannot sum non-numeric value: %s", valStr)
		}
		existing := 0.0
		if current.Exists() {
			existing = current.Float()
		}
		return sjson.SetBytes(data, path, existing+numVal)
	case "COLLECT_COUNT":
		existing := 0.0
		if current.Exists() {
			existing = current.Float()
		}
		return sjson.SetBytes(data, path, existing+1)
	default:
		if current.IsArray() {
			var arr []interface{}
			if v, ok := current.Value().([]interface{}); ok {
				arr = v
			}
			arr = append(arr, val)
			return sjson.SetBytes(data, path, arr)
		}
		return sjson.SetBytes(data, path, []interface{}{val})
	}
}

// evaluateCell evaluates a single decision table cell expression against the provided actual value.
func (e *Engine) evaluateCell(expression string, actual gjson.Result, logStack *[]models.LogStackEntry) (bool, []models.LogStackEntry) {
	localLogs := []models.LogStackEntry{}
	addInfoLog(logStack, fmt.Sprintf("Evaluating Cell: Expr='%s' Actual='%s'", expression, actual.String()))
	expr := strings.TrimSpace(expression)
	if expr == "" || expr == "-" {
		return true, localLogs
	}
	if strings.ToUpper(expr) == "NULL" {
		addInfoLog(logStack, "Matched NULL check")
		return actual.Type == gjson.Null, localLogs
	}
	if strings.ToUpper(expr) == "EMPTY" {
		addInfoLog(logStack, "Matched EMPTY check")
		return actual.String() == "", localLogs
	}
	actualStr := actual.String()
	if strings.HasPrefix(expr, "!") {
		addInfoLog(logStack, "Matched Negation check")
		return actualStr != strings.TrimSpace(expr[1:]), localLogs
	}
	if strings.Contains(expr, ",") && !strings.ContainsAny(expr, "[]{}()") {
		addInfoLog(logStack, "Matched List check")
		parts := strings.Split(expr, ",")
		for _, p := range parts {
			if matched, logs := e.evaluateCell(strings.TrimSpace(p), actual, logStack); matched {
				addLogs(logStack, logs)
				return true, localLogs
			}
		}
		return false, localLogs
	}
	isNumericExpr := strings.ContainsAny(expr, "><=") || strings.Contains(expr, "..")
	addInfoLog(logStack, fmt.Sprintf("isNumericExpr=%v, actual.Type=%v", isNumericExpr, actual.Type))
	if (actual.Type == gjson.Number) || isNumericExpr {
		actualNum := actual.Float()
		addInfoLog(logStack, fmt.Sprintf("Evaluating Numeric: ActualNum=%f", actualNum))
		if strings.Contains(expr, "..") {
			addInfoLog(logStack, "Matched Range check")
			parts := strings.Split(expr, "..")
			if len(parts) == 2 {
				addInfoLog(logStack, fmt.Sprintf("Range Parts: Min='%s' Max='%s'", parts[0], parts[1]))
				min, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
				max, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
				if err1 == nil && err2 == nil {
					addInfoLog(logStack, fmt.Sprintf("Parsed Range: Min=%f Max=%f", min, max))
					return actualNum >= min && actualNum <= max, localLogs
				}
			}
		}
		if strings.HasPrefix(expr, ">=") {
			v, err := strconv.ParseFloat(strings.TrimSpace(expr[2:]), 64)
			addInfoLog(logStack, fmt.Sprintf("Matched >= check against %f", v))
			return err == nil && actualNum >= v, localLogs
		}
		if strings.HasPrefix(expr, "<=") {
			v, err := strconv.ParseFloat(strings.TrimSpace(expr[2:]), 64)
			addInfoLog(logStack, fmt.Sprintf("Matched <= check against %f", v))
			return err == nil && actualNum <= v, localLogs
		}
		if strings.HasPrefix(expr, ">") {
			v, err := strconv.ParseFloat(strings.TrimSpace(expr[1:]), 64)
			addInfoLog(logStack, fmt.Sprintf("Matched > check against %f", v))
			return err == nil && actualNum > v, localLogs
		}
		if strings.HasPrefix(expr, "<") {
			v, err := strconv.ParseFloat(strings.TrimSpace(expr[1:]), 64)
			addInfoLog(logStack, fmt.Sprintf("Matched < check against %f", v))
			return err == nil && actualNum < v, localLogs
		}
	}
	if strings.HasPrefix(expr, "~") {
		addInfoLog(logStack, "Matched Regex check")
		pattern := strings.TrimSpace(expr[1:])
		var re *regexp.Regexp
		if cached, ok := e.regexCache.Load(pattern); ok {
			addInfoLog(logStack, "Using cached regex")
			re = cached.(*regexp.Regexp)
		} else {
			addInfoLog(logStack, "Compiling new regex")
			var err error
			re, err = regexp.Compile(pattern)
			if err != nil {
				addErrorLog(logStack, fmt.Sprintf("Regex compilation failed: %v", err))
				return false, localLogs
			}
			e.regexCache.Store(pattern, re)
		}
		return re.MatchString(actualStr), localLogs
	}
	if strings.HasPrefix(expr, "^") {
		addInfoLog(logStack, "Matched Starts With check")
		return strings.HasPrefix(actualStr, strings.TrimSpace(expr[1:])), localLogs
	}
	if strings.HasPrefix(expr, "$") {
		addInfoLog(logStack, "Matched Ends With check")
		return strings.HasSuffix(actualStr, strings.TrimSpace(expr[1:])), localLogs
	}
	addInfoLog(logStack, "Matched Exact Match check")
	return actualStr == expr, localLogs
}

// Helper
func toFloat(v interface{}) float64 {
	switch i := v.(type) {
	case float64:
		return i
	case int:
		return float64(i)
	case int64:
		return float64(i)
	case string:
		f, _ := strconv.ParseFloat(i, 64)
		return f
	}
	return 0
}

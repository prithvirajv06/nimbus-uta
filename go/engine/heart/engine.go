package heart

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

	"github.com/prithvirajv06/nimbus-uta/go/engine/models"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

func NewEngine() *Engine {
	return &Engine{}
}

// --- Decision Table Execution ---

func (e *Engine) ProcessDecisionTable(ctx context.Context, table models.DecisionTable, data []byte) ([]byte, []models.LogStackEntry, error) {
	output := data
	var matchedRows [][]models.Variables
	logStack := []models.LogStackEntry{}
	// 1. Collect ALL Matching Rows
	for _, ruleRow := range table.Rules {
		// Check context for cancellation
		if ctx.Err() != nil {
			addErrorLog(&logStack, "Context cancelled during decision table processing")
			return nil, logStack, ctx.Err()
		}

		isMatch := true
		for i, inputDef := range table.InputsColumns {
			addInfoLog(&logStack, fmt.Sprintf("Evaluating Input Column: %s", inputDef.VarKey))
			actualVal := gjson.GetBytes(output, inputDef.VarKey)
			if cond, logs := e.evaluateCell(ruleRow[i].Value, actualVal, &logStack); !cond {
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
		return nil, logStack, err
	}

	for variable, value := range finalValues {
		output, err = sjson.SetBytes(output, variable, value)
		if err != nil {
			return nil, logStack, err
		}
	}

	return output, logStack, nil
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
			valStr := row[valIndex].Value
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
		if outDef.IsPriority {
			idx := len(table.InputsColumns) + i
			pA, _ := strconv.ParseFloat(rowA[idx].Value, 64)
			pB, _ := strconv.ParseFloat(rowB[idx].Value, 64)
			if pA != pB {
				return pA > pB
			}
		}
	}
	return false
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

// --- RuleSet Execution (Reusing previous robust logic) ---

func (e *Engine) ExecuteRuleSet(ctx context.Context, rs models.LogicFlow, data []byte) ([]byte, []models.LogStackEntry, error) {
	currentData := data
	tempVariables := make(map[string]interface{})
	logStack := []models.LogStackEntry{}
	for _, rule := range rs.LogicalSteps {
		addInfoLog(&logStack, fmt.Sprintf("Processing Rule: %s", rule.OperationName))
		select {
		case <-ctx.Done():
			return nil, logStack, ctx.Err()
		default:
		}

		matched := false
		if rule.Condition.Operator != "" {
			addInfoLog(&logStack, fmt.Sprintf("Evaluating Condition for Rule: %s", rule.OperationName))
			matchedSub, logSub := e.evalCond(rule.Condition, currentData, &logStack)
			matched = matchedSub
			addLogs(&logStack, logSub)
			addResultLog(&logStack, fmt.Sprintf("Condition Result for Rule %s: %v", rule.OperationName, matchedSub))
		}

		ops := rule.OperationIfFalse
		if matched {
			addInfoLog(&logStack, fmt.Sprintf("Condition matched for Rule: %s", rule.OperationName))
			ops = rule.OperationIfTrue
		}

		for _, op := range ops {
			var err error
			if strings.Contains(op.Variable.VarKey, "[*]") {
				addInfoLog(&logStack, fmt.Sprintf("Applying Array Operation on Variable: %s", op.Variable))
				currentData, _, err = e.walkArray(op, currentData, &logStack, &tempVariables)
			} else {
				addInfoLog(&logStack, fmt.Sprintf("Applying Operation on Variable: %s", op.Variable))
				currentData, err = e.applyOp(op, currentData, op.Variable.VarKey, &tempVariables)
			}
			if err != nil {
				return nil, logStack, err
			}
		}
	}
	return currentData, logStack, nil
}

func (e *Engine) evalCond(c models.Condition, data []byte, logStack *[]models.LogStackEntry) (bool, []models.LogStackEntry) {
	localLogs := []models.LogStackEntry{}
	if len(c.Conditions) > 0 {
		switch strings.ToUpper(c.Operator) {
		case "AND":
			addInfoLog(logStack, "Evaluating AND condition group")
			for _, sub := range c.Conditions {
				addInfoLog(logStack, "Evaluating sub-condition of AND group")
				if res, logSub := e.evalCond(sub, data, logStack); !res {
					addLogs(logStack, logSub)
					return false, localLogs
				}
			}
			return true, localLogs
		case "OR":
			addInfoLog(logStack, "Evaluating OR condition group")
			for _, sub := range c.Conditions {
				addInfoLog(logStack, "Evaluating sub-condition of OR group")
				if res, logSub := e.evalCond(sub, data, logStack); res {
					addLogs(logStack, logSub)
					return true, localLogs
				}
			}
			return false, localLogs
		case "NOT":
			addInfoLog(logStack, "Evaluating NOT condition")
			res, logSub := e.evalCond(c.Conditions[0], data, logStack)
			addLogs(logStack, logSub)
			return !res, localLogs
		}
	}

	if c.Variable.VarKey != "" {
		addInfoLog(logStack, fmt.Sprintf("Evaluating Leaf Condition: Variable='%s' Logical='%s' Value='%v'", c.Variable, c.Logical, c.OpValue))
		if strings.Contains(c.Variable.VarKey, "[*]") {
			return e.evalArrayCondition(c, data, logStack)
		}
		actual := gjson.GetBytes(data, c.Variable.VarKey)
		result, logSub := e.compareDirect(actual, c.Logical, c.OpValue, logStack)
		addLogs(logStack, logSub)
		return result, localLogs
	}

	return true, localLogs
}

func (e *Engine) evalArrayCondition(c models.Condition, data []byte, logStack *[]models.LogStackEntry) (bool, []models.LogStackEntry) {
	localLogs := []models.LogStackEntry{}
	addInfoLog(logStack, fmt.Sprintf("Evaluating Array Condition: Variable='%s' Logical='%s' Value='%v'", c.Variable, c.Logical, c.OpValue))
	parts := strings.Split(c.Variable.VarKey, "[*]")

	var search func(d []byte, prefix string, idx int) (bool, []models.LogStackEntry)
	search = func(d []byte, prefix string, idx int) (bool, []models.LogStackEntry) {
		segment := strings.Trim(parts[idx], ".")
		fullPath := segment
		if prefix != "" {
			fullPath = prefix + "." + segment
		}

		if idx == len(parts)-1 {
			addInfoLog(logStack, fmt.Sprintf("Evaluating Leaf Condition at Path: %s", fullPath))
			actual := gjson.GetBytes(d, fullPath)
			result, logSub := e.compareDirect(actual, c.Logical, c.OpValue, logStack)
			addLogs(logStack, logSub)
			res := result
			addResultLog(logStack, fmt.Sprintf("Comparison Result: %v", res))
			if strings.ToUpper(c.Operator) == "NOT" {
				addInfoLog(logStack, "Applying NOT operator at leaf level")
				return !res, localLogs
			}
			addInfoLog(logStack, "Returning comparison result at leaf level")
			return res, localLogs
		}

		res := gjson.GetBytes(d, fullPath)
		if !res.IsArray() {
			return false, localLogs
		}

		matchesCount := 0
		totalElements := 0

		res.ForEach(func(key, value gjson.Result) bool {
			totalElements++
			itemPath := fmt.Sprintf("%s.%d", fullPath, key.Int())
			addInfoLog(logStack, fmt.Sprintf("Processing Array Item at Path: %s", itemPath))
			if res, subLog := search(d, itemPath, idx+1); res {
				matchesCount++
				addLogs(logStack, subLog)
				if strings.ToUpper(c.Operator) != "AND" && strings.ToUpper(c.Operator) != "ALL" {
					addInfoLog(logStack, "Short-circuiting on first match for OR/ANY operator")
					return false
				}
			} else {
				addLogs(logStack, subLog)
				if strings.ToUpper(c.Operator) == "AND" || strings.ToUpper(c.Operator) == "ALL" {
					addInfoLog(logStack, "Short-circuiting on first failure for AND/ALL operator")
					return false
				}
			}
			return true
		})

		switch strings.ToUpper(c.Operator) {
		case "AND", "ALL":
			addInfoLog(logStack, "Evaluating ALL condition result")
			return matchesCount == totalElements && totalElements > 0, localLogs
		case "NOT":
			addInfoLog(logStack, "Evaluating NOT condition result")
			return matchesCount == 0, localLogs
		default:
			addInfoLog(logStack, "Evaluating ANY condition result")
			return matchesCount > 0, localLogs
		}
	}

	return search(data, "", 0)
}

func (e *Engine) compareDirect(actual gjson.Result, logical string, expected interface{}, logStack *[]models.LogStackEntry) (bool, []models.LogStackEntry) {
	localLogs := []models.LogStackEntry{}
	addInfoLog(logStack, fmt.Sprintf("Comparing Actual='%s' Logical='%s' Expected='%v'", actual.String(), logical, expected))
	switch strings.ToLower(logical) {
	case "empty":
		result := actual.String() == "" || actual.Type == gjson.Null || (actual.Type == gjson.JSON && len(actual.Array()) == 0) || (actual.Type == gjson.JSON && len(actual.Map()) == 0) || (actual.Type == gjson.String && actual.String() == "null") || (actual.Type == gjson.String && actual.String() == "undefined")
		addResultLog(logStack, fmt.Sprintf("Expected EMPTY, actual input %v Case EMPTY: Result=%v", actual.String(), result))
		return result, localLogs
	case "equal", "eq":
		result := actual.String() == fmt.Sprintf("%v", expected)
		addResultLog(logStack, fmt.Sprintf("Expected %v, actual input %v Case EQUAL: Result=%v", expected, actual.String(), result))
		return result, localLogs
	case "notequal", "neq":
		result := actual.String() != fmt.Sprintf("%v", expected)
		addResultLog(logStack, fmt.Sprintf("Expected %v, actual input %v Case NOTEQUAL: Result=%v", expected, actual.String(), result))
		return result, localLogs
	case "gt":
		result := actual.Float() > toFloat(expected)
		addResultLog(logStack, fmt.Sprintf("Expected %v, actual input %v Case GT: Result=%v", expected, actual.String(), result))
		return result, localLogs
	case "lt":
		result := actual.Float() < toFloat(expected)
		addResultLog(logStack, fmt.Sprintf("Expected %v, actual input %v Case LT: Result=%v", expected, actual.String(), result))
		return result, localLogs
	case "gte":
		result := actual.Float() >= toFloat(expected)
		addResultLog(logStack, fmt.Sprintf("Expected %v, actual input %v Case GTE: Result=%v", expected, actual.String(), result))
		return result, localLogs
	case "lte":
		result := actual.Float() <= toFloat(expected)
		addResultLog(logStack, fmt.Sprintf("Expected %v, actual input %v Case LTE: Result=%v", expected, actual.String(), result))
		return result, localLogs
	case "contains":
		result := strings.Contains(actual.String(), fmt.Sprintf("%v", expected))
		addResultLog(logStack, fmt.Sprintf("Expected %v, actual input %v Case CONTAINS: Result=%v", expected, actual.String(), result))
		return result, localLogs
	case "startswith":
		result := strings.HasPrefix(actual.String(), fmt.Sprintf("%v", expected))
		addResultLog(logStack, fmt.Sprintf("Expected %v, actual input %v Case STARTSWITH: Result=%v", expected, actual.String(), result))
		return result, localLogs
	case "endswith":
		result := strings.HasSuffix(actual.String(), fmt.Sprintf("%v", expected))
		addResultLog(logStack, fmt.Sprintf("Expected %v, actual input %v Case ENDSWITH: Result=%v", expected, actual.String(), result))
		return result, localLogs
	case "matches":
		pattern := fmt.Sprintf("%v", expected)
		var re *regexp.Regexp
		if cached, ok := e.regexCache.Load(pattern); ok {
			re = cached.(*regexp.Regexp)
			addResultLog(logStack, fmt.Sprintf("Expected %v, actual input %v Case MATCHES: Using cached regex", expected, actual.String()))
		} else {
			var err error
			re, err = regexp.Compile(pattern)
			if err != nil {
				addErrorLog(logStack, fmt.Sprintf("Expected %v, actual input %v Case MATCHES: Regex compilation failed: %v", expected, actual.String(), err))
				return false, localLogs
			}
			e.regexCache.Store(pattern, re)
			addResultLog(logStack, fmt.Sprintf("Expected %v, actual input %v Case MATCHES: Compiled and cached new regex", expected, actual.String()))
		}
		result := re.MatchString(actual.String())
		addResultLog(logStack, fmt.Sprintf("Expected %v, actual input %v Case MATCHES: Result=%v", expected, actual.String(), result))
		return result, localLogs
	case "inarray":
		strVal := actual.String()
		if list, ok := expected.([]interface{}); ok {
			for _, item := range list {
				if fmt.Sprintf("%v", item) == strVal {
					addResultLog(logStack, fmt.Sprintf("Expected %v, actual input %v Case INARRAY: Value found in array", expected, actual.String()))
					return true, localLogs
				}
			}
		}
		addResultLog(logStack, fmt.Sprintf("Expected %v, actual input %v Case INARRAY: Value not found in array", expected, actual.String()))
		return false, localLogs
	case "notinarray":
		strVal := actual.String()
		if list, ok := expected.([]interface{}); ok {
			for _, item := range list {
				if fmt.Sprintf("%v", item) == strVal {
					addInfoLog(logStack, fmt.Sprintf("Expected %v, actual input %v Case NOTINARRAY: Value found in array", expected, actual.String()))
					return false, localLogs
				}
			}
		}
		addResultLog(logStack, fmt.Sprintf("Expected %v, actual input %v Case NOTINARRAY: Value not found in array", expected, actual.String()))
		return true, localLogs
	}
	addResultLog(logStack, "Case DEFAULT: No match, returning false")
	return false, localLogs
}

func (e *Engine) walkArray(op models.Operation, data []byte, logStack *[]models.LogStackEntry, tempVariables *map[string]interface{}) ([]byte, []models.LogStackEntry, error) {
	localLogs := []models.LogStackEntry{}
	parts := strings.Split(op.Variable.VarKey, "[*]")
	var walk func(d []byte, prefix string, idx int) ([]byte, []models.LogStackEntry, error)
	walk = func(d []byte, prefix string, idx int) ([]byte, []models.LogStackEntry, error) {
		addInfoLog(logStack, fmt.Sprintf("Walking Array: Prefix='%s' Index=%d", prefix, idx))
		segment := strings.Trim(parts[idx], ".")
		fullPath := segment
		if prefix != "" {
			fullPath = prefix + "." + segment
		}

		if idx == len(parts)-1 {
			result, err := e.applyOp(op, d, fullPath, tempVariables)
			if err != nil {
				return d, localLogs, err
			}
			return result, localLogs, nil
		}

		res := gjson.GetBytes(d, fullPath)
		if !res.IsArray() {
			addInfoLog(logStack, fmt.Sprintf("Path '%s' is not an array. Skipping.", fullPath))
			return d, localLogs, nil
		}

		updated := d
		res.ForEach(func(key, value gjson.Result) bool {
			itemPath := fmt.Sprintf("%s.%d", fullPath, key.Int())
			filterPassed, logSub := e.checkFilters(itemPath, segment, op.ArrayFilters, updated, logStack)
			addLogs(logStack, logSub)
			if !filterPassed {
				return true
			}
			var err error
			updated, logSub, err = walk(updated, itemPath, idx+1)
			addLogs(logStack, logSub)
			return err == nil
		})
		return updated, localLogs, nil
	}
	return walk(data, "", 0)
}

func (e *Engine) applyOp(op models.Operation, data []byte, path string, tempVariables *map[string]interface{}) ([]byte, error) {
	current := gjson.GetBytes(data, path)
	var val interface{} = op.OpValue
	//Set object, array, or primitive value
	if op.ValueIsPath {
		val = gjson.Get(string(data), op.OpValue.(string)).Value()
	}
	switch strings.ToUpper(op.Operation) {
	case "CREATE_TEMP_OBJ":
		tempObj := make(map[string]interface{})
		(*tempVariables)[path] = tempObj
		return data, nil
	case "SET":
		return sjson.SetBytes(data, path, val)
	case "ADD":
		return sjson.SetBytes(data, path, current.Float()+toFloat(val))
	case "MULT":
		return sjson.SetBytes(data, path, current.Float()*toFloat(val))
	case "SUB":
		return sjson.SetBytes(data, path, current.Float()-toFloat(val))
	case "SET_OBJ":
		return sjson.SetBytes(data, path, val)
	case "COLLECT", "COLLECT_SUM", "COLLECT_COUNT":
		return e.applyAggregation(data, path, fmt.Sprintf("%v", val), strings.ToUpper(op.Operation))
	case "DELETE":
		return sjson.DeleteBytes(data, path)
	case "PUSH":
		var arr []interface{}
		if current.IsArray() {
			if v, ok := current.Value().([]interface{}); ok {
				arr = v
			}
		}
		if strings.Contains(fmt.Sprintf("%v", val), "{") || strings.Contains(fmt.Sprintf("%v", val), "[") {
			var parsed interface{}
			err := json.Unmarshal([]byte(fmt.Sprintf("%v", val)), &parsed)
			if err == nil {
				val = parsed
			}
		}
		arr = append(arr, val)
		return sjson.SetBytes(data, path, arr)
	case "REMOVE":
		return sjson.DeleteBytes(data, path)
	case "CLEAR":
		if current.IsArray() {
			return sjson.SetBytes(data, path, []interface{}{})
		}
		return sjson.DeleteBytes(data, path)
	case "UPPERCASE":
		return sjson.SetBytes(data, path, cases.Upper(language.AmericanEnglish).String(current.String()))
	case "LOWERCASE":
		return sjson.SetBytes(data, path, cases.Lower(language.AmericanEnglish).String(current.String()))
	case "TRIM":
		return sjson.SetBytes(data, path, strings.TrimSpace(current.String()))
	case "APPEND":
		return sjson.SetBytes(data, path, current.String()+fmt.Sprintf("%v", val))
	case "PREPEND":
		return sjson.SetBytes(data, path, fmt.Sprintf("%v", val)+current.String())
	case "INCREMENT":
		return sjson.SetBytes(data, path, current.Int()+1)
	case "DECREMENT":
		return sjson.SetBytes(data, path, current.Int()-1)
	case "TOGGLE":
		if current.Type == gjson.True {
			return sjson.SetBytes(data, path, false)
		}
		return sjson.SetBytes(data, path, true)
	case "REVERSE":
		if current.IsArray() {
			var arr []interface{}
			if v, ok := current.Value().([]interface{}); ok {
				arr = v
			}
			for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
				arr[i], arr[j] = arr[j], arr[i]
			}
			return sjson.SetBytes(data, path, arr)
		}
		return data, nil
	case "SORT_ASC":
		if current.IsArray() {
			var arr []float64
			if v, ok := current.Value().([]interface{}); ok {
				for _, item := range v {
					arr = append(arr, toFloat(item))
				}
			}
			sort.Float64s(arr)
			return sjson.SetBytes(data, path, arr)
		}
	case "SORT_DESC":
		if current.IsArray() {
			var arr []float64
			if v, ok := current.Value().([]interface{}); ok {
				for _, item := range v {
					arr = append(arr, toFloat(item))
				}
			}
			sort.Sort(sort.Reverse(sort.Float64Slice(arr)))
			return sjson.SetBytes(data, path, arr)
		}
	}
	return data, nil
}

func (e *Engine) checkFilters(path, name string, filters []models.ArrayFilter, data []byte, logStack *[]models.LogStackEntry) (bool, []models.LogStackEntry) {
	localLogs := []models.LogStackEntry{}
	addInfoLog(logStack, fmt.Sprintf("Checking Filters for Path='%s' Name='%s'", path, name))
	for _, f := range filters {
		if strings.Contains(name, f.ArrayName) {
			addInfoLog(logStack, fmt.Sprintf("Applying Filter: Property='%s' Logical='%s' Value='%v'", f.Property, f.Logical, f.OpValue))
			val := gjson.GetBytes(data, path+"."+f.Property)
			compare, logSub := e.compareDirect(val, f.Logical, f.OpValue, logStack)
			addLogs(logStack, logSub)
			if !compare {
				return false, localLogs
			}
		}
	}
	return true, localLogs
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

package engine

import (
	"fmt"
	"strings"

	"github.com/prithvirajv06/nimbus-uta/go/engine/internal/models"
)

func processConditionConfig(re *RuleEngine, cond models.ConditionConfig) bool {
	// Process LeftVar
	leftPath := cond.LeftVar.VarKey
	re.AddLog(fmt.Sprintf("Processing LeftVar: %s", leftPath))
	leftValue, err := re.GetValue(leftPath)
	if err != nil {
		re.AddLog(fmt.Sprintf("Failed to get LeftVar value at %s: %v", leftPath, err))
		return false
	}
	re.AddLog(fmt.Sprintf("LeftVar value: %v", leftValue))

	// Process RightValue
	rightValue := cond.RightValue
	if rightStr, ok := cond.RightValue.(string); ok && re.isExpression(rightStr) {
		re.AddLog(fmt.Sprintf("Evaluating RightValue expression: %s", rightStr))
		rightValue, err = re.evaluateExpression(rightStr, nil)
		if err != nil {
			re.AddLog(fmt.Sprintf("Failed to evaluate RightValue expression '%s': %v", rightStr, err))
			return false
		}
		re.AddLog(fmt.Sprintf("RightValue evaluated to: %v", rightValue))
	}
	switch cond.Operator {
	case "contains":
		re.AddLog(fmt.Sprintf("Checking if LeftVar contains RightValue"))
		switch lv := leftValue.(type) {
		case string:
			rv, ok := rightValue.(string)
			if ok && strings.Contains(lv, rv) {
				re.AddLog("Condition met: LeftVar contains RightValue")
				return true
			} else {
				re.AddLog("Condition not met: LeftVar does not contain RightValue")
				return false
			}
		case []interface{}:
			found := false
			for _, item := range lv {
				if item == rightValue {
					found = true
					break
				}
			}
			if found {
				re.AddLog("Condition met: LeftVar slice contains RightValue")
				return true
			} else {
				re.AddLog("Condition not met: LeftVar slice does not contain RightValue")
				return false
			}
		default:
			re.AddLog("LeftVar is neither string nor slice for 'contains' operator")
			return false
		}
	case "not_contains":
		re.AddLog(fmt.Sprintf("Checking if LeftVar does not contain RightValue"))
		switch lv := leftValue.(type) {
		case string:
			rv, ok := rightValue.(string)
			if ok && !strings.Contains(lv, rv) {
				re.AddLog("Condition met: LeftVar does not contain RightValue")
				return true
			} else {
				return false
			}
		case []interface{}:
			found := false
			for _, item := range lv {
				if item == rightValue {
					found = true
					break
				}
				if !found {
					re.AddLog("Condition met: LeftVar slice does not contain RightValue")
					return true
				}
			}
		default:
			re.AddLog("LeftVar is neither string nor slice for 'not_contains' operator")
			return false
		}
	case "is_empty":
		re.AddLog(fmt.Sprintf("Checking if LeftVar is empty"))
		isEmpty := false
		switch lv := leftValue.(type) {
		case string:
			isEmpty = lv == ""
		case []interface{}:
			isEmpty = len(lv) == 0
		}
		if isEmpty {
			re.AddLog("Condition met: LeftVar is empty")
			return true
		}
	case "is_not_empty":
		re.AddLog(fmt.Sprintf("Checking if LeftVar is not empty"))
		isNotEmpty := false
		switch lv := leftValue.(type) {
		case string:
			isNotEmpty = lv != ""
		case []interface{}:
			isNotEmpty = len(lv) > 0
		}
		if isNotEmpty {
			re.AddLog("Condition met: LeftVar is not empty")
			return true
		}
	case "starts_with":
		re.AddLog(fmt.Sprintf("Checking if LeftVar starts with RightValue"))
		if lv, ok := leftValue.(string); ok {
			if rv, ok := rightValue.(string); ok && strings.HasPrefix(lv, rv) {
				re.AddLog("Condition met: LeftVar starts with RightValue")
				return true
			}
		}
	case "ends_with":
		re.AddLog(fmt.Sprintf("Checking if LeftVar ends with RightValue"))
		if lv, ok := leftValue.(string); ok {
			if rv, ok := rightValue.(string); ok && strings.HasSuffix(lv, rv) {
				re.AddLog("Condition met: LeftVar ends with RightValue")
				return true
			}
		}
	default:
		re.AddLog(fmt.Sprintf("Unsupported operator in ConditionConfig: %s", cond.Operator))
		return false
	}
	re.AddLog(fmt.Sprintf("ConditionConfig processed: LeftValue=%v, RightValue=%v, Operator=%s", leftValue, rightValue, cond.Operator))
	return false
}

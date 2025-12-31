package heart

import (
	"fmt"
	"strings"

	"github.com/prithvirajv06/nimbus-uta/go/engine/models"
)

// ParseLogicFlowFromText parses a ruleset text into a LogicFlow struct.
// Example rule text (one per line):
//
//	IF age > 18 AND status == 'active' THEN score = 10; flag = true
//	IF amount <= 1000 OR vip == true THEN discount = 5
//	ELSE score = 0
func ParseLogicFlowFromText(name, description, text string) (models.LogicFlow, error) {
	var steps []models.LogicalStep
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var condStr, thenStr string
		if strings.HasPrefix(strings.ToUpper(line), "IF ") {
			parts := strings.SplitN(line, "THEN", 2)
			if len(parts) != 2 {
				return models.LogicFlow{}, fmt.Errorf("invalid rule: %s", line)
			}
			condStr = strings.TrimSpace(parts[0][3:])
			thenStr = strings.TrimSpace(parts[1])
		} else if strings.HasPrefix(strings.ToUpper(line), "ELSE ") {
			condStr = ""
			thenStr = strings.TrimSpace(line[5:])
		} else if strings.ToUpper(line) == "ELSE" {
			condStr = ""
			thenStr = ""
		} else {
			return models.LogicFlow{}, fmt.Errorf("invalid rule: %s", line)
		}

		var cond models.Condition
		if condStr != "" {
			cond = parseCondition(condStr)
		}
		ops := parseOperations(thenStr)
		step := models.LogicalStep{
			OperationName:    line,
			Condition:        cond,
			OperationIfTrue:  ops,
			OperationIfFalse: nil,
		}
		steps = append(steps, step)
	}
	return models.LogicFlow{
		Name:         name,
		Description:  description,
		LogicalSteps: steps,
	}, nil
}

// parseOperations parses a string like "score = 10; flag = true" into []Operation
func parseOperations(thenStr string) []models.Operation {
	var ops []models.Operation
	stmts := strings.Split(thenStr, ";")
	for _, stmt := range stmts {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		tokens := tokenize(stmt)
		if len(tokens) < 3 {
			continue
		}
		varKey := tokens[0]
		op := tokens[1]
		opValue := strings.Join(tokens[2:], " ")
		opValue = strings.Trim(opValue, "'\"")
		ops = append(ops, models.Operation{
			Variable:  models.Variables{VarKey: varKey},
			Operation: mapAssignmentOp(op),
			OpValue:   opValue,
		})
	}
	return ops
}

// tokenize splits a string into tokens, handling quoted values.
func tokenize(s string) []string {
	var tokens []string
	var curr strings.Builder
	inQuote := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == ' ' && !inQuote {
			if curr.Len() > 0 {
				tokens = append(tokens, curr.String())
				curr.Reset()
			}
		} else if c == '\'' || c == '"' {
			inQuote = !inQuote
		} else {
			curr.WriteByte(c)
		}
	}
	if curr.Len() > 0 {
		tokens = append(tokens, curr.String())
	}
	return tokens
}

// ...existing code...

// mapOp maps text operators and wording to canonical logical names.
func mapOp(op string) string {
	switch strings.ToLower(op) {
	case ">", "greater", "greaterthan", "greater than":
		return "gt"
	case "<", "less", "lessthan", "less than":
		return "lt"
	case ">=", "greaterorequal", "greater or equal", "greaterthanorequal", "greater than or equal":
		return "gte"
	case "<=", "lessorequal", "less or equal", "lessthanorequal", "less than or equal":
		return "lte"
	case "==", "=", "equals", "equal", "is":
		return "equal"
	case "!=", "not", "not=", "notequals", "not equals", "does not equal":
		return "notequal"
	case "contains":
		return "contains"
	case "startswith", "starts with":
		return "startswith"
	case "endswith", "ends with":
		return "endswith"
	case "in", "inarray":
		return "inarray"
	case "notin", "not in":
		return "notinarray"
	default:
		return op
	}
}

// parseCondition parses a condition string into a Condition struct.
func parseCondition(condStr string) models.Condition {
	condStr = strings.TrimSpace(condStr)
	upper := strings.ToUpper(condStr)
	// Handle AND/OR/NOT at top level (wording and symbol)
	if idx := strings.Index(upper, " AND "); idx != -1 {
		parts := strings.Split(condStr, " AND ")
		var subs []models.Condition
		for _, p := range parts {
			subs = append(subs, parseCondition(p))
		}
		return models.Condition{Operator: "AND", Conditions: subs}
	}
	if idx := strings.Index(upper, " OR "); idx != -1 {
		parts := strings.Split(condStr, " OR ")
		var subs []models.Condition
		for _, p := range parts {
			subs = append(subs, parseCondition(p))
		}
		return models.Condition{Operator: "OR", Conditions: subs}
	}
	if strings.HasPrefix(upper, "NOT ") {
		sub := parseCondition(condStr[4:])
		return models.Condition{Operator: "NOT", Conditions: []models.Condition{sub}}
	}
	// Support symbol forms too
	if idx := strings.Index(condStr, "&&"); idx != -1 {
		parts := strings.Split(condStr, "&&")
		var subs []models.Condition
		for _, p := range parts {
			subs = append(subs, parseCondition(p))
		}
		return models.Condition{Operator: "AND", Conditions: subs}
	}
	if idx := strings.Index(condStr, "||"); idx != -1 {
		parts := strings.Split(condStr, "||")
		var subs []models.Condition
		for _, p := range parts {
			subs = append(subs, parseCondition(p))
		}
		return models.Condition{Operator: "OR", Conditions: subs}
	}
	if strings.HasPrefix(condStr, "!") {
		sub := parseCondition(condStr[1:])
		return models.Condition{Operator: "NOT", Conditions: []models.Condition{sub}}
	}
	// Parse leaf: var op value (support multi-word ops)
	tokens := tokenize(condStr)
	if len(tokens) < 3 {
		return models.Condition{}
	}
	varKey := tokens[0]
	// Try to match two-word operators
	logical := tokens[1]
	opValueStart := 2
	if len(tokens) > 3 {
		twoWordOp := strings.ToLower(tokens[1] + " " + tokens[2])
		if mapped := mapOp(twoWordOp); mapped != twoWordOp {
			logical = twoWordOp
			opValueStart = 3
		}
	}
	opValue := strings.Join(tokens[opValueStart:], " ")
	opValue = strings.Trim(opValue, "'\"")
	return models.Condition{
		Variable: models.Variables{VarKey: varKey},
		Logical:  mapOp(logical),
		OpValue:  opValue,
	}
}

// ...existing code...

// mapAssignmentOp maps assignment ops to operation names.
func mapAssignmentOp(op string) string {
	switch op {
	case "=":
		return "SET"
	case "+=":
		return "ADD"
	case "-=":
		return "SUB"
	case "*=":
		return "MULT"
	default:
		return op
	}
}

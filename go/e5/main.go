package main

import "smart-rule-engine-dt/models"

func main() {

	testDTRequest := models.NewDTRequestStruct()

	newRule := models.DTRule{
		RuleID:      "rule-001",
		Description: "This is a test rule",
		Severity:    "high",
		Action:      "block",
		Condition: Condition{
			Operator: "AND",
			Operands: map[string]interface{}{
				"ip": "192.168.1.1",
			},
		},
	}

	testDTRequest.AddDTRule(newRule)
}

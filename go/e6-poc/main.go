package main

import (
	"html/template"
	"smart-rule-engine-dt/models"
	"strings"

	"github.com/mailru/easyjson/buffer"
)

func main() {

	testDTRequest := models.NewDTRequestStruct()

	newRule := models.DTRule{
		RuleID:      "rule-001",
		Description: "This is a test rule",
		Severity:    "high",
		Action: []models.Action{
			{
				Type: "set",
				Target: models.Variables{
					VarKey: "profile.status",
					Type:   "string",
				},
				Value: "active",
			}},
		Condition: []models.Condition{
			{
				Operator: "AND",
				Operands: models.Variables{
					VarKey: "profile.credits",
					Type:   "number",
					Value:  18,
				},
			},
		},
	}

	arrayRule := models.DTRule{
		RuleID:      "rule-002",
		Description: "This is another test rule with array operands",
		Severity:    "medium",
		Action: []models.Action{
			{
				Type: "set",
				Target: models.Variables{
					VarKey: "users[*].note",
					Type:   "string",
				},
				Value: "User age condition met",
			}},
		Condition: []models.Condition{
			{
				Operator: "OR",
				Operands: models.Variables{
					VarKey: "users[*].roles",
					Type:   "array",
					Value:  []string{"admin", "editor"},
				},
			},
		},
	}

	testDTRequest.AddDTRule(newRule)
	testDTRequest.AddDTRule(arrayRule)
}

var (
	assignTmpl = template.Must(template.New("assign").Parse(`
{{.Pad}}addLog('Assigning {{.ValStr}} to {{.TargetStr}}');
{{.Pad}}{{.TargetStr}} = {{.ValStr}};
`))
	conditionTmpl = template.Must(template.New("condition").Parse(`
{{.Pad}}if ({{.LeftStr}} {{.Operator}} {{.RightStr}}) {
`))
	elseTmpl = template.Must(template.New("else").Parse(`
{{.Pad}}} else {
`))
	endIfTmpl = template.Must(template.New("endIf").Parse(`
{{.Pad}}}
`))

	actionArrayTmpl = template.Must(template.New("actionArray").Parse(`
{{.Pad}}for (var i = 0; i < {{.ArrayStr}}.length; i++) {
{{.Pad}}    var element = {{.ArrayStr}}[i];
{{.Pad}}    // Apply action to element
{{.Pad}}}
`))
)

func processDTRequest(dtRequest *models.DTRequestStruct) {
	// Example processing logic
	buffer.Buffer{}
	for _, rule := range dtRequest.DTRules {

		// If input Containms array operands need to procees for each element
		if strings.Contains(rule.Condition.Operands.VarKey, "[*]") {
			// For each element in the array, evaluate the condition, Will create a JS template String to handle this
			// Pseudo code:
			/*
				for each element in array:
					if evaluateCondition(element, rule.Condition):
						applyAction(element, rule.Action)
			*/
			arrayStep := strings.Split(rule.Condition.Operands.VarKey, "[*]")
			for as := range arrayStep {

			}
		}

	}

}

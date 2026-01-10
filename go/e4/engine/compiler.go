package engine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/prithvirajv06/nimbus-uta/go/engine/internal/models"
)

// GenerateScript compiles the JSON pipeline into a JS function
func GenerateScript(wfDef models.WorkflowDef) string {
	slog.Info("Generating Script for Workflow", "nimb_id", wfDef.NIMB_ID)
	var steps []models.PipelineStep = wfDef.Pipeline
	var buffer bytes.Buffer

	// 1. Inject Helper Functions (The "Standard Library" of our engine)
	buffer.WriteString(`
		// --- Helper Functions ---
		var log = [];

		function addLog(message) {
			log.push({timestamp: new Date().toISOString(), message: message});
		}	
	`)

	// 2. Start Workflow Execution
	buffer.WriteString("\n// --- Workflow Execution ---\n")
	buffer.WriteString("var data = {};\n") // Root data object
	buffer.WriteString("var log = [];\n")
	//3. Initialize Metadata Variables
	if len(wfDef.Metadata) > 0 {
		for _, v := range wfDef.Metadata {
			valStr, _ := json.Marshal(v.Value)
			buffer.WriteString(fmt.Sprintf("%s = %s;\n", v.Key, string(valStr)))
		}
	}
	// 4. Compile Pipeline Steps
	buffer.WriteString(compileSteps(steps, 0, true))
	slog.Info("Compiled Script", "script", buffer.String())
	return buffer.String()
}

func compileSteps(steps []models.PipelineStep, indent int, isRoot bool) string {
	var out bytes.Buffer
	pad := strings.Repeat("  ", indent)

	for _, step := range steps {
		out.WriteString("\n")
		slog.Info("Compiling Step", "type", step.Type, "statement", step.Statement, "target", step.Target)
		switch step.Type {
		case "assignment":
			assignmentStep(&out, step, pad, isRoot)
		case "push_array":
			pushArrayStep(&out, step, pad, isRoot)
		case "condition":
			conditionStep(&out, step, pad, indent, isRoot)
		case "network_call":
			networkCallStep(&out, step, pad, isRoot)
		case "for_each":
			forEachStep(&out, step, pad, indent, isRoot)
		default:
			slog.Warn("Unknown step type", "type", step.Type)
		}
	}
	return out.String()
}

func assignmentStep(out *bytes.Buffer, step models.PipelineStep, pad string, isRoot bool) {
	// Generate:
	// We strip "data." from the path for the helper, as we pass 'data' as root
	if isRoot {
		step.Target = "data." + step.Target
	}
	valStr := string(step.Value)
	//addLog(information about assignment)
	out.WriteString(fmt.Sprintf("\n%s addLog('Assigning %s to %s');\n", pad, valStr, step.Target))
	out.WriteString(fmt.Sprintf("%s %s = %s;", pad, step.Target, valStr))
}

func pushArrayStep(out *bytes.Buffer, step models.PipelineStep, pad string, isRoot bool) {
	// Generate: $push(data, "Transactions", {...});
	if isRoot {
		step.Target = "data." + step.Target
	}
	valStr := string(step.Value)
	out.WriteString(fmt.Sprintf("\n%s addLog('Pushing %s to %s');\n", pad, valStr, step.Target))
	out.WriteString(fmt.Sprintf("%s %s.push(JSON.parse(%s));", pad, step.Target, valStr))
}

func conditionStep(out *bytes.Buffer, step models.PipelineStep, pad string, indent int, isRoot bool) {
	// Generate: if (data.Age >= 21) { ... }
	if isRoot {
		step.Target = "data." + step.Target
	}
	out.WriteString(fmt.Sprintf("\n%s addLog('Evaluating condition: %s');\n", pad, step.Statement))
	// Eg Statment amount > 50 && discount < 20
	// need to convert data.amount > 50 && data.discount < 20
	if isRoot {
		step.Statement = strings.ReplaceAll(step.Statement, "data.", "")
	} else {
		step.Statement = strings.ReplaceAll(step.Statement, step.Target+".", "")
	}
	out.WriteString(fmt.Sprintf("%sif (%s) {", pad, step.Statement))
	// addLog is condition matching
	contextVar := step.ContextVar // Add ContextVar to PipelineStep struct
	if contextVar != "" {
		out.WriteString(fmt.Sprintf("\n%s  var %s = %s;", pad, contextVar, step.Target))

		for _, childStep := range step.Children {
			childStep.Target = strings.ReplaceAll(childStep.Target, strings.Split(childStep.Target, ".")[0]+".", contextVar+".")
		}
	}
	out.WriteString(fmt.Sprintf("\n%s addLog('Condition %s evaluated to true');\n", pad, step.Statement))
	out.WriteString(compileSteps(step.Children, indent+1, false))
	out.WriteString(fmt.Sprintf("\n%s}", pad))
}

func networkCallStep(out *bytes.Buffer, step models.PipelineStep, pad string, isRoot bool) {
	// Generate: var result = $http("GET", "url");
	out.WriteString(fmt.Sprintf("%svar %s = $http('%s', '%s');",
		pad, step.ContextVar, step.Method, step.URL))
}

func forEachStep(out *bytes.Buffer, step models.PipelineStep, pad string, indent int, isRoot bool) {
	indexVar := fmt.Sprintf("i%d", indent)
	var arrayPath string
	if isRoot {
		arrayPath = "data." + step.Target
	}
	out.WriteString(fmt.Sprintf("\n%s addLog('Iterating over array: %s');\n", pad, arrayPath))
	contextVar := step.ContextVar // Add ContextVar to PipelineStep struct
	out.WriteString(fmt.Sprintf("%sfor(var %s=0; %s<%s.length; %s++){", pad, indexVar, indexVar, arrayPath, indexVar))
	if contextVar != "" {
		out.WriteString(fmt.Sprintf("\n%s  var %s = %s[%s];", pad, contextVar, arrayPath, indexVar))

		for index, childStep := range step.Children {
			step.Children[index].Target = contextVar + "." + childStep.Target
		}
	}
	out.WriteString(compileSteps(step.Children, indent+1, false))
	out.WriteString(fmt.Sprintf("\n%s}", pad))
}

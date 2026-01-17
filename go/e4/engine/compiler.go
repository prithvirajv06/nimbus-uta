package engine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"text/template"

	"github.com/prithvirajv06/nimbus-uta/go/engine/internal/models"
)

// Templates for JS code generation
var (
	assignTmpl = template.Must(template.New("assign").Parse(`
{{.Pad}}addLog("Assigning {{.ValStr}} to {{.TargetStr}}");
{{.Pad}}{{.TargetStr}} = {{.ValStr}};
`))
	pushArrayTmpl = template.Must(template.New("pushArray").Parse(`
{{.Pad}}addLog("Pushing {{.ValStr}} to {{.TargetStr}}");
{{.Pad}}{{.TargetStr}}.push(JSON.parse({{.ValStr}}));
`))
	conditionTmpl = template.Must(template.New("condition").Parse(`
{{.Pad}}addLog("Evaluating condition: {{.Statement}}");
{{.Pad}}if ({{.Statement}}) {
{{.TrueBlock}}
{{.Pad}}  addLog("Condition {{.Statement}} evaluated to true");
{{.Pad}}}
{{if .HasFalse}}
{{.Pad}}else {
{{.FalseBlock}}
{{.Pad}}  addLog("Condition {{.Statement}} evaluated to false");
{{.Pad}}}
{{end}}
`))
	networkCallTmpl = template.Must(template.New("networkCall").Parse(`
{{.Pad}}var {{.ContextVar}} = $http('{{.Method}}', '{{.URL}}');
`))
	forEachTmpl = template.Must(template.New("forEach").Parse(`
{{.Pad}}addLog("Iterating over array: {{.ArrayPath}}");
{{.Pad}}for(var {{.IndexVar}}=0; {{.IndexVar}}<{{.ArrayPath}}.length; {{.IndexVar}}++){
{{.Pad}}  var {{.ContextVar}} = {{.ArrayPath}}[{{.IndexVar}}];
{{.ChildrenBlock}}
{{.Pad}}}
`))

	toUpperCaseTmpl = template.Must(template.New("toUpperCase").Parse(`
{{.Pad}}addLog("Converting {{.ValStr}} to uppercase and assigning to {{.TargetStr}}");
{{.Pad}}{{.TargetStr}} = {{.ValStr}}.toUpperCase();
`))

	toLowerCaseTmpl = template.Must(template.New("toLowerCase").Parse(`
{{.Pad}}addLog("Converting {{.ValStr}} to lowercase and assigning to {{.TargetStr}}");
{{.Pad}}{{.TargetStr}} = {{.ValStr}}.toLowerCase();
`))

	increementTmpl = template.Must(template.New("increment").Parse(`
{{.Pad}}addLog("Incrementing {{.TargetStr}} by {{.ValStr}}");
{{.Pad}}{{.TargetStr}} += {{.ValStr}};
`))

	decreementTmpl = template.Must(template.New("decrement").Parse(`
{{.Pad}}addLog("Decrementing {{.TargetStr}} by {{.ValStr}}");
{{.Pad}}{{.TargetStr}} -= {{.ValStr}};
`))

	appendPrexixTmpl = template.Must(template.New("appendPrefix").Parse(`
{{.Pad}}addLog("Appending prefix {{.ValStr}} to {{.TargetStr}}");
{{.Pad}}{{.TargetStr}} = {{.ValStr}} + {{.TargetStr}};
`))

	appendSuffixTmpl = template.Must(template.New("appendSuffix").Parse(`
{{.Pad}}addLog("Appending suffix {{.ValStr}} to {{.TargetStr}}");
{{.Pad}}{{.TargetStr}} = {{.TargetStr}} + {{.ValStr}};
`))

	multiplyTmpl = template.Must(template.New("multiply").Parse(`
{{.Pad}}addLog("Multiplying {{.TargetStr}} by {{.ValStr}}");
{{.Pad}}{{.TargetStr}} *= {{.ValStr}};
`))

	divideTmpl = template.Must(template.New("divide").Parse(`
{{.Pad}}addLog("Dividing {{.TargetStr}} by {{.ValStr}}");
{{.Pad}}{{.TargetStr}} /= {{.ValStr}};
`))
)

// Step handler registry
type stepHandler func(*bytes.Buffer, models.WorkflowStep, string, int, bool)

var stepRegistry map[string]stepHandler

// GenerateScript compiles the JSON pipeline into a JS function
func GenerateScript(wfDef models.LogicFlow) string {
	slog.Info("Generating Script for Workflow", "nimb_id", wfDef.NIMB_ID)
	var buffer bytes.Buffer

	// 1. Inject Helper Functions
	buffer.WriteString(`
        // --- Helper Functions ---
        var log = [];

        function addLog(message) {
            log.push({timestamp: new Date().toISOString(), message: message});
        }	
    `)

	// 2. Start Workflow Execution
	buffer.WriteString("\n// --- Workflow Execution ---\n")
	buffer.WriteString("function executeWorkflow(inputData) {\n")
	buffer.WriteString("  var data = inputData;\n")

	// 3. Initialize Metadata Variables
	if len(wfDef.Metadata) > 0 {
		for _, v := range wfDef.Metadata {
			valStr, err := json.Marshal(v.Value)
			if err != nil {
				slog.Warn("Failed to marshal metadata value", "varKey", v.VarKey, "error", err)
				continue
			}
			buffer.WriteString(fmt.Sprintf("  %s = %s;\n", v.VarKey, string(valStr)))
		}
	}

	// 4. Compile Pipeline Steps
	buffer.WriteString(compileSteps(wfDef.LogicalSteps, 1, true))
	buffer.WriteString("\n  return {data: data, log: log};\n")
	buffer.WriteString("}\n")
	slog.Info("Compiled Script", "script", buffer.String())
	return buffer.String()
}

func compileSteps(steps []models.WorkflowStep, indent int, isRoot bool) string {
	var out bytes.Buffer
	pad := strings.Repeat("  ", indent)
	for _, step := range steps {
		out.WriteString("\n")
		slog.Info("Compiling Step", "type", step.Type, "statement", step.Statement, "target", step.Target)
		if handler, ok := stepRegistry[step.Type]; ok {
			handler(&out, step, pad, indent, isRoot)
		} else {
			slog.Warn("Unknown step type", "type", step.Type)
			if (step.Type == "start" || step.Type == "end") && len(step.Children) > 0 {
				out.WriteString(compileSteps(step.Children, indent, isRoot))
			}
		}
	}
	return out.String()
}

func assignmentStep(out *bytes.Buffer, step models.WorkflowStep, pad string, indent int, isRoot bool) {
	targetStr := step.Target.VarKey
	valStr := string(step.Value)
	for _, mapping := range step.ContextMap {
		targetStr = strings.ReplaceAll(targetStr, mapping.VarKey, mapping.ContextKey)
		valStr = strings.ReplaceAll(valStr, mapping.VarKey, mapping.ContextKey)
	}

	if step.Target.Type == "string" && (!strings.Contains(valStr, "?") && !strings.Contains(valStr, ":")) {
		valStr = fmt.Sprintf("'%s'", strings.ReplaceAll(valStr, "'", "\\'"))
	}
	assignTmpl.Execute(out, map[string]string{
		"Pad":       pad,
		"TargetStr": targetStr,
		"ValStr":    valStr,
	})
}

func pushArrayStep(out *bytes.Buffer, step models.WorkflowStep, pad string, indent int, isRoot bool) {
	targetStr := step.Target.VarKey
	valStr := string(step.Value)
	for _, mapping := range step.ContextMap {
		targetStr = strings.ReplaceAll(targetStr, mapping.VarKey, mapping.ContextKey)
	}
	pushArrayTmpl.Execute(out, map[string]string{
		"Pad":       pad,
		"TargetStr": targetStr,
		"ValStr":    valStr,
	})
}

func conditionStep(out *bytes.Buffer, step models.WorkflowStep, pad string, indent int, isRoot bool) {
	statement := step.Statement
	for _, mapping := range step.ContextMap {
		statement = strings.ReplaceAll(statement, mapping.VarKey, mapping.ContextKey)
	}
	trueBlock := compileSteps(step.TrueChildren, indent+1, false)
	falseBlock := ""
	hasFalse := len(step.FalseChildren) > 0
	if hasFalse {
		falseBlock = compileSteps(step.FalseChildren, indent+1, false)
	}
	conditionTmpl.Execute(out, map[string]interface{}{
		"Pad":        pad,
		"Statement":  statement,
		"TrueBlock":  trueBlock,
		"HasFalse":   hasFalse,
		"FalseBlock": falseBlock,
	})
}

func networkCallStep(out *bytes.Buffer, step models.WorkflowStep, pad string, indent int, isRoot bool) {
	networkCallTmpl.Execute(out, map[string]string{
		"Pad":        pad,
		"ContextVar": step.ContextVar,
		"Method":     step.Type,
		"URL":        string(step.Value),
	})
}

func forEachStep(out *bytes.Buffer, step models.WorkflowStep, pad string, indent int, isRoot bool) {
	indexVar := fmt.Sprintf("i%d", indent)
	arrayPath := step.Target.VarKey
	for _, mapping := range step.ContextMap {
		arrayPath = strings.ReplaceAll(arrayPath, mapping.VarKey, mapping.ContextKey)
	}
	childrenBlock := compileSteps(step.Children, indent+2, false)
	forEachTmpl.Execute(out, map[string]string{
		"Pad":           pad,
		"ArrayPath":     arrayPath,
		"IndexVar":      indexVar,
		"ContextVar":    step.ContextVar,
		"ChildrenBlock": childrenBlock,
	})
}

func templateStepHandler(tmpl *template.Template, extraFields map[string]string) stepHandler {
	return func(out *bytes.Buffer, step models.WorkflowStep, pad string, indent int, isRoot bool) {
		targetStr := step.Target.VarKey
		valStr := string(step.Value)
		for _, mapping := range step.ContextMap {
			targetStr = strings.ReplaceAll(targetStr, mapping.VarKey, mapping.ContextKey)
		}
		data := map[string]string{
			"Pad":       pad,
			"TargetStr": targetStr,
			"ValStr":    valStr,
		}
		// Merge extra fields if any
		for k, v := range extraFields {
			data[k] = v
		}
		tmpl.Execute(out, data)
	}
}

// Initialize stepRegistry after all step handler functions are defined
func init() {
	stepRegistry = map[string]stepHandler{
		"condition":    conditionStep,
		"network_call": networkCallStep,
		"for_each":     forEachStep,
		"assignment":   assignmentStep,
		// Generic template-based handlers (One template per step type)
		"push_array":    templateStepHandler(pushArrayTmpl, nil),
		"to_uppercase":  templateStepHandler(toUpperCaseTmpl, nil),
		"to_lowercase":  templateStepHandler(toLowerCaseTmpl, nil),
		"increment":     templateStepHandler(increementTmpl, nil),
		"decrement":     templateStepHandler(decreementTmpl, nil),
		"append_prefix": templateStepHandler(appendPrexixTmpl, nil),
		"append_suffix": templateStepHandler(appendSuffixTmpl, nil),
		"multiply":      templateStepHandler(multiplyTmpl, nil),
		"divide":        templateStepHandler(divideTmpl, nil),
	}
}

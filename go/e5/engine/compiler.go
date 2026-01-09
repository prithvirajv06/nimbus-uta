package engine

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/prithvirajv06/nimbus-uta/go/engine/internal/models"
)

// GenerateScript compiles the JSON pipeline into a JS function
func GenerateScript(steps []models.PipelineStep) string {
	var buffer bytes.Buffer

	// 1. Inject Helper Functions (The "Standard Library" of our engine)
	buffer.WriteString(`
		// --- Helper: Safe Deep Set ---
		function $set(root, path, value) {
			var parts = path.split('.');
			var curr = root;
			for (var i = 0; i < parts.length - 1; i++) {
				var key = parts[i];
				if (!curr[key]) curr[key] = {}; // Auto-create path
				curr = curr[key];
			}
			curr[parts[parts.length - 1]] = value;
		}

		// --- Helper: Array Push ---
		function $push(root, path, value) {
			var parts = path.split('.');
			var curr = root;
			for (var i = 0; i < parts.length - 1; i++) {
				curr = curr[parts[i]];
			}
			var target = curr[parts[parts.length - 1]];
			if (!Array.isArray(target)) {
				// Initialize if missing
				target = [];
				curr[parts[parts.length - 1]] = target;
			}
			target.push(value);
		}
	`)

	// 2. Compile Body
	buffer.WriteString("\n// --- Workflow Execution ---\n")
	buffer.WriteString(compileSteps(steps, 0))

	return buffer.String()
}

func compileSteps(steps []models.PipelineStep, indent int) string {
	var out bytes.Buffer
	pad := strings.Repeat("  ", indent)

	for _, step := range steps {
		out.WriteString("\n")
		switch step.Type {

		case "assignment":
			// Generate: $set(data, "Customer.Name", "Alice");
			// We strip "data." from the path for the helper, as we pass 'data' as root
			cleanPath := strings.TrimPrefix(step.Target, "data.")
			valStr := string(step.Value)
			out.WriteString(fmt.Sprintf("%s$set(data, '%s', %s);", pad, cleanPath, valStr))

		case "push_array":
			// Generate: $push(data, "Transactions", {...});
			cleanPath := strings.TrimPrefix(step.Target, "data.")
			valStr := string(step.Value)
			out.WriteString(fmt.Sprintf("%s$push(data, '%s', %s);", pad, cleanPath, valStr))

		case "condition":
			// Generate: if (data.Age >= 21) { ... }
			out.WriteString(fmt.Sprintf("%sif (%s) {", pad, step.Statement))
			out.WriteString(compileSteps(step.Children, indent+1))
			out.WriteString(fmt.Sprintf("\n%s}", pad))

		case "network_call":
			// Generate: var fraudResult = $http("POST", "url");
			out.WriteString(fmt.Sprintf("%svar %s = $http('%s', '%s');",
				pad, step.ResultVar, step.Method, step.URL))
		case "for_each":
			indexVar := fmt.Sprintf("i%d", indent)
			arrayPath := strings.TrimPrefix(step.Target, "data.")
			contextVar := step.ContextVar // Add ContextVar to PipelineStep struct
			out.WriteString(fmt.Sprintf("%sfor(var %s=0; %s<%s.length; %s++){", pad, indexVar, indexVar, arrayPath, indexVar))
			if contextVar != "" {
				out.WriteString(fmt.Sprintf("\n%s  var %s = %s[%s];", pad, contextVar, arrayPath, indexVar))
			}
			out.WriteString(compileSteps(step.Children, indent+1))
			out.WriteString(fmt.Sprintf("\n%s}", pad))
		}
	}
	return out.String()
}

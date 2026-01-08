package engine

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"text/template"

	"github.com/prithvirajv06/nimbus-uta/go/engine/internal/models"
)

// --- 1. IR (Intermediate Representation) Interfaces ---

type Compiler struct {
	Structs []string
}

type VariableMeta struct {
	VarKey   string         `json:"var_key"`
	Type     string         `json:"type"`
	Children []VariableMeta `json:"children"`
}

// Conditional Logic
type IfStep struct {
	Condition string
	Body      []Step
}

// Looping Logic
type LoopStep struct {
	Slice    string
	Iterator string
	Body     []Step
}

// API/Network Request
type APIReqStep struct {
	ResultVar string
	URL       string
	Method    string
}

// Assignment Logic
type AssignStep struct {
	Target string
	Value  string
}

type Step interface {
	GenerateGo(indent string) string
}

type RuleService struct {
	NIMB_ID          string         `bson:"nimb_id" json:"nimb_id"`
	Name             string         `json:"engine_name"`
	Description      string         `json:"description"`
	VariableMetadata VariableMeta   `json:"variable_metadata"`
	Pipeline         []PipelineStep `json:"pipeline"`
	Port             int            `json:"port"`
	Audit            models.Audit   `json:"audit_fields"`
}

// PipelineStep matches the incoming JSON structure for logic
type PipelineStep struct {
	StepType  string         `json:"step_type"`
	Target    string         `json:"target,omitempty"`
	Iterator  string         `json:"iterator,omitempty"`
	Statement string         `json:"statement,omitempty"`
	Value     string         `json:"value,omitempty"`
	ResultVar string         `json:"result_var,omitempty"`
	URL       string         `json:"url,omitempty"`
	Method    string         `json:"method,omitempty"`
	Children  []PipelineStep `json:"children,omitempty"`
}

var serviceRegistry = make(map[string]string)

func RegisterService(nimbID, endpoint string) {
	serviceRegistry[nimbID] = endpoint
}

func GetServiceEndpoint(nimbID string) (string, bool) {
	endpoint, ok := serviceRegistry[nimbID]
	return endpoint, ok
}

func buildAndRunDockerService(engineName string, port int) error {
	imageName := fmt.Sprintf("%s-engine", strings.ToLower(engineName))
	containerName := fmt.Sprintf("%s-container", strings.ToLower(engineName))
	fileName := fmt.Sprintf("services/%s_main.go", strings.ToLower(engineName))

	// Build the binary
	buildCmd := exec.Command("go", "build", "-o", "services/main", fileName)
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		return err
	}

	// Build Docker image
	dockerBuild := exec.Command("docker", "build", "-t", imageName, "-f", "Dockerfile", ".")
	dockerBuild.Stdout = os.Stdout
	dockerBuild.Stderr = os.Stderr
	if err := dockerBuild.Run(); err != nil {
		return err
	}

	// Run Docker container
	dockerRun := exec.Command("docker", "run", "-d", "--name", containerName, "-p", fmt.Sprintf("%d:8081", port), imageName)
	dockerRun.Stdout = os.Stdout
	dockerRun.Stderr = os.Stderr
	return dockerRun.Run()
}

// ...existing code...
func killProcessOnPort(port int) error {
	// Find the PID using netstat
	cmd := exec.Command("netstat", "-ano")
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	re := regexp.MustCompile(fmt.Sprintf(`(?m)^.*:%d\s+.*\s+LISTENING\s+(\d+)$`, port))
	matches := re.FindStringSubmatch(string(out))
	if len(matches) < 2 {
		return nil // No process found
	}
	pid := matches[1]
	// Kill the process
	killCmd := exec.Command("taskkill", "/PID", pid, "/F")
	return killCmd.Run()
}

func (c *Compiler) CompileFromRequest(req RuleService) (bool, error) {

	// 1. Generate the Data Contracts (Structs)
	rootType := c.BuildType(req.VariableMetadata)

	// 2. Generate the Logic Blocks (Pipeline)
	pipelineIR := ParsePipeline(req.Pipeline)
	var logicBuf bytes.Buffer
	for _, s := range pipelineIR {
		logicBuf.WriteString(s.GenerateGo("\t"))
	}

	// 3. Finalize Source Code
	source := c.ExecuteTemplate(req.Name, rootType, logicBuf.String(), req.Port)
	// 4. Save/Deploy
	fileName := fmt.Sprintf("services/%s_main.go", strings.ToLower(req.Name))
	os.WriteFile(fileName, []byte(source), 0644)

	// 5. Build and run as Docker container
	port := 8081 // You may want to assign unique ports per service
	if err := buildAndRunDockerService(req.Name, port); err != nil {
		return false, err
	}

	// 6. Register service
	RegisterService(req.NIMB_ID, fmt.Sprintf("http://localhost:%d", port))
	return true, nil
}

// ...existing code...

func (s *AssignStep) GenerateGo(in string) string {
	escapedValue := strings.ReplaceAll(s.Value, "\"", "\\\"")
	return fmt.Sprintf("%slogger.Add(\"Assign: %s = %s\")\n%s%s = %s\n", in, s.Target, escapedValue, in, s.Target, s.Value)
}

func (s *APIReqStep) GenerateGo(in string) string {
	escapedURL := strings.ReplaceAll(s.URL, "\"", "\\\"")
	return fmt.Sprintf("%s%s, err := CallExternal(\"%s\", \"%s\", nil, logger)\n%sif err != nil { logger.Add(\"Error in network call:\", %s); return }\n", in, s.ResultVar, s.Method, escapedURL, in, s.ResultVar)
}

func (s *LoopStep) GenerateGo(in string) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%slogger.Add(\"Loop over: %s\")\n", in, s.Slice))
	buf.WriteString(fmt.Sprintf("%sfor i := range %s {\n", in, s.Slice))
	buf.WriteString(fmt.Sprintf("%s\t%s := &%s[i]\n", in, s.Iterator, s.Slice))
	for _, child := range s.Body {
		buf.WriteString(child.GenerateGo(in + "\t"))
	}
	buf.WriteString(in + "}\n")
	return buf.String()
}

func (s *IfStep) GenerateGo(in string) string {
	var buf bytes.Buffer
	escapedCond := strings.ReplaceAll(s.Condition, "\"", "\\\"")
	buf.WriteString(fmt.Sprintf("%slogger.Add(\"Condition: %s\")\n", in, escapedCond))
	buf.WriteString(fmt.Sprintf("%sif %s {\n", in, s.Condition))
	for _, child := range s.Body {
		buf.WriteString(child.GenerateGo(in + "\t"))
	}
	buf.WriteString(in + "} else {\n")
	buf.WriteString(fmt.Sprintf("%s\tlogger.Add(\"Condition failed: %s\")\n", in, escapedCond))
	buf.WriteString(in + "}\n")
	return buf.String()
}

// --- 2. Type Generator (Metadata to Structs) ---

func (c *Compiler) BuildType(v VariableMeta) string {
	raw := v.VarKey
	if strings.Contains(raw, ".") {
		parts := strings.Split(raw, ".")
		raw = parts[len(parts)-1]
	}
	name := strings.Title(strings.ReplaceAll(raw, "[*]", ""))

	if v.Type == "object" || v.Type == "array" {
		var fields []string
		for _, child := range v.Children {
			childType := c.BuildType(child)
			fName := strings.Title(strings.Split(child.VarKey, ".")[len(strings.Split(child.VarKey, "."))-1])
			fName = strings.ReplaceAll(fName, "[*]", "")
			fields = append(fields, fmt.Sprintf("\t%s %s `json:\"%s\"`", fName, childType, strings.ToLower(fName)))
		}
		c.Structs = append(c.Structs, fmt.Sprintf("type %s struct {\n%s\n}", name, strings.Join(fields, "\n")))
		if v.Type == "array" {
			return "[]" + name
		}
		return name
	}
	switch v.Type {
	case "number":
		return "float64"
	case "boolean":
		return "bool"
	default:
		return "string"
	}
}

// ExecuteTemplate: The missing link that generates the full file
func (c *Compiler) ExecuteTemplate(engineName, rootType, logicBody string, port int) string {

	tmpl := template.Must(template.New("svc").Parse(serviceTmpl))
	var out bytes.Buffer
	tmpl.Execute(&out, map[string]interface{}{
		"Types":      c.Structs,
		"RootType":   rootType,
		"Logic":      logicBody,
		"EngineName": engineName,
		"Port":       port,
	})
	return out.String()
}

// --- 3. Main Service Template ---

const serviceTmpl = `package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "io"
)

type Logger struct {
    Logs []string ` + "`json:\"logs\"`" + `
}

func (l *Logger) Add(msg string) {
    l.Logs = append(l.Logs, msg)
}

// --- DATA CONTRACTS ---
{{range .Types}}
{{.}}
{{end}}

type RuleService struct{}

// Helper for Network Calls
func CallExternal(url, method string, payload interface{}, logger *Logger) (map[string]interface{}, error) {
    logger.Add(fmt.Sprintf("Calling external: %s %s", method, url))
    // Standard library implementation
    return map[string]interface{}{"status": "success"}, nil
}

func (rs *RuleService) Process(w http.ResponseWriter, r *http.Request) {
    var data {{.RootType}}
    logger := &Logger{}
    body, _ := io.ReadAll(r.Body)
    if err := json.Unmarshal(body, &data); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // --- GENERATED LOGIC ---
{{.Logic}}

    if r.Header.Get("nimbus-debug") == "true" {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]interface{}{
            "log": logger.Logs,
            "response": data,
        })
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}

func main() {
    rs := &RuleService{}
    http.HandleFunc("/", rs.Process)
    fmt.Println("Rule Service starting on :{{.Port}}")
    http.ListenAndServe(":{{.Port}}", nil)
}
`

// ParsePipeline recursively converts JSON steps into Go IR Steps
func ParsePipeline(jsonSteps []PipelineStep) []Step {
	var steps []Step

	for _, js := range jsonSteps {
		switch js.StepType {
		case "assignment":
			steps = append(steps, &AssignStep{Target: js.Target, Value: js.Value})

		case "network_call":
			steps = append(steps, &APIReqStep{ResultVar: js.ResultVar, URL: js.URL, Method: js.Method})

		case "loop":
			steps = append(steps, &LoopStep{
				Slice:    js.Target,
				Iterator: js.Iterator,
				Body:     ParsePipeline(js.Children), // Recursive
			})

		case "condition":
			steps = append(steps, &IfStep{
				Condition: js.Statement,
				Body:      ParsePipeline(js.Children), // Recursive
			})
		}
	}
	return steps
}

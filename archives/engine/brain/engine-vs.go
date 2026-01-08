// Package ruleengine provides a flexible, high-performance rule engine
// for processing JSON data with complex rules and transformations.
package ruleengine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

// Common errors
var (
	ErrInvalidRule      = errors.New("invalid rule definition")
	ErrInvalidCondition = errors.New("invalid condition")
	ErrInvalidAction    = errors.New("invalid action")
	ErrPathNotFound     = errors.New("path not found")
	ErrExecutionTimeout = errors.New("execution timeout")
)

// Rule represents a business rule with conditions and actions
type Rule struct {
	ID                string     `json:"id"`
	Name              string     `json:"name"`
	Description       string     `json:"description"`
	Condition         *Condition `json:"condition"`
	Actions           []Action   `json:"actions"`
	Priority          int        `json:"priority"`
	Enabled           bool       `json:"enabled"`
	ContinueOnSuccess bool       `json:"continue_on_success"`
	StopOnError       bool       `json:"stop_on_error"`
	Tags              []string   `json:"tags,omitempty"`
}

// Condition represents a logical condition for rule evaluation
type Condition struct {
	Type     string       `json:"type"` // "and", "or", "not", "expr"
	Children []*Condition `json:"children,omitempty"`
	Path     string       `json:"path,omitempty"`
	Operator string       `json:"operator,omitempty"`
	Value    interface{}  `json:"value,omitempty"`
}

// Action represents an operation to perform on the data
type Action struct {
	Type      string      `json:"type"` // "set", "delete", "copy", "append", "transform"
	Path      string      `json:"path"`
	Value     interface{} `json:"value,omitempty"`
	Source    string      `json:"source,omitempty"`
	Function  string      `json:"function,omitempty"`
	Arguments []string    `json:"arguments,omitempty"`
}

// ExecutionResult contains the results of rule execution
type ExecutionResult struct {
	Output       map[string]interface{} `json:"output"`
	RulesFired   []string               `json:"rules_fired"`
	RulesFailed  []string               `json:"rules_failed"`
	ExecutionLog []ExecutionLogEntry    `json:"execution_log"`
	Duration     time.Duration          `json:"duration"`
	Error        error                  `json:"error,omitempty"`
}

// ExecutionLogEntry represents a single log entry during execution
type ExecutionLogEntry struct {
	Timestamp time.Time   `json:"timestamp"`
	RuleID    string      `json:"rule_id"`
	RuleName  string      `json:"rule_name"`
	Event     string      `json:"event"`
	Details   interface{} `json:"details,omitempty"`
}

// Engine is the main rule engine
type Engine struct {
	rules            []*Rule
	functionRegistry *FunctionRegistry
	operatorRegistry *OperatorRegistry
	maxConcurrency   int
	enableLogging    bool
	timeout          time.Duration
	mu               sync.RWMutex
	compiledPatterns map[string]*regexp.Regexp
	contextPool      *sync.Pool
}

// EngineOption is a functional option for configuring the engine
type EngineOption func(*Engine)

// NewEngine creates a new rule engine instance
func NewEngine(opts ...EngineOption) *Engine {
	e := &Engine{
		rules:            make([]*Rule, 0),
		functionRegistry: NewFunctionRegistry(),
		operatorRegistry: NewOperatorRegistry(),
		maxConcurrency:   10,
		enableLogging:    true,
		timeout:          30 * time.Second,
		compiledPatterns: make(map[string]*regexp.Regexp),
		contextPool: &sync.Pool{
			New: func() interface{} {
				return &ExecutionContext{
					Variables: make(map[string]interface{}),
					Log:       make([]ExecutionLogEntry, 0),
				}
			},
		},
	}

	for _, opt := range opts {
		opt(e)
	}

	e.registerDefaultFunctions()
	e.registerDefaultOperators()

	return e
}

// WithConcurrency sets the maximum concurrency level
func WithConcurrency(n int) EngineOption {
	return func(e *Engine) {
		e.maxConcurrency = n
	}
}

// WithLogging enables or disables execution logging
func WithLogging(enabled bool) EngineOption {
	return func(e *Engine) {
		e.enableLogging = enabled
	}
}

// WithTimeout sets the execution timeout
func WithTimeout(d time.Duration) EngineOption {
	return func(e *Engine) {
		e.timeout = d
	}
}

// ExecutionContext holds the state during rule execution
type ExecutionContext struct {
	Data      map[string]interface{}
	Variables map[string]interface{}
	Log       []ExecutionLogEntry
	mu        sync.Mutex
}

// AddLog adds a log entry to the execution context
func (ctx *ExecutionContext) AddLog(ruleID, ruleName, event string, details interface{}) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.Log = append(ctx.Log, ExecutionLogEntry{
		Timestamp: time.Now(),
		RuleID:    ruleID,
		RuleName:  ruleName,
		Event:     event,
		Details:   details,
	})
}

// LoadRules loads rules into the engine
func (e *Engine) LoadRules(rules []*Rule) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, rule := range rules {
		if err := e.validateRule(rule); err != nil {
			return fmt.Errorf("invalid rule %s: %w", rule.ID, err)
		}
	}

	e.rules = rules
	e.sortRulesByPriority()
	return nil
}

// LoadRulesFromJSON loads rules from JSON bytes
func (e *Engine) LoadRulesFromJSON(data []byte) error {
	var rules []*Rule
	if err := json.Unmarshal(data, &rules); err != nil {
		return fmt.Errorf("failed to unmarshal rules: %w", err)
	}
	return e.LoadRules(rules)
}

// Execute executes rules against the input data
func (e *Engine) Execute(ctx context.Context, input map[string]interface{}) (*ExecutionResult, error) {
	start := time.Now()

	timeoutCtx, cancel := context.WithTimeout(ctx, e.timeout)
	defer cancel()

	execCtx := e.contextPool.Get().(*ExecutionContext)
	defer e.contextPool.Put(execCtx)

	execCtx.Data = deepCopy(input)
	execCtx.Variables = make(map[string]interface{})
	execCtx.Log = execCtx.Log[:0]

	result := &ExecutionResult{
		Output:       execCtx.Data,
		RulesFired:   make([]string, 0),
		RulesFailed:  make([]string, 0),
		ExecutionLog: make([]ExecutionLogEntry, 0),
	}

	e.mu.RLock()
	rules := e.rules
	e.mu.RUnlock()

	for _, rule := range rules {
		if !rule.Enabled {
			continue
		}

		select {
		case <-timeoutCtx.Done():
			result.Error = ErrExecutionTimeout
			result.Duration = time.Since(start)
			return result, ErrExecutionTimeout
		default:
		}

		if err := e.executeRule(execCtx, rule); err != nil {
			result.RulesFailed = append(result.RulesFailed, rule.ID)
			if e.enableLogging {
				execCtx.AddLog(rule.ID, rule.Name, "error", err.Error())
			}
			if rule.StopOnError {
				result.Error = err
				break
			}
		} else {
			result.RulesFired = append(result.RulesFired, rule.ID)
			if !rule.ContinueOnSuccess {
				break
			}
		}
	}

	result.ExecutionLog = execCtx.Log
	result.Duration = time.Since(start)
	return result, nil
}

// executeRule executes a single rule
func (e *Engine) executeRule(ctx *ExecutionContext, rule *Rule) error {
	if e.enableLogging {
		ctx.AddLog(rule.ID, rule.Name, "evaluating", nil)
	}

	matched, err := e.evaluateCondition(ctx, rule.Condition)
	if err != nil {
		return fmt.Errorf("condition evaluation failed: %w", err)
	}

	if !matched {
		if e.enableLogging {
			ctx.AddLog(rule.ID, rule.Name, "condition_not_met", nil)
		}
		return nil
	}

	if e.enableLogging {
		ctx.AddLog(rule.ID, rule.Name, "condition_met", nil)
	}

	for i, action := range rule.Actions {
		if err := e.executeAction(ctx, &action); err != nil {
			return fmt.Errorf("action %d failed: %w", i, err)
		}
	}

	if e.enableLogging {
		ctx.AddLog(rule.ID, rule.Name, "executed", nil)
	}

	return nil
}

// evaluateCondition evaluates a condition tree
func (e *Engine) evaluateCondition(ctx *ExecutionContext, cond *Condition) (bool, error) {
	if cond == nil {
		return true, nil
	}

	switch cond.Type {
	case "and":
		for _, child := range cond.Children {
			result, err := e.evaluateCondition(ctx, child)
			if err != nil {
				return false, err
			}
			if !result {
				return false, nil
			}
		}
		return true, nil

	case "or":
		for _, child := range cond.Children {
			result, err := e.evaluateCondition(ctx, child)
			if err != nil {
				return false, err
			}
			if result {
				return true, nil
			}
		}
		return false, nil

	case "not":
		if len(cond.Children) != 1 {
			return false, ErrInvalidCondition
		}
		result, err := e.evaluateCondition(ctx, cond.Children[0])
		if err != nil {
			return false, err
		}
		return !result, nil

	case "expr":
		value, err := getValueAtPath(ctx.Data, cond.Path)
		if err != nil {
			return false, nil // Path not found = false
		}
		return e.operatorRegistry.Evaluate(cond.Operator, value, cond.Value)

	default:
		return false, fmt.Errorf("%w: unknown type %s", ErrInvalidCondition, cond.Type)
	}
}

// executeAction executes a single action
func (e *Engine) executeAction(ctx *ExecutionContext, action *Action) error {
	switch action.Type {
	case "set":
		return setValueAtPath(ctx.Data, action.Path, action.Value)

	case "delete":
		return deleteValueAtPath(ctx.Data, action.Path)

	case "copy":
		value, err := getValueAtPath(ctx.Data, action.Source)
		if err != nil {
			return err
		}
		return setValueAtPath(ctx.Data, action.Path, value)

	case "append":
		existing, err := getValueAtPath(ctx.Data, action.Path)
		if err != nil {
			return setValueAtPath(ctx.Data, action.Path, []interface{}{action.Value})
		}
		arr, ok := existing.([]interface{})
		if !ok {
			return fmt.Errorf("path %s is not an array", action.Path)
		}
		arr = append(arr, action.Value)
		return setValueAtPath(ctx.Data, action.Path, arr)

	case "transform":
		args := make([]interface{}, len(action.Arguments))
		for i, argPath := range action.Arguments {
			val, err := getValueAtPath(ctx.Data, argPath)
			if err != nil {
				return err
			}
			args[i] = val
		}
		result, err := e.functionRegistry.Call(action.Function, args...)
		if err != nil {
			return err
		}
		return setValueAtPath(ctx.Data, action.Path, result)

	default:
		return fmt.Errorf("%w: unknown type %s", ErrInvalidAction, action.Type)
	}
}

// validateRule validates a rule definition
func (e *Engine) validateRule(rule *Rule) error {
	if rule.ID == "" {
		return fmt.Errorf("%w: missing ID", ErrInvalidRule)
	}
	if rule.Condition == nil {
		return fmt.Errorf("%w: missing condition", ErrInvalidRule)
	}
	if len(rule.Actions) == 0 {
		return fmt.Errorf("%w: no actions defined", ErrInvalidRule)
	}
	return nil
}

// sortRulesByPriority sorts rules by priority (higher first)
func (e *Engine) sortRulesByPriority() {
	sort.Slice(e.rules, func(i, j int) bool {
		return e.rules[i].Priority > e.rules[j].Priority
	})
}

// FunctionRegistry manages custom functions
type FunctionRegistry struct {
	functions map[string]func(...interface{}) (interface{}, error)
	mu        sync.RWMutex
}

// NewFunctionRegistry creates a new function registry
func NewFunctionRegistry() *FunctionRegistry {
	return &FunctionRegistry{
		functions: make(map[string]func(...interface{}) (interface{}, error)),
	}
}

// Register registers a new function
func (fr *FunctionRegistry) Register(name string, fn func(...interface{}) (interface{}, error)) {
	fr.mu.Lock()
	defer fr.mu.Unlock()
	fr.functions[name] = fn
}

// Call calls a registered function
func (fr *FunctionRegistry) Call(name string, args ...interface{}) (interface{}, error) {
	fr.mu.RLock()
	fn, exists := fr.functions[name]
	fr.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("function %s not found", name)
	}
	return fn(args...)
}

// OperatorRegistry manages comparison operators
type OperatorRegistry struct {
	operators map[string]func(interface{}, interface{}) (bool, error)
	mu        sync.RWMutex
}

// NewOperatorRegistry creates a new operator registry
func NewOperatorRegistry() *OperatorRegistry {
	return &OperatorRegistry{
		operators: make(map[string]func(interface{}, interface{}) (bool, error)),
	}
}

// Register registers a new operator
func (or *OperatorRegistry) Register(name string, fn func(interface{}, interface{}) (bool, error)) {
	or.mu.Lock()
	defer or.mu.Unlock()
	or.operators[name] = fn
}

// Evaluate evaluates an operator
func (or *OperatorRegistry) Evaluate(op string, left, right interface{}) (bool, error) {
	or.mu.RLock()
	fn, exists := or.operators[op]
	or.mu.RUnlock()

	if !exists {
		return false, fmt.Errorf("operator %s not found", op)
	}
	return fn(left, right)
}

// registerDefaultFunctions registers built-in functions
func (e *Engine) registerDefaultFunctions() {

	e.functionRegistry.Register("len", func(args ...interface{}) (interface{}, error) {
		if len(args) != 1 {
			return nil, errors.New("len requires 1 argument")
		}
		switch v := args[0].(type) {
		case string:
			return len(v), nil
		case []interface{}:
			return len(v), nil
		case map[string]interface{}:
			return len(v), nil
		default:
			return nil, errors.New("len requires string, array, or map argument")
		}
	})

	e.functionRegistry.Register("trim", func(args ...interface{}) (interface{}, error) {
		if len(args) != 1 {
			return nil, errors.New("trim requires 1 argument")
		}
		str, ok := args[0].(string)
		if !ok {
			return nil, errors.New("trim requires string argument")
		}
		return strings.TrimSpace(str), nil
	})

	e.functionRegistry.Register("replace", func(args ...interface{}) (interface{}, error) {
		if len(args) != 3 {
			return nil, errors.New("replace requires 3 arguments")
		}
		str, ok1 := args[0].(string)
		old, ok2 := args[1].(string)
		newVal, ok3 := args[2].(string)
		if !ok1 || !ok2 || !ok3 {
			return nil, errors.New("replace requires string arguments")
		}
		return strings.ReplaceAll(str, old, newVal), nil
	})

	e.functionRegistry.Register("substr", func(args ...interface{}) (interface{}, error) {
		if len(args) < 2 || len(args) > 3 {
			return nil, errors.New("substr requires 2 or 3 arguments")
		}
		str, ok := args[0].(string)
		if !ok {
			return nil, errors.New("substr requires string as first argument")
		}
		start, ok := toInt(args[1])
		if !ok {
			return nil, errors.New("substr requires integer as second argument")
		}
		if start < 0 || start > len(str) {
			return nil, errors.New("substr start out of range")
		}
		if len(args) == 2 {
			return str[start:], nil
		}
		length, ok := toInt(args[2])
		if !ok {
			return nil, errors.New("substr requires integer as third argument")
		}
		end := start + length
		if end > len(str) {
			end = len(str)
		}
		return str[start:end], nil
	})

	e.functionRegistry.Register("now", func(args ...interface{}) (interface{}, error) {
		return time.Now().Format(time.RFC3339), nil
	})

	e.functionRegistry.Register("join", func(args ...interface{}) (interface{}, error) {
		if len(args) != 2 {
			return nil, errors.New("join requires 2 arguments")
		}
		arr, ok := args[0].([]interface{})
		if !ok {
			return nil, errors.New("join requires array as first argument")
		}
		sep, ok := args[1].(string)
		if !ok {
			return nil, errors.New("join requires string as second argument")
		}
		strs := make([]string, len(arr))
		for i, v := range arr {
			strs[i] = fmt.Sprint(v)
		}
		return strings.Join(strs, sep), nil
	})

	e.functionRegistry.Register("split", func(args ...interface{}) (interface{}, error) {
		if len(args) != 2 {
			return nil, errors.New("split requires 2 arguments")
		}
		str, ok := args[0].(string)
		if !ok {
			return nil, errors.New("split requires string as first argument")
		}
		sep, ok := args[1].(string)
		if !ok {
			return nil, errors.New("split requires string as second argument")
		}
		parts := strings.Split(str, sep)
		result := make([]interface{}, len(parts))
		for i, v := range parts {
			result[i] = v
		}
		return result, nil
	})

	e.functionRegistry.Register("has_prefix", func(args ...interface{}) (interface{}, error) {
		if len(args) != 2 {
			return nil, errors.New("has_prefix requires 2 arguments")
		}
		str, ok1 := args[0].(string)
		prefix, ok2 := args[1].(string)
		if !ok1 || !ok2 {
			return nil, errors.New("has_prefix requires string arguments")
		}
		return strings.HasPrefix(str, prefix), nil
	})

	e.functionRegistry.Register("has_suffix", func(args ...interface{}) (interface{}, error) {
		if len(args) != 2 {
			return nil, errors.New("has_suffix requires 2 arguments")
		}
		str, ok1 := args[0].(string)
		suffix, ok2 := args[1].(string)
		if !ok1 || !ok2 {
			return nil, errors.New("has_suffix requires string arguments")
		}
		return strings.HasSuffix(str, suffix), nil
	})

	// Existing built-ins
	e.functionRegistry.Register("upper", func(args ...interface{}) (interface{}, error) {
		if len(args) != 1 {
			return nil, errors.New("upper requires 1 argument")
		}
		str, ok := args[0].(string)
		if !ok {
			return nil, errors.New("upper requires string argument")
		}
		return strings.ToUpper(str), nil
	})

	e.functionRegistry.Register("lower", func(args ...interface{}) (interface{}, error) {
		if len(args) != 1 {
			return nil, errors.New("lower requires 1 argument")
		}
		str, ok := args[0].(string)
		if !ok {
			return nil, errors.New("lower requires string argument")
		}
		return strings.ToLower(str), nil
	})

	e.functionRegistry.Register("concat", func(args ...interface{}) (interface{}, error) {
		var result strings.Builder
		for _, arg := range args {
			result.WriteString(fmt.Sprint(arg))
		}
		return result.String(), nil
	})

	e.functionRegistry.Register("upper", func(args ...interface{}) (interface{}, error) {
		if len(args) != 1 {
			return nil, errors.New("upper requires 1 argument")
		}
		str, ok := args[0].(string)
		if !ok {
			return nil, errors.New("upper requires string argument")
		}
		return strings.ToUpper(str), nil
	})

	e.functionRegistry.Register("lower", func(args ...interface{}) (interface{}, error) {
		if len(args) != 1 {
			return nil, errors.New("lower requires 1 argument")
		}
		str, ok := args[0].(string)
		if !ok {
			return nil, errors.New("lower requires string argument")
		}
		return strings.ToLower(str), nil
	})

	e.functionRegistry.Register("concat", func(args ...interface{}) (interface{}, error) {
		var result strings.Builder
		for _, arg := range args {
			result.WriteString(fmt.Sprint(arg))
		}
		return result.String(), nil
	})
}

// Helper for substr
func toInt(v interface{}) (int, bool) {
	switch val := v.(type) {
	case int:
		return val, true
	case int32:
		return int(val), true
	case int64:
		return int(val), true
	case float64:
		return int(val), true
	case float32:
		return int(val), true
	default:
		return 0, false
	}
}

// registerDefaultOperators registers built-in operators
func (e *Engine) registerDefaultOperators() {
	e.operatorRegistry.Register("eq", func(left, right interface{}) (bool, error) {
		return fmt.Sprint(left) == fmt.Sprint(right), nil
	})

	e.operatorRegistry.Register("ne", func(left, right interface{}) (bool, error) {
		return fmt.Sprint(left) != fmt.Sprint(right), nil
	})

	e.operatorRegistry.Register("gt", func(left, right interface{}) (bool, error) {
		l, r, err := toFloat64Pair(left, right)
		if err != nil {
			return false, err
		}
		return l > r, nil
	})

	e.operatorRegistry.Register("lt", func(left, right interface{}) (bool, error) {
		l, r, err := toFloat64Pair(left, right)
		if err != nil {
			return false, err
		}
		return l < r, nil
	})

	e.operatorRegistry.Register("gte", func(left, right interface{}) (bool, error) {
		l, r, err := toFloat64Pair(left, right)
		if err != nil {
			return false, err
		}
		return l >= r, nil
	})

	e.operatorRegistry.Register("lte", func(left, right interface{}) (bool, error) {
		l, r, err := toFloat64Pair(left, right)
		if err != nil {
			return false, err
		}
		return l <= r, nil
	})

	e.operatorRegistry.Register("contains", func(left, right interface{}) (bool, error) {
		leftStr := fmt.Sprint(left)
		rightStr := fmt.Sprint(right)
		return strings.Contains(leftStr, rightStr), nil
	})

	e.operatorRegistry.Register("in", func(left, right interface{}) (bool, error) {
		arr, ok := right.([]interface{})
		if !ok {
			return false, errors.New("in operator requires array")
		}
		leftStr := fmt.Sprint(left)
		for _, item := range arr {
			if fmt.Sprint(item) == leftStr {
				return true, nil
			}
		}
		return false, nil
	})

	e.operatorRegistry.Register("exists", func(left, right interface{}) (bool, error) {
		return left != nil, nil
	})

	// Additional commonly used operators:

	e.operatorRegistry.Register("starts_with", func(left, right interface{}) (bool, error) {
		leftStr := fmt.Sprint(left)
		rightStr := fmt.Sprint(right)
		return strings.HasPrefix(leftStr, rightStr), nil
	})

	e.operatorRegistry.Register("ends_with", func(left, right interface{}) (bool, error) {
		leftStr := fmt.Sprint(left)
		rightStr := fmt.Sprint(right)
		return strings.HasSuffix(leftStr, rightStr), nil
	})

	e.operatorRegistry.Register("matches", func(left, right interface{}) (bool, error) {
		leftStr := fmt.Sprint(left)
		pattern, ok := right.(string)
		if !ok {
			return false, errors.New("matches operator requires string pattern")
		}
		re, err := regexp.Compile(pattern)
		if err != nil {
			return false, err
		}
		return re.MatchString(leftStr), nil
	})

	e.operatorRegistry.Register("empty", func(left, right interface{}) (bool, error) {
		if left == nil {
			return true, nil
		}
		switch v := left.(type) {
		case string:
			return v == "", nil
		case []interface{}:
			return len(v) == 0, nil
		case map[string]interface{}:
			return len(v) == 0, nil
		default:
			return false, nil
		}
	})

	e.operatorRegistry.Register("not_empty", func(left, right interface{}) (bool, error) {
		if left == nil {
			return false, nil
		}
		switch v := left.(type) {
		case string:
			return v != "", nil
		case []interface{}:
			return len(v) > 0, nil
		case map[string]interface{}:
			return len(v) > 0, nil
		default:
			return true, nil
		}
	})
}

// Helper functions

func deepCopy(src map[string]interface{}) map[string]interface{} {
	data, _ := json.Marshal(src)
	var dst map[string]interface{}
	json.Unmarshal(data, &dst)
	return dst
}

func toFloat64Pair(left, right interface{}) (float64, float64, error) {
	l, lok := toFloat64(left)
	r, rok := toFloat64(right)
	if !lok || !rok {
		return 0, 0, errors.New("cannot convert to numbers")
	}
	return l, r, nil
}

func toFloat64(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	case int32:
		return float64(val), true
	default:
		return 0, false
	}
}

func getValueAtPath(data map[string]interface{}, path string) (interface{}, error) {
	if path == "" || path == "$" {
		return data, nil
	}

	parts := strings.Split(strings.TrimPrefix(path, "$."), ".")
	return getValueAtPathRecursive(data, parts)
}

func getValueAtPathRecursive(current interface{}, parts []string) (interface{}, error) {
	if len(parts) == 0 {
		return current, nil
	}
	part := parts[0]

	// Wildcard array: e.g. items[]
	if strings.HasSuffix(part, "[]") {
		key := strings.TrimSuffix(part, "[]")
		m, ok := current.(map[string]interface{})
		if !ok {
			return nil, ErrPathNotFound
		}
		arr, ok := m[key].([]interface{})
		if !ok {
			return nil, ErrPathNotFound
		}
		results := make([]interface{}, 0, len(arr))
		for _, elem := range arr {
			val, err := getValueAtPathRecursive(elem, parts[1:])
			if err == nil {
				results = append(results, val)
			}
		}
		return results, nil
	}

	// Indexed array: e.g. items[0]
	if strings.Contains(part, "[") && strings.HasSuffix(part, "]") {
		key := part[:strings.Index(part, "[")]
		indexStr := part[strings.Index(part, "[")+1 : len(part)-1]
		m, ok := current.(map[string]interface{})
		if !ok {
			return nil, ErrPathNotFound
		}
		arr, ok := m[key].([]interface{})
		if !ok {
			return nil, ErrPathNotFound
		}
		var idx int
		fmt.Sscanf(indexStr, "%d", &idx)
		if idx < 0 || idx >= len(arr) {
			return nil, ErrPathNotFound
		}
		return getValueAtPathRecursive(arr[idx], parts[1:])
	}

	// Normal object key
	m, ok := current.(map[string]interface{})
	if !ok {
		return nil, ErrPathNotFound
	}
	val, exists := m[part]
	if !exists {
		return nil, ErrPathNotFound
	}
	return getValueAtPathRecursive(val, parts[1:])
}

// setValueAtPath supports array wildcards ([])
func setValueAtPath(data map[string]interface{}, path string, value interface{}) error {
	if path == "" || path == "$" {
		return errors.New("cannot set root path")
	}
	parts := strings.Split(strings.TrimPrefix(path, "$."), ".")
	return setValueAtPathRecursive(data, parts, value)
}

func setValueAtPathRecursive(current interface{}, parts []string, value interface{}) error {
	if len(parts) == 0 {
		return errors.New("invalid path")
	}
	part := parts[0]

	// Wildcard array: e.g. items[]
	if strings.HasSuffix(part, "[]") {
		key := strings.TrimSuffix(part, "[]")
		m, ok := current.(map[string]interface{})
		if !ok {
			return ErrPathNotFound
		}
		arr, ok := m[key].([]interface{})
		if !ok {
			return ErrPathNotFound
		}
		for i := range arr {
			if len(parts) == 1 {
				arr[i] = value
			} else {
				if err := setValueAtPathRecursive(arr[i], parts[1:], value); err != nil {
					return err
				}
			}
		}
		m[key] = arr
		return nil
	}

	// Indexed array: e.g. items[0]
	if strings.Contains(part, "[") && strings.HasSuffix(part, "]") {
		key := part[:strings.Index(part, "[")]
		indexStr := part[strings.Index(part, "[")+1 : len(part)-1]
		m, ok := current.(map[string]interface{})
		if !ok {
			return ErrPathNotFound
		}
		arr, ok := m[key].([]interface{})
		if !ok {
			return ErrPathNotFound
		}
		var idx int
		fmt.Sscanf(indexStr, "%d", &idx)
		if idx < 0 || idx >= len(arr) {
			return ErrPathNotFound
		}
		if len(parts) == 1 {
			arr[idx] = value
		} else {
			if err := setValueAtPathRecursive(arr[idx], parts[1:], value); err != nil {
				return err
			}
		}
		m[key] = arr
		return nil
	}

	// Normal object key
	m, ok := current.(map[string]interface{})
	if !ok {
		return ErrPathNotFound
	}
	if len(parts) == 1 {
		m[part] = value
		return nil
	}
	if _, exists := m[part]; !exists {
		m[part] = make(map[string]interface{})
	}
	return setValueAtPathRecursive(m[part], parts[1:], value)
}

// deleteValueAtPath supports array wildcards ([])
func deleteValueAtPath(data map[string]interface{}, path string) error {
	if path == "" || path == "$" {
		return errors.New("cannot delete root path")
	}
	parts := strings.Split(strings.TrimPrefix(path, "$."), ".")
	return deleteValueAtPathRecursive(data, parts)
}

func deleteValueAtPathRecursive(current interface{}, parts []string) error {
	if len(parts) == 0 {
		return errors.New("invalid path")
	}
	part := parts[0]

	// Wildcard array: e.g. items[]
	if strings.HasSuffix(part, "[]") {
		key := strings.TrimSuffix(part, "[]")
		m, ok := current.(map[string]interface{})
		if !ok {
			return ErrPathNotFound
		}
		arr, ok := m[key].([]interface{})
		if !ok {
			return ErrPathNotFound
		}
		for i := range arr {
			if len(parts) == 1 {
				arr[i] = nil
			} else {
				if err := deleteValueAtPathRecursive(arr[i], parts[1:]); err != nil {
					return err
				}
			}
		}
		m[key] = arr
		return nil
	}

	// Indexed array: e.g. items[0]
	if strings.Contains(part, "[") && strings.HasSuffix(part, "]") {
		key := part[:strings.Index(part, "[")]
		indexStr := part[strings.Index(part, "[")+1 : len(part)-1]
		m, ok := current.(map[string]interface{})
		if !ok {
			return ErrPathNotFound
		}
		arr, ok := m[key].([]interface{})
		if !ok {
			return ErrPathNotFound
		}
		var idx int
		fmt.Sscanf(indexStr, "%d", &idx)
		if idx < 0 || idx >= len(arr) {
			return ErrPathNotFound
		}
		if len(parts) == 1 {
			arr[idx] = nil
		} else {
			if err := deleteValueAtPathRecursive(arr[idx], parts[1:]); err != nil {
				return err
			}
		}
		m[key] = arr
		return nil
	}

	// Normal object key
	m, ok := current.(map[string]interface{})
	if !ok {
		return ErrPathNotFound
	}
	if len(parts) == 1 {
		delete(m, part)
		return nil
	}
	val, exists := m[part]
	if !exists {
		return ErrPathNotFound
	}
	return deleteValueAtPathRecursive(val, parts[1:])
}

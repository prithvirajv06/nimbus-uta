package engine

import (
	"fmt"
	"time"

	"github.com/dop251/goja"
)

// Execute runs the compiled script against the input data
func Execute(script string, inputData map[string]interface{}) (map[string]interface{}, []string, error) {
	vm := goja.New()

	// 1. Setup Input Data
	// We wrap input in a container so "data.Age" in script works naturally
	vm.Set("data", inputData)

	// 2. Register Capabilities (The "Bridge")
	vm.Set("$http", func(method, url string) map[string]interface{} {
		fmt.Printf("[Network] Calling %s %s...\n", method, url)
		// Mock Response for demo
		return map[string]interface{}{
			"status": "success",
			"score":  95,
		}
	})

	// 3. Run with Timeout
	time.AfterFunc(200*time.Millisecond, func() {
		vm.Interrupt("timeout")
	})

	_, err := vm.RunString(script)
	if err != nil {
		return nil, nil, err
	}

	// 4. Extract Modified Data
	res := vm.Get("data")
	logs := vm.Get("log")
	return res.Export().(map[string]interface{}), logs.Export().([]string), nil
}

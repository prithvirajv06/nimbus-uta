package engine

import (
	"time"

	"github.com/dop251/goja"
)

// Execute runs the compiled script against the input data
func Execute(script string, inputData map[string]interface{}) (map[string]interface{}, []interface{}, error) {
	vm := goja.New()

	time.AfterFunc(200*time.Millisecond, func() {
		vm.Interrupt("timeout")
	})

	_, err := vm.RunString(script)
	if err != nil {
		return nil, nil, err
	}
	executeWorkflow, ok := goja.AssertFunction(vm.Get("executeWorkflow"))
	if !ok {
		panic("Not a function")
	}
	response, err := executeWorkflow(goja.Undefined(), vm.ToValue(inputData))
	if err != nil {
		return nil, nil, err
	}
	calculatedResp := response.ToObject(vm).Export().(map[string]interface{})
	jsonResponse := calculatedResp["data"].(map[string]interface{})
	logsInterface := calculatedResp["log"].([]interface{})
	print(response)
	// 4. Extract Modified Data
	return jsonResponse, logsInterface, nil
}

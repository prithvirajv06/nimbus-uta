package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "io"
)

type Logger struct {
    Logs []string `json:"logs"`
}

func (l *Logger) Add(msg string) {
    l.Logs = append(l.Logs, msg)
}

// --- DATA CONTRACTS ---

type CollateralItems struct {
	Type string `json:"type"`
	MarketValue float64 `json:"marketvalue"`
	Status string `json:"status"`
}

type PolicyInfo struct {
	CollateralItems []CollateralItems `json:"collateralitems"`
}


type RuleService struct{}

// Helper for Network Calls
func CallExternal(url, method string, payload interface{}, logger *Logger) (map[string]interface{}, error) {
    logger.Add(fmt.Sprintf("Calling external: %s %s", method, url))
    // Standard library implementation
    return map[string]interface{}{"status": "success"}, nil
}

func (rs *RuleService) Process(w http.ResponseWriter, r *http.Request) {
    var data PolicyInfo
    logger := &Logger{}
    body, _ := io.ReadAll(r.Body)
    if err := json.Unmarshal(body, &data); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // --- GENERATED LOGIC ---
	logger.Add("Loop over: data.CollateralItems")
	for i := range data.CollateralItems {
		item := &data.CollateralItems[i]
		logger.Add("Condition: item.MarketValue > 100000 && item.Type == \"REAL_ESTATE\"")
		if item.MarketValue > 100000 && item.Type == "REAL_ESTATE" {
			logger.Add("Assign: item.Status = \"PREMIUM_ASSET\"")
			item.Status = "PREMIUM_ASSET"
		} else {
			logger.Add("Condition failed: item.MarketValue > 100000 && item.Type == \"REAL_ESTATE\"")
		}
	}


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
    fmt.Println("Rule Service starting on :8081")
    http.ListenAndServe(":8081", nil)
}

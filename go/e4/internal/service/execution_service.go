package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prithvirajv06/nimbus-uta/go/engine/config"
	"github.com/prithvirajv06/nimbus-uta/go/engine/engine"
	"github.com/prithvirajv06/nimbus-uta/go/engine/internal/models"
	"github.com/prithvirajv06/nimbus-uta/go/engine/internal/repository"
	"github.com/prithvirajv06/nimbus-uta/go/engine/pkg/cache"
	"github.com/prithvirajv06/nimbus-uta/go/engine/pkg/database"
	"github.com/prithvirajv06/nimbus-uta/go/engine/pkg/messaging"
)

type ExecutionService struct {
	mongo    *database.MongoDB
	rabbitMQ *messaging.RabbitMQ
	redis    *cache.RedisClient
	cfg      *config.Config
}

func NewExecutionService(db *database.MongoDB, rabbitMQ *messaging.RabbitMQ, cfg *config.Config,
	redisClient *cache.RedisClient) *ExecutionService {
	rabbitMQ.DeclareQueue("variablepackage.event", true)
	return &ExecutionService{
		mongo:    db,
		rabbitMQ: rabbitMQ,
		redis:    redisClient,
		cfg:      cfg,
	}
}

func (lfs *ExecutionService) HandleRuleExecution(c *gin.Context) {
	nimbID := c.Query("nimb_id")
	versionStr := c.Query("version")
	version, _ := strconv.Atoi(versionStr)
	repo := repository.NewGenericRepository[models.LogicFlow](c.Request.Context(), lfs.mongo.Database, "logic_flows")
	eng, err := repo.FindOne(map[string]interface{}{"nimb_id": nimbID, "audit.is_archived": false, "audit.version": version})
	if HandleError(c, err, "Failed to fetch logic flow for execution") {
		return
	}
	var req map[string]interface{}

	// 1. Parse the incoming JSON Facts (the data to be tested)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	start := time.Now()
	engine := engine.NewRuleEngine(req)
	err = engine.ExecuteWorkflow(*eng)
	if err != nil {
		fmt.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Execution failed",
			"details": err.Error(),
		})
		return
	}
	duration := time.Since(start)
	c.Header("X_TIME-TAKEN", strconv.FormatInt(duration.Milliseconds(), 10))
	if c.GetHeader("X-NIMBUS-DEBUG") == "YES" {
		c.JSON(http.StatusOK, gin.H{
			"data": engine.Input,
			"logs": engine.GetLog(),
		})
	} else {
		c.JSON(http.StatusOK, engine.Input)
	}
}

func (lfs *ExecutionService) HandleDTExecution(c *gin.Context) {
	nimbID := c.Query("nimb_id")
	versionStr := c.Query("version")
	version, _ := strconv.Atoi(versionStr)
	repo := repository.NewGenericRepository[models.DecisionTable](c.Request.Context(), lfs.mongo.Database, "decision_tables")
	table, err := repo.FindOne(map[string]interface{}{"nimb_id": nimbID, "audit.is_archived": false, "audit.version": version})
	if HandleError(c, err, "Failed to fetch decision table for execution") {
		return
	}
	var body []byte = make([]byte, c.Request.ContentLength)
	_, err = c.Request.Body.Read(body)
	start := time.Now()
	engine := engine.NewDTEngine()
	output, logs, err := engine.ProcessDecisionTable(c.Request.Context(), *table, body)
	if err != nil {
		fmt.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Execution failed",
			"details": err.Error(),
		})
		return
	}
	duration := time.Since(start)
	c.Header("X_TIME-TAKEN", strconv.FormatInt(duration.Milliseconds(), 10))
	var response interface{}
	json.Unmarshal(output, &response)
	if c.GetHeader("X-NIMBUS-DEBUG") == "YES" {
		c.JSON(http.StatusOK, gin.H{
			"data": response,
			"logs": logs,
		})
	} else {
		c.JSON(http.StatusOK, response)
	}
}

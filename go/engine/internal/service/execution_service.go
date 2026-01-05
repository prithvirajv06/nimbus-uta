package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prithvirajv06/nimbus-uta/go/engine/config"
	"github.com/prithvirajv06/nimbus-uta/go/engine/heart"
	"github.com/prithvirajv06/nimbus-uta/go/engine/internal/models"
	"github.com/prithvirajv06/nimbus-uta/go/engine/internal/repository"
	"github.com/prithvirajv06/nimbus-uta/go/engine/pkg/database"
	"github.com/prithvirajv06/nimbus-uta/go/engine/pkg/messaging"
)

type ExecutionService struct {
	mongo    *database.MongoDB
	rabbitMQ *messaging.RabbitMQ
	engine   *heart.Engine
	cfg      *config.Config
}

func NewExecutionService(db *database.MongoDB, rabbitMQ *messaging.RabbitMQ, cfg *config.Config) *ExecutionService {
	rabbitMQ.DeclareQueue("variablepackage.event", true)
	return &ExecutionService{
		mongo:    db,
		rabbitMQ: rabbitMQ,
		cfg:      cfg,
	}
}

func (s *ExecutionService) ExecuteDT(c *gin.Context) {
	nimbID := c.Query("nimb_id")
	versionStr := c.Query("version")
	isDebug := c.Query("is_debug")
	version, _ := strconv.Atoi(versionStr)
	repo := repository.NewGenericRepository[models.DecisionTable](c.Request.Context(), s.mongo.Database, "decision_tables")
	decisionTable, err := repo.FindOne(map[string]interface{}{"nimb_id": nimbID, "audit.is_archived": false, "audit.version": version})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch decision table"})
		return
	}
	rawBody, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	resp, logStack, err := s.engine.ProcessDecisionTable(c, *decisionTable, rawBody)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to process decision table"})
		return
	}
	if isDebug == "YES" {
		RespondJSON(c, 200, "success", "Execution Completed", gin.H{
			"response": resp,
			"logs":     logStack,
		})
	} else {
		c.JSON(200, resp)
	}
}

func (s *ExecutionService) ExecuteLogicalFlow(c *gin.Context) {

	nimbID := c.Query("nimb_id")
	versionStr := c.Query("version")
	isDebug := c.Query("is_debug")

	version, _ := strconv.Atoi(versionStr)
	repo := repository.NewGenericRepository[models.LogicFlow](c.Request.Context(), s.mongo.Database, "logical_flows")
	logicFlow, err := repo.FindOne(map[string]interface{}{"nimb_id": nimbID, "audit.is_archived": false, "audit.version": version})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch decision table"})
		return
	}
	rawBody, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	resp, logStack, err := s.engine.ExecuteRuleSet(c, *logicFlow, rawBody)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to process decision table"})
		return
	}
	if isDebug == "YES" {
		RespondJSON(c, 200, "success", "Execution Completed", gin.H{
			"response": resp,
			"logs":     logStack,
		})
	} else {
		c.JSON(200, resp)
	}
}

package service

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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

	// 3. Generate/Get Compiled Script
	// In production, we cache the generated jsCode string based on a version hash
	var redisCacheKey = "engine_script_" + eng.NIMB_ID + "_v" + strconv.Itoa(eng.Audit.Version)
	var jsCode string
	cached, err := lfs.redis.GetString(c.Request.Context(), redisCacheKey)
	if err == nil {
		jsCode = cached
	}
	// if err != nil || cached == "" {
	jsCode = engine.GenerateScript(*eng)
	// Cache for future use
	_ = lfs.redis.Set(c.Request.Context(), redisCacheKey, jsCode, 10*time.Minute)
	// }

	// 4. Execute in the Sandboxed Runtime
	start := time.Now()
	finalData, logs, err := engine.Execute(jsCode, req)
	duration := time.Since(start)
	//Save jsCode in file
	path1 := filepath.Join(os.TempDir(), "dat1")
	_ = os.WriteFile(path1, []byte(jsCode), 0644)
	if err != nil {
		fmt.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Execution failed",
			"details": err.Error(),
		})
		return
	}

	// 5. Return the modified facts and audit info
	c.Header("X_TIME-TAKEN", strconv.FormatInt(duration.Milliseconds(), 10))
	if c.GetHeader("X-NIMBUS-DEBUG") == "YES" {
		c.JSON(http.StatusOK, gin.H{
			"data": finalData,
			"logs": logs,
		})
		return
	} else {
		c.JSON(http.StatusOK, finalData)
	}
}

package service

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prithvirajv06/nimbus-uta/go/engine/config"
	"github.com/prithvirajv06/nimbus-uta/go/engine/engine"
	"github.com/prithvirajv06/nimbus-uta/go/engine/internal/models"
	"github.com/prithvirajv06/nimbus-uta/go/engine/internal/repository"
	"github.com/prithvirajv06/nimbus-uta/go/engine/internal/utils"
	"github.com/prithvirajv06/nimbus-uta/go/engine/pkg/database"
	"github.com/prithvirajv06/nimbus-uta/go/engine/pkg/messaging"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
)

type ExecutionService struct {
	mongo    *database.MongoDB
	rabbitMQ *messaging.RabbitMQ
	redis    *redis.Client
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

func (lfs *ExecutionService) CreateNewEngine(c *gin.Context) {
	var payload models.WorkflowDef
	err := c.ShouldBindJSON(&payload)
	if HandleError(c, err, "Unable to unmarshel payload") {
		return
	}
	repo := repository.NewGenericRepository[models.WorkflowDef](c.Request.Context(), lfs.mongo.Database, "engines")
	payload.NIMB_ID = utils.GenerateNIMBID("N_L_ENG_")
	payload.Audit.SetInitialAudit(c)
	payload.Audit.Version, _ = GetNextVersionNumber(c, lfs.mongo.Database, payload.NIMB_ID)
	_, err = repo.InsertOne(payload)
	if HandleError(c, err, "Failed to create logic flow") {
		return
	}
	RespondJSON(c, 201, "success", "Logic flow created", payload)
}

func (lfs *ExecutionService) UpdateEngine(c *gin.Context) {
	var payload models.WorkflowDef
	err := c.ShouldBindJSON(&payload)
	if HandleError(c, err, "Unable to unmarshel payload") {
		return
	}
	// Archive Old Version and Create New Version
	repo := repository.NewGenericRepository[models.WorkflowDef](c.Request.Context(), lfs.mongo.Database, "engines")
	err = ArchiveEntity(c, repo, payload.NIMB_ID, payload.Audit.Version)
	if HandleError(c, err, "Failed to archive old version of logic flow") {
		return
	}
	payload.Audit.MinorVersion += 1
	_, err = repo.InsertOne(payload)
	if HandleError(c, err, "Failed to create new version of logic flow") {
		return
	}
	RespondJSON(c, 200, "success", "Logic flow updated", payload)
}

func (lfs *ExecutionService) GetEngineByNIMBID(c *gin.Context) {
	nimbID := c.Query("nimb_id")
	versionStr := c.Query("version")
	version, _ := strconv.Atoi(versionStr)
	repo := repository.NewGenericRepository[models.WorkflowDef](c.Request.Context(), lfs.mongo.Database, "engines")
	engineFromDb, err := repo.FindOne(map[string]interface{}{"nimb_id": nimbID, "audit.is_archived": false, "audit.version": version})
	if HandleError(c, err, "Failed to fetch logic flow") {
		return
	}
	RespondJSON(c, 200, "success", "Logic flow retrieved", engineFromDb)
}

func (lfs *ExecutionService) GetAllEngines(c *gin.Context) {
	repo := repository.NewGenericRepository[models.WorkflowDef](c.Request.Context(), lfs.mongo.Database, "engines")
	option := GetCommonSortOption()
	option.SetProjection(bson.M{"variable_packages": 0, "logical_steps": 0})
	engines, err := repo.FindMany(bson.M{"audit.is_archived": false}, option)
	if HandleError(c, err, "Failed to fetch logic flows") {
		return
	}
	RespondJSON(c, 200, "success", "Logic flows retrieved", engines)
}

func (lfs *ExecutionService) ArchiveEngine(c *gin.Context) {
	repo := repository.NewGenericRepository[models.WorkflowDef](c.Request.Context(), lfs.mongo.Database, "engines")
	versionStr := c.Query("version")
	version, _ := strconv.Atoi(versionStr)
	err := ArchiveEntity(c, repo, c.Query("nimb_id"), version)
	if HandleError(c, err, "Failed to archive logic flow") {
		return
	}
	RespondJSON(c, 200, "success", "Logic flow archived", nil)
}

func (lfs *ExecutionService) CloneEngine(c *gin.Context) {
	nimbID := c.Query("nimb_id")
	repo := repository.NewGenericRepository[models.WorkflowDef](c.Request.Context(), lfs.mongo.Database, "engines")
	engineFromDb, err := repo.FindOne(map[string]interface{}{"nimb_id": nimbID, "audit.is_archived": false})
	if HandleError(c, err, "Failed to fetch logic flow to clone") {
		return
	}
	engineFromDb.Audit.SetInitialAudit(c)
	engineFromDb.Audit.Version, _ = GetNextVersionNumber(c, lfs.mongo.Database, engineFromDb.NIMB_ID)
	_, err = repo.InsertOne(*engineFromDb)
	if HandleError(c, err, "Failed to clone logic flow") {
		return
	}
	RespondJSON(c, 201, "success", "Logic flow cloned", engineFromDb)
}

func (lfs *ExecutionService) HandleRuleExecution(c *gin.Context) {
	nimbID := c.Query("nimb_id")
	versionStr := c.Query("version")
	version, _ := strconv.Atoi(versionStr)
	repo := repository.NewGenericRepository[models.WorkflowDef](c.Request.Context(), lfs.mongo.Database, "engines")
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
	cached, err := lfs.redis.Get(c.Request.Context(), redisCacheKey).Result()
	if err == nil {
		jsCode = cached
	}
	if err != nil || cached == "" {
		jsCode = engine.GenerateScript(eng.Pipeline)
		// Cache for future use
		_ = lfs.redis.Set(c.Request.Context(), redisCacheKey, jsCode, 10*time.Minute)
	}

	// 4. Execute in the Sandboxed Runtime
	start := time.Now()
	finalData, err := engine.Execute(jsCode, req)
	duration := time.Since(start)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Execution failed",
			"details": err.Error(),
		})
		return
	}

	// 5. Return the modified facts and audit info
	c.Header("X_TIME-TAKEN", strconv.FormatInt(duration.Milliseconds(), 10))
	c.JSON(http.StatusOK, finalData)
}

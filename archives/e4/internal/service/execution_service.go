package service

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prithvirajv06/nimbus-uta/go/engine/config"
	"github.com/prithvirajv06/nimbus-uta/go/engine/engine"
	"github.com/prithvirajv06/nimbus-uta/go/engine/internal/repository"
	"github.com/prithvirajv06/nimbus-uta/go/engine/internal/utils"
	"github.com/prithvirajv06/nimbus-uta/go/engine/pkg/database"
	"github.com/prithvirajv06/nimbus-uta/go/engine/pkg/messaging"
	"go.mongodb.org/mongo-driver/bson"
)

type ExecutionService struct {
	mongo    *database.MongoDB
	rabbitMQ *messaging.RabbitMQ
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
	var payload engine.RuleService
	err := c.ShouldBindJSON(&payload)
	if HandleError(c, err, "Unable to unmarshel payload") {
		return
	}
	repo := repository.NewGenericRepository[engine.RuleService](c.Request.Context(), lfs.mongo.Database, "engines")
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
	var payload engine.RuleService
	err := c.ShouldBindJSON(&payload)
	if HandleError(c, err, "Unable to unmarshel payload") {
		return
	}
	// Archive Old Version and Create New Version
	repo := repository.NewGenericRepository[engine.RuleService](c.Request.Context(), lfs.mongo.Database, "engines")
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
	repo := repository.NewGenericRepository[engine.RuleService](c.Request.Context(), lfs.mongo.Database, "engines")
	engineFromDb, err := repo.FindOne(map[string]interface{}{"nimb_id": nimbID, "audit.is_archived": false, "audit.version": version})
	if HandleError(c, err, "Failed to fetch logic flow") {
		return
	}
	comp := &engine.Compiler{}
	comp.CompileFromRequest(*engineFromDb)
	RespondJSON(c, 200, "success", "Logic flow retrieved", engineFromDb)
}

func (lfs *ExecutionService) GetAllEngines(c *gin.Context) {
	repo := repository.NewGenericRepository[engine.RuleService](c.Request.Context(), lfs.mongo.Database, "engines")
	option := GetCommonSortOption()
	option.SetProjection(bson.M{"variable_packages": 0, "logical_steps": 0})
	engines, err := repo.FindMany(bson.M{"audit.is_archived": false}, option)
	if HandleError(c, err, "Failed to fetch logic flows") {
		return
	}
	RespondJSON(c, 200, "success", "Logic flows retrieved", engines)
}

func (lfs *ExecutionService) ArchiveEngine(c *gin.Context) {
	repo := repository.NewGenericRepository[engine.RuleService](c.Request.Context(), lfs.mongo.Database, "engines")
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
	repo := repository.NewGenericRepository[engine.RuleService](c.Request.Context(), lfs.mongo.Database, "engines")
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

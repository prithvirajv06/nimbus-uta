package service

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prithvirajv06/nimbus-uta/go/core/config"
	"github.com/prithvirajv06/nimbus-uta/go/core/internal/models"
	"github.com/prithvirajv06/nimbus-uta/go/core/internal/repository"
	"github.com/prithvirajv06/nimbus-uta/go/core/internal/utils"
	"github.com/prithvirajv06/nimbus-uta/go/core/pkg/database"
	"github.com/prithvirajv06/nimbus-uta/go/core/pkg/messaging"
	"go.mongodb.org/mongo-driver/bson"
)

type LogicFlowService struct {
	mongo    *database.MongoDB
	cfg      *config.Config
	rabbitMQ *messaging.RabbitMQ
}

func NewLogicFlowService(db *database.MongoDB, rabbitMQ *messaging.RabbitMQ, cfg *config.Config) *LogicFlowService {
	rabbitMQ.DeclareQueue("logicflow.event", true)
	return &LogicFlowService{
		mongo:    db,
		rabbitMQ: rabbitMQ,
		cfg:      cfg,
	}
}

func (lfs *LogicFlowService) CreateNewLogicFlow(c *gin.Context) {
	var payload models.LogicFlow
	err := c.ShouldBindJSON(&payload)
	if HandleError(c, err, "Unable to unmarshel payload") {
		return
	}
	repo := repository.NewGenericRepository[models.LogicFlow](c.Request.Context(), lfs.mongo.Database, "logic_flows")
	payload.NIMB_ID = utils.GenerateNIMBID("N_L_FLOW")
	payload.Audit.SetInitialAudit(c)
	payload.Audit.Version, _ = GetNextVersionNumber(c, lfs.mongo.Database, payload.NIMB_ID)
	_, err = repo.InsertOne(payload)
	if HandleError(c, err, "Failed to create logic flow") {
		return
	}
	RespondJSON(c, 201, "success", "Logic flow created", payload)
}

func (lfs *LogicFlowService) UpdateLogicFlow(c *gin.Context) {
	var payload models.LogicFlow
	err := c.ShouldBindJSON(&payload)
	if HandleError(c, err, "Unable to unmarshel payload") {
		return
	}
	// Archive Old Version and Create New Version
	repo := repository.NewGenericRepository[models.LogicFlow](c.Request.Context(), lfs.mongo.Database, "logic_flows")
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

func (lfs *LogicFlowService) GetLogicFlowByNIMBID(c *gin.Context) {
	nimbID := c.Query("nimb_id")
	versionStr := c.Query("version")
	version, _ := strconv.Atoi(versionStr)
	repo := repository.NewGenericRepository[models.LogicFlow](c.Request.Context(), lfs.mongo.Database, "logic_flows")
	logicFlow, err := repo.FindOne(map[string]interface{}{"nimb_id": nimbID, "audit.is_archived": false, "audit.version": version})
	if HandleError(c, err, "Failed to fetch logic flow") {
		return
	}
	RespondJSON(c, 200, "success", "Logic flow retrieved", logicFlow)
}

func (lfs *LogicFlowService) GetAllLogicFlows(c *gin.Context) {
	repo := repository.NewGenericRepository[models.LogicFlow](c.Request.Context(), lfs.mongo.Database, "logic_flows")
	option := GetCommonSortOption()
	option.SetProjection(bson.M{"variable_packages": 0, "logical_steps": 0})
	logicFlows, err := repo.FindMany(bson.M{"audit.is_archived": false}, option)
	if HandleError(c, err, "Failed to fetch logic flows") {
		return
	}
	RespondJSON(c, 200, "success", "Logic flows retrieved", logicFlows)
}

func (lfs *LogicFlowService) ArchiveLogicFlow(c *gin.Context) {
	repo := repository.NewGenericRepository[models.LogicFlow](c.Request.Context(), lfs.mongo.Database, "logic_flows")
	versionStr := c.Query("version")
	version, _ := strconv.Atoi(versionStr)
	err := ArchiveEntity(c, repo, c.Query("nimb_id"), version)
	if HandleError(c, err, "Failed to archive logic flow") {
		return
	}
	RespondJSON(c, 200, "success", "Logic flow archived", nil)
}

func (lfs *LogicFlowService) CloneLogicFlow(c *gin.Context) {
	nimbID := c.Query("nimb_id")
	repo := repository.NewGenericRepository[models.LogicFlow](c.Request.Context(), lfs.mongo.Database, "logic_flows")
	logicFlow, err := repo.FindOne(map[string]interface{}{"nimb_id": nimbID, "audit.is_archived": false})
	if HandleError(c, err, "Failed to fetch logic flow to clone") {
		return
	}
	logicFlow.Audit.SetInitialAudit(c)
	logicFlow.Audit.Version, _ = GetNextVersionNumber(c, lfs.mongo.Database, logicFlow.NIMB_ID)
	_, err = repo.InsertOne(*logicFlow)
	if HandleError(c, err, "Failed to clone logic flow") {
		return
	}
	RespondJSON(c, 201, "success", "Logic flow cloned", logicFlow)
}

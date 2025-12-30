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

type DecisionTableService struct {
	mongo    *database.MongoDB
	cfg      *config.Config
	rabbitMQ *messaging.RabbitMQ
}

func NewDecisionTableService(db *database.MongoDB, cfg *config.Config, rabbitMQ *messaging.RabbitMQ) *DecisionTableService {
	return &DecisionTableService{
		mongo:    db,
		cfg:      cfg,
		rabbitMQ: rabbitMQ,
	}
}

func (dts *DecisionTableService) CreateNewDecisionTable(c *gin.Context) {
	var payload models.DecisionTable
	err := c.ShouldBindJSON(&payload)
	if HandleError(c, err, "Unable to unmarshal payload") {
		return
	}
	repo := repository.NewGenericRepository[models.DecisionTable](c.Request.Context(), dts.mongo.Database, "decision_tables")
	payload.NIMB_ID = utils.GenerateNIMBID("N_D_TABLE")
	payload.Audit.SetInitialAudit(c)
	_, err = repo.InsertOne(payload)
	if HandleError(c, err, "Failed to create decision table") {
		return
	}
	RespondJSON(c, 201, "success", "Decision table created", payload)
}

func (dts *DecisionTableService) UpdateDecisionTable(c *gin.Context) {
	var payload models.DecisionTable
	err := c.ShouldBindJSON(&payload)
	if HandleError(c, err, "Unable to unmarshel payload") {
		return
	}
	// Archive Old Version and Create New Version
	repo := repository.NewGenericRepository[models.DecisionTable](c.Request.Context(), dts.mongo.Database, "decision_tables")
	err = ArchiveEntity(c, repo, payload.NIMB_ID, payload.Audit.Version)
	if HandleError(c, err, "Failed to archive old version of decision table") {
		return
	}
	payload.Audit.SetModifiedAudit(c)
	_, err = repo.InsertOne(payload)
	if HandleError(c, err, "Failed to create new version of decision table") {
		return
	}
	RespondJSON(c, 201, "success", "Decision table updated", payload)
}

func (dts *DecisionTableService) GetDecisionTableByNIMBID(c *gin.Context) {
	nimbID := c.Query("nimb_id")
	versionStr := c.Query("version")
	version, _ := strconv.Atoi(versionStr)
	repo := repository.NewGenericRepository[models.DecisionTable](c.Request.Context(), dts.mongo.Database, "decision_tables")
	table, err := repo.FindOne(map[string]interface{}{"nimb_id": nimbID, "audit.is_archived": false, "audit.version": version})
	if HandleError(c, err, "Failed to fetch decision table") {
		return
	}
	RespondJSON(c, 200, "success", "Decision table retrieved", table)
}

func (dts *DecisionTableService) GetAllDecisionTables(c *gin.Context) {
	repo := repository.NewGenericRepository[models.DecisionTable](c.Request.Context(), dts.mongo.Database, "decision_tables")
	option := GetCommonSortOption()
	option.SetProjection(bson.M{"variable_package": 0, "rules": 0, "input_columns": 0, "output_columns": 0})
	tables, err := repo.FindMany(bson.M{"audit.is_archived": false}, option)
	if HandleError(c, err, "Failed to fetch decision tables") {
		return
	}
	RespondJSON(c, 200, "success", "Decision tables retrieved", tables)
}

func (dts *DecisionTableService) ArchiveDecisionTable(c *gin.Context) {
	repo := repository.NewGenericRepository[models.DecisionTable](c.Request.Context(), dts.mongo.Database, "decision_tables")
	versionStr := c.Query("version")
	version, _ := strconv.Atoi(versionStr)
	err := ArchiveEntity(c, repo, c.Query("nimb_id"), version)
	if HandleError(c, err, "Failed to archive decision table") {
		return
	}
	RespondJSON(c, 200, "success", "Decision table archived", nil)
}

func (dts *DecisionTableService) CloneDecisionTable(c *gin.Context) {
	var nimbID = c.Query("nimb_id")
	var versionStr = c.Query("version")
	version, _ := strconv.Atoi(versionStr)
	repo := repository.NewGenericRepository[models.DecisionTable](c.Request.Context(), dts.mongo.Database, "decision_tables")
	origDT, err := repo.FindOne(map[string]interface{}{"nimb_id": nimbID, "audit.is_archived": false, "audit.version": version})
	if HandleError(c, err, "Failed to fetch decision table to clone") {
		return
	}
	origDT.Audit.SetInitialAudit(c)
	origDT.Audit.Version, _ = GetNextVersionNumber(c, dts.mongo.Database, nimbID)
	_, err = repo.InsertOne(*origDT)
	if HandleError(c, err, "Failed to clone decision table") {
		return
	}
	RespondJSON(c, 201, "success", "Decision table cloned", origDT)
}

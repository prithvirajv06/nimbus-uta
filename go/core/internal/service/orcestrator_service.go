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

type OrcestartionService struct {
	mongo    *database.MongoDB
	cfg      *config.Config
	rabbitMQ *messaging.RabbitMQ
}

func NewOrcestartionService(db *database.MongoDB, cfg *config.Config, rabbitMQ *messaging.RabbitMQ) *OrcestartionService {
	return &OrcestartionService{
		mongo:    db,
		cfg:      cfg,
		rabbitMQ: rabbitMQ,
	}
}

func (os *OrcestartionService) CreateNewOrcestartion(c *gin.Context) {
	var payload models.Orcestartion
	err := c.ShouldBindJSON(&payload)
	if HandleError(c, err, "Unable to unmarshal payload") {
		return
	}
	repo := repository.NewGenericRepository[models.Orcestartion](c.Request.Context(), os.mongo.Database, "orcestartions")
	payload.NIMB_ID = utils.GenerateNIMBID("N_ORC")
	payload.Audit.SetInitialAudit(c)
	payload.Audit.Version, _ = GetNextVersionNumber(c, os.mongo.Database, payload.NIMB_ID)
	_, err = repo.InsertOne(payload)
	if HandleError(c, err, "Failed to create orcestartion") {
		return
	}
	RespondJSON(c, 201, "success", "Orcestartion created", payload)
}

func (os *OrcestartionService) UpdateOrcestartion(c *gin.Context) {
	var payload models.Orcestartion
	err := c.ShouldBindJSON(&payload)
	if HandleError(c, err, "Unable to unmarshel payload") {
		return
	}
	// Archive Old Version and Create New Version
	repo := repository.NewGenericRepository[models.Orcestartion](c.Request.Context(), os.mongo.Database, "orcestartions")
	err = ArchiveEntity(c, repo, payload.NIMB_ID, payload.Audit.Version)
	if HandleError(c, err, "Failed to archive old version of orcestartion") {
		return
	}
	payload.Audit.SetModifiedAudit(c)
	_, err = repo.InsertOne(payload)
	if HandleError(c, err, "Failed to create new version of orcestartion") {
		return
	}
	RespondJSON(c, 201, "success", "Orcestartion updated", payload)
}

func (os *OrcestartionService) GetOrcestartionByNIMBID(c *gin.Context) {
	nimbID := c.Query("nimb_id")
	versionStr := c.Query("version")
	version, _ := strconv.Atoi(versionStr)
	repo := repository.NewGenericRepository[models.Orcestartion](c.Request.Context(), os.mongo.Database, "orcestartions")
	orc, err := repo.FindOne(map[string]interface{}{"nimb_id": nimbID, "audit.is_archived": false, "audit.version": version})
	if HandleError(c, err, "Failed to fetch orcestartion") {
		return
	}
	RespondJSON(c, 200, "success", "Orcestartion retrieved", orc)
}

func (os *OrcestartionService) GetAllOrcestartions(c *gin.Context) {
	repo := repository.NewGenericRepository[models.Orcestartion](c.Request.Context(), os.mongo.Database, "orcestartions")
	option := GetCommonSortOption()
	option.SetProjection(bson.M{"workflow": 0})
	orcestartions, err := repo.FindMany(bson.M{"audit.is_archived": false}, option)
	if HandleError(c, err, "Failed to fetch orcestartions") {
		return
	}
	RespondJSON(c, 200, "success", "Orcestartions retrieved", orcestartions)
}

func (os *OrcestartionService) ArchiveOrcestartion(c *gin.Context) {
	repo := repository.NewGenericRepository[models.Orcestartion](c.Request.Context(), os.mongo.Database, "orcestartions")
	versionStr := c.Query("version")
	version, _ := strconv.Atoi(versionStr)
	err := ArchiveEntity(c, repo, c.Query("nimb_id"), version)
	if HandleError(c, err, "Failed to archive orcestartion") {
		return
	}
	RespondJSON(c, 200, "success", "Orcestartion archived", nil)
}

func (os *OrcestartionService) CloneOrcestartion(c *gin.Context) {
	nimbID := c.Query("nimb_id")
	versionStr := c.Query("version")
	version, _ := strconv.Atoi(versionStr)
	repo := repository.NewGenericRepository[models.Orcestartion](c.Request.Context(), os.mongo.Database, "orcestartions")
	origOrc, err := repo.FindOne(map[string]interface{}{"nimb_id": nimbID, "audit.is_archived": false, "audit.version": version})
	if HandleError(c, err, "Failed to fetch original orcestartion") {
		return
	}
	origOrc.Audit.SetInitialAudit(c)
	origOrc.Audit.Version, _ = GetNextVersionNumber(c, os.mongo.Database, origOrc.NIMB_ID)
	_, err = repo.InsertOne(*origOrc)
	if HandleError(c, err, "Failed to create orcestartion") {
		return
	}
	RespondJSON(c, 201, "success", "Orcestartion cloned", origOrc)
}

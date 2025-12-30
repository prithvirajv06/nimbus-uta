package service

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/prithvirajv06/nimbus-uta/go/core/internal/models"
	"github.com/prithvirajv06/nimbus-uta/go/core/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ArchiveEntity[T any](c *gin.Context, repo *repository.GenericRepository[T], nimbID string, version int) error {
	userID := c.GetHeader("user_id")
	slog.Info("Archiving entity", "nimbID", nimbID, "userID", userID)
	archiveResult, err := repo.UpdateOne(bson.M{"nimb_id": nimbID, "audit.is_archived": false, "audit.version": version},
		bson.M{"$set": bson.M{"audit.is_archived": true}})
	if err != nil || archiveResult.MatchedCount == 0 {
		return err
	}
	return nil
}

func RespondJSON(c *gin.Context, statusCode int, status, message string, data interface{}) {
	slog.Info("Responding JSON", "status", status, "message", message)
	c.JSON(statusCode, models.ApiResponse{
		Status:  status,
		Message: message,
		Data:    data,
	})
}

func HandleError(c *gin.Context, err error, message string) bool {
	slog.Info("Looking for error")
	if err != nil {
		slog.Error("error occurred", "error", err)
		RespondJSON(c, 500, "failure", message, nil)
		return true
	}
	return false
}

func GetNextVersionNumber(c *gin.Context, db *mongo.Database, nimbID string) (int, error) {
	slog.Info("Getting next version number for", "nimbID", nimbID)
	repo := repository.NewGenericRepository[models.IdVersionMapping](c.Request.Context(), db, "id_version_mappings")
	reposult, err := repo.FindOne(bson.M{"nimb_id": nimbID})
	if err != nil {
		repo.InsertOne(models.IdVersionMapping{NIMB_ID: nimbID, NextVersion: 2})
		return 1, err
	}
	reposult.NextVersion += 1
	_, err = repo.UpdateOne(bson.M{"nimb_id": nimbID}, bson.M{"$set": bson.M{"next_version": reposult.NextVersion}})
	return reposult.NextVersion, nil
}

func GetCommonSortOption() *options.FindOptions {
	return options.Find().SetSort(bson.D{{Key: "audit.version", Value: -1}, {Key: "audit.minor_version", Value: -1}, {Key: "audit.updated_at", Value: -1}})
}

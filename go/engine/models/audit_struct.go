package models

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Audit struct {
	CreatedAt       time.Time `bson:"created_at" json:"created_at"`
	CreatedBy       string    `bson:"created_by" json:"created_by"`
	ModifiedAt      time.Time `bson:"modified_at" json:"modified_at"`
	ModifiedBy      string    `bson:"modified_by" json:"modified_by"`
	IsProdCandidate bool      `bson:"is_prod_candidate" json:"is_prod_candidate"`

	//Status fields
	Active         string `bson:"active" json:"active"`
	IsArchived     bool   `bson:"is_archived" json:"is_archived,omitempty"`
	RestoreArchive bool   `bson:"-" json:"restore_archive,omitempty"`

	//Versioning fields
	Version      int `bson:"version" json:"version,omitempty"`
	MinorVersion int `bson:"minor_version" json:"minor_version,omitempty"`
}

func (a *Audit) SetInitialAudit(c *gin.Context) {
	userID := c.GetHeader("user_id")
	if userID == "" {
		userID = "SYSTEM"
	}
	a.CreatedBy = userID
	a.ModifiedBy = userID
	a.CreatedAt = time.Now()
	a.ModifiedAt = time.Now()
	a.Version = 1
	a.MinorVersion = 1
}

func (a *Audit) SetModifiedAudit(c *gin.Context) {
	userID := c.GetHeader("user_id")
	if userID == "" {
		userID = "SYSTEM"
	}
	a.ModifiedBy = userID
	a.ModifiedAt = time.Now()
	a.MinorVersion += 1
}

func (a *Audit) SetRestoreArchive(c *gin.Context) {
	userID := c.GetHeader("user_id")
	if userID == "" {
		userID = "SYSTEM"
	}
	a.ModifiedBy = userID
	a.ModifiedAt = time.Now()
	a.IsArchived = false
}

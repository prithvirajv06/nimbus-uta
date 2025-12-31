package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/prithvirajv06/nimbus-uta/go/core/internal/service"
)

type DecisionTableHandler struct {
	service *service.DecisionTableService
}

func NewDecisionTableHandler(svc *service.DecisionTableService) *DecisionTableHandler {
	return &DecisionTableHandler{
		service: svc,
	}
}

func (h *DecisionTableHandler) RegisterRoutes(router gin.IRoutes) {
	router.POST("/decision-tables", h.service.CreateNewDecisionTable)
	router.PUT("/decision-tables", h.service.UpdateDecisionTable)
	router.GET("/decision-table", h.service.GetDecisionTableByNIMBID)
	router.POST("/decision-tables/list", h.service.GetAllDecisionTables)
	router.DELETE("/decision-table", h.service.ArchiveDecisionTable)
	router.PUT("/decision-table/clone", h.service.CloneDecisionTable)
	//TODO: Add routes to export/import decision table
	//router.GET("/decision-table/export", h.service.ExportDecisionTable)
	//router.POST("/decision-table/import", h.service.ImportDecisionTable)
}

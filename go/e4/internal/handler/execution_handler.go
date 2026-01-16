package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/prithvirajv06/nimbus-uta/go/engine/internal/service"
)

type ExecutionHandler struct {
	service *service.ExecutionService
}

func NewExecutionHandler(svc *service.ExecutionService) *ExecutionHandler {
	return &ExecutionHandler{
		service: svc,
	}
}

func (h *ExecutionHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/engines", h.service.CreateNewEngine)
	router.PUT("/engines", h.service.UpdateEngine)
	router.GET("/engine", h.service.GetEngineByNIMBID)
	router.POST("/engines/list", h.service.GetAllEngines)
	router.DELETE("/engine", h.service.ArchiveEngine)
	router.PUT("/engine/clone", h.service.CloneEngine)
	router.POST("/engine/logic-flow", h.service.HandleRuleExecution)
}

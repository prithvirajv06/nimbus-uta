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
	router.POST("/execute/logic-flow", h.service.HandleRuleExecution)
}

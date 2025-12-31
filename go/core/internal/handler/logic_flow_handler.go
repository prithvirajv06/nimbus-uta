package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/prithvirajv06/nimbus-uta/go/core/internal/service"
)

type LogicFlowHandler struct {
	service *service.LogicFlowService
}

func NewLogicFlowHandler(svc *service.LogicFlowService) *LogicFlowHandler {
	return &LogicFlowHandler{
		service: svc,
	}
}

func (h *LogicFlowHandler) RegisterRoutes(router gin.IRoutes) {
	router.POST("/logic-flows", h.service.CreateNewLogicFlow)
	router.PUT("/logic-flows", h.service.UpdateLogicFlow)
	router.GET("/logic-flow", h.service.GetLogicFlowByNIMBID)
	router.POST("/logic-flows/list", h.service.GetAllLogicFlows)
	router.DELETE("/logic-flow", h.service.ArchiveLogicFlow)
	router.PUT("/logic-flow/clone", h.service.CloneLogicFlow)
}

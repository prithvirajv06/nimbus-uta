package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/prithvirajv06/nimbus-uta/go/core/internal/service"
)

type VariableHandler struct {
	service *service.VariablePackageService
}

func NewVariableHandler(svc *service.VariablePackageService) *VariableHandler {
	return &VariableHandler{
		service: svc,
	}
}

func (h *VariableHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/variables/packages", h.service.CreateNewVariablePackageFromJSON)
	router.PUT("/variables/packages", h.service.UpdateVariablePackage)
	router.GET("/variables/package", h.service.GetVariablePackageByNIMBID)
	router.POST("/variables/packages/list", h.service.GetAllVariablePackages)
	router.DELETE("/variables/package", h.service.ArchiveVariablePackage)
}

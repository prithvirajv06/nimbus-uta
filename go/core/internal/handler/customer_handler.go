package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/prithvirajv06/nimbus-uta/go/core/internal/service"
)

type CustomerHandler struct {
	service *service.CustomerService
}

func NewCustomerHandler(svc *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{
		service: svc,
	}
}

func (h *CustomerHandler) RegisterRoutes(router *gin.RouterGroup) {
	// Define your routes here, e.g.:
	router.POST("/customers/register", h.service.RegisterNewUser)
	router.POST("/customers/login", h.service.LoginUser)
}

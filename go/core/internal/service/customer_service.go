package service

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prithvirajv06/nimbus-uta/go/core/config"
	"github.com/prithvirajv06/nimbus-uta/go/core/internal/models"
	"github.com/prithvirajv06/nimbus-uta/go/core/internal/repository"
	"github.com/prithvirajv06/nimbus-uta/go/core/pkg/database"
	"github.com/prithvirajv06/nimbus-uta/go/core/pkg/messaging"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type CustomerService struct {
	mongo    *database.MongoDB
	rabbitMQ *messaging.RabbitMQ
	cfg      *config.Config
}

func NewCustomerService(db *database.MongoDB, rabbitMQ *messaging.RabbitMQ, cfg *config.Config) *CustomerService {
	rabbitMQ.DeclareQueue("user.event", true)
	return &CustomerService{
		mongo:    db,
		rabbitMQ: rabbitMQ,
		cfg:      cfg,
	}
}

func (s *CustomerService) RegisterNewUser(c *gin.Context) {
	var payload models.User
	if err := c.ShouldBindJSON(&payload); err != nil {
		return
	}

	repo := repository.NewGenericRepository[models.User](c.Request.Context(), s.mongo.Database, "users")

	// Check if user with the same email already exists in the organization
	existingUser, _ := repo.FindOne(bson.M{"email": payload.Email, "organization": payload.Organization})
	if existingUser != nil {
		apiReponse := models.ApiResponse{
			Status:  "failure",
			Message: "User with this email already exists in the organization",
		}
		c.JSON(http.StatusConflict, apiReponse)
		return
	}

	//Create Organization if not exists
	orgRepo := repository.NewGenericRepository[models.Organization](c.Request.Context(), s.mongo.Database, "organizations")
	existingOrg, _ := orgRepo.FindOne(bson.M{"name": payload.Organization.Name})
	if existingOrg == nil {
		newOrg := models.Organization{
			Name:   payload.Organization.Name,
			Active: true,
			Audit:  models.Audit{},
		}
		_, err := orgRepo.InsertOne(newOrg)
		if err != nil {
			apiResponse := models.ApiResponse{
				Status:  "failure",
				Message: "Failed to create organization",
			}
			c.JSON(http.StatusInternalServerError, apiResponse)
			return
		}
		payload.Organization.ID = newOrg.ID
	}

	// Insert new user

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	payload.Password = string(hashedPassword)
	payload.Active = true
	_, err = repo.InsertOne(payload)
	if err != nil {
		apiResponse := models.ApiResponse{
			Status:  "failure",
			Message: "Failed to create user",
		}
		c.JSON(http.StatusInternalServerError, apiResponse)
		return
	}
	apiReponse := models.ApiResponse{
		Status:  "success",
		Message: "User created successfully",
	}
	var notificationPayload models.MessageEvent = models.MessageEvent{
		EventType: "USER_CREATED",
		Payload:   payload,
		Timestamp: int64(time.Now().Unix()),
	}
	s.rabbitMQ.Publish(c.Request.Context(), payload.ID, notificationPayload)
	c.JSON(http.StatusCreated, apiReponse)
}

func (s *CustomerService) LoginUser(c *gin.Context) {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		return
	}
	repo := repository.NewGenericRepository[models.User](c.Request.Context(), s.mongo.Database, "users")
	user, err := repo.FindOne(bson.M{"email": payload.Email})
	if err != nil || user == nil {
		apiResponse := models.ApiResponse{
			Status:  "failure",
			Message: "Invalid email or password",
		}
		c.JSON(http.StatusUnauthorized, apiResponse)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		apiResponse := models.ApiResponse{
			Status:  "failure",
			Message: "Invalid email or password",
		}
		c.JSON(http.StatusUnauthorized, apiResponse)
		return
	}
	apiResponse := models.ApiResponse{
		Status:  "success",
		Message: "Login successful",
	}
	c.JSON(http.StatusOK, apiResponse)
}

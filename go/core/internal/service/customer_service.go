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
	err := c.ShouldBindJSON(&payload)
	if HandleError(c, err, "Unabe unmarshel JSON") {
		return
	}

	repo := repository.NewGenericRepository[models.User](c.Request.Context(), s.mongo.Database, "users")

	// Check if user with the same email already exists in the organization
	_, err = repo.FindOne(bson.M{"email": payload.Email, "organization": payload.Organization})
	if HandleError(c, err, "User with this email already exists in the organization") {
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	payload.Password = string(hashedPassword)
	defaultRole := models.CreateDefaultAdminRole()
	newUser := models.NewUser(payload.Fname, payload.Lname, payload.Email, payload.Password, payload.Organization.Name, defaultRole)
	payload = *newUser

	//Create Organization if not exists
	orgRepo := repository.NewGenericRepository[models.Organization](c.Request.Context(), s.mongo.Database, "organizations")
	existingOrg, _ := orgRepo.FindOne(bson.M{"name": payload.Organization.Name})
	if existingOrg == nil {
		newOrg := models.NewOrganization(payload.Organization.Name, payload.Organization.Address)
		_, err := orgRepo.InsertOne(*newOrg)
		if HandleError(c, err, "Failed to create organization") {
			return
		}
		payload.Organization = *newOrg
	}

	// Insert new user
	_, err = repo.InsertOne(payload)
	if HandleError(c, err, "Failed to create user") {
		return
	}
	var notificationPayload models.MessageEvent = models.MessageEvent{
		EventType: "USER_CREATED",
		Payload:   payload,
		Timestamp: int64(time.Now().Unix()),
	}
	s.rabbitMQ.Publish(c.Request.Context(), payload.NIMB_ID, notificationPayload)
	RespondJSON(c, http.StatusCreated, "success", "User created successfully", payload)
}

/*
*
LoginUser handles user login requests. It verifies the provided email and password,
and responds with a success message if the credentials are valid.
*/
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
	if HandleError(c, err, "Invalid email or password") {
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if HandleError(c, err, "Invalid email or password") {
		return
	}
	RespondJSON(c, http.StatusOK, "success", "Login successful", nil)
}

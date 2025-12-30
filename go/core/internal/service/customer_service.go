package service

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
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
	result, err := repo.FindOne(bson.M{"email": payload.Email, "organization": payload.Organization})
	if result != nil {
		RespondJSON(c, http.StatusConflict, "failure", "User with the same email already exists in the organization", nil)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	defaultRole := models.CreateDefaultAdminRole()
	newUser := models.NewUser(c, payload.Fname, payload.Lname, payload.Email, payload.Organization.Name, defaultRole)
	newUser.Password = string(hashedPassword)
	//Create Organization if not exists
	orgRepo := repository.NewGenericRepository[models.Organization](c.Request.Context(), s.mongo.Database, "organizations")
	existingOrg, _ := orgRepo.FindOne(bson.M{"name": payload.Organization.Name})
	if existingOrg == nil {
		newOrg := models.NewOrganization(c, payload.Organization.Name, payload.Organization.Address)
		_, err := orgRepo.InsertOne(newOrg)
		if HandleError(c, err, "Failed to create organization") {
			return
		}
		newUser.Organization = newOrg
	}

	// Insert new user
	_, err = repo.InsertOne(newUser)
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
	user.JWTToken, err = GenerateJWTToken(s.cfg.Auth.JWTToken, user.NIMB_ID, user.Email, user.Role.Name, s.cfg.Auth.JWTTokenExpiryHours)
	if HandleError(c, err, "Failed to generate JWT token") {
		return
	}
	RespondJSON(c, http.StatusOK, "success", "Login successful", user)
}

func GenerateJWTToken(secret, nimbID, email, role string, expiryHours int) (string, error) {
	claims := jwt.MapClaims{
		"nimb_id": nimbID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(time.Duration(expiryHours) * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

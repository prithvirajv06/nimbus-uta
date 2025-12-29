package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/prithvirajv06/nimbus-uta/go/notification/config"
	"github.com/prithvirajv06/nimbus-uta/go/notification/internal/service"
	"github.com/prithvirajv06/nimbus-uta/go/notification/pkg/cache"
	"github.com/prithvirajv06/nimbus-uta/go/notification/pkg/database"
	"github.com/prithvirajv06/nimbus-uta/go/notification/pkg/messaging"
)

func main() {
	file, err := os.OpenFile("application.json", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer file.Close()
	logger := slog.New(slog.NewJSONHandler(file, nil))
	slog.SetDefault(logger)

	// Load configuration
	cfg := config.Load()

	// Initialize MongoDB
	mongoDB, err := database.NewMongoDB(cfg.MongoDB)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoDB.Disconnect()

	// // Initialize Redis
	redisClient := cache.NewRedisClient(cfg.Redis)
	defer redisClient.Close()

	// // Initialize RabbitMQ
	rabbitMQ, err := messaging.NewRabbitMQ(cfg.RabbitMQ)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()
	startConsumers(rabbitMQ, &service.CustomerService{
		Cfg: cfg,
	})
	select {}
}

func startConsumers(rabbitMQ *messaging.RabbitMQ, customerService *service.CustomerService) {
	// Start user event consumer
	if err := rabbitMQ.Consume("user.event", customerService.HandleUserEvent); err != nil {
		log.Printf("Failed to start consumer: %v", err)
	}
}

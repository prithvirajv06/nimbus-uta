package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prithvirajv06/nimbus-uta/go/core/config"
	"github.com/prithvirajv06/nimbus-uta/go/core/internal/handler"
	"github.com/prithvirajv06/nimbus-uta/go/core/internal/service"
	"github.com/prithvirajv06/nimbus-uta/go/core/pkg/cache"
	"github.com/prithvirajv06/nimbus-uta/go/core/pkg/database"
	"github.com/prithvirajv06/nimbus-uta/go/core/pkg/messaging"
	"github.com/prithvirajv06/nimbus-uta/go/core/pkg/storage"
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

	// Initialize MinIO
	minioClient, err := storage.NewMinIOClient(cfg.MinIO)
	if err != nil {
		log.Fatalf("Failed to connect to MinIO: %v", err)
	}
	_ = minioClient // Use minioClient as needed

	// Set up Gin router and handlers
	router := gin.Default()

	// Customer Service and Handler
	customerService := service.NewCustomerService(mongoDB, rabbitMQ, cfg)
	customerHandler := handler.NewCustomerHandler(customerService)
	customerHandler.RegisterRoutes(router)

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// Start server in goroutine
	go func() {
		slog.Info("Starting server on port " + cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed: " + err.Error())
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown: " + err.Error())
	}

	slog.Info("Server exited")
}

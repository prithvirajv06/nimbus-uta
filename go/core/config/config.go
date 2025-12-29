// config/config.go
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server    ServerConfig
	MongoDB   MongoDBConfig
	Redis     RedisConfig
	RabbitMQ  RabbitMQConfig
	MinIO     MinIOConfig
	ZenEngine ZenEngineConfig
	Auth      AuthConfig
	PYWorker  string
}

type AuthConfig struct {
	JWTToken string
}
type ServerConfig struct {
	Port string
	Env  string
}

type MongoDBConfig struct {
	URI      string
	Database string
	Timeout  int
	Username string
	Password string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type RabbitMQConfig struct {
	URL      string
	Exchange string
	Queue    string
}

type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

type ZenEngineConfig struct {
	RulesPath string
}

func Load() *Config {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		Auth: AuthConfig{
			JWTToken: getEnv("JWT_TOKEN", "ABCDSASFAFHJKEHRJKHESKFHSIUIOASUDKLSAJDKKFJDKLJ"),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Env:  getEnv("SERVER_ENV", "development"),
		},
		MongoDB: MongoDBConfig{
			URI:      getEnv("MONGODB_URI", "mongodb://localhost:27017"),
			Database: getEnv("MONGODB_DATABASE", "myapp"),
			Timeout:  10,
			Username: getEnv("MONGODB_USERNAME", ""),
			Password: getEnv("MONGODB_PASSWORD", ""),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       0,
		},
		RabbitMQ: RabbitMQConfig{
			URL:      getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
			Exchange: getEnv("RABBITMQ_EXCHANGE", "main_exchange"),
			Queue:    getEnv("RABBITMQ_QUEUE", "main_queue"),
		},
		MinIO: MinIOConfig{
			Endpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
			AccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
			SecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
			Bucket:    getEnv("MINIO_BUCKET", "uploads"),
			UseSSL:    false,
		},
		ZenEngine: ZenEngineConfig{
			RulesPath: getEnv("ZEN_RULES_PATH", "./rules"),
		},
		PYWorker: getEnv("PYTHON_WORKER", "localhost:5005"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

package database

import (
	"context"
	"time"

	"github.com/prithvirajv06/nimbus-uta/go/engine/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Mongo DB cient to make use in repositories

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func NewMongoDB(cfg config.MongoDBConfig) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.URI)
	if cfg.Database != "" || cfg.Password != "" {
		clientOptions.SetAuth(options.Credential{
			Username: cfg.Username,
			Password: cfg.Password,
		})
	}
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return &MongoDB{
		Client:   client,
		Database: client.Database(cfg.Database),
	}, nil
}

func (m *MongoDB) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return m.Client.Disconnect(ctx)
}

func (m *MongoDB) Collection(name string) *mongo.Collection {
	return m.Database.Collection(name)
}

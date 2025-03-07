// setup_test.go
package tests

import (
	"context"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
)

func init() {
	// Load environment variables and connect to MongoDB for all tests

	projectRoot := filepath.Join("..", ".env")
	err := godotenv.Load(projectRoot)
	if err != nil {
		logrus.Fatalf("Error loading .env file: %v", err)
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		logrus.Fatalf("MONGO_URI not set in .env file")
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logrus.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	if err := client.Ping(context.TODO(), nil); err != nil {
		logrus.Fatalf("Failed to ping MongoDB: %v", err)
	}
}

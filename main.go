package main

import (
	"context"
	"os"

	"omhs-backend/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client // Declare client here

func main() {
	// Set up logger
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
	logrus.Info("Starting server...")

	// Load environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		logrus.Fatalf("Error loading .env file: %v", err)
	}

	// Retrieve the MongoDB URI from environment variables
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		logrus.Fatal("MONGO_URI not set in .env file")
	}

	// Set client options
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Connect to MongoDB
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logrus.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		logrus.Fatalf("Failed to ping MongoDB: %v", err)
	}

	logrus.Info("Connected to MongoDB Atlas!")

	// Set up Gin router
	r := gin.Default()

	// Set trusted proxies to only localhost
	err = r.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		logrus.Fatalf("Failed to set trusted proxies: %v", err)
	}

	// CORS setup
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Initialize authentication routes
	controllers.InitializeAuthRoutes(r, client)

	// Initialize other routes
	controllers.InitializeRoutes(r, client)

	// Start the server
	if err := r.Run(":8080"); err != nil {
		logrus.Fatalf("Failed to run server: %v", err)
	}
}

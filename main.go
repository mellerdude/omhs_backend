package main

import (
	"context"
	"errors"
	"os"

	"omhs-backend/controllers"
	"omhs-backend/utils"

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

	// Get the singleton instance of ProjectManager
	projectManager := utils.NewProjectManager()

	// Load environment variables from the .env file
	projectManager.Execute(func() (interface{}, error) { return nil, godotenv.Load() }, "Error loading .env file")

	// Retrieve the MongoDB URI from environment variables
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		projectManager.Execute(func() (interface{}, error) { return nil, errors.New("MONGO_URI not set") }, "MONGO_URI not set in .env file")
	}

	// Set client options
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Connect to MongoDB
	projectManager.Execute(func() (interface{}, error) {
		var err error
		client, err = mongo.Connect(context.TODO(), clientOptions)
		return client, err
	}, "Failed to connect to MongoDB")

	// Check the connection
	projectManager.Execute(func() (interface{}, error) { return nil, client.Ping(context.TODO(), nil) }, "Failed to ping MongoDB")

	logrus.Info("Connected to MongoDB Atlas!")

	// Set up Gin router
	r := gin.Default()

	// Set trusted proxies to only localhost
	projectManager.Execute(func() (interface{}, error) { return nil, r.SetTrustedProxies([]string{"127.0.0.1"}) }, "Failed to set trusted proxies")

	// CORS setup
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	authController := controllers.NewAuthController(projectManager)
	requestController := controllers.NewRequestController(projectManager)

	// Initialize authentication routes
	controllers.InitializeAuthRoutes(r, client, authController)

	// Initialize other routes
	controllers.InitializeRequestRoutes(r, client, requestController)

	// Start the server
	projectManager.Execute(func() (interface{}, error) { return nil, r.Run(":8080") }, "Failed to run server")
}

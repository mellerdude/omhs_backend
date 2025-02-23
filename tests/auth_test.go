package tests

import (
	"context"
	"os"
	"path/filepath"

	"omhs-backend/controllers"
	"omhs-backend/utils"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client         *mongo.Client
	projectManager *utils.ProjectManager
)

func init() {
	projectManager = utils.NewProjectManager()

	// Load environment variables from the .env file in the parent directory
	projectRoot := filepath.Join("..", ".env")
	err := godotenv.Load(projectRoot)
	projectManager.Execute(func() (interface{}, error) { return nil, godotenv.Load(projectRoot) }, "Error loading .env file")
	if err != nil {
		logrus.Fatalf("Error loading .env file: %v", err)
	}

	// Retrieve the MongoDB URI from environment variables
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		logrus.Fatalf("MONGO_URI not set in .env file")
	}

	// Initialize MongoDB client
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logrus.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ping the primary
	if err := client.Ping(context.TODO(), nil); err != nil {
		logrus.Fatalf("Failed to ping MongoDB: %v", err)
	}
}

func TestRegister(t *testing.T) {
	// Initialize router and controllers
	router := gin.Default()
	pm := utils.NewProjectManager()
	authController := controllers.NewAuthController(pm)
	requestController := controllers.NewRequestController(pm)
	controllers.InitializeAuthRoutes(router, client, authController)
	controllers.InitializeRequestRoutes(router, client, requestController)

	// Test data
	user := map[string]string{
		"username": os.Getenv("NON_ADMIN_USER"),
		"password": os.Getenv("NON_ADMIN_PASS"),
		"email":    os.Getenv("EMAIL_USER"),
	}

	// Register user
	registeredUser := RegisterUser(router, user, t)

	// Login user
	LoginUser(router, user["username"], user["password"], t)

	// Admin login
	adminToken := LoginUser(router, os.Getenv("ADMIN_USER"), os.Getenv("ADMIN_PASS"), t)

	// Delete registered user
	DeleteUser(router, registeredUser.ID.Hex(), adminToken, t)
}

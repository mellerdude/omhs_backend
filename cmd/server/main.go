package main

import (
	"context"
	"errors"
	"os"

	"omhs-backend/internal/auth"
	"omhs-backend/internal/requests"
	"omhs-backend/internal/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
	logrus.Info("Starting server...")

	pm := utils.NewProjectManager()

	pm.Execute(func() error { return godotenv.Load() }, "Error loading .env file")

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		pm.Execute(func() error { return errors.New("MONGO_URI not set") }, "Missing MONGO_URI in .env file")
	}

	client := connectMongo(mongoURI, pm)
	logrus.Info("Connected to MongoDB!")

	r := gin.Default()
	setupCORS(r, pm)
	initRoutes(r, client)

	pm.Execute(func() error { return r.Run(":8080") }, "Failed to run server")
}

func connectMongo(uri string, pm *utils.ProjectManager) *mongo.Client {
	clientOptions := options.Client().ApplyURI(uri)
	var client *mongo.Client

	pm.Execute(func() error {
		var err error
		client, err = mongo.Connect(context.TODO(), clientOptions)
		return err
	}, "Failed to connect to MongoDB")

	pm.Execute(func() error { return client.Ping(context.TODO(), nil) }, "Failed to ping MongoDB")

	return client
}

func setupCORS(r *gin.Engine, pm *utils.ProjectManager) {
	pm.Execute(func() error { return r.SetTrustedProxies([]string{"127.0.0.1"}) }, "Failed to set trusted proxies")

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
}

func initRoutes(r *gin.Engine, client *mongo.Client) {
	authRepo := auth.NewMongoUserRepository(client)
	authService := auth.NewAuthService(authRepo)
	authController := auth.NewAuthController(authService)
	auth.RegisterRoutes(r, authController)

	reqRepo := requests.NewMongoRequestRepository(client)
	reqService := requests.NewRequestService(reqRepo)
	reqController := requests.NewRequestController(reqService)
	requests.RegisterRoutes(r, reqController)
}

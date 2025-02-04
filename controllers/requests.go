package controllers

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"omhs-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var requestClient *mongo.Client

// InitializeRoutes sets up the CRUD routes for documents
// Parameters:
// - r: *gin.Engine - the Gin engine to which the routes are added
// - dbClient: *mongo.Client - the MongoDB client to be used for database operations
func InitializeRoutes(r *gin.Engine, dbClient *mongo.Client) {
	requestClient = dbClient
	r.POST("/:database/:collection", createDocument)
	r.GET("/:database/:collection/:id", getDocument)
	r.PUT("/:database/:collection/:id", updateDocument)
	r.DELETE("/:database/:collection/:id", deleteDocument)
	r.GET("/:database/:collection", getAllDocuments) // Route for fetching all documents
}

// authenticateToken validates the token in the Authorization header
// Parameters:
// - c: *gin.Context - the Gin context for the request
// Returns:
// - *models.User: the authenticated user
// - error: any error encountered during token validation
func authenticateToken(c *gin.Context) (*models.User, error) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		logrus.Warn("Missing Authorization header")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
		c.Abort()
		return nil, errors.New("missing Authorization header")
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	var user models.User
	err := requestClient.Database("users").Collection("authentication").FindOne(context.TODO(), bson.M{"token": tokenString}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		logrus.Warn("Invalid or expired token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		c.Abort()
		return nil, errors.New("invalid or expired token")
	} else if err != nil {
		logrus.Errorf("Failed to validate token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate token", "details": err.Error()})
		c.Abort()
		return nil, err
	}

	return &user, nil
}

// createDocument handles creating a new document
// Parameters:
// - c: *gin.Context - the Gin context for the request
//
// Route: /:database/:collection
// Method: POST
// Expected JSON Body:
//
//	{
//	  "data": {
//	    "name": "Alice",
//	    "age": 30,
//	    "hobbies": ["reading", "hiking"]
//	  }
//	}
func createDocument(c *gin.Context) {
	logrus.Info("Creating a new document")

	// Authenticate the request
	authenticatedUser, err := authenticateToken(c)
	if err != nil {
		return
	}

	// Check if the authenticated user is an admin
	if !authenticatedUser.IsAdmin {
		logrus.Warnf("Unauthorized attempt to create document by user: %s", authenticatedUser.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	databaseName := c.Param("database")
	collectionName := c.Param("collection")
	collection := requestClient.Database(databaseName).Collection(collectionName)

	var doc models.Document
	if err := c.ShouldBindJSON(&doc); err != nil {
		logrus.Errorf("Invalid JSON format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doc.ID = primitive.NewObjectID()
	_, err = collection.InsertOne(context.TODO(), doc)
	if err != nil {
		logrus.Errorf("Failed to create document: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logrus.Info("Document created successfully")
	c.JSON(http.StatusOK, doc)
}

// getDocument handles retrieving a document by ID
// Parameters:
// - c: *gin.Context - the Gin context for the request
//
// Route: /:database/:collection/:id
// Method: GET
func getDocument(c *gin.Context) {
	logrus.Info("Retrieving a document by ID")

	if _, err := authenticateToken(c); err != nil {
		return
	}

	databaseName := c.Param("database")
	collectionName := c.Param("collection")
	collection := requestClient.Database(databaseName).Collection(collectionName)

	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	var doc models.Document
	err := collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&doc)
	if err != nil {
		logrus.Errorf("Failed to retrieve document: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logrus.Info("Document retrieved successfully")
	c.JSON(http.StatusOK, doc)
}

// updateDocument handles updating a document by ID
// Parameters:
// - c: *gin.Context - the Gin context for the request
//
// Route: /:database/:collection/:id
// Method: PUT
// Expected JSON Body:
//
//	{
//	  "data": {
//	    "name": "Alice",
//	    "age": 31,
//	    "hobbies": ["reading", "hiking", "cooking"]
//	  }
//	}
func updateDocument(c *gin.Context) {
	logrus.Info("Updating a document by ID")

	// Authenticate the request
	authenticatedUser, err := authenticateToken(c)
	if err != nil {
		return
	}

	// Check if the authenticated user is an admin
	if !authenticatedUser.IsAdmin {
		logrus.Warnf("Unauthorized attempt to update document by user: %s", authenticatedUser.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	databaseName := c.Param("database")
	collectionName := c.Param("collection")
	collection := requestClient.Database(databaseName).Collection(collectionName)

	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	var doc models.Document
	if err := c.ShouldBindJSON(&doc); err != nil {
		logrus.Errorf("Invalid JSON format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": id}, bson.D{{"$set", doc.Data}})
	if err != nil {
		logrus.Errorf("Failed to update document: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logrus.Info("Document updated successfully")
	c.JSON(http.StatusOK, doc)
}

// deleteDocument handles deleting a document by ID
// Parameters:
// - c: *gin.Context - the Gin context for the request
//
// Route: /:database/:collection/:id
// Method: DELETE
func deleteDocument(c *gin.Context) {
	logrus.Info("Deleting a document by ID")

	// Authenticate the request
	authenticatedUser, err := authenticateToken(c)
	if err != nil {
		return
	}

	// Check if the authenticated user is an admin
	if !authenticatedUser.IsAdmin {
		logrus.Warnf("Unauthorized attempt to delete document by user: %s", authenticatedUser.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	databaseName := c.Param("database")
	collectionName := c.Param("collection")
	collection := requestClient.Database(databaseName).Collection(collectionName)

	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	_, err = collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		logrus.Errorf("Failed to delete document: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logrus.Info("Document deleted successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Document deleted"})
}

// getAllDocuments handles retrieving all documents from a collection
// Parameters:
// - c: *gin.Context - the Gin context for the request
//
// Route: /:database/:collection
// Method: GET
func getAllDocuments(c *gin.Context) {
	logrus.Info("Retrieving all documents from a collection")

	// Authenticate the request
	authenticatedUser, err := authenticateToken(c)
	if err != nil {
		return
	}

	databaseName := c.Param("database")
	collectionName := c.Param("collection")

	// Check if the path is users/authentication and if the user is not an admin
	if databaseName == "users" && collectionName == "authentication" && !authenticatedUser.IsAdmin {
		logrus.Warnf("Unauthorized attempt to retrieve all documents from users/authentication by user: %s", authenticatedUser.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	collection := requestClient.Database(databaseName).Collection(collectionName)

	var documents []models.Document
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		logrus.Errorf("Failed to retrieve documents: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var doc models.Document
		err := cursor.Decode(&doc)
		if err != nil {
			logrus.Errorf("Failed to decode document: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		documents = append(documents, doc)
	}
	if err := cursor.Err(); err != nil {
		logrus.Errorf("Cursor error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logrus.Info("All documents retrieved successfully")
	c.JSON(http.StatusOK, documents)
}

package controllers

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"omhs-backend/models"
	"omhs-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	requestClient *mongo.Client
)

type RequestController struct {
	pm *utils.ProjectManager
}

func NewRequestController(pm *utils.ProjectManager) *RequestController {
	return &RequestController{pm: pm}
}

// InitializeRequestRoutes sets up the CRUD routes for documents
func InitializeRequestRoutes(r *gin.Engine, dbClient *mongo.Client, requestController *RequestController) {
	requestClient = dbClient
	r.POST("/:database/:collection", requestController.createDocument)
	r.GET("/:database/:collection/:id", requestController.getDocument)
	r.PUT("/:database/:collection/:id", requestController.updateDocument)
	r.DELETE("/:database/:collection/:id", requestController.deleteDocument)
	r.GET("/:database/:collection", requestController.getAllDocuments)
}

// authenticateToken validates the token in the Authorization header
func (rc *RequestController) authenticateToken(c *gin.Context) (*models.User, error) {
	var user *models.User
	rc.pm.Execute(func() (interface{}, error) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			return nil, errors.New("missing Authorization header")
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		err := requestClient.Database("users").Collection("authentication").FindOne(context.TODO(), bson.M{"token": tokenString}).Decode(&user)
		return user, err
	}, "Failed to validate token")
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		c.Abort()
		return nil, errors.New("invalid or expired token")
	}
	return user, nil
}

// createDocument handles creating a new document
func (rc *RequestController) createDocument(c *gin.Context) {
	logrus.Info("Creating a new document")

	authenticatedUser, err := rc.authenticateToken(c)
	if err != nil || !authenticatedUser.IsAdmin {
		logrus.Warnf("Unauthorized attempt to create document by user: %s", authenticatedUser.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	databaseName := c.Param("database")
	collectionName := c.Param("collection")
	collection := requestClient.Database(databaseName).Collection(collectionName)

	var doc models.Document
	rc.pm.Execute(func() (interface{}, error) {
		if err := c.ShouldBindJSON(&doc); err != nil {
			return nil, err
		}
		doc.ID = primitive.NewObjectID()
		return collection.InsertOne(context.TODO(), doc)
	}, "Failed to create document")

	c.JSON(http.StatusOK, doc)
}

// getDocument handles retrieving a document by ID
func (rc *RequestController) getDocument(c *gin.Context) {
	logrus.Info("Retrieving a document by ID")

	authenticatedUser, err := rc.authenticateToken(c)
	if err != nil {
		return
	}

	if !authenticatedUser.IsAdmin {
		logrus.Warnf("Unauthorized attempt to retrieve document by user: %s", authenticatedUser.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	databaseName := c.Param("database")
	collectionName := c.Param("collection")
	collection := requestClient.Database(databaseName).Collection(collectionName)

	var doc bson.M
	rc.pm.Execute(func() (interface{}, error) {
		id, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			logrus.Errorf("Invalid document ID: %v", err)
			return nil, err
		}
		err = collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&doc)
		if err != nil {
			logrus.Errorf("Failed to retrieve document: %v", err)
			return nil, err
		}
		return doc, nil
	}, "Failed to retrieve document")

	logrus.Infof("Document retrieved: %+v", doc)
	c.JSON(http.StatusOK, doc)
}

// updateDocument handles updating a document by ID
func (rc *RequestController) updateDocument(c *gin.Context) {
	logrus.Info("Updating a document by ID")

	authenticatedUser, err := rc.authenticateToken(c)
	if err != nil || !authenticatedUser.IsAdmin {
		logrus.Warnf("Unauthorized attempt to update document by user: %s", authenticatedUser.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	databaseName := c.Param("database")
	collectionName := c.Param("collection")
	collection := requestClient.Database(databaseName).Collection(collectionName)

	var doc models.Document
	rc.pm.Execute(func() (interface{}, error) {
		id, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			return nil, err
		}
		if err := c.ShouldBindJSON(&doc); err != nil {
			return nil, err
		}
		result, err := collection.UpdateOne(context.TODO(), bson.M{"_id": id}, bson.D{{"$set", doc.Data}})
		return result, err
	}, "Failed to update document")

	c.JSON(http.StatusOK, doc)
}

// deleteDocument handles deleting a document by ID
func (rc *RequestController) deleteDocument(c *gin.Context) {
	logrus.Info("Deleting a document by ID")

	authenticatedUser, err := rc.authenticateToken(c)
	if err != nil || !authenticatedUser.IsAdmin {
		logrus.Warnf("Unauthorized attempt to delete document by user: %s", authenticatedUser.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	databaseName := c.Param("database")
	collectionName := c.Param("collection")
	collection := requestClient.Database(databaseName).Collection(collectionName)

	logrus.Infof("Deleting document from database: %s, collection: %s, with ID: %s", databaseName, collectionName, c.Param("id"))

	rc.pm.Execute(func() (interface{}, error) {
		id, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			logrus.Errorf("Invalid document ID: %v", err)
			return nil, err
		}
		result, err := collection.DeleteOne(context.TODO(), bson.M{"_id": id})
		if err != nil {
			logrus.Errorf("Failed to delete document: %v", err)
			return nil, err
		}
		if result.DeletedCount == 0 {
			logrus.Warnf("No document found with ID: %s", c.Param("id"))
			return nil, mongo.ErrNoDocuments
		}
		return result, nil
	}, "Failed to delete document")

	c.JSON(http.StatusOK, gin.H{"message": "Document deleted"})
}

// getAllDocuments handles retrieving all documents from a collection
func (rc *RequestController) getAllDocuments(c *gin.Context) {
	logrus.Info("Retrieving all documents from a collection")

	authenticatedUser, err := rc.authenticateToken(c)
	if err != nil {
		return
	}

	databaseName := c.Param("database")
	collectionName := c.Param("collection")
	if databaseName == "users" && collectionName == "authentication" && !authenticatedUser.IsAdmin {
		logrus.Warnf("Unauthorized attempt to retrieve all documents from users/authentication by user: %s", authenticatedUser.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	collection := requestClient.Database(databaseName).Collection(collectionName)

	var documents []models.Document
	rc.pm.Execute(func() (interface{}, error) {
		cursor, err := collection.Find(context.TODO(), bson.M{})
		if err != nil {
			return nil, err
		}
		defer cursor.Close(context.TODO())

		for cursor.Next(context.TODO()) {
			var doc models.Document
			err := cursor.Decode(&doc)
			if err != nil {
				return nil, err
			}
			documents = append(documents, doc)
		}
		return nil, cursor.Err()
	}, "Failed to retrieve documents")

	c.JSON(http.StatusOK, documents)
}

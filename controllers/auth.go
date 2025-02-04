package controllers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"omhs-backend/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var authClient *mongo.Client // Declare a local variable to store the passed client

// InitializeAuthRoutes sets up the authentication routes
// Parameters:
// - r: *gin.Engine - the Gin engine to which the routes are added
// - dbClient: *mongo.Client - the MongoDB client to be used for database operations
func InitializeAuthRoutes(r *gin.Engine, dbClient *mongo.Client) { // Corrected function signature
	authClient = dbClient // Assign the passed client to the local variable
	r.POST("/register", registerUser)
	r.POST("/login", loginUser)
	r.POST("/change-password", changePassword) // Add route for changing password
	r.POST("/reset-password", resetPassword)   // Add route for resetting password
}

// registerUser handles user registration
// Parameters:
// - c: *gin.Context - the Gin context for the request
//
// Route: /register
// Method: POST
// Expected JSON Body:
//
//	{
//	  "username": "new_user",
//	  "password": "newpassword123",
//	  "email": "user@example.com"
//	}
func registerUser(c *gin.Context) {
	logrus.Info("Attempting to register a new user")

	collection := authClient.Database("users").Collection("authentication") // Use the local client variable

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		logrus.Errorf("Invalid JSON format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format", "details": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Errorf("Failed to hash password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password", "details": err.Error()})
		return
	}
	user.Password = string(hashedPassword)
	user.ID = primitive.NewObjectID()
	user.IsAdmin = false // Default to false

	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		logrus.Errorf("Failed to register user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user", "details": err.Error()})
		return
	}

	logrus.Infof("User registered successfully: %s", user.Username)
	c.JSON(http.StatusCreated, user)
}

// loginUser handles user login
// Parameters:
// - c: *gin.Context - the Gin context for the request
//
// Route: /login
// Method: POST
// Expected JSON Body:
//
//	{
//	  "username": "existing_user",
//	  "password": "existingpassword123"
//	}
func loginUser(c *gin.Context) {
	logrus.Info("Attempting to log in user")

	collection := authClient.Database("users").Collection("authentication") // Use the local client variable

	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		logrus.Errorf("Invalid JSON format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format", "details": err.Error()})
		return
	}

	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{"username": input.Username}).Decode(&user)
	if err == mongo.ErrNoDocuments || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)) != nil {
		logrus.Warnf("Invalid username or password for user: %s", input.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	token := user.Token
	// Check if last login was more than 48 hours ago
	if time.Since(user.LastLogin) > 48*time.Hour {
		token, err = generateToken()
		if err != nil {
			logrus.Errorf("Failed to generate token: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token", "details": err.Error()})
			return
		}
	}

	user.Token = token
	user.LastLogin = time.Now()
	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": user.ID}, bson.D{{"$set", bson.M{"token": token, "lastLogin": user.LastLogin}}})
	if err != nil {
		logrus.Errorf("Failed to update login information: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update login information", "details": err.Error()})
		return
	}

	logrus.Infof("User logged in successfully: %s", user.Username)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// resetPassword handles the generation of a passkey for password reset
// Parameters:
// - c: *gin.Context - the Gin context for the request
//
// Route: /reset-password
// Method: POST
// Expected JSON Body:
//
//	{
//	  "email": "user@example.com"
//	}
func resetPassword(c *gin.Context) {
	logrus.Info("Attempting to handle forget password request")

	type ResetPasswordRequest struct {
		Email string `json:"email"`
	}

	var request ResetPasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logrus.Errorf("Invalid JSON format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format", "details": err.Error()})
		return
	}

	logrus.Infof("Received reset password request for email: %s", request.Email)

	collection := authClient.Database("users").Collection("authentication")

	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{"email": request.Email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		logrus.Warnf("Email not found: %s", request.Email)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email not found"})
		return
	}
	if err != nil {
		logrus.Errorf("Error finding user by email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
		return
	}

	passkey, err := generatePasskey()
	if err != nil {
		logrus.Errorf("Failed to generate passkey: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate passkey", "details": err.Error()})
		return
	}

	passkeyGeneratedAt := time.Now()
	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": user.ID}, bson.D{
		{"$set", bson.M{"passkey": passkey, "passkeyGeneratedAt": passkeyGeneratedAt}},
	})
	if err != nil {
		logrus.Errorf("Failed to save passkey: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save passkey", "details": err.Error()})
		return
	}

	// Start a goroutine to invalidate the passkey after 10 minutes
	go func() {
		time.Sleep(10 * time.Minute)
		_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": user.ID}, bson.D{
			{"$set", bson.M{"passkey": "NOT_PASSKEY"}},
		})
		if err != nil {
			logrus.Errorf("Failed to invalidate passkey: %v", err)
		} else {
			logrus.Infof("Passkey invalidated for user: %s", user.Email)
		}
	}()

	subject := "Password Reset Passkey"
	message := fmt.Sprintf("Your passkey for resetting your password is: %s", passkey)
	err = sendEmail(user.Email, subject, message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email", "details": err.Error()})
		return
	}

	logrus.Infof("Passkey sent successfully to: %s", user.Email)
	c.JSON(http.StatusOK, gin.H{"message": "Passkey sent successfully"})
}

// changePassword handles changing a user's password after passkey verification
// Parameters:
// - c: *gin.Context - the Gin context for the request
//
// Route: /change-password
// Method: POST
// Expected JSON Body:
//
//	{
//	  "email": "user@example.com",
//	  "passkey": "ABC123",
//	  "newPassword": "newpassword123"
//	}
func changePassword(c *gin.Context) {
	logrus.Info("Attempting to change password")

	type ChangePasswordRequest struct {
		Email       string `json:"email"`
		Passkey     string `json:"passkey"`
		NewPassword string `json:"newPassword"`
	}

	var request ChangePasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logrus.Errorf("Invalid JSON format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format", "details": err.Error()})
		return
	}

	collection := authClient.Database("users").Collection("authentication")

	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{"email": request.Email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		logrus.Warnf("Email not found: %s", request.Email)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email not found"})
		return
	}

	// Check if the passkey is still valid
	if !isPasskeyValid(user) {
		logrus.Warnf("Passkey expired for user: %s", request.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Passkey expired"})
		return
	}

	// Verify the passkey
	if user.Passkey != request.Passkey {
		logrus.Warnf("Invalid passkey for user: %s", request.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid passkey"})
		return
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		logrus.Errorf("Failed to hash new password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password", "details": err.Error()})
		return
	}

	// Update the password and clear the passkey
	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": user.ID}, bson.D{
		{"$set", bson.M{"password": hashedPassword}},
		{"$unset", bson.M{"passkey": "", "passkeyGeneratedAt": ""}},
	})
	if err != nil {
		logrus.Errorf("Failed to change password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to change password", "details": err.Error()})
		return
	}

	logrus.Infof("Password changed successfully for user: %s", request.Email)
	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

func isPasskeyValid(user models.User) bool {
	// Check if the passkey is older than 10 minutes
	if time.Since(user.PasskeyGeneratedAt) > 10*time.Minute {
		// Invalidate the passkey if more than 10 minutes have passed
		collection := authClient.Database("users").Collection("authentication")
		_, err := collection.UpdateOne(context.TODO(), bson.M{"_id": user.ID}, bson.D{
			{"$set", bson.M{"passkey": "NOT_PASSKEY"}},
		})
		if err != nil {
			logrus.Errorf("Failed to invalidate passkey: %v", err)
		}
		return false
	}
	return true
}

// generateToken generates a secure token for authentication
// Returns:
// - string: the generated token
// - error: any error encountered during token generation
func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		logrus.Errorf("Failed to generate token: %v", err)
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// generatePasskey generates a 6-letter passkey for password reset
// Returns:
// - string: the generated passkey
// - error: any error encountered during passkey generation
func generatePasskey() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		logrus.Errorf("Failed to generate passkey: %v", err)
		return "", err
	}
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return string(b), nil
}

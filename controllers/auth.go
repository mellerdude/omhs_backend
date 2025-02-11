package controllers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"omhs-backend/models"
	"omhs-backend/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var (
	authClient *mongo.Client
)

type AuthController struct {
	pm *utils.ProjectManager
}

func NewAuthController(pm *utils.ProjectManager) *AuthController {
	return &AuthController{pm: pm}
}

// InitializeAuthRoutes sets up the authentication routes
func InitializeAuthRoutes(r *gin.Engine, dbClient *mongo.Client, authController *AuthController) {
	authClient = dbClient // Assign the passed client to the local variable
	r.POST("/register", authController.registerUser)
	r.POST("/login", authController.loginUser)
	r.POST("/change-password", authController.changePassword) // Add route for changing password
	r.POST("/reset-password", authController.resetPassword)   // Add route for resetting password
}

// registerUser handles user registration
func (ac *AuthController) registerUser(c *gin.Context) {
	logrus.Info("Attempting to register a new user")

	collection := authClient.Database("users").Collection("authentication") // Use the local client variable

	var user models.User
	ac.pm.Execute(func() (interface{}, error) {
		if err := c.ShouldBindJSON(&user); err != nil {
			return nil, err
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedPassword)
		user.ID = primitive.NewObjectID()
		user.IsAdmin = false // Default to false
		return collection.InsertOne(context.TODO(), user)
	}, "Failed to register user")

	logrus.Infof("User registered successfully: %s", user.Username)
	c.JSON(http.StatusCreated, user)
}

// loginUser handles user login
func (ac *AuthController) loginUser(c *gin.Context) {
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
	ac.pm.Execute(func() (interface{}, error) {
		return collection.UpdateOne(context.TODO(), bson.M{"_id": user.ID}, bson.D{{"$set", bson.M{"token": token, "lastLogin": user.LastLogin}}})
	}, "Failed to update login information")

	logrus.Infof("User logged in successfully: %s", user.Username)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// resetPassword handles the generation of a passkey for password reset
func (ac *AuthController) resetPassword(c *gin.Context) {
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
	ac.pm.Execute(func() (interface{}, error) {
		err := collection.FindOne(context.TODO(), bson.M{"email": request.Email}).Decode(&user)
		return user, err
	}, "Failed to find user by email")
	if user == (models.User{}) {
		logrus.Warnf("Email not found: %s", request.Email)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email not found"})
		return
	}

	passkey, err := generatePasskey()
	if err != nil {
		logrus.Errorf("Failed to generate passkey: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate passkey", "details": err.Error()})
		return
	}

	passkeyGeneratedAt := time.Now()
	ac.pm.Execute(func() (interface{}, error) {
		return collection.UpdateOne(context.TODO(), bson.M{"_id": user.ID}, bson.D{
			{"$set", bson.M{"passkey": passkey, "passkeyGeneratedAt": passkeyGeneratedAt}},
		})
	}, "Failed to save passkey")

	// Start a goroutine to invalidate the passkey after 10 minutes
	go func() {
		time.Sleep(10 * time.Minute)
		ac.pm.Execute(func() (interface{}, error) {
			return collection.UpdateOne(context.TODO(), bson.M{"_id": user.ID}, bson.D{
				{"$set", bson.M{"passkey": "NOT_PASSKEY"}},
			})
		}, "Failed to invalidate passkey")
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
func (ac *AuthController) changePassword(c *gin.Context) {
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
	ac.pm.Execute(func() (interface{}, error) {
		err := collection.FindOne(context.TODO(), bson.M{"email": request.Email}).Decode(&user)
		return user, err
	}, "Failed to find user by email")
	if user == (models.User{}) {
		logrus.Warnf("Email not found: %s", request.Email)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email not found"})
		return
	}

	// Check if the passkey is still valid
	if !ac.isPasskeyValid(user) {
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
	ac.pm.Execute(func() (interface{}, error) {
		return collection.UpdateOne(context.TODO(), bson.M{"_id": user.ID}, bson.D{
			{"$set", bson.M{"password": hashedPassword}},
			{"$unset", bson.M{"passkey": "", "passkeyGeneratedAt": ""}},
		})
	}, "Failed to change password")

	logrus.Infof("Password changed successfully for user: %s", request.Email)
	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

func (ac *AuthController) isPasskeyValid(user models.User) bool {
	// Check if the passkey is older than 10 minutes
	if time.Since(user.PasskeyGeneratedAt) > 10*time.Minute {
		// Invalidate the passkey if more than 10 minutes have passed
		collection := authClient.Database("users").Collection("authentication")
		ac.pm.Execute(func() (interface{}, error) {
			return collection.UpdateOne(context.TODO(), bson.M{"_id": user.ID}, bson.D{
				{"$set", bson.M{"passkey": "NOT_PASSKEY"}},
			})
		}, "Failed to invalidate passkey")
		return false
	}
	return true
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		logrus.Errorf("Failed to generate token: %v", err)
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

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

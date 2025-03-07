package tests

import (
	"encoding/json"
	"os"
	"testing"

	"omhs-backend/controllers"
	"omhs-backend/models"
	"omhs-backend/utils"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	authTestManager.RegisterTest(t, "TestRegister")
	// Initialize router and controllers
	router := gin.Default()
	pm := utils.NewProjectManager()
	authController := controllers.NewAuthController(pm)
	requestController := controllers.NewRequestController(pm)
	controllers.InitializeAuthRoutes(router, client, authController)
	controllers.InitializeRequestRoutes(router, client, requestController)

	// Test data
	username := os.Getenv("NON_ADMIN_USER") + generateRandomString(5)
	user := map[string]string{
		"username": username,
		"password": os.Getenv("NON_ADMIN_PASS"),
		"email":    os.Getenv("EMAIL_USER"),
	}

	// Register user
	body, code := RegisterUser(router, user)
	assert.Equal(t, http.StatusCreated, code)

	var registeredUser models.User
	json.Unmarshal([]byte(body), &registeredUser)
	assert.Equal(t, user["username"], registeredUser.Username)
	assert.Equal(t, user["email"], registeredUser.Email)

	// Login user
	body, code = LoginUser(router, user["username"], user["password"])
	assert.Equal(t, http.StatusOK, code)

	var loginResponse map[string]string
	json.Unmarshal([]byte(body), &loginResponse)
	_, ok := loginResponse["token"]
	assert.True(t, ok, "Login response does not contain a token")

	// Admin login
	body, code = LoginUser(router, os.Getenv("ADMIN_USER"), os.Getenv("ADMIN_PASS"))
	assert.Equal(t, http.StatusOK, code)

	json.Unmarshal([]byte(body), &loginResponse)
	adminToken, ok := loginResponse["token"]
	assert.True(t, ok, "Admin login response does not contain a token")

	// Delete registered user
	_, code = DeleteUser(router, registeredUser.ID.Hex(), adminToken)
	assert.Equal(t, http.StatusOK, code)
}

func TestResetPassword(t *testing.T) {
	authTestManager.RegisterTest(t, "TestResetPassword")

	// Initialize router and controllers
	router := gin.Default()
	pm := utils.NewProjectManager()
	authController := controllers.NewAuthController(pm)
	requestController := controllers.NewRequestController(pm)
	controllers.InitializeAuthRoutes(router, client, authController)
	controllers.InitializeRequestRoutes(router, client, requestController)

	// Test data
	username := os.Getenv("NON_ADMIN_USER") + generateRandomString(5)
	user := map[string]string{
		"username": username,
		"password": os.Getenv("NON_ADMIN_PASS"),
		"email":    os.Getenv("EMAIL_USER"),
	}

	// Register user
	body, code := RegisterUser(router, user)
	assert.Equal(t, http.StatusCreated, code)

	var registeredUser models.User
	json.Unmarshal([]byte(body), &registeredUser)
	assert.Equal(t, user["username"], registeredUser.Username)
	assert.Equal(t, user["email"], registeredUser.Email)

	// Reset password
	_, code = ResetPassword(router, user["email"], user["username"])
	assert.Equal(t, http.StatusOK, code)

	// Admin login to retrieve passkey
	body, code = LoginUser(router, os.Getenv("ADMIN_USER"), os.Getenv("ADMIN_PASS"))
	assert.Equal(t, http.StatusOK, code)

	var loginResponse map[string]string
	json.Unmarshal([]byte(body), &loginResponse)
	adminToken, ok := loginResponse["token"]
	assert.True(t, ok, "Admin login response does not contain a token")

	// Retrieve the passkey from the user document
	body, code = GetPasskey(router, "users", "authentication", registeredUser.ID.Hex(), adminToken)
	assert.Equal(t, http.StatusOK, code)

	var userDoc map[string]interface{}
	json.Unmarshal([]byte(body), &userDoc)
	passkey, ok := userDoc["passkey"].(string)
	assert.True(t, ok, "Passkey should not be empty")

	// Change password using the retrieved passkey
	_, code = ChangePassword(router, user["email"], user["username"], passkey, "newPassword")
	assert.Equal(t, http.StatusOK, code)

	// Delete registered user
	_, code = DeleteUser(router, registeredUser.ID.Hex(), adminToken)
	assert.Equal(t, http.StatusOK, code)
}

func TestDuplicateUserRegistration(t *testing.T) {
	authTestManager.RegisterTest(t, "TestDuplicateUserRegistration")

	// Initialize router and controllers
	router := gin.Default()
	pm := utils.NewProjectManager()
	authController := controllers.NewAuthController(pm)
	controllers.InitializeAuthRoutes(router, client, authController)

	// Test data
	user := map[string]string{
		"username": os.Getenv("NON_ADMIN_USER"),
		"password": os.Getenv("NON_ADMIN_PASS"),
		"email":    os.Getenv("EMAIL_USER"),
	}

	// Attempt to register the same user again
	_, code := RegisterUser(router, user)
	assert.Equal(t, http.StatusConflict, code)
}

func TestInvalidLogin(t *testing.T) {
	authTestManager.RegisterTest(t, "TestInvalidLogin")

	// Initialize router and controllers
	router := gin.Default()
	pm := utils.NewProjectManager()
	authController := controllers.NewAuthController(pm)
	controllers.InitializeAuthRoutes(router, client, authController)

	// Test data
	user := map[string]string{
		"username": "invalid_user",
		"password": "invalid_pass",
	}

	// Attempt to login with invalid credentials
	_, code := LoginUser(router, user["username"], user["password"])
	assert.Equal(t, http.StatusUnauthorized, code)
}

func TestPasswordResetWithInvalidEmailOrUsername(t *testing.T) {
	authTestManager.RegisterTest(t, "TestPasswordResetWithInvalidEmailOrUsername")

	// Initialize router and controllers
	router := gin.Default()
	pm := utils.NewProjectManager()
	authController := controllers.NewAuthController(pm)
	controllers.InitializeAuthRoutes(router, client, authController)

	// Test data
	request := map[string]string{
		"email":    "invalid_email@example.com",
		"username": "invalid_user",
	}

	// Attempt to reset password with invalid email or username
	_, code := ResetPassword(router, request["email"], request["username"])
	assert.Equal(t, http.StatusBadRequest, code)
}

func TestPasswordChangeWithInvalidPasskey(t *testing.T) {
	authTestManager.RegisterTest(t, "TestPasswordChangeWithInvalidPasskey")

	// Initialize router and controllers
	router := gin.Default()
	pm := utils.NewProjectManager()
	authController := controllers.NewAuthController(pm)
	requestController := controllers.NewRequestController(pm)
	controllers.InitializeAuthRoutes(router, client, authController)
	controllers.InitializeRequestRoutes(router, client, requestController)

	// Test data
	username := os.Getenv("NON_ADMIN_USER") + generateRandomString(5)
	user := map[string]string{
		"username": username,
		"password": os.Getenv("NON_ADMIN_PASS"),
		"email":    os.Getenv("EMAIL_USER"),
	}

	// Register user
	body, code := RegisterUser(router, user)
	assert.Equal(t, http.StatusCreated, code)

	var registeredUser models.User
	json.Unmarshal([]byte(body), &registeredUser)
	assert.Equal(t, user["username"], registeredUser.Username)
	assert.Equal(t, user["email"], registeredUser.Email)

	// Reset password
	_, code = ResetPassword(router, user["email"], user["username"])
	assert.Equal(t, http.StatusOK, code)

	// Attempt to change password with invalid passkey
	_, code = ChangePassword(router, user["email"], user["username"], "invalid_passkey", "newPassword")
	assert.Equal(t, http.StatusUnauthorized, code)

	// Admin login to delete registered user
	body, code = LoginUser(router, os.Getenv("ADMIN_USER"), os.Getenv("ADMIN_PASS"))
	assert.Equal(t, http.StatusOK, code)

	var loginResponse map[string]string
	json.Unmarshal([]byte(body), &loginResponse)
	adminToken, ok := loginResponse["token"]
	assert.True(t, ok, "Admin login response does not contain a token")

	// Delete registered user
	_, code = DeleteUser(router, registeredUser.ID.Hex(), adminToken)
	assert.Equal(t, http.StatusOK, code)
}

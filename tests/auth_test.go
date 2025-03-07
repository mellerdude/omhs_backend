package tests

import (
	"encoding/json"
	"os"
	"testing"

	"omhs-backend/models"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func registerUserAndGetToken(t *testing.T, router *gin.Engine, user map[string]string) (models.User, string) {
	body, code := RegisterUser(router, user)
	assert.Equal(t, http.StatusCreated, code)

	var registeredUser models.User
	json.Unmarshal([]byte(body), &registeredUser)
	assert.Equal(t, user["username"], registeredUser.Username)
	assert.Equal(t, user["email"], registeredUser.Email)

	body, code = LoginUser(router, user["username"], user["password"])
	assert.Equal(t, http.StatusOK, code)

	var loginResponse map[string]string
	json.Unmarshal([]byte(body), &loginResponse)
	token, ok := loginResponse["token"]
	assert.True(t, ok, "Login response does not contain a token")

	return registeredUser, token
}

func TestRegister(t *testing.T) {
	authTestManager.RegisterTest(t, "TestRegister")
	router, _ := initializeRouterAndControllers(client)

	user := setupTestData()
	registeredUser, _ := registerUserAndGetToken(t, router, user)

	adminToken := AdminLogin(router, os.Getenv("ADMIN_USER"), os.Getenv("ADMIN_PASS"))

	_, code := DeleteUser(router, registeredUser.ID.Hex(), adminToken)
	assert.Equal(t, http.StatusOK, code)
}

func TestResetPassword(t *testing.T) {
	authTestManager.RegisterTest(t, "TestResetPassword")
	router, _ := initializeRouterAndControllers(client)

	user := setupTestData()
	registeredUser, _ := registerUserAndGetToken(t, router, user)

	_, code := ResetPassword(router, user["email"], user["username"])
	assert.Equal(t, http.StatusOK, code)

	adminToken := AdminLogin(router, os.Getenv("ADMIN_USER"), os.Getenv("ADMIN_PASS"))

	body, code := GetPasskey(router, "users", "authentication", registeredUser.ID.Hex(), adminToken)
	assert.Equal(t, http.StatusOK, code)

	var userDoc map[string]interface{}
	json.Unmarshal([]byte(body), &userDoc)
	passkey, ok := userDoc["passkey"].(string)
	assert.True(t, ok, "Passkey should not be empty")

	_, code = ChangePassword(router, user["email"], user["username"], passkey, "newPassword")
	assert.Equal(t, http.StatusOK, code)

	_, code = DeleteUser(router, registeredUser.ID.Hex(), adminToken)
	assert.Equal(t, http.StatusOK, code)
}

func TestDuplicateUserRegistration(t *testing.T) {
	authTestManager.RegisterTest(t, "TestDuplicateUserRegistration")
	router, _ := initializeRouterAndControllers(client)

	user := map[string]string{
		"username": os.Getenv("NON_ADMIN_USER"),
		"password": os.Getenv("NON_ADMIN_PASS"),
		"email":    os.Getenv("EMAIL_USER"),
	}

	_, code := RegisterUser(router, user)
	assert.Equal(t, http.StatusConflict, code)
}

func TestInvalidLogin(t *testing.T) {
	authTestManager.RegisterTest(t, "TestInvalidLogin")
	router, _ := initializeRouterAndControllers(client)

	user := map[string]string{
		"username": "invalid_user",
		"password": "invalid_pass",
	}

	_, code := LoginUser(router, user["username"], user["password"])
	assert.Equal(t, http.StatusUnauthorized, code)
}

func TestPasswordResetWithInvalidEmailOrUsername(t *testing.T) {
	authTestManager.RegisterTest(t, "TestPasswordResetWithInvalidEmailOrUsername")
	router, _ := initializeRouterAndControllers(client)

	request := map[string]string{
		"email":    "invalid_email@example.com",
		"username": "invalid_user",
	}

	_, code := ResetPassword(router, request["email"], request["username"])
	assert.Equal(t, http.StatusBadRequest, code)
}

func TestPasswordChangeWithInvalidPasskey(t *testing.T) {
	authTestManager.RegisterTest(t, "TestPasswordChangeWithInvalidPasskey")
	router, _ := initializeRouterAndControllers(client)

	user := setupTestData()
	registeredUser, _ := registerUserAndGetToken(t, router, user)

	_, code := ResetPassword(router, user["email"], user["username"])
	assert.Equal(t, http.StatusOK, code)

	_, code = ChangePassword(router, user["email"], user["username"], "invalid_passkey", "newPassword")
	assert.Equal(t, http.StatusUnauthorized, code)

	adminToken := AdminLogin(router, os.Getenv("ADMIN_USER"), os.Getenv("ADMIN_PASS"))

	_, code = DeleteUser(router, registeredUser.ID.Hex(), adminToken)
	assert.Equal(t, http.StatusOK, code)
}

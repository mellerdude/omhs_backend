package tests

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"omhs-backend/internal/auth"
)

func registerUserAndGetToken(t *testing.T, router *gin.Engine, user map[string]string) (auth.User, string) {
	body, code := RegisterUser(router, user)
	assert.Equal(t, http.StatusCreated, code)

	var registeredUser auth.User
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
	router, _ := initializeRouterAndControllers(client)

	user := setupTestData()
	registeredUser, _ := registerUserAndGetToken(t, router, user)

	adminToken := AdminLogin(router, os.Getenv("ADMIN_USER"), os.Getenv("ADMIN_PASS"))

	_, code := DeleteUser(router, registeredUser.ID.Hex(), adminToken)
	assert.Equal(t, http.StatusOK, code)

	authTestManager.RegisterTest(t, "TestRegister")
}

func TestResetPassword(t *testing.T) {
	router, _ := initializeRouterAndControllers(client)

	user := setupTestData()
	registeredUser, _ := registerUserAndGetToken(t, router, user)

	_, code := ResetPassword(router, user["email"], user["username"])
	assert.Equal(t, http.StatusOK, code)

	adminToken := AdminLogin(router, os.Getenv("ADMIN_USER"), os.Getenv("ADMIN_PASS"))

	body, code := GetPasskey(router, "users", "authentication", registeredUser.ID.Hex(), adminToken)
	assert.Equal(t, http.StatusOK, code)

	var root map[string]json.RawMessage
	json.Unmarshal([]byte(body), &root)

	// Unmarshal only the "data" field directly into userDoc
	var userDoc map[string]interface{}
	if data, ok := root["data"]; ok {
		json.Unmarshal(data, &userDoc)
	}
	passkey, ok := userDoc["passkey"].(string)
	assert.True(t, ok, "Passkey should not be empty")

	_, code = ChangePassword(router, user["email"], user["username"], passkey, "newPassword")
	assert.Equal(t, http.StatusOK, code)

	_, code = DeleteUser(router, registeredUser.ID.Hex(), adminToken)
	assert.Equal(t, http.StatusOK, code)

	authTestManager.RegisterTest(t, "TestResetPassword")
}

func TestDuplicateUserRegistration(t *testing.T) {
	router, _ := initializeRouterAndControllers(client)

	user := map[string]string{
		"username": os.Getenv("NON_ADMIN_USER"),
		"password": os.Getenv("NON_ADMIN_PASS"),
		"email":    os.Getenv("EMAIL_USER"),
	}

	_, code := RegisterUser(router, user)
	assert.Equal(t, http.StatusConflict, code)

	authTestManager.RegisterTest(t, "TestDuplicateUserRegistration")
}

func TestInvalidLogin(t *testing.T) {
	router, _ := initializeRouterAndControllers(client)

	user := map[string]string{
		"username": "invalid_user",
		"password": "invalid_pass",
	}

	_, code := LoginUser(router, user["username"], user["password"])
	assert.Equal(t, http.StatusUnauthorized, code)

	authTestManager.RegisterTest(t, "TestInvalidLogin")
}

func TestPasswordResetWithInvalidEmailOrUsername(t *testing.T) {
	router, _ := initializeRouterAndControllers(client)

	request := map[string]string{
		"email":    "invalid_email@example.com",
		"username": "invalid_user",
	}

	_, code := ResetPassword(router, request["email"], request["username"])
	assert.Equal(t, http.StatusBadRequest, code)

	authTestManager.RegisterTest(t, "TestPasswordResetWithInvalidEmailOrUsername")
}

func TestPasswordChangeWithInvalidPasskey(t *testing.T) {
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

	authTestManager.RegisterTest(t, "TestPasswordChangeWithInvalidPasskey")
}

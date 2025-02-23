package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"omhs-backend/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func RegisterUser(router *gin.Engine, user map[string]string, t *testing.T) models.User {
	userJSON, _ := json.Marshal(user)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	logrus.Infof("Register Response: %s", w.Body.String())
	assert.Equal(t, http.StatusCreated, w.Code)

	var registeredUser models.User
	json.Unmarshal(w.Body.Bytes(), &registeredUser)
	assert.Equal(t, user["username"], registeredUser.Username)
	assert.Equal(t, user["email"], registeredUser.Email)

	return registeredUser
}

func LoginUser(router *gin.Engine, username, password string, t *testing.T) string {
	loginReqJSON, _ := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	loginReq, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginReqJSON))
	loginReq.Header.Set("Content-Type", "application/json")

	loginRecorder := httptest.NewRecorder()
	router.ServeHTTP(loginRecorder, loginReq)

	logrus.Infof("Login Response: %s", loginRecorder.Body.String())
	assert.Equal(t, http.StatusOK, loginRecorder.Code)

	var loginResponse map[string]string
	json.Unmarshal(loginRecorder.Body.Bytes(), &loginResponse)
	token, ok := loginResponse["token"]
	assert.True(t, ok, "Login response does not contain a token")

	return token
}

func DeleteUser(router *gin.Engine, userID, adminToken string, t *testing.T) {
	deleteReq, _ := http.NewRequest("DELETE", "/users/authentication/"+userID, nil)
	deleteReq.Header.Set("Authorization", "Bearer "+adminToken)
	deleteRecorder := httptest.NewRecorder()
	router.ServeHTTP(deleteRecorder, deleteReq)

	logrus.Infof("Delete Response: %s", deleteRecorder.Body.String())
	assert.Equal(t, http.StatusOK, deleteRecorder.Code)
}

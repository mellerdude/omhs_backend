package tests

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"omhs-backend/controllers"
	"omhs-backend/models"
	"omhs-backend/utils"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func generateRandomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func RegisterUser(router *gin.Engine, user map[string]string) (string, int) {
	userJSON, _ := json.Marshal(user)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	logrus.Infof("Register Response: %s", w.Body.String())
	return w.Body.String(), w.Code
}

func LoginUser(router *gin.Engine, username, password string) (string, int) {
	loginReqJSON, _ := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	loginReq, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginReqJSON))
	loginReq.Header.Set("Content-Type", "application/json")

	loginRecorder := httptest.NewRecorder()
	router.ServeHTTP(loginRecorder, loginReq)

	logrus.Infof("Login Response: %s", loginRecorder.Body.String())
	return loginRecorder.Body.String(), loginRecorder.Code
}

func DeleteUser(router *gin.Engine, userID, adminToken string) (string, int) {
	deleteReq, _ := http.NewRequest("DELETE", "/users/authentication/"+userID, nil)
	deleteReq.Header.Set("Authorization", "Bearer "+adminToken)
	deleteRecorder := httptest.NewRecorder()
	router.ServeHTTP(deleteRecorder, deleteReq)

	logrus.Infof("Delete Response: %s", deleteRecorder.Body.String())
	return deleteRecorder.Body.String(), deleteRecorder.Code
}

func ResetPassword(router *gin.Engine, email, username string) (string, int) {
	resetReqJSON, _ := json.Marshal(map[string]string{
		"email":    email,
		"username": username,
	})
	resetReq, _ := http.NewRequest("POST", "/reset-password", bytes.NewBuffer(resetReqJSON))
	resetReq.Header.Set("Content-Type", "application/json")

	resetRecorder := httptest.NewRecorder()
	router.ServeHTTP(resetRecorder, resetReq)

	logrus.Infof("Reset Password Response: %s", resetRecorder.Body.String())
	return resetRecorder.Body.String(), resetRecorder.Code
}

func ChangePassword(router *gin.Engine, email, username, passkey, newPassword string) (string, int) {
	changeReqJSON, _ := json.Marshal(map[string]string{
		"email":       email,
		"username":    username,
		"passkey":     passkey,
		"newPassword": newPassword,
	})
	changeReq, _ := http.NewRequest("POST", "/change-password", bytes.NewBuffer(changeReqJSON))
	changeReq.Header.Set("Content-Type", "application/json")

	changeRecorder := httptest.NewRecorder()
	router.ServeHTTP(changeRecorder, changeReq)

	logrus.Infof("Change Password Response: %s", changeRecorder.Body.String())
	return changeRecorder.Body.String(), changeRecorder.Code
}

func GetPasskey(router *gin.Engine, database, collection, userID, adminToken string) (string, int) {
	getReq, _ := http.NewRequest("GET", "/"+database+"/"+collection+"/"+userID, nil)
	getReq.Header.Set("Authorization", "Bearer "+adminToken)
	getRecorder := httptest.NewRecorder()
	router.ServeHTTP(getRecorder, getReq)

	logrus.Infof("Get Document Response: %s", getRecorder.Body.String())
	return getRecorder.Body.String(), getRecorder.Code
}

func AdminLogin(router *gin.Engine, username, password string) string {
	body, _ := LoginUser(router, username, password)
	var loginResponse map[string]string
	json.Unmarshal([]byte(body), &loginResponse)
	return loginResponse["token"]
}

func setupTestData() map[string]string {
	username := os.Getenv("NON_ADMIN_USER") + generateRandomString(5)
	return map[string]string{
		"username": username,
		"password": os.Getenv("NON_ADMIN_PASS"),
		"email":    os.Getenv("EMAIL_USER"),
	}
}

func initializeRouterAndControllers(client *mongo.Client) (*gin.Engine, *utils.ProjectManager) {
	router := gin.Default()
	pm := utils.NewProjectManager()
	authController := controllers.NewAuthController(pm)
	requestController := controllers.NewRequestController(pm)
	controllers.InitializeAuthRoutes(router, client, authController)
	controllers.InitializeRequestRoutes(router, client, requestController)
	return router, pm
}

func createDocument(router *gin.Engine, database, collection, adminToken string, doc models.Document) (string, int) {
	docJSON, _ := json.Marshal(doc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/"+database+"/"+collection, bytes.NewBuffer(docJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)
	router.ServeHTTP(w, req)

	logrus.Infof("Create Document Response: %s", w.Body.String())
	return w.Body.String(), w.Code
}

func getDocument(router *gin.Engine, database, collection, docID, adminToken string) (string, int) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/"+database+"/"+collection+"/"+docID, nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	router.ServeHTTP(w, req)

	logrus.Infof("Get Document Response: %s", w.Body.String())
	return w.Body.String(), w.Code
}

func updateDocument(router *gin.Engine, database, collection, docID, adminToken string, doc models.Document) (string, int) {
	docJSON, _ := json.Marshal(doc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/"+database+"/"+collection+"/"+docID, bytes.NewBuffer(docJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)
	router.ServeHTTP(w, req)

	logrus.Infof("Update Document Response: %s", w.Body.String())
	return w.Body.String(), w.Code
}

func deleteDocument(router *gin.Engine, database, collection, docID, adminToken string) (string, int) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/"+database+"/"+collection+"/"+docID, nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	router.ServeHTTP(w, req)

	logrus.Infof("Delete Document Response: %s", w.Body.String())
	return w.Body.String(), w.Code
}

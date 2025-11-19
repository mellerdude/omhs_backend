package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"

	"omhs-backend/internal/auth"
	"omhs-backend/internal/kanban"
	"omhs-backend/internal/middleware"
	"omhs-backend/internal/requests"
)

// Build full router with Auth + Requests + Kanban + JWT middleware
func setupKanbanRouter(client *mongo.Client) *gin.Engine {
	router := gin.Default()
	api := router.Group("/api")

	// Auth
	authRepo := auth.NewMongoUserRepository(client)
	authService := auth.NewAuthService(authRepo)
	authController := auth.NewAuthController(authService)
	auth.RegisterRoutes(api, authController)

	// Requests
	requestRepo := requests.NewMongoRequestRepository(client)
	requestService := requests.NewRequestService(requestRepo)
	requestController := requests.NewRequestController(requestService)
	requests.RegisterRoutes(api, requestController)

	// Kanban
	kanbanRepo := kanban.NewKanbanRepository(requestRepo)
	kanbanService := kanban.NewKanbanService(*kanbanRepo)
	kanbanController := kanban.NewKanbanController(kanbanService)

	protected := api.Group("/")
	protected.Use(middleware.JWTMiddleware())
	kanban.RegisterRoutes(protected, kanbanController)

	return router
}

func TestKanbanUnauthorizedAccess(t *testing.T) {
	router := setupKanbanRouter(client)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/kanban", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestKanbanCreate(t *testing.T) {
	router := setupKanbanRouter(client)

	user := setupTestData()
	registeredUser, token := registerUserAndGetToken(t, router, user)

	body := map[string]interface{}{
		"title": "My First Board",
		"lists": []interface{}{},
	}

	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/kanban", bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// --- CLEANUP ---
	adminToken := AdminLogin(router, os.Getenv("ADMIN_USER"), os.Getenv("ADMIN_PASS"))

	_, code := deleteDocument(router, "data", "Kanbans", registeredUser.ID.Hex(), adminToken)
	assert.Equal(t, http.StatusOK, code)

	_, code = DeleteUser(router, registeredUser.ID.Hex(), adminToken)
	assert.Equal(t, http.StatusOK, code)
}

func TestKanbanGet(t *testing.T) {
	router := setupKanbanRouter(client)

	user := setupTestData()
	registeredUser, token := registerUserAndGetToken(t, router, user)

	// Create board
	createBody := map[string]interface{}{
		"title": "Board",
		"lists": []interface{}{},
	}

	jsonBody, _ := json.Marshal(createBody)
	req1, _ := http.NewRequest("POST", "/api/kanban", bytes.NewBuffer(jsonBody))
	req1.Header.Set("Authorization", "Bearer "+token)
	req1.Header.Set("Content-Type", "application/json")

	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	// GET board
	req2, _ := http.NewRequest("GET", "/api/kanban", nil)
	req2.Header.Set("Authorization", "Bearer "+token)

	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	var data map[string]interface{}
	json.Unmarshal(w2.Body.Bytes(), &data)
	assert.Equal(t, "Board", data["title"])

	// --- CLEANUP ---
	adminToken := AdminLogin(router, os.Getenv("ADMIN_USER"), os.Getenv("ADMIN_PASS"))

	_, code := deleteDocument(router, "data", "Kanbans", registeredUser.ID.Hex(), adminToken)
	assert.Equal(t, http.StatusOK, code)

	_, code = DeleteUser(router, registeredUser.ID.Hex(), adminToken)
	assert.Equal(t, http.StatusOK, code)
}

func TestKanbanUpdate(t *testing.T) {
	router := setupKanbanRouter(client)

	user := setupTestData()
	registeredUser, token := registerUserAndGetToken(t, router, user)

	// Create board
	createBody := map[string]interface{}{
		"title": "Initial",
		"lists": []interface{}{},
	}

	jsonBody, _ := json.Marshal(createBody)
	req, _ := http.NewRequest("POST", "/api/kanban", bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Update board
	updateBody := map[string]interface{}{
		"title": "Updated",
		"lists": []interface{}{
			map[string]interface{}{
				"name":  "Todo",
				"tasks": []interface{}{},
			},
		},
	}

	jsonUpdate, _ := json.Marshal(updateBody)
	req2, _ := http.NewRequest("PUT", "/api/kanban", bytes.NewBuffer(jsonUpdate))
	req2.Header.Set("Authorization", "Bearer "+token)
	req2.Header.Set("Content-Type", "application/json")

	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	// --- CLEANUP ---
	adminToken := AdminLogin(router, os.Getenv("ADMIN_USER"), os.Getenv("ADMIN_PASS"))

	_, code := deleteDocument(router, "data", "Kanbans", registeredUser.ID.Hex(), adminToken)
	assert.Equal(t, http.StatusOK, code)

	_, code = DeleteUser(router, registeredUser.ID.Hex(), adminToken)
	assert.Equal(t, http.StatusOK, code)
}

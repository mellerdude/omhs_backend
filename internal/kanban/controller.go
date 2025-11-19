package kanban

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type KanbanController struct {
	service *KanbanService
}

func NewKanbanController(s *KanbanService) *KanbanController {
	return &KanbanController{service: s}
}

func getUserId(c *gin.Context) (primitive.ObjectID, bool) {
	val, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no userId in token"})
		return primitive.NilObjectID, false
	}
	return val.(primitive.ObjectID), true
}

func (ctr *KanbanController) GetKanban(c *gin.Context) {
	userId, ok := getUserId(c)
	if !ok {
		return
	}

	data, err := ctr.service.GetKanban(userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (ctr *KanbanController) CreateKanban(c *gin.Context) {
	userId, ok := getUserId(c)
	if !ok {
		return
	}

	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	doc, err := ctr.service.CreateKanban(userId, body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, doc)
}

func (ctr *KanbanController) UpdateKanban(c *gin.Context) {
	userId, ok := getUserId(c)
	if !ok {
		return
	}

	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	if err := ctr.service.UpdateKanban(userId, body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

package requests

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RequestController struct {
	service *RequestService
}

func NewRequestController(s *RequestService) *RequestController {
	return &RequestController{service: s}
}

func (ctr *RequestController) Create(c *gin.Context) {
	db := c.Param("database")
	col := c.Param("collection")

	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	doc, err := ctr.service.Create(db, col, body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, doc)
}

func (ctr *RequestController) Get(c *gin.Context) {
	db := c.Param("database")
	col := c.Param("collection")
	id := c.Param("id")

	doc, err := ctr.service.Get(db, col, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, doc)
}

func (ctr *RequestController) Update(c *gin.Context) {
	db := c.Param("database")
	col := c.Param("collection")
	id := c.Param("id")

	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	doc, err := ctr.service.Update(db, col, id, body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, doc)
}

func (ctr *RequestController) Delete(c *gin.Context) {
	db := c.Param("database")
	col := c.Param("collection")
	id := c.Param("id")

	if err := ctr.service.Delete(db, col, id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (ctr *RequestController) GetAll(c *gin.Context) {
	db := c.Param("database")
	col := c.Param("collection")

	docs, err := ctr.service.GetAll(db, col)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, docs)
}

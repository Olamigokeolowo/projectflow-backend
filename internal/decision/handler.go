package decision

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) List(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"decisions": []string{}})
}

func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"id": id, "title": "stub decision"})
}

func (h *Handler) Create(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "created"})
}

func (h *Handler) ListTasks(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"decision_id": id, "tasks": []string{}})
}
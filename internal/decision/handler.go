package decision

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(c *gin.Context) {
	decisions, err := h.service.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list decisions"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"decisions": decisions})
}

func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")

	d, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "decision not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get decision"})
		return
	}
	c.JSON(http.StatusOK, d)
}

func (h *Handler) Create(c *gin.Context) {
	var req struct {
		Title  string `json:"title"`
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	d, err := h.service.Create(c.Request.Context(), req.Title, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create decision"})
		return
	}
	c.JSON(http.StatusCreated, d)
}

func (h *Handler) ListTasks(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"decision_id": id, "tasks": []string{}})
}
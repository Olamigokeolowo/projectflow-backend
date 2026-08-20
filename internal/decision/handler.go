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
	userID := c.GetString("user_id") // set by AuthRequired middleware

	d, err := h.service.Get(c.Request.Context(), id, userID)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "decision not found"})
		case errors.Is(err, ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": "you do not have access to this decision"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get decision"})
		}
		return
	}
	c.JSON(http.StatusOK, d)
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateDecisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("user_id") // set by AuthRequired middleware

	d, err := h.service.Create(c.Request.Context(), req.Title, req.Status, userID)
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

func (h *Handler) SlowOperation(c *gin.Context) {
	err := h.service.SlowOperation(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusRequestTimeout, gin.H{"error": "request cancelled or timed out"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "slow operation completed"})
}
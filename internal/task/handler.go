package task

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
	decisionID := c.Param("id")
	userID := c.GetString("user_id")

	tasks, err := h.service.ListByDecision(c.Request.Context(), decisionID, userID)
	if err != nil {
		if errors.Is(err, ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "you do not have access to this workspace"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list tasks"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func (h *Handler) Create(c *gin.Context) {
	decisionID := c.Param("id")
	userID := c.GetString("user_id")

	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	t, err := h.service.Create(c.Request.Context(), req.Title, decisionID, req.AssigneeID, userID)
	if err != nil {
		if errors.Is(err, ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "you do not have access to this workspace"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create task"})
		return
	}
	c.JSON(http.StatusCreated, t)
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")

	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	t, err := h.service.Update(c.Request.Context(), id, userID, req.Status, req.AssigneeID)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		case errors.Is(err, ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": "you do not have access to this workspace"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update task"})
		}
		return
	}
	c.JSON(http.StatusOK, t)
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")

	err := h.service.Delete(c.Request.Context(), id, userID)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		case errors.Is(err, ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": "you do not have access to this workspace"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete task"})
		}
		return
	}
	c.Status(http.StatusNoContent)
}
package comment

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

func (h *Handler) ListForDecision(c *gin.Context) {
	decisionID := c.Param("id")
	userID := c.GetString("user_id")

	comments, err := h.service.ListByTarget(c.Request.Context(), "decision", decisionID, userID)
	if err != nil {
		if errors.Is(err, ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "you do not have access to this workspace"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list comments"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"comments": comments})
}

func (h *Handler) CreateForDecision(c *gin.Context) {
	decisionID := c.Param("id")
	userID := c.GetString("user_id")

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cm, err := h.service.Create(c.Request.Context(), req.Body, userID, "decision", decisionID)
	if err != nil {
		if errors.Is(err, ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "you do not have access to this workspace"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create comment"})
		return
	}
	c.JSON(http.StatusCreated, cm)
}

func (h *Handler) ListForTask(c *gin.Context) {
	taskID := c.Param("id")
	userID := c.GetString("user_id")

	comments, err := h.service.ListByTarget(c.Request.Context(), "task", taskID, userID)
	if err != nil {
		if errors.Is(err, ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "you do not have access to this workspace"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list comments"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"comments": comments})
}

func (h *Handler) CreateForTask(c *gin.Context) {
	taskID := c.Param("id")
	userID := c.GetString("user_id")

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cm, err := h.service.Create(c.Request.Context(), req.Body, userID, "task", taskID)
	if err != nil {
		if errors.Is(err, ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "you do not have access to this workspace"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create comment"})
		return
	}
	c.JSON(http.StatusCreated, cm)
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")

	err := h.service.Delete(c.Request.Context(), id, userID)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
		case errors.Is(err, ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": "you do not have permission to delete this comment"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete comment"})
		}
		return
	}
	c.Status(http.StatusNoContent)
}
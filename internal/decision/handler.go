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
	workspaceID := c.Query("workspace_id")
	userID := c.GetString("user_id")

	decisions, err := h.service.ListByWorkspace(c.Request.Context(), workspaceID, userID)
	if err != nil {
		if errors.Is(err, ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "you do not have access to this workspace"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list decisions"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"decisions": decisions})
}

func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")

	d, err := h.service.Get(c.Request.Context(), id, userID)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "decision not found"})
		case errors.Is(err, ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": "you do not have access to this workspace"})
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

	userID := c.GetString("user_id")

	d, err := h.service.Create(c.Request.Context(), req.Title, req.Status, userID, req.WorkspaceID)
	if err != nil {
		if errors.Is(err, ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "you do not have access to this workspace"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create decision"})
		return
	}
	c.JSON(http.StatusCreated, d)
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")

	var req UpdateDecisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	d, err := h.service.Update(c.Request.Context(), id, userID, req.Title, req.Status)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "decision not found"})
		case errors.Is(err, ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": "you do not have access to this workspace"})
		case errors.Is(err, ErrNotOwner):
			c.JSON(http.StatusForbidden, gin.H{"error": "you do not have permission to update this decision"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update decision"})
		}
		return
	}
	c.JSON(http.StatusOK, d)
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")

	err := h.service.Delete(c.Request.Context(), id, userID)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "decision not found"})
		case errors.Is(err, ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": "you do not have access to this workspace"})
		case errors.Is(err, ErrNotOwner):
			c.JSON(http.StatusForbidden, gin.H{"error": "you do not have permission to delete this decision"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete decision"})
		}
		return
	}
	c.Status(http.StatusNoContent)
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
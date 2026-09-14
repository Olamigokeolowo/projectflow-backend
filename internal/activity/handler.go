package activity

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
	workspaceID := c.Param("id")
	userID := c.GetString("user_id")

	entries, err := h.service.ListByWorkspace(c.Request.Context(), workspaceID, userID)
	if err != nil {
		if errors.Is(err, ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "you do not have access to this workspace"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list activity"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"activity": entries})
}
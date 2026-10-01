package workspace

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

func (h *Handler) Create(c *gin.Context) {
	var req CreateWorkspaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetString("user_id")
	w, err := h.service.Create(c.Request.Context(), req.Name, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create workspace"})
		return
	}
	c.JSON(http.StatusCreated, w)
}

func (h *Handler) List(c *gin.Context) {
	userID := c.GetString("user_id")
	workspaces, err := h.service.ListForUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list workspaces"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"workspaces": workspaces})
}

func (h *Handler) AddMember(c *gin.Context) {
	workspaceID := c.Param("id")
	userID := c.GetString("user_id")

	var req AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.AddMember(c.Request.Context(), workspaceID, userID, req.Email)
	if err != nil {
		switch {
		case errors.Is(err, ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": "you are not a member of this workspace"})
		case errors.Is(err, ErrNotAdmin):
			c.JSON(http.StatusForbidden, gin.H{"error": "only workspace admins can add members"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add member"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "member added"})
}

func (h *Handler) RemoveMember(c *gin.Context) {
	workspaceID := c.Param("id")
	targetUserID := c.Param("userId")
	requestingUserID := c.GetString("user_id")

	err := h.service.RemoveMember(c.Request.Context(), workspaceID, requestingUserID, targetUserID)
	if err != nil {
		switch {
		case errors.Is(err, ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": "you are not a member of this workspace"})
		case errors.Is(err, ErrNotAdmin):
			c.JSON(http.StatusForbidden, gin.H{"error": "only workspace admins can remove members"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove member"})
		}
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) ListMembers(c *gin.Context) {
	workspaceID := c.Param("id")
	userID := c.GetString("user_id")

	members, err := h.service.ListMembers(c.Request.Context(), workspaceID, userID)
	if err != nil {
		if errors.Is(err, ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "you are not a member of this workspace"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list members"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"members": members})
}
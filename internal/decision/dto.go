package decision

type CreateDecisionRequest struct {
	Title       string `json:"title" binding:"required,min=3,max=200"`
	Status      string `json:"status" binding:"required,oneof=draft proposed accepted rejected"`
	WorkspaceID string `json:"workspace_id" binding:"required"`
}

type UpdateDecisionRequest struct {
	Title  *string `json:"title,omitempty" binding:"omitempty,min=3,max=200"`
	Status *string `json:"status,omitempty" binding:"omitempty,oneof=draft proposed accepted rejected"`
}
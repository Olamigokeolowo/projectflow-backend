package task

type CreateTaskRequest struct {
	Title      string `json:"title" binding:"required,min=3,max=200"`
	AssigneeID string `json:"assignee_id" binding:"required"`
}

type UpdateTaskRequest struct {
	Status     *string `json:"status,omitempty" binding:"omitempty,oneof=todo in_progress done"`
	AssigneeID *string `json:"assignee_id,omitempty"`
}
package workspace

import "time"

type Workspace struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

type Membership struct {
	UserID      string    `json:"user_id"`
	WorkspaceID string    `json:"workspace_id"`
	JoinedAt    time.Time `json:"joined_at"`
}
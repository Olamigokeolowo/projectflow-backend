package activity

import "time"

type Activity struct {
	ID          string    `json:"id"`
	WorkspaceID string    `json:"workspace_id"`
	ActorID     string    `json:"actor_id"`
	Verb        string    `json:"verb"`
	TargetType  string    `json:"target_type"`
	TargetID    string    `json:"target_id"`
	Summary     string    `json:"summary"`
	CreatedAt   time.Time `json:"created_at"`
}
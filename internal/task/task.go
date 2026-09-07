package task

import "time"

type Task struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	Status     string    `json:"status"`
	DecisionID string    `json:"decision_id"`
	AssigneeID string    `json:"assignee_id"`
	CreatedBy  string    `json:"created_by"`
	CreatedAt  time.Time `json:"created_at"`
}
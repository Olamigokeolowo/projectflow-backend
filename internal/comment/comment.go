package comment

import "time"

type Comment struct {
	ID         string    `json:"id"`
	Body       string    `json:"body"`
	AuthorID   string    `json:"author_id"`
	TargetType string    `json:"target_type"`
	TargetID   string    `json:"target_id"`
	CreatedAt  time.Time `json:"created_at"`
}
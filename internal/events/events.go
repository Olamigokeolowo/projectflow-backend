package events

import "context"

type DecisionCreated struct {
	DecisionID string
	OwnerID    string
	Title      string
}

type Publisher interface {
	Publish(ctx context.Context, event DecisionCreated) error
}

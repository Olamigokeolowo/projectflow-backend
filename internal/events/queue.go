package events

import (
	"context"
	"log"
)

// InMemoryQueue is a temporary stand-in for a real message broker.
type InMemoryQueue struct {
	events chan DecisionCreated
}

func NewInMemoryQueue(bufferSize int) *InMemoryQueue {
	return &InMemoryQueue{
		events: make(chan DecisionCreated, bufferSize),
	}
}

func (q *InMemoryQueue) Publish(ctx context.Context, event DecisionCreated) error {
	select {
	case q.events <- event:
		return nil
	default:
		log.Println("event queue full, dropping event:", event.DecisionID)
		return nil // in production, this might return an error or block instead
	}
}

func (q *InMemoryQueue) Consume() <-chan DecisionCreated {
	return q.events
}

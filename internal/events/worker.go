package events

import (
	"context"
	"log"
)

func StartWorker(ctx context.Context, queue *InMemoryQueue) {
	go func() {
		for {
			select {
			case event := <-queue.Consume():
				processEvent(event)
			case <-ctx.Done():
				log.Println("worker shutting down")
				return
			}
		}
	}()
}

func processEvent(event DecisionCreated) {
	// simulates background work: notifications, audit logs, activity feeds, etc.
	log.Printf("processing event: decision %s created by %s (title: %s)", event.DecisionID, event.OwnerID, event.Title)
}
package activity

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, workspaceID, actorID, verb, targetType, targetID, summary string) (*Activity, error)
	ListByWorkspace(ctx context.Context, workspaceID string) ([]*Activity, error)
}

type InMemoryRepository struct {
	mu   sync.RWMutex
	data map[string]*Activity
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{data: make(map[string]*Activity)}
}

func (r *InMemoryRepository) Create(ctx context.Context, workspaceID, actorID, verb, targetType, targetID, summary string) (*Activity, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	a := &Activity{
		ID:          uuid.NewString(),
		WorkspaceID: workspaceID,
		ActorID:     actorID,
		Verb:        verb,
		TargetType:  targetType,
		TargetID:    targetID,
		Summary:     summary,
		CreatedAt:   time.Now(),
	}
	r.data[a.ID] = a
	return a, nil
}

func (r *InMemoryRepository) ListByWorkspace(ctx context.Context, workspaceID string) ([]*Activity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*Activity
	for _, a := range r.data {
		if a.WorkspaceID == workspaceID {
			result = append(result, a)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.After(result[j].CreatedAt)
	})

	return result, nil
}
package task

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("task not found")

type Repository interface {
	ListByDecision(ctx context.Context, decisionID string) ([]*Task, error)
	GetByID(ctx context.Context, id string) (*Task, error)
	Create(ctx context.Context, title, decisionID, assigneeID, createdBy string) (*Task, error)
	Update(ctx context.Context, t *Task) (*Task, error)
	Delete(ctx context.Context, id string) error
}

// InMemoryRepository is a temporary stand-in for a real database.
type InMemoryRepository struct {
	mu   sync.RWMutex
	data map[string]*Task
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		data: make(map[string]*Task),
	}
}

func (r *InMemoryRepository) Create(ctx context.Context, title, decisionID, assigneeID, createdBy string) (*Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	t := &Task{
		ID:         uuid.NewString(),
		Title:      title,
		Status:     "todo",
		DecisionID: decisionID,
		AssigneeID: assigneeID,
		CreatedBy:  createdBy,
		CreatedAt:  time.Now(),
	}
	r.data[t.ID] = t
	return t, nil
}

func (r *InMemoryRepository) ListByDecision(ctx context.Context, decisionID string) ([]*Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*Task
	for _, t := range r.data {
		if t.DecisionID == decisionID {
			result = append(result, t)
		}
	}
	return result, nil
}

func (r *InMemoryRepository) GetByID(ctx context.Context, id string) (*Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	t, ok := r.data[id]
	if !ok {
		return nil, ErrNotFound
	}
	return t, nil
}

func (r *InMemoryRepository) Update(ctx context.Context, t *Task) (*Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[t.ID]; !ok {
		return nil, ErrNotFound
	}
	r.data[t.ID] = t
	return t, nil
}

func (r *InMemoryRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[id]; !ok {
		return ErrNotFound
	}
	delete(r.data, id)
	return nil
}
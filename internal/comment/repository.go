package comment

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("comment not found")

type Repository interface {
	ListByTarget(ctx context.Context, targetType, targetID string) ([]*Comment, error)
	GetByID(ctx context.Context, id string) (*Comment, error)
	Create(ctx context.Context, body, authorID, targetType, targetID string) (*Comment, error)
	Delete(ctx context.Context, id string) error
}

type InMemoryRepository struct {
	mu   sync.RWMutex
	data map[string]*Comment
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		data: make(map[string]*Comment),
	}
}

func (r *InMemoryRepository) Create(ctx context.Context, body, authorID, targetType, targetID string) (*Comment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	c := &Comment{
		ID:         uuid.NewString(),
		Body:       body,
		AuthorID:   authorID,
		TargetType: targetType,
		TargetID:   targetID,
		CreatedAt:  time.Now(),
	}
	r.data[c.ID] = c
	return c, nil
}

func (r *InMemoryRepository) ListByTarget(ctx context.Context, targetType, targetID string) ([]*Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*Comment
	for _, c := range r.data {
		if c.TargetType == targetType && c.TargetID == targetID {
			result = append(result, c)
		}
	}
	return result, nil
}

func (r *InMemoryRepository) GetByID(ctx context.Context, id string) (*Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	c, ok := r.data[id]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
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
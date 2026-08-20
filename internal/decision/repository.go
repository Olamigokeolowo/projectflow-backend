package decision

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("decision not found")

type Repository interface {
	List(ctx context.Context) ([]*Decision, error)
	GetByID(ctx context.Context, id string) (*Decision, error)
	Create(ctx context.Context, title, status, ownerID string) (*Decision, error)
	SlowOperation(ctx context.Context) error // simulates a slow query, for demonstrating cancellation
}

// InMemoryRepository is a temporary stand-in for a real database.
type InMemoryRepository struct {
	mu   sync.RWMutex
	data map[string]*Decision
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		data: make(map[string]*Decision),
	}
}

func (r *InMemoryRepository) Create(ctx context.Context, title, status, ownerID string) (*Decision, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	d := &Decision{
		ID:        uuid.NewString(),
		Title:     title,
		Status:    status,
		OwnerID:   ownerID,
		CreatedAt: time.Now(),
	}
	r.data[d.ID] = d
	return d, nil
}

func (r *InMemoryRepository) List(ctx context.Context) ([]*Decision, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*Decision, 0, len(r.data))
	for _, d := range r.data {
		result = append(result, d)
	}
	return result, nil
}

func (r *InMemoryRepository) GetByID(ctx context.Context, id string) (*Decision, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	d, ok := r.data[id]
	if !ok {
		return nil, ErrNotFound
	}
	return d, nil
}
func (r *InMemoryRepository) SlowOperation(ctx context.Context) error {
	select {
	case <-time.After(5 * time.Second):
		return nil
	case <-ctx.Done():
		log.Println("slow operation cancelled early:", ctx.Err())
		return ctx.Err()
	}
}
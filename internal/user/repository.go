package user

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("user not found")
var ErrEmailTaken = errors.New("email already registered")

type Repository interface {
	Create(ctx context.Context, email, passwordHash string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
}

type InMemoryRepository struct {
	mu   sync.RWMutex
	data map[string]*User // keyed by email
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{data: make(map[string]*User)}
}

func (r *InMemoryRepository) Create(ctx context.Context, email, passwordHash string) (*User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.data[email]; exists {
		return nil, ErrEmailTaken
	}

	u := &User{
		ID:           uuid.NewString(),
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
	}
	r.data[email] = u
	return u, nil
}

func (r *InMemoryRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.data[email]
	if !ok {
		return nil, ErrNotFound
	}
	return u, nil
}
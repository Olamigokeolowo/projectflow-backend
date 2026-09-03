package workspace

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("workspace not found")

type Repository interface {
	Create(ctx context.Context, name, createdBy string) (*Workspace, error)
	GetByID(ctx context.Context, id string) (*Workspace, error)
	ListForUser(ctx context.Context, userID string) ([]*Workspace, error)
	AddMember(ctx context.Context, workspaceID, userID string) error
	IsMember(ctx context.Context, workspaceID, userID string) (bool, error)
	ListMembers(ctx context.Context, workspaceID string) ([]string, error)
}

type InMemoryRepository struct {
	mu          sync.RWMutex
	workspaces  map[string]*Workspace
	memberships map[string][]string
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		workspaces:  make(map[string]*Workspace),
		memberships: make(map[string][]string),
	}
}

func (r *InMemoryRepository) Create(ctx context.Context, name, createdBy string) (*Workspace, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	w := &Workspace{
		ID:        uuid.NewString(),
		Name:      name,
		CreatedBy: createdBy,
		CreatedAt: time.Now(),
	}
	r.workspaces[w.ID] = w
	r.memberships[w.ID] = []string{createdBy}
	return w, nil
}

func (r *InMemoryRepository) GetByID(ctx context.Context, id string) (*Workspace, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	w, ok := r.workspaces[id]
	if !ok {
		return nil, ErrNotFound
	}
	return w, nil
}

func (r *InMemoryRepository) ListForUser(ctx context.Context, userID string) ([]*Workspace, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*Workspace
	for wsID, members := range r.memberships {
		for _, m := range members {
			if m == userID {
				result = append(result, r.workspaces[wsID])
				break
			}
		}
	}
	return result, nil
}

func (r *InMemoryRepository) AddMember(ctx context.Context, workspaceID, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.workspaces[workspaceID]; !ok {
		return ErrNotFound
	}
	for _, m := range r.memberships[workspaceID] {
		if m == userID {
			return nil
		}
	}
	r.memberships[workspaceID] = append(r.memberships[workspaceID], userID)
	return nil
}

func (r *InMemoryRepository) IsMember(ctx context.Context, workspaceID, userID string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, m := range r.memberships[workspaceID] {
		if m == userID {
			return true, nil
		}
	}
	return false, nil
}

func (r *InMemoryRepository) ListMembers(ctx context.Context, workspaceID string) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.memberships[workspaceID], nil
}
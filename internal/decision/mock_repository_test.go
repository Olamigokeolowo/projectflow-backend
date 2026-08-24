package decision

import (
	"context"

	"github.com/Olamigokeolowo/projectflow-backend/internal/events"
)

// mockRepository is a test double implementing the Repository interface.
type mockRepository struct {
	listFunc          func(ctx context.Context) ([]*Decision, error)
	getByIDFunc       func(ctx context.Context, id string) (*Decision, error)
	createFunc        func(ctx context.Context, title, status, ownerID string) (*Decision, error)
	slowOperationFunc func(ctx context.Context) error
}

func (m *mockRepository) List(ctx context.Context) ([]*Decision, error) {
	return m.listFunc(ctx)
}

func (m *mockRepository) GetByID(ctx context.Context, id string) (*Decision, error) {
	return m.getByIDFunc(ctx, id)
}

func (m *mockRepository) Create(ctx context.Context, title, status, ownerID string) (*Decision, error) {
	return m.createFunc(ctx, title, status, ownerID)
}

func (m *mockRepository) SlowOperation(ctx context.Context) error {
	return m.slowOperationFunc(ctx)
}

// mockPublisher is a test double implementing the events.Publisher interface.
type mockPublisher struct {
	publishFunc func(ctx context.Context, event events.DecisionCreated) error
}

func (m *mockPublisher) Publish(ctx context.Context, event events.DecisionCreated) error {
	if m.publishFunc != nil {
		return m.publishFunc(ctx, event)
	}
	return nil
}

type mockCache struct {
	store map[string]string
}

func newMockCache() *mockCache {
	return &mockCache{store: make(map[string]string)}
}

func (m *mockCache) Get(ctx context.Context, key string) (string, bool) {
	v, ok := m.store[key]
	return v, ok
}

func (m *mockCache) Set(ctx context.Context, key string, value string) {
	m.store[key] = value
}

func (m *mockCache) Delete(ctx context.Context, key string) {
	delete(m.store, key)
}

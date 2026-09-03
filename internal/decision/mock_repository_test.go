package decision

import (
	"context"

	"github.com/Olamigokeolowo/projectflow-backend/internal/events"
)

// mockRepository is a test double implementing the Repository interface.
type mockRepository struct {
	listByWorkspaceFunc func(ctx context.Context, workspaceID string) ([]*Decision, error)
	getByIDFunc         func(ctx context.Context, id string) (*Decision, error)
	createFunc          func(ctx context.Context, title, status, ownerID, workspaceID string) (*Decision, error)
	slowOperationFunc   func(ctx context.Context) error
}

func (m *mockRepository) ListByWorkspace(ctx context.Context, workspaceID string) ([]*Decision, error) {
	return m.listByWorkspaceFunc(ctx, workspaceID)
}

func (m *mockRepository) GetByID(ctx context.Context, id string) (*Decision, error) {
	return m.getByIDFunc(ctx, id)
}

func (m *mockRepository) Create(ctx context.Context, title, status, ownerID, workspaceID string) (*Decision, error) {
	return m.createFunc(ctx, title, status, ownerID, workspaceID)
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

// mockCache is a test double implementing the cache.Cache interface.
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

// mockMembershipChecker is a test double implementing the MembershipChecker interface.
type mockMembershipChecker struct {
	isMemberFunc func(ctx context.Context, workspaceID, userID string) (bool, error)
}

func (m *mockMembershipChecker) IsMember(ctx context.Context, workspaceID, userID string) (bool, error) {
	if m.isMemberFunc != nil {
		return m.isMemberFunc(ctx, workspaceID, userID)
	}
	return true, nil // default: always a member, unless a test overrides this
}
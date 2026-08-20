package decision

import "context"

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

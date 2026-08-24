package decision

import (
	"context"
	"testing"

	"github.com/Olamigokeolowo/projectflow-backend/internal/events"
	"github.com/stretchr/testify/assert"
)

func TestService_List(t *testing.T) {
	repo := &mockRepository{
		listFunc: func(ctx context.Context) ([]*Decision, error) {
			return []*Decision{{ID: "1", Title: "Test Decision"}}, nil
		},
	}
	service := NewService(repo, &mockPublisher{}, newMockCache())

	decisions, err := service.List(context.Background())

	assert.NoError(t, err)
	assert.Len(t, decisions, 1)
	assert.Equal(t, "Test Decision", decisions[0].Title)
}

func TestService_Get_Success(t *testing.T) {
	repo := &mockRepository{
		getByIDFunc: func(ctx context.Context, id string) (*Decision, error) {
			return &Decision{ID: id, OwnerID: "user-123"}, nil
		},
	}
	service := NewService(repo, &mockPublisher{}, newMockCache())

	d, err := service.Get(context.Background(), "decision-1", "user-123")

	assert.NoError(t, err)
	assert.Equal(t, "user-123", d.OwnerID)
}

func TestService_Get_Forbidden(t *testing.T) {
	repo := &mockRepository{
		getByIDFunc: func(ctx context.Context, id string) (*Decision, error) {
			return &Decision{ID: id, OwnerID: "user-123"}, nil // owned by a different user
		},
	}
	service := NewService(repo, &mockPublisher{}, newMockCache())

	_, err := service.Get(context.Background(), "decision-1", "user-999") // different requester

	assert.ErrorIs(t, err, ErrForbidden)
}

func TestService_Get_NotFound(t *testing.T) {
	repo := &mockRepository{
		getByIDFunc: func(ctx context.Context, id string) (*Decision, error) {
			return nil, ErrNotFound
		},
	}
	service := NewService(repo, &mockPublisher{}, newMockCache())

	_, err := service.Get(context.Background(), "missing-id", "user-123")

	assert.ErrorIs(t, err, ErrNotFound)
}

func TestService_Create(t *testing.T) {
	repo := &mockRepository{
		createFunc: func(ctx context.Context, title, status, ownerID string) (*Decision, error) {
			return &Decision{ID: "new-id", Title: title, Status: status, OwnerID: ownerID}, nil
		},
	}
	service := NewService(repo, &mockPublisher{}, newMockCache())

	d, err := service.Create(context.Background(), "New Decision", "draft", "user-123")

	assert.NoError(t, err)
	assert.Equal(t, "New Decision", d.Title)
	assert.Equal(t, "user-123", d.OwnerID)
}

func TestService_Create_PublishesEvent(t *testing.T) {
	var published bool

	repo := &mockRepository{
		createFunc: func(ctx context.Context, title, status, ownerID string) (*Decision, error) {
			return &Decision{ID: "new-id", Title: title, Status: status, OwnerID: ownerID}, nil
		},
	}
	publisher := &mockPublisher{
		publishFunc: func(ctx context.Context, event events.DecisionCreated) error {
			published = true
			assert.Equal(t, "new-id", event.DecisionID)
			return nil
		},
	}
	service := NewService(repo, publisher, newMockCache())

	_, err := service.Create(context.Background(), "New Decision", "draft", "user-123")

	assert.NoError(t, err)
	assert.True(t, published)
}
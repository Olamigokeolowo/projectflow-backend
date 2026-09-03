package decision

import (
	"context"
	"testing"

	"github.com/Olamigokeolowo/projectflow-backend/internal/events"
	"github.com/stretchr/testify/assert"
)

func TestService_ListByWorkspace(t *testing.T) {
	repo := &mockRepository{
		listByWorkspaceFunc: func(ctx context.Context, workspaceID string) ([]*Decision, error) {
			return []*Decision{{ID: "1", Title: "Test Decision", WorkspaceID: workspaceID}}, nil
		},
	}
	membership := &mockMembershipChecker{}
	service := NewService(repo, &mockPublisher{}, newMockCache(), membership)

	decisions, err := service.ListByWorkspace(context.Background(), "ws-1", "user-123")

	assert.NoError(t, err)
	assert.Len(t, decisions, 1)
	assert.Equal(t, "Test Decision", decisions[0].Title)
}

func TestService_Get_Success(t *testing.T) {
	repo := &mockRepository{
		getByIDFunc: func(ctx context.Context, id string) (*Decision, error) {
			return &Decision{ID: id, OwnerID: "user-123", WorkspaceID: "ws-1"}, nil
		},
	}
	membership := &mockMembershipChecker{}
	service := NewService(repo, &mockPublisher{}, newMockCache(), membership)

	d, err := service.Get(context.Background(), "decision-1", "user-123")

	assert.NoError(t, err)
	assert.Equal(t, "user-123", d.OwnerID)
}

func TestService_Get_Forbidden(t *testing.T) {
	repo := &mockRepository{
		getByIDFunc: func(ctx context.Context, id string) (*Decision, error) {
			return &Decision{ID: id, OwnerID: "user-123", WorkspaceID: "ws-1"}, nil
		},
	}
	membership := &mockMembershipChecker{
		isMemberFunc: func(ctx context.Context, workspaceID, userID string) (bool, error) {
			return false, nil
		},
	}
	service := NewService(repo, &mockPublisher{}, newMockCache(), membership)

	_, err := service.Get(context.Background(), "decision-1", "user-999")

	assert.ErrorIs(t, err, ErrForbidden)
}

func TestService_Get_NotFound(t *testing.T) {
	repo := &mockRepository{
		getByIDFunc: func(ctx context.Context, id string) (*Decision, error) {
			return nil, ErrNotFound
		},
	}
	membership := &mockMembershipChecker{}
	service := NewService(repo, &mockPublisher{}, newMockCache(), membership)

	_, err := service.Get(context.Background(), "missing-id", "user-123")

	assert.ErrorIs(t, err, ErrNotFound)
}

func TestService_Create(t *testing.T) {
	repo := &mockRepository{
		createFunc: func(ctx context.Context, title, status, ownerID, workspaceID string) (*Decision, error) {
			return &Decision{ID: "new-id", Title: title, Status: status, OwnerID: ownerID, WorkspaceID: workspaceID}, nil
		},
	}
	membership := &mockMembershipChecker{}
	service := NewService(repo, &mockPublisher{}, newMockCache(), membership)

	d, err := service.Create(context.Background(), "New Decision", "draft", "user-123", "ws-1")

	assert.NoError(t, err)
	assert.Equal(t, "New Decision", d.Title)
	assert.Equal(t, "user-123", d.OwnerID)
	assert.Equal(t, "ws-1", d.WorkspaceID)
}

func TestService_Create_PublishesEvent(t *testing.T) {
	var published bool

	repo := &mockRepository{
		createFunc: func(ctx context.Context, title, status, ownerID, workspaceID string) (*Decision, error) {
			return &Decision{ID: "new-id", Title: title, Status: status, OwnerID: ownerID, WorkspaceID: workspaceID}, nil
		},
	}
	publisher := &mockPublisher{
		publishFunc: func(ctx context.Context, event events.DecisionCreated) error {
			published = true
			assert.Equal(t, "new-id", event.DecisionID)
			return nil
		},
	}
	membership := &mockMembershipChecker{}
	service := NewService(repo, publisher, newMockCache(), membership)

	_, err := service.Create(context.Background(), "New Decision", "draft", "user-123", "ws-1")

	assert.NoError(t, err)
	assert.True(t, published)
}

func TestService_Create_ForbiddenWhenNotMember(t *testing.T) {
	repo := &mockRepository{}
	membership := &mockMembershipChecker{
		isMemberFunc: func(ctx context.Context, workspaceID, userID string) (bool, error) {
			return false, nil
		},
	}
	service := NewService(repo, &mockPublisher{}, newMockCache(), membership)

	_, err := service.Create(context.Background(), "New Decision", "draft", "user-123", "ws-1")

	assert.ErrorIs(t, err, ErrForbidden)
}
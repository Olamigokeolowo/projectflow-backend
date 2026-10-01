package workspace

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("workspace not found")

type Repository interface {
	Create(ctx context.Context, name, createdBy string) (*Workspace, error)
	GetByID(ctx context.Context, id string) (*Workspace, error)
	ListForUser(ctx context.Context, userID string) ([]*Workspace, error)
	AddMember(ctx context.Context, workspaceID, userID string) error
	RemoveMember(ctx context.Context, workspaceID, userID string) error
	IsMember(ctx context.Context, workspaceID, userID string) (bool, error)
	GetRole(ctx context.Context, workspaceID, userID string) (string, error)
	ListMembers(ctx context.Context, workspaceID string) ([]*Membership, error)
}
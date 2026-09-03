package workspace

import (
	"context"
	"errors"
)

var ErrForbidden = errors.New("you are not a member of this workspace")

type UserFinder interface {
	FindIDByEmail(ctx context.Context, email string) (string, error)
}

type Service struct {
	repo       Repository
	userFinder UserFinder
}

func NewService(repo Repository, userFinder UserFinder) *Service {
	return &Service{repo: repo, userFinder: userFinder}
}

func (s *Service) Create(ctx context.Context, name, createdBy string) (*Workspace, error) {
	return s.repo.Create(ctx, name, createdBy)
}

func (s *Service) ListForUser(ctx context.Context, userID string) ([]*Workspace, error) {
	return s.repo.ListForUser(ctx, userID)
}

func (s *Service) AddMember(ctx context.Context, workspaceID, requestingUserID, email string) error {
	isMember, err := s.repo.IsMember(ctx, workspaceID, requestingUserID)
	if err != nil {
		return err
	}
	if !isMember {
		return ErrForbidden
	}

	targetUserID, err := s.userFinder.FindIDByEmail(ctx, email)
	if err != nil {
		return err
	}
	return s.repo.AddMember(ctx, workspaceID, targetUserID)
}

func (s *Service) ListMembers(ctx context.Context, workspaceID, requestingUserID string) ([]string, error) {
	isMember, err := s.repo.IsMember(ctx, workspaceID, requestingUserID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, ErrForbidden
	}
	return s.repo.ListMembers(ctx, workspaceID)
}

func (s *Service) IsMember(ctx context.Context, workspaceID, userID string) (bool, error) {
	return s.repo.IsMember(ctx, workspaceID, userID)
}


package comment

import (
	"context"
	"errors"
)

var ErrForbidden = errors.New("you do not have access to this workspace")

type WorkspaceResolver interface {
	GetWorkspaceID(ctx context.Context, targetType, targetID string) (string, error)
}

type MembershipChecker interface {
	IsMember(ctx context.Context, workspaceID, userID string) (bool, error)
}

type Service struct {
	repo       Repository
	resolver   WorkspaceResolver
	membership MembershipChecker
}

func NewService(repo Repository, resolver WorkspaceResolver, membership MembershipChecker) *Service {
	return &Service{repo: repo, resolver: resolver, membership: membership}
}

func (s *Service) authorize(ctx context.Context, targetType, targetID, userID string) error {
	workspaceID, err := s.resolver.GetWorkspaceID(ctx, targetType, targetID)
	if err != nil {
		return err
	}

	isMember, err := s.membership.IsMember(ctx, workspaceID, userID)
	if err != nil {
		return err
	}
	if !isMember {
		return ErrForbidden
	}
	return nil
}

func (s *Service) ListByTarget(ctx context.Context, targetType, targetID, requestingUserID string) ([]*Comment, error) {
	if err := s.authorize(ctx, targetType, targetID, requestingUserID); err != nil {
		return nil, err
	}
	return s.repo.ListByTarget(ctx, targetType, targetID)
}

func (s *Service) Create(ctx context.Context, body, authorID, targetType, targetID string) (*Comment, error) {
	if err := s.authorize(ctx, targetType, targetID, authorID); err != nil {
		return nil, err
	}
	return s.repo.Create(ctx, body, authorID, targetType, targetID)
}

func (s *Service) Delete(ctx context.Context, id, requestingUserID string) error {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.authorize(ctx, c.TargetType, c.TargetID, requestingUserID); err != nil {
		return err
	}

	if c.AuthorID != requestingUserID {
		return ErrForbidden
	}

	return s.repo.Delete(ctx, id)
}
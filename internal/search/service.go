package search

import (
	"context"
	"errors"
)

var ErrForbidden = errors.New("you do not have access to this workspace")

type MembershipChecker interface {
	IsMember(ctx context.Context, workspaceID, userID string) (bool, error)
}

type Service struct {
	repo       *Repository
	membership MembershipChecker
}

func NewService(repo *Repository, membership MembershipChecker) *Service {
	return &Service{repo: repo, membership: membership}
}

func (s *Service) Search(ctx context.Context, workspaceID, query, userID string) (*Results, error) {
	isMember, err := s.membership.IsMember(ctx, workspaceID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, ErrForbidden
	}
	return s.repo.Search(ctx, workspaceID, query)
}
package activity

import (
	"context"
	"log"
)

type MembershipChecker interface {
	IsMember(ctx context.Context, workspaceID, userID string) (bool, error)
}

type Service struct {
	repo       Repository
	membership MembershipChecker
}

func NewService(repo Repository, membership MembershipChecker) *Service {
	return &Service{repo: repo, membership: membership}
}

// Record is best-effort — a failure here should never break the action that triggered it.
func (s *Service) Record(ctx context.Context, workspaceID, actorID, verb, targetType, targetID, summary string) {
	if _, err := s.repo.Create(ctx, workspaceID, actorID, verb, targetType, targetID, summary); err != nil {
		log.Println("failed to record activity:", err)
	}
}

var ErrForbidden = errForbidden{}

type errForbidden struct{}

func (errForbidden) Error() string { return "you do not have access to this workspace" }

func (s *Service) ListByWorkspace(ctx context.Context, workspaceID, requestingUserID string) ([]*Activity, error) {
	isMember, err := s.membership.IsMember(ctx, workspaceID, requestingUserID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, ErrForbidden
	}
	return s.repo.ListByWorkspace(ctx, workspaceID)
}
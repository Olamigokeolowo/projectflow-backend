package task

import (
	"context"
	"errors"
)

var ErrForbidden = errors.New("you do not have access to this workspace")

// DecisionFinder is the narrow slice of decision.Service that task actually
// needs — just enough to resolve which workspace a decision belongs to.
type DecisionFinder interface {
	GetWorkspaceID(ctx context.Context, decisionID string) (string, error)
}

// MembershipChecker mirrors the same interface decision.Service depends on.
type MembershipChecker interface {
	IsMember(ctx context.Context, workspaceID, userID string) (bool, error)
}

type Service struct {
	repo       Repository
	decisions  DecisionFinder
	membership MembershipChecker
}

func NewService(repo Repository, decisions DecisionFinder, membership MembershipChecker) *Service {
	return &Service{repo: repo, decisions: decisions, membership: membership}
}

func (s *Service) authorize(ctx context.Context, decisionID, userID string) error {
	workspaceID, err := s.decisions.GetWorkspaceID(ctx, decisionID)
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

func (s *Service) ListByDecision(ctx context.Context, decisionID, requestingUserID string) ([]*Task, error) {
	if err := s.authorize(ctx, decisionID, requestingUserID); err != nil {
		return nil, err
	}
	return s.repo.ListByDecision(ctx, decisionID)
}

func (s *Service) Create(ctx context.Context, title, decisionID, assigneeID, createdBy string) (*Task, error) {
	if err := s.authorize(ctx, decisionID, createdBy); err != nil {
		return nil, err
	}
	return s.repo.Create(ctx, title, decisionID, assigneeID, createdBy)
}

func (s *Service) Update(ctx context.Context, id, requestingUserID string, status, assigneeID *string) (*Task, error) {
	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := s.authorize(ctx, t.DecisionID, requestingUserID); err != nil {
		return nil, err
	}

	if status != nil {
		t.Status = *status
	}
	if assigneeID != nil {
		t.AssigneeID = *assigneeID
	}

	return s.repo.Update(ctx, t)
}

func (s *Service) Delete(ctx context.Context, id, requestingUserID string) error {
	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.authorize(ctx, t.DecisionID, requestingUserID); err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}
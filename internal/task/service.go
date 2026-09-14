package task

import (
	"context"
	"errors"

	"github.com/Olamigokeolowo/projectflow-backend/internal/activity"
)

var ErrForbidden = errors.New("you do not have access to this workspace")

type DecisionFinder interface {
	GetWorkspaceID(ctx context.Context, decisionID string) (string, error)
}

type MembershipChecker interface {
	IsMember(ctx context.Context, workspaceID, userID string) (bool, error)
}

type Service struct {
	repo       Repository
	decisions  DecisionFinder
	membership MembershipChecker
	activity   *activity.Service
}

func NewService(repo Repository, decisions DecisionFinder, membership MembershipChecker, act *activity.Service) *Service {
	return &Service{repo: repo, decisions: decisions, membership: membership, activity: act}
}

func (s *Service) authorize(ctx context.Context, decisionID, userID string) (string, error) {
	workspaceID, err := s.decisions.GetWorkspaceID(ctx, decisionID)
	if err != nil {
		return "", err
	}

	isMember, err := s.membership.IsMember(ctx, workspaceID, userID)
	if err != nil {
		return "", err
	}
	if !isMember {
		return "", ErrForbidden
	}
	return workspaceID, nil
}

func (s *Service) ListByDecision(ctx context.Context, decisionID, requestingUserID string) ([]*Task, error) {
	if _, err := s.authorize(ctx, decisionID, requestingUserID); err != nil {
		return nil, err
	}
	return s.repo.ListByDecision(ctx, decisionID)
}

func (s *Service) Create(ctx context.Context, title, decisionID, assigneeID, createdBy string) (*Task, error) {
	workspaceID, err := s.authorize(ctx, decisionID, createdBy)
	if err != nil {
		return nil, err
	}

	t, err := s.repo.Create(ctx, title, decisionID, assigneeID, createdBy)
	if err != nil {
		return nil, err
	}

	s.activity.Record(ctx, workspaceID, createdBy, "created", "task", t.ID, "created task \""+t.Title+"\"")

	return t, nil
}

func (s *Service) Update(ctx context.Context, id, requestingUserID string, status, assigneeID *string) (*Task, error) {
	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	workspaceID, err := s.authorize(ctx, t.DecisionID, requestingUserID)
	if err != nil {
		return nil, err
	}

	if status != nil {
		t.Status = *status
	}
	if assigneeID != nil {
		t.AssigneeID = *assigneeID
	}

	updated, err := s.repo.Update(ctx, t)
	if err != nil {
		return nil, err
	}

	s.activity.Record(ctx, workspaceID, requestingUserID, "updated", "task", t.ID, "updated task \""+t.Title+"\"")

	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id, requestingUserID string) error {
	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	workspaceID, err := s.authorize(ctx, t.DecisionID, requestingUserID)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	s.activity.Record(ctx, workspaceID, requestingUserID, "deleted", "task", t.ID, "deleted task \""+t.Title+"\"")

	return nil
}

func (s *Service) GetDecisionID(ctx context.Context, taskID string) (string, error) {
	t, err := s.repo.GetByID(ctx, taskID)
	if err != nil {
		return "", err
	}
	return t.DecisionID, nil
}
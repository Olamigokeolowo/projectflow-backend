package workspace

import (
	"context"
	"errors"

	"github.com/Olamigokeolowo/projectflow-backend/internal/activity"
)

var ErrForbidden = errors.New("you are not a member of this workspace")

type UserFinder interface {
	FindIDByEmail(ctx context.Context, email string) (string, error)
}

type Service struct {
	repo       Repository
	userFinder UserFinder
	activity   *activity.Service
}

func NewService(repo Repository, userFinder UserFinder, act *activity.Service) *Service {
	return &Service{repo: repo, userFinder: userFinder, activity: act}
}

func (s *Service) Create(ctx context.Context, name, createdBy string) (*Workspace, error) {
	w, err := s.repo.Create(ctx, name, createdBy)
	if err != nil {
		return nil, err
	}
	s.activity.Record(ctx, w.ID, createdBy, "created", "workspace", w.ID, "created workspace \""+w.Name+"\"")
	return w, nil
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

	if err := s.repo.AddMember(ctx, workspaceID, targetUserID); err != nil {
		return err
	}

	s.activity.Record(ctx, workspaceID, requestingUserID, "added", "member", targetUserID, "added "+email+" to the workspace")

	return nil
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
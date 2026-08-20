package decision

import (
	"context"
	"errors"
)

var ErrForbidden = errors.New("you do not have access to this decision")

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context) ([]*Decision, error) {
	return s.repo.List(ctx)
}

func (s *Service) Get(ctx context.Context, id, requestingUserID string) (*Decision, error) {
	d, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if d.OwnerID != requestingUserID {
		return nil, ErrForbidden
	}

	return d, nil
}

func (s *Service) Create(ctx context.Context, title, status, ownerID string) (*Decision, error) {
	return s.repo.Create(ctx, title, status, ownerID)
}

func (s *Service) SlowOperation(ctx context.Context) error {
	return s.repo.SlowOperation(ctx)
}
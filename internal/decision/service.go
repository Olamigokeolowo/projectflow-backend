package decision

import "context"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context) ([]*Decision, error) {
	return s.repo.List(ctx)
}

func (s *Service) Get(ctx context.Context, id string) (*Decision, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, title, status string) (*Decision, error) {
	return s.repo.Create(ctx, title, status)
}
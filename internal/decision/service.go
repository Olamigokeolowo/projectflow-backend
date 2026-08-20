package decision

import (
	"context"
	"errors"
	"log"

	"github.com/Olamigokeolowo/projectflow-backend/internal/events"
)

var ErrForbidden = errors.New("you do not have access to this decision")

type Service struct {
	repo      Repository
	publisher events.Publisher
}

func NewService(repo Repository, publisher events.Publisher) *Service {
	return &Service{repo: repo, publisher: publisher}
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
	d, err := s.repo.Create(ctx, title, status, ownerID)
	if err != nil {
		return nil, err
	}

	if pubErr := s.publisher.Publish(ctx, events.DecisionCreated{
		DecisionID: d.ID,
		OwnerID:    d.OwnerID,
		Title:      d.Title,
	}); pubErr != nil {
		log.Println("failed to publish DecisionCreated event:", pubErr)
	}

	return d, nil
}

func (s *Service) SlowOperation(ctx context.Context) error {
	return s.repo.SlowOperation(ctx)
}
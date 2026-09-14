package decision

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/Olamigokeolowo/projectflow-backend/internal/activity"
	"github.com/Olamigokeolowo/projectflow-backend/internal/cache"
	"github.com/Olamigokeolowo/projectflow-backend/internal/events"
)

var ErrForbidden = errors.New("you do not have access to this workspace")
var ErrNotOwner = errors.New("only the owner can modify this decision")

type MembershipChecker interface {
	IsMember(ctx context.Context, workspaceID, userID string) (bool, error)
}

type Service struct {
	repo       Repository
	publisher  events.Publisher
	cache      cache.Cache
	membership MembershipChecker
	activity   *activity.Service
}

func NewService(repo Repository, publisher events.Publisher, c cache.Cache, membership MembershipChecker, act *activity.Service) *Service {
	return &Service{repo: repo, publisher: publisher, cache: c, membership: membership, activity: act}
}

func (s *Service) ListByWorkspace(ctx context.Context, workspaceID, requestingUserID string) ([]*Decision, error) {
	isMember, err := s.membership.IsMember(ctx, workspaceID, requestingUserID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, ErrForbidden
	}
	return s.repo.ListByWorkspace(ctx, workspaceID)
}

func (s *Service) Get(ctx context.Context, id, requestingUserID string) (*Decision, error) {
	cacheKey := "decision:" + id

	if cached, found := s.cache.Get(ctx, cacheKey); found {
		var d Decision
		if err := json.Unmarshal([]byte(cached), &d); err == nil {
			isMember, err := s.membership.IsMember(ctx, d.WorkspaceID, requestingUserID)
			if err != nil {
				return nil, err
			}
			if !isMember {
				return nil, ErrForbidden
			}
			log.Println("cache hit:", cacheKey)
			return &d, nil
		}
	}

	log.Println("cache miss:", cacheKey)
	d, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	isMember, err := s.membership.IsMember(ctx, d.WorkspaceID, requestingUserID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, ErrForbidden
	}

	if serialized, err := json.Marshal(d); err == nil {
		s.cache.Set(ctx, cacheKey, string(serialized))
	}
	return d, nil
}

func (s *Service) Create(ctx context.Context, title, status, ownerID, workspaceID string) (*Decision, error) {
	isMember, err := s.membership.IsMember(ctx, workspaceID, ownerID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, ErrForbidden
	}

	d, err := s.repo.Create(ctx, title, status, ownerID, workspaceID)
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

	s.activity.Record(ctx, workspaceID, ownerID, "created", "decision", d.ID, "created decision \""+d.Title+"\"")

	return d, nil
}

func (s *Service) Update(ctx context.Context, id, requestingUserID string, title, status *string) (*Decision, error) {
	d, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	isMember, err := s.membership.IsMember(ctx, d.WorkspaceID, requestingUserID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, ErrForbidden
	}

	if d.OwnerID != requestingUserID {
		return nil, ErrNotOwner
	}

	if title != nil {
		d.Title = *title
	}
	if status != nil {
		d.Status = *status
	}

	updated, err := s.repo.Update(ctx, d)
	if err != nil {
		return nil, err
	}

	s.cache.Delete(ctx, "decision:"+id)
	s.activity.Record(ctx, d.WorkspaceID, requestingUserID, "updated", "decision", d.ID, "updated decision \""+d.Title+"\"")

	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id, requestingUserID string) error {
	d, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	isMember, err := s.membership.IsMember(ctx, d.WorkspaceID, requestingUserID)
	if err != nil {
		return err
	}
	if !isMember {
		return ErrForbidden
	}

	if d.OwnerID != requestingUserID {
		return ErrNotOwner
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	s.cache.Delete(ctx, "decision:"+id)
	s.activity.Record(ctx, d.WorkspaceID, requestingUserID, "deleted", "decision", d.ID, "deleted decision \""+d.Title+"\"")

	return nil
}

func (s *Service) GetWorkspaceID(ctx context.Context, decisionID string) (string, error) {
	d, err := s.repo.GetByID(ctx, decisionID)
	if err != nil {
		return "", err
	}
	return d.WorkspaceID, nil
}

func (s *Service) SlowOperation(ctx context.Context) error {
	return s.repo.SlowOperation(ctx)
}
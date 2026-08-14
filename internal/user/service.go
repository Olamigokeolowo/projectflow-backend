package user

import (
	"context"
	"errors"

	"github.com/Olamigokeolowo/projectflow-backend/internal/auth"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(ctx context.Context, email, password string) (*AuthResponse, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u, err := s.repo.Create(ctx, email, string(hash))
	if err != nil {
		return nil, err
	}

	token, err := auth.GenerateToken(u.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{Token: token, User: u}, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*AuthResponse, error) {
	u, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials // deliberately vague — see note below
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := auth.GenerateToken(u.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{Token: token, User: u}, nil
}
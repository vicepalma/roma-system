package service

import (
	"context"
	"errors"
	"strings"

	"github.com/vicepalma/roma-system/backend/internal/repository"
)

type AuthService interface {
	Signup(ctx context.Context, email, name, password, roleHint string) (*repository.UserBasic, error)
	Me(ctx context.Context, userID string) (*MeResponse, error)
}

type authService struct {
	repo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) AuthService {
	return &authService{repo: repo}
}

func (s *authService) Signup(ctx context.Context, email, name, password, roleHint string) (*repository.UserBasic, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	name = strings.TrimSpace(name)
	if email == "" || name == "" || password == "" {
		return nil, errors.New("email, name and password are required")
	}

	// repo ya valida único por email
	u, err := s.repo.CreateUser(ctx, email, name, password)
	if err != nil {
		return nil, err
	}
	return u, nil
}

type MeResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"` // coach | disciple (derivado)
}

func (s *authService) Me(ctx context.Context, userID string) (*MeResponse, error) {
	u, err := s.repo.GetUserBasic(ctx, userID)
	if err != nil {
		return nil, err
	}
	role, err := s.repo.DeriveRole(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &MeResponse{
		ID:    u.ID,
		Email: u.Email,
		Name:  u.Name,
		Role:  role,
	}, nil
}

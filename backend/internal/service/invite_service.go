package service

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"errors"
	"strings"
	"time"

	"github.com/vicepalma/roma-system/backend/internal/repository"
)

type InviteService interface {
	CreateInvite(ctx context.Context, coachID, email, name string, ttlHours int) (*InviteDTO, error)
	AcceptInvite(ctx context.Context, code string, discipleID string) (*AcceptResult, error)
}

type InviteDTO struct {
	Code      string    `json:"code"`
	InviteURL string    `json:"invite_url"`
	ExpiresAt time.Time `json:"expires_at"`
	Email     string    `json:"email"`
	Name      *string   `json:"name,omitempty"`
}

type AcceptResult struct {
	LinkID string `json:"link_id"`
	Status string `json:"status"` // accepted
}

type inviteService struct {
	inv     repository.InviteRepository
	coach   CoachService
	baseURL string // ej: http://localhost:5173/invite/  (opcional)
}

func NewInviteService(inv repository.InviteRepository, coach CoachService, baseURL string) InviteService {
	return &inviteService{inv: inv, coach: coach, baseURL: strings.TrimRight(baseURL, "/")}
}

func randCode(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// Base32 sin padding, mayúsculas; recorta a ~10-12 chars
	return strings.TrimRight(base32.StdEncoding.EncodeToString(b), "="), nil
}

func (s *inviteService) CreateInvite(ctx context.Context, coachID, email, name string, ttlHours int) (*InviteDTO, error) {
	if coachID == "" || email == "" {
		return nil, errors.New("coach_id and email required")
	}
	if ttlHours <= 0 {
		ttlHours = 72
	}
	code, err := randCode(8) // ~13 chars
	if err != nil {
		return nil, err
	}
	email = strings.ToLower(strings.TrimSpace(email))
	var namePtr *string
	if strings.TrimSpace(name) != "" {
		n := strings.TrimSpace(name)
		namePtr = &n
	}
	inv := &repository.Invitation{
		Code:      code,
		CoachID:   coachID,
		Email:     email,
		Name:      namePtr,
		Status:    "pending",
		ExpiresAt: time.Now().Add(time.Duration(ttlHours) * time.Hour),
	}
	if err := s.inv.Create(ctx, inv); err != nil {
		return nil, err
	}
	url := code
	if s.baseURL != "" {
		url = s.baseURL + "/" + code
	}
	return &InviteDTO{
		Code:      code,
		InviteURL: url,
		ExpiresAt: inv.ExpiresAt,
		Email:     email,
		Name:      namePtr,
	}, nil
}

func (s *inviteService) AcceptInvite(ctx context.Context, code string, discipleID string) (*AcceptResult, error) {
	if code == "" || discipleID == "" {
		return nil, errors.New("code and disciple_id required")
	}
	inv, err := s.inv.FindByCode(ctx, code)
	if err != nil {
		return nil, errors.New("invalid_code")
	}
	if inv.Status != "pending" {
		return nil, errors.New("invite_not_pending")
	}
	if time.Now().After(inv.ExpiresAt) {
		return nil, errors.New("invite_expired")
	}

	// Crea/acepta link coach-disciple
	link, err := s.coach.CreateLink(ctx, inv.CoachID, discipleID, true /*autoAccept*/)
	if err != nil {
		return nil, err
	}

	if err := s.inv.MarkAccepted(ctx, inv.ID, discipleID, time.Now()); err != nil {
		return nil, err
	}
	return &AcceptResult{LinkID: link.ID, Status: "accepted"}, nil
}

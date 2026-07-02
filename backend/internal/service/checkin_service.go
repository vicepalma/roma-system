package service

import (
	"context"
	"errors"
	"time"

	"github.com/vicepalma/roma-system/backend/internal/domain"
	"github.com/vicepalma/roma-system/backend/internal/repository"
)

var ErrInvalidCheckin = errors.New("invalid_checkin")

type CheckinService interface {
	Create(ctx context.Context, discipleID string, checkedAt time.Time, weightKG *float64, notes *string) (*domain.Checkin, error)
	List(ctx context.Context, discipleID string, limit, offset int) ([]domain.Checkin, int64, error)
	Get(ctx context.Context, id string) (*domain.Checkin, error)
}

type checkinService struct{ repo repository.CheckinRepository }

func NewCheckinService(repo repository.CheckinRepository) CheckinService {
	return &checkinService{repo: repo}
}

func (s *checkinService) Create(ctx context.Context, discipleID string, checkedAt time.Time, weightKG *float64, notes *string) (*domain.Checkin, error) {
	if discipleID == "" || checkedAt.IsZero() {
		return nil, ErrInvalidCheckin
	}
	if weightKG != nil && *weightKG <= 0 {
		return nil, ErrInvalidCheckin
	}
	checkin := &domain.Checkin{
		DiscipleID: discipleID,
		CheckedAt:  checkedAt,
		WeightKG:   weightKG,
		Notes:      notes,
	}
	if err := s.repo.Create(ctx, checkin); err != nil {
		return nil, err
	}
	return checkin, nil
}

func (s *checkinService) List(ctx context.Context, discipleID string, limit, offset int) ([]domain.Checkin, int64, error) {
	return s.repo.ListByDisciple(ctx, discipleID, limit, offset)
}

func (s *checkinService) Get(ctx context.Context, id string) (*domain.Checkin, error) {
	return s.repo.FindByID(ctx, id)
}

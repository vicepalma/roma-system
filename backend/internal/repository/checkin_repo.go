package repository

import (
	"context"

	"github.com/vicepalma/roma-system/backend/internal/domain"
	"gorm.io/gorm"
)

type CheckinRepository interface {
	Create(ctx context.Context, checkin *domain.Checkin) error
	ListByDisciple(ctx context.Context, discipleID string, limit, offset int) ([]domain.Checkin, int64, error)
	FindByID(ctx context.Context, id string) (*domain.Checkin, error)
}

type checkinRepository struct{ db *gorm.DB }

func NewCheckinRepository(db *gorm.DB) CheckinRepository { return &checkinRepository{db: db} }

func (r *checkinRepository) Create(ctx context.Context, checkin *domain.Checkin) error {
	return r.db.WithContext(ctx).Create(checkin).Error
}

func (r *checkinRepository) ListByDisciple(ctx context.Context, discipleID string, limit, offset int) ([]domain.Checkin, int64, error) {
	q := r.db.WithContext(ctx).Model(&domain.Checkin{}).Where("disciple_id = ?", discipleID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	var rows []domain.Checkin
	err := q.Order("checked_at DESC").Order("created_at DESC").Limit(limit).Offset(offset).Find(&rows).Error
	return rows, total, err
}

func (r *checkinRepository) FindByID(ctx context.Context, id string) (*domain.Checkin, error) {
	var out domain.Checkin
	if err := r.db.WithContext(ctx).First(&out, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &out, nil
}

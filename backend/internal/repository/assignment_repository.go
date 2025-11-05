package repository

import (
	"context"
	"time"

	"github.com/vicepalma/roma-system/backend/internal/domain"
	"gorm.io/gorm"
)

type AssignmentRepository interface {
	WithTx(tx *gorm.DB) AssignmentRepository
	FindByID(ctx context.Context, id string) (*domain.Assignment, error)
	DeactivateAllForDisciple(ctx context.Context, discipleID string) error
	ActivateOne(ctx context.Context, assignmentID, discipleID string) error
}

type assignmentRepository struct {
	db *gorm.DB
}

func NewAssignmentRepository(db *gorm.DB) AssignmentRepository { return &assignmentRepository{db: db} }
func (r *assignmentRepository) WithTx(tx *gorm.DB) AssignmentRepository {
	return &assignmentRepository{db: tx}
}

func (r *assignmentRepository) FindByID(ctx context.Context, id string) (*domain.Assignment, error) {
	var a domain.Assignment
	if err := r.db.WithContext(ctx).First(&a, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *assignmentRepository) DeactivateAllForDisciple(ctx context.Context, discipleID string) error {
	return r.db.WithContext(ctx).
		Model(&domain.Assignment{}).
		Where("disciple_id = ? AND is_active = true", discipleID).
		Updates(map[string]any{
			"is_active": false,
			"end_date":  gorm.Expr("COALESCE(end_date, ?::date)", time.Now()),
		}).Error
}

func (r *assignmentRepository) ActivateOne(ctx context.Context, assignmentID, discipleID string) error {
	return r.db.WithContext(ctx).
		Model(&domain.Assignment{}).
		Where("id = ? AND disciple_id = ?", assignmentID, discipleID).
		Updates(map[string]any{
			"is_active":  true,
			"end_date":   nil,
			"start_date": gorm.Expr("COALESCE(start_date, CURRENT_DATE)"),
		}).Error
}

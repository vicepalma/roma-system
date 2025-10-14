package repository

import (
	"context"
	"strings"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Exercise struct {
	ID            string         `gorm:"type:uuid;primaryKey" json:"id"`
	Name          string         `gorm:"not null" json:"name"`
	PrimaryMuscle string         `gorm:"not null" json:"primary_muscle"`
	Equipment     *string        `json:"equipment"`                         // NULL permitido
	Tags          pq.StringArray `gorm:"type:text[]" json:"tags,omitempty"` // default '{}' en BD
	Notes         *string        `json:"notes,omitempty"`
}

type ExerciseFilter struct {
	Query     string
	Muscle    string
	Equipment string
	Limit     int
	Offset    int
}

type ExerciseRepository interface {
	Search(ctx context.Context, f ExerciseFilter) ([]Exercise, int64, error)
	Create(ctx context.Context, e *Exercise) error
	Get(ctx context.Context, id string) (*Exercise, error)
	Update(ctx context.Context, id string, upd *Exercise) (*Exercise, error)
	Delete(ctx context.Context, id string) error
}

type exerciseRepository struct{ db *gorm.DB }

func NewExerciseRepository(db *gorm.DB) ExerciseRepository { return &exerciseRepository{db: db} }

func (r *exerciseRepository) Search(ctx context.Context, f ExerciseFilter) ([]Exercise, int64, error) {
	q := r.db.WithContext(ctx).Table("exercises")
	if s := strings.TrimSpace(f.Query); s != "" {
		ilike := "%" + strings.ToLower(s) + "%"
		q = q.Where("lower(name) LIKE ? OR lower(primary_muscle) LIKE ? OR lower(equipment) LIKE ?", ilike, ilike, ilike)
	}
	if s := strings.TrimSpace(f.Muscle); s != "" {
		q = q.Where("lower(primary_muscle) = ?", strings.ToLower(s))
	}
	if s := strings.TrimSpace(f.Equipment); s != "" {
		q = q.Where("lower(coalesce(equipment,'')) = ?", strings.ToLower(s))
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if f.Limit > 0 {
		q = q.Limit(f.Limit)
	}
	if f.Offset > 0 {
		q = q.Offset(f.Offset)
	}
	q = q.Order("lower(name) ASC, id ASC")

	var items []Exercise
	if err := q.Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *exerciseRepository) Create(ctx context.Context, e *Exercise) error {
	return r.db.WithContext(ctx).Create(e).Error
}

func (r *exerciseRepository) Get(ctx context.Context, id string) (*Exercise, error) {
	var ex Exercise
	if err := r.db.WithContext(ctx).First(&ex, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &ex, nil
}

func (r *exerciseRepository) Update(ctx context.Context, id string, upd *Exercise) (*Exercise, error) {
	var ex Exercise
	if err := r.db.WithContext(ctx).First(&ex, "id = ?", id).Error; err != nil {
		return nil, err
	}
	// Solo campos mutables
	ex.Name = upd.Name
	ex.PrimaryMuscle = upd.PrimaryMuscle
	ex.Equipment = upd.Equipment
	ex.Tags = upd.Tags
	ex.Notes = upd.Notes

	if err := r.db.WithContext(ctx).Save(&ex).Error; err != nil {
		return nil, err
	}
	return &ex, nil
}

func (r *exerciseRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&Exercise{}, "id = ?", id).Error
}

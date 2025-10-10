package repository

import (
	"context"
	"strings"

	"github.com/lib/pq"
	"github.com/vicepalma/roma-system/backend/internal/domain"
	"gorm.io/gorm"
)

type ExerciseRepository interface {
	Search(ctx context.Context, q string, tags []string, match string, limit, offset int) ([]domain.Exercise, int64, error)
}

type exerciseRepository struct{ db *gorm.DB }

func NewExerciseRepository(db *gorm.DB) ExerciseRepository {
	return &exerciseRepository{db: db}
}

func normalizeTags(tags []string) []string {
	out := make([]string, 0, len(tags))
	for _, t := range tags {
		t = strings.TrimSpace(t)
		if t != "" {
			out = append(out, t)
		}
	}
	return out
}

func applyFilters(db *gorm.DB, q string, tags []string, match string) *gorm.DB {
	// Texto: name o primary_muscle
	if q != "" {
		like := "%" + q + "%"
		db = db.Where(
			"lower(name) LIKE lower(?) OR lower(primary_muscle) LIKE lower(?)",
			like, like,
		)
	}

	// Tags: any/all
	tags = normalizeTags(tags)
	if len(tags) > 0 {
		arr := pq.StringArray(tags)
		if strings.EqualFold(match, "all") {
			// Contiene TODOS los tags
			db = db.Where("tags @> ?::text[]", arr)
		} else {
			// Contiene ALGUNO de los tags (default)
			db = db.Where("tags && ?::text[]", arr)
		}
	}

	return db
}

func (r *exerciseRepository) Search(ctx context.Context, q string, tags []string, match string, limit, offset int) ([]domain.Exercise, int64, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	base := r.db.WithContext(ctx).Model(&domain.Exercise{})
	base = applyFilters(base, q, tags, match)

	// total
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// items
	var rows []domain.Exercise
	if err := base.Order("name ASC").Limit(limit).Offset(offset).Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

package service

import (
	"context"

	"github.com/vicepalma/roma-system/backend/internal/domain"
	"github.com/vicepalma/roma-system/backend/internal/repository"
)

type ExerciseService interface {
	Search(ctx context.Context, q string, tags []string, match string, limit, offset int) ([]domain.Exercise, int64, error)
}

type exerciseService struct {
	repo repository.ExerciseRepository
}

func NewExerciseService(r repository.ExerciseRepository) ExerciseService {
	return &exerciseService{repo: r}
}

func (s *exerciseService) Search(ctx context.Context, q string, tags []string, match string, limit, offset int) ([]domain.Exercise, int64, error) {
	return s.repo.Search(ctx, q, tags, match, limit, offset)
}

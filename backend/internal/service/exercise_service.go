package service

import (
	"context"
	"errors"
	"strings"

	"github.com/lib/pq"
	"github.com/vicepalma/roma-system/backend/internal/repository"
)

type ExerciseService interface {
	List(ctx context.Context, f repository.ExerciseFilter) ([]repository.Exercise, int64, error)
	Create(ctx context.Context, in CreateExercise) (*repository.Exercise, error)
	Get(ctx context.Context, id string) (*repository.Exercise, error)
	Update(ctx context.Context, id string, in UpdateExercise) (*repository.Exercise, error)
	Delete(ctx context.Context, id string) error
}

type exerciseService struct {
	repo repository.ExerciseRepository
}

func NewExerciseService(r repository.ExerciseRepository) ExerciseService {
	return &exerciseService{repo: r}
}

// DTOs
type CreateExercise struct {
	Name          string   `json:"name"`
	PrimaryMuscle string   `json:"primary_muscle"`
	Equipment     *string  `json:"equipment"`
	Tags          []string `json:"tags"`
	Notes         *string  `json:"notes"`
}

type UpdateExercise = CreateExercise

func (s *exerciseService) List(ctx context.Context, f repository.ExerciseFilter) ([]repository.Exercise, int64, error) {
	// sane defaults
	if f.Limit < 0 {
		f.Limit = 0
	}
	if f.Offset < 0 {
		f.Offset = 0
	}
	return s.repo.Search(ctx, f)
}

func (s *exerciseService) Create(ctx context.Context, in CreateExercise) (*repository.Exercise, error) {
	if strings.TrimSpace(in.Name) == "" || strings.TrimSpace(in.PrimaryMuscle) == "" {
		return nil, errors.New("name and primary_muscle are required")
	}
	ex := &repository.Exercise{
		Name:          strings.TrimSpace(in.Name),
		PrimaryMuscle: strings.TrimSpace(in.PrimaryMuscle),
		Equipment:     normalizePtr(in.Equipment),
		Tags:          pq.StringArray(uniqueLower(in.Tags)),
		Notes:         normalizePtr(in.Notes),
	}
	if err := s.repo.Create(ctx, ex); err != nil {
		return nil, err
	}
	return ex, nil
}

func (s *exerciseService) Get(ctx context.Context, id string) (*repository.Exercise, error) {
	return s.repo.Get(ctx, id)
}

func (s *exerciseService) Update(ctx context.Context, id string, in UpdateExercise) (*repository.Exercise, error) {
	if strings.TrimSpace(in.Name) == "" || strings.TrimSpace(in.PrimaryMuscle) == "" {
		return nil, errors.New("name and primary_muscle are required")
	}
	upd := &repository.Exercise{
		Name:          strings.TrimSpace(in.Name),
		PrimaryMuscle: strings.TrimSpace(in.PrimaryMuscle),
		Equipment:     normalizePtr(in.Equipment),
		Tags:          pq.StringArray(uniqueLower(in.Tags)),
		Notes:         normalizePtr(in.Notes),
	}
	return s.repo.Update(ctx, id, upd)
}

func (s *exerciseService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func uniqueLower(ss []string) []string {
	seen := make(map[string]struct{}, len(ss))
	out := make([]string, 0, len(ss))
	for _, v := range ss {
		v = strings.TrimSpace(strings.ToLower(v))
		if v == "" {
			continue
		}
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			out = append(out, v)
		}
	}
	return out
}

func normalizePtr(p *string) *string {
	if p == nil {
		return nil
	}
	s := strings.TrimSpace(*p)
	if s == "" {
		return nil
	}
	return &s
}

package service

import (
	"context"
	"time"

	"github.com/vicepalma/roma-system/backend/internal/domain"
	"github.com/vicepalma/roma-system/backend/internal/repository"
)

type ProgramService interface {
	CreateProgram(ctx context.Context, ownerID, title string, notes *string) (*domain.Program, error)
	ListMyPrograms(ctx context.Context, ownerID string, limit, offset int) ([]domain.Program, int64, error)
	AddWeek(ctx context.Context, programID string, weekIndex int) (*domain.ProgramWeek, error)
	AddDay(ctx context.Context, weekID string, dayIndex int, notes *string) (*domain.ProgramDay, error)
	AddPrescription(ctx context.Context, p *domain.Prescription) (*domain.Prescription, error)
	Assign(ctx context.Context, programID, discipleID, assignedBy string, start time.Time, end *time.Time) (*domain.Assignment, error)
	MyToday(ctx context.Context, discipleID string, date time.Time) (*domain.ProgramDay, []domain.Prescription, error)
}

type programService struct{ repo repository.ProgramRepository }

func NewProgramService(r repository.ProgramRepository) ProgramService {
	return &programService{repo: r}
}

func (s *programService) CreateProgram(ctx context.Context, ownerID, title string, notes *string) (*domain.Program, error) {
	p := &domain.Program{OwnerID: ownerID, Title: title, Notes: notes}
	err := s.repo.CreateProgram(ctx, p)
	return p, err
}

func (s *programService) ListMyPrograms(ctx context.Context, ownerID string, limit, offset int) ([]domain.Program, int64, error) {
	return s.repo.ListMyPrograms(ctx, ownerID, limit, offset)
}

func (s *programService) AddWeek(ctx context.Context, programID string, weekIndex int) (*domain.ProgramWeek, error) {
	w := &domain.ProgramWeek{ProgramID: programID, WeekIndex: weekIndex}
	return w, s.repo.AddWeek(ctx, w)
}

func (s *programService) AddDay(ctx context.Context, weekID string, dayIndex int, notes *string) (*domain.ProgramDay, error) {
	d := &domain.ProgramDay{WeekID: weekID, DayIndex: dayIndex, Notes: notes}
	return d, s.repo.AddDay(ctx, d)
}

func (s *programService) AddPrescription(ctx context.Context, p *domain.Prescription) (*domain.Prescription, error) {
	return p, s.repo.AddPrescription(ctx, p)
}

func (s *programService) Assign(ctx context.Context, programID, discipleID, assignedBy string, start time.Time, end *time.Time) (*domain.Assignment, error) {
	version, err := s.repo.GetProgramVersion(ctx, programID)
	if err != nil {
		return nil, err
	}
	a := &domain.Assignment{
		ProgramID:      programID,
		ProgramVersion: version,
		DiscipleID:     discipleID,
		AssignedBy:     assignedBy,
		StartDate:      start,
		EndDate:        end,
		IsActive:       true,
	}
	return a, s.repo.Assign(ctx, a)
}

func (s *programService) MyToday(ctx context.Context, discipleID string, date time.Time) (*domain.ProgramDay, []domain.Prescription, error) {
	asg, err := s.repo.FindActiveAssignmentForDate(ctx, discipleID, date)
	if err != nil {
		return nil, nil, err
	}
	day, err := s.repo.FindDayForDate(ctx, asg, date)
	if err != nil {
		return nil, nil, err
	}
	prescs, err := s.repo.ListPrescriptionsByDay(ctx, day.ID)
	return day, prescs, err
}

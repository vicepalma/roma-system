package service

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/vicepalma/roma-system/backend/internal/domain"
	"github.com/vicepalma/roma-system/backend/internal/repository"
)

type CreateProgram struct {
	Title      string  `json:"title" binding:"required,min=2"`
	Notes      *string `json:"notes"`
	Visibility *string `json:"visibility"` // private|unlisted|public (futuro)
}

type UpdateProgram struct {
	Title      *string `json:"title"`
	Notes      *string `json:"notes"`
	Visibility *string `json:"visibility"`
}

type CreatePrescription struct {
	ExerciseID string   `json:"exercise_id" binding:"required"`
	Series     int      `json:"series" binding:"required,min=1"`
	Reps       string   `json:"reps" binding:"required"`
	RestSec    *int     `json:"rest_sec"`
	ToFailure  bool     `json:"to_failure"`
	Tempo      *string  `json:"tempo"`
	RIR        *int     `json:"rir"`
	RPE        *float32 `json:"rpe"`
	MethodID   *string  `json:"method_id"`
	Notes      *string  `json:"notes"`
	Position   *int     `json:"position"`
}

type UpdatePrescription struct {
	ExerciseID *string  `json:"exercise_id"`
	Series     *int     `json:"series"`
	Reps       *string  `json:"reps"`
	RestSec    *int     `json:"rest_sec"`
	ToFailure  *bool    `json:"to_failure"`
	Tempo      *string  `json:"tempo"`
	RIR        *int     `json:"rir"`
	RPE        *float32 `json:"rpe"`
	MethodID   *string  `json:"method_id"`
	Notes      *string  `json:"notes"`
	Position   *int     `json:"position"`
}

type ProgramService interface {
	CreateProgram(ctx context.Context, ownerID, title string, notes *string) (*domain.Program, error)
	ListMyPrograms(ctx context.Context, ownerID string, limit, offset int) ([]domain.Program, int64, error)
	AddWeek(ctx context.Context, programID string, weekIndex int) (*domain.ProgramWeek, error)
	DeleteWeek(ctx context.Context, programID, weekID string) error
	AddDay(ctx context.Context, weekID string, dayIndex int, notes *string) (*domain.ProgramDay, error)
	AddPrescription(ctx context.Context, p *domain.Prescription) (*domain.Prescription, error)
	Assign(ctx context.Context, programID, discipleID, assignedBy string, start time.Time, end *time.Time) (*domain.Assignment, error)
	MyToday(ctx context.Context, discipleID string, date time.Time) (*domain.ProgramDay, []domain.Prescription, error)

	List(ctx context.Context, f repository.ProgramFilter) ([]repository.Program, int64, error)
	Create(ctx context.Context, ownerID string, in CreateProgram) (*repository.Program, error)
	Get(ctx context.Context, id string) (*repository.Program, error)
	Update(ctx context.Context, id string, in UpdateProgram) (*repository.Program, error)
	Delete(ctx context.Context, id string) error

	NewVersion(ctx context.Context, programID string) (*repository.ProgramVersion, error)
	ListVersions(ctx context.Context, programID string) ([]repository.ProgramVersion, error)

	ListWeeks(ctx context.Context, programID string) ([]repository.ProgramWeek, error)
	ListDays(ctx context.Context, weekID string) ([]repository.ProgramDay, error)
	UpdateDay(ctx context.Context, dayID string, notes *string, dayIndex *int) (*repository.ProgramDay, error)
	DeleteDay(ctx context.Context, dayID string) error

	ListPrescriptions(ctx context.Context, dayID string) ([]repository.PrescriptionRow, error)
	UpdatePrescription(ctx context.Context, id string, in UpdatePrescription) (*repository.Prescription, error)
	DeletePrescription(ctx context.Context, id string) error
	ReorderPrescriptions(ctx context.Context, dayID string, orderedIDs []string) error

	GetProgram(ctx context.Context, id string) (*repository.ProgramRow, error)
	UpdateProgram(ctx context.Context, id string, title *string, notes *string, visibility *string) (*repository.ProgramRow, error)
	DeleteProgram(ctx context.Context, id string) error
	CreateNextVersionClone(ctx context.Context, programID string) (*repository.ProgramRow, error)
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

func (s *programService) DeleteWeek(ctx context.Context, programID, weekID string) error {
	// primero borrar los días asociados
	if err := s.repo.DeleteDaysByWeek(ctx, weekID); err != nil {
		return err
	}
	// luego la semana
	return s.repo.DeleteWeek(ctx, programID, weekID)
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

func (s *programService) List(ctx context.Context, f repository.ProgramFilter) ([]repository.Program, int64, error) {
	return s.repo.Search(ctx, f)
}

func (s *programService) Create(ctx context.Context, ownerID string, in CreateProgram) (*repository.Program, error) {
	p := &repository.Program{
		ID:         uuid.NewString(),
		OwnerID:    ownerID,
		Title:      strings.TrimSpace(in.Title),
		Notes:      normalizePtr(in.Notes),
		Visibility: ternaryString(in.Visibility, "private"),
		Version:    1,
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	// guarda versión 1
	ver := &repository.ProgramVersion{
		ID:        uuid.NewString(),
		ProgramID: p.ID,
		Version:   1,
		Title:     p.Title,
		Notes:     p.Notes,
	}
	_ = s.repo.CreateVersion(ctx, ver)
	return p, nil
}

func (s *programService) Get(ctx context.Context, id string) (*repository.Program, error) {
	return s.repo.Get(ctx, id)
}

func (s *programService) Update(ctx context.Context, id string, in UpdateProgram) (*repository.Program, error) {
	patch := map[string]any{}
	if in.Title != nil {
		patch["title"] = strings.TrimSpace(*in.Title)
	}
	if in.Notes != nil {
		patch["notes"] = strings.TrimSpace(*in.Notes)
	}
	if in.Visibility != nil {
		patch["visibility"] = strings.TrimSpace(*in.Visibility)
	}
	return s.repo.Update(ctx, id, patch)
}

func (s *programService) Delete(ctx context.Context, id string) error { return s.repo.Delete(ctx, id) }

func (s *programService) NewVersion(ctx context.Context, programID string) (*repository.ProgramVersion, error) {
	p, err := s.repo.Get(ctx, programID)
	if err != nil {
		return nil, err
	}
	n, err := s.repo.NextVersionNumber(ctx, programID)
	if err != nil {
		return nil, err
	}
	v := &repository.ProgramVersion{
		ID:        uuid.NewString(),
		ProgramID: programID,
		Version:   n,
		Title:     p.Title,
		Notes:     p.Notes,
	}
	if err := s.repo.CreateVersion(ctx, v); err != nil {
		return nil, err
	}
	// opcional: subir Program.version
	_, _ = s.repo.Update(ctx, programID, map[string]any{"version": n})
	return v, nil
}

func (s *programService) ListVersions(ctx context.Context, programID string) ([]repository.ProgramVersion, error) {
	return s.repo.ListVersions(ctx, programID)
}

func (s *programService) ListWeeks(ctx context.Context, programID string) ([]repository.ProgramWeek, error) {
	return s.repo.ListWeeks(ctx, programID)
}

func (s *programService) ListDays(ctx context.Context, weekID string) ([]repository.ProgramDay, error) {
	return s.repo.ListDays(ctx, weekID)
}

func (s *programService) UpdateDay(ctx context.Context, dayID string, notes *string, dayIndex *int) (*repository.ProgramDay, error) {
	patch := map[string]any{}
	if notes != nil {
		patch["notes"] = strings.TrimSpace(*notes)
	}
	if dayIndex != nil {
		patch["day_index"] = *dayIndex
	}
	return s.repo.UpdateDay(ctx, dayID, patch)
}

func (s *programService) DeleteDay(ctx context.Context, dayID string) error {
	return s.repo.DeleteDay(ctx, dayID)
}

func (s *programService) ListPrescriptions(ctx context.Context, dayID string) ([]repository.PrescriptionRow, error) {
	return s.repo.ListPrescriptions(ctx, dayID)
}

func (s *programService) UpdatePrescription(ctx context.Context, id string, in UpdatePrescription) (*repository.Prescription, error) {
	patch := map[string]any{}
	if in.ExerciseID != nil {
		patch["exercise_id"] = *in.ExerciseID
	}
	if in.Series != nil {
		patch["series"] = *in.Series
	}
	if in.Reps != nil {
		patch["reps"] = strings.TrimSpace(*in.Reps)
	}
	if in.RestSec != nil {
		patch["rest_sec"] = *in.RestSec
	}
	if in.ToFailure != nil {
		patch["to_failure"] = *in.ToFailure
	}
	if in.Tempo != nil {
		patch["tempo"] = strings.TrimSpace(*in.Tempo)
	}
	if in.RIR != nil {
		patch["rir"] = *in.RIR
	}
	if in.RPE != nil {
		patch["rpe"] = *in.RPE
	}
	if in.MethodID != nil {
		patch["method_id"] = *in.MethodID
	}
	if in.Notes != nil {
		patch["notes"] = strings.TrimSpace(*in.Notes)
	}
	if in.Position != nil {
		patch["position"] = *in.Position
	}
	return s.repo.UpdatePrescription(ctx, id, patch)
}

func (s *programService) DeletePrescription(ctx context.Context, id string) error {
	return s.repo.DeletePrescription(ctx, id)
}

func (s *programService) ReorderPrescriptions(ctx context.Context, dayID string, orderedIDs []string) error {
	return s.repo.ReorderPrescriptions(ctx, dayID, orderedIDs)
}

func ternaryString(p *string, def string) string {
	if p == nil || strings.TrimSpace(*p) == "" {
		return def
	}
	return strings.TrimSpace(*p)
}

func (s *programService) GetProgram(ctx context.Context, id string) (*repository.ProgramRow, error) {
	return s.repo.GetProgram(ctx, id)
}

func (s *programService) UpdateProgram(ctx context.Context, id string, title *string, notes *string, visibility *string) (*repository.ProgramRow, error) {
	patch := map[string]any{"updated_at": time.Now()}
	if title != nil {
		patch["title"] = *title
	}
	if notes != nil {
		patch["notes"] = *notes
	}
	if visibility != nil {
		patch["visibility"] = *visibility
	}
	if err := s.repo.UpdateProgram(ctx, id, patch); err != nil {
		return nil, err
	}
	return s.repo.GetProgram(ctx, id)
}

func (s *programService) DeleteProgram(ctx context.Context, id string) error {
	return s.repo.DeleteProgram(ctx, id)
}

func (s *programService) CreateNextVersionClone(ctx context.Context, programID string) (*repository.ProgramRow, error) {
	return s.repo.CreateNextVersionClone(ctx, programID)
}

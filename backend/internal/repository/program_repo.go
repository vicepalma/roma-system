package repository

import (
	"context"
	"time"

	"github.com/vicepalma/roma-system/backend/internal/domain"
	"gorm.io/gorm"
)

type ProgramRepository interface {
	CreateProgram(ctx context.Context, p *domain.Program) error
	ListMyPrograms(ctx context.Context, ownerID string, limit, offset int) ([]domain.Program, int64, error)
	AddWeek(ctx context.Context, w *domain.ProgramWeek) error
	AddDay(ctx context.Context, d *domain.ProgramDay) error
	AddPrescription(ctx context.Context, p *domain.Prescription) error
	GetProgramVersion(ctx context.Context, programID string) (int, error)
	Assign(ctx context.Context, a *domain.Assignment) error

	FindActiveAssignmentForDate(ctx context.Context, discipleID string, date time.Time) (*domain.Assignment, error)
	FindDayForDate(ctx context.Context, assignment *domain.Assignment, date time.Time) (*domain.ProgramDay, error)
	ListPrescriptionsByDay(ctx context.Context, dayID string) ([]domain.Prescription, error)
}

type programRepository struct{ db *gorm.DB }

func NewProgramRepository(db *gorm.DB) ProgramRepository { return &programRepository{db: db} }

func (r *programRepository) CreateProgram(ctx context.Context, p *domain.Program) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *programRepository) ListMyPrograms(ctx context.Context, ownerID string, limit, offset int) ([]domain.Program, int64, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	var rows []domain.Program
	tx := r.db.WithContext(ctx).Model(&domain.Program{}).Where("owner_id = ?", ownerID)
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := tx.Order("created_at DESC").Limit(limit).Offset(offset).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *programRepository) AddWeek(ctx context.Context, w *domain.ProgramWeek) error {
	return r.db.WithContext(ctx).Create(w).Error
}
func (r *programRepository) AddDay(ctx context.Context, d *domain.ProgramDay) error {
	return r.db.WithContext(ctx).Create(d).Error
}
func (r *programRepository) AddPrescription(ctx context.Context, p *domain.Prescription) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *programRepository) GetProgramVersion(ctx context.Context, programID string) (int, error) {
	var prog domain.Program
	if err := r.db.WithContext(ctx).First(&prog, "id = ?", programID).Error; err != nil {
		return 0, err
	}
	return prog.Version, nil
}

func (r *programRepository) Assign(ctx context.Context, a *domain.Assignment) error {
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *programRepository) FindActiveAssignmentForDate(ctx context.Context, discipleID string, date time.Time) (*domain.Assignment, error) {
	var asg domain.Assignment
	if err := r.db.WithContext(ctx).
		Where("disciple_id = ? AND is_active = true AND start_date <= ? AND (end_date IS NULL OR end_date >= ?)",
			discipleID, date, date).
		Order("created_at DESC").First(&asg).Error; err != nil {
		return nil, err
	}
	return &asg, nil
}

func (r *programRepository) FindDayForDate(ctx context.Context, asg *domain.Assignment, date time.Time) (*domain.ProgramDay, error) {
	// Día relativo dentro del programa: (0..)
	days := int(date.Sub(asg.StartDate).Hours() / 24)
	if days < 0 {
		days = 0
	}
	// obtenemos semana/día index (simplificado: 7 días por semana)
	weekIndex := days/7 + 1
	dayIndex := days%7 + 1

	var wk domain.ProgramWeek
	if err := r.db.WithContext(ctx).
		First(&wk, "program_id = ? AND week_index = ?", asg.ProgramID, weekIndex).Error; err != nil {
		return nil, err
	}
	var day domain.ProgramDay
	if err := r.db.WithContext(ctx).
		First(&day, "week_id = ? AND day_index = ?", wk.ID, dayIndex).Error; err != nil {
		return nil, err
	}
	return &day, nil
}

func (r *programRepository) ListPrescriptionsByDay(ctx context.Context, dayID string) ([]domain.Prescription, error) {
	var rows []domain.Prescription
	err := r.db.WithContext(ctx).Where("day_id = ?", dayID).Order("position ASC").Find(&rows).Error
	return rows, err
}

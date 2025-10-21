package repository

import (
	"context"
	"errors"
	"time"

	"github.com/vicepalma/roma-system/backend/internal/domain"
	"gorm.io/gorm"
)

type Program struct {
	ID         string `gorm:"type:uuid;primaryKey"`
	OwnerID    string `gorm:"type:uuid;not null"`
	Title      string `gorm:"not null"`
	Notes      *string
	Visibility string `gorm:"not null;default:'private'"`
	Version    int    `gorm:"not null;default:1"`
	CreatedAt  int64  `gorm:"autoCreateTime"`
	UpdatedAt  int64  `gorm:"autoUpdateTime"`
}

type ProgramVersion struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	ProgramID string `gorm:"type:uuid;not null;index"`
	Version   int    `gorm:"not null"`
	Title     string `gorm:"not null"`
	Notes     *string
	CreatedAt int64 `gorm:"autoCreateTime"`
}

type ProgramWeek struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	ProgramID string `gorm:"type:uuid;index;not null"`
	WeekIndex int    `gorm:"not null"`
}

type ProgramDay struct {
	ID       string `gorm:"type:uuid;primaryKey"`
	WeekID   string `gorm:"type:uuid;index;not null"`
	DayIndex int    `gorm:"not null"`
	Notes    *string
}

type Prescription struct {
	ID           string `gorm:"type:uuid;primaryKey"`
	DayID        string `gorm:"type:uuid;index;not null"`
	ExerciseID   string `gorm:"type:uuid;index;not null"`
	Series       int    `gorm:"not null"`
	Reps         string `gorm:"not null"`
	RestSec      *int
	ToFailure    bool
	Tempo        *string
	RIR          *int
	RPE          *float32 `gorm:"type:numeric(3,1)"`
	MethodID     *string  `gorm:"type:uuid"`
	Notes        *string
	Position     int    `gorm:"not null;default:1"`
	ExerciseName string `gorm:"-"` // opc para joins
}

type ProgramFilter struct {
	Query  string
	Limit  int
	Offset int
	Owner  string // opcional
}

type ProgramRow struct {
	ID         string `gorm:"type:uuid;primaryKey"`
	OwnerID    string `gorm:"type:uuid;not null"`
	Title      string `gorm:"not null"`
	Notes      *string
	Visibility string `gorm:"not null;default:private"`
	Version    int    `gorm:"not null;default:1"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type ProgramRepository interface {
	CreateProgram(ctx context.Context, p *domain.Program) error
	ListMyPrograms(ctx context.Context, ownerID string, limit, offset int) ([]domain.Program, int64, error)
	GetProgramVersion(ctx context.Context, programID string) (int, error)
	Assign(ctx context.Context, a *domain.Assignment) error

	FindActiveAssignmentForDate(ctx context.Context, discipleID string, date time.Time) (*domain.Assignment, error)
	FindDayForDate(ctx context.Context, assignment *domain.Assignment, date time.Time) (*domain.ProgramDay, error)
	ListPrescriptionsByDay(ctx context.Context, dayID string) ([]domain.Prescription, error)

	Search(ctx context.Context, f ProgramFilter) ([]Program, int64, error)
	Create(ctx context.Context, p *Program) error
	Get(ctx context.Context, id string) (*Program, error)
	Update(ctx context.Context, id string, patch map[string]any) (*Program, error)
	Delete(ctx context.Context, id string) error

	CreateVersion(ctx context.Context, v *ProgramVersion) error
	ListVersions(ctx context.Context, programID string) ([]ProgramVersion, error)
	NextVersionNumber(ctx context.Context, programID string) (int, error)

	// weeks & days
	AddWeek(ctx context.Context, w *domain.ProgramWeek) error
	ListWeeks(ctx context.Context, programID string) ([]ProgramWeek, error)
	AddDay(ctx context.Context, d *domain.ProgramDay) error
	ListDays(ctx context.Context, weekID string) ([]ProgramDay, error)
	UpdateDay(ctx context.Context, id string, patch map[string]any) (*ProgramDay, error)
	DeleteDay(ctx context.Context, id string) error
	DeleteDaysByWeek(ctx context.Context, weekID string) error
	DeleteWeek(ctx context.Context, programID, weekID string) error

	// prescriptions
	ListPrescriptions(ctx context.Context, dayID string) ([]Prescription, error)
	AddPrescription(ctx context.Context, p *domain.Prescription) error
	UpdatePrescription(ctx context.Context, id string, patch map[string]any) (*Prescription, error)
	DeletePrescription(ctx context.Context, id string) error
	ReorderPrescriptions(ctx context.Context, dayID string, orderedIDs []string) error

	GetProgram(ctx context.Context, id string) (*ProgramRow, error)
	UpdateProgram(ctx context.Context, id string, patch map[string]any) error
	DeleteProgram(ctx context.Context, id string) error
	CreateNextVersionClone(ctx context.Context, programID string) (*ProgramRow, error)
}

type programRepository struct{ db *gorm.DB }

func NewProgramRepository(db *gorm.DB) ProgramRepository { return &programRepository{db: db} }

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
func (r *programRepository) DeleteDaysByWeek(ctx context.Context, weekID string) error {
	return r.db.WithContext(ctx).
		Exec(`DELETE FROM program_days WHERE week_id = ?`, weekID).Error
}

func (r *programRepository) DeleteWeek(ctx context.Context, programID, weekID string) error {
	return r.db.WithContext(ctx).
		Exec(`DELETE FROM program_weeks WHERE id = ? AND program_id = ?`, weekID, programID).Error
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

func (r *programRepository) Search(ctx context.Context, f ProgramFilter) ([]Program, int64, error) {
	q := r.db.WithContext(ctx).Model(&Program{})
	if f.Owner != "" {
		q = q.Where("owner_id = ?", f.Owner)
	}
	if f.Query != "" {
		like := "%" + f.Query + "%"
		q = q.Where("LOWER(title) LIKE LOWER(?)", like)
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
	q = q.Order("created_at DESC")
	var items []Program
	return items, 0, q.Find(&items).Error
}

func (r *programRepository) Create(ctx context.Context, p *Program) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *programRepository) Get(ctx context.Context, id string) (*Program, error) {
	var p Program
	if err := r.db.WithContext(ctx).First(&p, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *programRepository) Update(ctx context.Context, id string, patch map[string]any) (*Program, error) {
	if err := r.db.WithContext(ctx).Model(&Program{}).Where("id = ?", id).Updates(patch).Error; err != nil {
		return nil, err
	}
	return r.Get(ctx, id)
}

func (r *programRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&Program{}, "id = ?", id).Error
}

// versions
func (r *programRepository) CreateVersion(ctx context.Context, v *ProgramVersion) error {
	return r.db.WithContext(ctx).Create(v).Error
}
func (r *programRepository) ListVersions(ctx context.Context, programID string) ([]ProgramVersion, error) {
	var items []ProgramVersion
	return items, r.db.WithContext(ctx).
		Where("program_id = ?", programID).
		Order("version DESC").
		Find(&items).Error
}
func (r *programRepository) NextVersionNumber(ctx context.Context, programID string) (int, error) {
	var max int
	err := r.db.WithContext(ctx).
		Model(&ProgramVersion{}).
		Select("COALESCE(MAX(version), 0)").
		Where("program_id = ?", programID).
		Scan(&max).Error
	return max + 1, err
}

// weeks & days
func (r *programRepository) ListWeeks(ctx context.Context, programID string) ([]ProgramWeek, error) {
	var items []ProgramWeek
	return items, r.db.WithContext(ctx).
		Where("program_id = ?", programID).
		Order("week_index ASC").
		Find(&items).Error
}

func (r *programRepository) ListDays(ctx context.Context, weekID string) ([]ProgramDay, error) {
	var items []ProgramDay
	return items, r.db.WithContext(ctx).
		Where("week_id = ?", weekID).
		Order("day_index ASC, id ASC").
		Find(&items).Error
}
func (r *programRepository) UpdateDay(ctx context.Context, id string, patch map[string]any) (*ProgramDay, error) {
	if err := r.db.WithContext(ctx).Model(&ProgramDay{}).Where("id = ?", id).Updates(patch).Error; err != nil {
		return nil, err
	}
	var d ProgramDay
	if err := r.db.WithContext(ctx).First(&d, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &d, nil
}
func (r *programRepository) DeleteDay(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&ProgramDay{}, "id = ?", id).Error
}

// prescriptions
func (r *programRepository) ListPrescriptions(ctx context.Context, dayID string) ([]Prescription, error) {
	var items []Prescription
	return items, r.db.WithContext(ctx).
		Model(&Prescription{}).
		Where("day_id = ?", dayID).
		Order("position ASC, id ASC").
		Find(&items).Error
}

func (r *programRepository) UpdatePrescription(ctx context.Context, id string, patch map[string]any) (*Prescription, error) {
	if err := r.db.WithContext(ctx).Model(&Prescription{}).Where("id = ?", id).Updates(patch).Error; err != nil {
		return nil, err
	}
	var pr Prescription
	if err := r.db.WithContext(ctx).First(&pr, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &pr, nil
}

func (r *programRepository) DeletePrescription(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&Prescription{}, "id = ?", id).Error
}

func (r *programRepository) ReorderPrescriptions(ctx context.Context, dayID string, orderedIDs []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for idx, id := range orderedIDs {
			if err := tx.Model(&Prescription{}).
				Where("id = ? AND day_id = ?", id, dayID).
				Update("position", idx+1).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetProgram
func (r *programRepository) GetProgram(ctx context.Context, id string) (*ProgramRow, error) {
	var p ProgramRow
	if err := r.db.WithContext(ctx).First(&p, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

// UpdateProgram
func (r *programRepository) UpdateProgram(ctx context.Context, id string, patch map[string]any) error {
	return r.db.WithContext(ctx).Model(&ProgramRow{}).Where("id = ?", id).Updates(patch).Error
}

// DeleteProgram
func (r *programRepository) DeleteProgram(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&ProgramRow{}, "id = ?", id).Error
}

// CreateNextVersionClone — clona programa + semanas + días + prescripciones (versión+1)
func (r *programRepository) CreateNextVersionClone(ctx context.Context, programID string) (*ProgramRow, error) {
	tx := r.db.WithContext(ctx).Begin()

	// 1) programa base
	var base ProgramRow
	if err := tx.Raw(`SELECT id, owner_id, title, notes, visibility, version FROM programs WHERE id = ?`, programID).
		Scan(&base).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if base.ID == "" {
		tx.Rollback()
		return nil, errors.New("program_not_found")
	}

	// 2) crea nuevo programa version+1
	var newProg ProgramRow
	err := tx.Raw(`
		INSERT INTO programs (owner_id, title, notes, visibility, version)
		VALUES (?, ?, ?, ?, ?)
		RETURNING id, owner_id, title, notes, visibility, version, created_at, updated_at
	`, base.OwnerID, base.Title, base.Notes, base.Visibility, base.Version+1).Scan(&newProg).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 3) mapear semanas
	type wkmap struct {
		OldID, NewID string
		WeekIndex    int
	}
	var wks []wkmap
	if err := tx.Raw(`SELECT id AS old_id, week_index FROM program_weeks WHERE program_id = ? ORDER BY week_index, id`, base.ID).
		Scan(&wks).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	for i := range wks {
		var id string
		if err := tx.Raw(`
			INSERT INTO program_weeks (program_id, week_index)
			VALUES (?, ?)
			RETURNING id
		`, newProg.ID, wks[i].WeekIndex).Row().Scan(&id); err != nil {
			tx.Rollback()
			return nil, err
		}
		wks[i].NewID = id
	}

	// 4) mapear días
	type dymap struct {
		OldID, NewID, OldWeekID, NewWeekID string
		DayIndex                           int
		Notes                              *string
	}
	var days []dymap
	if err := tx.Raw(`SELECT d.id AS old_id, d.week_id AS old_week_id, d.day_index, d.notes
	                   FROM program_days d
	                   JOIN program_weeks w ON w.id = d.week_id
	                   WHERE w.program_id = ?
	                   ORDER BY d.day_index, d.id`, base.ID).Scan(&days).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	// asignar week nuevo
	wkIndexByOld := map[string]string{}
	for _, wk := range wks {
		wkIndexByOld[wk.OldID] = wk.NewID
	}
	for i := range days {
		days[i].NewWeekID = wkIndexByOld[days[i].OldWeekID]
		var id string
		if err := tx.Raw(`
			INSERT INTO program_days (week_id, day_index, notes)
			VALUES (?, ?, ?)
			RETURNING id
		`, days[i].NewWeekID, days[i].DayIndex, days[i].Notes).Row().Scan(&id); err != nil {
			tx.Rollback()
			return nil, err
		}
		days[i].NewID = id
	}

	// 5) clonar prescripciones
	dayMap := map[string]string{}
	for _, d := range days {
		dayMap[d.OldID] = d.NewID
	}
	type pres struct {
		ID         string
		DayID      string
		ExerciseID string
		Series     int
		Reps       string
		RestSec    *int
		ToFailure  bool
		Tempo      *string
		Rir        *int
		Rpe        *float64
		MethodID   *string
		Notes      *string
		Position   int
	}
	var presc []pres
	if err := tx.Raw(`SELECT id, day_id, exercise_id, series, reps, rest_sec, to_failure, tempo,
	                          rir, rpe, method_id, notes, position
	                   FROM prescriptions
	                   WHERE day_id IN (SELECT d.id FROM program_days d
	                                    JOIN program_weeks w ON w.id=d.week_id
	                                    WHERE w.program_id = ?)`, base.ID).Scan(&presc).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	for _, p := range presc {
		newDay := dayMap[p.DayID]
		if newDay == "" {
			continue
		}
		if err := tx.Exec(`
			INSERT INTO prescriptions
			(day_id, exercise_id, series, reps, rest_sec, to_failure, tempo, rir, rpe, method_id, notes, position)
			VALUES (?,?,?,?,?,?,?,?,?,?,?,?)
		`, newDay, p.ExerciseID, p.Series, p.Reps, p.RestSec, p.ToFailure, p.Tempo, p.Rir, p.Rpe, p.MethodID, p.Notes, p.Position).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return &newProg, nil
}

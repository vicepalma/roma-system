package repository

import (
	"context"
	"time"

	"github.com/vicepalma/roma-system/backend/internal/domain"
	"gorm.io/gorm"
)

// Datos mínimos para listar discípulos del coach
type DiscipleRow struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type AssignmentMinimal struct {
	ID             string    `json:"id"`
	ProgramID      string    `json:"program_id"`
	ProgramVersion int       `json:"program_version"`
	DiscipleID     string    `json:"disciple_id"`
	AssignedBy     string    `json:"assigned_by"`
	StartDate      time.Time `json:"start_date"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
}

type AssignmentListRow struct {
	ID             string     `json:"id"`
	DiscipleID     string     `json:"disciple_id"`
	DiscipleName   string     `json:"disciple_name"`
	DiscipleEmail  string     `json:"disciple_email"`
	ProgramID      string     `json:"program_id"`
	ProgramTitle   string     `json:"program_title"`
	ProgramVersion int        `json:"program_version"`
	StartDate      time.Time  `json:"start_date"`
	EndDate        *time.Time `json:"end_date,omitempty"`
	IsActive       bool       `json:"is_active"`
	CreatedAt      time.Time  `json:"created_at"`
}

type AssignmentRow struct {
	ID             string    `gorm:"type:uuid;primaryKey"`
	ProgramID      string    `gorm:"type:uuid;not null"`
	ProgramVersion int       `gorm:"not null"`
	DiscipleID     string    `gorm:"type:uuid;not null"`
	AssignedBy     string    `gorm:"type:uuid;not null"`
	StartDate      time.Time `gorm:"type:date;not null"`
	EndDate        *time.Time
	IsActive       bool `gorm:"not null;default:true"`
	CreatedAt      time.Time
}

type ProgramDayLite struct {
	ID       string `gorm:"type:uuid;primaryKey"`
	WeekID   string `gorm:"type:uuid;index"`
	DayIndex int
	Notes    *string
}

type CoachRepository interface {
	CreateLink(ctx context.Context, coachID, discipleID string, autoAccept bool) (*domain.CoachLink, error)
	UpdateStatus(ctx context.Context, id, newStatus string, actorID string) (*domain.CoachLink, error)
	ListLinksForUser(ctx context.Context, userID string) (incoming, outgoing []domain.CoachLink, err error)
	CanCoach(ctx context.Context, coachID, discipleID string) (bool, error)

	ListDisciples(ctx context.Context, coachID string) ([]DiscipleRow, error)
	CreateAssignment(ctx context.Context, coachID, discipleID, programID string, startDate time.Time) (*AssignmentMinimal, error)

	ListAssignmentsForCoach(ctx context.Context, coachID string, discipleID *string, limit, offset int) ([]AssignmentListRow, int64, error)

	GetAssignmentByID(ctx context.Context, id string) (*AssignmentRow, error)
	UpdateAssignment(ctx context.Context, id string, patch map[string]any) error
	ListProgramDaysByProgramWeek(ctx context.Context, programID string, weekIndex int) ([]ProgramDayLite, error)
}

type coachRepository struct{ db *gorm.DB }

func NewCoachRepository(db *gorm.DB) CoachRepository { return &coachRepository{db: db} }

func (r *coachRepository) CreateLink(ctx context.Context, coachID, discipleID string, autoAccept bool) (*domain.CoachLink, error) {
	status := "pending"
	if coachID == discipleID || autoAccept {
		status = "accepted"
	}
	link := &domain.CoachLink{
		CoachID:    coachID,
		DiscipleID: discipleID,
		Status:     status,
	}
	if err := r.db.WithContext(ctx).Create(link).Error; err != nil {
		return nil, err
	}
	return link, nil
}

// Solo el DISCÍPULO puede aceptar/rechazar (regla en service, este método solo guarda)
func (r *coachRepository) UpdateStatus(ctx context.Context, id, newStatus string, _ string) (*domain.CoachLink, error) {
	var link domain.CoachLink
	if err := r.db.WithContext(ctx).First(&link, "id = ?", id).Error; err != nil {
		return nil, err
	}
	link.Status = newStatus
	if err := r.db.WithContext(ctx).Save(&link).Error; err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *coachRepository) ListLinksForUser(ctx context.Context, userID string) (incoming, outgoing []domain.CoachLink, err error) {
	if err = r.db.WithContext(ctx).
		Where("disciple_id = ?", userID).
		Order("created_at DESC").
		Find(&incoming).Error; err != nil {
		return
	}
	if err = r.db.WithContext(ctx).
		Where("coach_id = ?", userID).
		Order("created_at DESC").
		Find(&outgoing).Error; err != nil {
		return
	}
	return
}

func (r *coachRepository) CanCoach(ctx context.Context, coachID, discipleID string) (bool, error) {
	if coachID == discipleID {
		return true, nil
	}
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.CoachLink{}).
		Where("coach_id = ? AND disciple_id = ? AND status = 'accepted'", coachID, discipleID).
		Count(&count).Error
	return count > 0, err
}

func (r *coachRepository) ListDisciples(ctx context.Context, coachID string) ([]DiscipleRow, error) {
	var rows []DiscipleRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT u.id, u.name, u.email
		FROM coach_links cl
		JOIN users u ON u.id = cl.disciple_id
		WHERE cl.coach_id = ? AND cl.status = 'accepted'
		ORDER BY u.name ASC
	`, coachID).Scan(&rows).Error
	return rows, err
}

func (r *coachRepository) CreateAssignment(ctx context.Context, coachID, discipleID, programID string, startDate time.Time) (*AssignmentMinimal, error) {
	// program_version = 1 (simple). Si ya versionas, aquí puedes consultar la actual.
	var row AssignmentMinimal
	err := r.db.WithContext(ctx).Raw(`
		INSERT INTO assignments (program_id, program_version, disciple_id, assigned_by, start_date, is_active)
		VALUES (?, 1, ?, ?, ?, true)
		RETURNING id, program_id, program_version, disciple_id, assigned_by, start_date, is_active, created_at
	`, programID, discipleID, coachID, startDate).Scan(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *coachRepository) ListAssignmentsForCoach(ctx context.Context, coachID string, discipleID *string, limit, offset int) ([]AssignmentListRow, int64, error) {
	// Conjunto de discípulos que el coach puede ver: vínculos aceptados + self-coach
	// Nota: incluimos self-coach con OR a.disciple_id = coachID
	base := `
SELECT 
  a.id, a.disciple_id, u.name AS disciple_name, u.email AS disciple_email,
  a.program_id, p.title AS program_title, a.program_version,
  a.start_date, a.end_date, a.is_active, a.created_at
FROM assignments a
JOIN users u ON u.id = a.disciple_id
JOIN programs p ON p.id = a.program_id
WHERE (
  a.disciple_id IN (
    SELECT cl.disciple_id 
    FROM coach_links cl
    WHERE cl.coach_id = ? AND cl.status = 'accepted'
  )
  OR a.disciple_id = ?
)
`

	args := []interface{}{coachID, coachID}

	if discipleID != nil && *discipleID != "" {
		base += ` AND a.disciple_id = ?`
		args = append(args, *discipleID)
	}

	baseOrderLimit := ` ORDER BY a.created_at DESC`
	if limit > 0 {
		baseOrderLimit += ` LIMIT ? OFFSET ?`
		args = append(args, limit, offset)
	}

	var rows []AssignmentListRow
	if err := r.db.WithContext(ctx).Raw(base+baseOrderLimit, args...).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	// total (sin limit/offset)
	countSQL := `
SELECT COUNT(*)
FROM assignments a
WHERE (
  a.disciple_id IN (
    SELECT cl.disciple_id 
    FROM coach_links cl
    WHERE cl.coach_id = ? AND cl.status = 'accepted'
  )
  OR a.disciple_id = ?
)`
	countArgs := []interface{}{coachID, coachID}
	if discipleID != nil && *discipleID != "" {
		countSQL += ` AND a.disciple_id = ?`
		countArgs = append(countArgs, *discipleID)
	}

	var total int64
	if err := r.db.WithContext(ctx).Raw(countSQL, countArgs...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

// GetAssignmentByID
func (r *coachRepository) GetAssignmentByID(ctx context.Context, id string) (*AssignmentRow, error) {
	var a AssignmentRow
	if err := r.db.WithContext(ctx).
		First(&a, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

// UpdateAssignment
func (r *coachRepository) UpdateAssignment(ctx context.Context, id string, patch map[string]any) error {
	return r.db.WithContext(ctx).
		Model(&AssignmentRow{}).
		Where("id = ?", id).
		Updates(patch).Error
}

// ListProgramDaysByProgramWeek (week_index=1 para el MVP)
func (r *coachRepository) ListProgramDaysByProgramWeek(ctx context.Context, programID string, weekIndex int) ([]ProgramDayLite, error) {
	type weekRow struct {
		ID string
	}
	var w weekRow
	if err := r.db.WithContext(ctx).
		Raw(`SELECT id FROM program_weeks WHERE program_id = ? AND week_index = ? LIMIT 1`, programID, weekIndex).
		Scan(&w).Error; err != nil {
		return nil, err
	}
	var days []ProgramDayLite
	if err := r.db.WithContext(ctx).
		Raw(`SELECT id, week_id, day_index, notes
		     FROM program_days WHERE week_id = ? ORDER BY day_index ASC, id ASC`, w.ID).
		Scan(&days).Error; err != nil {
		return nil, err
	}
	return days, nil
}

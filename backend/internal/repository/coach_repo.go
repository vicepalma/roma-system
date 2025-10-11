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

type CoachRepository interface {
	CreateLink(ctx context.Context, coachID, discipleID string, autoAccept bool) (*domain.CoachLink, error)
	UpdateStatus(ctx context.Context, id, newStatus string, actorID string) (*domain.CoachLink, error)
	ListLinksForUser(ctx context.Context, userID string) (incoming, outgoing []domain.CoachLink, err error)
	CanCoach(ctx context.Context, coachID, discipleID string) (bool, error)

	ListDisciples(ctx context.Context, coachID string) ([]DiscipleRow, error)
	CreateAssignment(ctx context.Context, coachID, discipleID, programID string, startDate time.Time) (*AssignmentMinimal, error)
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

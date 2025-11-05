package repository

import (
	"context"
	"errors"
	"time"

	"github.com/vicepalma/roma-system/backend/internal/domain"
	"gorm.io/gorm"
)

type SessionMeta struct {
	ID           string     `json:"id"`
	AssignmentID string     `json:"assignment_id"`
	DiscipleID   string     `json:"disciple_id"`
	DayID        string     `json:"day_id"`
	PerformedAt  time.Time  `json:"performed_at"`
	Status       string     `json:"status"`
	EndedAt      *time.Time `json:"ended_at,omitempty"`
	Notes        *string    `json:"notes,omitempty"`
}

type SessionRepository interface {
	CreateSession(ctx context.Context, s *domain.SessionLog) error
	GetSession(ctx context.Context, id, discipleID string) (*domain.SessionLog, error)
	ListSets(ctx context.Context, sessionID string) ([]domain.SetRow, error)
	AddSet(ctx context.Context, set *domain.SetLog) error
	AddCardio(ctx context.Context, seg *CardioSegment) error
	ListCardio(ctx context.Context, sessionID string) ([]CardioSegment, error)

	GetSessionByID(ctx context.Context, id string) (*domain.SessionLog, error)
	ListSessionSets(ctx context.Context, sessionID string, prescriptionID *string, limit, offset int) ([]domain.SetLog, int64, error)

	GetSessionMeta(ctx context.Context, id string) (*SessionMeta, error)
	UpdateSession(ctx context.Context, id string, patch map[string]any) error
	UpdateSet(ctx context.Context, setID string, patch map[string]any) error
	DeleteSet(ctx context.Context, setID string) error
	GetLatestOpenByDisciple(ctx context.Context, discipleID string) (*domain.SessionLog, error)
}

type sessionRepository struct{ db *gorm.DB }

func NewSessionRepository(db *gorm.DB) SessionRepository { return &sessionRepository{db: db} }

func (r *sessionRepository) CreateSession(ctx context.Context, s *domain.SessionLog) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *sessionRepository) GetSession(ctx context.Context, id, discipleID string) (*domain.SessionLog, error) {
	var s domain.SessionLog
	if err := r.db.WithContext(ctx).
		First(&s, "id = ? AND disciple_id = ?", id, discipleID).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *sessionRepository) ListSets(ctx context.Context, sessionID string) ([]domain.SetRow, error) {
	var rows []domain.SetRow
	err := r.db.WithContext(ctx).
		Table("set_logs AS s").
		Select(`
			s.id,
			s.session_id,
			s.prescription_id,
			s.set_index,
			s.weight,
			s.reps,
			s.rpe,
			s.to_failure,
			p.day_id,
			p.exercise_id,
			COALESCE(e.name, '') AS exercise_name
		`).
		Joins(`JOIN prescriptions AS p ON p.id = s.prescription_id`).
		Joins(`LEFT JOIN exercises AS e ON e.id = p.exercise_id`).
		Where("s.session_id = ?", sessionID).
		Order("s.set_index ASC, s.id ASC").
		Scan(&rows).Error
	return rows, err
}

func (r *sessionRepository) AddSet(ctx context.Context, set *domain.SetLog) error {
	return r.db.WithContext(ctx).Create(set).Error
}

/* -------- Cardio (session) -------- */
type CardioSegment struct {
	ID          string  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	SessionID   string  `gorm:"type:uuid;not null;index" json:"session_id"`
	Modality    string  `gorm:"type:text;not null" json:"modality"`
	Minutes     int     `gorm:"not null" json:"minutes"`
	TargetHRMin *int    `json:"target_hr_min,omitempty"`
	TargetHRMax *int    `json:"target_hr_max,omitempty"`
	Notes       *string `json:"notes,omitempty"`
}

func (CardioSegment) TableName() string { return "cardio_segments" }

func (r *sessionRepository) AddCardio(ctx context.Context, seg *CardioSegment) error {
	return r.db.WithContext(ctx).Create(seg).Error
}
func (r *sessionRepository) ListCardio(ctx context.Context, sessionID string) ([]CardioSegment, error) {
	var rows []CardioSegment
	err := r.db.WithContext(ctx).Where("session_id = ?", sessionID).Order("id ASC").Find(&rows).Error
	return rows, err
}

func (r *sessionRepository) GetSessionByID(ctx context.Context, id string) (*domain.SessionLog, error) {
	var s domain.SessionLog
	err := r.db.WithContext(ctx).First(&s, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *sessionRepository) ListSessionSets(ctx context.Context, sessionID string, prescriptionID *string, limit, offset int) ([]domain.SetLog, int64, error) {
	q := r.db.WithContext(ctx).Model(&domain.SetLog{}).Where("session_id = ?", sessionID)
	if prescriptionID != nil && *prescriptionID != "" {
		q = q.Where("prescription_id = ?", *prescriptionID)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	itemsQ := q.Order("set_index ASC").Order("id ASC")
	if limit > 0 {
		itemsQ = itemsQ.Limit(limit).Offset(offset)
	}

	var items []domain.SetLog
	if err := itemsQ.Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *sessionRepository) UpdateSet(ctx context.Context, setID string, patch map[string]any) error {
	return r.db.WithContext(ctx).Table("set_logs").Where("id = ?", setID).Updates(patch).Error
}

func (r *sessionRepository) DeleteSet(ctx context.Context, setID string) error {
	return r.db.WithContext(ctx).Exec(`DELETE FROM set_logs WHERE id = ?`, setID).Error
}

func (r *sessionRepository) UpdateSession(ctx context.Context, id string, patch map[string]any) error {
	return r.db.WithContext(ctx).Table("session_logs").Where("id = ?", id).Updates(patch).Error
}

func (r *sessionRepository) GetSessionMeta(ctx context.Context, id string) (*SessionMeta, error) {
	var out SessionMeta
	err := r.db.WithContext(ctx).
		Raw(`SELECT id, assignment_id, disciple_id, day_id, performed_at, status, ended_at, notes
		     FROM session_logs WHERE id = ?`, id).
		Scan(&out).Error
	if err != nil {
		return nil, err
	}
	if out.ID == "" {
		return nil, gorm.ErrRecordNotFound
	}
	return &out, nil
}

func (r *sessionRepository) GetLatestOpenByDisciple(ctx context.Context, discipleID string) (*domain.SessionLog, error) {
	var s domain.SessionLog
	q := `
	SELECT id, assignment_id, day_id, disciple_id, performed_at, notes, created_at, updated_at, status, ended_at
	FROM session_logs
	WHERE disciple_id = ?
	  AND status = 'open'
	ORDER BY performed_at DESC, created_at DESC
	LIMIT 1;
	`
	if err := r.db.WithContext(ctx).Raw(q, discipleID).Scan(&s).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	if s.ID == "" {
		return nil, nil
	}
	return &s, nil
}

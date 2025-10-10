package repository

import (
	"context"
	"time"

	"github.com/vicepalma/roma-system/backend/internal/domain"
	"gorm.io/gorm"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, s *domain.SessionLog) error
	GetSession(ctx context.Context, id, discipleID string) (*domain.SessionLog, error)
	ListSets(ctx context.Context, sessionID string) ([]domain.SetLog, error)
	AddSet(ctx context.Context, set *domain.SetLog) error
	UpdateSession(ctx context.Context, id string, performedAt *time.Time, notes *string) error
	AddCardio(ctx context.Context, seg *CardioSegment) error
	ListCardio(ctx context.Context, sessionID string) ([]CardioSegment, error)
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

func (r *sessionRepository) ListSets(ctx context.Context, sessionID string) ([]domain.SetLog, error) {
	var rows []domain.SetLog
	err := r.db.WithContext(ctx).Where("session_id = ?", sessionID).
		Order("set_index ASC, id ASC").Find(&rows).Error
	return rows, err
}

func (r *sessionRepository) AddSet(ctx context.Context, set *domain.SetLog) error {
	return r.db.WithContext(ctx).Create(set).Error
}

func (r *sessionRepository) UpdateSession(ctx context.Context, id string, performedAt *time.Time, notes *string) error {
	up := map[string]interface{}{}
	if performedAt != nil {
		up["performed_at"] = *performedAt
	}
	if notes != nil {
		up["notes"] = notes
	}
	return r.db.WithContext(ctx).Model(&domain.SessionLog{}).Where("id = ?", id).Updates(up).Error
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

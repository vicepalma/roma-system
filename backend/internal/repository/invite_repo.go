package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Invitation struct {
	ID         string `gorm:"type:uuid;primaryKey"`
	Code       string `gorm:"uniqueIndex"`
	CoachID    string `gorm:"type:uuid;not null"`
	Email      string `gorm:"not null"`
	Name       *string
	Status     string    `gorm:"not null;default:pending"`
	ExpiresAt  time.Time `gorm:"not null"`
	AcceptedBy *string
	AcceptedAt *time.Time
	CreatedAt  time.Time `gorm:"not null;default:now()"`
}

func (Invitation) TableName() string { return "invite_codes" }

type InviteRepository interface {
	Create(ctx context.Context, inv *Invitation) error
	FindByCode(ctx context.Context, code string) (*Invitation, error)
	MarkAccepted(ctx context.Context, id string, userID string, at time.Time) error
	MarkRevoked(ctx context.Context, id string) error
}

type inviteRepository struct{ db *gorm.DB }

func NewInviteRepository(db *gorm.DB) InviteRepository { return &inviteRepository{db: db} }

func (r *inviteRepository) Create(ctx context.Context, inv *Invitation) error {
	return r.db.WithContext(ctx).Create(inv).Error
}

func (r *inviteRepository) FindByCode(ctx context.Context, code string) (*Invitation, error) {
	var inv Invitation
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&inv).Error; err != nil {
		return nil, err
	}
	return &inv, nil
}

func (r *inviteRepository) MarkAccepted(ctx context.Context, id string, userID string, at time.Time) error {
	return r.db.WithContext(ctx).
		Model(&Invitation{}).
		Where("id = ? AND status = 'pending'", id).
		Updates(map[string]any{
			"status":      "accepted",
			"accepted_by": userID,
			"accepted_at": at,
		}).Error
}

func (r *inviteRepository) MarkRevoked(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Model(&Invitation{}).
		Where("id = ? AND status = 'pending'", id).
		Update("status", "revoked").Error
}

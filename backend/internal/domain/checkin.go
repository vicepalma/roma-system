package domain

import "time"

type Checkin struct {
	ID         string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	DiscipleID string    `gorm:"type:uuid;not null;index" json:"disciple_id"`
	CheckedAt  time.Time `gorm:"type:date;not null;default:CURRENT_DATE" json:"checked_at"`
	CreatedAt  time.Time `gorm:"not null;default:now()" json:"created_at"`
	WeightKG   *float64  `gorm:"column:weight_kg" json:"weight_kg,omitempty"`
	Notes      *string   `gorm:"type:text" json:"notes,omitempty"`
}

func (Checkin) TableName() string { return "checkins" }

package domain

import (
	"time"

	"github.com/lib/pq"
)

type Exercise struct {
	ID            string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name          string         `gorm:"type:text;not null" json:"name"`
	PrimaryMuscle string         `gorm:"type:text;not null" json:"primary_muscle"`
	Equipment     *string        `gorm:"type:text" json:"equipment,omitempty"`
	Tags          pq.StringArray `gorm:"type:text[]" json:"tags"`
	Notes         *string        `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Exercise) TableName() string { return "exercises" }

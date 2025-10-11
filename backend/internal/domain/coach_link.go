package domain

import "time"

type CoachLink struct {
	ID         string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CoachID    string    `gorm:"type:uuid;not null" json:"coach_id"`
	DiscipleID string    `gorm:"type:uuid;not null" json:"disciple_id"`
	Status     string    `gorm:"type:text;default:pending" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

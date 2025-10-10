package domain

import "time"

type SessionLog struct {
	ID           string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	AssignmentID string    `gorm:"type:uuid;not null;index" json:"assignment_id"`
	DiscipleID   string    `gorm:"type:uuid;not null;index" json:"disciple_id"`
	DayID        string    `gorm:"type:uuid;not null;index" json:"day_id"`
	PerformedAt  time.Time `gorm:"not null;default:now()" json:"performed_at"`
	Notes        *string   `gorm:"type:text" json:"notes,omitempty"`
}

func (SessionLog) TableName() string { return "session_logs" }

type SetLog struct {
	ID             string   `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	SessionID      string   `gorm:"type:uuid;not null;index" json:"session_id"`
	PrescriptionID string   `gorm:"type:uuid;not null;index" json:"prescription_id"`
	SetIndex       int      `gorm:"not null" json:"set_index"`
	Weight         *float64 `json:"weight,omitempty"`
	Reps           int      `gorm:"not null" json:"reps"`
	RPE            *float32 `json:"rpe,omitempty"`
	ToFailure      bool     `gorm:"not null;default:false" json:"to_failure"`
}

func (SetLog) TableName() string { return "set_logs" }

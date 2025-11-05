package domain

import "time"

type SessionLog struct {
	ID           string     `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	AssignmentID string     `gorm:"type:uuid;not null;index" json:"assignment_id"`
	DayID        string     `gorm:"type:uuid;not null;index" json:"day_id"`
	DiscipleID   string     `gorm:"type:uuid;not null;index" json:"disciple_id"`
	PerformedAt  time.Time  `gorm:"not null;default:now()" json:"performed_at"`
	Notes        *string    `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt    time.Time  `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"not null;default:now()" json:"updated_at"`
	Status       string     `gorm:"type:varchar(20);not null;default:'open'" json:"status"`
	EndedAt      *time.Time `json:"ended_at,omitempty"`
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

type SetRow struct {
	ID             string   `json:"id"`
	SessionID      string   `json:"session_id"`
	PrescriptionID string   `json:"prescription_id"`
	SetIndex       int      `json:"set_index"`
	Weight         *float64 `json:"weight,omitempty"`
	Reps           int      `json:"reps"`
	RPE            *float32 `json:"rpe,omitempty"`
	ToFailure      bool     `json:"to_failure"`

	// Campos enriquecidos desde prescription/exercise
	DayID        string `json:"day_id"`
	ExerciseID   string `json:"exercise_id"`
	ExerciseName string `json:"exercise_name"`
}

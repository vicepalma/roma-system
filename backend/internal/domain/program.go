package domain

import "time"

type Program struct {
	ID         string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	OwnerID    string    `gorm:"type:uuid;not null;index" json:"owner_id"`
	Title      string    `gorm:"type:text;not null" json:"title"`
	Notes      *string   `gorm:"type:text" json:"notes,omitempty"`
	Visibility string    `gorm:"type:text;not null;default:'private'" json:"visibility"`
	Version    int       `gorm:"not null;default:1" json:"version"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Program) TableName() string { return "programs" }

type ProgramWeek struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ProgramID string `gorm:"type:uuid;not null;index" json:"program_id"`
	WeekIndex int    `gorm:"not null" json:"week_index"`
}

func (ProgramWeek) TableName() string { return "program_weeks" }

type ProgramDay struct {
	ID       string  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	WeekID   string  `gorm:"type:uuid;not null;index" json:"week_id"`
	DayIndex int     `gorm:"not null" json:"day_index"`
	Title    *string `gorm:"type:text" json:"title,omitempty"`
	Notes    *string `gorm:"type:text" json:"notes,omitempty"`
}

func (ProgramDay) TableName() string { return "program_days" }

type Prescription struct {
	ID         string   `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	DayID      string   `gorm:"type:uuid;not null;index" json:"day_id"`
	ExerciseID string   `gorm:"type:uuid;not null;index" json:"exercise_id"`
	Series     int      `gorm:"not null" json:"series"`
	Reps       string   `gorm:"type:text;not null" json:"reps"`
	RestSec    *int     `json:"rest_sec,omitempty"`
	ToFailure  bool     `gorm:"not null;default:false" json:"to_failure"`
	Tempo      *string  `json:"tempo,omitempty"`
	RIR        *int     `json:"rir,omitempty"`
	RPE        *float32 `json:"rpe,omitempty"`
	MethodID   *string  `json:"method_id,omitempty"`
	Notes      *string  `json:"notes,omitempty"`
	Position   int      `gorm:"not null;default:1" json:"position"`
}

func (Prescription) TableName() string { return "prescriptions" }

type Assignment struct {
	ID             string     `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ProgramID      string     `gorm:"type:uuid;not null;index" json:"program_id"`
	ProgramVersion int        `gorm:"not null" json:"program_version"`
	DiscipleID     string     `gorm:"type:uuid;not null;index" json:"disciple_id"`
	AssignedBy     string     `gorm:"type:uuid;not null" json:"assigned_by"`
	StartDate      time.Time  `gorm:"type:date;not null" json:"start_date"`
	EndDate        *time.Time `gorm:"type:date" json:"end_date,omitempty"`
	IsActive       bool       `gorm:"not null;default:true" json:"is_active"`
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

func (Assignment) TableName() string { return "assignments" }
